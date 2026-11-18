---
title: Fields — The OpenFjord modeling language by example
---

# Fields

> Part of [the language by example](README.md).

A field is a value of the thing itself, for example a name or a price. A
kind lists its fields, and each field has a type.

<table>
<tr>
<td>

`fields` lists the fields of the kind, by name. `type` gives the type of
the value. A `string` is text.

</td>
<td>

```yaml
# defs/dish.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    name:
      type: string
```

</td>
</tr>
<tr>
<td>

A `number` can have decimals, and an `integer` is a whole number. A
`boolean` is `true` or `false`.

</td>
<td>

```yaml
    price:
      type: number
    servings:
      type: integer
    vegetarian:
      type: boolean
```

</td>
</tr>
<tr>
<td>

A field is required, unless it says something else. With `optional: true`,
a fact can omit the field.

</td>
<td>

```yaml
    chefNote:
      type: string
      optional: true
```

</td>
</tr>
<tr>
<td>

A test file states two dishes. The first one gives a value for each
field.

</td>
<td>

```yaml
# tests/examples/two-dishes.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    price: 24.5
    servings: 2
    vegetarian: true
    chefNote: Stir the cheese until the guests come.
```

</td>
</tr>
<tr>
<td>

The second one omits `chefNote`, which is the only optional field.

</td>
<td>

```yaml
dish/carbonara:
  is: Fact
  kind: Dish
  fields:
    name: carbonara
    price: 14
    servings: 1
    vegetarian: false
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so each
dish has every required field, and each value has the type of its field.

</td>
<td>

```yaml
twoDishes:
  is: Test
  says: A dish with a note from the chef and a dish without one are both valid.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test fields
pass  test twoDishes  (fields/tests/examples/two-dishes.of.yaml)
1 passed, 0 failed
```

More about this: [fields](../../language/reference/kind.md#fields) and
[the rules for facts](../../language/reference/fact.md#rules).
