// SPDX-License-Identifier: Apache-2.0

package check

import (
	"maps"
	"slices"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/engine"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/value"
	"github.com/openfjord/of-lang/cli/internal/world"
)

// use counts where a variable occurs outside its kind item.
type use struct {
	pos   diag.Pos
	ends  int
	other bool
}

// stmt holds what checking one statement has found so far.
type stmt struct {
	w     model.World
	st    model.Statement
	rules model.Rules
	// declared holds every variable of the statement, in any block.
	declared map[string]*use
	ds       []diag.Diagnostic
}

func (c *stmt) add(pos diag.Pos, rule, format string, args ...any) {
	c.ds = append(c.ds, diag.New(pos, format+" ("+rule+")", args...))
}

// Statement checks a constraint, a derivative or a question. expect lists the
// variables a question's expect names, which count as uses.
func Statement(w model.World, st model.Statement, expect []string) []diag.Diagnostic {
	c := &stmt{w: w, st: st, rules: model.RulesFor(st.Sort), declared: map[string]*use{}}
	if !lowerCamel.MatchString(st.Q.Name) {
		c.add(st.Pos, "FL7", "%s must be lowerCamelCase", st.Q.Name)
	}
	env := c.block(st.Match, nil)
	c.conditions(st.Unless, env)
	for _, m := range engine.MessageVar.FindAllStringSubmatch(st.Message, -1) {
		kind, ok := c.lookup(st.Pos, m[1], env, c.rules.Vars)
		if !ok {
			continue
		}
		if m[2] != "" {
			if _, _, ok := world.Field(w, kind, m[2]); !ok {
				c.add(st.Pos, c.rules.Fields, "the message names %s of %s, which is not a field of %s", m[2], m[1], kind.Name)
			}
		}
	}
	for _, d := range st.Derive {
		c.derivation(d, env)
	}
	for _, v := range expect {
		if u, ok := c.declared[v]; ok {
			u.other = true
		}
	}
	for _, v := range slices.Sorted(maps.Keys(c.declared)) {
		u := c.declared[v]
		if u.ends < 2 && !u.other {
			c.add(u.pos, c.rules.Unused, "%s is used at one place only; give that end its kind alone", v)
		}
	}
	return c.ds
}

// lookup finds the kind of a variable in scope, and counts the use.
func (c *stmt) lookup(pos diag.Pos, v string, env map[string]model.QName, rule string) (model.QName, bool) {
	kind, ok := env[v]
	switch {
	case ok:
		c.declared[v].other = true
		return kind, true
	case c.declared[v] != nil:
		c.add(pos, c.rules.Scope, "%s belongs to a nested block and is not used outside it", v)
	default:
		c.add(pos, rule, "%s has no kind item", v)
	}
	return model.QName{}, false
}

// block checks a match and its filter, and returns the variables in scope
// inside it: the enclosing ones and its own.
func (c *stmt) block(b model.Block, outer map[string]model.QName) map[string]model.QName {
	env := maps.Clone(outer)
	if env == nil {
		env = map[string]model.QName{}
	}
	for _, ki := range b.Kinds {
		if c.declared[ki.Var] != nil {
			c.add(ki.Pos, c.rules.OneKindItem, "%s already has a kind item", ki.Var)
			continue
		}
		c.declared[ki.Var] = &use{pos: ki.Pos}
		env[ki.Var] = ki.Kind
	}
	for _, li := range b.Links {
		var want [2]model.QName
		want[0], want[1] = world.LinkEnds(c.w, li.Link, li.Inverted)
		for i, e := range []model.End{li.From, li.To} {
			kind := e.Kind
			if e.Var != "" {
				k, ok := env[e.Var]
				switch {
				case !ok && c.declared[e.Var] != nil:
					c.add(e.Pos, c.rules.Scope, "%s belongs to a nested block and is not used outside it", e.Var)
					continue
				case !ok:
					c.add(e.Pos, c.rules.Ends, "%s has no kind item in this match or one that encloses it", e.Var)
					continue
				}
				c.declared[e.Var].ends++
				kind = k
			}
			if !world.IsA(c.w, kind, want[i]) {
				c.add(e.Pos, c.rules.EndKinds, "%s joins a %s here, and %s is not a %s", li.Name, want[i].Name, kind.Name, want[i].Name)
			}
		}
	}
	c.conditions(b.Filter, env)
	return env
}

func (c *stmt) conditions(conds []model.Condition, env map[string]model.QName) {
	for _, cond := range conds {
		c.condition(cond, env)
	}
}

func (c *stmt) condition(cond model.Condition, env map[string]model.QName) {
	switch cond.Type {
	case model.CondAnyOf:
		c.conditions(cond.Any, env)
		return
	case model.CondAtLeastOne, model.CondNone, model.CondCount:
		c.block(*cond.Block, env)
		return
	default:
		// A condition on a node, checked below.
	}
	kind, ok := c.lookup(cond.Pos, cond.Node, env, c.rules.Vars)
	if !ok {
		return
	}
	var schema *model.Schema
	if cond.Field != "" {
		f, _, ok := world.Field(c.w, kind, cond.Field)
		if !ok {
			c.add(cond.Pos, c.rules.Fields, "%s is not a field of %s", cond.Field, kind.Name)
			return
		}
		s := f.Schema
		if len(cond.Part) > 0 {
			ps, ok := world.Part(c.w, s, cond.Part)
			if !ok {
				c.add(cond.Pos, c.rules.Parts, "the schema of %s has no part %v", cond.Field, cond.Part)
				return
			}
			s = ps
		}
		schema = &s
	}
	if cond.Type != model.CondCompare {
		return
	}
	c.compared(cond, kind, schema, env)
}

// compared checks what a comparison compares with.
func (c *stmt) compared(cond model.Condition, kind model.QName, schema *model.Schema, env map[string]model.QName) {
	than := cond.Than
	literals := than.Values
	if than.Value != nil {
		literals = append(literals, *than.Value)
	}
	for _, v := range literals {
		switch {
		case schema != nil:
			for _, msg := range value.Validate(c.w, *schema, v) {
				c.add(than.Pos, c.rules.Types, "%s does not allow %s: %s", cond.Field, value.Format(v), msg)
			}
		case v.Kind != model.ValueString:
			c.add(than.Pos, c.rules.Types, "%s without a field is compared by its id, so the value must be the id of a fact, not %s", cond.Node, value.Format(v))
		}
	}
	if than.Node == "" {
		return
	}
	otherKind, ok := c.lookup(than.Pos, than.Node, env, c.rules.Vars)
	if !ok {
		return
	}
	var other *model.Schema
	if than.Field != "" {
		f, _, ok := world.Field(c.w, otherKind, than.Field)
		if !ok {
			c.add(than.Pos, c.rules.Fields, "%s is not a field of %s", than.Field, otherKind.Name)
			return
		}
		other = &f.Schema
	}
	switch {
	case schema == nil && other == nil:
		if !world.Compatible(c.w, kind, otherKind) {
			c.add(cond.Pos, c.rules.NodeKinds, "%s is a %s and %s is a %s, and neither kind is the other", cond.Node, kind.Name, than.Node, otherKind.Name)
		}
	case schema == nil || other == nil:
		c.add(cond.Pos, c.rules.Types, "a comparison compares a thing with a thing, or a field with a field")
	default:
		a, okA := world.TypeOf(*schema)
		b, okB := world.TypeOf(*other)
		if okA && okB && a != b {
			c.add(cond.Pos, c.rules.Types, "%s is a %s and %s is a %s; compare values of one type", cond.Field, world.TypeText(a), than.Field, world.TypeText(b))
		}
	}
}

func (c *stmt) derivation(d model.Derivation, env map[string]model.QName) {
	kindOf := func(v string) (model.QName, bool) {
		k, ok := env[v]
		switch {
		case ok:
			c.declared[v].other = true
			return k, true
		case c.declared[v] != nil:
			c.add(d.Pos, "DV14", "%s has its kind item in a nested block; derive uses only variables of the derivative's own match", v)
		default:
			c.add(d.Pos, "DV5", "%s has no kind item", v)
		}
		return model.QName{}, false
	}
	switch d.Type {
	case model.DeriveLink:
		from, okF := kindOf(d.From)
		to, okT := kindOf(d.To)
		if d.Link.Pkg != c.st.Q.Pkg {
			c.add(d.Pos, "PK12", "%s is a link of package %s; a derivative adds only links of its own package", d.Name, d.Link.Pkg)
		}
		if !okF || !okT {
			return
		}
		wantFrom, wantTo := world.LinkEnds(c.w, d.Link, d.Inv)
		if !world.IsA(c.w, from, wantFrom) || !world.IsA(c.w, to, wantTo) {
			c.add(d.Pos, "DV13", "%s joins a %s to a %s, not a %s to a %s", d.Name, wantFrom.Name, wantTo.Name, from.Name, to.Name)
		}
	case model.DeriveKind:
		kind, ok := kindOf(d.Node)
		if d.Kind.Pkg != c.st.Q.Pkg {
			c.add(d.Pos, "PK12", "%s is a kind of package %s; a derivative adds only kinds of its own package", d.Kind.Name, d.Kind.Pkg)
		}
		if world.IsUnion(c.w, d.Kind) {
			c.add(d.Pos, "KD10", "%s is a union; becomes names one of its listed kinds", d.Kind.Name)
			return
		}
		if ok && (d.Kind == kind || !slices.Contains(world.Ancestors(c.w, d.Kind), kind)) {
			c.add(d.Pos, "DV13", "%s does not extend %s, the kind of %s", d.Kind.Name, kind.Name, d.Node)
		}
	case model.DeriveField:
		kind, ok := kindOf(d.Node)
		if !ok {
			return
		}
		f, owner, ok := world.Field(c.w, kind, d.Field)
		if !ok {
			c.add(d.Pos, "DV13", "%s is not a field of %s", d.Field, kind.Name)
			return
		}
		if owner.Pkg != c.st.Q.Pkg {
			c.add(d.Pos, "PK12", "%s is a field of %s, a kind of package %s; a derivative sets only fields of its own package", d.Field, owner.Name, owner.Pkg)
		}
		if d.SetTo.Value != nil {
			for _, msg := range value.Validate(c.w, f.Schema, *d.SetTo.Value) {
				c.add(d.Pos, "DV13", "%s does not allow %s: %s", d.Field, value.Format(*d.SetTo.Value), msg)
			}
			return
		}
		otherKind, ok := c.lookup(d.Pos, d.SetTo.Node, env, "DV5")
		if !ok {
			return
		}
		if d.SetTo.Field == "" {
			c.add(d.Pos, "DV13", "set to with same as names a field of %s", d.SetTo.Node)
			return
		}
		of, _, ok := world.Field(c.w, otherKind, d.SetTo.Field)
		if !ok {
			c.add(d.Pos, "DV9", "%s is not a field of %s", d.SetTo.Field, otherKind.Name)
			return
		}
		a, okA := world.TypeOf(f.Schema)
		b, okB := world.TypeOf(of.Schema)
		if okA && okB && a != b {
			c.add(d.Pos, "DV13", "%s is a %s and %s is a %s; a field takes a value of its own type", d.Field, world.TypeText(a), d.SetTo.Field, world.TypeText(b))
		}
	}
}

// checkFieldWriters reports two derivatives that can give different values
// to the same field of the same kind.
func checkFieldWriters(w model.World) []diag.Diagnostic {
	type writer struct {
		deriv model.QName
		value *model.Value
	}
	writers := map[string][]writer{}
	var ds []diag.Diagnostic
	for _, q := range w.Order {
		st := w.Statements[q]
		if st.Sort != model.SortDerivative || q.Pkg != w.Root {
			continue
		}
		env := map[string]model.QName{}
		for _, ki := range st.Match.Kinds {
			env[ki.Var] = ki.Kind
		}
		for _, d := range st.Derive {
			if d.Type != model.DeriveField {
				continue
			}
			_, owner, ok := world.Field(w, env[d.Node], d.Field)
			if !ok {
				continue
			}
			key := owner.Pkg + "." + owner.Name + "." + d.Field
			for _, prev := range writers[key] {
				same := prev.value != nil && d.SetTo.Value != nil && value.Equal(*prev.value, *d.SetTo.Value)
				if prev.deriv != q && !same {
					ds = append(ds, diag.New(d.Pos, "%s and %s can both set %s of a %s, to different values; a field has one value (DV15)", prev.deriv.Name, q.Name, d.Field, owner.Name))
				}
			}
			writers[key] = append(writers[key], writer{deriv: q, value: d.SetTo.Value})
		}
	}
	return ds
}
