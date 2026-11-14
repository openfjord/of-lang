---
title: Kinds — The OpenFjord modeling language
---

# Kinds

> Part of [the reference](README.md).

A kind is a kind of thing: a database, a table, an application. It reads like
an object — what it is, examples of it, its fields, and the links it carries
with how many of each.

```yaml
Table:
  is: Kind
  definition: A set of rows with one schema, stored in a database.
  examples:
    - the order_lines table in the sales database
  notExamples:
    - a temporary copy a test run makes and throws away
    - a view computed on read
  fields:
    name:
      type: Name
    visibility:
      enum: [public, internal]
      default: internal
    retentionDays:
      type: RetentionDays
      optional: true
  links:
    inDatabase: exactly 1
    belongsToTeam: at most 1
    feedsTable: any number

CriticalTable:
  is: Kind
  extends: Table
  definition: A table read by so many jobs that a change to it needs extra care.
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Kind` |
| `definition` | yes | one sentence saying what the kind is |
| `extends` | no | the kind this one is a narrower sort of |
| `examples` | no | a list of things that are of this kind, in words |
| `notExamples` | no | a list of things that could be mistaken for it but are not |
| `fields` | no | the kind's [fields](field.md), by name |
| `links` | no | the links a thing of this kind carries, each with a count |

## Rules

- **KD1.** A kind MUST have a `definition`.
- **KD2.** `extends` MUST name a kind, and a kind MUST NOT extend itself,
  directly or through others.
- **KD3.** A kind that extends another has all of its fields and link counts,
  and MAY add its own. A thing of the narrower kind is also a thing of the
  broader one.
- **KD4.** Each key under `links` MUST be a link, or a link's inverse name,
  that starts at this kind or a kind it extends. For an inverse name, the
  link must end there.
- **KD5.** Each count MUST be one of `exactly N`, `at least N`, `at most N`
  or `any number`, where N is a whole number.
- **KD6.** No field MAY be named `kind`: facts use that key to name their
  kind.

## What counts mean

A count is a check, never an inference. A database that names two owning
teams is rejected; it is never read as "the two teams are the same". A link a
kind does not list under `links` cannot be stated on a fact of that kind.

Every field is required unless it says `optional: true` or has a `default`
([FD3](field.md#rules)). Together, these mean a fact must have everything its
kind asks for, and nothing else ([FC3](fact.md#rules)).
