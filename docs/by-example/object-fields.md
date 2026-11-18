---
title: Objects and arrays — The OpenFjord modeling language by example
---

# Objects and arrays

> Part of [the language by example](README.md).

A field can hold more than one value. An `object` has named parts, and an
`array` is a list of values.

<table>
<tr>
<td>

Each dish has a `name`, which is one value.

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

`properties` lists the parts of an `object`, by name. Each part has a
schema, like a field. A part is required, unless it has `optional: true`.

</td>
<td>

```yaml
    plating:
      type: object
      properties:
        plate:
          type: string
        garnish:
          type: string
          optional: true
```

</td>
</tr>
<tr>
<td>

`items` gives the schema of each value of an `array`. Here, each value is
text.

</td>
<td>

```yaml
    allergens:
      type: array
      items:
        type: string
```

</td>
</tr>
<tr>
<td>

A test file states two dishes. The first one gives the two parts of
`plating`, and a list of two allergens.

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
    plating:
      plate: wooden board
      garnish: basil
    allergens: [gluten, milk]
```

</td>
</tr>
<tr>
<td>

The second one omits the optional part, and its list is empty.

</td>
<td>

```yaml
dish/minestrone:
  is: Fact
  kind: Dish
  fields:
    name: minestrone
    plating:
      plate: deep bowl
    allergens: []

twoDishes:
  is: Test
  says: A dish with a garnish and a dish without one are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A condition names one part of an `object` field with `part`. This question
asks which dishes have a deep bowl as their plate.

</td>
<td>

```yaml
# tests/questions/which-dishes-go-in-a-deep-bowl.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesGoInADeepBowl:
  is: Question
  asks: Which dishes go in a deep bowl?
  given: file:../examples/two-dishes.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
  filter:
    - node: $dish
      field: plating
      part: plate
      equals:
        value: deep bowl
  expect:
    $dish: [dish/minestrone]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test object-fields
pass  test twoDishes  (object-fields/tests/examples/two-dishes.of.yaml)
pass  question whichDishesGoInADeepBowl  (object-fields/tests/questions/which-dishes-go-in-a-deep-bowl.of.yaml)
2 passed, 0 failed
```

More about this: [the keywords of a field](../../language/reference/kind.md#keywords)
and [conditions on a part](../../language/reference/constraint.md#filter-and-unless).
