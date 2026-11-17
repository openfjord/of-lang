// SPDX-License-Identifier: Apache-2.0

// Package check finds problems in definitions that loading alone does not:
// names, cycles, unions, links that cannot fit their kinds, values that break
// their own schema, variables and comparisons in constraints and derivatives,
// derivatives that can never settle, tags, exports, and docs that point at
// nothing. Each message ends with the rule of the specification it enforces.
package check

import (
	"cmp"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/engine"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
	"github.com/openfjord/of-lang/cli/internal/value"
	"github.com/openfjord/of-lang/cli/internal/world"
)

var (
	upperCamel = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	lowerCamel = regexp.MustCompile(`^[a-z][A-Za-z0-9]*$`)
	linkName   = regexp.MustCompile(`^[a-z][a-z0-9]*([A-Z][a-z0-9]*)+$`)
	docLink    = regexp.MustCompile(`!?\[[^\]]*\]\(([^)\s]+)\)`)
)

// LowerCamel says whether a name is lowerCamelCase, as constraint,
// derivative, test and question names are.
func LowerCamel(s string) bool { return lowerCamel.MatchString(s) }

// World checks every definition of the root package.
func World(w model.World) []diag.Diagnostic {
	var ds []diag.Diagnostic
	root := w.Root
	pkg := w.Packages[root]
	for _, k := range sorted(w.Kinds, func(k model.Kind) diag.Pos { return k.Pos }) {
		if k.Q.Pkg == root {
			ds = append(ds, checkKind(w, k)...)
		}
	}
	for _, l := range sorted(w.Links, func(l model.Link) diag.Pos { return l.Pos }) {
		if l.Q.Pkg == root {
			ds = append(ds, checkLink(w, l)...)
		}
	}
	for _, vt := range sorted(w.ValueTypes, func(v model.ValueType) diag.Pos { return v.Pos }) {
		if vt.Q.Pkg != root {
			continue
		}
		if !upperCamel.MatchString(vt.Q.Name) {
			ds = append(ds, diag.New(vt.Pos, "value type %s must be UpperCamelCase (FL7)", vt.Q.Name))
		}
		ds = append(ds, checkSchema(w, vt.Schema, vt.Q.Name)...)
	}
	for _, q := range w.Order {
		if q.Pkg == root {
			ds = append(ds, Statement(w, w.Statements[q], nil)...)
		}
	}
	ds = append(ds, checkFieldWriters(w)...)
	_, sds := engine.Stratify(w)
	ds = append(ds, sds...)
	ds = append(ds, checkTags(w)...)
	ds = append(ds, checkExports(w, pkg)...)
	for _, doc := range pkg.Docs {
		ds = append(ds, checkDoc(doc)...)
	}
	return diag.Sorted(ds)
}

func sorted[V any](m map[model.QName]V, pos func(V) diag.Pos) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	slices.SortFunc(out, func(a, b V) int {
		pa, pb := pos(a), pos(b)
		return cmp.Or(cmp.Compare(pa.File, pb.File), cmp.Compare(pa.Line, pb.Line))
	})
	return out
}

func checkKind(w model.World, k model.Kind) []diag.Diagnostic {
	var ds []diag.Diagnostic
	add := func(format string, args ...any) { ds = append(ds, diag.New(k.Pos, format, args...)) }
	if !upperCamel.MatchString(k.Q.Name) {
		add("kind %s must be UpperCamelCase (FL7)", k.Q.Name)
	}
	if k.Extends != nil {
		switch {
		case slices.Contains(world.Ancestors(w, *k.Extends), k.Q):
			add("%s extends itself through %s (KD2)", k.Q.Name, k.Extends.Name)
		case world.IsUnion(w, *k.Extends):
			add("%s extends %s, which is a union; no kind extends a union (KD9)", k.Q.Name, k.Extends.Name)
		}
		if k.Extends.Pkg != k.Q.Pkg && !opens(w, k.Extends.Pkg, func(p model.OpenPoint) bool { return p.Extend == k.Extends.Name }) {
			add("%s extends %s, which package %s does not open with openFor: extend (PK12)", k.Q.Name, k.Extends.Name, k.Extends.Pkg)
		}
	}
	if len(k.OneOf) > 0 {
		if len(k.OneOf) < 2 {
			add("oneOf of %s lists one kind; a union lists two or more (KD7)", k.Q.Name)
		}
		for _, a := range k.OneOf {
			for _, b := range k.OneOf {
				if a != b && slices.Contains(world.Ancestors(w, a), b) {
					add("%s lists %s and %s, but %s extends %s, so a thing would be two listed kinds (KD7)", k.Q.Name, a.Name, b.Name, a.Name, b.Name)
				}
			}
		}
		if len(k.Fields) > 0 {
			add("%s is a union, so it has no fields; its listed kinds keep their own (KD8)", k.Q.Name)
		}
		if k.Extends != nil {
			add("%s is a union, so it extends no kind (KD12)", k.Q.Name)
		}
	}
	var inherited []model.Field
	if k.Extends != nil {
		inherited = world.Fields(w, *k.Extends)
	}
	for _, f := range k.Fields {
		if slices.ContainsFunc(inherited, func(i model.Field) bool { return i.Name == f.Name }) {
			ds = append(ds, diag.New(f.Pos, "%s declares %s, which it already has from a kind it extends; put a stricter rule in a constraint (KD13)", k.Q.Name, f.Name))
		}
		ds = append(ds, checkSchema(w, f.Schema, k.Q.Name+"."+f.Name)...)
	}
	return ds
}

// opens says whether a package opens an extension point.
func opens(w model.World, pkg string, match func(model.OpenPoint) bool) bool {
	return slices.ContainsFunc(w.Packages[pkg].Manifest.OpenFor, match)
}

func checkLink(w model.World, l model.Link) []diag.Diagnostic {
	var ds []diag.Diagnostic
	add := func(format string, args ...any) { ds = append(ds, diag.New(l.Pos, format, args...)) }
	if !linkName.MatchString(l.Q.Name) {
		add("link %s needs two or more words in lowerCamelCase, such as usesRecipe (LK1)", l.Q.Name)
	}
	if l.Inverse != "" && !linkName.MatchString(l.Inverse) {
		add("inverse %s needs two or more words in lowerCamelCase (LK1)", l.Inverse)
	}
	if l.Under != nil {
		if *l.Under == l.Q || slices.Contains(world.SuperLinks(w, *l.Under), l.Q) {
			add("%s is under itself through %s (LK4)", l.Q.Name, l.Under.Name)
		} else if super, ok := w.Links[*l.Under]; ok {
			if !world.IsA(w, l.From, super.From) || !world.IsA(w, l.To, super.To) {
				add("%s joins %s to %s, which does not fit under %s, which joins %s to %s (LK4)", l.Q.Name, l.From.Name, l.To.Name, super.Q.Name, super.From.Name, super.To.Name)
			}
		}
		if l.Under.Pkg != l.Q.Pkg && !opens(w, l.Under.Pkg, func(p model.OpenPoint) bool { return p.Under == l.Under.Name }) {
			add("%s is under %s, which package %s does not open with openFor: under (PK12)", l.Q.Name, l.Under.Name, l.Under.Pkg)
		}
	}
	if l.Chains && !world.IsA(w, l.To, l.From) {
		add("%s chains, so its to must be its from or a kind that extends it; it joins %s to %s (LK5)", l.Q.Name, l.From.Name, l.To.Name)
	}
	return ds
}

func checkSchema(w model.World, s model.Schema, owner string) []diag.Diagnostic {
	var ds []diag.Diagnostic
	full := world.Resolved(w, s)
	base := model.Schema{Type: full.Type, Ref: s.Ref, HasRef: s.HasRef}
	for _, e := range s.Enum {
		for _, msg := range value.Validate(w, base, e) {
			ds = append(ds, diag.New(s.Pos, "%s: enum value %s (FD5)", owner, msg))
		}
	}
	for _, kv := range []struct {
		name string
		v    *model.Value
	}{{"default", s.Default}, {"const", s.Const}} {
		if kv.v == nil {
			continue
		}
		for _, msg := range value.Validate(w, s, *kv.v) {
			ds = append(ds, diag.New(s.Pos, "%s: %s %s (FD5)", owner, kv.name, msg))
		}
	}
	if s.HasRef {
		ds = append(ds, checkNarrowing(w, s, owner)...)
		if vt, ok := w.ValueTypes[s.Ref]; ok && len(vt.Schema.OneOf) > 0 {
			for _, k := range s.Keywords {
				if !slices.Contains([]string{"type", "optional", "default"}, k) {
					ds = append(ds, diag.New(s.Pos, "%s names %s, which has oneOf, so it adds only optional or default, not %s (VT6)", owner, vt.Q.Name, k))
				}
			}
		}
	}
	if s.Items != nil {
		ds = append(ds, checkSchema(w, *s.Items, owner+" items")...)
	}
	for _, p := range s.Properties {
		ds = append(ds, checkSchema(w, p.Schema, owner+"."+p.Name)...)
	}
	for _, c := range s.OneOf {
		ds = append(ds, checkSchema(w, c, owner+" choice")...)
	}
	return ds
}

// checkNarrowing reports a field that widens the value type it names.
func checkNarrowing(w model.World, s model.Schema, owner string) []diag.Diagnostic {
	vt, ok := w.ValueTypes[s.Ref]
	if !ok {
		return nil
	}
	base := world.Resolved(w, vt.Schema)
	var ds []diag.Diagnostic
	widen := func(what string) {
		ds = append(ds, diag.New(s.Pos, "%s widens the %s of %s; a field may only narrow a value type (VT3)", owner, what, vt.Q.Name))
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

func checkExports(w model.World, pkg model.Package) []diag.Diagnostic {
	var ds []diag.Diagnostic
	ex := pkg.Manifest.Exports
	for _, list := range []struct {
		sort  model.Sort
		names []string
	}{{model.SortKind, ex.Kinds}, {model.SortLink, ex.Links}, {model.SortValueType, ex.ValueTypes}} {
		for _, n := range list.names {
			d, ok := pkg.Defs[n]
			switch {
			case !ok:
				ds = append(ds, diag.New(pkg.Manifest.Pos, "exports %s, which the package does not define (PK11)", n))
			case d.InverseOf != "":
				ds = append(ds, diag.New(pkg.Manifest.Pos, "exports %s, an inverse name; export its link %s, which exports both (PK11)", n, d.InverseOf))
			case d.Sort != list.sort:
				ds = append(ds, diag.New(pkg.Manifest.Pos, "exports %s as a %s, but it is a %s (PK11)", n, list.sort, d.Sort))
			}
		}
	}
	for _, p := range pkg.Manifest.OpenFor {
		switch {
		case p.Extend != "":
			d, ok := pkg.Defs[p.Extend]
			switch {
			case !ok || d.Sort != model.SortKind:
				ds = append(ds, diag.New(p.Pos, "openFor: extend names %s, which is not a kind of the package (PK13)", p.Extend))
			case world.IsUnion(w, model.QName{Pkg: pkg.Name, Name: p.Extend}):
				ds = append(ds, diag.New(p.Pos, "openFor: extend names %s, which is a union; no kind extends a union (PK13)", p.Extend))
			}
		default:
			d, ok := pkg.Defs[p.Under]
			if !ok || d.Sort != model.SortLink || d.InverseOf != "" {
				ds = append(ds, diag.New(p.Pos, "openFor: under names %s, which is not a link of the package (PK13)", p.Under))
			}
		}
	}
	return ds
}

// checkDoc reports links and embeds in a doc that point at a missing file,
// or at a name its file does not define.
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
				ds = append(ds, diag.New(pos, "%s does not exist (DC2)", file))
				continue
			}
			if _, ok := source.FormOf(file); !ok || name == "" {
				continue
			}
			if !definesName(full, name) {
				ds = append(ds, diag.New(pos, "%s has no statement %s (DC2)", file, name))
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
