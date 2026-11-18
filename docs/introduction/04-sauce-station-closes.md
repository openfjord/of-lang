---
title: Find what depends on what — Introduction to the OpenFjord modeling language
---

# Find what depends on what

> Part of [the introduction](README.md).

The kitchen has stations: the pizza station, the sauce station, the pastry
station. One morning, the fridge of the sauce station breaks, and the
station closes for the day. The head cook must know quickly which other
stations cannot make all of their recipes. The answer is in the recipes: a
station stops when one of its recipes uses a recipe of the closed station.

## Where the cooks make each recipe

The model does not know the stations yet. The kitchen adds a kind for them,
in a file of its own:

```yaml
# defs/station.of.yaml
Station:
  is: Kind
  definition: A part of the kitchen where a group of cooks makes some of the recipes.
  examples:
    - the pizza station
    - the pastry station
  fields:
    name:
      type: Name
```

Cooks make each recipe at a station. That is a link:

```yaml
# defs/recipe.of.yaml
madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.
  inverse: makesRecipe
```

`inverse` gives a name to the link in the opposite direction. `makesRecipe`
from a station to a recipe is the same relationship as `madeAtStation` from
the recipe to the station. A match can use the name that reads best.

`Dish` extends `Recipe`, so a dish can have each link that starts at
`Recipe`. The kitchen does not declare the link a second time for dishes.

In the files of facts, the kitchen adds the pizza station and the sauce
station, and each recipe gets its station. The pesto and the tomato sauce
get `madeAtStation: station/sauce`. The two dishes get
`madeAtStation: station/pizza`.

## A station depends on a station

No person states which station depends on which. That follows from the
recipes, so a derivative adds it:

```yaml
# defs/station.of.yaml
dependsOnStation:
  is: Link
  from: Station
  to: Station
  reads: "{from} depends on {to}"
  definition: The first station cannot make all of its recipes when the second station stops work.

stationDependsOnStation:
  is: Derivative
  says: >
    A station depends on each different station that makes a recipe that one
    of its own recipes uses.
  match:
    - kind: { name: Station, as: $station }
    - kind: { name: Recipe, as: $recipe }
    - kind: { name: Recipe, as: $part }
    - kind: { name: Station, as: $other }
    - link: makesRecipe
      from: $station
      to:   $recipe
    - link: usesRecipe
      from: $recipe
      to:   $part
    - link: madeAtStation
      from: $part
      to:   $other
  filter:
    - node: $other
      does not equal:
        same as: $station
  derive:
    - link: dependsOnStation
      from: $station
      to:   $other
```

The match follows three links: from a station to one of its recipes, from
that recipe to a recipe that it uses, and from there to the station that
makes that recipe. The first link item uses the inverse name `makesRecipe`,
so each item reads from left to right.

Without the filter, each station that uses its own recipes depends on
itself. The condition compares two things: `same as` names the variable to
compare with, and `does not equal` keeps only the matches where the two
stations are different.

`defs/station.of.yaml` imports from `defs/recipe.of.yaml`, and
`defs/recipe.of.yaml` imports from it. Files of one package can import each
other in a circle.

## The answer

```console
$ of-lang run step-04 step-04/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/milk  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  station/pizza depends on station/sauce  (derivative stationDependsOnStation)
rejected:
  nothing
```

The last line under `derived` is the answer for this kitchen: the pizza
station depends on the sauce station. `explain` shows which recipes cause
it:

```console
$ of-lang explain step-04 step-04/tests/examples/allergen-through-a-sauce.of.yaml station/pizza dependsOnStation station/sauce
station/pizza depends on station/sauce  (derivative stationDependsOnStation)
  A station depends on each different station that makes a recipe that one of its own recipes uses.
  with $other = station/sauce, $part = recipe/tomato_sauce, $recipe = dish/margherita, $station = station/pizza
  because:
    dish/margherita is made at station/pizza  (written at step-04/tests/examples/allergen-through-a-sauce.of.yaml:113:5)
    dish/margherita uses recipe/tomato_sauce  (written at step-04/tests/examples/allergen-through-a-sauce.of.yaml:114:5)
    recipe/tomato_sauce is made at station/sauce  (written at step-04/tests/examples/allergen-through-a-sauce.of.yaml:103:5)
```

The head cook asks the model in the same way as in
[chapter 2](02-dishes-with-nuts.md). The question
[`whichStationsStopWithTheSauceStation`](steps/step-04/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
finds each station with a `dependsOnStation` link to `station/sauce`, and
expects the answer `[station/pizza]`.

`usesRecipe` chains. Thus, the answer also includes a station whose recipe
uses a recipe of the sauce station through other recipes.

## Each recipe has one station

The answer is only complete if each recipe tells where cooks make it. A link
alone does not make sure of this: a thing can have a link any number of
times, or not at all. A constraint sets the number:

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

`count` counts the matches of a nested block, and the condition compares the
number with a literal `value`. In the nested block, `to: Station` is a kind
alone: the station at that end only has to exist, so it needs no variable.
This constraint has no `unless`. Each match that passes the filter is a
violation: a recipe without a station, or a recipe with two stations.

## The story of the margherita, complete

The test from [chapter 1](01-margherita-and-nuts.md) expects one derived
fact. The kitchen makes
[the test](steps/step-04/tests/examples/allergen-through-a-sauce.of.yaml)
tell all that the example shows so far:

- Under `derived` are the facts that the tool must add: the margherita
  contains nuts and milk, it uses parmesan, and the pizza station depends on
  the sauce station.
- Under `notDerived` are facts that the tool must not add: the marinara
  does not contain nuts, and the sauce station does not depend on the pizza
  station. A dependency has a direction.
- `rejected: []` tells that the facts break no constraint.

```console
$ of-lang test step-04
pass  test allergenThroughASauce  (step-04/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test dishWithoutAllergenWarning  (step-04/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-04/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-04/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
4 passed, 0 failed
```

When a station closes, the head cook knows what stops, and why. When a cook
moves a recipe to a different station, the answer follows.

Some recipes are in so many dishes that each change to them is a risk:
[Find what many things share](05-recipes-that-need-care.md)

More about this: [links](../../language/reference/link.md),
[derivatives](../../language/reference/derivative.md),
[constraints](../../language/reference/constraint.md) and
[tests](../../language/reference/test.md#examples) in the reference. In the
language by example: [Inverse names](../by-example/inverse-names.md),
[Links that chain](../by-example/chains.md),
[Same as](../by-example/same-as.md), [Count](../by-example/count.md),
[A kind alone](../by-example/kind-alone.md) and
[Expect derived facts](../by-example/expect-derived.md).
