---
title: Pending changes — Introduction to the OpenFjord modeling language
---

# Pending changes

> Part of [the introduction](README.md).

The package is finished, but a kitchen changes. The restaurant sells more
pizzas, and the pastry station starts to make pizza dough too. The package
rejects this: `recipeAtOneStation` allows exactly one station for each
recipe. A cook proposes to change the rule. The owner of the package must
approve the change first, and wants to see what it does.

A pending change is a change that nobody has accepted yet. It is in the
files of the package, as marks on the statements that it adds, removes or
updates.

## The marks

The cook writes the change in `defs/recipe.of.yaml`:

```yaml
# defs/recipe.of.yaml
recipeAtOneStation:          # the change removes this statement
  is: Constraint
  pending: removed
  says: Cooks make each recipe at exactly one station.
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - count:
        match:
          - link: madeAtStation
            from: $recipe
            to:   Station
      does not equal:
        value: 1
  message: "{$recipe} must be made at exactly one station."

recipeAtAStation:            # the change adds this statement
  is: Constraint
  pending: added
  says: Cooks make each recipe at one station or more.
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - count:
        match:
          - link: madeAtStation
            from: $recipe
            to:   Station
      is below:
        value: 1
  message: "{$recipe} must be made at a station."

madeAtStation:               # the change updates this statement
  is: Link
  pending: updated
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: A station where cooks make the recipe.
  inverse: makesRecipe
  was:
    from: Recipe
    to: Station
    reads: "{from} is made at {to}"
    definition: The station where cooks make the recipe.
    inverse: makesRecipe
```

`pending` can be on each statement, next to the keys of the statement
itself. It has one of three values:

- `removed`: the statement stays in the file as it is, with the mark.
- `added`: the statement is in the file, with the mark.
- `updated`: the statement shows its new version. `was` keeps the whole old
  version, without `is`. Here, the update changes one word of the
  definition, from "the station" to "a station".

Facts and tests can have the mark too. The cook adds an example that proves
the change, in a file of its own:

```yaml
# tests/examples/dough-at-two-stations.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Station]
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, madeAtStation]

doughAtTwoStations:
  is: Test
  pending: added
  says: >
    The pizza station and the pastry station each make the pizza dough. The
    facts break no constraint.
  expect:
    rejected: []

station/pizza:
  is: Fact
  pending: added
  kind: Station
  fields:
    name: pizza

station/pastry:
  is: Fact
  pending: added
  kind: Station
  fields:
    name: pastry

recipe/pizza_dough:
  is: Fact
  pending: added
  kind: Recipe
  fields:
    name: pizza_dough
    steps: https://recipes.example.com/pizza_dough
  links:
    madeAtStation: [station/pizza, station/pastry]
```

## Two versions of one package

All the marks in a package together are its pending change. A package has
one pending change or none. With the marks, the files hold two versions of
the package:

| Version | Added statements | Removed statements | Updated statements |
|---|---|---|---|
| without the change | not included | included | the `was` version |
| with the change | included | not included | the new version |

`check` checks the two versions. Each one must obey each rule of the
language:

```console
$ of-lang check --packages ../../../examples step-07
ok
```

`test`, `run` and `explain` use the version with the change. Thus, the
owner sees that the tests of the package still pass, and that the example
of the cook passes too:

```console
$ of-lang test --packages ../../../examples step-07
pass  test allergenThroughASauce  (step-07/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (step-07/tests/examples/base-recipe.of.yaml)
pass  test batchKeptTooWarm  (step-07/tests/examples/batch-kept-too-warm.of.yaml)
pass  test dishWithoutAllergenWarning  (step-07/tests/examples/dish-without-allergen-warning.of.yaml)
pass  test doughAtTwoStations  (step-07/tests/examples/dough-at-two-stations.of.yaml)
pass  question whichDishesContainNuts  (step-07/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-07/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (step-07/tests/terms/garnish.of.yaml)
8 passed, 0 failed
```

The package in [`steps/step-07`](steps/step-07/) is the finished package, without this
change. To get the output above, put the marks from this chapter in a copy
of it.

Other packages see only the version without the change. A pending change
has no effect outside its package until somebody accepts it.

## Accept or reject

The owner decides about the change as a whole.

To accept the change is one operation:

1. Delete each removed statement.
2. Delete `pending` and `was` from each added statement and each updated
   statement.

Then `recipeAtOneStation` is gone, `recipeAtAStation` and the example are
ordinary statements, and `madeAtStation` has its new definition.

To reject the change is also one operation:

1. Delete each added statement.
2. Replace each updated statement with its `was` version.
3. Delete `pending` from each removed statement.

Then the package is again as it is in `steps/step-07`.

In the two cases, no mark is left, and the package has no pending change.
The manifest and the Markdown files in `docs/` have no marks.

## Where to go from here

This is the end of the introduction. The package of the kitchen uses most
of the language, but not all of it.

- [The language by example](../by-example/README.md) has a short page for
  each thing in the language, also for the things that the kitchen does not
  use.
- [The cheat sheet](../../language/cheatsheet.md) shows all of the syntax on
  one page.
- [The tool's page](../../cli/README.md) describes each command and its
  output.

More about this: [pending changes](../../language/reference/pending.md) in the
reference. In the language by example:
[Pending changes](../by-example/pending-changes.md).
