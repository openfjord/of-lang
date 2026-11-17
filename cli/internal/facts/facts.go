// SPDX-License-Identifier: Apache-2.0

// Package facts reads a test file: its facts, checked against their kinds,
// and the tests and questions written beside them.
package facts

import (
	"fmt"
	"path/filepath"

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
	Tests     []source.Entry
	Questions []source.Entry
}

// Read reads a test file against the root package of a loaded world. The
// error is for a file that cannot be read at all.
func Read(l load.Loaded, path string) (File, []diag.Diagnostic, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return File{}, nil, fmt.Errorf("resolve %s: %w", path, err)
	}
	sf, ds, err := source.ReadFile(path)
	if err != nil {
		return File{}, nil, err
	}
	pkg := l.World.Packages[l.World.Root]
	scope, sds := load.Scope(l.World.Packages, pkg, sf)
	ds = append(ds, sds...)
	out := File{Path: path, Scope: scope, Graph: engine.NewGraph()}
	var links []factLink
	seen := map[string]bool{}
	for _, e := range sf.Entries {
		if _, ok := pkg.Defs[e.Name]; ok {
			ds = append(ds, diag.New(e.Pos, "%s is already a definition in package %s; entries in tests need their own names", e.Name, pkg.Name))
			continue
		}
		if seen[e.Name] {
			ds = append(ds, diag.New(e.Pos, "%s appears twice", e.Name))
			continue
		}
		seen[e.Name] = true
		switch model.Sort(e.Is) {
		case model.SortFact:
			m, ok := source.AsMap(path, e.Node, &ds)
			if !ok {
				continue
			}
			links = append(links, readFact(l.World, &out, e, m, &ds)...)
		case model.SortTest:
			out.Tests = append(out.Tests, e)
		case model.SortQuestion:
			out.Questions = append(out.Questions, e)
		default:
			ds = append(ds, diag.New(e.Pos, "%s says is: %s; a test file holds Fact, Test and Question entries", e.Name, e.Is))
		}
	}
	for _, pl := range links {
		if !engine.HasNode(out.Graph, pl.to) {
			ds = append(ds, diag.New(pl.pos, "%s is not a fact in this file", pl.to))
			continue
		}
		from, to := pl.from, pl.to
		if pl.ref.Inverted {
			from, to = to, from
		}
		engine.Add(&out.Graph, engine.LinkFact(from, pl.ref.Q, to), engine.Reason{Type: engine.ReasonAsserted, Pos: pl.pos})
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
	kindName, ok := source.GetScalar(m, "kind", ds)
	if !ok {
		*ds = append(*ds, diag.New(e.Pos, "fact %s needs kind:", e.Name))
		return nil
	}
	kref, ok := load.Resolve(out.Scope, kindName, model.SortKind, e.Pos, ds)
	if !ok {
		return nil
	}
	engine.AddNode(&out.Graph, e.Name, e.Pos)
	engine.Add(&out.Graph, engine.KindFact(e.Name, kref.Q), engine.Reason{Type: engine.ReasonAsserted, Pos: e.Pos})
	var links []factLink
	given := map[string]bool{}
	for _, p := range m.Pairs {
		if p.Key == "is" || p.Key == "kind" {
			continue
		}
		if f, ok := world.Field(w, kref.Q, p.Key); ok {
			given[p.Key] = true
			v, ok := load.Value(path, p.Value, ds)
			if !ok {
				continue
			}
			for _, msg := range value.Validate(w, f.Schema, v) {
				*ds = append(*ds, diag.New(p.Pos, "%s: %s", p.Key, msg))
			}
			engine.Add(&out.Graph, engine.FieldFact(e.Name, p.Key, v), engine.Reason{Type: engine.ReasonAsserted, Pos: p.Pos})
			continue
		}
		ref, ok := out.Scope[p.Key]
		if !ok || ref.Sort != model.SortLink {
			*ds = append(*ds, diag.New(p.Pos, "%s is not a field of %s, nor an imported link", p.Key, kindName))
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
	for _, f := range world.Fields(w, kref.Q) {
		if given[f.Name] {
			continue
		}
		switch {
		case f.Schema.Default != nil:
			engine.Add(&out.Graph, engine.FieldFact(e.Name, f.Name, *f.Schema.Default), engine.Reason{Type: engine.ReasonAsserted, Pos: f.Pos})
		case !f.Optional:
			*ds = append(*ds, diag.New(e.Pos, "fact %s needs %s, which every %s has", e.Name, f.Name, kindName))
		}
	}
	return links
}
