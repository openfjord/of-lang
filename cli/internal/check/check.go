// SPDX-License-Identifier: Apache-2.0

// Package check finds problems in definitions that loading alone does not:
// names, cycles, links that cannot fit their kinds, variables missing from a
// sentence, rules that can never settle, values that break their own schema,
// and docs that point at nothing.
package check

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/engine"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/value"
	"github.com/openfjord/of-lang/cli/internal/world"
)

var (
	kindName      = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	linkName      = regexp.MustCompile(`^[a-z][a-z0-9]*([A-Z][a-z0-9]*)+$`)
	statementName = regexp.MustCompile(`^[a-z][A-Za-z0-9]*$`)
	docLink       = regexp.MustCompile(`!?\[[^\]]*\]\(([^)\s]+)\)`)
	varInText     = regexp.MustCompile(`\$[a-z][A-Za-z0-9]*`)
)

// World checks every definition of the root package.
func World(w model.World) []diag.Diagnostic {
	var ds []diag.Diagnostic
	root := w.Root
	for _, k := range sortedByPos(w.Kinds, func(k model.Kind) (model.QName, diag.Pos) { return k.Q, k.Pos }) {
		if k.Q.Pkg == root {
			ds = append(ds, checkKind(w, k)...)
		}
	}
	for _, l := range sortedByPos(w.Links, func(l model.Link) (model.QName, diag.Pos) { return l.Q, l.Pos }) {
		if l.Q.Pkg == root {
			ds = append(ds, checkLink(w, l)...)
		}
	}
	for _, vt := range sortedByPos(w.ValueTypes, func(v model.ValueType) (model.QName, diag.Pos) { return v.Q, v.Pos }) {
		if vt.Q.Pkg != root {
			continue
		}
		if !kindName.MatchString(vt.Q.Name) {
			ds = append(ds, diag.New(vt.Pos, "value type %s should be written in UpperCamelCase", vt.Q.Name))
		}
		ds = append(ds, checkSchema(w, vt.Schema, vt.Q.Name)...)
	}
	for _, q := range w.Order {
		if q.Pkg == root {
			ds = append(ds, checkStatement(w.Statements[q])...)
		}
	}
	_, sds := engine.Stratify(w)
	ds = append(ds, sds...)
	pkg := w.Packages[root]
	ds = append(ds, checkExports(pkg)...)
	for _, doc := range pkg.Docs {
		ds = append(ds, checkDoc(doc)...)
	}
	return diag.Sorted(ds)
}

func sortedByPos[V any](m map[model.QName]V, key func(V) (model.QName, diag.Pos)) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	slices.SortFunc(out, func(a, b V) int {
		_, pa := key(a)
		_, pb := key(b)
		if c := strings.Compare(pa.File, pb.File); c != 0 {
			return c
		}
		return pa.Line - pb.Line
	})
	return out
}

func checkKind(w model.World, k model.Kind) []diag.Diagnostic {
	var ds []diag.Diagnostic
	if !kindName.MatchString(k.Q.Name) {
		ds = append(ds, diag.New(k.Pos, "kind %s should be written in UpperCamelCase", k.Q.Name))
	}
	if k.Extends != nil && slices.Contains(world.Ancestors(w, *k.Extends), k.Q) {
		ds = append(ds, diag.New(k.Pos, "%s extends itself through %s", k.Q.Name, k.Extends.Name))
	}
	for _, f := range k.Fields {
		ds = append(ds, checkSchema(w, f.Schema, k.Q.Name+"."+f.Name)...)
		if f.Optional && f.Schema.Default != nil {
			ds = append(ds, diag.New(f.Pos, "%s has a default, so optional: true says nothing", f.Name))
		}
	}
	for _, lc := range k.Links {
		link, ok := w.Links[lc.Ref.Q]
		if !ok {
			continue
		}
		end := link.From
		if lc.Ref.Inverted {
			end = link.To
		}
		if !world.IsA(w, k.Q, end) {
			ds = append(ds, diag.New(lc.Pos, "%s starts at a %s, so a %s cannot have it", lc.Name, end.Name, k.Q.Name))
		}
	}
	return ds
}

func checkLink(w model.World, l model.Link) []diag.Diagnostic {
	var ds []diag.Diagnostic
	if !linkName.MatchString(l.Q.Name) {
		ds = append(ds, diag.New(l.Pos, "link %s needs two or more words in lowerCamelCase, such as ownedByTeam", l.Q.Name))
	}
	if l.Inverse != "" && !linkName.MatchString(l.Inverse) {
		ds = append(ds, diag.New(l.Pos, "inverse %s needs two or more words in lowerCamelCase", l.Inverse))
	}
	if !strings.Contains(l.Reads, "{from}") || !strings.Contains(l.Reads, "{to}") {
		ds = append(ds, diag.New(l.Pos, "reads: of %s must use both {from} and {to}", l.Q.Name))
	}
	if l.Under != nil {
		if slices.Contains(world.SuperLinks(w, *l.Under), l.Q) || *l.Under == l.Q {
			ds = append(ds, diag.New(l.Pos, "%s sits under itself through %s", l.Q.Name, l.Under.Name))
		} else if super, ok := w.Links[*l.Under]; ok {
			if !world.IsA(w, l.From, super.From) || !world.IsA(w, l.To, super.To) {
				ds = append(ds, diag.New(l.Pos, "%s joins %s to %s, which does not fit under %s, which joins %s to %s", l.Q.Name, l.From.Name, l.To.Name, super.Q.Name, super.From.Name, super.To.Name))
			}
		}
	}
	if l.Chains && !world.IsA(w, l.To, l.From) && !world.IsA(w, l.From, l.To) {
		ds = append(ds, diag.New(l.Pos, "%s chains, so it must join a kind to the same kind; it joins %s to %s", l.Q.Name, l.From.Name, l.To.Name))
	}
	return ds
}

func checkSchema(w model.World, s model.Schema, owner string) []diag.Diagnostic {
	var ds []diag.Diagnostic
	base := model.Schema{Type: s.Type, Ref: s.Ref, HasRef: s.HasRef}
	for _, e := range s.Enum {
		for _, msg := range value.Validate(w, base, e) {
			ds = append(ds, diag.New(s.Pos, "%s: enum value %s", owner, msg))
		}
	}
	for name, v := range map[string]*model.Value{"default": s.Default, "const": s.Const} {
		if v == nil {
			continue
		}
		for _, msg := range value.Validate(w, s, *v) {
			ds = append(ds, diag.New(s.Pos, "%s: %s %s", owner, name, msg))
		}
	}
	if s.HasRef {
		ds = append(ds, checkNarrowing(w, s, owner)...)
	}
	if s.Items != nil {
		ds = append(ds, checkSchema(w, *s.Items, owner+" items")...)
	}
	for _, p := range s.Properties {
		ds = append(ds, checkSchema(w, p.Schema, owner+"."+p.Name)...)
	}
	return ds
}

// checkNarrowing reports a field that loosens the value type it names.
func checkNarrowing(w model.World, s model.Schema, owner string) []diag.Diagnostic {
	vt, ok := w.ValueTypes[s.Ref]
	if !ok {
		return nil
	}
	base := vt.Schema
	var ds []diag.Diagnostic
	widen := func(what string) {
		ds = append(ds, diag.New(s.Pos, "%s widens %s's %s; a field may only narrow a value type", owner, vt.Q.Name, what))
	}
	if s.MinLength != nil && base.MinLength != nil && *s.MinLength < *base.MinLength {
		widen("minLength")
	}
	if s.MaxLength != nil && base.MaxLength != nil && *s.MaxLength > *base.MaxLength {
		widen("maxLength")
	}
	if s.Minimum != nil && base.Minimum != nil && s.Minimum.Cmp(base.Minimum) < 0 {
		widen("minimum")
	}
	if s.Maximum != nil && base.Maximum != nil && s.Maximum.Cmp(base.Maximum) > 0 {
		widen("maximum")
	}
	if s.Format != "" && base.Format != "" && s.Format != base.Format {
		widen("format")
	}
	return ds
}

func checkStatement(st model.Statement) []diag.Diagnostic {
	var ds []diag.Diagnostic
	if !statementName.MatchString(st.Q.Name) {
		ds = append(ds, diag.New(st.Pos, "%s should be written in lowerCamelCase", st.Q.Name))
	}
	used := usedVars(st)
	for _, v := range used {
		if !slices.Contains(st.Vars, v) {
			ds = append(ds, diag.New(st.Pos, "%s uses %s, which its sentence does not mention; every variable appears in the sentence", st.Q.Name, v))
		}
	}
	for _, v := range varInText.FindAllString(st.Message, -1) {
		if !slices.Contains(engine.OuterVars(st), v) {
			ds = append(ds, diag.New(st.Pos, "the message of %s uses %s, which no match binds", st.Q.Name, v))
		}
	}
	return ds
}

func usedVars(st model.Statement) []string {
	var out []string
	add := func(v string) {
		if v != "" && !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	addLinks := func(pls []model.PatternLink) {
		for _, pl := range pls {
			add(pl.From)
			add(pl.To)
		}
	}
	var addConds func([]model.Condition)
	addConds = func(cs []model.Condition) {
		for _, c := range cs {
			add(c.Node)
			add(c.Than.Node)
			addLinks(c.Links)
			addConds(c.Where)
			addConds(c.Any)
		}
	}
	addLinks(st.Links)
	addConds(st.Where)
	addConds(st.Require)
	for _, d := range st.Derive {
		add(d.From)
		add(d.To)
		add(d.Node)
		add(d.SetTo.Node)
	}
	return out
}

func checkExports(pkg model.Package) []diag.Diagnostic {
	var ds []diag.Diagnostic
	ex := pkg.Manifest.Exports
	for sort, names := range map[model.Sort][]string{model.SortKind: ex.Kinds, model.SortLink: ex.Links, model.SortValueType: ex.ValueTypes} {
		for _, n := range names {
			d, ok := pkg.Defs[n]
			switch {
			case !ok:
				ds = append(ds, diag.New(pkg.Manifest.Pos, "exports %s, which the package does not define", n))
			case d.InverseOf != "":
				ds = append(ds, diag.New(pkg.Manifest.Pos, "exports %s, an inverse name; export its link %s instead", n, d.InverseOf))
			case d.Sort != sort:
				ds = append(ds, diag.New(pkg.Manifest.Pos, "exports %s as a %s, but it is a %s", n, sort, d.Sort))
			}
		}
	}
	for _, p := range pkg.Manifest.OpenFor {
		name := p.Extend + p.Under
		if _, ok := pkg.Defs[name]; !ok {
			ds = append(ds, diag.New(pkg.Manifest.Pos, "openFor names %s, which the package does not define", name))
		}
	}
	return ds
}

// checkDoc reports links and embeds in a doc that point at a missing file, or
// at a name its YAML file does not define.
func checkDoc(path string) []diag.Diagnostic {
	data, err := os.ReadFile(path)
	if err != nil {
		return []diag.Diagnostic{diag.New(diag.Pos{File: path}, "cannot read: %v", err)}
	}
	var ds []diag.Diagnostic
	for i, line := range strings.Split(string(data), "\n") {
		for _, m := range docLink.FindAllStringSubmatchIndex(line, -1) {
			target := line[m[2]:m[3]]
			if strings.Contains(target, "://") || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") {
				continue
			}
			pos := diag.Pos{File: path, Line: i + 1, Column: m[2] + 1}
			file, name, _ := strings.Cut(target, "#")
			full := filepath.Join(filepath.Dir(path), file)
			if _, err := os.Stat(full); err != nil {
				ds = append(ds, diag.New(pos, "%s does not exist", file))
				continue
			}
			if name == "" || !strings.HasSuffix(file, ".yaml") {
				continue
			}
			if !definesName(full, name) {
				ds = append(ds, diag.New(pos, "%s does not define %s", file, name))
			}
		}
	}
	return ds
}

func definesName(path, name string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil || len(doc.Content) == 0 {
		return false
	}
	root := doc.Content[0]
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == name && name != "import" {
			return true
		}
	}
	return false
}
