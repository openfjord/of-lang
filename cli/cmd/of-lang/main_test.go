// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const example = "../../../examples/restaurant"

func runTool(t *testing.T, args ...string) (int, string) {
	t.Helper()
	var out bytes.Buffer
	code, err := run(args, &out)
	if err != nil {
		t.Fatalf("of-lang %s: %v", strings.Join(args, " "), err)
	}
	return code, out.String()
}

func TestExamplePasses(t *testing.T) {
	t.Parallel()
	if code, out := runTool(t, "check", example); code != exitOK || out != "ok\n" {
		t.Errorf("check: exit %d, output:\n%s", code, out)
	}
	code, out := runTool(t, "test", example)
	if code != exitOK || !strings.HasSuffix(out, "5 passed, 0 failed\n") {
		t.Errorf("test: exit %d, output:\n%s", code, out)
	}
}

// copyExample copies the example package to a temporary folder and applies
// one replacement to one of its files.
func copyExample(t *testing.T, file, old, replacement string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "restaurant")
	if err := os.CopyFS(dir, os.DirFS(example)); err != nil {
		t.Fatalf("copy example: %v", err)
	}
	path := filepath.Join(dir, file)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
	if !bytes.Contains(data, []byte(old)) {
		t.Fatalf("%s does not contain %q", file, old)
	}
	if err := os.WriteFile(path, bytes.Replace(data, []byte(old), []byte(replacement), 1), 0o600); err != nil {
		t.Fatalf("write %s: %v", file, err)
	}
	return dir
}

func TestBrokenPackages(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, file, old, replacement, command, want string
	}{
		{
			name: "reads without to", file: "defs/recipe.of.yaml",
			old: `reads: "{from} is made at {to}"`, replacement: `reads: "{from} is made"`,
			command: "check", want: "recipe.of.yaml:49:3: the reads of madeAtStation must use both {from} and {to} (LK2)",
		},
		{
			name: "missing garnish", file: "tests/examples/allergen-through-a-sauce.of.yaml",
			old: "    usesGarnish: ingredient/parmesan\n", replacement: "",
			command: "test", want: "expected derived, but not: dish/margherita contains allergen/milk",
		},
		{
			name: "wrong answer", file: "tests/questions/which-dishes-contain-nuts.of.yaml",
			old: "$dish: [dish/margherita]", replacement: "$dish: [dish/marinara]",
			command: "test", want: "$dish: expected [dish/marinara], got [dish/margherita]",
		},
		{
			name: "term that does not hold", file: "tests/terms/garnish.of.yaml",
			old: "kind: Dish\n        is a: Recipe", replacement: "kind: Recipe\n        is a: Dish",
			command: "test", want: "expected derived, but not: Recipe is a Dish",
		},
		{
			name: "fact of an unknown kind", file: "tests/examples/base-recipe.of.yaml",
			old: "kind: Station\n", replacement: "kind: Kitchen\n",
			command: "test", want: "base-recipe.of.yaml:20:1: Kitchen is not defined in this file or imported (FL9)\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			dir := copyExample(t, c.file, c.old, c.replacement)
			code, out := runTool(t, c.command, dir)
			if code != exitProblems || !strings.Contains(out, c.want) {
				t.Errorf("%s: exit %d, want %q in output:\n%s", c.command, code, c.want, out)
			}
		})
	}
}
