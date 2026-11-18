---
title: Any of — The OpenFjord modeling language by example
---

# Any of

> Part of [the language by example](README.md).

Each condition in `filter` must hold. `any of` lists conditions, and holds
when one or more of them hold.

<table>
<tr>
<td>

Each dish has a price and a status.

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
    price:
      type: integer
    status:
      enum: [active, retired]
```

</td>
</tr>
<tr>
<td>

A test file states three dishes. The lobster is expensive, and the fondue is
retired.

</td>
<td>

```yaml
# tests/examples/three-dishes.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12
    status: active

dish/lobster:
  is: Fact
  kind: Dish
  fields:
    name: lobster
    price: 45
    status: active

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    price: 28
    status: retired
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threeDishes:
  is: Test
  says: Active dishes and a retired dish are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes are retired or cost more than 30. Its `match`
finds each dish.

</td>
<td>

```yaml
# tests/questions/which-dishes-are-retired-or-expensive.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesAreRetiredOrExpensive:
  is: Question
  asks: Which dishes are retired or cost more than 30?
  given: file:../examples/three-dishes.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

`any of` is one condition of the `filter`. Its first condition holds for a
retired dish.

</td>
<td>

```yaml
  filter:
    - any of:
        - node: $dish
          field: status
          equals:
            value: retired
```

</td>
</tr>
<tr>
<td>

Its second condition holds for a dish that costs more than 30.

</td>
<td>

```yaml
        - node: $dish
          field: price
          is above:
            value: 30
```

</td>
</tr>
<tr>
<td>

One condition holds for the fondue, and the other holds for the lobster. No
condition holds for the margherita.

</td>
<td>

```yaml
  expect:
    $dish: [dish/lobster, dish/fondue]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test any-of
pass  test threeDishes  (any-of/tests/examples/three-dishes.of.yaml)
pass  question whichDishesAreRetiredOrExpensive  (any-of/tests/questions/which-dishes-are-retired-or-expensive.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
