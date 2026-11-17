// SPDX-License-Identifier: Apache-2.0

package load

import (
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
)

var varPattern = regexp.MustCompile(`\$[a-z][A-Za-z0-9]*`)

// BuiltinTypes are the JSON Schema types a field may name without importing.
var BuiltinTypes = []string{"string", "integer", "number", "boolean", "array", "object"}

func decodeDefinition(w *model.World, pkg model.Package, path string, e source.Entry, scope map[string]model.Ref) []diag.Diagnostic {
	var ds []diag.Diagnostic
	m, ok := source.AsMap(path, e.Node, &ds)
	if !ok {
		return ds
	}
	q := model.QName{Pkg: pkg.Name, Name: e.Name}
	switch model.Sort(e.Is) {
	case model.SortKind:
		w.Kinds[q] = decodeKind(q, path, m, scope, &ds)
	case model.SortLink:
		w.Links[q] = decodeLink(q, path, m, scope, &ds)
	case model.SortValueType:
		w.ValueTypes[q] = decodeValueType(q, path, m, scope, &ds)
	case model.SortDoc:
		w.Docs[q] = decodeDoc(q, pkg, path, m, &ds)
	case model.SortConstraint, model.SortDerivative:
		w.Statements[q] = Statement(q, model.Sort(e.Is), path, m, scope, &ds)
		w.Order = append(w.Order, q)
	default:
		// Reported when the package's files were collected.
	}
	return ds
}

// decodeDoc reads a Doc: the path of a Markdown file in the docs/ folder of
// its package, relative to the file of the Doc.
func decodeDoc(q model.QName, pkg model.Package, path string, m source.Map, ds *[]diag.Diagnostic) model.Doc {
	allow(m, []string{"is", "path"}, "FL2", ds)
	d := model.Doc{Q: q, Pos: m.Pos}
	p, ok := pairOf(m, "path")
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "doc %s needs path: (DC3)", q.Name))
		return d
	}
	text, ok := source.Scalar(path, p.Value, ds)
	if !ok {
		return d
	}
	rel, found := strings.CutPrefix(text, "file:")
	if !found {
		*ds = append(*ds, diag.New(p.Pos, "path of %s starts with file: (DC3)", q.Name))
		return d
	}
	full := filepath.Join(filepath.Dir(path), rel)
	inDocs, err := filepath.Rel(filepath.Join(pkg.Dir, "docs"), full)
	switch {
	case err != nil || inDocs == ".." || strings.HasPrefix(inDocs, ".."+string(filepath.Separator)):
		*ds = append(*ds, diag.New(p.Pos, "%s is not in the docs/ folder of the package (DC3)", rel))
	case filepath.Ext(full) != ".md":
		*ds = append(*ds, diag.New(p.Pos, "%s is not a Markdown file (.md) (DC3)", rel))
	default:
		if _, err := os.Stat(full); err != nil {
			*ds = append(*ds, diag.New(p.Pos, "%s does not exist (DC3)", rel))
		} else {
			d.Path = full
		}
	}
	return d
}

// allow reports every key that is not in the allowed list, naming the rule.
func allow(m source.Map, allowed []string, rule string, ds *[]diag.Diagnostic) {
	for _, p := range m.Pairs {
		if !slices.Contains(allowed, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not allowed here; allowed: %v (%s)", p.Key, allowed, rule))
		}
	}
}

func decodeKind(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.Kind {
	allow(m, []string{"is", "definition", "extends", "oneOf", "examples", "notExamples", "fields", "tags", "docs"}, "FL2", ds)
	k := model.Kind{Q: q, Pos: m.Pos}
	k.Definition = requiredText(path, m, "definition", q.Name, "KD1", ds)
	if n, ok := source.Get(m, "extends"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			if ref, ok := resolve(scope, name, model.SortKind, source.PosOf(path, n), "KD2", ds); ok {
				k.Extends = &ref.Q
			}
		}
	}
	if n, ok := source.Get(m, "oneOf"); ok {
		names, _ := source.Strings(path, n, ds)
		for _, name := range names {
			if ref, ok := resolve(scope, name, model.SortKind, source.PosOf(path, n), "KD7", ds); ok {
				k.OneOf = append(k.OneOf, ref.Q)
			}
		}
	}
	k.Examples = optionalStrings(path, m, "examples", ds)
	k.NotExamples = optionalStrings(path, m, "notExamples", ds)
	if n, ok := source.Get(m, "fields"); ok {
		k.Fields = decodeFields(path, n, scope, ds)
	}
	k.Tags = decodeTags(path, m, ds)
	if n, ok := source.Get(m, "docs"); ok {
		names, _ := source.Strings(path, n, ds)
		for _, name := range names {
			ref, ok := resolve(scope, name, model.SortDoc, source.PosOf(path, n), "KD14", ds)
			switch {
			case !ok:
			case ref.Q.Pkg != q.Pkg:
				*ds = append(*ds, diag.New(source.PosOf(path, n), "%s is a doc of package %s; a kind lists the docs of its own package (KD14)", name, ref.Q.Pkg))
			default:
				k.Docs = append(k.Docs, ref.Q)
			}
		}
	}
	return k
}

func decodeLink(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.Link {
	allow(m, []string{"is", "from", "to", "reads", "definition", "under", "inverse", "chains", "tags"}, "FL2", ds)
	l := model.Link{Q: q, Pos: m.Pos}
	l.From = kindRef(path, m, "from", q.Name, scope, ds)
	l.To = kindRef(path, m, "to", q.Name, scope, ds)
	l.Reads = requiredText(path, m, "reads", q.Name, "LK2", ds)
	if p, ok := pairOf(m, "reads"); ok && l.Reads != "" && (!strings.Contains(l.Reads, "{from}") || !strings.Contains(l.Reads, "{to}")) {
		*ds = append(*ds, diag.New(p.Pos, "the reads of %s must use both {from} and {to} (LK2)", q.Name))
	}
	l.Definition = requiredText(path, m, "definition", q.Name, "LK2", ds)
	if n, ok := source.Get(m, "under"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			if ref, ok := resolve(scope, name, model.SortLink, source.PosOf(path, n), "LK4", ds); ok {
				if ref.Inverted {
					*ds = append(*ds, diag.New(source.PosOf(path, n), "%s is an inverse name; under: names a declared link (LK4)", name))
				}
				l.Under = &ref.Q
			}
		}
	}
	l.Inverse, _ = source.GetScalar(m, "inverse", ds)
	if n, ok := source.Get(m, "chains"); ok {
		l.Chains, _ = source.Bool(path, n, ds)
	}
	l.Tags = decodeTags(path, m, ds)
	return l
}

func kindRef(path string, m source.Map, key, owner string, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.QName {
	n, ok := source.Get(m, key)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "%s needs %s: (LK3)", owner, key))
		return model.QName{}
	}
	name, ok := source.Scalar(path, n, ds)
	if !ok {
		return model.QName{}
	}
	ref, ok := resolve(scope, name, model.SortKind, source.PosOf(path, n), "LK3", ds)
	if !ok {
		return model.QName{}
	}
	return ref.Q
}

func decodeValueType(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.ValueType {
	if slices.Contains(BuiltinTypes, q.Name) {
		*ds = append(*ds, diag.New(m.Pos, "%s is a built-in type and cannot name a value type (VT2)", q.Name))
	}
	vt := model.ValueType{Q: q, Pos: m.Pos}
	vt.Definition = requiredText(path, m, "definition", q.Name, "VT1", ds)
	vt.Schema = decodeSchema(path, m, scope, []string{"is", "definition", "tags"}, ds)
	vt.Tags = decodeTags(path, m, ds)
	return vt
}

// decodeTags reads a definition's tags. Whether the manifest declares them is
// a check of its own.
func decodeTags(path string, m source.Map, ds *[]diag.Diagnostic) []model.Tag {
	n, ok := source.Get(m, "tags")
	if !ok {
		return nil
	}
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(source.PosOf(path, n), "tags is a list of tag names (TG3)"))
		return nil
	}
	var out []model.Tag
	for _, item := range n.Content {
		name, ok := source.Scalar(path, item, ds)
		if !ok {
			continue
		}
		if slices.ContainsFunc(out, func(t model.Tag) bool { return t.Name == name }) {
			*ds = append(*ds, diag.New(source.PosOf(path, item), "the tag %s is listed twice (TG4)", name))
			continue
		}
		out = append(out, model.Tag{Name: name, Pos: source.PosOf(path, item)})
	}
	return out
}

func decodeFields(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) []model.Field {
	fm, ok := source.AsMap(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.Field
	for _, p := range fm.Pairs {
		sm, ok := source.AsMap(path, p.Value, ds)
		if !ok {
			continue
		}
		f := model.Field{Name: p.Key, Pos: p.Pos}
		if on, ok := source.Get(sm, "optional"); ok {
			f.Optional, _ = source.Bool(path, on, ds)
		}
		f.Schema = decodeSchema(path, sm, scope, []string{"optional"}, ds)
		out = append(out, f)
	}
	return out
}

var schemaKeys = []string{"type", "enum", "const", "default", "pattern", "minLength", "maxLength", "minimum", "maximum", "format", "items", "properties", "oneOf"}

var droppedKeys = map[string]string{
	"required":             "every field and part is required unless it says optional: true",
	"$ref":                 "name an imported value type with type:",
	"null":                 "there is no null; leave the value out",
	"anyOf":                "only oneOf combines schemas; put other logic in a constraint",
	"allOf":                "only oneOf combines schemas; put other logic in a constraint",
	"not":                  "only oneOf combines schemas; put other logic in a constraint",
	"if":                   "conditional schemas are not available; use a constraint",
	"then":                 "conditional schemas are not available; use a constraint",
	"else":                 "conditional schemas are not available; use a constraint",
	"additionalProperties": "parts that are not listed are always rejected",
}

func decodeSchema(path string, m source.Map, scope map[string]model.Ref, extra []string, ds *[]diag.Diagnostic) model.Schema {
	s := model.Schema{Pos: m.Pos}
	for _, p := range m.Pairs {
		if why, ok := droppedKeys[p.Key]; ok {
			*ds = append(*ds, diag.New(p.Pos, "%s is not part of the language: %s (FD2)", p.Key, why))
			continue
		}
		if !slices.Contains(schemaKeys, p.Key) && !slices.Contains(extra, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not a schema keyword the language uses (FD1)", p.Key))
			continue
		}
		if slices.Contains(schemaKeys, p.Key) || p.Key == "optional" {
			s.Keywords = append(s.Keywords, p.Key)
		}
	}
	if n, ok := source.Get(m, "type"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			if slices.Contains(BuiltinTypes, name) {
				s.Type = name
			} else if ref, ok := resolve(scope, name, model.SortValueType, source.PosOf(path, n), "FD1", ds); ok {
				s.Ref, s.HasRef = ref.Q, true
			}
		}
	}
	if n, ok := source.Get(m, "enum"); ok {
		if n.Kind != yaml.SequenceNode || len(n.Content) == 0 {
			*ds = append(*ds, diag.New(source.PosOf(path, n), "enum is a list of allowed values (FD1)"))
		} else {
			for _, item := range n.Content {
				if v, ok := Value(path, item, ds); ok {
					s.Enum = append(s.Enum, v)
				}
			}
		}
	}
	s.Const = optionalValue(path, m, "const", ds)
	s.Default = optionalValue(path, m, "default", ds)
	s.Pattern, _ = source.GetScalar(m, "pattern", ds)
	if s.Pattern != "" {
		if _, err := regexp.Compile(s.Pattern); err != nil {
			*ds = append(*ds, diag.New(m.Pos, "pattern %q is not a valid regular expression: %v (FD4)", s.Pattern, err))
		}
	}
	s.MinLength = optionalInt(path, m, "minLength", ds)
	s.MaxLength = optionalInt(path, m, "maxLength", ds)
	s.Minimum = optionalNumber(path, m, "minimum", ds)
	s.Maximum = optionalNumber(path, m, "maximum", ds)
	s.Format, _ = source.GetScalar(m, "format", ds)
	if s.Format != "" && !slices.Contains([]string{"date", "date-time", "duration", "uri"}, s.Format) {
		*ds = append(*ds, diag.New(m.Pos, "format %q is not one of date, date-time, duration, uri (FD4)", s.Format))
	}
	if n, ok := source.Get(m, "items"); ok {
		if im, ok := source.AsMap(path, n, ds); ok {
			items := decodeSchema(path, im, scope, nil, ds)
			s.Items = &items
		}
	}
	if n, ok := source.Get(m, "properties"); ok {
		s.Properties = decodeFields(path, n, scope, ds)
	}
	if n, ok := source.Get(m, "oneOf"); ok {
		s.OneOf = decodeChoices(path, n, scope, ds)
	}
	return s
}

// decodeChoices reads the choices of a oneOf: each a schema without optional
// or default.
func decodeChoices(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) []model.Schema {
	items, ok := sequence(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.Schema
	for _, item := range items {
		cm, ok := source.AsMap(path, item, ds)
		if !ok {
			continue
		}
		for _, p := range cm.Pairs {
			if p.Key == "optional" || p.Key == "default" {
				*ds = append(*ds, diag.New(p.Pos, "a choice of oneOf has no %s (FD7)", p.Key))
			}
		}
		out = append(out, decodeSchema(path, cm, scope, []string{"optional"}, ds))
	}
	return out
}

// Value reads a literal value. There is no null: a value is always something.
func Value(path string, n *yaml.Node, ds *[]diag.Diagnostic) (model.Value, bool) {
	pos := source.PosOf(path, n)
	switch n.Kind {
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!null":
			*ds = append(*ds, diag.New(pos, "there is no null; leave the value out instead (FD6)"))
			return model.Value{}, false
		case "!!bool":
			return model.Value{Kind: model.ValueBoolean, Bool: n.Value == "true"}, true
		case "!!int":
			r, ok := new(big.Rat).SetString(n.Value)
			if !ok {
				*ds = append(*ds, diag.New(pos, "%s is not a number", n.Value))
				return model.Value{}, false
			}
			return model.Value{Kind: model.ValueInteger, Num: r, Text: n.Value}, true
		case "!!float":
			r, ok := new(big.Rat).SetString(n.Value)
			if !ok {
				*ds = append(*ds, diag.New(pos, "%s is not an exact number", n.Value))
				return model.Value{}, false
			}
			return model.Value{Kind: model.ValueNumber, Num: r, Text: n.Value}, true
		default:
			return model.Value{Kind: model.ValueString, Text: n.Value}, true
		}
	case yaml.SequenceNode:
		v := model.Value{Kind: model.ValueArray}
		for _, item := range n.Content {
			iv, ok := Value(path, item, ds)
			if !ok {
				return model.Value{}, false
			}
			v.Items = append(v.Items, iv)
		}
		return v, true
	case yaml.MappingNode:
		m, ok := source.AsMap(path, n, ds)
		if !ok {
			return model.Value{}, false
		}
		v := model.Value{Kind: model.ValueObject, Props: map[string]model.Value{}}
		for _, p := range m.Pairs {
			pv, ok := Value(path, p.Value, ds)
			if !ok {
				return model.Value{}, false
			}
			v.Props[p.Key] = pv
		}
		return v, true
	default:
		*ds = append(*ds, diag.New(pos, "expected a value"))
		return model.Value{}, false
	}
}

func optionalValue(path string, m source.Map, key string, ds *[]diag.Diagnostic) *model.Value {
	n, ok := source.Get(m, key)
	if !ok {
		return nil
	}
	v, ok := Value(path, n, ds)
	if !ok {
		return nil
	}
	return &v
}

func optionalInt(path string, m source.Map, key string, ds *[]diag.Diagnostic) *int {
	n, ok := source.Get(m, key)
	if !ok {
		return nil
	}
	i, err := strconv.Atoi(n.Value)
	if n.Kind != yaml.ScalarNode || err != nil || i < 0 {
		*ds = append(*ds, diag.New(source.PosOf(path, n), "%s must be a whole number of zero or more (FD1)", key))
		return nil
	}
	return &i
}

func optionalNumber(path string, m source.Map, key string, ds *[]diag.Diagnostic) *big.Rat {
	n, ok := source.Get(m, key)
	if !ok {
		return nil
	}
	r, ok := new(big.Rat).SetString(n.Value)
	if n.Kind != yaml.ScalarNode || !ok {
		*ds = append(*ds, diag.New(source.PosOf(path, n), "%s must be a number (FD1)", key))
		return nil
	}
	return r
}

func requiredText(path string, m source.Map, key, owner, rule string, ds *[]diag.Diagnostic) string {
	n, ok := source.Get(m, key)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "%s needs %s: (%s)", owner, key, rule))
		return ""
	}
	s, _ := source.Scalar(path, n, ds)
	return s
}

func sequence(path string, n *yaml.Node, ds *[]diag.Diagnostic) ([]*yaml.Node, bool) {
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(source.PosOf(path, n), "expected a list"))
		return nil, false
	}
	return n.Content, true
}

// resolve looks a name up in a file's scope. An empty sort accepts any.
func resolve(scope map[string]model.Ref, name string, sort model.Sort, pos diag.Pos, rule string, ds *[]diag.Diagnostic) (model.Ref, bool) {
	ref, ok := scope[name]
	if !ok {
		*ds = append(*ds, diag.New(pos, "%s is not defined in this file or imported (FL9)", name))
		return model.Ref{}, false
	}
	if sort != "" && ref.Sort != sort {
		*ds = append(*ds, diag.New(pos, "%s is a %s, not a %s (%s)", name, ref.Sort, sort, rule))
		return model.Ref{}, false
	}
	return ref, true
}

// Resolve is resolve for other packages: test files resolve their names the
// same way definitions do.
func Resolve(scope map[string]model.Ref, name string, sort model.Sort, pos diag.Pos, rule string, ds *[]diag.Diagnostic) (model.Ref, bool) {
	return resolve(scope, name, sort, pos, rule, ds)
}
