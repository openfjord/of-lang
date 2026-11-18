---
title: Unions — The OpenFjord modeling language by example
---

# Unions

> Part of [the language by example](README.md).

Sometimes a kind is exactly one of a fixed list of kinds, and nothing else.
`oneOf` on a kind states this. Such a kind is a union.

<table>
<tr>
<td>

The kitchen has recipes and ingredients. A recipe uses the two of them to
make its food.

</td>
<td>

```yaml
# defs/food.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
  fields:
    name:
      type: string

Ingredient:
  is: Kind
  definition: A food that the restaurant buys and does not make.
  fields:
    name:
      type: string
```

</td>
</tr>
<tr>
<td>

`oneOf` lists the kinds that together make up `Food`. Each food is a recipe
or an ingredient, and never both. A union has no fields. The listed kinds
keep their own fields.

</td>
<td>

```yaml
Food:
  is: Kind
  definition: Something that a recipe uses to make its food.
  oneOf: [Recipe, Ingredient]
```

</td>
</tr>
<tr>
<td>

A link can end at the union. Then it can end at a recipe or at an
ingredient.

</td>
<td>

```yaml
usesFood:
  is: Link
  from: Recipe
  to: Food
  reads: "{from} uses {to}"
  definition: The recipe needs the food.
```

</td>
</tr>
<tr>
<td>

A test file states one ingredient and two recipes. A fact names a listed
kind as its `kind`, and never the union.

</td>
<td>

```yaml
# tests/examples/a-pizza.of.yaml
import:
  - from: file:../../defs/food.of.yaml
    use: [Recipe, Ingredient, Food, usesFood]

ingredient/mozzarella:
  is: Fact
  kind: Ingredient
  fields:
    name: mozzarella

recipe/tomato_sauce:
  is: Fact
  kind: Recipe
  fields:
    name: tomato_sauce
```

</td>
</tr>
<tr>
<td>

The margherita uses a recipe and an ingredient through the same link.

</td>
<td>

```yaml
recipe/margherita:
  is: Fact
  kind: Recipe
  fields:
    name: margherita
  links:
    usesFood: [recipe/tomato_sauce, ingredient/mozzarella]
```

</td>
</tr>
<tr>
<td>

The test expects what the definitions derive about the kinds: a recipe is
a food, and a recipe is never an ingredient.

</td>
<td>

```yaml
aPizza:
  is: Test
  says: The margherita uses a recipe and an ingredient, and each one is a food.
  expect:
    derived:
      - kind: Recipe
        is a: Food
      - kind: Recipe
        never is: Ingredient
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question can match the union. It asks for each food that the margherita
uses, and the answer has a recipe and an ingredient.

</td>
<td>

```yaml
# tests/questions/which-foods-does-the-margherita-use.of.yaml
import:
  - from: file:../../defs/food.of.yaml
    use: [Recipe, Food, usesFood]

whichFoodsDoesTheMargheritaUse:
  is: Question
  asks: Which foods does the margherita use?
  given: file:../examples/a-pizza.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
    - kind: { name: Food, as: $food }
    - link: usesFood
      from: $recipe
      to:   $food
  filter:
    - node: $recipe
      equals:
        value: recipe/margherita
  expect:
    $food: [recipe/tomato_sauce, ingredient/mozzarella]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test unions
pass  test aPizza  (unions/tests/examples/a-pizza.of.yaml)
pass  question whichFoodsDoesTheMargheritaUse  (unions/tests/questions/which-foods-does-the-margherita-use.of.yaml)
2 passed, 0 failed
```

More about this: [unions](../../language/reference/kind.md#unions) and
[derived terms](../../language/reference/test.md#derived-terms).
