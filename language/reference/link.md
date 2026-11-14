---
title: Links — The OpenFjord modeling language
---

# Links

> Part of [the reference](README.md).

A link is a named, directed relationship from one kind of thing to another.
It is declared once, as its own entry, usually in the same file as the kind
it comes out of; that kind says how many of it each thing carries
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
| `from` | yes | the kind the link starts at |
| `to` | yes | the kind the link ends at |
| `reads` | yes | the sentence product sees, with `{from}` and `{to}` |
| `definition` | yes | one sentence saying what the link means |
| `under` | no | a broader link this one is a narrower sort of |
| `inverse` | no | a name for reading the link the other way |
| `chains` | no | `true` when a chain of the link implies the link |

A link carries nothing but its two ends: it has no fields. Anything to say
about a link becomes a thing in the middle, such as a `Read` between a job
and a table.

## Rules

- **LK1.** A link's name, and its inverse name, MUST be lowerCamelCase with
  two words or more, saying what it links and how: `readsTable`, never
  `reads` (reason 2 in [the research](../../README.md)).
- **LK2.** A link MUST have a `definition` and a `reads` sentence, and
  `reads` MUST contain both `{from}` and `{to}` (reason 7).
- **LK3.** `from` and `to` MUST name kinds. They are checks: a fact linking a
  team with `readsTable` is rejected. No link ever makes something into a
  kind.
- **LK4.** `under` MUST name a declared link, not an inverse name. The
  narrower link's `from` MUST be the broader link's `from` or extend it, and
  the same for `to`. A link MUST NOT sit under itself, directly or through
  others.
- **LK5.** A link with `chains: true` MUST join a kind to the same kind, or
  to a kind that extends it.

## What links imply

These are the only inferences the language makes without a rule. Each fact
they add is a [derived fact](fact.md#derived-facts) with its reason.

| Declared | Adds |
|---|---|
| `inverse: hasTable` on `inDatabase` | nothing new: `hasTable` reads every `inDatabase` the other way, in patterns, facts and tests alike |
| `under: usesTable` on `readsTable` | a `usesTable` for every `readsTable` |
| `chains: true` on `feedsTable` | `a feedsTable c` whenever `a feedsTable b` and `b feedsTable c` |
