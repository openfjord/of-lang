---
title: Patterns — A modeling language for software people
---

# Patterns

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

Counts and allowed values sit on kinds. Anything that looks across more than
one thing is a [constraint](constraint.md) (it rejects) or a [rule](rule.md)
(it derives). Both are written the same way: match a pattern of links, narrow
it with conditions, then require something of every match or derive new facts
from it. A [question](test.md) in a test is a pattern too, with the answer it
expects.

Four principles shape how they are written:

1. **Readability comes before brevity.** Every part is spelled out in YAML
   keys and words. There are no expression strings to parse.
2. **No implicit variables.** Every thing a statement talks about has a name.
   Nothing is reached through a dotted path that skips over an unnamed thing.
3. **Variables start with `$`,** so a variable is never mistaken for a kind, a
   link, a field or a value.
4. **The sentence is plain words.** `says` is what product reads, so it
   names no variables; a `$` word in it is an error. Variables belong to the
   logic.

## The shape of a statement

| Key | Holds | In |
|---|---|---|
| `says` | the sentence, in plain words | both |
| `match` | the pattern of links to find, each end with its kind | both (may be left out when there is one variable) |
| `filter` | conditions that narrow the matches | both, optional |
| `require` | conditions every match must satisfy; without it, every match is a violation | constraints, optional |
| `message` | the broken rule, worded with "must" or "must not" | constraints |
| `derive` | facts to add for every match | rules |

Worked examples are in [constraints](constraint.md) and [rules](rule.md).

## Variables

- **Not in `says`.** The sentence is plain words, with no variables in it.
  A variable is a `$` word used in the logic.
- **Named** `$` then a lowercase letter, then letters and digits: `$job`,
  `$writerApp`.
- **Each variable stands for one thing** in a match. Two variables may match
  the same thing unless a condition says otherwise.
- **Every mention gives its kind.** In a pattern, `from` and `to` are
  written `{$job: Job}`: the variable and the kind it must be. Every mention
  says it, in nested blocks too, and all mentions of a variable must give the
  same kind. Conditions, derivations and messages write the variable bare.
- **The kind must fit the link.** It is the kind the link joins, or one that
  extends it: `{$table: CriticalTable}` on `readsTable` narrows the matches to
  critical tables. A kind the link cannot join is an error.
- **Scope.** A variable that appears only inside a nested block (`at least
  one`, `none`, `count`) belongs to that block and is not visible outside it.

## Links in a pattern

Each entry in `match` is one link that must exist. The link must be declared
in the package or exported by a dependency, and `from` and `to` each give a
variable and its kind. An inverse may be used in either direction:
`link: writtenByJob, from: {$table: Table}, to: {$writer: Job}` matches
exactly what `link: writesTable, from: {$writer: Job}, to: {$table: Table}`
does.

`repeated: true` matches a link one or more times in a row. The things in
between are not named: they are the steps of the chain, and explanations list
them.

```yaml
- link: feedsTable
  from: {$source: Table}
  to:   {$target: Table}
  repeated: true
```

## Conditions

Conditions appear in `filter` and `require`. All entries in a list must hold.

A condition is about a field (`node` and `field`) or about a thing itself
(`node` alone), followed by a comparison word and then `value:` (a literal) or
`same as:` (a variable, or a variable and a field):

```yaml
- node:  $table
  field: retentionDays
  is at most:
    same as:
      node:  $database
      field: retentionDays
```

| Word | Meaning |
|---|---|
| `equals`, `does not equal` | equal, not equal |
| `is above`, `is below` | greater, smaller |
| `is at least`, `is at most` | greater or equal, smaller or equal |
| `is one of` | equal to one of a list (`values:` instead of `value:`) |
| `is a` | the thing is of this kind (takes a kind name directly) |
| `is set`, `is not set` | the field has a value, or has none |

A field with no value fails every comparison except `is not set`.

Patterns can nest inside conditions:

| Word | Meaning |
|---|---|
| `at least one` | the nested `match` (and its own `filter`) finds at least one |
| `none` | they never match |
| `count` | how many times they match, followed by a comparison word |
| `any of` | at least one of the listed conditions holds |

There is no general `not`: `does not equal`, `is not set` and `none` cover the
negative cases while keeping each condition readable on its own.

## How they run

1. **Rules run over the stated facts,** each adding what it derives for every
   match, then again over the result, until a pass adds nothing. Rules never
   create things, merge things or remove facts, so this always ends.
2. **`none`, `count`, `does not equal` and `is not set` over derived facts wait.** A
   rule that uses them over something another rule derives runs only after
   that rule has finished. Rules that depend on each other through these words
   in a circle are rejected. (In Datalog terms, the rules must be stratified.)
3. **Constraints run last,** over stated and derived facts alike, and add
   nothing.
4. **Derived facts live in their own graph,** each recording the rule and the
   match that produced it.

Closed world: a link that is neither stated nor derived does not exist, and a
field with no value has none.

## Explanations

An explanation is the statement's match, filled in: each `match` entry becomes
a fact line and each condition a checked line. A line that is itself derived
unfolds into its own explanation, down to stated facts. A failed `at least
one` also lists what was searched.

## What is checked before anything runs

- A `$` word in `says` or `asks`.
- A link, kind or field that is not declared.
- A `from` or `to` whose kind the link cannot join, or a variable whose
  mentions give different kinds.
- A comparison between values of different types, or with a value a field does
  not allow.
- A derived field that two rules could set differently.
- Rules that depend on each other in a circle through negative words or
  `count`.
- A message quoting a variable that belongs to a nested block.
