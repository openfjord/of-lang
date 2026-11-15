---
title: Links — The OpenFjord modeling language
---

# Links

> Part of [the reference](README.md).

A link is a named relationship from a kind to a kind, in one direction. A link
is declared one time, as its own statement. Usually, it is in the same file as
the kind where it starts. That kind tells how many of the link each thing has
([KD4](kind.md#rules)).

```yaml
inDatabase:
  is: Link
  from: Table
  to: Database
  reads: "{from} is in {to}"
  definition: The database that stores the table.
  inverse: hasTable

readsTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} reads {to}"
  definition: The job reads rows from the table when it runs.
  under: usesTable                   # every readsTable is also a usesTable
  inverse: readByJob

feedsTable:
  is: Link
  from: Table
  to: Table
  reads: "{from} feeds {to}"
  definition: >
    Rows of the first table flow, directly or through other tables, into the
    second.
  chains: true                       # a feeds b and b feeds c means a feeds c
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Link` |
| `from` | yes | the kind where the link starts |
| `to` | yes | the kind where the link ends |
| `reads` | yes | the sentence that product sees, with `{from}` and `{to}` |
| `definition` | yes | one sentence that tells what the link means |
| `under` | no | the broader link; this link is a narrower sort of it |
| `inverse` | no | a name for the link in the opposite direction |
| `chains` | no | `true` when a chain of the link implies the link |
| `tags` | no | the [tags](tag.md) that the link belongs to |

A link has only its two ends. It has no fields. Any data about a link becomes
a thing in the middle, for example a `Read` between a job and a table.

## Rules

- **LK1.** A link's name, and its inverse name, MUST be lowerCamelCase, with
  two or more words that tell what the link connects and how: `readsTable`,
  never `reads` (reason 2 in [the research](../../README.md)).
- **LK2.** A link MUST have a `definition` and a `reads` sentence. `reads`
  MUST contain both `{from}` and `{to}` (reason 7).
- **LK3.** `from` and `to` MUST name kinds. They are checks: the language
  rejects a fact that links a team with `readsTable`. A link never gives a
  thing a kind.
- **LK4.** `under` MUST name a declared link, not an inverse name. The `from`
  of the narrower link MUST be the `from` of the broader link, or a kind that
  extends it. The same rule applies to `to`. A link MUST NOT be under itself,
  directly or through other links.
- **LK5.** A link with `chains: true` MUST connect a kind to the same kind, or
  to a kind that extends it.

## What links imply

These are the only inferences that the language makes without a derivative.
Each fact that they add is a [derived fact](fact.md#derived-facts), and
records its reason.

| Declared | Adds |
|---|---|
| `inverse: hasTable` on `inDatabase` | nothing new: `hasTable` is every `inDatabase` in the opposite direction, in all derivatives, constraints, questions and facts |
| `under: usesTable` on `readsTable` | a `usesTable` for every `readsTable` |
| `chains: true` on `feedsTable` | `a feedsTable c` each time that `a feedsTable b` and `b feedsTable c` |
