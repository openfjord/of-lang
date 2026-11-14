---
title: Tests — The OpenFjord modeling language
---

# Tests

> Part of [the reference](README.md).

A package proves itself with tests in its `tests/` folder (reason 16 in
[the research](../../README.md)). There are two kinds: an **example** says
what a set of facts must derive and which constraints it must break; a
**question** asks a set of facts something and says what the answer must be.

## Examples

An example is a file of facts with one `Test` statement:

```yaml
# tests/examples/table-inherits-database-team.of.yaml
import:
  - from: file:../../defs/team.of.yaml
    use: [Team]
  - from: file:../../defs/database.of.yaml
    use: [Database, ownedByTeam]
  - from: file:../../defs/table.of.yaml
    use: [Table, inDatabase, belongsToTeam]

tableInheritsDatabaseTeam:
  is: Test
  says: A table belongs to the team that owns its database.
  expect:
    derived:
      - link: belongsToTeam
        from: db/sales/table/order_lines
        to:   team/commerce
    rejected: []

team/commerce:
  is: Fact
  kind: Team
  name: commerce

db/sales:
  is: Fact
  kind: Database
  name: sales
  ownedByTeam: team/commerce

db/sales/table/order_lines:
  is: Fact
  kind: Table
  name: order_lines
  inDatabase: db/sales
```

| Key under `expect` | Value |
|---|---|
| `derived` | facts that must hold after derivatives run |
| `notDerived` | facts that must not hold |
| `rejected` | every constraint the facts must break |

A fact under `derived` or `notDerived` is written like a derivation, with ids
in place of variables:

| Fact | Written |
|---|---|
| a link | `link`, `from`, `to` |
| a kind | `node`, `becomes` |
| a field value | `node`, `field`, `set to` with `value:` |

- **TS1.** An example file MUST hold exactly one `Test` statement, which MUST
  have `says` and `expect`.
- **TS2.** The test passes only if every fact under `derived` holds, no fact
  under `notDerived` holds, the facts break exactly the constraints under
  `rejected`, and every fact obeys its kind ([FC3](fact.md#rules)). Leaving
  `rejected` out means none.
- **TS3.** Every link, kind and constraint the test names MUST be imported.

## Questions

A question asks the facts of an example something, and says what the answer
must be. It finds things and the links between them, narrows what it found
with conditions, and compares the things a variable stood for with the
expected list.

```yaml
# tests/questions/who-breaks-if-order-intake-goes.of.yaml
import:
  - from: file:../../defs/application.of.yaml
    use: [Application, dependsOnApp]

whoBreaksIfOrderIntakeGoes:
  is: Question
  asks: Which applications depend on order intake?
  given: file:../examples/two-apps-sharing-a-table.of.yaml
  match:
    - link: dependsOnApp
      from: {$app: Application}
      to:   {$target: Application}
  filter:
    - node: $target
      equals:
        value: app/order_intake
  expect:
    $app: [app/revenue_report]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Question` |
| `asks` | yes | the question, in plain words |
| `given` | no | `file:` and an example file whose facts to ask; without it, the facts in this file |
| `match` | yes | the links to find |
| `filter` | no | conditions that narrow the matches |
| `expect` | yes | for each variable listed, the ids it must stand for |

Each item in `match` is one link that must exist, with `link`, `from` and
`to`; `from` and `to` each give a variable and its kind, as
`{$variable: Kind}`. `match` finds every combination of things for which all
its links hold; each is one match, giving every variable one thing.
`repeated: true` on an item follows the link one or more steps.

`filter` keeps the matches for which every condition holds. A condition
names a `node`, optionally a `field`, a comparison word, and what to compare
with:

| Word | Takes |
|---|---|
| `equals`, `does not equal`, `is above`, `is below`, `is at least`, `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set`, `is not set` | nothing more |

A condition may also nest a block with its own `match` and `filter`:
`at least one`, `none`, or `count` with a comparison word; `any of` holds
when one of its listed conditions does.

- **TS4.** A question MUST have an `asks` sentence, and the sentence MUST NOT
  contain a `$` word.
- **TS5.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file imports.
- **TS6.** `from` and `to` MUST each be written `{$variable: Kind}`, every
  time, in nested blocks too; every mention of a variable MUST give the same
  kind, and the kind MUST be the link's kind at that end or one extending it.
- **TS7.** Conditions MUST write a variable bare (`$target`), and every
  variable they use MUST appear in a `match` item. A variable that appears
  only inside a nested block MUST NOT be used outside it.
- **TS8.** A comparison MUST compare values of the same type (two fields
  that name the same value type, or two fields that name none and have the
  same built-in type); a literal MUST be a value the field allows; a field
  MUST be one of its node's kind; two nodes compared without fields MUST
  have the same kind, or one extending the other.
- **TS9.** Each key under `expect` MUST be a variable of the question's
  `match`.
- **TS10.** The question passes only if, for each variable listed, the
  distinct things it stands for across all matches are exactly the ids
  listed.
