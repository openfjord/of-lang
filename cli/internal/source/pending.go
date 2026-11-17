// SPDX-License-Identifier: Apache-2.0

package source

import (
	"slices"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
)

// Version picks one version of a package that has a pending change.
type Version string

// The two versions of a package.
const (
	WithChange    Version = "with the change"
	WithoutChange Version = "without the change"
)

// The values of pending.
const (
	markAdded   = "added"
	markRemoved = "removed"
	markUpdated = "updated"
)

// HasChange says whether any statement of a file carries a pending mark.
func HasChange(f File) bool {
	return slices.ContainsFunc(f.Entries, func(e Entry) bool {
		m, ok := AsMap(f.Path, e.Node, new([]diag.Diagnostic))
		if !ok {
			return false
		}
		_, ok = Get(m, "pending")
		return ok
	})
}

// Select returns a file as one version of its package reads it. With the
// change, it leaves out removed statements; without it, it leaves out added
// statements and reads each updated statement as its was version. Either
// way, the statements it returns carry no pending or was.
func Select(f File, v Version) (File, []diag.Diagnostic) {
	out := File{Path: f.Path, Imports: f.Imports}
	var ds []diag.Diagnostic
	for _, e := range f.Entries {
		m, ok := AsMap(f.Path, e.Node, &ds)
		if !ok {
			continue
		}
		mark, was := pendingMarks(f.Path, e, m, &ds)
		switch {
		case v == WithChange && mark == markRemoved, v == WithoutChange && mark == markAdded:
			continue
		case v == WithoutChange && mark == markUpdated && was != nil:
			e.Node = wasNode(e.Node, was)
		default:
			e.Node = withoutKeys(e.Node, "pending", "was")
		}
		out.Entries = append(out.Entries, e)
	}
	return out, ds
}

func pendingMarks(path string, e Entry, m Map, ds *[]diag.Diagnostic) (string, *yaml.Node) {
	mark := ""
	if n, ok := Get(m, "pending"); ok {
		s, ok := Scalar(path, n, ds)
		if ok && !slices.Contains([]string{markAdded, markRemoved, markUpdated}, s) {
			*ds = append(*ds, diag.New(PosOf(path, n), "pending is added, removed or updated, not %q (PD1)", s))
		}
		mark = s
	}
	was, hasWas := Get(m, "was")
	switch {
	case mark == markUpdated && !hasWas:
		*ds = append(*ds, diag.New(e.Pos, "%s has pending: updated, so it needs was: with the version before the change (PD2)", e.Name))
		return mark, nil
	case mark != markUpdated && hasWas:
		*ds = append(*ds, diag.New(PosOf(path, was), "was belongs only on a statement with pending: updated (PD2)"))
		return mark, nil
	case !hasWas:
		return mark, nil
	}
	wm, ok := AsMap(path, was, ds)
	if !ok {
		return mark, nil
	}
	for _, p := range wm.Pairs {
		switch p.Key {
		case "is":
			*ds = append(*ds, diag.New(p.Pos, "was leaves out is: an update does not change what a statement is (PD3)"))
		case "pending", "was":
			*ds = append(*ds, diag.New(p.Pos, "was holds the statement without pending and was (PD2)"))
		}
	}
	return mark, was
}

// wasNode builds the old version of an updated statement: its is, then the
// keys under was.
func wasNode(n, was *yaml.Node) *yaml.Node {
	out := &yaml.Node{Kind: yaml.MappingNode, Tag: n.Tag, Line: n.Line, Column: n.Column}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == "is" {
			out.Content = append(out.Content, n.Content[i], n.Content[i+1])
		}
	}
	for i := 0; i+1 < len(was.Content); i += 2 {
		if !slices.Contains([]string{"is", "pending", "was"}, was.Content[i].Value) {
			out.Content = append(out.Content, was.Content[i], was.Content[i+1])
		}
	}
	return out
}

func withoutKeys(n *yaml.Node, keys ...string) *yaml.Node {
	if n.Kind != yaml.MappingNode {
		return n
	}
	out := *n
	out.Content = nil
	for i := 0; i+1 < len(n.Content); i += 2 {
		if !slices.Contains(keys, n.Content[i].Value) {
			out.Content = append(out.Content, n.Content[i], n.Content[i+1])
		}
	}
	return &out
}
