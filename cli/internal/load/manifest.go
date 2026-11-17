// SPDX-License-Identifier: Apache-2.0

package load

import (
	"regexp"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"

	"oflang/internal/diag"
	"oflang/internal/model"
	"oflang/internal/source"
)

var (
	packageNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	versionPattern     = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
)

func parseManifest(path string, data []byte) (model.Manifest, []diag.Diagnostic) {
	root, ds := source.Document(path, data)
	if root == nil {
		return model.Manifest{}, append(ds, diag.New(diag.Pos{File: path}, "package.yaml is empty"))
	}
	mm, ok := source.AsMap(path, root, &ds)
	if !ok {
		return model.Manifest{}, ds
	}
	source.Allow(mm, []string{"package", "version", "namespace", "dependsOn", "exports", "openFor"}, &ds)
	m := model.Manifest{Pos: mm.Pos, DependsOn: map[string]string{}}
	m.Package, _ = source.GetScalar(mm, "package", &ds)
	m.Version, _ = source.GetScalar(mm, "version", &ds)
	m.Namespace, _ = source.GetScalar(mm, "namespace", &ds)
	switch {
	case m.Package == "":
		ds = append(ds, diag.New(mm.Pos, "package.yaml needs package:"))
	case !packageNamePattern.MatchString(m.Package):
		ds = append(ds, diag.New(mm.Pos, "package name %s must be lowercase letters, digits and hyphens", m.Package))
	}
	if !versionPattern.MatchString(m.Version) {
		ds = append(ds, diag.New(mm.Pos, "version must look like 1.2.3, not %q", m.Version))
	}
	if m.Namespace == "" {
		ds = append(ds, diag.New(mm.Pos, "package.yaml needs namespace:, used when exporting to RDF"))
	}
	if n, ok := source.Get(mm, "dependsOn"); ok {
		if dm, ok := source.AsMap(path, n, &ds); ok {
			for _, p := range dm.Pairs {
				r, ok := source.Scalar(path, p.Value, &ds)
				if !ok {
					continue
				}
				if !validRange(r) {
					ds = append(ds, diag.New(p.Pos, "version range %q must be 1.2.3, ^1.2, ^1.2.3, ~1.2 or ~1.2.3", r))
				}
				m.DependsOn[p.Key] = r
			}
		}
	}
	if n, ok := source.Get(mm, "exports"); ok {
		if em, ok := source.AsMap(path, n, &ds); ok {
			source.Allow(em, []string{"kinds", "links", "valueTypes"}, &ds)
			m.Exports.Kinds = optionalStrings(path, em, "kinds", &ds)
			m.Exports.Links = optionalStrings(path, em, "links", &ds)
			m.Exports.ValueTypes = optionalStrings(path, em, "valueTypes", &ds)
		}
	}
	if n, ok := source.Get(mm, "openFor"); ok {
		m.OpenFor = parseOpenFor(path, n, &ds)
	}
	return m, ds
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
		source.Allow(om, []string{"extend", "under"}, ds)
		var p model.OpenPoint
		p.Extend, _ = source.GetScalar(om, "extend", ds)
		p.Under, _ = source.GetScalar(om, "under", ds)
		if (p.Extend == "") == (p.Under == "") {
			*ds = append(*ds, diag.New(om.Pos, "an extension point is either extend: a kind or under: a link"))
			continue
		}
		out = append(out, p)
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
