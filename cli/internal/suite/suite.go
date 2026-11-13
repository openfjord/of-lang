// Package suite runs a package's tests: examples that expect some facts to be
// derived and some constraints to reject, and questions that expect answers.
package suite

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"oflang/internal/diag"
	"oflang/internal/engine"
	"oflang/internal/facts"
	"oflang/internal/load"
	"oflang/internal/model"
	"oflang/internal/source"
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

// Run reads every YAML file under the package's tests folder and runs what it
// finds. Diagnostics are problems in the test files themselves.
func Run(l load.Loaded, strata [][]model.QName) ([]Outcome, []diag.Diagnostic, error) {
	dir := filepath.Join(l.World.Packages[l.World.Root].Dir, "tests")
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".yaml") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("read tests: %w", err)
	}
	var outcomes []Outcome
	var ds []diag.Diagnostic
	for _, path := range paths {
		f, fds, err := facts.Read(l, path)
		if err != nil {
			return nil, nil, err
		}
		ds = append(ds, fds...)
		if len(f.Tests) > 1 {
			ds = append(ds, diag.New(f.Tests[1].Pos, "an example file holds one Test; %s is the second", f.Tests[1].Name))
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
	source.Allow(m, []string{"is", "says", "expect"}, ds)
	o.Says, _ = source.GetScalar(m, "says", ds)
	for _, p := range res.Problems {
		o.Failures = append(o.Failures, "problem: "+diag.Format(p))
	}
	en, ok := source.Get(m, "expect")
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "test %s needs expect:", e.Name))
		return o
	}
	em, ok := source.AsMap(f.Path, en, ds)
	if !ok {
		return o
	}
	source.Allow(em, []string{"derived", "notDerived", "rejected"}, ds)
	for _, key := range []string{"derived", "notDerived"} {
		n, ok := source.Get(em, key)
		if !ok {
			continue
		}
		for _, fact := range expectedFacts(f, n, ds) {
			has := engine.Has(res.Graph, fact)
			switch {
			case key == "derived" && !has:
				o.Failures = append(o.Failures, "expected derived, but not: "+engine.Sentence(w, fact))
			case key == "notDerived" && has:
				o.Failures = append(o.Failures, "expected not derived, but was: "+engine.Sentence(w, fact))
			}
		}
	}
	var want []string
	if n, ok := source.Get(em, "rejected"); ok {
		names, _ := source.Strings(f.Path, n, ds)
		for _, name := range names {
			ref, ok := load.Resolve(f.Scope, name, model.SortConstraint, source.PosOf(f.Path, n), ds)
			if ok {
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

func expectedFacts(f facts.File, n *yaml.Node, ds *[]diag.Diagnostic) []engine.Fact {
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(source.PosOf(f.Path, n), "expected a list of facts"))
		return nil
	}
	items := n.Content
	var out []engine.Fact
	for _, item := range items {
		m, ok := source.AsMap(f.Path, item, ds)
		if !ok {
			continue
		}
		pos := m.Pos
		switch {
		case has(m, "link"):
			source.Allow(m, []string{"from", "link", "to"}, ds)
			from, _ := source.GetScalar(m, "from", ds)
			to, _ := source.GetScalar(m, "to", ds)
			name, _ := source.GetScalar(m, "link", ds)
			ref, ok := load.Resolve(f.Scope, name, model.SortLink, pos, ds)
			if !ok {
				continue
			}
			if ref.Inverted {
				from, to = to, from
			}
			out = append(out, engine.LinkFact(from, ref.Q, to))
		case has(m, "becomes"):
			source.Allow(m, []string{"node", "becomes"}, ds)
			node, _ := source.GetScalar(m, "node", ds)
			name, _ := source.GetScalar(m, "becomes", ds)
			ref, ok := load.Resolve(f.Scope, name, model.SortKind, pos, ds)
			if !ok {
				continue
			}
			out = append(out, engine.KindFact(node, ref.Q))
		case has(m, "set to"):
			source.Allow(m, []string{"node", "field", "set to"}, ds)
			node, _ := source.GetScalar(m, "node", ds)
			field, _ := source.GetScalar(m, "field", ds)
			vn, _ := source.Get(m, "set to")
			v, ok := load.Value(f.Path, vn, ds)
			if !ok {
				continue
			}
			out = append(out, engine.FieldFact(node, field, v))
		default:
			*ds = append(*ds, diag.New(pos, "an expected fact is a link (from, link, to), a kind (node, becomes) or a field (node, field, set to)"))
		}
	}
	return out
}

func runQuestion(l load.Loaded, strata [][]model.QName, f facts.File, e source.Entry, ds *[]diag.Diagnostic) (Outcome, error) {
	o := Outcome{Name: e.Name, Sort: string(model.SortQuestion), File: f.Path, Failures: []string{}}
	m, ok := source.AsMap(f.Path, e.Node, ds)
	if !ok {
		return o, nil
	}
	st := load.Statement(model.QName{Pkg: l.World.Root, Name: e.Name}, model.SortQuestion, f.Path, m, f.Scope, "asks", []string{"given", "expect"}, ds)
	o.Says = st.Says
	given := f
	if gn, ok := source.GetScalar(m, "given", ds); ok {
		rel, found := strings.CutPrefix(gn, "file:")
		if !found {
			*ds = append(*ds, diag.New(m.Pos, "given: names a file with file:"))
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
	en, ok := source.Get(m, "expect")
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "question %s needs expect:", e.Name))
		return o, nil
	}
	em, ok := source.AsMap(f.Path, en, ds)
	if !ok {
		return o, nil
	}
	var vars []string
	for _, p := range em.Pairs {
		if !slices.Contains(st.Vars, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not a variable of the question", p.Key))
			continue
		}
		vars = append(vars, p.Key)
	}
	answers := Answers(l.World, res.Graph, st, vars)
	for _, p := range em.Pairs {
		if !slices.Contains(vars, p.Key) {
			continue
		}
		want, _ := source.Strings(f.Path, p.Value, ds)
		slices.Sort(want)
		want = slices.Compact(want)
		got := answers[p.Key]
		if !slices.Equal(want, got) {
			o.Failures = append(o.Failures, fmt.Sprintf("%s: expected [%s], got [%s]", p.Key, strings.Join(want, ", "), strings.Join(got, ", ")))
		}
	}
	return o, nil
}

// Answers returns, for each variable, the sorted distinct things it stands for
// across every match of a question.
func Answers(w model.World, g engine.Graph, st model.Statement, vars []string) map[string][]string {
	out := map[string][]string{}
	all := append(engine.OuterVars(st), vars...)
	for _, m := range engine.Find(w, g, st.Links, st.Where, all, nil) {
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

func has(m source.Map, key string) bool {
	_, ok := source.Get(m, key)
	return ok
}
