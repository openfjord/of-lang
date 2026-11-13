// Package source reads the language's YAML files into entries and import
// blocks, keeping positions so later stages can point at what they reject.
package source

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"go.yaml.in/yaml/v3"

	"oflang/internal/diag"
)

// Use is one name an import block brings in, possibly renamed.
type Use struct {
	Name string
	As   string
	Pos  diag.Pos
}

// Import is one entry of a file's import block.
type Import struct {
	From string
	Uses []Use
	Pos  diag.Pos
}

// Entry is one top-level entry of a file: a definition, a fact, a test or a
// question, keyed by its name.
type Entry struct {
	Name string
	Is   string
	Pos  diag.Pos
	Node *yaml.Node
}

// File is a parsed YAML file of the language.
type File struct {
	Path    string
	Imports []Import
	Entries []Entry
}

// Pair is one key and value of a YAML mapping.
type Pair struct {
	Key   string
	Pos   diag.Pos
	Value *yaml.Node
}

// Map is a YAML mapping with its pairs in file order.
type Map struct {
	File  string
	Pos   diag.Pos
	Pairs []Pair
}

// ReadFile reads and parses one file. The error is for files that cannot be
// read at all; problems in the content are diagnostics.
func ReadFile(path string) (File, []diag.Diagnostic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, nil, fmt.Errorf("read %s: %w", path, err)
	}
	f, ds := Parse(path, data)
	return f, ds, nil
}

// Parse splits a file into its import block and its entries.
func Parse(path string, data []byte) (File, []diag.Diagnostic) {
	root, ds := Document(path, data)
	if root == nil {
		return File{Path: path}, ds
	}
	m, ok := AsMap(path, root, &ds)
	if !ok {
		return File{Path: path}, ds
	}
	f := File{Path: path}
	for _, p := range m.Pairs {
		if p.Key == "import" {
			f.Imports = parseImports(path, p.Value, &ds)
			continue
		}
		em, ok := AsMap(path, p.Value, &ds)
		if !ok {
			continue
		}
		is, ok := Get(em, "is")
		if !ok {
			ds = append(ds, diag.New(p.Pos, "%s does not say what it is: add an is: line", p.Key))
			continue
		}
		isText, ok := Scalar(path, is, &ds)
		if !ok {
			continue
		}
		f.Entries = append(f.Entries, Entry{Name: p.Key, Is: isText, Pos: p.Pos, Node: p.Value})
	}
	return f, ds
}

// Document parses YAML text into its single top-level node.
func Document(path string, data []byte) (*yaml.Node, []diag.Diagnostic) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var doc yaml.Node
	if err := dec.Decode(&doc); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, []diag.Diagnostic{diag.New(diag.Pos{File: path}, "not valid YAML: %v", err)}
	}
	if len(doc.Content) == 0 {
		return nil, nil
	}
	return doc.Content[0], nil
}

func parseImports(path string, n *yaml.Node, ds *[]diag.Diagnostic) []Import {
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(PosOf(path, n), "import is a list of entries, each with from: and use:"))
		return nil
	}
	var out []Import
	for _, item := range n.Content {
		m, ok := AsMap(path, item, ds)
		if !ok {
			continue
		}
		Allow(m, []string{"from", "use"}, ds)
		fromNode, ok := Get(m, "from")
		if !ok {
			*ds = append(*ds, diag.New(m.Pos, "an import needs from:"))
			continue
		}
		from, ok := Scalar(path, fromNode, ds)
		if !ok {
			continue
		}
		useNode, ok := Get(m, "use")
		if !ok {
			*ds = append(*ds, diag.New(m.Pos, "an import needs use: with the names it brings in"))
			continue
		}
		out = append(out, Import{From: from, Uses: parseUses(path, useNode, ds), Pos: m.Pos})
	}
	return out
}

func parseUses(path string, n *yaml.Node, ds *[]diag.Diagnostic) []Use {
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(PosOf(path, n), "use: is a list of names"))
		return nil
	}
	var out []Use
	for _, item := range n.Content {
		if item.Kind == yaml.ScalarNode {
			out = append(out, Use{Name: item.Value, As: item.Value, Pos: PosOf(path, item)})
			continue
		}
		m, ok := AsMap(path, item, ds)
		if !ok {
			continue
		}
		Allow(m, []string{"name", "as"}, ds)
		name, okName := GetScalar(m, "name", ds)
		as, okAs := GetScalar(m, "as", ds)
		if !okName || !okAs {
			*ds = append(*ds, diag.New(m.Pos, "a renamed import needs name: and as:"))
			continue
		}
		out = append(out, Use{Name: name, As: as, Pos: m.Pos})
	}
	return out
}

// PosOf returns where a node starts.
func PosOf(path string, n *yaml.Node) diag.Pos {
	return diag.Pos{File: path, Line: n.Line, Column: n.Column}
}

// AsMap reads a mapping node, reporting a duplicate key or a node of another
// shape.
func AsMap(path string, n *yaml.Node, ds *[]diag.Diagnostic) (Map, bool) {
	if n.Kind != yaml.MappingNode {
		*ds = append(*ds, diag.New(PosOf(path, n), "expected a mapping of keys and values"))
		return Map{}, false
	}
	m := Map{File: path, Pos: PosOf(path, n)}
	seen := map[string]bool{}
	for i := 0; i+1 < len(n.Content); i += 2 {
		k := n.Content[i]
		if seen[k.Value] {
			*ds = append(*ds, diag.New(PosOf(path, k), "%s appears twice", k.Value))
			continue
		}
		seen[k.Value] = true
		m.Pairs = append(m.Pairs, Pair{Key: k.Value, Pos: PosOf(path, k), Value: n.Content[i+1]})
	}
	return m, true
}

// Get returns the value under a key.
func Get(m Map, key string) (*yaml.Node, bool) {
	for _, p := range m.Pairs {
		if p.Key == key {
			return p.Value, true
		}
	}
	return nil, false
}

// GetScalar returns the text under a key, when it is a single value.
func GetScalar(m Map, key string, ds *[]diag.Diagnostic) (string, bool) {
	n, ok := Get(m, key)
	if !ok {
		return "", false
	}
	return Scalar(m.File, n, ds)
}

// Allow reports every key that is not in the allowed list.
func Allow(m Map, allowed []string, ds *[]diag.Diagnostic) {
	for _, p := range m.Pairs {
		if !slices.Contains(allowed, p.Key) {
			*ds = append(*ds, diag.New(p.Pos, "%s is not allowed here; allowed: %v", p.Key, allowed))
		}
	}
}

// Scalar reads a single value as text.
func Scalar(path string, n *yaml.Node, ds *[]diag.Diagnostic) (string, bool) {
	if n.Kind != yaml.ScalarNode {
		*ds = append(*ds, diag.New(PosOf(path, n), "expected a single value"))
		return "", false
	}
	return n.Value, true
}

// Strings reads one value or a list of values as text.
func Strings(path string, n *yaml.Node, ds *[]diag.Diagnostic) ([]string, bool) {
	if n.Kind == yaml.ScalarNode {
		return []string{n.Value}, true
	}
	if n.Kind != yaml.SequenceNode {
		*ds = append(*ds, diag.New(PosOf(path, n), "expected a value or a list of values"))
		return nil, false
	}
	var out []string
	for _, item := range n.Content {
		s, ok := Scalar(path, item, ds)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// Bool reads true or false.
func Bool(path string, n *yaml.Node, ds *[]diag.Diagnostic) (bool, bool) {
	if n.Kind != yaml.ScalarNode || n.ShortTag() != "!!bool" {
		*ds = append(*ds, diag.New(PosOf(path, n), "expected true or false"))
		return false, false
	}
	return n.Value == "true", true
}
