---
title: Constraints — The OpenFjord modeling language
---

# Constraints

> Part of [the reference](README.md).

A constraint rejects facts. It finds things and the links between them, and
it uses conditions to narrow the matches. It reports each match as a
violation. A match that satisfies the conditions under `unless` is not a
violation.

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

Read it as follows. A job reads a table that a job of a different
application writes. This is a violation, unless two things are true: the
table is public, and the application of the first job consumes a contract
that covers the table.

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

Each match is a violation. The link `feedsTable` chains, so this constraint
also finds a table that feeds itself through other tables.

Four principles apply to constraints:

1. **Readability comes before brevity.** YAML keys and words state each part
   in full. There are no expression strings to parse.
2. **No implicit variables.** Each thing that a constraint refers to has a
   name and a kind. A constraint does not reach a thing through a dotted path
   that skips an unnamed thing.
3. **Variables start with `$`,** so a reader cannot confuse a variable with a
   kind, a link, a field or a value.
4. **The sentence is plain words.** Product reads `says`, so `says` names no
   variables. Variables are part of the logic.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Constraint` |
| `says` | yes | the rule, in plain words |
| `match` | yes, unless [CN11](#rules) makes it optional | the links to find |
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

`from` and `to` each give a variable and the kind that the variable must
have. A variable is `$`, then a lowercase letter, then letters and digits:
`$job`, `$writerApp`.

`match` finds each combination of things for which all its links hold. Each
combination is one match: a row that gives one thing to each variable. Two
matches that give the same things to the same variables are one match. Two
variables can represent the same thing, unless a condition prevents it. One
variable that occurs two times, as in `noFeedbackLoop`, is the same thing at
each occurrence.

The inverse name of a link reads the link in the opposite direction: `link: writtenByJob, from:
{$table: Table}, to: {$writer: Job}` matches exactly the same things as `link: writesTable,
from: {$writer: Job}, to: {$table: Table}`.

`repeated: true` matches the link one or more times in a sequence. The things
between the two ends have no names. Thus a longer chain between the same two
things is still one match. Explanations list the steps.

A narrower kind limits the matches. For example, `{$table: CriticalTable}`
on `readsTable` matches only jobs that read critical tables.

## Filter and unless

`filter` and `unless` both contain conditions, and both receive the matches
one at a time. `filter` selects the matches that the constraint applies to.
`unless` tells which of these matches to excuse. Each condition in a list must
hold. Thus `unless` excuses a match only when all its conditions hold.

A condition applies to a thing (`node` alone) or to one of its fields (`node`
and `field`). Then it gives a comparison word and what to compare with:

- `value:` (a literal)
- `values:` (a list, for `is one of`)
- `same as:` (a variable, or a variable and a field).

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
| `equals`, `does not equal` | the two are equal, or are not equal |
| `is above`, `is below` | the first is greater, or smaller |
| `is at least`, `is at most` | the first is greater or equal, or smaller or equal |
| `is one of` | the first equals one of `values` |
| `is a` | the thing has the named kind (`is a: CriticalTable`) |
| `is set`, `is not set` | the field has a value, or has no value |

When a condition compares a thing without a field, it compares the id of the
thing. A field with no value fails all comparisons except `is not set`.

A condition cannot see links directly. It reaches links through a nested
block. A nested block has its own `match` and an optional `filter`. It runs
one time for each match of the constraint:

| Word | Holds when |
|---|---|
| `at least one` | the nested block has one or more matches |
| `none` | the nested block has no matches |
| `count` | the number of its matches satisfies the comparison word next to it |
| `any of` | one or more of the listed conditions hold |

There is no general `not`. `does not equal`, `is not set` and `none` cover the
negative cases, and each condition stays easy to read alone.

## Messages

A message states the rule that the match breaks. It uses the words "must" or
"must not". It quotes variables in braces, and the match gives their values.
`{$job}` gives the id of the thing. `{$table: visibility}` gives the value of
a field.

## Rules

- **CN1.** A constraint MUST have a `says` sentence, and the sentence MUST
  NOT contain a `$` word.
- **CN2.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file defines or imports.
- **CN3.** Each `from` and each `to` MUST use the form `{$variable: Kind}`,
  every time and also in nested blocks. Every mention of a variable MUST give
  the same kind.
- **CN4.** The kind MUST be the kind that the link joins at that end, or a
  kind that extends it.
- **CN5.** Conditions and messages MUST write a variable without its kind
  (`$table`), and every variable that they use MUST appear in a `match` item.
- **CN6.** A variable that appears only in a nested block belongs to that
  block. The constraint MUST NOT use it outside that block, and this includes
  the message.
- **CN7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **CN8.** A comparison MUST compare values of the same type: two fields
  that name the same value type, or two fields that name no value type and
  have the same built-in type. A literal MUST be a value that the field
  allows.
- **CN9.** A field can occur in a condition, next to `node` or in `same as`.
  A field can also occur in a message. Such a field MUST be a field of that
  variable's kind, or of a kind that it extends.
- **CN10.** When a condition compares two nodes without fields, the two
  variables MUST have the same kind, or kinds where one extends the other.
- **CN11.** A constraint MAY have no `match` when each variable that it uses
  appears in a nested block. Then each such variable represents, in turn,
  every thing of its kind.
- **CN12.** A constraint MUST have a `message`. The message SHOULD use the
  words "must" or "must not".
- **CN13.** A constraint MUST NOT have `derive`. Constraints add no facts.

## How constraints run

Constraints run after each [derivative](derivative.md) finishes. They run over
stated facts and derived facts, and they add no facts. A link that is not
stated and not derived does not exist. A field with no value is not set, not unknown.

- **CN14.** The tool reports each match that is a violation one time, with
  its message and its explanation. The explanation gives the conditions under
  `unless` that failed, if any, and the facts that the match used.
