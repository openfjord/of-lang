---
title: Parts — The OpenFjord modeling language by example
---

# Parts

> Part of [the language by example](README.md).

A condition on an `object` field names one part of the value with `part`.

<table>
<tr>
<td>

The field `plating` is an object with two parts. Each dish must have a
`plate`. The `garnish` is optional.

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
    plating:
      type: object
      properties:
        plate:   { type: string }
        garnish: { type: string, optional: true }
```

</td>
</tr>
<tr>
<td>

A test file states a dish that cooks serve on a board.

</td>
<td>

```yaml
# tests/examples/three-platings.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    plating:
      plate: board
      garnish: basil
```

</td>
</tr>
<tr>
<td>

Cooks serve the marinara on a board too, and the fondue in a pot.

</td>
<td>

```yaml
dish/marinara:
  is: Fact
  kind: Dish
  fields:
    name: marinara
    plating:
      plate: board

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    plating:
      plate: pot
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threePlatings:
  is: Test
  says: Dishes on a board and in a pot are valid, with or without a garnish.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes cooks serve on a board.

</td>
<td>

```yaml
# tests/questions/which-dishes-are-on-a-board.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesAreOnABoard:
  is: Question
  asks: Which dishes do cooks serve on a board?
  given: file:../examples/three-platings.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

`part` goes next to `field`. The condition compares the `plate` part of
`plating`, not the whole object. For a part in a part, `part` is a list of
the names.

</td>
<td>

```yaml
  filter:
    - node: $dish
      field: plating
      part: plate
      equals:
        value: board
```

</td>
</tr>
<tr>
<td>

The answer is the two dishes on a board.

</td>
<td>

```yaml
  expect:
    $dish: [dish/margherita, dish/marinara]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test parts
pass  test threePlatings  (parts/tests/examples/three-platings.of.yaml)
pass  question whichDishesAreOnABoard  (parts/tests/questions/which-dishes-are-on-a-board.of.yaml)
2 passed, 0 failed
```

More about this: [filter and
unless](../../language/reference/constraint.md#filter-and-unless),
[fields](../../language/reference/kind.md#fields).
