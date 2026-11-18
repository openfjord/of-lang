---
title: At least one — The OpenFjord modeling language by example
---

# At least one

> Part of [the language by example](README.md).

A condition cannot see links directly. `at least one` holds a nested block
with its own `match` and `filter`, and holds when the block has one or more
matches.

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
# tests/examples/fresh-and-frozen.of.yaml
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

A file of facts holds exactly one test.

</td>
<td>

```yaml
freshAndFrozen:
  is: Test
  says: A dish with fresh ingredients and a dish with a frozen ingredient are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes use a frozen ingredient. Its `match` finds each
dish.

</td>
<td>

```yaml
# tests/questions/which-dishes-use-a-frozen-ingredient.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Ingredient, usesIngredient]

whichDishesUseAFrozenIngredient:
  is: Question
  asks: Which dishes use a frozen ingredient?
  given: file:../examples/fresh-and-frozen.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

The nested block runs one time for each dish. Its `match` finds the
ingredients that `$dish` uses. The variable `$ingredient` belongs to the
nested block, and nothing outside the block can use it.

</td>
<td>

```yaml
  filter:
    - at least one:
        match:
          - kind: { name: Ingredient, as: $ingredient }
          - link: usesIngredient
            from: $dish
            to:   $ingredient
```

</td>
</tr>
<tr>
<td>

The `filter` of the nested block keeps the frozen ingredients.

</td>
<td>

```yaml
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

For the scampi, the nested block has one match. For the margherita, it has
no matches.

</td>
<td>

```yaml
  expect:
    $dish: [dish/scampi]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test at-least-one
pass  test freshAndFrozen  (at-least-one/tests/examples/fresh-and-frozen.of.yaml)
pass  question whichDishesUseAFrozenIngredient  (at-least-one/tests/questions/which-dishes-use-a-frozen-ingredient.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
