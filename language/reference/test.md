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
# tests/examples/table-inherits-database-team.yaml
import:
  - from: file:../../defs/team.yaml
    use: [Team]
  - from: file:../../defs/database.yaml
    use: [Database, ownedByTeam]
  - from: file:../../defs/table.yaml
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
| `derived` | facts that must hold after rules run |
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

A question is a [pattern](pattern.md) with the answer it expects for some of
its variables, asked of the facts in an example file:

```yaml
# tests/questions/who-breaks-if-order-intake-goes.yaml
import:
  - from: file:../../defs/application.yaml
    use: [Application, dependsOnApp]

whoBreaksIfOrderIntakeGoes:
  is: Question
  asks: Which applications depend on order intake?
  given: file:../examples/two-apps-sharing-a-table.yaml
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
| `match`, `filter` | as for any [pattern](pattern.md) | the pattern |
| `expect` | yes | for each variable listed, the ids it must stand for |

- **TS4.** Each key under `expect` MUST be a variable of the pattern.
- **TS5.** The question passes only if, for each variable listed, the
  distinct things it stands for across all matches are exactly the ids
  listed.
