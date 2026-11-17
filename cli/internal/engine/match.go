// SPDX-License-Identifier: Apache-2.0

package engine

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/value"
	"github.com/openfjord/of-lang/cli/internal/world"
)

// Match is one row of a match: what each variable stands for, and the facts
// that made it fit.
type Match struct {
	Env      map[string]string
	Premises []Fact
}

// item is a link item with a variable at each end. An end with a kind alone
// gets a variable of its own, starting with #, that no one else can name.
type item struct {
	from, to string
	link     model.QName
	inverted bool
	repeated bool
}

// anonymous says whether a variable stands for an end with a kind alone.
func anonymous(v string) bool { return strings.HasPrefix(v, "#") }

// Find returns every row of a block: each combination of things for which
// its items hold, then its filter. start holds what an enclosing block
// already bound. Rows that give the same things to the same variables are one
// row; withAnonymous counts the ends with a kind alone as variables too, as a
// nested block does.
func Find(w model.World, g Graph, b model.Block, start map[string]string, withAnonymous bool) []Match {
	kinds := map[string]model.QName{}
	for _, ki := range b.Kinds {
		kinds[ki.Var] = ki.Kind
	}
	items := make([]item, 0, len(b.Links))
	for i, li := range b.Links {
		it := item{link: li.Link, inverted: li.Inverted, repeated: li.Repeated}
		it.from = endVar(li.From, fmt.Sprintf("#%d.from", i), kinds)
		it.to = endVar(li.To, fmt.Sprintf("#%d.to", i), kinds)
		items = append(items, it)
	}
	var out []Match
	seen := map[string]bool{}
	for _, m := range matchItems(w, g, items, kinds, Match{Env: maps.Clone(start)}) {
		for _, m := range bindRest(w, g, m, b.Kinds) {
			fm, ok := allHold(w, g, b.Filter, m)
			if !ok {
				continue
			}
			k := rowKey(fm.Env, withAnonymous)
			if seen[k] {
				continue
			}
			seen[k] = true
			out = append(out, fm)
		}
	}
	return out
}

func endVar(e model.End, anon string, kinds map[string]model.QName) string {
	if e.Var != "" {
		return e.Var
	}
	kinds[anon] = e.Kind
	return anon
}

func rowKey(env map[string]string, withAnonymous bool) string {
	var b strings.Builder
	for _, k := range slices.Sorted(maps.Keys(env)) {
		if anonymous(k) && !withAnonymous {
			continue
		}
		b.WriteString(k + "=" + env[k] + "\x00")
	}
	return b.String()
}

// thingIsA says whether a thing has a kind, stated or derived, that is the
// wanted kind.
func thingIsA(w model.World, g Graph, id string, want model.QName) bool {
	return slices.ContainsFunc(g.Kinds[id], func(k model.QName) bool { return world.IsA(w, k, want) })
}

// kindPremise returns the kind fact that makes a thing a thing of a kind.
func kindPremise(w model.World, g Graph, id string, want model.QName) Fact {
	for _, k := range g.Kinds[id] {
		if world.IsA(w, k, want) {
			return KindFact(id, k)
		}
	}
	return KindFact(id, want)
}

func bindRest(w model.World, g Graph, m Match, kis []model.KindItem) []Match {
	for _, ki := range kis {
		if _, ok := m.Env[ki.Var]; ok {
			continue
		}
		var out []Match
		for _, id := range g.Nodes {
			if !thingIsA(w, g, id, ki.Kind) {
				continue
			}
			env, _ := bind(m.Env, ki.Var, id)
			out = append(out, bindRest(w, g, Match{Env: env, Premises: m.Premises}, kis)...)
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

func matchItems(w model.World, g Graph, items []item, kinds map[string]model.QName, m Match) []Match {
	if len(items) == 0 {
		return []Match{m}
	}
	i := mostBound(items, m.Env)
	rest := slices.Delete(slices.Clone(items), i, i+1)
	var out []Match
	for _, step := range stepMatches(w, g, items[i], kinds, m) {
		out = append(out, matchItems(w, g, rest, kinds, step)...)
	}
	return out
}

func mostBound(items []item, env map[string]string) int {
	best, bestScore := 0, -1
	for i, it := range items {
		score := 0
		if _, ok := env[it.from]; ok {
			score++
		}
		if _, ok := env[it.to]; ok {
			score++
		}
		if score > bestScore {
			best, bestScore = i, score
		}
	}
	return best
}

// stepMatches extends a row by one link item.
func stepMatches(w model.World, g Graph, it item, kinds map[string]model.QName, m Match) []Match {
	fromVar, toVar := it.from, it.to
	if it.inverted {
		fromVar, toVar = toVar, fromVar
	}
	from, fromBound := m.Env[fromVar]
	to, toBound := m.Env[toVar]
	var paths [][]Edge
	switch {
	case it.repeated:
		paths = repeatedPaths(g, it.link, from, fromBound, to, toBound)
	case fromBound && toBound:
		if g.Edges[Edge{From: from, Link: it.link, To: to}] {
			paths = append(paths, []Edge{{From: from, Link: it.link, To: to}})
		}
	case fromBound:
		for _, t := range g.out[from][it.link] {
			paths = append(paths, []Edge{{From: from, Link: it.link, To: t}})
		}
	case toBound:
		for _, f := range g.in[to][it.link] {
			paths = append(paths, []Edge{{From: f, Link: it.link, To: to}})
		}
	default:
		for _, e := range g.byLink[it.link] {
			paths = append(paths, []Edge{e})
		}
	}
	var out []Match
	for _, path := range paths {
		first, last := path[0].From, path[len(path)-1].To
		if !fits(w, g, kinds, fromVar, first) || !fits(w, g, kinds, toVar, last) {
			continue
		}
		env, ok := bind(m.Env, fromVar, first)
		if !ok {
			continue
		}
		env, ok = bind(env, toVar, last)
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

// fits says whether a thing can stand for a variable: it has the kind of the
// variable's kind item, or of the end's kind alone.
func fits(w model.World, g Graph, kinds map[string]model.QName, v, id string) bool {
	k, ok := kinds[v]
	return !ok || thingIsA(w, g, id, k)
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
		left, lp, ok := nodeValue(g, env, c.Node, c.Field, c.Part)
		if !ok {
			return false, nil
		}
		rights, rp, ok := operandValues(g, env, c.Than)
		if !ok || !compare(c.Op, left, rights) {
			return false, nil
		}
		return true, append(lp, rp...)
	case model.CondIsA:
		id := env[c.Node]
		if !thingIsA(w, g, id, c.Kind) {
			return false, nil
		}
		return true, []Fact{kindPremise(w, g, id, c.Kind)}
	case model.CondIsSet:
		_, ps, ok := nodeValue(g, env, c.Node, c.Field, c.Part)
		return ok, ps
	case model.CondIsNotSet:
		_, _, ok := nodeValue(g, env, c.Node, c.Field, c.Part)
		return !ok, nil
	case model.CondAtLeastOne:
		ms := Find(w, g, *c.Block, env, true)
		if len(ms) == 0 {
			return false, nil
		}
		return true, ms[0].Premises
	case model.CondNone:
		return len(Find(w, g, *c.Block, env, true)) == 0, nil
	case model.CondCount:
		n := len(Find(w, g, *c.Block, env, true))
		rights, _, ok := operandValues(g, env, c.Than)
		return ok && compare(c.Op, value.Integer(n), rights), nil
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

// nodeValue returns the value a condition reads: the id of a thing, or the
// value of one of its fields, or of a part of that value. ok is false when
// the field or the part has no value.
func nodeValue(g Graph, env map[string]string, node, field string, part []string) (model.Value, []Fact, bool) {
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
	premise := FieldFact(id, field, v)
	for _, p := range part {
		next, ok := v.Props[p]
		if v.Kind != model.ValueObject || !ok {
			return model.Value{}, nil, false
		}
		v = next
	}
	return v, []Fact{premise}, true
}

func operandValues(g Graph, env map[string]string, o model.Operand) ([]model.Value, []Fact, bool) {
	switch {
	case o.Value != nil:
		return []model.Value{*o.Value}, nil, true
	case o.Node != "":
		v, ps, ok := nodeValue(g, env, o.Node, o.Field, nil)
		return []model.Value{v}, ps, ok
	default:
		return o.Values, nil, true
	}
}

// compare applies a comparison word. A comparison holds only when the two
// values have the same type.
func compare(op model.Op, left model.Value, rights []model.Value) bool {
	if op == model.OpIsOneOf {
		return slices.ContainsFunc(rights, func(r model.Value) bool { return value.Equal(left, r) })
	}
	if len(rights) != 1 || !value.SameType(left, rights[0]) {
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

// BlockVars lists the variables of a block's own kind items.
func BlockVars(b model.Block) []string {
	out := make([]string, 0, len(b.Kinds))
	for _, ki := range b.Kinds {
		out = append(out, ki.Var)
	}
	return out
}
