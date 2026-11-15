---
title: Tests — The OpenFjord modeling language
---

# Tests

> Part of [the reference](README.md).

A package uses tests in its `tests/` folder to prove that it is correct
(reason 16 in [the research](../../README.md)). A test is an example or a
question:

- An **example** tells what a set of facts must derive, and which constraints
  the facts must break.
- A **question** asks something about a set of facts, and tells what the
  answer must be.

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
| `derived` | facts that must hold after the derivatives run |
| `notDerived` | facts that must not hold |
| `rejected` | all constraints that the facts must break |

A fact under `derived` or `notDerived` has the same form as a derivation,
with ids in place of variables:

| Fact | Written |
|---|---|
| a link | `link`, `from`, `to` |
| a kind | `node`, `becomes` |
| a field value | `node`, `field`, `set to` with `value:` |

- **TS1.** An example file MUST hold exactly one `Test` statement, which MUST
  have `says` and `expect`.
- **TS2.** The test passes only if all of these are true:
  - every fact under `derived` holds
  - no fact under `notDerived` holds
  - the facts break exactly the constraints under `rejected`
  - every fact obeys its kind ([FC3](fact.md#rules)).

  A test with no `rejected` is the same as a test that expects no broken
  constraints.
- **TS3.** The file MUST import each link, kind and constraint that the test
  names.

## Questions

A question asks something about the facts of an example, and tells what the
answer must be. It finds things and the links between them, and it uses
conditions to narrow the matches. Then it compares the things that a variable
represented with the expected list.

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
| `given` | no | `file:` and an example file with the facts to ask about; if not given, the facts in this file |
| `match` | yes | the links to find |
| `filter` | no | conditions that narrow the matches |
| `expect` | yes | for each listed variable, the ids that it must represent |

Each item in `match` is one link that must exist, with `link`, `from` and
`to`. `from` and `to` each give a variable and its kind, as
`{$variable: Kind}`. `match` finds each combination of things for which all
its links hold. Each combination is one match, which gives one thing to each
variable. `repeated: true` on an item follows the link for one or more steps.

`filter` keeps the matches for which each condition holds. A condition names
a `node`, an optional `field`, a comparison word and what to compare with:

| Word | Takes |
|---|---|
| `equals`, `does not equal`, `is above`, `is below`, `is at least`, `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set`, `is not set` | nothing more |

A condition can also contain a nested block with its own `match` and
`filter`. The words for this are `at least one`, `none`, or `count` with a
comparison word. `any of` holds when one of its listed conditions holds.

- **TS4.** A question MUST have an `asks` sentence, and the sentence MUST NOT
  contain a `$` word.
- **TS5.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file imports.
- **TS6.** Each `from` and each `to` MUST use the form `{$variable: Kind}`,
  every time and also in nested blocks. Every mention of a variable MUST give
  the same kind. The kind MUST be the kind of the link at that end, or a kind
  that extends it.
- **TS7.** Conditions MUST write a variable without its kind (`$target`), and
  every variable that they use MUST appear in a `match` item. A question
  MUST NOT use a variable outside a nested block if the variable appears only
  in that block.
- **TS8.** A comparison MUST compare values of the same type: two fields that
  name the same value type, or two fields that name no value type and have
  the same built-in type. A literal MUST be a value that the field allows. A
  field MUST be a field of the kind of its node. Two nodes that a condition
  compares without fields MUST have the same kind, or kinds where one
  extends the other.
- **TS9.** Each key under `expect` MUST be a variable of the question's
  `match`.
- **TS10.** The question passes only if, for each listed variable, the
  distinct things that it represents in all matches are exactly the listed
  ids.
