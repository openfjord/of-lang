// SPDX-License-Identifier: Apache-2.0

// Package load reads a package and its dependencies from disk, resolves every
// import block, and turns each statement into the model's definitions.
package load

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
)

// Options says where to look for packages named in pkg: imports, and which
// version of the root package to read when it has a pending change. Other
// packages always read without their change.
type Options struct {
	SearchDirs []string
	Version    source.Version
}

// Loaded is a world together with the parsed files it came from, which test
// files need to resolve their own imports.
type Loaded struct {
	World   model.World
	Files   map[string]source.File
	Version source.Version
	// HasChange says whether the root package has a pending change.
	HasChange bool
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
	if opts.Version == "" {
		opts.Version = source.WithChange
	}
	st := state{
		opts:     opts,
		packages: map[string]model.Package{},
		files:    map[string]source.File{},
		loading:  map[string]bool{},
	}
	root, err := loadPackage(&st, abs, "", diag.Pos{}, opts.Version)
	if err != nil {
		return Loaded{}, nil, err
	}
	w := model.World{
		Root:       root,
		Packages:   st.packages,
		Kinds:      map[model.QName]model.Kind{},
		Links:      map[model.QName]model.Link{},
		ValueTypes: map[model.QName]model.ValueType{},
		Docs:       map[model.QName]model.Doc{},
		Statements: map[model.QName]model.Statement{},
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
	l := Loaded{World: w, Files: st.files, Version: opts.Version, HasChange: st.rootChange}
	return l, diag.Sorted(st.diags), nil
}

type state struct {
	opts       Options
	packages   map[string]model.Package
	files      map[string]source.File
	loading    map[string]bool
	diags      []diag.Diagnostic
	rootChange bool
}

// ManifestNames are the names a manifest may have.
var ManifestNames = []string{"package.yaml", "package.yml"}

var lockNames = []string{"package.lock.yaml", "package.lock.yml"}

// manifestPath finds the manifest of the package in dir.
func manifestPath(dir string) (string, bool) {
	for _, name := range ManifestNames {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

func loadPackage(st *state, dir, wantName string, from diag.Pos, version source.Version) (string, error) {
	path, ok := manifestPath(dir)
	if !ok {
		return "", fmt.Errorf("%s has no package.yaml or package.yml", dir)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read package manifest: %w", err)
	}
	m, ds := parseManifest(path, data)
	st.diags = append(st.diags, ds...)
	if m.Package == "" {
		return "", fmt.Errorf("%s: package: is missing", path)
	}
	if wantName != "" && m.Package != wantName {
		st.diags = append(st.diags, diag.New(from, "found %s, which is package %s, not %s (PK9)", dir, m.Package, wantName))
	}
	if _, ok := st.packages[m.Package]; ok {
		return m.Package, nil
	}
	st.loading[m.Package] = true
	pkg := model.Package{Name: m.Package, Dir: dir, Manifest: m, Defs: map[string]model.Definition{}}
	change, err := collectFiles(st, &pkg, version)
	if err != nil {
		return "", err
	}
	if wantName == "" {
		st.rootChange = change
	}
	st.packages[m.Package] = pkg
	for _, dep := range sortedKeys(m.DependsOn) {
		if st.loading[dep] {
			st.diags = append(st.diags, diag.New(m.Pos, "packages depend on each other in a circle through %s (PK10)", dep))
			continue
		}
		depDir, ok := findPackage(st.opts.SearchDirs, dep)
		if !ok {
			st.diags = append(st.diags, diag.New(m.Pos, "dependency %s not found in %v (PK9)", dep, st.opts.SearchDirs))
			continue
		}
		if _, err := loadPackage(st, depDir, dep, m.Pos, source.WithoutChange); err != nil {
			return "", err
		}
		depVersion := st.packages[dep].Manifest.Version
		if !versionMatches(m.DependsOn[dep], depVersion) {
			st.diags = append(st.diags, diag.New(m.Pos, "%s needs %s %s, but found version %s (PK9)", m.Package, dep, m.DependsOn[dep], depVersion))
		}
	}
	st.loading[m.Package] = false
	return m.Package, nil
}

// collectFiles checks the layout of a package, reads every file in its defs/,
// and lists its tests and docs. It reports whether the package has a pending
// change, in its definitions or its tests.
func collectFiles(st *state, pkg *model.Package, version source.Version) (bool, error) {
	entries, err := os.ReadDir(pkg.Dir)
	if err != nil {
		return false, fmt.Errorf("read package %s: %w", pkg.Name, err)
	}
	var manifests, locks []string
	hasDefs := false
	for _, e := range entries {
		name := e.Name()
		switch {
		case strings.HasPrefix(name, "."):
			// Hidden files belong to editors and version control, not the package.
		case slices.Contains(ManifestNames, name) && !e.IsDir():
			manifests = append(manifests, name)
		case slices.Contains(lockNames, name) && !e.IsDir():
			locks = append(locks, name)
		case name == "defs" && e.IsDir():
			hasDefs = true
		case (name == "docs" || name == "tests") && e.IsDir():
			// Read below.
		default:
			st.diags = append(st.diags, diag.New(diag.Pos{File: filepath.Join(pkg.Dir, name)}, "the top level of a package holds only the manifest, the lock file, defs/, docs/ and tests/ (PK1)"))
		}
	}
	if len(manifests) > 1 {
		st.diags = append(st.diags, diag.New(diag.Pos{File: filepath.Join(pkg.Dir, manifests[1])}, "a package has package.yaml or package.yml, not both (PK1)"))
	}
	if len(locks) > 1 {
		st.diags = append(st.diags, diag.New(diag.Pos{File: filepath.Join(pkg.Dir, locks[1])}, "a package has package.lock.yaml or package.lock.yml, not both (PK16)"))
	}
	if !hasDefs {
		st.diags = append(st.diags, diag.New(diag.Pos{File: pkg.Dir}, "a package needs a defs/ folder (PK1)"))
	}
	var form source.Form
	var formPath string
	languageFile := func(path string) bool {
		f, ok := source.FormOf(path)
		if !ok {
			st.diags = append(st.diags, diag.New(diag.Pos{File: path}, "files in defs/ and tests/ end in .of.yaml, .of.yml or .of.json (FL15)"))
			return false
		}
		switch {
		case form == "":
			form, formPath = f, path
		case f != form:
			st.diags = append(st.diags, diag.New(diag.Pos{File: path}, "this file is %s, but %s is %s; all the files of a package use one form (FL18)", f, relPath(pkg.Dir, formPath), form))
		}
		return true
	}
	change := false
	err = walkFiles(filepath.Join(pkg.Dir, "defs"), func(path string) error {
		if !languageFile(path) {
			return nil
		}
		raw, ds, err := source.ReadFile(path)
		if err != nil {
			return err
		}
		st.diags = append(st.diags, ds...)
		st.diags = append(st.diags, source.UnusedImports(raw)...)
		change = change || source.HasChange(raw)
		f, ds := source.Select(raw, version)
		st.diags = append(st.diags, ds...)
		st.files[path] = f
		pkg.Files = append(pkg.Files, path)
		for _, e := range f.Entries {
			addDefinition(st, pkg, e, path)
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("read package %s: %w", pkg.Name, err)
	}
	err = walkFiles(filepath.Join(pkg.Dir, "tests"), func(path string) error {
		if !languageFile(path) {
			return nil
		}
		pkg.Tests = append(pkg.Tests, path)
		raw, _, err := source.ReadFile(path)
		if err != nil {
			return err
		}
		change = change || source.HasChange(raw)
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("read tests of %s: %w", pkg.Name, err)
	}
	err = walkFiles(filepath.Join(pkg.Dir, "docs"), func(path string) error {
		if strings.HasSuffix(path, ".md") {
			pkg.Docs = append(pkg.Docs, path)
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("read docs of %s: %w", pkg.Name, err)
	}
	return change, nil
}

// walkFiles calls fn for every file under dir, in name order. A missing dir
// has no files.
func walkFiles(dir string, fn func(path string) error) error {
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		return fn(path)
	})
}

func addDefinition(st *state, pkg *model.Package, e source.Entry, path string) {
	sort := model.Sort(e.Is)
	if !slices.Contains(model.DefinitionSorts, sort) {
		if slices.Contains(model.TestSorts, sort) {
			st.diags = append(st.diags, diag.New(e.Pos, "%s is a %s; those belong in tests/ (FL4)", e.Name, sort))
			return
		}
		st.diags = append(st.diags, diag.New(e.Pos, "%s says is: %s; a statement in defs/ is one of %v (FL1)", e.Name, e.Is, model.DefinitionSorts))
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
		st.diags = append(st.diags, diag.New(d.Pos, "import cannot name a statement (FL5)"))
		return
	}
	if prev, ok := pkg.Defs[name]; ok {
		st.diags = append(st.diags, diag.New(d.Pos, "%s is already defined at %s (FL6)", name, diag.Format(diag.Diagnostic{Pos: prev.Pos})))
		return
	}
	for other, prev := range pkg.Defs {
		if strings.EqualFold(other, name) {
			st.diags = append(st.diags, diag.New(d.Pos, "%s differs from %s, defined at %s, only by case (FL6)", name, other, diag.Format(diag.Diagnostic{Pos: prev.Pos})))
		}
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
			path, ok := manifestPath(candidate)
			if !ok {
				continue
			}
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			m, _ := parseManifest(path, data)
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
			case !ok && file != "":
				ds = append(ds, diag.New(u.Pos, "%s is not defined in %s (FL10)", u.Name, imp.From))
				continue
			case !ok:
				ds = append(ds, diag.New(u.Pos, "%s is not defined in package %s (FL11)", u.Name, target.Name))
				continue
			case file != "" && d.File != file:
				ds = append(ds, diag.New(u.Pos, "%s is not defined in %s; it is in %s (FL10)", u.Name, imp.From, relPath(pkg.Dir, d.File)))
				continue
			case file == "" && !exported(target, u.Name, d):
				ds = append(ds, diag.New(u.Pos, "%s is not exported by package %s (FL11)", u.Name, target.Name))
				continue
			}
			if prev, ok := scope[u.As]; ok && prev.Q != (model.QName{Pkg: target.Name, Name: u.Name}) {
				ds = append(ds, diag.New(u.Pos, "%s is already a name in this file; rename one with as: (FL12)", u.As))
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
			*ds = append(*ds, diag.New(imp.Pos, "%s is outside package %s; import from another package with pkg: (FL10)", imp.From, pkg.Name))
			return model.Package{}, "", false
		}
		if _, err := os.Stat(target); err != nil {
			*ds = append(*ds, diag.New(imp.Pos, "%s does not exist (FL10)", imp.From))
			return model.Package{}, "", false
		}
		return pkg, target, true
	case strings.HasPrefix(imp.From, "pkg:"):
		name := strings.TrimPrefix(imp.From, "pkg:")
		if _, ok := pkg.Manifest.DependsOn[name]; !ok {
			*ds = append(*ds, diag.New(imp.Pos, "package %s is not in dependsOn (FL11)", name))
			return model.Package{}, "", false
		}
		target, ok := pkgs[name]
		if !ok {
			*ds = append(*ds, diag.New(imp.Pos, "package %s was not found (PK9)", name))
			return model.Package{}, "", false
		}
		return target, "", true
	default:
		*ds = append(*ds, diag.New(imp.Pos, "%s must start with a scheme: file: or pkg: (FL17)", imp.From))
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
