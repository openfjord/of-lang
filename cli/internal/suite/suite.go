// SPDX-License-Identifier: Apache-2.0

// Package suite runs a package's tests: examples that expect facts and terms
// to be derived and constraints to reject, and questions that expect answers.
package suite

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/check"
	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/engine"
	"github.com/openfjord/of-lang/cli/internal/facts"
	"github.com/openfjord/of-lang/cli/internal/load"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
	"github.com/openfjord/of-lang/cli/internal/world"
)

// Outcome is the result of one test or question.
type Outcome struct {
	Name     string   `json:"name"`
	Sort     string   `json:"sort"`
	File     string   `json:"file"`
	Says     string   `json:"says"`
	Failures []string `json:"failures"`
}

// Passed says whether a test or question met every expectation.
func Passed(o Outcome) bool { return len(o.Failures) == 0 }

// Run runs every test file of the root package. Diagnostics are problems in
// the test files themselves.
func Run(l load.Loaded, strata [][]model.QName) ([]Outcome, []diag.Diagnostic, error) {
	var outcomes []Outcome
	var ds []diag.Diagnostic
	for _, path := range l.World.Packages[l.World.Root].Tests {
		f, fds, err := facts.Read(l, path)
		if err != nil {
			return nil, nil, err
		}
		ds = append(ds, fds...)
		for _, e := range append(slices.Clone(f.Tests), f.Questions...) {
			if !check.LowerCamel(e.Name) {
				ds = append(ds, diag.New(e.Pos, "%s is not lowerCamelCase (FL7)", e.Name))
			}
		}
		switch {
		case len(f.Tests) > 1:
			ds = append(ds, diag.New(f.Tests[1].Pos, "an example file holds one Test; %s is the second (TS1)", f.Tests[1].Name))
		case len(f.Tests) == 0 && len(f.Questions) == 0:
			ds = append(ds, diag.New(diag.Pos{File: f.Path, Line: 1, Column: 1}, "a test file holds one Test, or one or more Questions (TS1)"))
		}
		if len(f.Tests) > 0 {
			res := engine.Run(l.World, strata, f.Graph)
			outcomes = append(outcomes, runTest(l.World, f, f.Tests[0], res, &ds))
		}
		for _, q := range f.Questions {
			o, err := runQuestion(l, strata, f, q, &ds)
			if err != nil {
				return nil, nil, err
			}
			outcomes = append(outcomes, o)
		}
	}
	return outcomes, diag.Sorted(ds), nil
}

func runTest(w model.World, f facts.File, e source.Entry, res engine.Result, ds *[]diag.Diagnostic) Outcome {
	o := Outcome{Name: e.Name, Sort: string(model.SortTest), File: f.Path, Failures: []string{}}
	m, ok := source.AsMap(f.Path, e.Node, ds)
	if !ok {
		return o
	}
	for _, p := range m.Pairs {
		if !slices.Contains([]string{"is", "says", "expect"}, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not allowed here; a test has is, says and expect (FL2)", p.Key))
		}
	}
	says, ok := source.GetScalar(m, "says", ds)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "test %s needs says: (TS1)", e.Name))
	}
	o.Says = strings.TrimSpace(says)
	for _, p := range res.Problems {
		o.Failures = append(o.Failures, "problem: "+diag.Format(p))
	}
	en, ok := source.Get(m, "expect")
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "test %s needs expect: (TS1)", e.Name))
		return o
	}
	em, ok := source.AsMap(f.Path, en, ds)
	if !ok {
		return o
	}
	for _, p := range em.Pairs {
		if !slices.Contains([]string{"derived", "notDerived", "rejected"}, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not allowed here; expect has derived, notDerived and rejected (FL2)", p.Key))
		}
	}
	noFacts := f.Facts == 0
	for _, key := range []string{"derived", "notDerived"} {
		n, ok := source.Get(em, key)
		if !ok {
			continue
		}
		for _, x := range expected(w, f, n, noFacts, ds) {
			holds := x.holds(res.Graph)
			switch {
			case key == "derived" && !holds:
				o.Failures = append(o.Failures, "expected derived, but not: "+x.text)
			case key == "notDerived" && holds:
				o.Failures = append(o.Failures, "expected not derived, but was: "+x.text)
			}
		}
	}
	var want []string
	if p, ok := pair(em, "rejected"); ok {
		if noFacts {
			*ds = append(*ds, diag.New(p.Pos, "a test with no facts lists only terms, and no rejected (TS15)"))
		}
		names, _ := source.Strings(f.Path, p.Value, ds)
		for _, name := range names {
			if ref, ok := load.Resolve(f.Scope, name, model.SortConstraint, p.Pos, "TS3", ds); ok {
				want = append(want, ref.Q.Name)
			}
		}
	}
	var got []string
	for _, v := range res.Violations {
		if !slices.Contains(got, v.Constraint.Name) {
			got = append(got, v.Constraint.Name)
		}
		if !slices.Contains(want, v.Constraint.Name) {
			o.Failures = append(o.Failures, "rejected unexpectedly by "+v.Constraint.Name+": "+v.Message)
		}
	}
	for _, name := range want {
		if !slices.Contains(got, name) {
			o.Failures = append(o.Failures, "expected rejection by "+name+", but nothing broke it")
		}
	}
	return o
}

// expectation is one fact or term under derived or notDerived.
type expectation struct {
	fact *engine.Fact
	term func() bool
	text string
}

func (x expectation) holds(g engine.Graph) bool {
	if x.fact != nil {
		return engine.Has(g, *x.fact)
	}
	return x.term()
}

func expected(w model.World, f facts.File, n *yaml.Node, noFacts bool, ds *[]diag.Diagnostic) []expectation {
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(source.PosOf(f.Path, n), "expected a list of facts and terms (TS15)"))
		return nil
	}
	var out []expectation
	for _, item := range n.Content {
		m, ok := source.AsMap(f.Path, item, ds)
		if !ok {
			continue
		}
		x, isFact, ok := expectationOf(w, f, m, ds)
		if !ok {
			continue
		}
		if isFact && noFacts {
			*ds = append(*ds, diag.New(m.Pos, "a test with no facts lists only terms (TS15)"))
			continue
		}
		out = append(out, x)
	}
	return out
}

// expectationOf reads one item by the keys it has, and says whether it is a
// fact rather than a term.
func expectationOf(w model.World, f facts.File, m source.Map, ds *[]diag.Diagnostic) (expectation, bool, bool) {
	keys := make([]string, 0, len(m.Pairs))
	for _, p := range m.Pairs {
		keys = append(keys, p.Key)
	}
	slices.Sort(keys)
	shape := strings.Join(keys, ",")
	text := func(key string) string {
		s, _ := source.GetScalar(m, key, ds)
		return s
	}
	ref := func(key string, sort model.Sort) (model.Ref, bool) {
		return load.Resolve(f.Scope, text(key), sort, m.Pos, "TS3", ds)
	}
	switch shape {
	case "from,link,to":
		l, ok := ref("link", model.SortLink)
		if !ok {
			return expectation{}, true, false
		}
		from, to := text("from"), text("to")
		if l.Inverted {
			from, to = to, from
		}
		return factExpectation(w, engine.LinkFact(from, l.Q, to)), true, true
	case "becomes,node":
		k, ok := ref("becomes", model.SortKind)
		if !ok {
			return expectation{}, true, false
		}
		return factExpectation(w, engine.KindFact(text("node"), k.Q)), true, true
	case "field,node,set to":
		vn, _ := source.Get(m, "set to")
		vm, ok := source.AsMap(f.Path, vn, ds)
		if !ok {
			return expectation{}, true, false
		}
		n, ok := source.Get(vm, "value")
		if !ok || len(vm.Pairs) != 1 {
			*ds = append(*ds, diag.New(vm.Pos, "set to in a test takes value: (TS15)"))
			return expectation{}, true, false
		}
		v, ok := load.Value(f.Path, n, ds)
		if !ok {
			return expectation{}, true, false
		}
		return factExpectation(w, engine.FieldFact(text("node"), text("field"), v)), true, true
	case "is a,kind", "kind,never is":
		key := "is a"
		if shape == "kind,never is" {
			key = "never is"
		}
		a, ok1 := ref("kind", model.SortKind)
		b, ok2 := ref(key, model.SortKind)
		if !ok1 || !ok2 {
			return expectation{}, false, false
		}
		term := func() bool { return world.IsA(w, a.Q, b.Q) }
		if key == "never is" {
			term = func() bool { return neverShare(w, a.Q, b.Q) }
		}
		return expectation{term: term, text: a.Q.Name + " " + key + " " + b.Q.Name}, false, true
	case "can have,kind":
		k, ok1 := ref("kind", model.SortKind)
		l, ok2 := ref("can have", model.SortLink)
		if !ok1 || !ok2 {
			return expectation{}, false, false
		}
		term := func() bool {
			start, _ := world.LinkEnds(w, l.Q, l.Inverted)
			return world.IsA(w, k.Q, start)
		}
		return expectation{term: term, text: k.Q.Name + " can have " + text("can have")}, false, true
	case "implies,link", "inverse of,link":
		key := strings.TrimSuffix(shape, ",link")
		a, ok1 := ref("link", model.SortLink)
		b, ok2 := ref(key, model.SortLink)
		if !ok1 || !ok2 {
			return expectation{}, false, false
		}
		term := func() bool {
			if key == "inverse of" {
				return a.Q == b.Q && a.Inverted != b.Inverted
			}
			return a.Inverted == b.Inverted && slices.Contains(world.SuperLinks(w, a.Q), b.Q)
		}
		return expectation{term: term, text: text("link") + " " + key + " " + text(key)}, false, true
	case "chains,link":
		l, ok := ref("link", model.SortLink)
		if !ok {
			return expectation{}, false, false
		}
		cn, _ := source.Get(m, "chains")
		want, ok := source.Bool(f.Path, cn, ds)
		if !ok || !want {
			*ds = append(*ds, diag.New(m.Pos, "a chains term is written chains: true (TS15)"))
			return expectation{}, false, false
		}
		return expectation{term: func() bool { return w.Links[l.Q].Chains }, text: text("link") + " chains"}, false, true
	default:
		*ds = append(*ds, diag.New(m.Pos, "an expected item is a fact (link, from, to; node, becomes; node, field, set to) or a term (kind with is a, never is or can have; link with implies, chains or inverse of) (TS15)"))
		return expectation{}, false, false
	}
}

func factExpectation(w model.World, f engine.Fact) expectation {
	return expectation{fact: &f, text: engine.Sentence(w, f)}
}

// neverShare says whether a union keeps two kinds apart: each is, or
// extends, a different kind that the union lists.
func neverShare(w model.World, a, b model.QName) bool {
	if world.IsA(w, a, b) || world.IsA(w, b, a) {
		return false
	}
	for _, u := range world.Unions(w) {
		for _, x := range u.OneOf {
			for _, y := range u.OneOf {
				if x != y && world.IsA(w, a, x) && world.IsA(w, b, y) {
					return true
				}
			}
		}
	}
	return false
}

func runQuestion(l load.Loaded, strata [][]model.QName, f facts.File, e source.Entry, ds *[]diag.Diagnostic) (Outcome, error) {
	o := Outcome{Name: e.Name, Sort: string(model.SortQuestion), File: f.Path, Failures: []string{}}
	m, ok := source.AsMap(f.Path, e.Node, ds)
	if !ok {
		return o, nil
	}
	st := load.Statement(model.QName{Pkg: l.World.Root, Name: e.Name}, model.SortQuestion, f.Path, m, f.Scope, ds)
	o.Says = st.Says
	ep, ok := pair(m, "expect")
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "question %s needs expect: (TS9)", e.Name))
		return o, nil
	}
	em, ok := source.AsMap(f.Path, ep.Value, ds)
	if !ok {
		return o, nil
	}
	matchVars := engine.BlockVars(st.Match)
	var vars []string
	for _, p := range em.Pairs {
		if !slices.Contains(matchVars, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not a variable of the match of %s (TS9)", p.Key, e.Name))
			continue
		}
		vars = append(vars, p.Key)
	}
	*ds = append(*ds, check.Statement(l.World, st, vars)...)
	given := f
	if gn, ok := source.GetScalar(m, "given", ds); ok {
		rel, found := strings.CutPrefix(gn, "file:")
		if !found {
			*ds = append(*ds, diag.New(m.Pos, "given: names an example file with file: (FL2)"))
			return o, nil
		}
		g, gds, err := facts.Read(l, filepath.Join(filepath.Dir(f.Path), rel))
		if err != nil {
			return o, err
		}
		*ds = append(*ds, gds...)
		given = g
	}
	res := engine.Run(l.World, strata, given.Graph)
	answers := Answers(l.World, res.Graph, st, vars)
	for _, p := range em.Pairs {
		if !slices.Contains(vars, p.Key) {
			continue
		}
		want, _ := source.Strings(f.Path, p.Value, ds)
		slices.Sort(want)
		want = slices.Compact(want)
		if got := answers[p.Key]; !slices.Equal(want, got) {
			o.Failures = append(o.Failures, fmt.Sprintf("%s: expected [%s], got [%s]", p.Key, strings.Join(want, ", "), strings.Join(got, ", ")))
		}
	}
	return o, nil
}

// Answers returns, for each variable, the sorted distinct things it stands for
// across every match of a question.
func Answers(w model.World, g engine.Graph, st model.Statement, vars []string) map[string][]string {
	out := map[string][]string{}
	for _, m := range engine.Find(w, g, st.Match, nil, false) {
		for _, v := range vars {
			if id, ok := m.Env[v]; ok && !slices.Contains(out[v], id) {
				out[v] = append(out[v], id)
			}
		}
	}
	for v := range out {
		slices.Sort(out[v])
	}
	return out
}

func pair(m source.Map, key string) (source.Pair, bool) {
	for _, p := range m.Pairs {
		if p.Key == key {
			return p, true
		}
	}
	return source.Pair{}, false
}
