---
title: None — The OpenFjord modeling language by example
---

# None

> Part of [the language by example](README.md).

`none` holds a nested block, and holds when the block has no matches.

<table>
<tr>
<td>

An ingredient is fresh or frozen.

</td>
<td>

```yaml
# defs/kitchen.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    name:
      type: string

Ingredient:
  is: Kind
  definition: A food that the kitchen buys and does not make.
  fields:
    name:
      type: string
    storage:
      enum: [fresh, frozen]
```

</td>
</tr>
<tr>
<td>

A dish uses ingredients.

</td>
<td>

```yaml
usesIngredient:
  is: Link
  from: Dish
  to: Ingredient
  reads: "{from} uses {to}"
  definition: The dish has the ingredient in it.
```

</td>
</tr>
<tr>
<td>

A test file states two fresh ingredients and one frozen ingredient.

</td>
<td>

```yaml
# tests/examples/three-dishes.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Ingredient, usesIngredient]

ingredient/tomatoes:
  is: Fact
  kind: Ingredient
  fields:
    name: tomatoes
    storage: fresh

ingredient/basil:
  is: Fact
  kind: Ingredient
  fields:
    name: basil
    storage: fresh

ingredient/shrimp:
  is: Fact
  kind: Ingredient
  fields:
    name: shrimp
    storage: frozen
```

</td>
</tr>
<tr>
<td>

The margherita uses only fresh ingredients. The scampi uses shrimp, which is
frozen.

</td>
<td>

```yaml
dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
  links:
    usesIngredient: [ingredient/tomatoes, ingredient/basil]

dish/scampi:
  is: Fact
  kind: Dish
  fields:
    name: scampi
  links:
    usesIngredient: [ingredient/tomatoes, ingredient/shrimp]
```

</td>
</tr>
<tr>
<td>

The bread states no ingredients.

</td>
<td>

```yaml
dish/bread:
  is: Fact
  kind: Dish
  fields:
    name: bread
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threeDishes:
  is: Test
  says: Dishes with fresh ingredients, with a frozen ingredient and with no ingredients are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes use no frozen ingredient. Its `match` finds
each dish.

</td>
<td>

```yaml
# tests/questions/which-dishes-use-no-frozen-ingredient.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Ingredient, usesIngredient]

whichDishesUseNoFrozenIngredient:
  is: Question
  asks: Which dishes use no frozen ingredient?
  given: file:../examples/three-dishes.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

The nested block runs one time for each dish. It finds the frozen
ingredients that `$dish` uses. `none` holds when there are no such
ingredients.

</td>
<td>

```yaml
  filter:
    - none:
        match:
          - kind: { name: Ingredient, as: $ingredient }
          - link: usesIngredient
            from: $dish
            to:   $ingredient
        filter:
          - node: $ingredient
            field: storage
            equals:
              value: frozen
```

</td>
</tr>
<tr>
<td>

The answer has the margherita, and also the bread. A link that is not stated
does not exist, so the nested block has no matches for the bread.

</td>
<td>

```yaml
  expect:
    $dish: [dish/margherita, dish/bread]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test none
pass  test threeDishes  (none/tests/examples/three-dishes.of.yaml)
pass  question whichDishesUseNoFrozenIngredient  (none/tests/questions/which-dishes-use-no-frozen-ingredient.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
