---
title: Defaults — The OpenFjord modeling language by example
---

# Defaults

> Part of [the language by example](README.md).

A field with a `default` has that value when a fact does not give one.

<table>
<tr>
<td>

A field is required, unless it says something else. Each dish must have a
`name`.

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

`status` has a default. A dish that gives no status is `active`.

</td>
<td>

```yaml
    status:
      enum: [active, retired]
      default: active
```

</td>
</tr>
<tr>
<td>

To see the default at work, a test file states two dishes. The first one
gives no status.

</td>
<td>

```yaml
# tests/examples/two-dishes.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
```

</td>
</tr>
<tr>
<td>

The second one gives a status of its own.

</td>
<td>

```yaml
dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    status: retired
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test. This one expects that the package
rejects nothing: a dish without a status is valid.

</td>
<td>

```yaml
twoDishes:
  is: Test
  says: A dish without a status and a retired dish are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes are active, given those facts.

</td>
<td>

```yaml
# tests/questions/which-dishes-are-active.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesAreActive:
  is: Question
  asks: Which dishes are active?
  given: file:../examples/two-dishes.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
  filter:
    - node: $dish
      field: status
      equals:
        value: active
```

</td>
</tr>
<tr>
<td>

The answer is the margherita, although its fact does not state a status.

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
$ of-lang test defaults
pass  test twoDishes  (defaults/tests/examples/two-dishes.of.yaml)
pass  question whichDishesAreActive  (defaults/tests/questions/which-dishes-are-active.of.yaml)
2 passed, 0 failed
```

More about this: [fields](../../language/reference/kind.md#fields).
