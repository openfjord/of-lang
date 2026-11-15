---
title: Kinds — The OpenFjord modeling language
---

# Kinds

> Part of [the reference](README.md).

A kind is a category of things, for example a database, a table or an
application. A kind tells what the thing is and gives examples of it. It also
lists the fields of the thing, and the links that the thing has, with a count
for each link.

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
| `definition` | yes | one sentence that tells what the kind is |
| `extends` | no | the broader kind; this kind is a narrower sort of it |
| `examples` | no | a list of things of this kind, in words |
| `notExamples` | no | a list of things that look like things of this kind, but are not |
| `fields` | no | the [fields](field.md) of the kind, by name |
| `links` | no | the links that a thing of this kind has, each with a count |
| `tags` | no | the [tags](tag.md) that the kind belongs to |

## Rules

- **KD1.** A kind MUST have a `definition`.
- **KD2.** `extends` MUST name a kind. A kind MUST NOT extend itself,
  directly or through other kinds.
- **KD3.** A kind that extends a different kind has all the fields and link
  counts of that kind. It MAY add its own fields and link counts. A thing of
  the narrower kind is also a thing of the broader kind.
- **KD4.** Each key under `links` MUST be a link, or the inverse name of a
  link, that starts at this kind or at a kind that this kind extends. For an
  inverse name, the link must end at that kind.
- **KD5.** Each count MUST be one of `exactly N`, `at least N`, `at most N`
  or `any number`, where N is a whole number.
- **KD6.** No field MAY have the name `kind`, because facts use that key to
  name their kind.

## What counts mean

A count is a check. It is never an inference. A database that names two
owning teams is an error. The language never infers that "the two teams are
the same". A fact cannot state a link that its kind does not list under
`links`.

Every field is required, unless it has `optional: true` or a `default`
([FD3](field.md#rules)). These two rules together mean that a fact must have
all the items that its kind asks for, and nothing more
([FC3](fact.md#rules)).
