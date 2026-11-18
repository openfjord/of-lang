---
title: Is the food kept cold enough? — Introduction to the OpenFjord modeling language
---

# Is the food kept cold enough?

> Part of [the introduction](README.md).

Food spoils when cooks keep it too warm or too long. On a busy day, a cook
puts a delivery of parmesan in the dry store, because the fridge is full.
Nobody sees it until the cheese is bad. The kitchen wants the tool to find
a batch in the wrong place on the day that it arrives.

## Batches and storage zones

This is a second topic, so the kitchen starts a file for it,
[`defs/spoilage.of.yaml`](steps/step-06/defs/spoilage.of.yaml). The kitchen
receives each ingredient in batches, and keeps each batch in a storage zone:

```yaml
# defs/spoilage.of.yaml
StorageZone:
  is: Kind
  definition: A place in the kitchen where cooks keep food at one temperature.
  examples:
    - the fridge
    - the dry store
  notExamples:
    - a station, where cooks make food and do not keep it
  fields:
    name:
      type: Name
    maxCelsius:
      type: Celsius

Batch:
  is: Kind
  definition: An amount of one ingredient that the kitchen received on one day.
  examples:
    - the parmesan that arrived on Monday
  notExamples:
    - parmesan in general, which is an ingredient
  fields:
    receivedOn:
      type: Day
    useBy:
      type: Day
    keepAtMostCelsius:
      type: Celsius
```

The file also defines the two value types of these fields. `Celsius` is a
whole number from -40 to 40. `Day` is text with `format: date`, a day in the
calendar such as `2026-11-17`. Two links connect the kinds to the model:
`batchOf` goes from a batch to its ingredient, and `keptIn` goes from a
batch to its zone.

## The rule

A batch is in the wrong place when its zone is warmer than the batch can
stay at. That compares a field of one thing with a field of a different
thing:

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

## A batch in the dry store

The example has two zones and two batches of parmesan. The batch of Monday
is in the fridge. The batch of Tuesday is in the dry store:

```yaml
# tests/examples/batch-kept-too-warm.of.yaml
zone/dry_store:
  is: Fact
  kind: StorageZone
  fields:
    name: dry_store
    maxCelsius: 20

batch/parmesan_tuesday:
  is: Fact
  kind: Batch
  fields:
    receivedOn: 2026-11-17
    useBy: 2026-12-15
    keepAtMostCelsius: 8
  links:
    batchOf: ingredient/parmesan
    keptIn: zone/dry_store
```

```console
$ of-lang run step-06 step-06/tests/examples/batch-kept-too-warm.of.yaml
derived:
  nothing
rejected:
  batchKeptColdEnough: batch/parmesan_tuesday must be kept at 8 °C or colder, but zone/dry_store is at 20 °C.
```

The tool rejects the batch of Tuesday, and the message gives the two
temperatures. The batch of Monday breaks no rule: the fridge is at 4 °C,
which is below 8 °C. The test of the file,
[`batchKeptTooWarm`](steps/step-06/tests/examples/batch-kept-too-warm.of.yaml),
expects exactly this rejection.

## A use-by day that cannot be correct

The kitchen also wants to find a mistake in the days of a batch. A second
constraint compares two fields of one batch:

```yaml
# defs/spoilage.of.yaml
useByAfterArrival:
  is: Constraint
  says: The use-by day of a batch is on or after the day that it arrived.
  match:
    - kind: { name: Batch, as: $batch }
  filter:
    - node:  $batch
      field: useBy
      is below:
        same as:
          node:  $batch
          field: receivedOn
  message: "{$batch} must not have a use-by day before the day that it arrived."
```

`is below` compares two days as it compares two numbers: the earlier day is
below the later day. In the example, each use-by day is after the day that
the batch arrived, so this constraint rejects nothing.

```console
$ of-lang test step-06
pass  test allergenThroughASauce  (step-06/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (step-06/tests/examples/base-recipe.of.yaml)
pass  test batchKeptTooWarm  (step-06/tests/examples/batch-kept-too-warm.of.yaml)
pass  test dishWithoutAllergenWarning  (step-06/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-06/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-06/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (step-06/tests/terms/garnish.of.yaml)
7 passed, 0 failed
```

A cook who records a delivery gets the message at once, and moves the batch
before the food spoils.

The model of the kitchen is complete. But the kitchen is not alone in the
restaurant: [Share the model with purchasing](07-share-with-purchasing.md).

More about this: [constraints](../../language/reference/constraint.md) and
[fields](../../language/reference/kind.md#fields) in the reference. In the
language by example: [Same as](../by-example/same-as.md),
[Comparison words](../by-example/comparison-words.md),
[Messages](../by-example/messages.md),
[Narrowing a value](../by-example/narrowing.md) and
[Expect rejected facts](../by-example/expect-rejected.md).
