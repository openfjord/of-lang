// SPDX-License-Identifier: Apache-2.0

// Package load reads a package and its dependencies from disk, resolves every
// import block, and turns each entry into the model's definitions.
package load

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/openfjord/of-lang/internal/diag"
	"github.com/openfjord/of-lang/internal/model"
	"github.com/openfjord/of-lang/internal/source"
)

// Options says where to look for packages named in pkg: imports.
type Options struct {
	SearchDirs []string
}

// Loaded is a world together with the parsed files it came from, which test
// files need to resolve their own imports.
type Loaded struct {
	World model.World
	Files map[string]source.File
}

// Package loads the package in dir, every package it depends on, and resolves
// all of their definitions. The error is for a package that cannot be read at
// all; everything else is a diagnostic.
func Package(dir string, opts Options) (Loaded, []diag.Diagnostic, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return Loaded{}, nil, fmt.Errorf("resolve %s: %w", dir, err)
	}
	if len(opts.SearchDirs) == 0 {
		opts.SearchDirs = []string{filepath.Dir(abs)}
	}
	st := state{
		opts:     opts,
		packages: map[string]model.Package{},
		files:    map[string]source.File{},
		loading:  map[string]bool{},
	}
	root, err := loadPackage(&st, abs, "", diag.Pos{})
	if err != nil {
		return Loaded{}, nil, err
	}
	w := model.World{
		Root:       root,
		Packages:   st.packages,
		Kinds:      map[model.QName]model.Kind{},
		Links:      map[model.QName]model.Link{},
		ValueTypes: map[model.QName]model.ValueType{},
		Statements: map[model.QName]model.Statement{},
		Mappings:   map[model.QName]model.Mapping{},
		Inverses:   map[model.QName]model.QName{},
	}
	for _, name := range sortedKeys(st.packages) {
		pkg := st.packages[name]
		for defName, d := range pkg.Defs {
			if d.InverseOf != "" {
				w.Inverses[model.QName{Pkg: name, Name: defName}] = model.QName{Pkg: name, Name: d.InverseOf}
			}
		}
		for _, path := range pkg.Files {
			f := st.files[path]
			scope, ds := Scope(st.packages, pkg, f)
			st.diags = append(st.diags, ds...)
			for _, e := range f.Entries {
				st.diags = append(st.diags, decodeDefinition(&w, pkg, f.Path, e, scope)...)
			}
		}
	}
	return Loaded{World: w, Files: st.files}, diag.Sorted(st.diags), nil
}

type state struct {
	opts     Options
	packages map[string]model.Package
	files    map[string]source.File
	loading  map[string]bool
	diags    []diag.Diagnostic
}

func loadPackage(st *state, dir, wantName string, from diag.Pos) (string, error) {
	manifestPath := filepath.Join(dir, "package.yaml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", fmt.Errorf("read package manifest: %w", err)
	}
	m, ds := parseManifest(manifestPath, data)
	st.diags = append(st.diags, ds...)
	if m.Package == "" {
		return "", fmt.Errorf("%s: package: is missing", manifestPath)
	}
	if wantName != "" && m.Package != wantName {
		st.diags = append(st.diags, diag.New(from, "found %s, which is package %s, not %s", dir, m.Package, wantName))
	}
	if _, ok := st.packages[m.Package]; ok {
		return m.Package, nil
	}
	st.loading[m.Package] = true
	pkg := model.Package{Name: m.Package, Dir: dir, Manifest: m, Defs: map[string]model.Definition{}}
	if err := collectFiles(st, &pkg); err != nil {
		return "", err
	}
	st.packages[m.Package] = pkg
	for _, dep := range sortedKeys(m.DependsOn) {
		if st.loading[dep] {
			st.diags = append(st.diags, diag.New(m.Pos, "packages depend on each other in a circle through %s", dep))
			continue
		}
		depDir, ok := findPackage(st.opts.SearchDirs, dep)
		if !ok {
			st.diags = append(st.diags, diag.New(m.Pos, "dependency %s not found in %v", dep, st.opts.SearchDirs))
			continue
		}
		if _, err := loadPackage(st, depDir, dep, m.Pos); err != nil {
			return "", err
		}
		depVersion := st.packages[dep].Manifest.Version
		if !versionMatches(m.DependsOn[dep], depVersion) {
			st.diags = append(st.diags, diag.New(m.Pos, "%s needs %s %s, but found version %s", m.Package, dep, m.DependsOn[dep], depVersion))
		}
	}
	st.loading[m.Package] = false
	return m.Package, nil
}

// collectFiles reads every definition file of a package, skipping its tests,
// and records where each name is defined.
func collectFiles(st *state, pkg *model.Package) error {
	testsDir := filepath.Join(pkg.Dir, "tests")
	err := filepath.WalkDir(pkg.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == testsDir {
				return filepath.SkipDir
			}
			return nil
		}
		switch {
		case strings.HasSuffix(path, ".md"):
			pkg.Docs = append(pkg.Docs, path)
			return nil
		case !strings.HasSuffix(path, ".yaml") || filepath.Base(path) == "package.yaml" && filepath.Dir(path) == pkg.Dir:
			return nil
		}
		f, ds, err := source.ReadFile(path)
		if err != nil {
			return err
		}
		st.diags = append(st.diags, ds...)
		st.files[path] = f
		pkg.Files = append(pkg.Files, path)
		for _, e := range f.Entries {
			addDefinition(st, pkg, e, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("read package %s: %w", pkg.Name, err)
	}
	return nil
}

func addDefinition(st *state, pkg *model.Package, e source.Entry, path string) {
	sort := model.Sort(e.Is)
	if !slices.Contains(model.DefinitionSorts, sort) {
		if slices.Contains(model.TestSorts, sort) {
			st.diags = append(st.diags, diag.New(e.Pos, "%s is a %s; those belong in the package's tests/ folder", e.Name, sort))
			return
		}
		st.diags = append(st.diags, diag.New(e.Pos, "%s says is: %s; a definition is one of %v", e.Name, e.Is, model.DefinitionSorts))
		return
	}
	define(st, pkg, e.Name, model.Definition{Sort: sort, Pos: e.Pos, File: path})
	if sort != model.SortLink {
		return
	}
	m, ok := source.AsMap(path, e.Node, &st.diags)
	if !ok {
		return
	}
	var ignored []diag.Diagnostic
	if inv, ok := source.GetScalar(m, "inverse", &ignored); ok {
		define(st, pkg, inv, model.Definition{Sort: model.SortLink, Pos: e.Pos, File: path, InverseOf: e.Name})
	}
}

func define(st *state, pkg *model.Package, name string, d model.Definition) {
	if name == "import" {
		st.diags = append(st.diags, diag.New(d.Pos, "import is reserved and cannot name an entry"))
		return
	}
	if prev, ok := pkg.Defs[name]; ok {
		st.diags = append(st.diags, diag.New(d.Pos, "%s is already defined at %s", name, diag.Format(diag.Diagnostic{Pos: prev.Pos})))
		return
	}
	pkg.Defs[name] = d
}

func findPackage(dirs []string, name string) (string, bool) {
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			candidate := filepath.Join(dir, e.Name())
			data, err := os.ReadFile(filepath.Join(candidate, "package.yaml"))
			if err != nil {
				continue
			}
			m, _ := parseManifest(filepath.Join(candidate, "package.yaml"), data)
			if m.Package == name {
				return candidate, true
			}
		}
	}
	return "", false
}

// Scope resolves a file's import block into the names the file may use: its
// own definitions, and what it imports.
func Scope(pkgs map[string]model.Package, pkg model.Package, f source.File) (map[string]model.Ref, []diag.Diagnostic) {
	scope := map[string]model.Ref{}
	var ds []diag.Diagnostic
	for name, d := range pkg.Defs {
		if d.File == f.Path {
			scope[name] = refOf(pkg.Name, name, d)
		}
	}
	for _, imp := range f.Imports {
		target, file, ok := importSource(pkgs, pkg, f.Path, imp, &ds)
		if !ok {
			continue
		}
		for _, u := range imp.Uses {
			d, ok := target.Defs[u.Name]
			switch {
			case !ok:
				ds = append(ds, diag.New(u.Pos, "%s is not defined in %s", u.Name, imp.From))
				continue
			case file != "" && d.File != file:
				ds = append(ds, diag.New(u.Pos, "%s is not defined in %s; it is in %s", u.Name, imp.From, relPath(pkg.Dir, d.File)))
				continue
			case file == "" && !exported(target, u.Name, d):
				ds = append(ds, diag.New(u.Pos, "%s is not exported by package %s", u.Name, target.Name))
				continue
			}
			if prev, ok := scope[u.As]; ok && prev.Q != (model.QName{Pkg: target.Name, Name: u.Name}) {
				ds = append(ds, diag.New(u.Pos, "%s is already a name in this file; rename one with as:", u.As))
				continue
			}
			scope[u.As] = refOf(target.Name, u.Name, d)
		}
	}
	return scope, ds
}

func importSource(pkgs map[string]model.Package, pkg model.Package, path string, imp source.Import, ds *[]diag.Diagnostic) (model.Package, string, bool) {
	switch {
	case strings.HasPrefix(imp.From, "file:"):
		rel := strings.TrimPrefix(imp.From, "file:")
		target := filepath.Clean(filepath.Join(filepath.Dir(path), rel))
		if !strings.HasPrefix(target, pkg.Dir+string(filepath.Separator)) {
			*ds = append(*ds, diag.New(imp.Pos, "%s is outside package %s; import from another package with pkg:", imp.From, pkg.Name))
			return model.Package{}, "", false
		}
		if _, err := os.Stat(target); err != nil {
			*ds = append(*ds, diag.New(imp.Pos, "%s does not exist", imp.From))
			return model.Package{}, "", false
		}
		return pkg, target, true
	case strings.HasPrefix(imp.From, "pkg:"):
		name := strings.TrimPrefix(imp.From, "pkg:")
		if _, ok := pkg.Manifest.DependsOn[name]; !ok && name != pkg.Name {
			*ds = append(*ds, diag.New(imp.Pos, "package %s is not in dependsOn", name))
			return model.Package{}, "", false
		}
		target, ok := pkgs[name]
		if !ok {
			*ds = append(*ds, diag.New(imp.Pos, "package %s was not found", name))
			return model.Package{}, "", false
		}
		return target, "", true
	default:
		*ds = append(*ds, diag.New(imp.Pos, "%s must start with file: or pkg:", imp.From))
		return model.Package{}, "", false
	}
}

func exported(pkg model.Package, name string, d model.Definition) bool {
	ex := pkg.Manifest.Exports
	if d.InverseOf != "" {
		return slices.Contains(ex.Links, d.InverseOf)
	}
	return slices.Contains(ex.Kinds, name) || slices.Contains(ex.Links, name) || slices.Contains(ex.ValueTypes, name)
}

func refOf(pkg, name string, d model.Definition) model.Ref {
	if d.InverseOf != "" {
		return model.Ref{Sort: model.SortLink, Q: model.QName{Pkg: pkg, Name: d.InverseOf}, Inverted: true}
	}
	return model.Ref{Sort: d.Sort, Q: model.QName{Pkg: pkg, Name: name}}
}

func relPath(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
