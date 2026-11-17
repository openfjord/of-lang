// SPDX-License-Identifier: Apache-2.0

package load

import (
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
)

// decoder reads the parts that constraints, derivatives and questions share,
// naming the rules of the sort it reads.
type decoder struct {
	path  string
	scope map[string]model.Ref
	rules model.Rules
	ds    *[]diag.Diagnostic
}

// Statement decodes a constraint, a derivative or a question. A question's
// sentence is asks; the caller reads its given and expect.
func Statement(q model.QName, sort model.Sort, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.Statement {
	d := decoder{path: path, scope: scope, rules: model.RulesFor(sort), ds: ds}
	sentenceKey := "says"
	var allowed []string
	switch sort {
	case model.SortConstraint:
		allowed = []string{"is", "tags", "says", "match", "filter", "unless", "message"}
		if p, ok := pairOf(m, "derive"); ok {
			*ds = append(*ds, diag.New(p.Pos, "a constraint adds no facts, so it has no derive (CN13)"))
		}
	case model.SortDerivative:
		allowed = []string{"is", "tags", "says", "match", "filter", "derive"}
		for _, key := range []string{"unless", "message"} {
			if p, ok := pairOf(m, key); ok {
				*ds = append(*ds, diag.New(p.Pos, "a derivative rejects nothing, so it has no %s (DV16)", key))
			}
		}
	default:
		sentenceKey = "asks"
		allowed = []string{"is", "asks", "given", "match", "filter", "expect"}
	}
	for _, p := range m.Pairs {
		if !slices.Contains(allowed, p.Key) && !slices.Contains([]string{"derive", "unless", "message"}, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not allowed here; allowed: %v (FL2)", p.Key, allowed))
		}
	}
	st := model.Statement{Q: q, Sort: sort, Pos: m.Pos}
	st.Says = strings.TrimSpace(requiredText(path, m, sentenceKey, q.Name, d.rules.Sentence, ds))
	if v := varPattern.FindString(st.Says); v != "" {
		*ds = append(*ds, diag.New(m.Pos, "%s names %s; the sentence is plain words, and variables belong to the logic (%s)", sentenceKey, v, d.rules.Sentence))
	}
	if sort != model.SortQuestion {
		st.Tags = decodeTags(path, m, ds)
	}
	st.Match = d.block(m, m.Pos)
	if n, ok := source.Get(m, "unless"); ok && sort == model.SortConstraint {
		st.Unless = d.conditions(n)
	}
	if sort == model.SortConstraint {
		st.Message = requiredText(path, m, "message", q.Name, "CN12", ds)
	}
	if sort == model.SortDerivative {
		n, ok := source.Get(m, "derive")
		if !ok {
			*ds = append(*ds, diag.New(m.Pos, "derivative %s needs derive: (DV12)", q.Name))
		} else {
			st.Derive = d.derivations(n)
			if len(st.Derive) == 0 {
				*ds = append(*ds, diag.New(m.Pos, "derive of %s needs at least one item (DV12)", q.Name))
			}
		}
	}
	return st
}

func pairOf(m source.Map, key string) (source.Pair, bool) {
	for _, p := range m.Pairs {
		if p.Key == key {
			return p, true
		}
	}
	return source.Pair{}, false
}

func (d decoder) diag(pos diag.Pos, rule, format string, args ...any) {
	*d.ds = append(*d.ds, diag.New(pos, format+" ("+rule+")", args...))
}

// block reads the match and filter keys of a mapping: a statement or a
// nested block.
func (d decoder) block(m source.Map, pos diag.Pos) model.Block {
	b := model.Block{Pos: pos}
	n, ok := source.Get(m, "match")
	if !ok {
		d.diag(pos, d.rules.Match, "match: is missing")
	} else {
		d.matchItems(n, &b)
		if len(b.Kinds)+len(b.Links) == 0 {
			d.diag(source.PosOf(d.path, n), d.rules.Match, "match needs at least one item")
		}
	}
	if fn, ok := source.Get(m, "filter"); ok {
		b.Filter = d.conditions(fn)
	}
	return b
}

func (d decoder) matchItems(n *yaml.Node, b *model.Block) {
	items, ok := sequence(d.path, n, d.ds)
	if !ok {
		return
	}
	for _, item := range items {
		m, ok := source.AsMap(d.path, item, d.ds)
		if !ok {
			continue
		}
		switch {
		case has(m, "kind"):
			if ki, ok := d.kindItem(m); ok {
				b.Kinds = append(b.Kinds, ki)
			}
		case has(m, "link"):
			if li, ok := d.linkItem(m); ok {
				b.Links = append(b.Links, li)
			}
		default:
			d.diag(m.Pos, d.rules.Items, "a match item is a kind item (kind:) or a link item (link:)")
		}
	}
}

func (d decoder) kindItem(m source.Map) (model.KindItem, bool) {
	allow(m, []string{"kind"}, d.rules.Items, d.ds)
	n, _ := source.Get(m, "kind")
	km, ok := source.AsMap(d.path, n, d.ds)
	if !ok {
		return model.KindItem{}, false
	}
	allow(km, []string{"name", "as"}, d.rules.Items, d.ds)
	name, okName := source.GetScalar(km, "name", d.ds)
	as, okAs := source.GetScalar(km, "as", d.ds)
	if !okName || !okAs {
		d.diag(km.Pos, d.rules.Items, "a kind item needs name: and as:")
		return model.KindItem{}, false
	}
	if !IsVariable(as) {
		d.diag(km.Pos, d.rules.Items, "as: must be a variable: $, a lowercase letter, then letters and digits; not %q", as)
		return model.KindItem{}, false
	}
	ref, ok := resolve(d.scope, name, model.SortKind, km.Pos, d.rules.Items, d.ds)
	if !ok {
		return model.KindItem{}, false
	}
	return model.KindItem{Pos: m.Pos, Var: as, Kind: ref.Q}, true
}

func (d decoder) linkItem(m source.Map) (model.LinkItem, bool) {
	allow(m, []string{"link", "from", "to", "repeated"}, d.rules.Items, d.ds)
	li := model.LinkItem{Pos: m.Pos}
	name, ok := source.GetScalar(m, "link", d.ds)
	if !ok {
		return li, false
	}
	ref, ok := resolve(d.scope, name, model.SortLink, m.Pos, d.rules.Items, d.ds)
	if !ok {
		return li, false
	}
	li.Name, li.Link, li.Inverted = name, ref.Q, ref.Inverted
	from, okFrom := d.end(m, "from")
	to, okTo := d.end(m, "to")
	if !okFrom || !okTo {
		return li, false
	}
	li.From, li.To = from, to
	if rn, ok := source.Get(m, "repeated"); ok {
		li.Repeated, _ = source.Bool(d.path, rn, d.ds)
	}
	return li, true
}

func (d decoder) end(m source.Map, key string) (model.End, bool) {
	n, ok := source.Get(m, key)
	if !ok {
		d.diag(m.Pos, d.rules.Items, "a link item needs %s:", key)
		return model.End{}, false
	}
	s, ok := source.Scalar(d.path, n, d.ds)
	if !ok {
		return model.End{}, false
	}
	pos := source.PosOf(d.path, n)
	if strings.HasPrefix(s, "$") {
		if !IsVariable(s) {
			d.diag(pos, d.rules.Ends, "%q is not a variable: $, a lowercase letter, then letters and digits", s)
			return model.End{}, false
		}
		return model.End{Pos: pos, Var: s}, true
	}
	ref, ok := resolve(d.scope, s, model.SortKind, pos, d.rules.Ends, d.ds)
	if !ok {
		return model.End{}, false
	}
	return model.End{Pos: pos, Kind: ref.Q}, true
}

// IsVariable says whether a word is a variable: $, a lowercase letter, then
// letters and digits.
func IsVariable(s string) bool {
	return s != "" && varPattern.FindString(s) == s
}

func (d decoder) conditions(n *yaml.Node) []model.Condition {
	items, ok := sequence(d.path, n, d.ds)
	if !ok {
		return nil
	}
	var out []model.Condition
	for _, item := range items {
		if c, ok := d.condition(item); ok {
			out = append(out, c)
		}
	}
	return out
}

var nestedWords = []model.CondType{model.CondAtLeastOne, model.CondNone, model.CondCount}

func (d decoder) condition(n *yaml.Node) (model.Condition, bool) {
	m, ok := source.AsMap(d.path, n, d.ds)
	if !ok {
		return model.Condition{}, false
	}
	c := model.Condition{Pos: m.Pos}
	if an, ok := source.Get(m, "any of"); ok {
		allow(m, []string{"any of"}, d.rules.Conds, d.ds)
		c.Type = model.CondAnyOf
		c.Any = d.conditions(an)
		return c, true
	}
	for _, nested := range nestedWords {
		nn, ok := source.Get(m, string(nested))
		if !ok {
			continue
		}
		c.Type = nested
		nm, ok := source.AsMap(d.path, nn, d.ds)
		if !ok {
			return c, false
		}
		allow(nm, []string{"match", "filter"}, d.rules.Conds, d.ds)
		b := d.block(nm, nm.Pos)
		c.Block = &b
		if nested != model.CondCount {
			allow(m, []string{string(nested)}, d.rules.Conds, d.ds)
			return c, true
		}
		allow(m, append([]string{"count"}, countOps()...), d.rules.Count, d.ds)
		op, than, ok := d.comparison(m)
		switch {
		case !ok:
			d.diag(m.Pos, d.rules.Count, "count needs one comparison word: %s", strings.Join(countOps(), ", "))
			return c, false
		case op == model.OpIsOneOf:
			d.diag(m.Pos, d.rules.Count, "count takes %s, not is one of", strings.Join(countOps(), ", "))
			return c, false
		case than.Value == nil || than.Value.Kind != model.ValueInteger || than.Value.Num.Sign() < 0:
			d.diag(m.Pos, d.rules.Count, "count compares with value: and a whole number, 0 or more")
			return c, false
		}
		c.Op, c.Than = op, than
		return c, true
	}
	node, ok := source.GetScalar(m, "node", d.ds)
	if !ok {
		d.diag(m.Pos, d.rules.Conds, "a condition starts with node:, any of:, at least one:, none: or count:")
		return c, false
	}
	if !IsVariable(node) {
		d.diag(m.Pos, d.rules.Vars, "node: names a variable alone, not %q", node)
		return c, false
	}
	c.Node = node
	c.Field, _ = source.GetScalar(m, "field", d.ds)
	if pn, ok := source.Get(m, "part"); ok {
		c.Part, _ = source.Strings(d.path, pn, d.ds)
		if c.Field == "" {
			d.diag(source.PosOf(d.path, pn), d.rules.Parts, "part goes only next to field")
		}
	}
	allow(m, append([]string{"node", "field", "part", "is a", "is set", "is not set"}, opNames()...), d.rules.Conds, d.ds)
	if kn, ok := source.Get(m, "is a"); ok {
		c.Type = model.CondIsA
		name, ok := source.Scalar(d.path, kn, d.ds)
		if !ok {
			return c, false
		}
		ref, ok := resolve(d.scope, name, model.SortKind, source.PosOf(d.path, kn), d.rules.Conds, d.ds)
		if !ok {
			return c, false
		}
		c.Kind = ref.Q
		return c, true
	}
	for _, t := range []model.CondType{model.CondIsSet, model.CondIsNotSet} {
		if _, ok := source.Get(m, string(t)); ok {
			if c.Field == "" {
				d.diag(m.Pos, d.rules.Conds, "%s needs field:", t)
			}
			c.Type = t
			return c, true
		}
	}
	op, than, ok := d.comparison(m)
	if !ok {
		d.diag(m.Pos, d.rules.Conds, "a condition needs a comparison word: one of %v, is a, is set or is not set", model.Ops)
		return c, false
	}
	c.Type, c.Op, c.Than = model.CondCompare, op, than
	return c, true
}

func opNames() []string {
	out := make([]string, 0, len(model.Ops))
	for _, op := range model.Ops {
		out = append(out, string(op))
	}
	return out
}

func countOps() []string {
	var out []string
	for _, op := range model.Ops {
		if op != model.OpIsOneOf {
			out = append(out, string(op))
		}
	}
	return out
}

func (d decoder) comparison(m source.Map) (model.Op, model.Operand, bool) {
	var found []model.Op
	for _, op := range model.Ops {
		if has(m, string(op)) {
			found = append(found, op)
		}
	}
	if len(found) == 0 {
		return "", model.Operand{}, false
	}
	if len(found) > 1 {
		d.diag(m.Pos, d.rules.Conds, "a condition has one comparison word, not %v", found)
	}
	n, _ := source.Get(m, string(found[0]))
	than, ok := d.operand(n, found[0] == model.OpIsOneOf)
	return found[0], than, ok
}

func (d decoder) operand(n *yaml.Node, list bool) (model.Operand, bool) {
	m, ok := source.AsMap(d.path, n, d.ds)
	if !ok {
		return model.Operand{}, false
	}
	out := model.Operand{Pos: m.Pos}
	if list {
		allow(m, []string{"values"}, d.rules.Conds, d.ds)
		vn, ok := source.Get(m, "values")
		if !ok {
			d.diag(m.Pos, d.rules.Conds, "is one of needs values: with a list")
			return out, false
		}
		items, ok := sequence(d.path, vn, d.ds)
		if !ok {
			return out, false
		}
		for _, item := range items {
			if v, ok := Value(d.path, item, d.ds); ok {
				out.Values = append(out.Values, v)
			}
		}
		return out, true
	}
	allow(m, []string{"value", "same as"}, d.rules.Conds, d.ds)
	if vn, ok := source.Get(m, "value"); ok {
		v, ok := Value(d.path, vn, d.ds)
		out.Value = &v
		return out, ok
	}
	sn, ok := source.Get(m, "same as")
	if !ok {
		d.diag(m.Pos, d.rules.Conds, "a comparison needs value: or same as:")
		return out, false
	}
	if sn.Kind == yaml.ScalarNode {
		if !IsVariable(sn.Value) {
			d.diag(source.PosOf(d.path, sn), d.rules.Vars, "same as: names a variable alone, not %q", sn.Value)
			return out, false
		}
		out.Node = sn.Value
		return out, true
	}
	sm, ok := source.AsMap(d.path, sn, d.ds)
	if !ok {
		return out, false
	}
	allow(sm, []string{"node", "field"}, d.rules.Conds, d.ds)
	node, _ := source.GetScalar(sm, "node", d.ds)
	if !IsVariable(node) {
		d.diag(sm.Pos, d.rules.Vars, "same as: names a variable alone in node:, not %q", node)
		return out, false
	}
	out.Node = node
	out.Field, _ = source.GetScalar(sm, "field", d.ds)
	return out, true
}

func (d decoder) variable(m source.Map, key string) string {
	s, ok := source.GetScalar(m, key, d.ds)
	if !ok {
		d.diag(m.Pos, "DV5", "%s: is missing", key)
		return ""
	}
	if !IsVariable(s) {
		d.diag(m.Pos, "DV5", "%s: names a variable alone, not %q", key, s)
		return ""
	}
	return s
}

func (d decoder) derivations(n *yaml.Node) []model.Derivation {
	items, ok := sequence(d.path, n, d.ds)
	if !ok {
		return nil
	}
	var out []model.Derivation
	for _, item := range items {
		m, ok := source.AsMap(d.path, item, d.ds)
		if !ok {
			continue
		}
		dv := model.Derivation{Pos: m.Pos}
		switch {
		case has(m, "link"):
			allow(m, []string{"link", "from", "to"}, "DV12", d.ds)
			dv.Type = model.DeriveLink
			dv.From, dv.To = d.variable(m, "from"), d.variable(m, "to")
			name, _ := source.GetScalar(m, "link", d.ds)
			ref, ok := resolve(d.scope, name, model.SortLink, m.Pos, "DV13", d.ds)
			if !ok || dv.From == "" || dv.To == "" {
				continue
			}
			dv.Name, dv.Link, dv.Inv = name, ref.Q, ref.Inverted
		case has(m, "becomes"):
			allow(m, []string{"node", "becomes"}, "DV12", d.ds)
			dv.Type = model.DeriveKind
			dv.Node = d.variable(m, "node")
			name, _ := source.GetScalar(m, "becomes", d.ds)
			ref, ok := resolve(d.scope, name, model.SortKind, m.Pos, "DV13", d.ds)
			if !ok || dv.Node == "" {
				continue
			}
			dv.Kind = ref.Q
		case has(m, "set to"):
			allow(m, []string{"node", "field", "set to"}, "DV12", d.ds)
			dv.Type = model.DeriveField
			dv.Node = d.variable(m, "node")
			dv.Field, _ = source.GetScalar(m, "field", d.ds)
			if dv.Field == "" {
				d.diag(m.Pos, "DV13", "set to needs field:")
				continue
			}
			sn, _ := source.Get(m, "set to")
			than, ok := d.operand(sn, false)
			if !ok || dv.Node == "" {
				continue
			}
			dv.SetTo = than
		default:
			d.diag(m.Pos, "DV12", "a derive item adds a link (link:), a kind (becomes:) or a field value (set to:)")
			continue
		}
		out = append(out, dv)
	}
	return out
}

func has(m source.Map, key string) bool {
	_, ok := source.Get(m, key)
	return ok
}
