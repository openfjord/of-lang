---
title: Count — The OpenFjord modeling language by example
---

# Count

> Part of [the language by example](README.md).

`count` holds a nested block, and compares the number of its matches with a
number.

<table>
<tr>
<td>

The package has dishes and ingredients.

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

A test file states three ingredients.

</td>
<td>

```yaml
# tests/examples/two-and-three-ingredients.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Ingredient, usesIngredient]

ingredient/tomatoes:
  is: Fact
  kind: Ingredient
  fields:
    name: tomatoes

ingredient/mozzarella:
  is: Fact
  kind: Ingredient
  fields:
    name: mozzarella

ingredient/basil:
  is: Fact
  kind: Ingredient
  fields:
    name: basil
```

</td>
</tr>
<tr>
<td>

The margherita uses three ingredients, and the marinara uses two.

</td>
<td>

```yaml
dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
  links:
    usesIngredient: [ingredient/tomatoes, ingredient/mozzarella, ingredient/basil]

dish/marinara:
  is: Fact
  kind: Dish
  fields:
    name: marinara
  links:
    usesIngredient: [ingredient/tomatoes, ingredient/basil]
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
twoAndThreeIngredients:
  is: Test
  says: A dish with three ingredients and a dish with two ingredients are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes use more than two ingredients. Its `match`
finds each dish.

</td>
<td>

```yaml
# tests/questions/which-dishes-use-more-than-two-ingredients.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Ingredient, usesIngredient]

whichDishesUseMoreThanTwoIngredients:
  is: Question
  asks: Which dishes use more than two ingredients?
  given: file:../examples/two-and-three-ingredients.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

The nested block runs one time for each dish. It has one match for each
ingredient that `$dish` uses.

</td>
<td>

```yaml
  filter:
    - count:
        match:
          - link: usesIngredient
            from: $dish
            to:   Ingredient
```

</td>
</tr>
<tr>
<td>

The comparison word goes next to `count`. It takes `value:` with a whole
number. The condition holds when the number of matches is above 2.

</td>
<td>

```yaml
      is above:
        value: 2
```

</td>
</tr>
<tr>
<td>

The answer is the margherita. The count for the marinara is 2.

</td>
<td>

```yaml
  expect:
    $dish: [dish/margherita]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test count
pass  test twoAndThreeIngredients  (count/tests/examples/two-and-three-ingredients.of.yaml)
pass  question whichDishesUseMoreThanTwoIngredients  (count/tests/questions/which-dishes-use-more-than-two-ingredients.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
