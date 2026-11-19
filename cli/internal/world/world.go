// SPDX-License-Identifier: Apache-2.0

// Package world answers questions about loaded definitions that several
// stages ask: which kinds a kind is, which fields it has, which links a link
// sits under, and how a field's schema reads.
package world

import (
	"cmp"
	"slices"

	"github.com/openfjord/of-lang/cli/internal/model"
)

// Ancestors returns a kind and every kind it extends, nearest first. A cycle
// stops the walk; the checker reports it.
func Ancestors(w model.World, k model.QName) []model.QName {
	var out []model.QName
	for cur := &k; cur != nil; {
		if slices.Contains(out, *cur) {
			break
		}
		out = append(out, *cur)
		kind, ok := w.Kinds[*cur]
		if !ok {
			break
		}
		cur = kind.Extends
	}
	return out
}

// Unions returns every union kind, in a stable order.
func Unions(w model.World) []model.Kind {
	var out []model.Kind
	for _, k := range w.Kinds {
		if len(k.OneOf) > 0 {
			out = append(out, k)
		}
	}
	slices.SortFunc(out, func(a, b model.Kind) int {
		return cmp.Or(cmp.Compare(a.Q.Pkg, b.Q.Pkg), cmp.Compare(a.Q.Name, b.Q.Name))
	})
	return out
}

// Supers returns every kind that a thing of kind k is: k, the kinds it
// extends, and the unions that list any of them.
func Supers(w model.World, k model.QName) []model.QName {
	out := Ancestors(w, k)
	unions := Unions(w)
	for grew := true; grew; {
		grew = false
		for _, u := range unions {
			if slices.Contains(out, u.Q) {
				continue
			}
			if slices.ContainsFunc(u.OneOf, func(l model.QName) bool { return slices.Contains(out, l) }) {
				out = append(out, u.Q)
				grew = true
			}
		}
	}
	return out
}

// IsA says whether a thing of kind k is also a thing of kind want: through
// extends, or through a union that lists it.
func IsA(w model.World, k, want model.QName) bool {
	return slices.Contains(Supers(w, k), want)
}

// Compatible says whether two kinds can name the same thing: one is the
// other, or extends it.
func Compatible(w model.World, a, b model.QName) bool {
	return IsA(w, a, b) || IsA(w, b, a)
}

// IsUnion says whether a kind is a union.
func IsUnion(w model.World, k model.QName) bool {
	return len(w.Kinds[k].OneOf) > 0
}

// Fields returns a kind's fields with those it inherits, the furthest
// ancestor's first.
func Fields(w model.World, k model.QName) []model.Field {
	anc := Ancestors(w, k)
	var out []model.Field
	for i := len(anc) - 1; i >= 0; i-- {
		out = append(out, w.Kinds[anc[i]].Fields...)
	}
	return out
}

// Field finds one field of a kind, inherited or its own, and the kind that
// declares it.
func Field(w model.World, k model.QName, name string) (model.Field, model.QName, bool) {
	for _, a := range Ancestors(w, k) {
		for _, f := range w.Kinds[a].Fields {
			if f.Name == name {
				return f, a, true
			}
		}
	}
	return model.Field{}, model.QName{}, false
}

// SuperLinks returns the links a link sits under, nearest first, not
// including the link itself.
func SuperLinks(w model.World, l model.QName) []model.QName {
	var out []model.QName
	cur := w.Links[l].Under
	for cur != nil && *cur != l && !slices.Contains(out, *cur) {
		out = append(out, *cur)
		cur = w.Links[*cur].Under
	}
	return out
}

// LinkEnds returns the kinds a link item joins, read in the direction its name
// goes: an inverse name swaps them.
func LinkEnds(w model.World, l model.QName, inverted bool) (from, to model.QName) {
	link := w.Links[l]
	if inverted {
		return link.To, link.From
	}
	return link.From, link.To
}

// Resolved returns a schema with the value type it names folded in: the value
// type's keywords, then the field's own on top. A field only narrows its value
// type, so the stricter field keywords win.
func Resolved(w model.World, s model.Schema) model.Schema {
	if !s.HasRef {
		return s
	}
	vt, ok := w.ValueTypes[s.Ref]
	if !ok {
		return s
	}
	out := Resolved(w, vt.Schema)
	out.Pos = s.Pos
	if s.Enum != nil {
		out.Enum = s.Enum
	}
	if s.Const != nil {
		out.Const = s.Const
	}
	if s.Default != nil {
		out.Default = s.Default
	}
	if s.Pattern != "" {
		out.Pattern = s.Pattern
	}
	if s.MinLength != nil {
		out.MinLength = s.MinLength
	}
	if s.MaxLength != nil {
		out.MaxLength = s.MaxLength
	}
	if s.Minimum != nil {
		out.Minimum = s.Minimum
	}
	if s.Maximum != nil {
		out.Maximum = s.Maximum
	}
	if s.Format != "" {
		out.Format = s.Format
	}
	if s.Items != nil {
		out.Items = s.Items
	}
	if s.Properties != nil {
		out.Properties = s.Properties
	}
	return out
}

// Type is a field's type in the sense of VT5: a value type, or a built-in
// type for a field that names none. Two fields compare only when their types
// are equal.
type Type struct {
	ValueType model.QName
	Builtin   string
}

// TypeOf returns the type of a schema. A schema with only enum or const takes
// the type of its values; integer and number compare with each other. ok is
// false when nothing says the type.
func TypeOf(s model.Schema) (Type, bool) {
	if s.HasRef {
		return Type{ValueType: s.Ref}, true
	}
	if s.Type != "" {
		if s.Type == "integer" {
			return Type{Builtin: "number"}, true
		}
		return Type{Builtin: s.Type}, true
	}
	vals := slices.Clone(s.Enum)
	if s.Const != nil {
		vals = append(vals, *s.Const)
	}
	if len(vals) == 0 {
		return Type{}, false
	}
	k := vals[0].Kind
	if k == model.ValueInteger {
		k = model.ValueNumber
	}
	return Type{Builtin: string(k)}, true
}

// TypeText renders a type for messages.
func TypeText(t Type) string {
	if t.Builtin != "" {
		return t.Builtin
	}
	return t.ValueType.Name
}

// Part finds the schema of a part of an object field, following a path of
// part names from the outer part to the inner one. For a oneOf, a part that
// one or more choices define counts. The schema it returns keeps the value
// type it names, so its type can be compared.
func Part(w model.World, s model.Schema, path []string) (model.Schema, bool) {
	cur := s
	for _, name := range path {
		next, ok := partOf(w, Resolved(w, cur), name)
		if !ok {
			return model.Schema{}, false
		}
		cur = next
	}
	return cur, true
}

func partOf(w model.World, s model.Schema, name string) (model.Schema, bool) {
	all := partsOf(w, s, name)
	if len(all) == 0 {
		return model.Schema{}, false
	}
	return all[0], true
}

// Parts finds each schema that a part of an object field can have, following
// a path of part names like Part. A part that several choices of a oneOf
// define has one schema for each of those choices.
func Parts(w model.World, s model.Schema, path []string) []model.Schema {
	cur := []model.Schema{s}
	for _, name := range path {
		var next []model.Schema
		for _, c := range cur {
			next = append(next, partsOf(w, Resolved(w, c), name)...)
		}
		cur = next
	}
	return cur
}

func partsOf(w model.World, s model.Schema, name string) []model.Schema {
	var out []model.Schema
	for _, p := range s.Properties {
		if p.Name == name {
			out = append(out, p.Schema)
		}
	}
	for _, c := range s.OneOf {
		out = append(out, partsOf(w, Resolved(w, c), name)...)
	}
	return out
}

// Reads renders a link's reads sentence for two things.
func Reads(w model.World, l model.QName, from, to string) string {
	link, ok := w.Links[l]
	if !ok || link.Reads == "" {
		return from + " " + l.Name + " " + to
	}
	out := []byte{}
	s := link.Reads
	for i := 0; i < len(s); i++ {
		switch {
		case hasAt(s, i, "{from}"):
			out = append(out, from...)
			i += len("{from}") - 1
		case hasAt(s, i, "{to}"):
			out = append(out, to...)
			i += len("{to}") - 1
		default:
			out = append(out, s[i])
		}
	}
	return string(out)
}

func hasAt(s string, i int, sub string) bool {
	return len(s)-i >= len(sub) && s[i:i+len(sub)] == sub
}
