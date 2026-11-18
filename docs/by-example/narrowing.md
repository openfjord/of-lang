---
title: Narrowing a value — The OpenFjord modeling language by example
---

# Narrowing a value

> Part of [the language by example](README.md).

A type allows many values. Keywords next to the type narrow the field to
the values that are correct for it.

<table>
<tr>
<td>

`pattern` is a regular expression that the text must match. The code of a
recipe starts with a lowercase letter. Then it has lowercase letters,
digits and underscores.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients.
  fields:
    code:
      type: string
      pattern: "^[a-z][a-z0-9_]*$"
```

</td>
</tr>
<tr>
<td>

`minLength` and `maxLength` set how many characters the text can have.

</td>
<td>

```yaml
    title:
      type: string
      minLength: 3
      maxLength: 40
```

</td>
</tr>
<tr>
<td>

`minimum` and `maximum` set the smallest number and the largest number.

</td>
<td>

```yaml
    servings:
      type: integer
      minimum: 1
      maximum: 12
```

</td>
</tr>
<tr>
<td>

`format` tells what the text holds. With `uri`, it is an absolute URI
(Uniform Resource Identifier), such as an address that starts with
`https:`. The other formats are `date`, `date-time` and `duration`.

</td>
<td>

```yaml
    steps:
      type: string
      format: uri
```

</td>
</tr>
<tr>
<td>

`const` allows exactly one value. Each recipe gives its amounts in metric
units.

</td>
<td>

```yaml
    units:
      const: metric
```

</td>
</tr>
<tr>
<td>

A test file states one recipe. Each value satisfies the keywords of its
field.

</td>
<td>

```yaml
# tests/examples/one-recipe.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe]

recipe/tomato_sauce:
  is: Fact
  kind: Recipe
  fields:
    code: tomato_sauce
    title: Tomato sauce
    servings: 8
    steps: https://recipes.example.com/tomato_sauce
    units: metric
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so the
fields allow all five values.

</td>
<td>

```yaml
oneRecipe:
  is: Test
  says: Each value of the tomato sauce is a value that its field allows.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test narrowing
pass  test oneRecipe  (narrowing/tests/examples/one-recipe.of.yaml)
1 passed, 0 failed
```

More about this: [the keywords of a field](../../language/reference/kind.md#keywords)
and [the field rules](../../language/reference/kind.md#field-rules).
