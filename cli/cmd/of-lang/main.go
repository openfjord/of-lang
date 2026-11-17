// SPDX-License-Identifier: Apache-2.0

// Command of-lang checks packages written in the modeling language, runs their
// rules over facts, explains what it derived, compares two sets of facts, and
// runs a package's tests.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"oflang/internal/check"
	"oflang/internal/diag"
	"oflang/internal/engine"
	"oflang/internal/facts"
	"oflang/internal/load"
	"oflang/internal/model"
	"oflang/internal/suite"
)

const usage = `usage: of-lang <command> [flags] <args>

commands:
  check <package-dir>...                                  check definitions and docs
  run <package-dir> <facts.yaml>                          derive facts and check constraints
  explain <package-dir> <facts.yaml> <subject> <link-or-kind> [object]
                                                          show why a fact holds
  diff <package-dir> <before.yaml> <after.yaml>           show what changes between two sets of facts
  test <package-dir>                                      run the package's tests

flags:
  --packages <dir>      where to look for dependencies (repeatable; default: the package's parent folder)
  --format text|json    output format (default: text)

exit codes: 0 everything passes, 1 problems found, 2 the tool could not run
`

// Exit codes.
const (
	exitOK       = 0
	exitProblems = 1
	exitBroken   = 2
)

// errUsage marks a command line the tool cannot make sense of.
var errUsage = errors.New("usage")

type options struct {
	format   string
	packages []string
	args     []string
}

func main() {
	code, err := run(os.Args[1:], os.Stdout)
	if err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprint(os.Stderr, usage)
		}
		fmt.Fprintln(os.Stderr, "of-lang:", err)
		os.Exit(exitBroken)
	}
	os.Exit(code)
}

func parseArgs(args []string) (options, error) {
	o := options{format: "text"}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--format", "--packages":
			if i+1 >= len(args) {
				return o, fmt.Errorf("%s needs a value: %w", a, errUsage)
			}
			i++
			if a == "--format" {
				o.format = args[i]
			} else {
				o.packages = append(o.packages, args[i])
			}
		case "-h", "--help":
			return o, errUsage
		default:
			o.args = append(o.args, a)
		}
	}
	if o.format != "text" && o.format != "json" {
		return o, fmt.Errorf("--format is text or json, not %q: %w", o.format, errUsage)
	}
	return o, nil
}

func run(args []string, out io.Writer) (int, error) {
	o, err := parseArgs(args)
	if err != nil {
		return exitBroken, err
	}
	if len(o.args) == 0 {
		return exitBroken, fmt.Errorf("no command: %w", errUsage)
	}
	cmd, rest := o.args[0], o.args[1:]
	want := map[string][2]int{"check": {1, -1}, "run": {2, 2}, "explain": {4, 5}, "diff": {3, 3}, "test": {1, 1}}
	n, ok := want[cmd]
	if !ok {
		return exitBroken, fmt.Errorf("unknown command %q: %w", cmd, errUsage)
	}
	if len(rest) < n[0] || n[1] >= 0 && len(rest) > n[1] {
		return exitBroken, fmt.Errorf("wrong number of arguments for %s: %w", cmd, errUsage)
	}
	r := reporter{out: out, json: o.format == "json"}
	opts := load.Options{SearchDirs: o.packages}
	switch cmd {
	case "check":
		return cmdCheck(r, opts, rest)
	case "run":
		return cmdRun(r, opts, rest[0], rest[1])
	case "explain":
		return cmdExplain(r, opts, rest)
	case "diff":
		return cmdDiff(r, opts, rest[0], rest[1], rest[2])
	default:
		return cmdTest(r, opts, rest[0])
	}
}

// prepared is a loaded package ready to run.
type prepared struct {
	loaded load.Loaded
	strata [][]model.QName
}

// prepare loads and checks a package. Diagnostics mean it cannot be run.
func prepare(opts load.Options, dir string) (prepared, []diag.Diagnostic, error) {
	l, ds, err := load.Package(dir, opts)
	if err != nil {
		return prepared{}, nil, err
	}
	strata, sds := engine.Stratify(l.World)
	ds = diag.Sorted(append(ds, sds...))
	return prepared{loaded: l, strata: strata}, ds, nil
}

func cmdCheck(r reporter, opts load.Options, dirs []string) (int, error) {
	var all []diag.Diagnostic
	for _, dir := range dirs {
		l, ds, err := load.Package(dir, opts)
		if err != nil {
			return exitBroken, err
		}
		all = append(all, ds...)
		if len(ds) == 0 {
			all = append(all, check.World(l.World)...)
		}
	}
	all = relative(diag.Sorted(all))
	if err := r.diagnostics(all); err != nil {
		return exitBroken, err
	}
	if len(all) > 0 {
		return exitProblems, nil
	}
	return exitOK, nil
}

// runFacts loads a package and one facts file, then runs it.
func runFacts(opts load.Options, dir, path string) (prepared, facts.File, engine.Result, []diag.Diagnostic, error) {
	p, ds, err := prepare(opts, dir)
	if err != nil || len(ds) > 0 {
		return p, facts.File{}, engine.Result{}, ds, err
	}
	f, ds, err := facts.Read(p.loaded, path)
	if err != nil || len(ds) > 0 {
		return p, f, engine.Result{}, diag.Sorted(ds), err
	}
	return p, f, engine.Run(p.loaded.World, p.strata, f.Graph), nil, nil
}

func cmdRun(r reporter, opts load.Options, dir, path string) (int, error) {
	p, _, res, ds, err := runFacts(opts, dir, path)
	if err != nil {
		return exitBroken, err
	}
	if len(ds) > 0 {
		return exitProblems, r.diagnostics(relative(ds))
	}
	if err := r.result(p.loaded.World, res); err != nil {
		return exitBroken, err
	}
	if len(res.Violations) > 0 || len(res.Problems) > 0 {
		return exitProblems, nil
	}
	return exitOK, nil
}

func cmdExplain(r reporter, opts load.Options, args []string) (int, error) {
	p, f, res, ds, err := runFacts(opts, args[0], args[1])
	if err != nil {
		return exitBroken, err
	}
	if len(ds) > 0 {
		return exitProblems, r.diagnostics(relative(ds))
	}
	subject, name := args[2], args[3]
	ref, ok := f.Scope[name]
	if !ok {
		return exitBroken, fmt.Errorf("%s is not imported by %s", name, args[1])
	}
	var fact engine.Fact
	switch ref.Sort {
	case model.SortLink:
		if len(args) < 5 {
			return exitBroken, fmt.Errorf("%s is a link, so explain needs the object too: %w", name, errUsage)
		}
		from, to := subject, args[4]
		if ref.Inverted {
			from, to = to, from
		}
		fact = engine.LinkFact(from, ref.Q, to)
	case model.SortKind:
		fact = engine.KindFact(subject, ref.Q)
	default:
		return exitBroken, fmt.Errorf("%s is a %s; explain takes a link or a kind", name, ref.Sort)
	}
	if !engine.Has(res.Graph, fact) {
		if err := r.notHeld(p.loaded.World, fact); err != nil {
			return exitBroken, err
		}
		return exitProblems, nil
	}
	if err := r.explanation(p.loaded.World, res.Graph, fact); err != nil {
		return exitBroken, err
	}
	return exitOK, nil
}

func cmdDiff(r reporter, opts load.Options, dir, before, after string) (int, error) {
	p, _, a, ds, err := runFacts(opts, dir, before)
	if err != nil {
		return exitBroken, err
	}
	if len(ds) > 0 {
		return exitProblems, r.diagnostics(relative(ds))
	}
	f, ds, err := facts.Read(p.loaded, after)
	if err != nil {
		return exitBroken, err
	}
	if len(ds) > 0 {
		return exitProblems, r.diagnostics(relative(diag.Sorted(ds)))
	}
	b := engine.Run(p.loaded.World, p.strata, f.Graph)
	if err := r.diff(p.loaded.World, a, b); err != nil {
		return exitBroken, err
	}
	return exitOK, nil
}

func cmdTest(r reporter, opts load.Options, dir string) (int, error) {
	p, ds, err := prepare(opts, dir)
	if err != nil {
		return exitBroken, err
	}
	if len(ds) > 0 {
		return exitProblems, r.diagnostics(relative(ds))
	}
	outcomes, ds, err := suite.Run(p.loaded, p.strata)
	if err != nil {
		return exitBroken, err
	}
	if len(ds) > 0 {
		return exitProblems, r.diagnostics(relative(ds))
	}
	for i := range outcomes {
		outcomes[i].File = relPath(outcomes[i].File)
	}
	if err := r.tests(outcomes); err != nil {
		return exitBroken, err
	}
	if slices.ContainsFunc(outcomes, func(o suite.Outcome) bool { return !suite.Passed(o) }) {
		return exitProblems, nil
	}
	return exitOK, nil
}

// relative shortens diagnostic paths to the working directory, so output reads
// the way the user typed the paths.
func relative(ds []diag.Diagnostic) []diag.Diagnostic {
	out := slices.Clone(ds)
	for i := range out {
		out[i].Pos.File = relPath(out[i].Pos.File)
	}
	return out
}

func relPath(path string) string {
	if path == "" {
		return path
	}
	wd, err := os.Getwd()
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(wd, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return path
	}
	return rel
}
