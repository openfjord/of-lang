// SPDX-License-Identifier: Apache-2.0

package load

import (
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/internal/diag"
	"github.com/openfjord/of-lang/internal/model"
	"github.com/openfjord/of-lang/internal/source"
)

// Statement decodes a constraint, a rule or a question. sentenceKey is the
// key whose sentence declares the variables: says, or asks for a question.
// extra lists further keys the caller reads itself.
func Statement(q model.QName, sort model.Sort, path string, m source.Map, scope map[string]model.Ref, sentenceKey string, extra []string, ds *[]diag.Diagnostic) model.Statement {
	allowed := []string{"is", sentenceKey, "links", "where"}
	switch sort {
	case model.SortConstraint:
		allowed = append(allowed, "require", "message")
	case model.SortRule:
		allowed = append(allowed, "derive")
	default:
		// A question has only a pattern, and the caller's extra keys.
	}
	source.Allow(m, append(allowed, extra...), ds)
	st := model.Statement{Q: q, Sort: sort, Pos: m.Pos}
	st.Says = strings.TrimSpace(requiredText(path, m, sentenceKey, q.Name, ds))
	for _, loc := range varPattern.FindAllStringIndex(st.Says, -1) {
		v := st.Says[loc[0]:loc[1]]
		if loc[0] == 0 || st.Says[loc[0]-1] != '(' || loc[1] == len(st.Says) || st.Says[loc[1]] != ')' {
			*ds = append(*ds, diag.New(m.Pos, "%s: write %s as (%s) in the sentence", sentenceKey, v, v))
			continue
		}
		if !slices.Contains(st.Vars, v) {
			st.Vars = append(st.Vars, v)
		}
	}
	if n, ok := source.Get(m, "links"); ok {
		st.Links = patternLinks(path, n, scope, ds)
	}
	if n, ok := source.Get(m, "where"); ok {
		st.Where = conditions(path, n, scope, ds)
	}
	if n, ok := source.Get(m, "require"); ok {
		st.Require = conditions(path, n, scope, ds)
	} else if sort == model.SortConstraint {
		*ds = append(*ds, diag.New(m.Pos, "constraint %s needs require:", q.Name))
	}
	if sort == model.SortConstraint {
		st.Message = requiredText(path, m, "message", q.Name, ds)
	}
	if n, ok := source.Get(m, "derive"); ok {
		st.Derive = derivations(path, n, scope, ds)
	} else if sort == model.SortRule {
		*ds = append(*ds, diag.New(m.Pos, "rule %s needs derive:", q.Name))
	}
	return st
}

func patternLinks(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) []model.PatternLink {
	items, ok := sequence(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.PatternLink
	for _, item := range items {
		m, ok := source.AsMap(path, item, ds)
		if !ok {
			continue
		}
		source.Allow(m, []string{"from", "link", "to", "repeated"}, ds)
		pl := model.PatternLink{Pos: m.Pos}
		pl.From = variable(path, m, "from", ds)
		pl.To = variable(path, m, "to", ds)
		name, ok := source.GetScalar(m, "link", ds)
		if !ok {
			*ds = append(*ds, diag.New(m.Pos, "a link in a pattern needs link:"))
			continue
		}
		ref, ok := resolve(scope, name, model.SortLink, m.Pos, ds)
		if !ok {
			continue
		}
		pl.Name, pl.Link, pl.Inverted = name, ref.Q, ref.Inverted
		if rn, ok := source.Get(m, "repeated"); ok {
			pl.Repeated, _ = source.Bool(path, rn, ds)
		}
		out = append(out, pl)
	}
	return out
}

func variable(path string, m source.Map, key string, ds *[]diag.Diagnostic) string {
	n, ok := source.Get(m, key)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "%s: is missing", key))
		return ""
	}
	s, ok := source.Scalar(path, n, ds)
	if !ok {
		return ""
	}
	if !isVariable(s) {
		*ds = append(*ds, diag.New(source.PosOf(path, n), "%s must be a variable starting with $, not %q", key, s))
		return ""
	}
	return s
}

func isVariable(s string) bool {
	return varPattern.FindString(s) == s && s != ""
}

func conditions(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) []model.Condition {
	items, ok := sequence(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.Condition
	for _, item := range items {
		if c, ok := condition(path, item, scope, ds); ok {
			out = append(out, c)
		}
	}
	return out
}

func condition(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) (model.Condition, bool) {
	m, ok := source.AsMap(path, n, ds)
	if !ok {
		return model.Condition{}, false
	}
	c := model.Condition{Pos: m.Pos}
	if an, ok := source.Get(m, "any of"); ok {
		source.Allow(m, []string{"any of"}, ds)
		c.Type = model.CondAnyOf
		c.Any = conditions(path, an, scope, ds)
		return c, true
	}
	for _, nested := range []model.CondType{model.CondAtLeastOne, model.CondNone, model.CondCount} {
		nn, ok := source.Get(m, string(nested))
		if !ok {
			continue
		}
		c.Type = nested
		nm, ok := source.AsMap(path, nn, ds)
		if !ok {
			return c, false
		}
		source.Allow(nm, []string{"links", "where"}, ds)
		if ln, ok := source.Get(nm, "links"); ok {
			c.Links = patternLinks(path, ln, scope, ds)
		} else {
			*ds = append(*ds, diag.New(nm.Pos, "%s needs links:", nested))
		}
		if wn, ok := source.Get(nm, "where"); ok {
			c.Where = conditions(path, wn, scope, ds)
		}
		if nested != model.CondCount {
			source.Allow(m, []string{string(nested)}, ds)
			return c, true
		}
		source.Allow(m, append([]string{"count"}, opNames()...), ds)
		op, than, ok := comparison(path, m, ds)
		if !ok {
			*ds = append(*ds, diag.New(m.Pos, "count needs a comparison word such as is above:"))
			return c, false
		}
		c.Op, c.Than = op, than
		return c, true
	}
	node, ok := source.GetScalar(m, "node", ds)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "a condition starts with node:, any of:, at least one:, none: or count:"))
		return c, false
	}
	if !isVariable(node) {
		*ds = append(*ds, diag.New(m.Pos, "node: must be a variable starting with $, not %q", node))
		return c, false
	}
	c.Node = node
	c.Field, _ = source.GetScalar(m, "field", ds)
	source.Allow(m, append([]string{"node", "field", "is a", "is set", "is not set"}, opNames()...), ds)
	if kn, ok := source.Get(m, "is a"); ok {
		c.Type = model.CondIsA
		if name, ok := source.Scalar(path, kn, ds); ok {
			if ref, ok := resolve(scope, name, model.SortKind, source.PosOf(path, kn), ds); ok {
				c.Kind = ref.Q
			}
		}
		return c, true
	}
	for _, t := range []model.CondType{model.CondIsSet, model.CondIsNotSet} {
		if _, ok := source.Get(m, string(t)); ok {
			if c.Field == "" {
				*ds = append(*ds, diag.New(m.Pos, "%s needs field:", t))
			}
			c.Type = t
			return c, true
		}
	}
	op, than, ok := comparison(path, m, ds)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "a condition needs a comparison word: one of %v, is a, is set or is not set", model.Ops))
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

func comparison(path string, m source.Map, ds *[]diag.Diagnostic) (model.Op, model.Operand, bool) {
	for _, op := range model.Ops {
		n, ok := source.Get(m, string(op))
		if !ok {
			continue
		}
		than, ok := operand(path, n, op == model.OpIsOneOf, ds)
		return op, than, ok
	}
	return "", model.Operand{}, false
}

func operand(path string, n *yaml.Node, list bool, ds *[]diag.Diagnostic) (model.Operand, bool) {
	m, ok := source.AsMap(path, n, ds)
	if !ok {
		return model.Operand{}, false
	}
	source.Allow(m, []string{"value", "values", "same as"}, ds)
	if list {
		vn, ok := source.Get(m, "values")
		if !ok {
			*ds = append(*ds, diag.New(m.Pos, "is one of needs values: with a list"))
			return model.Operand{}, false
		}
		items, ok := sequence(path, vn, ds)
		if !ok {
			return model.Operand{}, false
		}
		var out model.Operand
		for _, item := range items {
			if v, ok := Value(path, item, ds); ok {
				out.Values = append(out.Values, v)
			}
		}
		return out, true
	}
	if vn, ok := source.Get(m, "value"); ok {
		v, ok := Value(path, vn, ds)
		return model.Operand{Value: &v}, ok
	}
	sn, ok := source.Get(m, "same as")
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "a comparison needs value: or same as:"))
		return model.Operand{}, false
	}
	if sn.Kind == yaml.ScalarNode {
		if !isVariable(sn.Value) {
			*ds = append(*ds, diag.New(source.PosOf(path, sn), "same as: names a variable starting with $, not %q", sn.Value))
			return model.Operand{}, false
		}
		return model.Operand{Node: sn.Value}, true
	}
	sm, ok := source.AsMap(path, sn, ds)
	if !ok {
		return model.Operand{}, false
	}
	source.Allow(sm, []string{"node", "field"}, ds)
	node := variable(path, sm, "node", ds)
	field, _ := source.GetScalar(sm, "field", ds)
	return model.Operand{Node: node, Field: field}, node != ""
}

func derivations(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) []model.Derivation {
	items, ok := sequence(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.Derivation
	for _, item := range items {
		m, ok := source.AsMap(path, item, ds)
		if !ok {
			continue
		}
		d := model.Derivation{Pos: m.Pos}
		switch {
		case has(m, "link"):
			source.Allow(m, []string{"from", "link", "to"}, ds)
			d.Type = model.DeriveLink
			d.From = variable(path, m, "from", ds)
			d.To = variable(path, m, "to", ds)
			name, _ := source.GetScalar(m, "link", ds)
			ref, ok := resolve(scope, name, model.SortLink, m.Pos, ds)
			if !ok {
				continue
			}
			d.Name, d.Link, d.Inv = name, ref.Q, ref.Inverted
		case has(m, "becomes"):
			source.Allow(m, []string{"node", "becomes"}, ds)
			d.Type = model.DeriveKind
			d.Node = variable(path, m, "node", ds)
			name, _ := source.GetScalar(m, "becomes", ds)
			ref, ok := resolve(scope, name, model.SortKind, m.Pos, ds)
			if !ok {
				continue
			}
			d.Kind = ref.Q
		case has(m, "set to"):
			source.Allow(m, []string{"node", "field", "set to"}, ds)
			d.Type = model.DeriveField
			d.Node = variable(path, m, "node", ds)
			d.Field, _ = source.GetScalar(m, "field", ds)
			sn, _ := source.Get(m, "set to")
			than, ok := operand(path, sn, false, ds)
			if !ok {
				continue
			}
			d.SetTo = than
		default:
			*ds = append(*ds, diag.New(m.Pos, "a derivation adds a link (link:), a kind (becomes:) or a field value (set to:)"))
			continue
		}
		out = append(out, d)
	}
	return out
}

func has(m source.Map, key string) bool {
	_, ok := source.Get(m, key)
	return ok
}
