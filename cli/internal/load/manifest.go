// SPDX-License-Identifier: Apache-2.0

package load

import (
	"regexp"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
	"github.com/openfjord/of-lang/cli/internal/model"
	"github.com/openfjord/of-lang/cli/internal/source"
)

var (
	packageNamePattern = regexp.MustCompile(`^(@[a-z][a-z0-9-]*/)?[a-z][a-z0-9-]*$`)
	versionPattern     = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	tagNamePattern     = regexp.MustCompile(`^[a-z][A-Za-z0-9]*$`)
)

func parseManifest(path string, data []byte) (model.Manifest, []diag.Diagnostic) {
	root, ds := source.Document(path, data)
	if root == nil {
		return model.Manifest{}, append(ds, diag.New(diag.Pos{File: path}, "the manifest is empty (PK4)"))
	}
	mm, ok := source.AsMap(path, root, &ds)
	if !ok {
		return model.Manifest{}, ds
	}
	keys := []string{"package", "version", "owner", "rdfNamespace", "dependsOn", "exports", "openFor", "tags"}
	allow(mm, keys, "PK4", &ds)
	m := model.Manifest{Pos: mm.Pos, DependsOn: map[string]string{}}
	m.Package = manifestText(mm, "package", &ds)
	m.Version = manifestText(mm, "version", &ds)
	m.Owner = manifestText(mm, "owner", &ds)
	m.RDFNamespace = manifestText(mm, "rdfNamespace", &ds)
	if m.Package != "" && !packageNamePattern.MatchString(m.Package) {
		ds = append(ds, diag.New(mm.Pos, "package name %s must be lowercase letters, digits and hyphens, with a letter first, and an optional @namespace/ before it (PK6)", m.Package))
	}
	if m.Version != "" && !versionPattern.MatchString(m.Version) {
		ds = append(ds, diag.New(mm.Pos, "version must be MAJOR.MINOR.PATCH, such as 1.2.3, not %q (PK5)", m.Version))
	}
	if n, ok := source.Get(mm, "dependsOn"); ok {
		if dm, ok := source.AsMap(path, n, &ds); ok {
			for _, p := range dm.Pairs {
				r, ok := source.Scalar(path, p.Value, &ds)
				if !ok {
					continue
				}
				if !packageNamePattern.MatchString(p.Key) {
					ds = append(ds, diag.New(p.Pos, "%s is not a full package name (PK7)", p.Key))
				}
				if !validRange(r) {
					ds = append(ds, diag.New(p.Pos, "version range %q must be 1.2.3, ^1.2, ^1.2.3, ~1.2 or ~1.2.3 (PK9)", r))
				}
				m.DependsOn[p.Key] = r
			}
		}
	}
	if n, ok := source.Get(mm, "exports"); ok {
		if em, ok := source.AsMap(path, n, &ds); ok {
			allow(em, []string{"kinds", "links", "valueTypes"}, "PK4", &ds)
			m.Exports.Kinds = optionalStrings(path, em, "kinds", &ds)
			m.Exports.Links = optionalStrings(path, em, "links", &ds)
			m.Exports.ValueTypes = optionalStrings(path, em, "valueTypes", &ds)
		}
	}
	if n, ok := source.Get(mm, "openFor"); ok {
		m.OpenFor = parseOpenFor(path, n, &ds)
	}
	if n, ok := source.Get(mm, "tags"); ok {
		m.Tags = parseTags(path, n, &ds)
	}
	return m, ds
}

func manifestText(m source.Map, key string, ds *[]diag.Diagnostic) string {
	s, ok := source.GetScalar(m, key, ds)
	if !ok || s == "" {
		*ds = append(*ds, diag.New(m.Pos, "the manifest needs %s: (PK4)", key))
		return ""
	}
	return s
}

func parseOpenFor(path string, n *yaml.Node, ds *[]diag.Diagnostic) []model.OpenPoint {
	items, ok := sequence(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.OpenPoint
	for _, item := range items {
		om, ok := source.AsMap(path, item, ds)
		if !ok {
			continue
		}
		allow(om, []string{"extend", "under"}, "PK13", ds)
		p := model.OpenPoint{Pos: om.Pos}
		p.Extend, _ = source.GetScalar(om, "extend", ds)
		p.Under, _ = source.GetScalar(om, "under", ds)
		if (p.Extend == "") == (p.Under == "") {
			*ds = append(*ds, diag.New(om.Pos, "an extension point is either extend: a kind or under: a link (PK13)"))
			continue
		}
		out = append(out, p)
	}
	return out
}

func parseTags(path string, n *yaml.Node, ds *[]diag.Diagnostic) []model.TagDecl {
	items, ok := sequence(path, n, ds)
	if !ok {
		return nil
	}
	var out []model.TagDecl
	seen := map[string]bool{}
	for _, item := range items {
		tm, ok := source.AsMap(path, item, ds)
		if !ok {
			continue
		}
		allow(tm, []string{"name", "description"}, "PK15", ds)
		t := model.TagDecl{Pos: tm.Pos}
		t.Name, _ = source.GetScalar(tm, "name", ds)
		t.Description, _ = source.GetScalar(tm, "description", ds)
		switch {
		case t.Name == "" || t.Description == "":
			*ds = append(*ds, diag.New(tm.Pos, "a tag needs a name and a description (PK15)"))
			continue
		case seen[t.Name]:
			*ds = append(*ds, diag.New(tm.Pos, "the tag %s is declared twice (PK15)", t.Name))
			continue
		case !tagNamePattern.MatchString(t.Name):
			*ds = append(*ds, diag.New(tm.Pos, "the tag %s must be lowerCamelCase (TG2)", t.Name))
		}
		seen[t.Name] = true
		out = append(out, t)
	}
	return out
}

func optionalStrings(path string, m source.Map, key string, ds *[]diag.Diagnostic) []string {
	n, ok := source.Get(m, key)
	if !ok {
		return nil
	}
	s, _ := source.Strings(path, n, ds)
	return s
}

func validRange(r string) bool {
	_, _, ok := parseRange(r)
	return ok
}

// parseRange turns a version range into an inclusive lower and an exclusive
// upper version.
func parseRange(r string) ([3]int, [3]int, bool) {
	switch {
	case strings.HasPrefix(r, "^"):
		v, parts, ok := parseVersion(strings.TrimPrefix(r, "^"))
		if !ok {
			return v, v, false
		}
		hi := [3]int{v[0] + 1, 0, 0}
		if v[0] == 0 && parts > 1 {
			hi = [3]int{0, v[1] + 1, 0}
		}
		return v, hi, true
	case strings.HasPrefix(r, "~"):
		v, _, ok := parseVersion(strings.TrimPrefix(r, "~"))
		return v, [3]int{v[0], v[1] + 1, 0}, ok
	default:
		v, parts, ok := parseVersion(r)
		if parts != 3 {
			return v, v, false
		}
		return v, [3]int{v[0], v[1], v[2] + 1}, ok
	}
}

func parseVersion(s string) ([3]int, int, bool) {
	var v [3]int
	parts := strings.Split(s, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return v, 0, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return v, 0, false
		}
		v[i] = n
	}
	return v, len(parts), true
}

func versionMatches(r, version string) bool {
	lo, hi, ok := parseRange(r)
	if !ok {
		return false
	}
	v, parts, ok := parseVersion(version)
	if !ok || parts != 3 {
		return false
	}
	return compareVersions(v, lo) >= 0 && compareVersions(v, hi) < 0
}

func compareVersions(a, b [3]int) int {
	for i := range a {
		if a[i] != b[i] {
			if a[i] < b[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}
