// Package diag holds the problems the tool reports about a user's files, so
// every stage can report several at once instead of stopping at the first.
package diag

import (
	"cmp"
	"fmt"
	"slices"
)

// Pos is a place in a file, kept so every problem points at what caused it.
type Pos struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// Diagnostic is one problem found in a user's files.
type Diagnostic struct {
	Pos     Pos    `json:"pos"`
	Message string `json:"message"`
}

// New builds a diagnostic with a formatted message.
func New(pos Pos, format string, args ...any) Diagnostic {
	return Diagnostic{Pos: pos, Message: fmt.Sprintf(format, args...)}
}

// Format renders a diagnostic the way compilers do, so editors can jump to it.
func Format(d Diagnostic) string {
	if d.Pos.File == "" {
		return d.Message
	}
	if d.Pos.Line == 0 {
		return fmt.Sprintf("%s: %s", d.Pos.File, d.Message)
	}
	return fmt.Sprintf("%s:%d:%d: %s", d.Pos.File, d.Pos.Line, d.Pos.Column, d.Message)
}

// Sorted returns the diagnostics in file and line order, without duplicates,
// so output is stable from run to run.
func Sorted(ds []Diagnostic) []Diagnostic {
	out := slices.Clone(ds)
	slices.SortFunc(out, func(a, b Diagnostic) int {
		return cmp.Or(
			cmp.Compare(a.Pos.File, b.Pos.File),
			cmp.Compare(a.Pos.Line, b.Pos.Line),
			cmp.Compare(a.Pos.Column, b.Pos.Column),
			cmp.Compare(a.Message, b.Message),
		)
	})
	return slices.Compact(out)
}
