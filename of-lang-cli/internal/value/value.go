// SPDX-License-Identifier: Apache-2.0

// Package value compares literal values, renders them, and checks them
// against the language's subset of JSON Schema.
package value

import (
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/openfjord/of-lang/internal/model"
)

// Text makes a string value, which is how a thing's id is compared.
func Text(s string) model.Value {
	return model.Value{Kind: model.ValueString, Text: s}
}

// Integer makes an integer value.
func Integer(n int) model.Value {
	return model.Value{Kind: model.ValueInteger, Num: big.NewRat(int64(n), 1), Text: fmt.Sprint(n)}
}

// Format renders a value the way it would be written in YAML.
func Format(v model.Value) string {
	switch v.Kind {
	case model.ValueString:
		return v.Text
	case model.ValueInteger, model.ValueNumber:
		if v.Text != "" {
			return v.Text
		}
		if v.Num.IsInt() {
			return v.Num.Num().String()
		}
		return strings.TrimRight(v.Num.FloatString(20), "0")
	case model.ValueBoolean:
		return fmt.Sprint(v.Bool)
	case model.ValueArray:
		parts := make([]string, 0, len(v.Items))
		for _, it := range v.Items {
			parts = append(parts, Format(it))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case model.ValueObject:
		keys := make([]string, 0, len(v.Props))
		for k := range v.Props {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+": "+Format(v.Props[k]))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return ""
	}
}

// Plain turns a value into plain Go data for JSON output.
func Plain(v model.Value) any {
	switch v.Kind {
	case model.ValueString:
		return v.Text
	case model.ValueInteger, model.ValueNumber:
		return Format(v)
	case model.ValueBoolean:
		return v.Bool
	case model.ValueArray:
		out := make([]any, 0, len(v.Items))
		for _, it := range v.Items {
			out = append(out, Plain(it))
		}
		return out
	case model.ValueObject:
		out := make(map[string]any, len(v.Props))
		for k, p := range v.Props {
			out[k] = Plain(p)
		}
		return out
	default:
		return nil
	}
}

func isNumber(v model.Value) bool {
	return v.Kind == model.ValueInteger || v.Kind == model.ValueNumber
}

// Equal says whether two values are the same. Integers and numbers compare by
// their value.
func Equal(a, b model.Value) bool {
	if isNumber(a) && isNumber(b) {
		return a.Num.Cmp(b.Num) == 0
	}
	if a.Kind != b.Kind {
		return false
	}
	switch a.Kind {
	case model.ValueString:
		return a.Text == b.Text
	case model.ValueBoolean:
		return a.Bool == b.Bool
	case model.ValueArray:
		return slices.EqualFunc(a.Items, b.Items, Equal)
	case model.ValueObject:
		if len(a.Props) != len(b.Props) {
			return false
		}
		for k, av := range a.Props {
			bv, ok := b.Props[k]
			if !ok || !Equal(av, bv) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// Compare orders two numbers, or two strings, which is how dates in ISO form
// compare too. ok is false for values with no order.
func Compare(a, b model.Value) (int, bool) {
	switch {
	case isNumber(a) && isNumber(b):
		return a.Num.Cmp(b.Num), true
	case a.Kind == model.ValueString && b.Kind == model.ValueString:
		return strings.Compare(a.Text, b.Text), true
	default:
		return 0, false
	}
}

var durationPattern = regexp.MustCompile(`^P(\d+Y)?(\d+M)?(\d+W)?(\d+D)?(T(\d+H)?(\d+M)?(\d+(\.\d+)?S)?)?$`)

// Validate checks a value against a schema and returns what is wrong with it.
func Validate(w model.World, s model.Schema, v model.Value) []string {
	var errs []string
	if s.HasRef {
		if vt, ok := w.ValueTypes[s.Ref]; ok {
			for _, e := range Validate(w, vt.Schema, v) {
				errs = append(errs, e+" (from "+s.Ref.Name+")")
			}
		}
	}
	if s.Type != "" && !hasType(s.Type, v) {
		return append(errs, fmt.Sprintf("%s is not of type %s", Format(v), s.Type))
	}
	if len(s.Enum) > 0 && !slices.ContainsFunc(s.Enum, func(e model.Value) bool { return Equal(e, v) }) {
		allowed := make([]string, 0, len(s.Enum))
		for _, e := range s.Enum {
			allowed = append(allowed, Format(e))
		}
		errs = append(errs, fmt.Sprintf("%s is not one of %s", Format(v), strings.Join(allowed, ", ")))
	}
	if s.Const != nil && !Equal(*s.Const, v) {
		errs = append(errs, fmt.Sprintf("%s must be %s", Format(v), Format(*s.Const)))
	}
	if v.Kind == model.ValueString {
		errs = append(errs, validateString(s, v.Text)...)
	}
	if isNumber(v) {
		if s.Minimum != nil && v.Num.Cmp(s.Minimum) < 0 {
			errs = append(errs, fmt.Sprintf("%s is below the minimum %s", Format(v), s.Minimum.RatString()))
		}
		if s.Maximum != nil && v.Num.Cmp(s.Maximum) > 0 {
			errs = append(errs, fmt.Sprintf("%s is above the maximum %s", Format(v), s.Maximum.RatString()))
		}
	}
	if v.Kind == model.ValueArray && s.Items != nil {
		for i, it := range v.Items {
			for _, e := range Validate(w, *s.Items, it) {
				errs = append(errs, fmt.Sprintf("item %d: %s", i+1, e))
			}
		}
	}
	if v.Kind == model.ValueObject && len(s.Properties) > 0 {
		errs = append(errs, validateObject(w, s.Properties, v)...)
	}
	return errs
}

func validateString(s model.Schema, text string) []string {
	var errs []string
	n := utf8.RuneCountInString(text)
	if s.MinLength != nil && n < *s.MinLength {
		errs = append(errs, fmt.Sprintf("%q is shorter than %d characters", text, *s.MinLength))
	}
	if s.MaxLength != nil && n > *s.MaxLength {
		errs = append(errs, fmt.Sprintf("%q is longer than %d characters", text, *s.MaxLength))
	}
	if s.Pattern != "" {
		if re, err := regexp.Compile(s.Pattern); err == nil && !re.MatchString(text) {
			errs = append(errs, fmt.Sprintf("%q does not match %s", text, s.Pattern))
		}
	}
	if s.Format != "" && !hasFormat(s.Format, text) {
		errs = append(errs, fmt.Sprintf("%q is not a %s", text, s.Format))
	}
	return errs
}

func validateObject(w model.World, props []model.Field, v model.Value) []string {
	var errs []string
	for _, p := range props {
		pv, ok := v.Props[p.Name]
		if !ok {
			if !p.Optional && p.Schema.Default == nil {
				errs = append(errs, fmt.Sprintf("%s is missing", p.Name))
			}
			continue
		}
		for _, e := range Validate(w, p.Schema, pv) {
			errs = append(errs, p.Name+": "+e)
		}
	}
	for k := range v.Props {
		if !slices.ContainsFunc(props, func(f model.Field) bool { return f.Name == k }) {
			errs = append(errs, fmt.Sprintf("%s is not a listed part", k))
		}
	}
	slices.Sort(errs)
	return errs
}

func hasType(t string, v model.Value) bool {
	switch t {
	case "string":
		return v.Kind == model.ValueString
	case "integer":
		return isNumber(v) && v.Num.IsInt()
	case "number":
		return isNumber(v)
	case "boolean":
		return v.Kind == model.ValueBoolean
	case "array":
		return v.Kind == model.ValueArray
	case "object":
		return v.Kind == model.ValueObject
	default:
		return false
	}
}

func hasFormat(format, text string) bool {
	switch format {
	case "date":
		_, err := time.Parse(time.DateOnly, text)
		return err == nil
	case "date-time":
		_, err := time.Parse(time.RFC3339, text)
		return err == nil
	case "duration":
		return text != "P" && !strings.HasSuffix(text, "T") && durationPattern.MatchString(text)
	default:
		return true
	}
}
