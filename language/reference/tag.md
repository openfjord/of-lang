---
title: Tags — The OpenFjord modeling language
---

# Tags

> Part of [the reference](README.md).

A tag marks definitions that belong together for one purpose, so that only
those are exported to RDF (Resource Description Framework). A lineage tool
needs tables, jobs and the links between them; a catalog needs tables and
databases. Each purpose gets a tag, and each definition lists the tags it
belongs to.

```yaml
# defs/tags.of.yaml
lineage:
  is: Tag
  definition: What a lineage tool needs to trace data from table to table.
```

```yaml
# defs/table.of.yaml
import:
  - from: file:tags.of.yaml
    use: [lineage]

Table:
  is: Kind
  definition: A set of rows with one schema, stored in a database.
  tags: [lineage]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Tag` |
| `definition` | yes | one sentence saying what the tagged definitions are for |

A tag is declared and imported like any other name, so a misspelled tag is
an error rather than a new subset. A package may [export](package.md#exports)
its tags for other packages to use on their own definitions.

## Tagging a definition

Kinds, links and value types take an optional `tags` key: a list of tags.
Nothing else carries tags. Constraints, derivatives and single fields are
always exported with the package as a whole, never in a subset.

Tags only choose what is exported. They change nothing about checking,
derivatives, constraints or tests.

Tags are not passed on: a kind that extends a tagged kind is not tagged
unless it lists the tag itself, and neither is a link `under` a tagged link.

A package tags only its own definitions. `tags` is a key on the definition
itself, so there is no way to tag another package's definitions.

## Exporting a subset

Exporting with tags exports every definition of the package that carries at
least one of them. Exporting without tags exports the whole package.

A definition in the subset may refer to one that is not: a link whose
`from` or `to` kind is untagged, a kind that extends an untagged kind, a
field whose value type is untagged, a link `under` an untagged link, a count
for a link that is not exported. Each such reference is left out of the
export, and a link whose `from` or `to` kind is not exported is left out
whole. The export lists everything it left out.

## Rules

- **TG1.** A tag MUST have a `definition`.
- **TG2.** A tag's name MUST be lowerCamelCase.
- **TG3.** `tags` MAY appear only on a kind, a link or a value type, and
  every tag it lists MUST be a tag the file defines or imports.
- **TG4.** A definition MUST NOT list the same tag twice.
