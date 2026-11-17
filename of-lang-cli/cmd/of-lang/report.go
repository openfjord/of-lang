// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/openfjord/of-lang/internal/diag"
	"github.com/openfjord/of-lang/internal/engine"
	"github.com/openfjord/of-lang/internal/model"
	"github.com/openfjord/of-lang/internal/suite"
	"github.com/openfjord/of-lang/internal/value"
)

// reporter writes results as text for people or JSON for tools.
type reporter struct {
	out  io.Writer
	json bool
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}

func writeLines(w io.Writer, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	if _, err := io.WriteString(w, strings.Join(lines, "\n")+"\n"); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

func (r reporter) diagnostics(ds []diag.Diagnostic) error {
	if r.json {
		if ds == nil {
			ds = []diag.Diagnostic{}
		}
		return writeJSON(r.out, map[string]any{"problems": ds})
	}
	lines := make([]string, 0, len(ds)+1)
	for _, d := range ds {
		lines = append(lines, diag.Format(d))
	}
	if len(ds) == 0 {
		lines = append(lines, "ok")
	}
	return writeLines(r.out, lines)
}

type jsonFact struct {
	Type     string `json:"type"`
	Sentence string `json:"sentence"`
	From     string `json:"from,omitzero"`
	Link     string `json:"link,omitzero"`
	To       string `json:"to,omitzero"`
	Node     string `json:"node,omitzero"`
	Kind     string `json:"kind,omitzero"`
	Field    string `json:"field,omitzero"`
	Value    any    `json:"value,omitzero"`
}

func factJSON(w model.World, f engine.Fact) jsonFact {
	j := jsonFact{Type: string(f.Type), Sentence: engine.Sentence(w, f)}
	switch f.Type {
	case engine.FactLink:
		j.From, j.Link, j.To = f.Edge.From, f.Edge.Link.Name, f.Edge.To
	case engine.FactKind:
		j.Node, j.Kind = f.Node, f.Kind.Name
	default:
		j.Node, j.Field, j.Value = f.Node, f.Field, value.Plain(f.Value)
	}
	return j
}

func reasonText(r engine.Reason) string {
	switch r.Type {
	case engine.ReasonRule:
		return "rule " + r.Rule.Name
	case engine.ReasonUnder:
		return "under"
	case engine.ReasonChains:
		return "chains"
	default:
		return "asserted"
	}
}

// derivedFacts lists the facts the run added, in the order it added them.
func derivedFacts(g engine.Graph) []engine.Fact {
	var out []engine.Fact
	for _, f := range g.Order {
		if g.Reasons[engine.Key(f)].Type != engine.ReasonAsserted {
			out = append(out, f)
		}
	}
	return out
}

type jsonViolation struct {
	Constraint string           `json:"constraint"`
	Message    string           `json:"message"`
	Bindings   []engine.Binding `json:"bindings"`
}

func violationJSON(v engine.Violation) jsonViolation {
	return jsonViolation{Constraint: v.Constraint.Name, Message: v.Message, Bindings: v.Bindings}
}

func (r reporter) result(w model.World, res engine.Result) error {
	derived := derivedFacts(res.Graph)
	if r.json {
		out := struct {
			Derived  []jsonFact        `json:"derived"`
			Rejected []jsonViolation   `json:"rejected"`
			Problems []diag.Diagnostic `json:"problems"`
		}{Derived: []jsonFact{}, Rejected: []jsonViolation{}, Problems: relative(res.Problems)}
		for _, f := range derived {
			j := factJSON(w, f)
			out.Derived = append(out.Derived, j)
		}
		for _, v := range res.Violations {
			out.Rejected = append(out.Rejected, violationJSON(v))
		}
		if out.Problems == nil {
			out.Problems = []diag.Diagnostic{}
		}
		return writeJSON(r.out, out)
	}
	lines := []string{"derived:"}
	if len(derived) == 0 {
		lines = append(lines, "  nothing")
	}
	for _, f := range derived {
		lines = append(lines, fmt.Sprintf("  %s  (%s)", engine.Sentence(w, f), reasonText(res.Graph.Reasons[engine.Key(f)])))
	}
	lines = append(lines, "rejected:")
	if len(res.Violations) == 0 {
		lines = append(lines, "  nothing")
	}
	for _, v := range res.Violations {
		lines = append(lines, fmt.Sprintf("  %s: %s", v.Constraint.Name, v.Message))
	}
	if len(res.Problems) > 0 {
		lines = append(lines, "problems:")
		for _, d := range relative(res.Problems) {
			lines = append(lines, "  "+diag.Format(d))
		}
	}
	return writeLines(r.out, lines)
}

func (r reporter) notHeld(w model.World, f engine.Fact) error {
	if r.json {
		return writeJSON(r.out, map[string]any{"fact": factJSON(w, f), "holds": false})
	}
	return writeLines(r.out, []string{"does not hold: " + engine.Sentence(w, f)})
}

type jsonExplanation struct {
	Fact     jsonFact          `json:"fact"`
	Reason   string            `json:"reason"`
	At       string            `json:"at,omitzero"`
	Rule     string            `json:"rule,omitzero"`
	Says     string            `json:"says,omitzero"`
	Bindings []engine.Binding  `json:"bindings,omitzero"`
	Because  []jsonExplanation `json:"because,omitzero"`
}

func explain(w model.World, g engine.Graph, f engine.Fact, seen []string) jsonExplanation {
	k := engine.Key(f)
	r := g.Reasons[k]
	j := jsonExplanation{Fact: factJSON(w, f), Reason: string(r.Type)}
	if r.Type == engine.ReasonAsserted && r.Pos.File != "" {
		j.At = fmt.Sprintf("%s:%d:%d", relPath(r.Pos.File), r.Pos.Line, r.Pos.Column)
	}
	if r.Type == engine.ReasonRule {
		j.Rule, j.Says, j.Bindings = r.Rule.Name, w.Statements[r.Rule].Says, r.Bindings
	}
	if slices.Contains(seen, k) {
		return j
	}
	seen = append(seen, k)
	for _, p := range r.Premises {
		j.Because = append(j.Because, explain(w, g, p, seen))
	}
	return j
}

func (r reporter) explanation(w model.World, g engine.Graph, f engine.Fact) error {
	e := explain(w, g, f, nil)
	if r.json {
		return writeJSON(r.out, e)
	}
	var lines []string
	var walk func(e jsonExplanation, depth int)
	walk = func(e jsonExplanation, depth int) {
		pad := strings.Repeat("  ", depth)
		head := pad + e.Fact.Sentence
		switch {
		case e.Rule != "":
			head += "  (rule " + e.Rule + ")"
		case e.At != "":
			head += "  (written at " + e.At + ")"
		default:
			head += "  (" + e.Reason + ")"
		}
		lines = append(lines, head)
		if e.Rule != "" {
			lines = append(lines, pad+"  "+strings.TrimSpace(e.Says))
			lines = append(lines, pad+"  with "+engine.BindingsText(e.Bindings))
			if len(e.Because) > 0 {
				lines = append(lines, pad+"  because:")
			}
		}
		for _, c := range e.Because {
			walk(c, depth+2)
		}
	}
	walk(e, 0)
	return writeLines(r.out, lines)
}

func violationKey(v engine.Violation) string {
	return v.Constraint.Name + "\x00" + engine.BindingsText(v.Bindings)
}

func (r reporter) diff(w model.World, a, b engine.Result) error {
	type change struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
	}
	factChanges := func(asserted bool) change {
		c := change{Added: []string{}, Removed: []string{}}
		pick := func(res engine.Result) map[string]string {
			out := map[string]string{}
			for _, f := range res.Graph.Order {
				k := engine.Key(f)
				if (res.Graph.Reasons[k].Type == engine.ReasonAsserted) == asserted {
					// A field's key names the field, not its value; a changed
					// value is a removal and an addition.
					out[k+"\x00"+value.Format(f.Value)] = engine.Sentence(w, f)
				}
			}
			return out
		}
		before, after := pick(a), pick(b)
		for k, s := range after {
			if _, ok := before[k]; !ok {
				c.Added = append(c.Added, s)
			}
		}
		for k, s := range before {
			if _, ok := after[k]; !ok {
				c.Removed = append(c.Removed, s)
			}
		}
		slices.Sort(c.Added)
		slices.Sort(c.Removed)
		return c
	}
	rejected := change{Added: []string{}, Removed: []string{}}
	vs := func(res engine.Result) map[string]string {
		out := map[string]string{}
		for _, v := range res.Violations {
			out[violationKey(v)] = v.Constraint.Name + ": " + v.Message
		}
		return out
	}
	before, after := vs(a), vs(b)
	for k, s := range after {
		if _, ok := before[k]; !ok {
			rejected.Added = append(rejected.Added, s)
		}
	}
	for k, s := range before {
		if _, ok := after[k]; !ok {
			rejected.Removed = append(rejected.Removed, s)
		}
	}
	slices.Sort(rejected.Added)
	slices.Sort(rejected.Removed)
	out := struct {
		Asserted change `json:"asserted"`
		Derived  change `json:"derived"`
		Rejected change `json:"rejected"`
	}{factChanges(true), factChanges(false), rejected}
	if r.json {
		return writeJSON(r.out, out)
	}
	var lines []string
	for _, sec := range []struct {
		name string
		c    change
	}{{"asserted", out.Asserted}, {"derived", out.Derived}, {"rejected", out.Rejected}} {
		lines = append(lines, sec.name+":")
		if len(sec.c.Added)+len(sec.c.Removed) == 0 {
			lines = append(lines, "  no change")
		}
		for _, s := range sec.c.Added {
			lines = append(lines, "  + "+s)
		}
		for _, s := range sec.c.Removed {
			lines = append(lines, "  - "+s)
		}
	}
	return writeLines(r.out, lines)
}

func (r reporter) tests(outcomes []suite.Outcome) error {
	if r.json {
		if outcomes == nil {
			outcomes = []suite.Outcome{}
		}
		return writeJSON(r.out, map[string]any{"tests": outcomes})
	}
	var lines []string
	failed := 0
	for _, o := range outcomes {
		status := "pass"
		if !suite.Passed(o) {
			status = "FAIL"
			failed++
		}
		lines = append(lines, fmt.Sprintf("%s  %s %s  (%s)", status, strings.ToLower(o.Sort), o.Name, o.File))
		for _, f := range o.Failures {
			lines = append(lines, "      "+f)
		}
	}
	lines = append(lines, fmt.Sprintf("%d passed, %d failed", len(outcomes)-failed, failed))
	return writeLines(r.out, lines)
}
