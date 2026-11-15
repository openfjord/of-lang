---
title: Constraints — The OpenFjord modeling language
---

# Constraints

> Part of [the reference](README.md).

A constraint rejects facts. It finds things and the links between them,
narrows what it found with conditions, and reports every match as a
violation, unless the match satisfies the conditions under `unless`.

## With `unless`

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
  unless:
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

Read it as: a job reading a table that a job of another application writes
is a violation, unless the table is public and a contract the reading
application consumes covers it.

## Without `unless`

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

Four principles shape how constraints are written:

1. **Readability comes before brevity.** Every part is spelled out in YAML
   keys and words. There are no expression strings to parse.
2. **No implicit variables.** Every thing a constraint talks about has a name
   and a kind. Nothing is reached through a dotted path that skips over an
   unnamed thing.
3. **Variables start with `$`,** so a variable is never mistaken for a kind, a
   link, a field or a value.
4. **The sentence is plain words.** `says` is what product reads, so it names
   no variables. Variables belong to the logic.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Constraint` |
| `says` | yes | the rule, in plain words |
| `match` | yes, unless [CN11](#rules) allows leaving it out | the links to find |
| `filter` | no | conditions that narrow the matches |
| `unless` | no | conditions that excuse a match |
| `message` | yes | what a violation says |

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
stand for the same thing unless a condition says otherwise; one variable used
twice, as in `noFeedbackLoop`, is the same thing both times.

A link's inverse name reads it the other way: `link: writtenByJob, from:
{$table: Table}, to: {$writer: Job}` matches exactly what `link: writesTable,
from: {$writer: Job}, to: {$table: Table}` does.

`repeated: true` matches the link one or more times in a row. The things in
between are not named, so a longer chain between the same two things is
still one match; explanations list the steps.

A narrower kind narrows the matches: `{$table: CriticalTable}` on
`readsTable` matches only jobs reading critical tables.

## Filter and unless

`filter` and `unless` both hold conditions and receive the matches one at a
time. `filter` decides which matches the constraint is about; `unless` says
what excuses one of them. Every condition in a list must hold: a match is
excused only when all of `unless` holds.

A condition is about a thing itself (`node` alone) or one of its fields
(`node` and `field`), followed by a comparison word and what to compare with:
`value:` (a literal), `values:` (a list, for `is one of`) or `same as:` (a
variable, or a variable and a field).

```yaml
unless:
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
match of the constraint:

| Word | Holds when |
|---|---|
| `at least one` | the nested block matches at least once |
| `none` | the nested block never matches |
| `count` | the number of its matches compares as the comparison word beside it says |
| `any of` | at least one of the listed conditions holds |

There is no general `not`: `does not equal`, `is not set` and `none` cover
the negative cases while keeping each condition readable on its own.

## Messages

A message states the broken rule, worded with "must" or "must not". It quotes
variables in braces, filled in from the violating match: `{$job}` gives the
thing's id, `{$table: visibility}` a field's value.

## Rules

- **CN1.** A constraint MUST have a `says` sentence, and the sentence MUST
  NOT contain a `$` word.
- **CN2.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file defines or imports.
- **CN3.** `from` and `to` MUST each be written `{$variable: Kind}`, every
  time, in nested blocks too, and every mention of a variable MUST give the
  same kind.
- **CN4.** The kind MUST be the kind the link joins at that end, or a kind
  that extends it.
- **CN5.** Conditions and messages MUST write a variable bare (`$table`), and
  every variable they use MUST appear in a `match` item.
- **CN6.** A variable that appears only inside a nested block belongs to that
  block and MUST NOT be used outside it, in the message included.
- **CN7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **CN8.** A comparison MUST compare values of the same type: two fields
  that name the same value type, or two fields that name none and have the
  same built-in type. A literal MUST be a value the field allows.
- **CN9.** A field named in a condition, beside `node` or inside `same as`,
  or quoted in a message, MUST be a field of that variable's kind or of a
  kind it extends.
- **CN10.** A condition comparing two nodes without fields MUST compare
  variables whose kinds are the same, or where one kind extends the other.
- **CN11.** `match` MAY be left out when every variable the constraint uses
  appears in a nested block; the constraint then ranges over every thing of
  each such variable's kind.
- **CN12.** A constraint MUST have a `message`, which SHOULD be worded with
  "must" or "must not".
- **CN13.** A constraint MUST NOT have `derive`. Constraints add nothing.

## How constraints run

Constraints run after every [derivative](derivative.md) has finished, over
stated and derived facts alike, and add nothing. A link that is neither
stated nor derived does not exist, and a field with no value has none.

- **CN14.** Each violating match is reported once, with its message and its
  explanation: which conditions under `unless` failed, if any, and the facts
  the match used.
