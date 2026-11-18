---
title: Find what many things share — Introduction to the OpenFjord modeling language
---

# Find what many things share

> Part of [the introduction](README.md).

A cook wants to put less salt in the tomato sauce. That looks like a small
change. But the tomato sauce is in the margherita, the marinara and the
lasagne, so the change is a change to three dishes. The kitchen wants to
know which recipes are like this, before a cook changes one of them.

## A name for such a recipe

The kitchen calls a recipe that many recipes use a base recipe. No person
decides which recipes are base recipes. The model counts:

```yaml
# defs/recipe.of.yaml
BaseRecipe:
  is: Kind
  extends: Recipe
  definition: A recipe that many other recipes use, so a change to it needs extra care.

baseRecipeFromUses:
  is: Derivative
  says: >
    A recipe that more than two recipes use, directly or through other
    recipes, is a base recipe.
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - count:
        match:
          - link: usesRecipe
            from: Recipe
            to:   $recipe
      is above:
        value: 2
  derive:
    - node: $recipe
      becomes: BaseRecipe
```

This derivative adds no link. `node` and `becomes` add a kind to a thing.
The kind must extend the kind of the variable, as `BaseRecipe` extends
`Recipe`.

The filter counts the recipes that use `$recipe`, and keeps the recipes
where the number is above 2. `usesRecipe` chains, so the count includes the
recipes that use `$recipe` through other recipes.

## Three dishes with the tomato sauce

A new file of facts,
[`tests/examples/base-recipe.of.yaml`](steps/step-05/tests/examples/base-recipe.of.yaml),
has three dishes that use the tomato sauce. Only the margherita uses the
pesto. The answer:

```console
$ of-lang run step-05 step-05/tests/examples/base-recipe.of.yaml
derived:
  recipe/tomato_sauce is a BaseRecipe  (derivative baseRecipeFromUses)
  station/pizza depends on station/sauce  (derivative stationDependsOnStation)
rejected:
  nothing
```

The tomato sauce is a base recipe, and the pesto is not. The second line
comes from [chapter 4](04-sauce-station-closes.md): each derivative of the
package runs on each file of facts.

```console
$ of-lang explain step-05 step-05/tests/examples/base-recipe.of.yaml recipe/tomato_sauce BaseRecipe
recipe/tomato_sauce is a BaseRecipe  (derivative baseRecipeFromUses)
  A recipe that more than two recipes use, directly or through other recipes, is a base recipe.
  with $recipe = recipe/tomato_sauce
```

The test of the file keeps the two halves of the answer:

```yaml
# tests/examples/base-recipe.of.yaml
tomatoSauceInThreeDishes:
  is: Test
  says: >
    A recipe that three dishes use is a base recipe. A recipe that only one
    dish uses is not a base recipe.
  expect:
    derived:
      - node:    recipe/tomato_sauce
        becomes: BaseRecipe
    notDerived:
      - node:    recipe/pesto
        becomes: BaseRecipe
```

## What a cook reads before a change

Extra care starts with the steps of the recipe as they are, and with the
time that the work takes. The kitchen keeps the steps of each recipe on a
page on the web. Each recipe tells where:

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
  examples:
    - the tomato sauce
    - the pizza dough
  notExamples:
    - a food that the restaurant buys and does not make, such as flour
    - a new dish that a cook tries one time and does not keep
  fields:
    name:
      type: Name
    steps:
      type: URL
    prepMinutes:
      type: Minutes
      optional: true
```

`steps` is required, so each recipe in the facts gets one, such as
`steps: https://recipes.example.com/pesto`. `optional: true` makes a field
optional, so a recipe can be without `prepMinutes`.

The two fields have value types of their own:

```yaml
# defs/names.of.yaml
URL:
  is: ValueType
  definition: The address of a page on the web.
  type: string
  format: uri

Minutes:
  is: ValueType
  definition: A number of minutes of work.
  type: integer
  minimum: 1
  maximum: 1440
```

A value type has a built-in `type`, such as `string` or `integer`, and
keywords that narrow it. `format: uri` allows only an absolute URI (Uniform
Resource Identifier), such as the address of a page on the web. A value of
`Minutes` is a whole number from 1 to 1440. The tool does not accept a fact
with a value that does not fit.

## The answers depend on the definitions

Each answer so far depends on a few words in the definitions. The dishes
with nuts depend on `extends`: a dish is a recipe. The milk in the
margherita depends on `under`. The base recipes and the stations depend on
`chains`. If somebody deletes one of these words, the answers change, and
possibly no example shows it. The kitchen writes down what the definitions
must tell, in a test without facts:

```yaml
# tests/terms/garnish.of.yaml
garnishIsAnIngredient:
  is: Test
  says: >
    A garnish is one of the ingredients that a recipe uses. A dish is a
    recipe, so a dish can have a station. A recipe can use a recipe that uses
    another recipe.
  expect:
    derived:
      - kind: Dish
        is a: Recipe
      - kind: Dish
        can have: madeAtStation
      - link: usesGarnish
        implies: usesIngredient
      - link: usesRecipe
        chains: true
    notDerived:
      - link: usesIngredient
        implies: usesGarnish
```

Each item is a derived term: something that the definitions tell about
kinds and links, before there are any things. The item under `notDerived`
tells that `under` goes in one direction only: an ingredient of a recipe is
not always a garnish.

```console
$ of-lang test step-05
pass  test allergenThroughASauce  (step-05/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (step-05/tests/examples/base-recipe.of.yaml)
pass  test dishWithoutAllergenWarning  (step-05/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-05/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-05/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (step-05/tests/terms/garnish.of.yaml)
6 passed, 0 failed
```

The kitchen knows which recipes need extra care, and the list follows the
recipes: when two more dishes use the pesto, the pesto becomes a base recipe
without a decision from a person.

A wrong recipe is not the only risk in a kitchen. Food also spoils:
[Check one value against another](06-food-kept-cold.md)

More about this: [derivatives](../../language/reference/derivative.md),
[value types](../../language/reference/kind.md#value-types) and
[derived terms](../../language/reference/test.md#derived-terms) in the
reference. In the language by example:
[Derive a kind](../by-example/derive-kind.md),
[Count](../by-example/count.md), [Fields](../by-example/fields.md),
[Value types](../by-example/value-types.md),
[Narrowing a value](../by-example/narrowing.md) and
[Derived terms](../by-example/derived-terms.md).
