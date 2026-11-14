---
title: 2. Facts and a test — Tour
---

# 2. Facts and a test

> Part of [the tour](README.md). Previous: [1. A kind](1-a-kind.md).
> Next: [3. Links](3-links.md).

Definitions say what can exist. Facts say what does. In a package, facts
appear only in tests, as examples that prove the definitions do what you
mean.

## A fact

A test is a file in `tests/`. Write `tests/examples/a-team.yaml`:

```yaml
import:
  - from: file:../../defs/team.yaml
    use: [Team]

aTeam:
  is: Test
  says: A team with a valid name is accepted.
  expect:
    rejected: []

team/commerce:
  is: Fact
  kind: Team
  name: commerce
```

`team/commerce` is the fact's id: any text, and a path reads best. `kind`
says what it is, and the other keys are its fields. A fact must have every
field its kind requires and nothing its kind does not list.

The `Test` statement says what the facts must lead to. `rejected: []` means no
rule may be broken. There is one `Test` per file.

## Run the tests

```
$ of-lang test data-model
pass  test aTeam  (data-model/tests/examples/a-team.yaml)
1 passed, 0 failed
```

Give the team a name its value type does not allow:

```yaml
team/commerce:
  is: Fact
  kind: Team
  name: Commerce
```

```
$ of-lang test data-model
FAIL  test aTeam  (data-model/tests/examples/a-team.yaml)
      team/commerce: name Commerce does not match ^[a-z][a-z0-9_]*$ (FC4)
0 passed, 1 failed
```

A test passes only if every fact obeys its kind and the facts lead to exactly
what `expect` says. Tests run in CI (continuous integration) on every change,
like any other tests.
