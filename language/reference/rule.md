---
title: Rules — The OpenFjord modeling language
---

# Rules

> Part of [the reference](README.md).

A rule derives facts. It is a [pattern](pattern.md) with what to add for
every match: a link, a kind, or a field value.

```yaml
appDependsOnWriter:
  is: Rule
  says: >
    An application depends on every other application whose jobs write a
    table that one of its own jobs reads.
  match:
    - link: hasJob
      from: {$app: Application}
      to:   {$reader: Job}
    - link: readsTable
      from: {$reader: Job}
      to:   {$table: Table}
    - link: writesTable
      from: {$writer: Job}
      to:   {$table: Table}
    - link: jobOf
      from: {$writer: Job}
      to:   {$other: Application}
  filter:
    - node: $other
      does not equal:
        same as: $app
  derive:
    - link: dependsOnApp
      from: $app
      to:   $other
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Rule` |
| `says` | yes | what the rule concludes, in plain words |
| `match` | see [PT9](pattern.md#rules) | the pattern |
| `filter` | no | conditions that narrow the matches |
| `derive` | yes | a list of facts to add for every match |

## What a rule can derive

| Derives | Written | Example |
|---|---|---|
| a link | `link`, `from`, `to` | `link: dependsOnApp, from: $app, to: $other` |
| a kind | `node`, `becomes` | `node: $table, becomes: CriticalTable` |
| a field value | `node`, `field`, `set to` with `value:` or `same as:` | `node: $table, field: visibility, set to: { value: internal }` |

```yaml
criticalTable:
  is: Rule
  says: A table read by more than two jobs is critical.
  filter:
    - count:
        match:
          - link: readsTable
            from: {$reader: Job}
            to:   {$table: Table}
      is above:
        value: 2
  derive:
    - node: $table
      becomes: CriticalTable
```

`criticalTable` has no top-level `match`: its one variable outside the
nested block, `$table`, ranges over every table ([PT9](pattern.md#rules)).

A rule derives only the direct step when the link chains: `tableFeedsTable`
below adds "a feeds b" for each job that reads a and writes b, and
`chains: true` on `feedsTable` adds the longer chains.

```yaml
tableFeedsTable:
  is: Rule
  says: A table feeds every table that a job reading it writes.
  match:
    - link: readsTable
      from: {$job: Job}
      to:   {$source: Table}
    - link: writesTable
      from: {$job: Job}
      to:   {$target: Table}
  derive:
    - link: feedsTable
      from: $source
      to:   $target
```

## Rules

- **RL1.** A rule MUST have `derive`, with at least one entry.
- **RL2.** A derived link MUST join variables whose kinds the link allows
  ([LK3](link.md#rules)). A derived kind MUST be a kind. A derived field MUST
  be a field of the thing's kind, and its value MUST be one the field allows.
- **RL3.** A rule MUST NOT create things: every variable in `derive` MUST be
  bound by the pattern.
- **RL4.** A field MUST end up with one value. Two rules that could give the
  same field different values are rejected before anything runs. If one
  rule's matches still give a thing's field two values, running stops with an
  error that shows both explanations.
- **RL5.** A rule MUST NOT have `require` or `message`. Rules reject nothing.
