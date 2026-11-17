// SPDX-License-Identifier: Apache-2.0

// Package facts reads a test file: its facts, checked against their kinds,
// and the tests and questions written beside them.
package facts

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/engine"
	"github.com/openfjord/of-lang/cli/internal/load"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
	"github.com/openfjord/of-lang/cli/internal/value"
	"github.com/openfjord/of-lang/cli/internal/world"
)

// File is one test file once read.
type File struct {
	Path      string
	Scope     map[string]model.Ref
	Graph     engine.Graph
	Facts     int
	Tests     []source.Entry
	Questions []source.Entry
}

// Read reads a test file against the root package of a loaded world, as the
// same version of a pending change. The error is for a file that cannot be
// read at all.
func Read(l load.Loaded, path string) (File, []diag.Diagnostic, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return File{}, nil, fmt.Errorf("resolve %s: %w", path, err)
	}
	raw, ds, err := source.ReadFile(path)
	if err != nil {
		return File{}, nil, err
	}
	sf, sds := source.Select(raw, l.Version)
	ds = append(ds, sds...)
	pkg := l.World.Packages[l.World.Root]
	scope, sds := load.Scope(l.World.Packages, pkg, sf)
	ds = append(ds, sds...)
	out := File{Path: path, Scope: scope, Graph: engine.NewGraph()}
	var links []factLink
	for _, e := range sf.Entries {
		if _, ok := pkg.Defs[e.Name]; ok {
			ds = append(ds, diag.New(e.Pos, "%s is already a definition in package %s; statements in a test file need their own names (FL8)", e.Name, pkg.Name))
			continue
		}
		switch model.Sort(e.Is) {
		case model.SortFact:
			m, ok := source.AsMap(path, e.Node, &ds)
			if !ok {
				continue
			}
			out.Facts++
			links = append(links, readFact(l.World, &out, e, m, &ds)...)
		case model.SortTest:
			out.Tests = append(out.Tests, e)
		case model.SortQuestion:
			out.Questions = append(out.Questions, e)
		default:
			ds = append(ds, diag.New(e.Pos, "%s says is: %s; a statement in tests/ is a Fact, a Test or a Question (FL4)", e.Name, e.Is))
		}
	}
	for _, fl := range links {
		if !engine.HasNode(out.Graph, fl.to) {
			ds = append(ds, diag.New(fl.pos, "%s is not a fact in this file (FC5)", fl.to))
			continue
		}
		from, to := fl.from, fl.to
		if fl.ref.Inverted {
			from, to = to, from
		}
		engine.Add(&out.Graph, engine.LinkFact(from, fl.ref.Q, to), engine.Reason{Type: engine.ReasonStated, Pos: fl.pos})
	}
	return out, ds, nil
}

type factLink struct {
	from string
	ref  model.Ref
	to   string
	pos  diag.Pos
}

// readFact records a fact's kind and fields, and returns its links to add
// once every fact in the file is known.
func readFact(w model.World, out *File, e source.Entry, m source.Map, ds *[]diag.Diagnostic) []factLink {
	path := out.Path
	for _, p := range m.Pairs {
		if !slices.Contains([]string{"is", "kind", "fields", "links"}, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not allowed here; a fact has is, kind, fields and links (FL2)", p.Key))
		}
	}
	// The id is known even when the kind is not, so links to it report
	// nothing more.
	engine.AddNode(&out.Graph, e.Name, e.Pos)
	kindName, ok := source.GetScalar(m, "kind", ds)
	if !ok {
		*ds = append(*ds, diag.New(e.Pos, "fact %s needs kind: (FC1)", e.Name))
		return nil
	}
	kref, ok := load.Resolve(out.Scope, kindName, model.SortKind, e.Pos, "FC1", ds)
	if !ok {
		return nil
	}
	if world.IsUnion(w, kref.Q) {
		*ds = append(*ds, diag.New(e.Pos, "%s is a union; a fact names one of its listed kinds (KD10)", kindName))
		return nil
	}
	kindFact := engine.KindFact(e.Name, kref.Q)
	engine.Add(&out.Graph, kindFact, engine.Reason{Type: engine.ReasonStated, Pos: e.Pos})
	if n, ok := source.Get(m, "fields"); ok {
		if fm, ok := source.AsMap(path, n, ds); ok {
			for _, p := range fm.Pairs {
				f, _, ok := world.Field(w, kref.Q, p.Key)
				if !ok {
					*ds = append(*ds, diag.New(p.Pos, "%s is not a field of %s (FC2)", p.Key, kindName))
					continue
				}
				v, ok := load.Value(path, p.Value, ds)
				if !ok {
					continue
				}
				for _, msg := range value.Validate(w, f.Schema, v) {
					*ds = append(*ds, diag.New(p.Pos, "%s: %s (FC4)", p.Key, msg))
				}
				engine.Add(&out.Graph, engine.FieldFact(e.Name, p.Key, v), engine.Reason{Type: engine.ReasonStated, Pos: p.Pos})
			}
		}
	}
	engine.AddDefaults(w, &out.Graph, e.Name, kref.Q, engine.Reason{Type: engine.ReasonDefault})
	var links []factLink
	n, ok := source.Get(m, "links")
	if !ok {
		return nil
	}
	lm, ok := source.AsMap(path, n, ds)
	if !ok {
		return nil
	}
	for _, p := range lm.Pairs {
		ref, ok := out.Scope[p.Key]
		if !ok || ref.Sort != model.SortLink {
			*ds = append(*ds, diag.New(p.Pos, "%s is not a link that this file imports (FC2)", p.Key))
			continue
		}
		start, _ := world.LinkEnds(w, ref.Q, ref.Inverted)
		if !world.IsA(w, kref.Q, start) {
			*ds = append(*ds, diag.New(p.Pos, "%s starts at a %s, so a %s cannot have it (FC2)", p.Key, start.Name, kindName))
			continue
		}
		targets, ok := source.Strings(path, p.Value, ds)
		if !ok {
			continue
		}
		for _, t := range targets {
			links = append(links, factLink{from: e.Name, ref: ref, to: t, pos: p.Pos})
		}
	}
	return links
}
