---
title: Packages — The OpenFjord modeling language
---

# Packages

> Part of [the reference](README.md).

A package is the unit that is versioned, owned, reviewed and depended on. It
holds definitions — kinds, links, value types, constraints and rules — with
the docs that explain them and the tests that prove them.

## Layout

```
data-model/
  package.yaml                     # the manifest
  defs/                            # definitions
    names.yaml                     # value types
    team.yaml                      # Team
    database.yaml                  # Database, and the links that come out of it
    table.yaml                     # Table, its links, the rules about tables
    job.yaml
    application.yaml
    reading-across-applications.yaml   # a constraint spanning several kinds
  docs/                            # Markdown that explains the definitions
    tables.md
  tests/                           # examples and questions
    examples/
    questions/
```

- **PK1.** A package's top level MUST hold `package.yaml` and `defs/`, MAY
  hold `docs/` and `tests/`, and MUST NOT hold anything else.
- **PK2.** Inside `defs/`, `docs/` and `tests/`, the layout is the author's.
  A file in `defs/` MAY hold definitions of any sort in any mix.
- **PK3.** Inside `defs/`, a definition SHOULD sit next to what it is about:
  a link with the kind it comes out of, a constraint or rule with the kind it
  concerns, and one that concerns several kinds a level up, beside them all.

## The manifest

```yaml
package: "@openfjord/data-model"  # namespace, then name; quoted, since YAML reserves a leading @
version: 0.3.0
owner: product                    # who must approve any change to the package
rdfNamespace: data                # used only in exported RDF, never in the files
dependsOn:
  "@openfjord/org-model": ^1.2    # semantic-version range, pinned in a lock file
exports:                          # what other packages may use; the rest is internal
  kinds: [Database, Table, Application, Job]
  links: [inDatabase, readsTable, writesTable, jobOf]
  valueTypes: [Name]
openFor:                          # extension points: what other packages may add to
  - extend: Table                 #   another package may declare a kind of table
  - under: readsTable             #   another package may declare a narrower link under readsTable
```

| Key | Required | Value |
|---|---|---|
| `package` | yes | the package's full name |
| `version` | yes | a semantic version, `MAJOR.MINOR.PATCH` |
| `owner` | yes | who must approve any change to the package |
| `rdfNamespace` | yes | the prefix used to build IRIs in the exported RDF |
| `dependsOn` | no | a mapping from package names to version ranges |
| `exports` | no | the `kinds`, `links` and `valueTypes` other packages may import |
| `openFor` | no | a list of extension points, each `extend: <kind>` or `under: <link>` |

- **PK4.** `package.yaml` MUST have `package`, `version`, `owner` and
  `rdfNamespace`, and MUST NOT have any other key than those in the table.
- **PK5.** `version` MUST be `MAJOR.MINOR.PATCH`, three whole numbers.

## Names and namespaces

A package name may start with a namespace, as in npm: `@openfjord/data-model`
is the package `data-model` in the namespace `@openfjord`.

- **PK6.** A package name MUST be lowercase letters, digits and hyphens,
  starting with a letter, optionally preceded by a namespace: `@`, then the
  same characters, then `/`.
- **PK7.** `dependsOn` and `pkg:` imports MUST use the full name, namespace
  included.

## Owners

A package has one owner, and everything inside it — definitions, docs and
tests — is owned by that owner, who must approve any change to it, product or
not (reasons 7 and 8 in [the research](../../README.md)).

- **PK8.** Nothing inside a package has an owner of its own. A product term
  and a technical term therefore MUST live in different packages.

## Dependencies

A version range is one of:

| Range | Means |
|---|---|
| `1.2.3` | exactly 1.2.3 |
| `^1.2`, `^1.2.3` | at least that version, below the next major (below the next minor while the major is 0) |
| `~1.2`, `~1.2.3` | at least that version, below the next minor |

- **PK9.** Every package in `dependsOn` MUST be found, and its version MUST
  fall in the range.
- **PK10.** Packages MUST NOT depend on each other in a circle.

## Exports

A package's definitions are internal unless it exports them.

- **PK11.** Every name in `exports` MUST be a definition of the package, of
  the sort it is listed under. A link's inverse name MUST NOT be listed; it is
  exported with its link.
- **PK12.** A package MUST state things only about its own definitions,
  except through the extension points its dependencies open in `openFor`.
- **PK13.** Every name in `openFor` MUST be a definition of the package:
  `extend` names a kind, `under` names a link.
- **PK14.** A new version MUST match what changed: removing or renaming a
  definition, or tightening a constraint, needs a new major version; adding a
  definition or loosening a constraint does not.

## Docs

Why, how and context live in Markdown files in `docs/` (reason 5). A doc
links to a definition or embeds it, so prose and definitions read side by
side:

````markdown
# Tables

A table holds rows with one schema. Every table sits in exactly one database:

![](../defs/table.yaml#Table)

Teams kept breaking each other by reading internal tables, so reading across
applications needs a contract — see
[crossAppReadNeedsContract](../defs/reading-across-applications.yaml#crossAppReadNeedsContract).
````

| Written | Means |
|---|---|
| `[text](path/file.yaml#Name)` | a link to the entry `Name` |
| `![](path/file.yaml#Name)` | the entry `Name`, embedded where it stands |
| `![](path/file.yaml)` | the whole file, embedded |

Both are standard Markdown, so a viewer that does not know the convention
still shows a working link. The part after `#` is an entry's name, not a line
number, so an embed follows its definition wherever it moves in the file.
Each definition's generated page lists the docs that link to it or embed it.

- **DC1.** Paths MUST be relative to the doc.
- **DC2.** A link or embed MUST point at a file that exists, and the part
  after `#`, when present, MUST be the name of an entry in that file.
