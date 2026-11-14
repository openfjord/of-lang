---
title: Constraints — The OpenFjord modeling language
---

# Constraints

> Part of [the reference](README.md).

A constraint rejects facts. It is a [pattern](pattern.md) with what every
match must satisfy, or, without `require`, a pattern whose every match is a
violation.

## With `require`: what must hold

```yaml
crossAppReadNeedsContract:
  is: Constraint
  says: >
    A job of one application may read a table that a job of another
    application writes, only if the table is public and the reading
    application consumes a contract that covers the table.
  match:
    - link: jobOf
      from: {$job: Job}
      to:   {$jobApp: Application}
    - link: readsTable
      from: {$job: Job}
      to:   {$table: Table}
    - link: writesTable
      from: {$writer: Job}
      to:   {$table: Table}
    - link: jobOf
      from: {$writer: Job}
      to:   {$writerApp: Application}
  filter:
    - node: $writerApp
      does not equal:
        same as: $jobApp
  require:
    - node:  $table
      field: visibility
      equals:
        value: public
    - at least one:
        match:
          - link: consumesContract
            from: {$jobApp: Application}
            to:   {$contract: Contract}
          - link: coversTable
            from: {$contract: Contract}
            to:   {$table: Table}
  message: "{$job} must not read {$table} from {$writerApp} without a contract."
```

Read it as: for every job reading a table that a job of another application
writes, the table must be public and a contract the reading application
consumes must cover it. A match that fails either requirement is a
violation.

## Without `require`: what must not happen

```yaml
noFeedbackLoop:
  is: Constraint
  says: A table must not feed itself, directly or through other tables.
  match:
    - link: feedsTable
      from: {$table: Table}
      to:   {$table: Table}
  message: "{$table} must not feed itself."
```

Every match is a violation. Because `feedsTable` chains, this also catches a
table that feeds itself through others.

The two forms mean the same: "every match must satisfy this" is "nothing
matches without satisfying this". Use `require` when the rule reads as what
must hold, and leave it out when it reads as what must not happen.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Constraint` |
| `says` | yes | the rule, in plain words |
| `match` | see [PT9](pattern.md#rules) | the pattern |
| `filter` | no | conditions that narrow the matches |
| `require` | no | conditions every match must satisfy |
| `message` | yes | what a violation says |

## Messages

A message states the broken rule, worded with "must" or "must not". It quotes
variables in braces, filled in from the violating match: `{$job}` gives the
thing's id, `{$table: visibility}` a field's value.

## Rules

- **CN1.** A constraint MUST have a `message`.
- **CN2.** A message SHOULD be worded with "must" or "must not".
- **CN3.** A message MUST quote only variables bound by the top-level
  pattern, not variables that belong to a nested block.
- **CN4.** A constraint MUST NOT have `derive`. Constraints add nothing.
- **CN5.** Each match that fails a requirement is reported once, with its
  message and its explanation: which requirement failed and the facts the
  match used.
