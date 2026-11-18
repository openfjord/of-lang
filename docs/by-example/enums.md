---
title: Enums — The OpenFjord modeling language by example
---

# Enums

> Part of [the language by example](README.md).

A field with `enum` has a fixed list of values. A fact gives one value from
the list.

<table>
<tr>
<td>

Each dish has a `name`, which can be any text.

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

`enum` lists the values that `course` can have. A dish is a starter, a
main course or a dessert, and nothing else.

</td>
<td>

```yaml
    course:
      enum: [starter, main, dessert]
```

</td>
</tr>
<tr>
<td>

A test file states two dishes. Each one gives a course from the list.

</td>
<td>

```yaml
# tests/examples/two-dishes.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/bruschetta:
  is: Fact
  kind: Dish
  fields:
    name: bruschetta
    course: starter

dish/tiramisu:
  is: Fact
  kind: Dish
  fields:
    name: tiramisu
    course: dessert
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so the
field allows the two courses.

</td>
<td>

```yaml
twoDishes:
  is: Test
  says: A starter and a dessert are both valid dishes.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes are desserts, given those facts. A condition
compares the field with a value from the list.

</td>
<td>

```yaml
# tests/questions/which-dishes-are-desserts.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesAreDesserts:
  is: Question
  asks: Which dishes are desserts?
  given: file:../examples/two-dishes.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
  filter:
    - node: $dish
      field: course
      equals:
        value: dessert
```

</td>
</tr>
<tr>
<td>

The answer is the tiramisu.

</td>
<td>

```yaml
  expect:
    $dish: [dish/tiramisu]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test enums
pass  test twoDishes  (enums/tests/examples/two-dishes.of.yaml)
pass  question whichDishesAreDesserts  (enums/tests/questions/which-dishes-are-desserts.of.yaml)
2 passed, 0 failed
```

More about this: [the keywords of a field](../../language/reference/kind.md#keywords)
and [the field rules](../../language/reference/kind.md#field-rules).
