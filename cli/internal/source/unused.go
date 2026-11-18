// SPDX-License-Identifier: Apache-2.0

package source

import (
	"go.yaml.in/yaml/v3"

	"github.com/openfjord/of-lang/cli/internal/diag"
)

// UnusedImports reports each name that a file imports and does not use
// (FL20). It reads the file as written, before a version of a pending change
// is selected, so a name that only one of the two versions uses counts as
// used. A name counts as used when a statement of the file holds it as a key
// or as a value; that is where every use of a kind, a link, a value type, a
// constraint, a derivative or a doc is written.
func UnusedImports(f File) []diag.Diagnostic {
	written := map[string]struct{}{}
	for _, e := range f.Entries {
		collectScalars(e.Node, written)
	}
	var ds []diag.Diagnostic
	for _, imp := range f.Imports {
		for _, u := range imp.Uses {
			if _, ok := written[u.As]; !ok {
				ds = append(ds, diag.New(u.Pos, "%s is imported and not used (FL20)", u.As))
			}
		}
	}
	return ds
}

// collectScalars adds every scalar under a node to the set.
func collectScalars(n *yaml.Node, into map[string]struct{}) {
	if n == nil {
		return
	}
	if n.Kind == yaml.ScalarNode {
		into[n.Value] = struct{}{}
		return
	}
	for _, child := range n.Content {
		collectScalars(child, into)
	}
}
