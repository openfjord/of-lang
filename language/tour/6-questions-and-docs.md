---
title: 6. Questions and docs — Tour
---

# 6. Questions and docs

> Part of [the tour](README.md). Previous: [5. Constraints](5-constraints.md).
> Next: [7. Many packages](7-many-packages.md).

## A question

An application depends on another when its jobs read what the other's jobs
write. `application.yaml` has a rule for it, `appDependsOnWriter`, deriving
`dependsOnApp`. Now ask the model a question someone would really ask, and
fix the answer you expect. Write
`tests/questions/who-breaks-if-order-intake-goes.yaml`:

```yaml
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

A question is a pattern, like a rule's, asked of the facts in an example
file. `expect` lists, for a variable, every thing it must stand for, and
nothing else.

```
$ of-lang test data-model
pass  test tableReadByThreeJobs  (data-model/tests/examples/critical-table.yaml)
pass  test crossAppReadWithoutContract  (data-model/tests/examples/cross-app-read-without-contract.yaml)
pass  test twoAppsSharingATable  (data-model/tests/examples/two-apps-sharing-a-table.yaml)
pass  question whoBreaksIfOrderIntakeGoes  (data-model/tests/questions/who-breaks-if-order-intake-goes.yaml)
4 passed, 0 failed
```

Questions like this are how you know the model answers what people need.
Write them before the rules, and let them fail until the rules are right.

## Docs

The model holds what can be checked. Why, how and history go in Markdown, in
`docs/`. A doc links to a definition or embeds it, so prose and definitions
read side by side:

```markdown
# Tables

A table holds rows with one schema. Every table sits in exactly one database,
and belongs to the team that owns that database:

![](../defs/table.yaml#Table)

![](../defs/table.yaml#tableTeamFromDatabase)

Teams kept breaking each other by reading internal tables, so reading across
applications needs a contract — see
[crossAppReadNeedsContract](../defs/reading-across-applications.yaml#crossAppReadNeedsContract).
```

`![](file#Name)` embeds the statement; `[text](file#Name)` links to it. Both are
plain Markdown, so any viewer shows a working link. The check fails if a link
points at a statement that does not exist, so docs cannot drift silently. Each
definition's generated page lists the docs that mention it.

The package now has its full shape:

```
data-model/
  package.yaml
  defs/
  docs/
  tests/
```

Nothing else may sit at the top.
