package engine

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"oflang/internal/diag"
	"oflang/internal/model"
	"oflang/internal/value"
	"oflang/internal/world"
)

var messageVar = regexp.MustCompile(`\{(\$[a-z][A-Za-z0-9]*)\}`)

// Violation is one match that breaks a constraint.
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
	// Problems are facts that break what the kinds declare: link counts, the
	// kinds a link joins, and field values a rule set.
	Problems []diag.Diagnostic
}

// Run applies the built-in inference and the rules, stratum by stratum, until
// nothing new appears; then it checks constraints and link counts.
func Run(w model.World, strata [][]model.QName, g Graph) Result {
	var problems []diag.Diagnostic
	builtins(w, &g)
	for _, stratum := range strata {
		for {
			changed := false
			for _, q := range stratum {
				added, ps := applyRule(w, &g, w.Statements[q])
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
	res := Result{Graph: g, Problems: problems}
	for _, q := range w.Order {
		st := w.Statements[q]
		if st.Sort == model.SortConstraint {
			res.Violations = append(res.Violations, violations(w, g, st)...)
		}
	}
	res.Problems = diag.Sorted(append(res.Problems, kindProblems(w, g)...))
	return res
}

// builtins adds what links' under and chains declarations imply, until
// nothing new appears.
func builtins(w model.World, g *Graph) bool {
	grew := false
	for {
		changed := false
		edges := slices.Clone(g.Order)
		for _, f := range edges {
			if f.Type != FactLink {
				continue
			}
			e := f.Edge
			link := w.Links[e.Link]
			if link.Under != nil {
				if Add(g, LinkFact(e.From, *link.Under, e.To), Reason{Type: ReasonUnder, Premises: []Fact{f}}) {
					changed = true
				}
			}
			if link.Chains {
				for _, next := range slices.Clone(g.out[e.To][e.Link]) {
					if Add(g, LinkFact(e.From, e.Link, next), Reason{Type: ReasonChains, Premises: []Fact{f, LinkFact(e.To, e.Link, next)}}) {
						changed = true
					}
				}
			}
		}
		if !changed {
			return grew
		}
		grew = true
	}
}

func applyRule(w model.World, g *Graph, st model.Statement) (bool, []diag.Diagnostic) {
	added := false
	var problems []diag.Diagnostic
	vars := OuterVars(st)
	for _, m := range Find(w, *g, st.Links, st.Where, vars, nil) {
		reason := Reason{Type: ReasonRule, Rule: st.Q, Bindings: sortedBindings(m.Env, vars), Premises: m.Premises}
		for _, d := range st.Derive {
			f, ok := derived(*g, d, m.Env)
			if !ok {
				continue
			}
			if f.Type == FactField {
				if cur, set := g.Fields[f.Node][f.Field]; set {
					if !value.Equal(cur, f.Value) {
						problems = append(problems, diag.New(d.Pos, "rule %s sets %s of %s to %s, but it is already %s", st.Q.Name, f.Field, f.Node, value.Format(f.Value), value.Format(cur)))
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
	vars := OuterVars(st)
	for _, m := range Find(w, g, st.Links, st.Where, vars, nil) {
		if _, ok := allHold(w, g, st.Require, m); ok {
			continue
		}
		out = append(out, Violation{
			Constraint: st.Q,
			Pos:        st.Pos,
			Bindings:   sortedBindings(m.Env, vars),
			Message:    Render(st.Message, m.Env),
			Premises:   m.Premises,
		})
	}
	return out
}

// Render fills a message's {$var} placeholders with what the variables stand
// for.
func Render(msg string, env map[string]string) string {
	return messageVar.ReplaceAllStringFunc(msg, func(s string) string {
		v := messageVar.FindStringSubmatch(s)[1]
		if id, ok := env[v]; ok {
			return id
		}
		return s
	})
}

// kindProblems checks every thing against its kinds: how many of each link it
// has, and that each link joins the kinds it declares.
func kindProblems(w model.World, g Graph) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, id := range g.Nodes {
		seen := map[diag.Pos]bool{}
		for _, k := range g.Kinds[id] {
			for _, lc := range world.LinkCounts(w, k) {
				if seen[lc.Pos] {
					continue
				}
				seen[lc.Pos] = true
				n := len(g.out[id][lc.Ref.Q])
				if lc.Ref.Inverted {
					n = len(g.in[id][lc.Ref.Q])
				}
				if n < lc.Count.Min || lc.Count.Max >= 0 && n > lc.Count.Max {
					out = append(out, diag.New(g.Pos[id], "%s has %d %s, but a %s has %s", id, n, lc.Name, k.Name, countText(lc.Count)))
				}
			}
		}
	}
	for _, f := range g.Order {
		if f.Type != FactLink {
			continue
		}
		e := f.Edge
		link := w.Links[e.Link]
		for _, end := range []struct {
			id, role string
			want     model.QName
		}{{e.From, "from", link.From}, {e.To, "to", link.To}} {
			if !isA(w, g, end.id, end.want) {
				out = append(out, diag.New(g.Pos[end.id], "%s links %s to %s, but %s is not a %s", e.Link.Name, e.From, e.To, end.id, end.want.Name))
			}
		}
	}
	return out
}

func isA(w model.World, g Graph, id string, want model.QName) bool {
	return slices.ContainsFunc(g.Kinds[id], func(k model.QName) bool { return world.IsA(w, k, want) })
}

func countText(c model.Count) string {
	switch {
	case c.Max < 0 && c.Min == 0:
		return "any number"
	case c.Max < 0:
		return fmt.Sprintf("at least %d", c.Min)
	case c.Min == c.Max:
		return fmt.Sprintf("exactly %d", c.Min)
	default:
		return fmt.Sprintf("at most %d", c.Max)
	}
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
