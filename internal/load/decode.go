// SPDX-License-Identifier: Apache-2.0

package load

import (
	"math/big"
	"regexp"
	"slices"
	"strconv"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/internal/diag"
	"github.com/openfjord/of-lang/internal/model"
	"github.com/openfjord/of-lang/internal/source"
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
	case model.SortConstraint, model.SortRule:
		w.Statements[q] = Statement(q, model.Sort(e.Is), path, m, scope, "says", nil, &ds)
		w.Order = append(w.Order, q)
	case model.SortMapping:
		w.Mappings[q] = decodeMapping(q, path, m, scope, &ds)
	default:
		// Reported when the package's files were collected.
	}
	return ds
}

func decodeKind(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.Kind {
	source.Allow(m, []string{"is", "definition", "extends", "examples", "notExamples", "fields", "links"}, ds)
	k := model.Kind{Q: q, Pos: m.Pos}
	k.Definition = requiredText(path, m, "definition", q.Name, ds)
	if n, ok := source.Get(m, "extends"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			if ref, ok := resolve(scope, name, model.SortKind, source.PosOf(path, n), ds); ok {
				k.Extends = &ref.Q
			}
		}
	}
	k.Examples = optionalStrings(path, m, "examples", ds)
	k.NotExamples = optionalStrings(path, m, "notExamples", ds)
	if n, ok := source.Get(m, "fields"); ok {
		k.Fields = decodeFields(path, n, scope, ds)
	}
	if n, ok := source.Get(m, "links"); ok {
		lm, ok := source.AsMap(path, n, ds)
		if ok {
			for _, p := range lm.Pairs {
				ref, ok := resolve(scope, p.Key, model.SortLink, p.Pos, ds)
				if !ok {
					continue
				}
				text, ok := source.Scalar(path, p.Value, ds)
				if !ok {
					continue
				}
				c, ok := parseCount(text)
				if !ok {
					*ds = append(*ds, diag.New(p.Pos, "%q is not a count; write exactly N, at least N, at most N or any number", text))
					continue
				}
				k.Links = append(k.Links, model.LinkCount{Ref: ref, Name: p.Key, Pos: p.Pos, Count: c})
			}
		}
	}
	return k
}

func parseCount(text string) (model.Count, bool) {
	if text == "any number" {
		return model.Count{Min: 0, Max: -1}, true
	}
	for _, form := range []struct {
		prefix   string
		min, max func(int) int
	}{
		{"exactly ", func(n int) int { return n }, func(n int) int { return n }},
		{"at least ", func(n int) int { return n }, func(int) int { return -1 }},
		{"at most ", func(int) int { return 0 }, func(n int) int { return n }},
	} {
		if len(text) > len(form.prefix) && text[:len(form.prefix)] == form.prefix {
			n, err := strconv.Atoi(text[len(form.prefix):])
			if err != nil || n < 0 {
				return model.Count{}, false
			}
			return model.Count{Min: form.min(n), Max: form.max(n)}, true
		}
	}
	return model.Count{}, false
}

func decodeLink(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.Link {
	source.Allow(m, []string{"is", "from", "to", "reads", "definition", "under", "inverse", "chains"}, ds)
	l := model.Link{Q: q, Pos: m.Pos}
	l.From = kindRef(path, m, "from", q.Name, scope, ds)
	l.To = kindRef(path, m, "to", q.Name, scope, ds)
	l.Reads = requiredText(path, m, "reads", q.Name, ds)
	l.Definition = requiredText(path, m, "definition", q.Name, ds)
	if n, ok := source.Get(m, "under"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			if ref, ok := resolve(scope, name, model.SortLink, source.PosOf(path, n), ds); ok {
				if ref.Inverted {
					*ds = append(*ds, diag.New(source.PosOf(path, n), "%s is an inverse name; under: names a declared link", name))
				}
				l.Under = &ref.Q
			}
		}
	}
	l.Inverse, _ = source.GetScalar(m, "inverse", ds)
	if n, ok := source.Get(m, "chains"); ok {
		l.Chains, _ = source.Bool(path, n, ds)
	}
	return l
}

func kindRef(path string, m source.Map, key, owner string, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.QName {
	n, ok := source.Get(m, key)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "%s needs %s:", owner, key))
		return model.QName{}
	}
	name, ok := source.Scalar(path, n, ds)
	if !ok {
		return model.QName{}
	}
	ref, ok := resolve(scope, name, model.SortKind, source.PosOf(path, n), ds)
	if !ok {
		return model.QName{}
	}
	return ref.Q
}

func decodeValueType(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.ValueType {
	if slices.Contains(BuiltinTypes, q.Name) {
		*ds = append(*ds, diag.New(m.Pos, "%s is a built-in type and cannot name a value type", q.Name))
	}
	vt := model.ValueType{Q: q, Pos: m.Pos}
	vt.Definition = requiredText(path, m, "definition", q.Name, ds)
	vt.Schema = decodeSchema(path, m, scope, []string{"is", "definition"}, ds)
	return vt
}

func decodeMapping(q model.QName, path string, m source.Map, scope map[string]model.Ref, ds *[]diag.Diagnostic) model.Mapping {
	source.Allow(m, []string{"is", "term", "equivalent"}, ds)
	mp := model.Mapping{Q: q, Pos: m.Pos}
	if n, ok := source.Get(m, "term"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			mp.Term, _ = resolve(scope, name, "", source.PosOf(path, n), ds)
		}
	} else {
		*ds = append(*ds, diag.New(m.Pos, "%s needs term: naming what it maps", q.Name))
	}
	mp.Equivalent = requiredText(path, m, "equivalent", q.Name, ds)
	return mp
}

func decodeFields(path string, n *yaml.Node, scope map[string]model.Ref, ds *[]diag.Diagnostic) []model.Field {
	fm, ok := source.AsMap(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.Field
	for _, p := range fm.Pairs {
		if p.Key == "kind" {
			*ds = append(*ds, diag.New(p.Pos, "kind is reserved: facts use it to name their kind"))
			continue
		}
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

var schemaKeys = []string{"type", "enum", "const", "default", "pattern", "minLength", "maxLength", "minimum", "maximum", "format", "items", "properties"}

var droppedKeys = map[string]string{
	"required":             "every field is required unless it says optional: true",
	"$ref":                 "name an imported value type with type:",
	"oneOf":                "combining schemas is not available; use a constraint",
	"anyOf":                "combining schemas is not available; use a constraint",
	"allOf":                "combining schemas is not available; use a constraint",
	"not":                  "combining schemas is not available; use a constraint",
	"if":                   "conditional schemas are not available; use a constraint",
	"then":                 "conditional schemas are not available; use a constraint",
	"else":                 "conditional schemas are not available; use a constraint",
	"additionalProperties": "unlisted parts are always rejected",
}

func decodeSchema(path string, m source.Map, scope map[string]model.Ref, extra []string, ds *[]diag.Diagnostic) model.Schema {
	for _, p := range m.Pairs {
		if why, ok := droppedKeys[p.Key]; ok {
			*ds = append(*ds, diag.New(p.Pos, "%s is not part of the language: %s", p.Key, why))
			continue
		}
		if !slices.Contains(schemaKeys, p.Key) && !slices.Contains(extra, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not a schema keyword the language uses", p.Key))
		}
	}
	s := model.Schema{Pos: m.Pos}
	if n, ok := source.Get(m, "type"); ok {
		if name, ok := source.Scalar(path, n, ds); ok {
			if slices.Contains(BuiltinTypes, name) {
				s.Type = name
			} else if ref, ok := resolve(scope, name, model.SortValueType, source.PosOf(path, n), ds); ok {
				s.Ref, s.HasRef = ref.Q, true
			}
		}
	}
	if n, ok := source.Get(m, "enum"); ok {
		if n.Kind != yaml.SequenceNode || len(n.Content) == 0 {
			*ds = append(*ds, diag.New(source.PosOf(path, n), "enum is a list of allowed values"))
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
			*ds = append(*ds, diag.New(m.Pos, "pattern %q is not a valid regular expression: %v", s.Pattern, err))
		}
	}
	s.MinLength = optionalInt(path, m, "minLength", ds)
	s.MaxLength = optionalInt(path, m, "maxLength", ds)
	s.Minimum = optionalNumber(path, m, "minimum", ds)
	s.Maximum = optionalNumber(path, m, "maximum", ds)
	s.Format, _ = source.GetScalar(m, "format", ds)
	if s.Format != "" && !slices.Contains([]string{"date", "date-time", "duration"}, s.Format) {
		*ds = append(*ds, diag.New(m.Pos, "format %q is not one of date, date-time, duration", s.Format))
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
	return s
}

// Value reads a literal value. There is no null: a value is always something.
func Value(path string, n *yaml.Node, ds *[]diag.Diagnostic) (model.Value, bool) {
	pos := source.PosOf(path, n)
	switch n.Kind {
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!null":
			*ds = append(*ds, diag.New(pos, "there is no null; leave the field out instead"))
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
		*ds = append(*ds, diag.New(source.PosOf(path, n), "%s must be a whole number of zero or more", key))
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
		*ds = append(*ds, diag.New(source.PosOf(path, n), "%s must be a number", key))
		return nil
	}
	return r
}

func requiredText(path string, m source.Map, key, owner string, ds *[]diag.Diagnostic) string {
	n, ok := source.Get(m, key)
	if !ok {
		*ds = append(*ds, diag.New(m.Pos, "%s needs %s:", owner, key))
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
func resolve(scope map[string]model.Ref, name string, sort model.Sort, pos diag.Pos, ds *[]diag.Diagnostic) (model.Ref, bool) {
	ref, ok := scope[name]
	if !ok {
		*ds = append(*ds, diag.New(pos, "%s is not defined in this file or imported", name))
		return model.Ref{}, false
	}
	if sort != "" && ref.Sort != sort {
		*ds = append(*ds, diag.New(pos, "%s is a %s, not a %s", name, ref.Sort, sort))
		return model.Ref{}, false
	}
	return ref, true
}

// Resolve is resolve for other packages: test files resolve their names the
// same way definitions do.
func Resolve(scope map[string]model.Ref, name string, sort model.Sort, pos diag.Pos, ds *[]diag.Diagnostic) (model.Ref, bool) {
	return resolve(scope, name, sort, pos, ds)
}
