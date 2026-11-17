// SPDX-License-Identifier: Apache-2.0

package engine

import (
	"slices"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/world"
)

type read struct {
	pred     string
	negative bool
}

// Stratify orders derivatives so that one that uses none, count or is not
// set on what others add runs only after them. Derivatives that depend on
// each other in a circle through those words can never settle, and are
// reported.
func Stratify(w model.World) ([][]model.QName, []diag.Diagnostic) {
	var derivs []model.QName
	for _, q := range w.Order {
		if w.Statements[q].Sort == model.SortDerivative {
			derivs = append(derivs, q)
		}
	}
	writers := map[string][]model.QName{}
	reads := map[model.QName][]read{}
	for _, q := range derivs {
		st := w.Statements[q]
		for _, p := range writes(w, st) {
			writers[p] = append(writers[p], q)
		}
		rs := blockReads(st.Match, false)
		for _, d := range st.Derive {
			if d.Type == model.DeriveField && d.SetTo.Field != "" {
				rs = append(rs, read{pred: fieldPred(d.SetTo.Field)})
			}
		}
		reads[q] = rs
	}
	level := map[model.QName]int{}
	limit := len(derivs)
	for range limit + 1 {
		changed := false
		for _, q := range derivs {
			for _, r := range reads[q] {
				for _, src := range writers[r.pred] {
					want := level[src]
					if r.negative {
						want++
					}
					if want > level[q] {
						level[q] = want
						changed = true
					}
				}
			}
		}
		if !changed {
			break
		}
	}
	var ds []diag.Diagnostic
	var strata [][]model.QName
	for _, q := range derivs {
		if level[q] > limit {
			ds = append(ds, diag.New(w.Statements[q].Pos, "derivative %s depends, through none, count or is not set, on what it adds itself; it can never settle (DV17)", q.Name))
			continue
		}
		for len(strata) <= level[q] {
			strata = append(strata, nil)
		}
		strata[level[q]] = append(strata[level[q]], q)
	}
	return strata, ds
}

func linkPred(q model.QName) string { return "link:" + q.Pkg + "." + q.Name }
func kindPred(q model.QName) string { return "kind:" + q.Pkg + "." + q.Name }
func fieldPred(f string) string     { return "field:" + f }

// writes lists what a derivative can add, including what links and kinds
// imply: the links a derived link sits under, and every kind a derived kind is.
func writes(w model.World, st model.Statement) []string {
	var out []string
	for _, d := range st.Derive {
		switch d.Type {
		case model.DeriveLink:
			out = append(out, linkPred(d.Link))
			for _, s := range world.SuperLinks(w, d.Link) {
				out = append(out, linkPred(s))
			}
		case model.DeriveKind:
			// The node already is the kind of its kind item, and every kind
			// that one extends; only the other kinds are new.
			var has model.QName
			for _, ki := range st.Match.Kinds {
				if ki.Var == d.Node {
					has = ki.Kind
				}
			}
			for _, k := range world.Supers(w, d.Kind) {
				if has == (model.QName{}) || !world.IsA(w, has, k) {
					out = append(out, kindPred(k))
				}
			}
		case model.DeriveField:
			out = append(out, fieldPred(d.Field))
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func blockReads(b model.Block, negative bool) []read {
	var out []read
	for _, ki := range b.Kinds {
		out = append(out, read{pred: kindPred(ki.Kind), negative: negative})
	}
	for _, li := range b.Links {
		out = append(out, read{pred: linkPred(li.Link), negative: negative})
		for _, e := range []model.End{li.From, li.To} {
			if e.Var == "" {
				out = append(out, read{pred: kindPred(e.Kind), negative: negative})
			}
		}
	}
	return append(out, conditionReads(b.Filter, negative)...)
}

func conditionReads(conds []model.Condition, negative bool) []read {
	var out []read
	for _, c := range conds {
		switch c.Type {
		case model.CondCompare:
			if c.Field != "" {
				out = append(out, read{pred: fieldPred(c.Field), negative: negative})
			}
			if c.Than.Field != "" {
				out = append(out, read{pred: fieldPred(c.Than.Field), negative: negative})
			}
		case model.CondIsA:
			out = append(out, read{pred: kindPred(c.Kind), negative: negative})
		case model.CondIsSet:
			out = append(out, read{pred: fieldPred(c.Field), negative: negative})
		case model.CondIsNotSet:
			out = append(out, read{pred: fieldPred(c.Field), negative: true})
		case model.CondAtLeastOne, model.CondNone, model.CondCount:
			out = append(out, blockReads(*c.Block, negative || c.Type != model.CondAtLeastOne)...)
		case model.CondAnyOf:
			out = append(out, conditionReads(c.Any, negative)...)
		}
	}
	return out
}
