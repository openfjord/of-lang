---
title: Tests — A modeling language for software people
---

# Tests

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A package proves itself with two kinds of test (reason 16).

**An example** is a file of facts with one `Test` entry saying what they
must derive and which constraints they must break:

```yaml
# tests/examples/table-inherits-database-team.yaml
import:
  - from: file:../../defs/table.yaml
    use: [Table, inDatabase, belongsToTeam]
  - from: file:../../defs/database.yaml
    use: [Database, ownedByTeam, Team]

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
  …
```

- `derived` lists links, kinds (`node`, `becomes`) or field values (`node`,
  `field`, `set to`) that must be derived. `notDerived` lists ones that must
  not be.
- `rejected` lists every constraint the facts must break; breaking any other
  fails the test. Leaving it out means none.

**A question** asks the facts of an example something, written as a pattern,
with the answer it expects for each variable:

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
      to:   {$target: Table}
  filter:
    - node: $target
      equals:
        value: app/order-intake
  expect:
    $app: [app/revenue-report]
```

`asks`, like `says`, is plain words. Each key of `expect` is a variable of
the pattern, and the answer must be exactly the listed things.
