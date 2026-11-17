// SPDX-License-Identifier: Apache-2.0

// Package world answers questions about loaded definitions that several
// stages ask: which kinds a kind extends, which fields and link counts it
// carries, and which links a link sits under.
package world

import (
	"slices"

	"oflang/internal/model"
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

// IsA says whether a thing of kind k is also a thing of kind want.
func IsA(w model.World, k, want model.QName) bool {
	return slices.Contains(Ancestors(w, k), want)
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

// Field finds one field of a kind, inherited or its own.
func Field(w model.World, k model.QName, name string) (model.Field, bool) {
	for _, f := range Fields(w, k) {
		if f.Name == name {
			return f, true
		}
	}
	return model.Field{}, false
}

// LinkCounts returns a kind's link counts with those it inherits.
func LinkCounts(w model.World, k model.QName) []model.LinkCount {
	var out []model.LinkCount
	for _, a := range Ancestors(w, k) {
		out = append(out, w.Kinds[a].Links...)
	}
	return out
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
