---
title: Value types — The OpenFjord modeling language by example
---

# Value types

> Part of [the language by example](README.md).

A value type gives a name to a schema one time, so that many fields can use
the schema again.

<table>
<tr>
<td>

A value type is a statement with `is: ValueType`. Like a kind, it needs a
`definition`.

</td>
<td>

```yaml
# defs/types.of.yaml
Minutes:
  is: ValueType
  definition: A number of minutes of work.
```

</td>
</tr>
<tr>
<td>

The other keys are the schema: a type and the keywords that go with it. A
number of minutes is a whole number from 1 to 1440.

</td>
<td>

```yaml
  type: integer
  minimum: 1
  maximum: 1440
```

</td>
</tr>
<tr>
<td>

A file imports each name that it uses but does not define. The file of the
recipe imports `Minutes`.

</td>
<td>

```yaml
# defs/recipe.of.yaml
import:
  - from: file:types.of.yaml
    use: [Minutes]
```

</td>
</tr>
<tr>
<td>

A field names the value type as its `type`. Two fields of the recipe use
`Minutes`.

</td>
<td>

```yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients.
  fields:
    name:
      type: string
    prepMinutes:
      type: Minutes
    cookMinutes:
      type: Minutes
```

</td>
</tr>
<tr>
<td>

A test file states two recipes. Each value is a valid number of minutes,
so the test passes.

</td>
<td>

```yaml
# tests/examples/two-recipes.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

recipe/pesto:
  is: Fact
  kind: Recipe
  fields:
    name: pesto
    prepMinutes: 15
    cookMinutes: 1

recipe/ragu:
  is: Fact
  kind: Recipe
  fields:
    name: ragu
    prepMinutes: 20
    cookMinutes: 180

twoRecipes:
  is: Test
  says: The times of the pesto and of the ragu are valid numbers of minutes.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A value type is a type of its own. A comparison needs two values of the
same type, and the two fields both name `Minutes`. Thus, a question can
compare them.

</td>
<td>

```yaml
# tests/questions/which-recipes-cook-longer-than-they-prepare.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

whichRecipesCookLongerThanTheyPrepare:
  is: Question
  asks: Which recipes take more time to cook than to prepare?
  given: file:../examples/two-recipes.of.yaml
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - node: $recipe
      field: cookMinutes
      is above:
        same as:
          node: $recipe
          field: prepMinutes
  expect:
    $recipe: [recipe/ragu]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test value-types
pass  test twoRecipes  (value-types/tests/examples/two-recipes.of.yaml)
pass  question whichRecipesCookLongerThanTheyPrepare  (value-types/tests/questions/which-recipes-cook-longer-than-they-prepare.of.yaml)
2 passed, 0 failed
```

More about this: [value types](../../language/reference/kind.md#value-types).
