---
title: Patterns — The OpenFjord modeling language
---

# Patterns

> Part of [the reference](README.md).

Counts and allowed values sit on [kinds](kind.md). Anything that looks across
more than one thing is a [constraint](constraint.md), which rejects, or a
[rule](rule.md), which derives. Both are written the same way: find a pattern
of things and the links between them, narrow it with conditions, then require
something of every match or derive new facts from it. A
[question](test.md#questions) in a test is a pattern too, with the answer it
expects.

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

Four principles shape how patterns are written:

1. **Readability comes before brevity.** Every part is spelled out in YAML
   keys and words. There are no expression strings to parse.
2. **No implicit variables.** Every thing a pattern talks about has a name
   and a kind. Nothing is reached through a dotted path that skips over an
   unnamed thing.
3. **Variables start with `$`,** so a variable is never mistaken for a kind, a
   link, a field or a value.
4. **The sentence is plain words.** `says` is what product reads, so it names
   no variables. Variables belong to the logic.

## The keys of a statement

| Key | Holds | In |
|---|---|---|
| `says` | the sentence, in plain words (`asks` in a question) | all |
| `match` | the links to find, each end with its variable and kind | all; may be left out, see [PT9](#rules) |
| `filter` | conditions that narrow the matches | all, optional |
| `require` | conditions every match must satisfy | [constraints](constraint.md), optional |
| `message` | the broken rule, worded with "must" or "must not" | [constraints](constraint.md) |
| `derive` | facts to add for every match | [rules](rule.md) |

## Match

Each entry in `match` is one link that must exist between two things:

```yaml
match:
  - link: readsTable
    from: {$job: Job}
    to:   {$table: Table}
```

`from` and `to` each give a variable and the kind it must be. A variable is
`$` followed by a lowercase letter, then letters and digits: `$job`,
`$writerApp`. Each variable stands for one thing in a match; two variables
may stand for the same thing unless a condition says otherwise.

A link's inverse name reads it the other way: `link: writtenByJob, from:
{$table: Table}, to: {$writer: Job}` matches exactly what `link: writesTable,
from: {$writer: Job}, to: {$table: Table}` does.

`repeated: true` matches the link one or more times in a row. The things in
between are not named; explanations list them as the steps of the chain.

```yaml
- link: feedsTable
  from: {$source: Table}
  to:   {$target: Table}
  repeated: true
```

A narrower kind narrows the matches: `{$table: CriticalTable}` on
`readsTable` matches only jobs reading critical tables.

## Conditions

Conditions appear in `filter` and `require`. Every condition in a list must
hold.

A condition is about a thing itself (`node` alone) or one of its fields
(`node` and `field`), followed by a comparison word and what to compare with:
`value:` (a literal), `values:` (a list, for `is one of`) or `same as:` (a
variable, or a variable and a field).

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

Patterns nest inside conditions. A nested block has its own `match` and
optional `filter`:

| Word | Holds when |
|---|---|
| `at least one` | the nested pattern matches at least once |
| `none` | the nested pattern never matches |
| `count` | the number of matches compares as the comparison word beside it says |
| `any of` | at least one of the listed conditions holds |

```yaml
filter:
  - count:
      match:
        - link: readsTable
          from: {$reader: Job}
          to:   {$table: Table}
    is above:
      value: 2
```

There is no general `not`: `does not equal`, `is not set` and `none` cover
the negative cases while keeping each condition readable on its own.

## Rules

- **PT1.** A statement MUST have a `says` sentence (`asks` in a question),
  and the sentence MUST NOT contain a `$` word.
- **PT2.** Every entry in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file defines or imports.
- **PT3.** `from` and `to` MUST each be written `{$variable: Kind}`, every
  time, in nested blocks too, and every mention of a variable MUST give the
  same kind.
- **PT4.** The kind MUST be the kind the link joins at that end, or a kind
  that extends it.
- **PT5.** Conditions, derivations and messages MUST write a variable bare
  (`$table`), and every variable they use MUST appear in a `match` entry.
- **PT6.** A variable that appears only inside a nested block belongs to that
  block and MUST NOT be used outside it.
- **PT7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **PT8.** A comparison MUST compare values of the same type, and a literal
  MUST be a value the field allows.
- **PT9.** `match` MAY be left out when every variable the statement uses
  appears in a nested block; the statement then ranges over every thing of
  each such variable's kind.

## How they run

1. **Rules run over the stated facts,** each adding what it derives for
   every match, then again over the result, until a pass adds nothing. Rules
   never create things, merge things or remove facts, so this always ends.
2. **`none`, `count` and `is not set` over derived facts wait.** A rule that
   uses them over something another rule derives runs only after that rule
   has finished.
3. **Constraints run last,** over stated and derived facts alike, and add
   nothing.
4. **Closed world.** A link that is neither stated nor derived does not
   exist, and a field with no value has none.

- **PT10.** Rules MUST NOT depend on each other in a circle through `none`,
  `count` or `is not set`. Such rules could never settle, and are rejected.
  (In Datalog terms, the rules must be stratified.)

## Explanations

Every derived fact and every rejection can be explained. An explanation is
the statement's match, filled in: the rule or constraint, what each variable
stood for, and each fact the match used. A fact that is itself derived
unfolds into its own explanation, down to facts written in a file.
