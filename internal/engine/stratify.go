// SPDX-License-Identifier: Apache-2.0

package engine

import (
	"slices"
	"strings"

	"github.com/openfjord/of-lang/internal/diag"
	"github.com/openfjord/of-lang/internal/model"
	"github.com/openfjord/of-lang/internal/world"
)

type read struct {
	pred     string
	negative bool
}

// Stratify orders rules so a rule that asks whether something is absent, or
// counts it, runs only after every rule that could add it. Rules that depend
// on their own absence cannot be ordered and are reported.
func Stratify(w model.World) ([][]model.QName, []diag.Diagnostic) {
	var rules []model.QName
	for _, q := range w.Order {
		if w.Statements[q].Sort == model.SortRule {
			rules = append(rules, q)
		}
	}
	writers := map[string][]model.QName{}
	reads := map[model.QName][]read{}
	for _, q := range rules {
		st := w.Statements[q]
		for _, p := range writes(w, st) {
			writers[p] = append(writers[p], q)
		}
		var rs []read
		for _, pl := range st.Links {
			rs = append(rs, read{pred: linkPred(pl.Link)})
		}
		rs = append(rs, conditionReads(st.Where, false)...)
		for _, d := range st.Derive {
			if d.Type == model.DeriveField && d.SetTo.Node != "" && d.SetTo.Field != "" {
				rs = append(rs, read{pred: "field:" + d.SetTo.Field})
			}
		}
		reads[q] = rs
	}
	level := map[model.QName]int{}
	limit := len(rules)
	for range limit + 1 {
		changed := false
		for _, q := range rules {
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
	for _, q := range rules {
		if level[q] > limit {
			ds = append(ds, diag.New(w.Statements[q].Pos, "rule %s depends, through none:, count: or is not set:, on what it adds itself; it can never settle", q.Name))
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

// writes lists what a rule can add, including the links its links sit under
// and the kinds its kinds extend, since built-in inference adds those too.
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
			for _, a := range world.Ancestors(w, d.Kind) {
				out = append(out, kindPred(a))
			}
		case model.DeriveField:
			out = append(out, "field:"+d.Field)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func conditionReads(conds []model.Condition, negative bool) []read {
	var out []read
	for _, c := range conds {
		switch c.Type {
		case model.CondCompare:
			if c.Field != "" {
				out = append(out, read{pred: "field:" + c.Field, negative: negative})
			}
			if c.Than.Field != "" {
				out = append(out, read{pred: "field:" + c.Than.Field, negative: negative})
			}
		case model.CondIsA:
			out = append(out, read{pred: kindPred(c.Kind), negative: negative})
		case model.CondIsSet:
			out = append(out, read{pred: "field:" + c.Field, negative: negative})
		case model.CondIsNotSet:
			out = append(out, read{pred: "field:" + c.Field, negative: true})
		case model.CondAtLeastOne, model.CondNone, model.CondCount:
			neg := negative || c.Type != model.CondAtLeastOne
			for _, pl := range c.Links {
				out = append(out, read{pred: linkPred(pl.Link), negative: neg})
			}
			out = append(out, conditionReads(c.Where, neg)...)
		case model.CondAnyOf:
			out = append(out, conditionReads(c.Any, negative)...)
		}
	}
	return out
}

// RuleNames renders rule names for messages.
func RuleNames(qs []model.QName) string {
	names := make([]string, 0, len(qs))
	for _, q := range qs {
		names = append(names, q.Name)
	}
	return strings.Join(names, ", ")
}
