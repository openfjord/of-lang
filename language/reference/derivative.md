---
title: Derivatives — The OpenFjord modeling language
---

# Derivatives

> Part of [the reference](README.md).

A derivative adds facts that follow from other facts. It finds things and the
links between them, narrows what it found with conditions, and for every
match adds a link, a kind or a field value.

```yaml
tableTeamFromDatabase:
  is: Derivative
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

Four principles shape how derivatives are written:

1. **Readability comes before brevity.** Every part is spelled out in YAML
   keys and words. There are no expression strings to parse.
2. **No implicit variables.** Every thing a derivative talks about has a name
   and a kind. Nothing is reached through a dotted path that skips over an
   unnamed thing.
3. **Variables start with `$`,** so a variable is never mistaken for a kind, a
   link, a field or a value.
4. **The sentence is plain words.** `says` is what product reads, so it names
   no variables. Variables belong to the logic.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Derivative` |
| `says` | yes | what the derivative concludes, in plain words |
| `match` | yes, unless [DV11](#rules) allows leaving it out | the links to find |
| `filter` | no | conditions that narrow the matches |
| `derive` | yes | what to add for every match |

## Match

Each item in `match` is one link that must exist between two things:

```yaml
match:
  - link: readsTable
    from: {$job: Job}
    to:   {$table: Table}
```

`from` and `to` each give a variable and the kind it must be. A variable is
`$` followed by a lowercase letter, then letters and digits: `$job`,
`$writerApp`.

`match` finds every combination of things for which all its links hold. Each
combination is one match: a row giving every variable one thing. Two matches
that give the same things to the same variables are one. Two variables may
stand for the same thing unless a condition says otherwise.

A link's inverse name reads it the other way: `link: writtenByJob, from:
{$table: Table}, to: {$writer: Job}` matches exactly what `link: writesTable,
from: {$writer: Job}, to: {$table: Table}` does.

`repeated: true` matches the link one or more times in a row. The things in
between are not named, so a longer chain between the same two things is
still one match; explanations list the steps.

```yaml
- link: feedsTable
  from: {$source: Table}
  to:   {$target: Table}
  repeated: true
```

A narrower kind narrows the matches: `{$table: CriticalTable}` on
`readsTable` matches only jobs reading critical tables.

## Filter

`filter` receives the matches one at a time and keeps those for which every
condition holds. A condition is about a thing itself (`node` alone) or one of
its fields (`node` and `field`), followed by a comparison word and what to
compare with: `value:` (a literal), `values:` (a list, for `is one of`) or
`same as:` (a variable, or a variable and a field).

```yaml
filter:
  - node: $writerApp
    does not equal:
      same as: $jobApp
  - node:  $table
    field: retentionDays
    is at most:
      same as:
        node:  $database
        field: retentionDays
  - node:  $table
    field: visibility
    is one of:
      values: [public, internal]
```

| Word | Holds when |
|---|---|
| `equals`, `does not equal` | the two are equal, or not |
| `is above`, `is below` | the first is greater, or smaller |
| `is at least`, `is at most` | the first is greater or equal, or smaller or equal |
| `is one of` | the first equals one of `values` |
| `is a` | the thing is of the kind named (`is a: CriticalTable`) |
| `is set`, `is not set` | the field has a value, or has none |

A thing compared without a field compares as its id. A field with no value
fails every comparison except `is not set`.

A condition cannot see links directly. It reaches them through a nested
block, which has its own `match` and optional `filter` and runs once for each
match of the derivative:

| Word | Holds when |
|---|---|
| `at least one` | the nested block matches at least once |
| `none` | the nested block never matches |
| `count` | the number of its matches compares as the comparison word beside it says |
| `any of` | at least one of the listed conditions holds |

There is no general `not`: `does not equal`, `is not set` and `none` cover
the negative cases while keeping each condition readable on its own.

## Derive

| Adds | Written | Example |
|---|---|---|
| a link | `link`, `from`, `to` | `link: dependsOnApp, from: $app, to: $other` |
| a kind | `node`, `becomes` | `node: $table, becomes: CriticalTable` |
| a field value | `node`, `field`, `set to` with `value:` or `same as:` | `node: $table, field: visibility, set to: { value: internal }` |

```yaml
appDependsOnWriter:
  is: Derivative
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

```yaml
criticalTable:
  is: Derivative
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

`criticalTable` has no `match` of its own: its one variable outside the
nested block, `$table`, ranges over every table ([DV11](#rules)).

When a link chains, a derivative adds only the direct step: `tableFeedsTable`
adds "a feeds b" for each job that reads a and writes b, and `chains: true`
on `feedsTable` adds the longer chains.

```yaml
tableFeedsTable:
  is: Derivative
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

- **DV1.** A derivative MUST have a `says` sentence, and the sentence MUST
  NOT contain a `$` word.
- **DV2.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file defines or imports.
- **DV3.** `from` and `to` MUST each be written `{$variable: Kind}`, every
  time, in nested blocks too, and every mention of a variable MUST give the
  same kind.
- **DV4.** The kind MUST be the kind the link joins at that end, or a kind
  that extends it.
- **DV5.** Conditions and `derive` MUST write a variable bare (`$table`), and
  every variable they use MUST appear in a `match` item.
- **DV6.** A variable that appears only inside a nested block belongs to that
  block and MUST NOT be used outside it.
- **DV7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **DV8.** A comparison MUST compare values of the same type: two fields
  that name the same value type, or two fields that name none and have the
  same built-in type. A literal MUST be a value the field allows.
- **DV9.** A field named in a condition, beside `node` or inside `same as`,
  MUST be a field of that node's kind or of a kind it extends.
- **DV10.** A condition comparing two nodes without fields MUST compare
  variables whose kinds are the same, or where one kind extends the other.
- **DV11.** `match` MAY be left out when every variable the derivative uses
  appears in a nested block; the derivative then ranges over every thing of
  each such variable's kind.
- **DV12.** A derivative MUST have `derive`, with at least one item.
- **DV13.** A derived link MUST join variables whose kinds the link allows
  ([LK3](link.md#rules)). A derived kind MUST be a kind. A derived field MUST
  be a field of the thing's kind, and its value MUST be one the field allows.
- **DV14.** A derivative MUST NOT create things: every variable in `derive`
  MUST be one the derivative binds outside any nested block.
- **DV15.** A field MUST end up with one value. Two derivatives that could
  give the same field different values are rejected before anything runs. If
  one derivative's matches still give a thing's field two values, running
  stops with an error that shows both explanations.
- **DV16.** A derivative MUST NOT have `require` or `message`. Derivatives
  reject nothing.

## How derivatives run

1. **Derivatives run over the stated facts,** each adding what it derives for
   every match, then again over the result, until a pass adds nothing. They
   never create things, merge things or remove facts, so this always ends.
2. **`none`, `count` and `is not set` over derived facts wait.** A derivative
   that uses them over something another derivative adds runs only after that
   one has finished.
3. **Closed world.** A link that is neither stated nor derived does not
   exist, and a field with no value has none.

- **DV17.** Derivatives MUST NOT depend on each other in a circle through
  `none`, `count` or `is not set`. Such derivatives could never settle, and
  are rejected. (In Datalog terms, they must be stratified.)

[Constraints](constraint.md) run after every derivative has finished.

## Explanations

Every derived fact can be explained: the derivative that added it, what each
variable stood for, and each fact the match used. A fact that is itself
derived unfolds into its own explanation, down to facts written in a file.
