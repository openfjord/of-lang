---
title: Variants — The OpenFjord modeling language by example
---

# Variants

> Part of [the language by example](README.md).

Some values have one of a small number of shapes. `oneOf` on a field lists
the shapes, and each value matches exactly one of them.

<table>
<tr>
<td>

A recipe tells how much food it makes. The kitchen counts some foods in
portions, and it weighs other foods. Thus, `yield` has two shapes.

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
    yield:
      oneOf:
```

</td>
</tr>
<tr>
<td>

Each choice is a schema. The first choice is an object with a number of
portions. Its `variant` part has a fixed value, which `const` gives.

</td>
<td>

```yaml
        - type: object
          properties:
            variant:
              const: portions
            count:
              type: integer
```

</td>
</tr>
<tr>
<td>

The second choice is an object with a weight. Its `variant` part has a
different fixed value, so a reader can see which shape a value has.

</td>
<td>

```yaml
        - type: object
          properties:
            variant:
              const: weight
            grams:
              type: integer
```

</td>
</tr>
<tr>
<td>

A test file states two recipes. The risotto gives its yield in portions,
and the pizza dough gives it as a weight.

</td>
<td>

```yaml
# tests/examples/two-recipes.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

recipe/risotto:
  is: Fact
  kind: Recipe
  fields:
    name: risotto
    yield:
      variant: portions
      count: 4

recipe/pizza_dough:
  is: Fact
  kind: Recipe
  fields:
    name: pizza_dough
    yield:
      variant: weight
      grams: 900
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so each
yield matches exactly one choice.

</td>
<td>

```yaml
twoRecipes:
  is: Test
  says: A yield in portions and a yield as a weight are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A condition reads a part of the value with `part`. This question uses the
`variant` part to find the recipes that give portions.

</td>
<td>

```yaml
# tests/questions/which-recipes-give-portions.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

whichRecipesGivePortions:
  is: Question
  asks: Which recipes give their yield in portions?
  given: file:../examples/two-recipes.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - node: $recipe
      field: yield
      part: variant
      equals:
        value: portions
  expect:
    $recipe: [recipe/risotto]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test variants
pass  test twoRecipes  (variants/tests/examples/two-recipes.of.yaml)
pass  question whichRecipesGivePortions  (variants/tests/questions/which-recipes-give-portions.of.yaml)
2 passed, 0 failed
```

More about this: [variants](../../language/reference/kind.md#variants) and
[the field rules](../../language/reference/kind.md#field-rules).
