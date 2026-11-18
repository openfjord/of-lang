---
title: Extends — The OpenFjord modeling language by example
---

# Extends

> Part of [the language by example](README.md).

A kind can extend a broader kind. It has every field of the broader kind,
and each of its things is also a thing of the broader kind.

<table>
<tr>
<td>

`Recipe` is the broader kind. Each recipe has a name and a preparation
time.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients.
  fields:
    name:
      type: string
    prepMinutes:
      type: integer
```

</td>
</tr>
<tr>
<td>

`extends` names the broader kind. A dish is a narrower sort of recipe.

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

A dish has the fields of a recipe, and it adds a field of its own. It
cannot declare `name` or `prepMinutes` again.

</td>
<td>

```yaml
  fields:
    price:
      type: number
```

</td>
</tr>
<tr>
<td>

A test file states one recipe and one dish. The dish gives the two fields
of a recipe, and also its price.

</td>
<td>

```yaml
# tests/examples/a-recipe-and-a-dish.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish]

recipe/tomato_sauce:
  is: Fact
  kind: Recipe
  fields:
    name: tomato_sauce
    prepMinutes: 40

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    prepMinutes: 12
    price: 12
```

</td>
</tr>
<tr>
<td>

The test expects what the definitions derive about the two kinds: a dish
is a recipe, but a recipe is not a dish.

</td>
<td>

```yaml
aRecipeAndADish:
  is: Test
  says: A dish is a recipe, and it has the fields of a recipe.
  expect:
    derived:
      - kind: Dish
        is a: Recipe
    notDerived:
      - kind: Recipe
        is a: Dish
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks for each recipe, given those facts. The answer includes
the margherita, because each dish is also a recipe.

</td>
<td>

```yaml
# tests/questions/which-recipes-are-there.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

whichRecipesAreThere:
  is: Question
  asks: Which recipes are there?
  given: file:../examples/a-recipe-and-a-dish.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
  expect:
    $recipe: [recipe/tomato_sauce, dish/margherita]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test extends
pass  test aRecipeAndADish  (extends/tests/examples/a-recipe-and-a-dish.of.yaml)
pass  question whichRecipesAreThere  (extends/tests/questions/which-recipes-are-there.of.yaml)
2 passed, 0 failed
```

More about this: [the kind rules](../../language/reference/kind.md#kind-rules)
and [derived terms](../../language/reference/test.md#derived-terms).
