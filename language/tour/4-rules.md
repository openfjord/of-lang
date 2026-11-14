---
title: 4. Rules — Tour
---

# 4. Rules

> Part of [the tour](README.md). Previous: [3. Links](3-links.md).
> Next: [5. Constraints](5-constraints.md).

A rule derives what follows from the facts, so nobody has to write it down.

## A first rule

A table belongs to the team that owns its database. First declare the link,
in `table.yaml`, and let tables carry it:

```yaml
belongsToTeam:
  is: Link
  from: Table
  to: Team
  reads: "{from} belongs to {to}"
  definition: The team that answers for the table, through its database.
```

```yaml
  links:
    inDatabase: exactly 1
    belongsToTeam: at most 1
```

`table.yaml` now uses `Team` and `ownedByTeam`, so it imports them:

```yaml
import:
  - from: file:names.yaml
    use: [Name, RetentionDays]
  - from: file:team.yaml
    use: [Team]
  - from: file:database.yaml
    use: [Database, ownedByTeam]
```

Then the rule, also in `table.yaml`, next to what it is about:

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
```

Read it as: for every table `$table` in a database `$database` that a team
`$team` owns, add "`$table` belongs to `$team`".

- `says` is the sentence product reads. It is plain words, with no variables.
- `match` lists the links that must exist. Each end names a variable and its
  kind, every time it appears.
- `derive` lists what to add for every match. Here a variable is written bare:
  its kind is already known.

Run the rules over the example facts:

```
$ of-lang run data-model data-model/tests/examples/a-database.yaml
derived:
  db/sales/table/order_lines belongs to team/commerce  (rule tableTeamFromDatabase)
rejected:
  nothing
```

and ask why:

```
$ of-lang explain data-model data-model/tests/examples/a-database.yaml db/sales/table/order_lines belongsToTeam team/commerce
db/sales/table/order_lines belongs to team/commerce  (rule tableTeamFromDatabase)
  A table belongs to the team that owns its database.
  with $database = db/sales, $table = db/sales/table/order_lines, $team = team/commerce
  because:
    db/sales/table/order_lines is in db/sales  (written at data-model/tests/examples/a-database.yaml:20:3)
    db/sales is owned by team/commerce  (written at data-model/tests/examples/a-database.yaml:14:3)
```

Every derived fact has an explanation like this, down to facts written in a
file. Test it by adding it to the test's `expect`:

```yaml
  expect:
    derived:
      - link: belongsToTeam
        from: db/sales/table/order_lines
        to:   team/commerce
```

## Jobs, `under` and `chains`

`defs/job.yaml` adds jobs, which read and write tables:

```yaml
usesTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} uses {to}"
  definition: The job reads or writes the table when it runs.

readsTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} reads {to}"
  definition: The job reads rows from the table when it runs.
  under: usesTable
  inverse: readByJob
```

`writesTable` is the same, `under: usesTable` too. `under` makes every
`readsTable` a `usesTable` as well: a specific link grouped under a general
one. A question about "what does this job use" now finds both.

Rows flow from table to table through jobs. Declare the link with
`chains: true`, and a rule for the direct step:

```yaml
feedsTable:
  is: Link
  from: Table
  to: Table
  reads: "{from} feeds {to}"
  definition: Rows of the first table flow, directly or through other tables, into the second.
  chains: true

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

The rule adds each direct step; `chains` adds the longer ones. If a feeds b
and b feeds c, a feeds c.

## A rule that counts

A rule can also give a thing a narrower kind. Declare the kind:

```yaml
CriticalTable:
  is: Kind
  extends: Table
  definition: A table read by so many jobs that a change to it needs extra care.
```

and the rule:

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

`filter` holds conditions that narrow the matches. `count` counts how often a
nested pattern matches, here for each table. The rule has no top-level
`match`, so `$table` ranges over every table.

```
$ of-lang run data-model data-model/tests/examples/critical-table.yaml
derived:
  db/sales/table/customers is a CriticalTable  (rule criticalTable)
  ...
```

Rules repeat until nothing new follows, so one rule's result can feed
another. They never remove a fact or create a thing, so this always ends.
