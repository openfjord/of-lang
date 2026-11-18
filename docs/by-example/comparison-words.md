---
title: Comparison words — The OpenFjord modeling language by example
---

# Comparison words

> Part of [the language by example](README.md).

Some comparison words do not take `value:` or `same as:`. `is one of` takes
a list, `is a` takes a kind, and `is set` and `is not set` take only a
field.

<table>
<tr>
<td>

A recipe has a status, and can have a number of minutes of work.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food.
  fields:
    name:
      type: string
    status:
      enum: [active, seasonal, retired]
    prepMinutes:
      type: integer
      optional: true
```

</td>
</tr>
<tr>
<td>

A dish is a recipe.

</td>
<td>

```yaml
Dish:
  is: Kind
  extends: Recipe
  definition: A recipe that the restaurant sells to guests.
```

</td>
</tr>
<tr>
<td>

A test file states two recipes that are not dishes. The stock gives no
`prepMinutes`.

</td>
<td>

```yaml
# tests/examples/four-recipes.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish]

recipe/pesto:
  is: Fact
  kind: Recipe
  fields:
    name: pesto
    status: active
    prepMinutes: 15

recipe/stock:
  is: Fact
  kind: Recipe
  fields:
    name: stock
    status: retired
```

</td>
</tr>
<tr>
<td>

It also states two dishes. The fondue gives no `prepMinutes`.

</td>
<td>

```yaml
dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    status: seasonal
    prepMinutes: 20

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    status: active
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
fourRecipes:
  is: Test
  says: Two recipes and two dishes are valid, with or without a number of minutes.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

`is one of` takes `values:` with a list. It holds when the field equals one
of the values. Only the stock has a different status.

</td>
<td>

```yaml
# tests/questions/which-recipes-are-active-or-seasonal.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

whichRecipesAreActiveOrSeasonal:
  is: Question
  asks: Which recipes are active or seasonal?
  given: file:../examples/four-recipes.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - node: $recipe
      field: status
      is one of:
        values: [active, seasonal]
  expect:
    $recipe: [recipe/pesto, dish/margherita, dish/fondue]
```

</td>
</tr>
<tr>
<td>

`is a` takes a kind. It applies to the thing, so the condition has no
`field`. It holds when the thing has that kind.

</td>
<td>

```yaml
# tests/questions/which-recipes-are-dishes.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish]

whichRecipesAreDishes:
  is: Question
  asks: Which recipes are dishes?
  given: file:../examples/four-recipes.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - node: $recipe
      is a: Dish
  expect:
    $recipe: [dish/margherita, dish/fondue]
```

</td>
</tr>
<tr>
<td>

`is not set` holds when the field has no value. `is set` holds when the
field has a value. The two words compare with nothing.

</td>
<td>

```yaml
# tests/questions/which-recipes-have-no-minutes.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

whichRecipesHaveNoMinutes:
  is: Question
  asks: Which recipes do not give a number of minutes?
  given: file:../examples/four-recipes.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - node: $recipe
      field: prepMinutes
      is not set: true
  expect:
    $recipe: [recipe/stock, dish/fondue]
```

</td>
</tr>
</table>

`test` runs the test and the three questions:

```console
$ of-lang test comparison-words
pass  test fourRecipes  (comparison-words/tests/examples/four-recipes.of.yaml)
pass  question whichRecipesAreActiveOrSeasonal  (comparison-words/tests/questions/which-recipes-are-active-or-seasonal.of.yaml)
pass  question whichRecipesAreDishes  (comparison-words/tests/questions/which-recipes-are-dishes.of.yaml)
pass  question whichRecipesHaveNoMinutes  (comparison-words/tests/questions/which-recipes-have-no-minutes.of.yaml)
4 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
