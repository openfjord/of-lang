// SPDX-License-Identifier: Apache-2.0

package engine

import (
	"regexp"
	"slices"
	"strings"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/value"
	"github.com/openfjord/of-lang/cli/internal/world"
)

// MessageVar finds the variables a message quotes: {$var}, or {$var: field}.
var MessageVar = regexp.MustCompile(`\{(\$[a-z][A-Za-z0-9]*)(?:\s*:\s*([^}\s]+)\s*)?\}`)

// Violation is one row that breaks a constraint.
type Violation struct {
	Constraint model.QName
	Pos        diag.Pos
	Bindings   []Binding
	Message    string
	Premises   []Fact
}

// Result is what running a package over facts produces.
type Result struct {
	Graph      Graph
	Violations []Violation
	// Problems are facts that break what the kinds declare: the kinds a link
	// joins, the fields a kind requires or does not list, a union, and a
	// field that derivatives give two values.
	Problems []diag.Diagnostic
}

// Run adds what links imply and what the derivatives derive, stratum by
// stratum, until nothing new appears; then it checks the constraints and
// what the kinds require.
func Run(w model.World, strata [][]model.QName, g Graph) Result {
	var problems []diag.Diagnostic
	builtins(w, &g)
	for _, stratum := range strata {
		for {
			changed := false
			for _, q := range stratum {
				added, ps := applyDerivative(w, &g, w.Statements[q])
				problems = append(problems, ps...)
				changed = changed || added
			}
			if builtins(w, &g) {
				changed = true
			}
			if !changed {
				break
			}
		}
	}
	res := Result{Graph: g}
	for _, q := range w.Order {
		st := w.Statements[q]
		if st.Sort == model.SortConstraint {
			res.Violations = append(res.Violations, violations(w, g, st)...)
		}
	}
	problems = append(problems, linkProblems(w, g)...)
	problems = append(problems, fieldProblems(w, g)...)
	problems = append(problems, unionProblems(w, g)...)
	res.Problems = diag.Sorted(problems)
	return res
}

// builtins adds what under and chains imply, and the defaults of fields
// that a derived kind brings, until nothing new appears.
func builtins(w model.World, g *Graph) bool {
	grew := false
	for {
		changed := false
		for _, f := range slices.Clone(g.Order) {
			switch f.Type {
			case FactLink:
				e := f.Edge
				link := w.Links[e.Link]
				if link.Under != nil && Add(g, LinkFact(e.From, *link.Under, e.To), Reason{Type: ReasonUnder, Premises: []Fact{f}}) {
					changed = true
				}
				if !link.Chains {
					continue
				}
				for _, next := range slices.Clone(g.out[e.To][e.Link]) {
					if Add(g, LinkFact(e.From, e.Link, next), Reason{Type: ReasonChains, Premises: []Fact{f, LinkFact(e.To, e.Link, next)}}) {
						changed = true
					}
				}
			case FactKind:
				if AddDefaults(w, g, f.Node, f.Kind, Reason{Type: ReasonDefault, Premises: []Fact{f}}) {
					changed = true
				}
			default:
				// A field value implies nothing.
			}
		}
		if !changed {
			return grew
		}
		grew = true
	}
}

// AddDefaults gives a thing the default of each field of a kind that it has
// no value for.
func AddDefaults(w model.World, g *Graph, id string, kind model.QName, r Reason) bool {
	added := false
	for _, f := range world.Fields(w, kind) {
		if f.Schema.Default == nil {
			continue
		}
		if _, ok := g.Fields[id][f.Name]; ok {
			continue
		}
		if Add(g, FieldFact(id, f.Name, *f.Schema.Default), Reason{Type: r.Type, Pos: f.Pos, Premises: r.Premises}) {
			added = true
		}
	}
	return added
}

func applyDerivative(w model.World, g *Graph, st model.Statement) (bool, []diag.Diagnostic) {
	added := false
	var problems []diag.Diagnostic
	vars := BlockVars(st.Match)
	for _, m := range Find(w, *g, st.Match, nil, false) {
		reason := Reason{Type: ReasonDerivative, Derivative: st.Q, Bindings: sortedBindings(m.Env, vars), Premises: m.Premises}
		for _, d := range st.Derive {
			f, ok := derived(*g, d, m.Env)
			if !ok {
				continue
			}
			if f.Type == FactField {
				if cur, set := g.Fields[f.Node][f.Field]; set {
					switch {
					case value.Equal(cur, f.Value):
						// The derivative agrees with the value the field has.
					case SetDefault(g, f, reason):
						added = true
					default:
						problems = append(problems, diag.New(d.Pos, "derivative %s sets %s of %s to %s, but it is already %s; a field has one value (DV15)", st.Q.Name, f.Field, f.Node, value.Format(f.Value), value.Format(cur)))
					}
					continue
				}
			}
			if Add(g, f, reason) {
				added = true
			}
		}
	}
	return added, problems
}

func derived(g Graph, d model.Derivation, env map[string]string) (Fact, bool) {
	switch d.Type {
	case model.DeriveLink:
		from, to := env[d.From], env[d.To]
		if d.Inv {
			from, to = to, from
		}
		return LinkFact(from, d.Link, to), from != "" && to != ""
	case model.DeriveKind:
		return KindFact(env[d.Node], d.Kind), env[d.Node] != ""
	case model.DeriveField:
		vs, _, ok := operandValues(g, env, d.SetTo)
		if !ok || len(vs) != 1 {
			return Fact{}, false
		}
		return FieldFact(env[d.Node], d.Field, vs[0]), env[d.Node] != ""
	default:
		return Fact{}, false
	}
}

func violations(w model.World, g Graph, st model.Statement) []Violation {
	var out []Violation
	vars := BlockVars(st.Match)
	for _, m := range Find(w, g, st.Match, nil, false) {
		if len(st.Unless) > 0 {
			if _, ok := allHold(w, g, st.Unless, m); ok {
				continue
			}
		}
		out = append(out, Violation{
			Constraint: st.Q,
			Pos:        st.Pos,
			Bindings:   sortedBindings(m.Env, vars),
			Message:    Render(g, st.Message, m.Env),
			Premises:   m.Premises,
		})
	}
	return out
}

// Render fills a message's {$var} and {$var: field} with what the row gives:
// the id of the thing, or the value of the field, or "not set".
func Render(g Graph, msg string, env map[string]string) string {
	return MessageVar.ReplaceAllStringFunc(msg, func(s string) string {
		sub := MessageVar.FindStringSubmatch(s)
		id, ok := env[sub[1]]
		switch {
		case !ok:
			return s
		case sub[2] == "":
			return id
		}
		v, ok := g.Fields[id][sub[2]]
		if !ok {
			return "not set"
		}
		return value.Format(v)
	})
}

// linkProblems reports links whose ends are not things of the kinds the link
// joins.
func linkProblems(w model.World, g Graph) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, f := range g.Order {
		if f.Type != FactLink {
			continue
		}
		e := f.Edge
		link := w.Links[e.Link]
		for _, end := range []struct {
			id   string
			want model.QName
		}{{e.From, link.From}, {e.To, link.To}} {
			if !thingIsA(w, g, end.id, end.want) {
				out = append(out, diag.New(g.Pos[end.id], "%s links %s to %s, but %s is not a %s (FC5)", e.Link.Name, e.From, e.To, end.id, end.want.Name))
			}
		}
	}
	return out
}

// fieldProblems reports, for each thing, a field that one of its kinds
// requires and it has no value for, and a value of a field that none of its
// kinds lists.
func fieldProblems(w model.World, g Graph) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, id := range g.Nodes {
		listed := map[string]bool{}
		for _, k := range g.Kinds[id] {
			for _, f := range world.Fields(w, k) {
				listed[f.Name] = true
				if f.Optional || f.Schema.Default != nil {
					continue
				}
				if _, ok := g.Fields[id][f.Name]; !ok {
					out = append(out, diag.New(g.Pos[id], "%s needs %s, which every %s has (FC3)", id, f.Name, k.Name))
				}
			}
		}
		for name := range g.Fields[id] {
			if !listed[name] {
				out = append(out, diag.New(g.Pos[id], "%s has %s, which none of its kinds lists (FC3)", id, name))
			}
		}
	}
	return out
}

// unionProblems reports a thing that is two kinds that one union lists.
func unionProblems(w model.World, g Graph) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, u := range world.Unions(w) {
		for _, id := range g.Nodes {
			var is []string
			for _, l := range u.OneOf {
				if thingIsA(w, g, id, l) {
					is = append(is, l.Name)
				}
			}
			if len(is) > 1 {
				out = append(out, diag.New(g.Pos[id], "%s is a %s, so it must be exactly one of its kinds, but it is %s (KD11)", id, u.Q.Name, strings.Join(is, " and ")))
			}
		}
	}
	return out
}

// Sentence renders a fact in plain words, using links' reads sentences.
func Sentence(w model.World, f Fact) string {
	switch f.Type {
	case FactLink:
		return world.Reads(w, f.Edge.Link, f.Edge.From, f.Edge.To)
	case FactKind:
		return f.Node + " is a " + f.Kind.Name
	default:
		return f.Node + " has " + f.Field + " " + value.Format(f.Value)
	}
}

// BindingsText renders bindings as $var = id pairs.
func BindingsText(bs []Binding) string {
	parts := make([]string, 0, len(bs))
	for _, b := range bs {
		parts = append(parts, b.Var+" = "+b.ID)
	}
	return strings.Join(parts, ", ")
}
