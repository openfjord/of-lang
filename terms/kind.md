---
title: Kinds — A modeling language for software people
---

# Kinds

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A kind reads like an object. It says what it is, gives examples, lists its
fields, and lists the links it carries with how many of each.

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
    format:
      enum: [row, columnar]
    visibility:
      enum: [public, internal]
      default: internal
    retentionDays:
      type: RetentionDays
      optional: true
  links:
    inDatabase: exactly 1
    hasColumn:  at least 1

Database:
  is: Kind
  definition: A named store of tables, owned by one team.
  fields:
    name:
      type: Name
  links:
    ownedByTeam: exactly 1

Job:
  is: Kind
  definition: A unit of work an application runs, reading and writing tables.
  fields:
    name:
      type: Name
  links:
    jobOf: exactly 1
```

A sub-kind names its parent: `PartitionedTable: { is: Kind, extends: Table, … }`.

Every field is required unless it says `optional: true` or has a `default`;
every link is required unless its count says otherwise; anything not listed
is rejected. The counts are constraints, not inference: a database naming two
teams is rejected, never read as "the two teams are the same".
