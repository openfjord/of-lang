---
title: Constraints — A modeling language for software people
---

# Constraints

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A constraint rejects facts. It is a [pattern](pattern.md) with what every
match must satisfy, written under `require`:

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

Read it as: for every `$job`, `$jobApp`, `$table`, `$writer` and `$writerApp`
joined by these links, where the two applications differ, the table must be
public and at least one contract must link them.

**Without `require`, every match is a violation.** The pattern then describes
the forbidden situation itself. Here `offersContract` links an application to
the contracts it publishes:

```yaml
noSelfConsumption:
  is: Constraint
  says: An application must not consume a contract it offers itself.

  match:
    - link: consumesContract
      from: {$app: Application}
      to:   {$contract: Contract}
    - link: offersContract
      from: {$app: Application}
      to:   {$contract: Contract}

  message: "{$app} must not consume {$contract}, a contract it offers itself."
```

The two forms mean the same: "every match must satisfy this" is "nothing
matches without satisfying this". Use `require` when the rule reads as what
must hold, and leave it out when it reads as what must not happen.

## Messages

`message` states the broken rule, worded with "must" or "must not", and
quotes variables in braces, filled in from the failing match:
`{$job}` gives the thing's id, `{$table: visibility}` a field's value. Only
variables of the top-level match may be quoted.
