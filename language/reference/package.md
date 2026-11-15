---
title: Packages — The OpenFjord modeling language
---

# Packages

> Part of [the reference](README.md).

A package is the unit that people version, own, review and depend on. It
holds definitions: kinds, links, value types, constraints and derivatives. It
also holds the docs that explain these definitions and the tests that prove
them.

## Layout

```
data-model/
  package.yaml                     # the manifest
  defs/                            # definitions
    names.of.yaml                     # value types
    team.of.yaml                      # Team
    database.of.yaml                  # Database, and the links that come out of it
    table.of.yaml                     # Table, its links, the derivatives about tables
    job.of.yaml
    application.of.yaml
    reading-across-applications.of.yaml   # a constraint spanning several kinds
  docs/                            # Markdown that explains the definitions
    tables.md
  tests/                           # examples and questions
    examples/
    questions/
```

- **PK1.** The top level of a package MUST hold `package.yaml` and `defs/`, MAY
  hold `docs/` and `tests/`, and MUST NOT hold anything else.
- **PK2.** In `defs/`, `docs/` and `tests/`, the author chooses the layout.
  Files in `defs/` and `tests/` have names that match `*.of.yaml`
  ([FL15](files.md#files-and-statements)). A file in `defs/` MAY hold
  definitions of any sort, in any combination.
- **PK3.** In `defs/`, a definition SHOULD be next to the item that it is
  about:
  - a link with the kind where the link starts
  - a constraint or derivative with the kind that it concerns
  - a constraint or derivative that concerns several kinds one level higher,
    next to all of those kinds.

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
  tags: [lineage]
openFor:                          # extension points: what other packages may add to
  - extend: Table                 #   another package may declare a kind of table
  - under: readsTable             #   another package may declare a narrower link under readsTable
```

| Key | Required | Value |
|---|---|---|
| `package` | yes | the full name of the package |
| `version` | yes | a semantic version, `MAJOR.MINOR.PATCH` |
| `owner` | yes | who must approve any change to the package |
| `rdfNamespace` | yes | the prefix that the tool uses to build IRIs in the exported RDF |
| `dependsOn` | no | a mapping from package names to version ranges |
| `exports` | no | the `kinds`, `links`, `valueTypes` and `tags` that other packages can import |
| `openFor` | no | a list of extension points, each `extend: <kind>` or `under: <link>` |

- **PK4.** `package.yaml` MUST have `package`, `version`, `owner` and
  `rdfNamespace`, and MUST NOT have a key that is not in the table.
- **PK5.** `version` MUST be `MAJOR.MINOR.PATCH`, three whole numbers.

## Names and namespaces

A package name can start with a namespace, as in npm: `@openfjord/data-model`
is the package `data-model` in the namespace `@openfjord`.

- **PK6.** A package name MUST be lowercase letters, digits and hyphens,
  with a letter first. An optional namespace can come before the name: `@`,
  then the same characters, then `/`.
- **PK7.** `dependsOn` and `pkg:` imports MUST use the full name, with the
  namespace.

## Owners

A package has one owner. That owner owns everything in the package:
definitions, docs and tests. The owner must approve each change to the package
(reasons 7 and 8 in [the research](../../README.md)). The users choose which
packages hold which terms.

- **PK8.** No item in a package has its own owner. Thus, terms with different
  owners MUST be in different packages.

## Dependencies

A version range is one of:

| Range | Means |
|---|---|
| `1.2.3` | exactly 1.2.3 |
| `^1.2`, `^1.2.3` | at least that version, below the next major version (below the next minor version when the major version is 0) |
| `~1.2`, `~1.2.3` | at least that version, below the next minor version |

- **PK9.** The tool MUST find every package in `dependsOn`, and the version
  of each package MUST be in the range.
- **PK10.** Packages MUST NOT depend on each other in a circle.

## Exports

The definitions of a package are internal if the package does not export
them.

- **PK11.** Every name in `exports` MUST be a definition of the package, and
  of the sort that its list shows. The inverse name of a link MUST NOT be in
  `exports`. The package exports it with its link.
- **PK12.** A package MUST state things only about its own definitions,
  except through the extension points that its dependencies open in
  `openFor`.
- **PK13.** Every name in `openFor` MUST be a definition of the package:
  `extend` identifies a kind, and `under` identifies a link.
- **PK14.** A new version MUST match the change:
  - A change that removes or renames a definition, or makes a constraint
    stricter, needs a new major version.
  - A change that adds a definition or makes a constraint less strict does not
    need a new major version.

## Docs

Markdown files in `docs/` give the reasons, the methods and the context
(reason 5). A doc has a Markdown link to a definition or embeds it. Thus, a
reader sees the text and the definitions together:

````markdown
# Tables

A table holds rows with one schema. Every table sits in exactly one database:

![](../defs/table.of.yaml#Table)

Teams kept breaking each other by reading internal tables, so reading across
applications needs a contract — see
[crossAppReadNeedsContract](../defs/reading-across-applications.of.yaml#crossAppReadNeedsContract).
````

| Written | Means |
|---|---|
| `[text](path/file.of.yaml#Name)` | a Markdown link to the statement `Name` |
| `![](path/file.of.yaml#Name)` | the statement `Name`, embedded at that position |
| `![](path/file.of.yaml)` | the full file, embedded |

Both forms are standard Markdown. Thus, a viewer that does not know this
convention still shows a Markdown link that works. The part after `#` is the
name of a statement, not a line number. Thus, an embed always shows its
definition, also when the definition moves in the file. The generated page of
each definition lists the docs that have a Markdown link to it or embed it.

- **DC1.** Paths MUST be relative to the doc.
- **DC2.** A Markdown link or an embed MUST point to a file that exists. If
  there is a part after `#`, it MUST be the name of a statement in that file.
