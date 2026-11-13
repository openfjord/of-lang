package engine

import (
	"maps"
	"slices"
	"strings"

	"oflang/internal/model"
	"oflang/internal/value"
	"oflang/internal/world"
)

// Match is one way a pattern fits the facts: what each variable stands for,
// and the facts that made it fit.
type Match struct {
	Env      map[string]string
	Premises []Fact
}

// Find returns every way the links and conditions fit the graph. vars are
// further variables the caller needs bound; any not bound by a link ranges
// over every thing.
func Find(w model.World, g Graph, links []model.PatternLink, where []model.Condition, vars []string, start map[string]string) []Match {
	var out []Match
	for _, m := range matchLinks(w, g, links, Match{Env: start}) {
		for _, m := range bindRest(g, m, vars) {
			if fm, ok := allHold(w, g, where, m); ok {
				out = append(out, fm)
			}
		}
	}
	return out
}

func bindRest(g Graph, m Match, vars []string) []Match {
	for _, v := range vars {
		if _, ok := m.Env[v]; ok {
			continue
		}
		var out []Match
		for _, id := range g.Nodes {
			env, _ := bind(m.Env, v, id)
			out = append(out, bindRest(g, Match{Env: env, Premises: m.Premises}, vars)...)
		}
		return out
	}
	return []Match{m}
}

func allHold(w model.World, g Graph, conds []model.Condition, m Match) (Match, bool) {
	premises := slices.Clone(m.Premises)
	for _, c := range conds {
		ok, ps := holds(w, g, c, m.Env)
		if !ok {
			return Match{}, false
		}
		premises = append(premises, ps...)
	}
	return Match{Env: m.Env, Premises: premises}, true
}

func matchLinks(w model.World, g Graph, links []model.PatternLink, m Match) []Match {
	if len(links) == 0 {
		return []Match{m}
	}
	i := mostBound(links, m.Env)
	pl := links[i]
	rest := slices.Delete(slices.Clone(links), i, i+1)
	var out []Match
	for _, step := range stepMatches(g, pl, m) {
		out = append(out, matchLinks(w, g, rest, step)...)
	}
	return out
}

func mostBound(links []model.PatternLink, env map[string]string) int {
	best, bestScore := 0, -1
	for i, pl := range links {
		score := 0
		if _, ok := env[pl.From]; ok {
			score++
		}
		if _, ok := env[pl.To]; ok {
			score++
		}
		if score > bestScore {
			best, bestScore = i, score
		}
	}
	return best
}

// stepMatches extends a match by one pattern link.
func stepMatches(g Graph, pl model.PatternLink, m Match) []Match {
	fromVar, toVar := pl.From, pl.To
	if pl.Inverted {
		fromVar, toVar = toVar, fromVar
	}
	from, fromBound := m.Env[fromVar]
	to, toBound := m.Env[toVar]
	var paths [][]Edge
	switch {
	case pl.Repeated:
		paths = repeatedPaths(g, pl.Link, from, fromBound, to, toBound)
	case fromBound && toBound:
		if g.Edges[Edge{From: from, Link: pl.Link, To: to}] {
			paths = append(paths, []Edge{{From: from, Link: pl.Link, To: to}})
		}
	case fromBound:
		for _, t := range g.out[from][pl.Link] {
			paths = append(paths, []Edge{{From: from, Link: pl.Link, To: t}})
		}
	case toBound:
		for _, f := range g.in[to][pl.Link] {
			paths = append(paths, []Edge{{From: f, Link: pl.Link, To: to}})
		}
	default:
		for _, e := range g.byLink[pl.Link] {
			paths = append(paths, []Edge{e})
		}
	}
	var out []Match
	for _, path := range paths {
		env, ok := bind(m.Env, fromVar, path[0].From)
		if !ok {
			continue
		}
		env, ok = bind(env, toVar, path[len(path)-1].To)
		if !ok {
			continue
		}
		premises := slices.Clone(m.Premises)
		for _, e := range path {
			premises = append(premises, LinkFact(e.From, e.Link, e.To))
		}
		out = append(out, Match{Env: env, Premises: premises})
	}
	return out
}

func bind(env map[string]string, v, id string) (map[string]string, bool) {
	if cur, ok := env[v]; ok {
		return env, cur == id
	}
	out := maps.Clone(env)
	if out == nil {
		out = map[string]string{}
	}
	out[v] = id
	return out, true
}

// repeatedPaths finds paths of one or more steps along a link.
func repeatedPaths(g Graph, link model.QName, from string, fromBound bool, to string, toBound bool) [][]Edge {
	var starts []string
	switch {
	case fromBound:
		starts = []string{from}
	case toBound:
		return reversePaths(walk(g.in, link, to))
	default:
		for _, e := range g.byLink[link] {
			if !slices.Contains(starts, e.From) {
				starts = append(starts, e.From)
			}
		}
	}
	var out [][]Edge
	for _, s := range starts {
		for _, p := range walk(g.out, link, s) {
			if !toBound || p[len(p)-1].To == to {
				out = append(out, p)
			}
		}
	}
	return out
}

// walk returns a shortest path from start to every thing reachable along a
// link, following the given index.
func walk(index map[string]map[model.QName][]string, link model.QName, start string) [][]Edge {
	paths := map[string][]Edge{}
	queue := []string{start}
	var order []string
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range index[cur][link] {
			if _, seen := paths[next]; seen {
				continue
			}
			paths[next] = append(slices.Clone(paths[cur]), Edge{From: cur, Link: link, To: next})
			order = append(order, next)
			queue = append(queue, next)
		}
	}
	out := make([][]Edge, 0, len(order))
	for _, n := range order {
		out = append(out, paths[n])
	}
	return out
}

func reversePaths(paths [][]Edge) [][]Edge {
	out := make([][]Edge, 0, len(paths))
	for _, p := range paths {
		r := make([]Edge, 0, len(p))
		for i := len(p) - 1; i >= 0; i-- {
			r = append(r, Edge{From: p[i].To, Link: p[i].Link, To: p[i].From})
		}
		out = append(out, r)
	}
	return out
}

// holds evaluates one condition and returns the facts it relied on.
func holds(w model.World, g Graph, c model.Condition, env map[string]string) (bool, []Fact) {
	switch c.Type {
	case model.CondCompare:
		left, lp, ok := nodeValue(g, env, c.Node, c.Field)
		if !ok {
			return false, nil
		}
		rights, rp, ok := operandValues(g, env, c.Than)
		if !ok {
			return false, nil
		}
		if !compare(c.Op, left, rights) {
			return false, nil
		}
		return true, append(lp, rp...)
	case model.CondIsA:
		id := env[c.Node]
		for _, k := range g.Kinds[id] {
			if world.IsA(w, k, c.Kind) {
				return true, []Fact{KindFact(id, k)}
			}
		}
		return false, nil
	case model.CondIsSet:
		v, ok := g.Fields[env[c.Node]][c.Field]
		if !ok {
			return false, nil
		}
		return true, []Fact{FieldFact(env[c.Node], c.Field, v)}
	case model.CondIsNotSet:
		_, ok := g.Fields[env[c.Node]][c.Field]
		return !ok, nil
	case model.CondAtLeastOne:
		ms := nested(w, g, c, env)
		if len(ms) == 0 {
			return false, nil
		}
		return true, ms[0].Premises
	case model.CondNone:
		return len(nested(w, g, c, env)) == 0, nil
	case model.CondCount:
		seen := map[string]bool{}
		for _, m := range nested(w, g, c, env) {
			seen[envKey(m.Env)] = true
		}
		rights, _, ok := operandValues(g, env, c.Than)
		return ok && compare(c.Op, value.Integer(len(seen)), rights), nil
	case model.CondAnyOf:
		for _, a := range c.Any {
			if ok, ps := holds(w, g, a, env); ok {
				return true, ps
			}
		}
		return false, nil
	default:
		return false, nil
	}
}

func nested(w model.World, g Graph, c model.Condition, env map[string]string) []Match {
	return Find(w, g, c.Links, c.Where, conditionVars(c.Where), maps.Clone(env))
}

func envKey(env map[string]string) string {
	keys := slices.Sorted(maps.Keys(env))
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k + "=" + env[k] + "\x00")
	}
	return b.String()
}

func nodeValue(g Graph, env map[string]string, node, field string) (model.Value, []Fact, bool) {
	id, ok := env[node]
	if !ok {
		return model.Value{}, nil, false
	}
	if field == "" {
		return value.Text(id), nil, true
	}
	v, ok := g.Fields[id][field]
	if !ok {
		return model.Value{}, nil, false
	}
	return v, []Fact{FieldFact(id, field, v)}, true
}

func operandValues(g Graph, env map[string]string, o model.Operand) ([]model.Value, []Fact, bool) {
	switch {
	case o.Value != nil:
		return []model.Value{*o.Value}, nil, true
	case o.Node != "":
		v, ps, ok := nodeValue(g, env, o.Node, o.Field)
		return []model.Value{v}, ps, ok
	default:
		return o.Values, nil, true
	}
}

func compare(op model.Op, left model.Value, rights []model.Value) bool {
	if op == model.OpIsOneOf {
		return slices.ContainsFunc(rights, func(r model.Value) bool { return value.Equal(left, r) })
	}
	if len(rights) != 1 {
		return false
	}
	right := rights[0]
	switch op {
	case model.OpEquals:
		return value.Equal(left, right)
	case model.OpDoesNotEqual:
		return !value.Equal(left, right)
	default:
		c, ok := value.Compare(left, right)
		if !ok {
			return false
		}
		switch op {
		case model.OpIsAbove:
			return c > 0
		case model.OpIsBelow:
			return c < 0
		case model.OpIsAtLeast:
			return c >= 0
		case model.OpIsAtMost:
			return c <= 0
		default:
			return false
		}
	}
}

// conditionVars lists the variables conditions use outside their nested
// patterns.
func conditionVars(conds []model.Condition) []string {
	var out []string
	add := func(v string) {
		if v != "" && !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	for _, c := range conds {
		add(c.Node)
		add(c.Than.Node)
		for _, v := range conditionVars(c.Any) {
			add(v)
		}
	}
	return out
}

// OuterVars lists the variables a statement binds for each match: those in
// its links, its conditions outside nested patterns, its derivations and its
// message. Variables that appear only inside a nested pattern belong to it.
func OuterVars(st model.Statement) []string {
	var out []string
	add := func(v string) {
		if v != "" && !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	for _, pl := range st.Links {
		add(pl.From)
		add(pl.To)
	}
	for _, v := range conditionVars(st.Where) {
		add(v)
	}
	for _, v := range conditionVars(st.Require) {
		add(v)
	}
	for _, d := range st.Derive {
		add(d.From)
		add(d.To)
		add(d.Node)
		add(d.SetTo.Node)
	}
	for _, m := range messageVar.FindAllStringSubmatch(st.Message, -1) {
		add(m[1])
	}
	return out
}
