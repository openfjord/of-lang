---
title: of-lang — The OpenFjord modeling language
---

# `of-lang`

A command-line tool for [the OpenFjord modeling language](../language/README.md):
it checks packages, runs their derivatives over facts, explains what it derived,
compares two sets of facts, and runs a package's tests.

## Install

With [Homebrew](https://brew.sh), from the
[OpenFjord tap](https://github.com/openfjord/homebrew-tap):

```bash
brew install openfjord/tap/of-lang
```

With Go:

```bash
go install github.com/openfjord/of-lang/cli/cmd/of-lang@latest
```

Or download a binary for macOS or Linux from the
[releases](https://github.com/openfjord/of-lang/releases). `of-lang version`
prints the installed version, and the [changelog](../CHANGELOG.md) lists what
each version changed.

## Build from source

With Go installed, from the repository root:

```bash
cd cli && go build -o of-lang ./cmd/of-lang
```

## Commands

```
of-lang check <package-dir>...
of-lang run <package-dir> <facts-file>
of-lang explain <package-dir> <facts-file> <subject> <link-or-kind> [object]
of-lang diff <package-dir> <before-file> <after-file>
of-lang test <package-dir>
of-lang version
```

| Command | Does |
|---|---|
| `check` | checks definitions and docs of one or more packages |
| `run` | derives what follows from a facts file, then checks constraints |
| `explain` | shows why one fact holds: a link (`subject link object`) or a kind (`subject Kind`) |
| `diff` | shows what changes between two facts files: stated, derived and rejected |
| `test` | runs every test and question in the package's `tests/`, including tests of derived terms |
| `version` | prints the tool's version (also `--version`) |

| Flag | Means |
|---|---|
| `--packages <dir>` | where to look for dependencies; repeatable; by default the folder that holds the package |
| `--format text\|json` | the output format; `text` by default |

| Exit code | Means |
|---|---|
| 0 | everything passes |
| 1 | problems found: a check failed, a constraint rejected something, or a test failed |
| 2 | the tool could not run: a bad command line, or a file it cannot read |

## Output

All examples below run in the folder that holds the package.

### `check`

```
$ of-lang check platform
ok
```

Each problem is one line, with the file, line and column:

```
$ of-lang check platform
platform/defs/service.of.yaml:53:3: the reads of ownedBy must use both {from} and {to} (LK2)
```

Each problem names the rule it breaks, such as `(LK2)`. A package with a
[pending change](../language/reference/pending.md) is checked twice, with the
change and without it; a problem that only the version without the change has
starts with `without the pending change:`. `run`, `explain`, `diff` and `test`
use the version with the change.

### `run`

Every derived fact as its `reads` sentence, with what derived it: a derivative, or
`under` or `chains` on a link. Then every rejection, with its message.

```
$ of-lang run platform platform/tests/examples/personal-data-through-a-call.of.yaml
derived:
  service/checkout reads table/payments  (under)
  service/accounts handles data/personal  (derivative dataFromTable)
  service/checkout handles data/payment  (derivative dataFromTable)
  service/checkout handles data/personal  (derivative dataFromService)
  team/web depends on team/core  (derivative teamDependsOnTeam)
rejected:
  nothing
```

```
$ of-lang run platform platform/tests/examples/service-without-declared-data.of.yaml
derived:
  service/accounts handles data/personal  (derivative dataFromTable)
  service/checkout handles data/personal  (derivative dataFromService)
  team/web depends on team/core  (derivative teamDependsOnTeam)
rejected:
  publicServiceDeclaresItsData: The privacy notice must tell that service/checkout handles data/personal.
```

A rejection makes `run` exit with 1.

### `explain`

The fact, the derivative's `says`, what each variable stood for, and the facts
the match used. A derived fact among them unfolds into its own explanation,
indented below it.

```
$ of-lang explain platform platform/tests/examples/personal-data-through-a-call.of.yaml service/checkout handlesData data/personal
service/checkout handles data/personal  (derivative dataFromService)
  A service that depends on a second service handles each data class that the second service handles.
  with $called = service/accounts, $data = data/personal, $service = service/checkout
  because:
    service/checkout depends on service/accounts  (written at platform/tests/examples/personal-data-through-a-call.of.yaml:119:5)
    service/accounts handles data/personal  (derivative dataFromTable)
      A service handles each data class that one of its tables holds.
      with $data = data/personal, $service = service/accounts, $table = table/customers
      because:
        service/accounts reads table/customers  (written at platform/tests/examples/personal-data-through-a-call.of.yaml:98:5)
        table/customers holds data/personal  (written at platform/tests/examples/personal-data-through-a-call.of.yaml:74:5)
```

```
$ of-lang explain platform platform/tests/examples/shared-service.of.yaml service/catalog SharedService
service/catalog is a SharedService  (derivative sharedServiceFromCalls)
  A service that more than two services depend on, directly or through other services, is a shared service.
  with $service = service/catalog
```

### `diff`

What changes from one facts file to the other, in three parts: facts stated,
facts derived, and rejections. `+` is added, `-` is gone. A field value is a
fact of its own, so a changed value shows as one `-` and one `+`.

In this example, a copy of the test file makes the search service depend on the
accounts service and increases its availability. The file is
`search-with-accounts.of.yaml`.

```
$ of-lang diff platform platform/tests/examples/personal-data-through-a-call.of.yaml search-with-accounts.of.yaml
stated:
  + service/search depends on service/accounts
  + service/search has availability 99.9
  - service/search has availability 99.5
derived:
  + service/search handles data/personal
rejected:
  + publicServiceDeclaresItsData: The privacy notice must tell that service/search handles data/personal.
```

`diff` always exits with 0: a difference is not a problem.

### `test`

One line per test or question, then a total. A failure says what was
expected and did not happen. In this example, somebody made the checkout
service write the products table, not the payments table, in
`personal-data-through-a-call.of.yaml`:

```
$ of-lang test platform
pass  test backupInWrongRegion  (platform/tests/examples/backup-in-wrong-region.of.yaml)
FAIL  test personalDataThroughACall  (platform/tests/examples/personal-data-through-a-call.of.yaml)
      expected derived, but not: service/checkout handles data/payment
      expected derived, but not: service/checkout reads table/payments
pass  test serviceWithoutDeclaredData  (platform/tests/examples/service-without-declared-data.of.yaml)
pass  test catalogCalledByThreeServices  (platform/tests/examples/shared-service.of.yaml)
pass  question whichPublicServicesHandlePersonalData  (platform/tests/questions/which-public-services-handle-personal-data.of.yaml)
pass  question whichTeamsStopWithTheCoreTeam  (platform/tests/questions/which-teams-stop-with-the-core-team.of.yaml)
pass  test writeIsARead  (platform/tests/terms/write.of.yaml)
6 passed, 1 failed
```

### JSON

`--format json` prints the same results as one JSON object, for other tools
to read. A derived link, for example:

```json
{
  "type": "link",
  "sentence": "service/checkout handles data/personal",
  "from": "service/checkout",
  "link": "handlesData",
  "to": "data/personal"
}
```

and a rejection carries the constraint, its message, and what each variable
stood for:

```json
{
  "constraint": "publicServiceDeclaresItsData",
  "message": "The privacy notice must tell that service/checkout handles data/personal.",
  "bindings": [
    { "var": "$data", "id": "data/personal" },
    { "var": "$service", "id": "service/checkout" }
  ]
}
```

### Derived facts

When the tool writes facts out, a derived fact is a fact with
`derived: true`; `derived` is `false` by default. The flag exists only in
the tool's output: the language has no such key, so a file that writes it is
rejected like any unknown key.

A thing that is partly stated and partly derived has two facts under one id:
the stated one, as written, and one derived fact beside it holding only what
was derived, never repeating the stated one. Its `kind` is the kind a
derivative made the thing, or the stated kind if no derivative did.

```yaml
service/catalog:
  is: Fact
  kind: SharedService              # a derivative made it a shared service
  derived: true

service/accounts:
  is: Fact
  kind: Service                    # the stated kind; no derivative changed it
  derived: true
  links:
    handlesData: [data/personal]   # from dataFromTable
```

The tool does not write facts in this form yet.

## Planned outputs

Not built yet. This page describes them so that the language can be judged
by what it produces.

### The change report

The [pending change](../language/reference/pending.md) of a package returns
one report, for the facts in its tests and examples. The report compares the
package without the change and with the change, so git is not necessary. It joins `diff` and `explain`: what was stated,
what newly follows and why, what no longer follows, and what is rejected.

```
Stated (2)
  + service/checkout depends on service/accounts
  + service/search depends on service/accounts

Derived (3)
  + team/web depends on team/core
      derivative  teamDependsOnTeam
      match   $team = team/web, $service = service/checkout,
              $called = service/accounts, $other = team/core
      links   service/checkout is owned by team/web
              service/checkout depends on service/accounts
              service/accounts is owned by team/core
      filter  team/core does not equal team/web
  + service/checkout handles data/personal
      derivative  dataFromService
  + (1 more, grouped: handlesData)

No longer derived (0)

Rejected (2)
  ✗ publicServiceDeclaresItsData
    The privacy notice must tell that service/checkout handles data/personal.
  ✗ publicServiceDeclaresItsData
    The privacy notice must tell that service/search handles data/personal.
```

Large results are grouped by link, with new kinds of conclusion and
rejections listed first.

### Product's pages

A page per definition for product: the definition, examples and
non-examples, and every constraint and derivative as its
`says`, shown as written.

> **Service** — A program that runs in production and answers requests.
> - Exactly one team owns each service.
> - A service handles each data class that one of its tables holds.
>   *(derivative)*
> - The privacy notice must name each data class that an active public
>   service handles. *(constraint)*
>
> Example: the checkout service. Not a service: a script that a developer
> runs one time and does not keep.

A change to a product-owned package shows the sentences that changed, for
approval on that page.

### The RDF export

RDF (Resource Description Framework), OWL (Web Ontology Language) and SHACL
(Shapes Constraint Language) are generated for exchange and for existing
tools. Nobody writes them.

The export does not translate each statement into one RDF item. One
statement can give many RDF items, or none, and several statements can give
one. The export holds the same knowledge as the package and its
dependencies. The check for this: each question and each test of the
package gets the same result from the exported RDF, with the same example
facts, as from the package.

The export keeps these properties:

- The `from` and `to` of a link are checks. They never give a thing a kind.
- Constraints reject. They never add facts.
- The statements of each package are in their own named graph.
- Each derived fact keeps the reason that it is true.

Each IRI (Internationalized Resource Identifier) is built from the package's
`rdfNamespace` and the definition's name.

`of-lang export <package-dir> [--tags a,b]` writes the export. `--tags`
exports only the definitions carrying at least one of the named
[tags](../language/reference/tag.md); without it, the whole package is
exported. A tagged subset is always complete, because each definition that
a tagged definition refers to has the same tag
([TG5](../language/reference/tag.md#rules)). `of-lang check` rejects a
package that breaks this rule. For example, in the example package,
`PublicService` has the tag `privacy`, so `Percent`, the type of its
`availability` field, must have it too. If somebody removed the tag from
`Percent`:

```
$ of-lang check platform
platform/defs/service.of.yaml:30:10: PublicService has the tag privacy and refers to Percent, which does not have it (TG5)
```

## Planned: a language server

`of-lang lsp` will be a language server, speaking LSP (the Language Server
Protocol) over standard input and output, so any editor that speaks LSP gets
the language with one command to launch. It is one server, depending on no
other language server, including the YAML one.

### Two layers in one server

- **Structure, from a schema.** One JSON Schema per statement type, chosen by
  `is:`, plus one for the manifest (`package.yaml` or `package.yml`), hand-written and embedded in the
  binary. The schemas are the source of truth for the shape of statements:
  every key, its description, and the fixed words it allows (`is at most`,
  `becomes`). The field keywords are part of them, so a field
  definition is checked and completed like any other part of a statement.
  The checker's structural errors (unknown keys, wrong shapes, words not in
  the list) are this validation, mapped back to a line and column.
- **Meaning, from the checker.** Everything that needs names (imports,
  kinds, links, fields, variables, fact ids, docs links) comes from the same
  loader and checker as `of-lang check`.

```
internal/schema/
  package.schema.json
  kind.schema.json
  link.schema.json
  valuetype.schema.json
  constraint.schema.json
  derivative.schema.json
  fact.schema.json
  test.schema.json
  question.schema.json
```

The schemas can also be published for other tools to use; the server never
needs them to.

### What it offers

| Feature | What it does |
|---|---|
| Diagnostics | the check's errors as you type, each with its rule code, linked to that rule in the reference |
| Go to definition | from any name to where it is defined: names in `use:`, `kind:`, `link:`, `extends`, `under`, `type:`, `kind: { name: … }`, fact ids, docs links such as `service.of.yaml#Service`, and a `$variable` to its kind item |
| Find references | where a definition is used, including the docs that link to it or embed it |
| Hover | a key's description from the schema; a definition's `definition`; a link's `reads`; a variable's kind; on a fact, what derivatives add to it |
| Completion | the keys and words the schema allows at that place; after `link:`, the links the file imports; in a link item's `from` and `to`, only variables and kinds that fit that end of the link; in `kind: { name: … }`, only kinds; after `field:`, only fields of that node's kind; a name not yet imported, with its import |
| Code actions | add a missing import; rename a definition and the imports that use it, warning when the name is exported |
| Inlay hints | in test files, the derived facts beside the facts they are about |
| Code lens | run a `Test` or `Question` from above it |

### What it needs from the tool

1. **Positions for every name,** where it is defined and wherever it is used.
2. **A loader that keeps going after errors,** giving a partial model of a
   broken file, since a file being typed is broken most of the time.
3. **Unsaved files,** read from the editor instead of the disk.
4. **Re-checking the package** after each pause in typing; packages are
   small enough that incremental checking can wait.

### Order

1. The schemas, with structural diagnostics, key completion and hover.
2. Diagnostics from the checker, go to definition, hover on names.
3. Completion of names, adding imports, rename.
4. Inlay hints for derived facts, and running tests from the editor.

## Planned: a user interface

`of-lang ui` will start a local server that serves a web app. The web app
edits packages. Each change goes to the files of the package, in whichever
form those files are. The rest of the design is not yet discussed.
