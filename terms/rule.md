---
title: Rules — A modeling language for software people
---

# Rules

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A rule derives facts. It is a [pattern](pattern.md) with what to add for
every match.

A rule's `derive` lists facts to add for every match: a link (`link`, `from`,
`to`), a kind (`node` and `becomes`), or a field value (`node`, `field` and
`set to`). A derived field must end up with one value; two rules giving the
same thing two different values for one field is an error, reported with both
explanations.

```yaml
tableTeamFromDatabase:
  is: Rule
  says: A table belongs to the team that owns its database.

  match:
    - link: inDatabase
      from: {$table: Table}
      to:   {$database: Database}
    - link: ownedByTeam
      from: {$database: Database}
      to:   {$team: Team}

  derive:
    - link: belongsToTeam
      from: $table
      to:   $team

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

criticalTable:
  is: Rule
  says: A table read by more than ten jobs is critical.

  filter:
    - count:
        match:
          - link: readsTable
            from: {$reader: Job}
            to:   {$table: Table}
      is above:
        value: 10

  derive:
    - node: $table
      becomes: CriticalTable
```

`feedsTable` is declared with `chains: true`, so the rule needs to derive only
the direct step; the chain follows. `criticalTable` has no top-level `match`:
with one variable outside the nested block, it matches every thing of the
kind that variable's mentions give.
