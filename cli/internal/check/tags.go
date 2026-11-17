// SPDX-License-Identifier: Apache-2.0

package check

import (
	"slices"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
)

// tagged is a definition of the root package with its tags and the
// definitions it refers to.
type tagged struct {
	name string
	pos  diag.Pos
	tags []model.Tag
	refs []model.QName
}

// checkTags reports a tag the manifest does not declare (TG3), and a tagged
// definition that refers to a definition without the same tag (TG5).
func checkTags(w model.World) []diag.Diagnostic {
	root := w.Packages[w.Root]
	var defs []tagged
	for _, k := range sorted(w.Kinds, func(k model.Kind) diag.Pos { return k.Pos }) {
		if k.Q.Pkg != w.Root {
			continue
		}
		refs := slices.Clone(k.OneOf)
		if k.Extends != nil {
			refs = append(refs, *k.Extends)
		}
		for _, f := range k.Fields {
			refs = append(refs, schemaRefs(f.Schema)...)
		}
		defs = append(defs, tagged{name: k.Q.Name, pos: k.Pos, tags: k.Tags, refs: refs})
	}
	for _, l := range sorted(w.Links, func(l model.Link) diag.Pos { return l.Pos }) {
		if l.Q.Pkg != w.Root {
			continue
		}
		refs := []model.QName{l.From, l.To}
		if l.Under != nil {
			refs = append(refs, *l.Under)
		}
		defs = append(defs, tagged{name: l.Q.Name, pos: l.Pos, tags: l.Tags, refs: refs})
	}
	for _, vt := range sorted(w.ValueTypes, func(v model.ValueType) diag.Pos { return v.Pos }) {
		if vt.Q.Pkg == w.Root {
			defs = append(defs, tagged{name: vt.Q.Name, pos: vt.Pos, tags: vt.Tags, refs: schemaRefs(vt.Schema)})
		}
	}
	for _, q := range w.Order {
		if q.Pkg != w.Root {
			continue
		}
		st := w.Statements[q]
		refs := blockRefs(st.Match)
		refs = append(refs, conditionRefs(st.Unless)...)
		for _, d := range st.Derive {
			switch d.Type {
			case model.DeriveLink:
				refs = append(refs, d.Link)
			case model.DeriveKind:
				refs = append(refs, d.Kind)
			default:
				// A field belongs to its kind, which the match names.
			}
		}
		defs = append(defs, tagged{name: q.Name, pos: st.Pos, tags: st.Tags, refs: refs})
	}
	var ds []diag.Diagnostic
	for _, d := range defs {
		for _, t := range d.tags {
			if !slices.ContainsFunc(root.Manifest.Tags, func(td model.TagDecl) bool { return td.Name == t.Name }) {
				ds = append(ds, diag.New(t.Pos, "the tag %s is not in the tags of the manifest (TG3)", t.Name))
				continue
			}
			for _, r := range uniq(d.refs) {
				if !slices.ContainsFunc(tagsOf(w, r), func(rt model.Tag) bool { return rt.Name == t.Name }) {
					ds = append(ds, diag.New(t.Pos, "%s has the tag %s and refers to %s, which does not have it (TG5)", d.name, t.Name, r.Name))
				}
			}
		}
	}
	return ds
}

func uniq(qs []model.QName) []model.QName {
	var out []model.QName
	for _, q := range qs {
		if q.Name != "" && !slices.Contains(out, q) {
			out = append(out, q)
		}
	}
	return out
}

// tagsOf returns the tags of any definition.
func tagsOf(w model.World, q model.QName) []model.Tag {
	if k, ok := w.Kinds[q]; ok {
		return k.Tags
	}
	if l, ok := w.Links[q]; ok {
		return l.Tags
	}
	if vt, ok := w.ValueTypes[q]; ok {
		return vt.Tags
	}
	return w.Statements[q].Tags
}

func schemaRefs(s model.Schema) []model.QName {
	var out []model.QName
	if s.HasRef {
		out = append(out, s.Ref)
	}
	if s.Items != nil {
		out = append(out, schemaRefs(*s.Items)...)
	}
	for _, p := range s.Properties {
		out = append(out, schemaRefs(p.Schema)...)
	}
	for _, c := range s.OneOf {
		out = append(out, schemaRefs(c)...)
	}
	return out
}

func blockRefs(b model.Block) []model.QName {
	var out []model.QName
	for _, ki := range b.Kinds {
		out = append(out, ki.Kind)
	}
	for _, li := range b.Links {
		out = append(out, li.Link)
		for _, e := range []model.End{li.From, li.To} {
			if e.Var == "" {
				out = append(out, e.Kind)
			}
		}
	}
	return append(out, conditionRefs(b.Filter)...)
}

func conditionRefs(conds []model.Condition) []model.QName {
	var out []model.QName
	for _, c := range conds {
		switch c.Type {
		case model.CondIsA:
			out = append(out, c.Kind)
		case model.CondAtLeastOne, model.CondNone, model.CondCount:
			out = append(out, blockRefs(*c.Block)...)
		case model.CondAnyOf:
			out = append(out, conditionRefs(c.Any)...)
		default:
			// A comparison names fields of kinds the match already names.
		}
	}
	return out
}
