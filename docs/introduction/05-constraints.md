---
title: Constraints — Introduction to the OpenFjord modeling language
---

# Constraints

> Part of [the introduction](README.md).

In the past, a guest of the restaurant became ill. The pesto in a dish
contained nuts, and the menu did not tell the guest. The package can
conclude that a dish contains nuts. But it accepts a menu that does not give
the warning. The kitchen needs a rule that rejects such a menu. A constraint
is that rule: it finds things by their links and conditions, and rejects
each match that breaks it.

## A constraint with an exception

```yaml
# defs/allergens.of.yaml
dishListsItsAllergens:
  is: Constraint
  says: >
    The menu must give a warning about each allergen in an active dish.
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Allergen, as: $allergen }
    - link: hasAllergen
      from: $dish
      to:   $allergen
  filter:
    - node:  $dish
      field: status
      equals:
        value: active
  unless:
    - at least one:
        match:
          - link: listsAllergen
            from: $dish
            to:   $allergen
  message: "The menu must warn that {$dish} contains {$allergen}."
```

`says`, `match` and `filter` are the same as in a derivative. The match
finds each dish with each allergen that it contains. The filter keeps the
dishes with the status `active`: a condition can name a `field` of a thing
and compare its value.

A constraint has no `derive`. It has two keys of its own:

- `unless` lists conditions. A match is a violation unless all of them hold.
  Here, the match is correct when the menu lists the allergen for the dish.
  `at least one` holds when its nested block matches one time or more.
- `message` is the text that the tool shows for each violation. It can name
  variables.

The match uses `hasAllergen`, a link that only derivatives add. All the
derivatives finish first, and then the constraints check the stated facts
and the derived facts together.

## A test that expects a rejection

The kitchen writes down the case of the guest who became ill. The facts are
a margherita that uses the pesto, and no warning on the menu. The test
expects that the facts break the constraint:

```yaml
# tests/examples/dish-without-allergen-warning.of.yaml
dishWithoutAllergenWarning:
  is: Test
  says: >
    The margherita uses the pesto, and the pesto contains nuts. The menu does
    not give a warning about nuts, so the constraint rejects the dish.
  expect:
    rejected: [dishListsItsAllergens]
```

`rejected` lists all the constraints that the facts must break. The test
fails if the facts break a different constraint, or none. The file imports
the constraint by its name, as it imports kinds and links.

`run` shows the message of each rejection:

```console
$ of-lang run step-05 step-05/tests/examples/dish-without-allergen-warning.of.yaml
derived:
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  station/pizza depends on station/sauce  (derivative stationDependsOnStation)
rejected:
  dishListsItsAllergens: The menu must warn that dish/margherita contains allergen/nuts.
```

The example from chapter 4 has `rejected: []`. In that example, the menu
lists nuts and milk for the margherita, so the facts break no constraint.

## A constraint that counts

Each recipe has one station. A link alone does not tell this: a thing can
have a link any number of times, or not at all. A constraint sets the
number:

```yaml
# defs/recipe.of.yaml
recipeAtOneStation:
  is: Constraint
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
```

This constraint has no `unless`. Each match that passes the filter is a
violation: a recipe without a station, or a recipe with two stations.

## A constraint that compares fields

The kitchen has a second topic: food spoils when cooks keep it too warm. The
kitchen starts a file for it, `defs/spoilage.of.yaml`, with statements of
the sorts that the chapters before this one explain:

- two value types: `Celsius`, a temperature, and `Day`, a day in the
  calendar
- two kinds: `StorageZone`, such as the fridge, with a field `maxCelsius`,
  and `Batch`, an amount of one ingredient that arrived on one day, with the
  fields `receivedOn`, `useBy` and `keepAtMostCelsius`
- two links: `batchOf` from a batch to its ingredient, and `keptIn` from a
  batch to its zone.

[The file in `steps/step-05`](steps/step-05/defs/spoilage.of.yaml) shows them. The
rule of the topic is a constraint:

```yaml
# defs/spoilage.of.yaml
batchKeptColdEnough:
  is: Constraint
  says: The cooks keep each batch in a zone that is cold enough for it.
  match:
    - kind: { name: Batch, as: $batch }
    - kind: { name: StorageZone, as: $zone }
    - link: keptIn
      from: $batch
      to:   $zone
  filter:
    - node:  $zone
      field: maxCelsius
      is above:
        same as:
          node:  $batch
          field: keepAtMostCelsius
  message: "{$batch} must be kept at {$batch: keepAtMostCelsius} °C or colder, but {$zone} is at {$zone: maxCelsius} °C."
```

`same as` compares a field with a field of a different thing, where `value`
compares it with a literal. The two fields have the same value type,
`Celsius`. In a message, `{$batch: keepAtMostCelsius}` shows the value of a
field.

An example states two batches of parmesan, one in the fridge and one in the
dry store:

```console
$ of-lang run step-05 step-05/tests/examples/batch-kept-too-warm.of.yaml
derived:
  nothing
rejected:
  batchKeptColdEnough: batch/parmesan_tuesday must be kept at 8 °C or colder, but zone/dry_store is at 20 °C.
```

A second constraint in the file, `useByAfterArrival`, compares two fields of
one batch: the use-by day must not be before the day that the batch arrived.

All the tests pass. The two tests that expect a rejection pass because the
rejection occurs:

```console
$ of-lang test step-05
pass  test allergenThroughASauce  (step-05/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (step-05/tests/examples/base-recipe.of.yaml)
pass  test batchKeptTooWarm  (step-05/tests/examples/batch-kept-too-warm.of.yaml)
pass  test dishWithoutAllergenWarning  (step-05/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-05/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-05/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (step-05/tests/terms/garnish.of.yaml)
7 passed, 0 failed
```

Derivatives derive, and constraints reject. No statement does both.

The model of the kitchen is complete. But the kitchen is not alone: a
different team owns the suppliers, and an app for guests wants to read the
allergens. The next chapter is about the package and the world around it:
[Packages](06-packages.md).

More about this: [constraints](../../language/reference/constraint.md) and
[tests](../../language/reference/test.md#examples) in the reference. In the
language by example: [Constraints](../by-example/constraints.md),
[Unless](../by-example/unless.md),
[At least one](../by-example/at-least-one.md),
[Messages](../by-example/messages.md), [Same as](../by-example/same-as.md) and
[Expect rejected facts](../by-example/expect-rejected.md).
