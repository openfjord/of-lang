---
title: Files and statements — The OpenFjord modeling language
---

# Files and statements

> Part of [the reference](README.md).

Everything in the language is written in YAML files whose names end in
`.of.yaml`, such as `table.of.yaml`. The extension tells editors and tools
which files are in the language. The one exception is a package's manifest,
which is always `package.yaml`. A file is a
mapping: an optional `import` block, then statements, each keyed by its name.

- **FL15.** Every file in `defs/` and `tests/` MUST have a name ending in
  `.of.yaml`. Any other file there is an error, never ignored.

```yaml
import:
  - from: file:names.of.yaml
    use: [Name]

Table:
  is: Kind
  definition: A set of rows with one schema, stored in a database.
  fields:
    name:
      type: Name

inDatabase:
  is: Link
  from: Table
  to: Database
  reads: "{from} is in {to}"
  definition: The database that stores the table.
```

## Statements

Every statement says what it is with `is:`.

| `is:` | Where | Described in |
|---|---|---|
| `Kind` | `defs/` | [Kinds](kind.md) |
| `Link` | `defs/` | [Links](link.md) |
| `ValueType` | `defs/` | [Fields and value types](field.md#value-types) |
| `Constraint` | `defs/` | [Constraints](constraint.md) |
| `Derivative` | `defs/` | [Derivatives](derivative.md) |
| `Fact` | `tests/` | [Facts](fact.md) |
| `Test` | `tests/` | [Tests](test.md) |
| `Question` | `tests/` | [Tests](test.md#questions) |

- **FL1.** Every top-level key other than `import` MUST be a statement: a
  mapping with an `is:` key whose value is one of the table above.
- **FL2.** A statement MUST NOT have a key its page does not list. Anything not
  listed is rejected, so a misspelled key is an error, never ignored.
- **FL3.** A YAML file MUST NOT repeat a key in the same mapping.
- **FL4.** Statements in `defs/` MUST be `Kind`, `Link`, `ValueType`,
  `Constraint` or `Derivative`. Statements in `tests/` MUST be `Fact`, `Test` or
  `Question`.
- **FL5.** No statement MAY be named `import`.

## Names

| Statement | Written as | Examples |
|---|---|---|
| kind, value type | UpperCamelCase | `Table`, `CriticalTable`, `RetentionDays` |
| link, link inverse | lowerCamelCase, two words or more | `inDatabase`, `readsTable`, `hasJob` |
| constraint, derivative, test, question | lowerCamelCase | `crossAppReadNeedsContract` |
| fact | an id: any text, usually a path | `db/sales/table/order_lines` |

- **FL6.** A definition's name MUST be unique within its package, and the
  inverse name a link declares counts as a name of that package.
- **FL7.** Kind and value type names MUST be UpperCamelCase; constraint,
  derivative, test and question names MUST be lowerCamelCase. Link names are
  covered by [LK1](link.md#rules).
- **FL8.** The name of a statement in a test file MUST NOT be the name of a
  definition in the package. Statements in a test file are local to that file.

## Imports

A file may use a name only if it defines it or imports it. The `import`
block sits at the top of the file and lists where each name comes from:

```yaml
import:
  - from: file:table.of.yaml              # a file in this package, relative to this one
    use: [Table, inDatabase]
  - from: file:../shared/team.of.yaml     # going up a folder
    use: [Team]
  - from: pkg:@openfjord/org-model     # another package
    use: [Contract, { name: Team, as: OrgTeam }]
```

| Key | Value |
|---|---|
| `from` | `file:` and a path relative to the importing file, or `pkg:` and a package's full name |
| `use` | a list of names; a name may be `{ name: …, as: … }` to rename it in this file |

- **FL9.** Every name a file uses that it does not define MUST be imported,
  including names from other files of the same package. The top of a file
  shows everything it depends on.
- **FL10.** A `file:` import MUST point at a file inside the same package,
  and each name it uses MUST be defined in that file. `./` is implied; `../`
  is written only to go up a folder.
- **FL11.** A `pkg:` import MUST name a package listed in the importing
  package's `dependsOn`, and MUST use only names that package
  [exports](package.md#exports).
- **FL12.** Two names in one file MUST NOT clash; `as` renames one of them.
- **FL13.** A link's inverse name MAY be imported on its own, like any
  other name: `hasJob`, declared as the inverse of `jobOf`, can be imported
  without `jobOf`.
- **FL14.** Field names are not imported: a field belongs to its kind.

## Ids

Files use plain names — `Table` for a kind, `db/sales/table/order_lines` for
a fact — and imports say where each comes from. An id is a name, not a link
(reason 10 in [the research](../../README.md)): the full IRI (Internationalized
Resource Identifier), built from the package's `rdfNamespace`, appears only
in the exported RDF.
