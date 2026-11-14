---
title: Packages — A modeling language for software people
---

# Packages

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A package holds definitions: kinds, links, value types,
constraints, rules and mappings. Facts appear only in its tests and examples
(reason 14); facts written at runtime are deferred.

**Ownership is at the package level.** A package has one owner, and
everything inside it — definitions, docs and tests — is owned by that owner,
who must approve any change to it, product or not. Nothing inside a package
has an owner of its own, so a product term and a technical term live in
different packages (reasons 7 and 8). The owner is named in the package's
manifest, `package.yaml`.

**The top level of a package is fixed:** `package.yaml` and three folders.

| At the top | Holds |
|---|---|
| `package.yaml` | the manifest: name, version, owner, dependencies, exports |
| `defs/` | definitions: kinds, links, value types, constraints, rules and mappings |
| `docs/` | Markdown that explains the definitions |
| `tests/` | examples and questions, with the facts they use |

Inside each folder the layout is the author's. A file in `defs/` may hold
definitions of any sort in any mix. For example:

```
data-model/
  package.yaml
  defs/
    job.yaml                       # Job, and the links that come out of it
    table.yaml                     # Table, its links, the rules about tables
    database.yaml
    application.yaml               # Application, its links, the rules about applications
    reading-across-applications.yaml
                                   # a constraint spanning jobs, tables and applications
  docs/
    tables.md                      # why tables work as they do; embeds parts of table.yaml
    reading-across-applications.md
  tests/
    questions/                     # questions with expected answers
    examples/                      # facts with what they must derive and reject
```

**Inside `defs/`, place things next to what they are about.** A link sits
with the kind it comes out of. A constraint, rule or mapping sits with the
kind it concerns; one that concerns several sits a level up, where it can be
seen beside all of them. This is guidance for the author, not a check: only
the top level is fixed.

## Docs

**Docs live in `docs/`.** Why, how and context live in Markdown files in
`docs/`, beside `defs/`, laid out however the author likes (reason 5). A doc
links to a definition or embeds it, so prose and definitions read side by
side, like a notebook:

````markdown
# Tables

A table holds rows with one schema. Every table sits in exactly one database:

![](../defs/table.yaml#Table)

Teams kept breaking each other by reading internal tables, so reading across
applications needs a contract — see
[crossAppReadNeedsContract](../defs/reading-across-applications.yaml#crossAppReadNeedsContract).
````

- `[text](path/to/file.yaml#Name)` links to a definition.
- `![](path/to/file.yaml#Name)` embeds that entry where it stands;
  `![](path/to/file.yaml)` embeds the whole file.
- Paths are relative to the doc, as in any Markdown file.
- The part after `#` is an entry's name, not a line number, so an embed
  follows its definition wherever it moves in the file.
- Both are standard Markdown. A viewer that does not know the convention still
  shows a working link to the file.
- A link or embed pointing at a missing file or entry is an error. Each
  definition's generated page lists the docs that link to it or embed it.

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
  links: [inDatabase, hasTable, readsTable, writesTable, jobOf]
  valueTypes: [Name]
openFor:                          # extension points: what other packages may add to
  - extend: Table                 #   another package may declare a kind of table
  - under: readsTable             #   another package may declare a narrower link under readsTable
```

**A package name may start with a namespace,** as in npm:
`@openfjord/data-model` is the package `data-model` in the namespace
`@openfjord`. The full name, namespace included, is what `dependsOn` and
`pkg:` imports use. Names and namespaces are lowercase letters, digits and
hyphens.

Checks on packages, all in CI:

- Every package names its owner.
- A package's top level is `package.yaml`, `defs/`, `docs/` and `tests/`.
- A package uses only terms it defines or that a dependency exports.
- A package states things only about its own terms, except through the
  extension points its dependencies open (reason 11).
- Dependencies follow allowed directions; a cycle is an error.
- Version changes match what changed: removing or renaming a term, or
  tightening a constraint, is a breaking change; adding a term or loosening a
  constraint is not.
