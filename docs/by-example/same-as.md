---
title: Same as — The OpenFjord modeling language by example
---

# Same as

> Part of [the language by example](README.md).

`same as:` compares with a variable, or with a field of a variable, in place
of a literal.

<table>
<tr>
<td>

Each dish has a name and a price.

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
```

</td>
</tr>
<tr>
<td>

A test file states three dishes. Two of them have the same price.

</td>
<td>

```yaml
# tests/examples/two-equal-prices.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12

dish/marinara:
  is: Fact
  kind: Dish
  fields:
    name: marinara
    price: 12

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    price: 28
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
twoEqualPrices:
  is: Test
  says: Two dishes with the same price and a third dish are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes have the same price as another dish. It has two
variables for dishes. `match` finds each combination of two dishes.

</td>
<td>

```yaml
# tests/questions/which-dishes-share-a-price.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesShareAPrice:
  is: Question
  asks: Which dishes have the same price as another dish?
  given: file:../examples/two-equal-prices.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Dish, as: $other }
```

</td>
</tr>
<tr>
<td>

Two variables can represent the same thing. This condition prevents it:
`same as:` with a variable compares the two things.

</td>
<td>

```yaml
  filter:
    - node: $dish
      does not equal:
        same as: $other
```

</td>
</tr>
<tr>
<td>

`same as:` with a `node` and a `field` compares with the value of that
field. The condition holds when the two prices are equal.

</td>
<td>

```yaml
    - node: $dish
      field: price
      equals:
        same as:
          node: $other
          field: price
```

</td>
</tr>
<tr>
<td>

The margherita and the marinara cost the same. No other dish has the price
of the fondue.

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
$ of-lang test same-as
pass  test twoEqualPrices  (same-as/tests/examples/two-equal-prices.of.yaml)
pass  question whichDishesShareAPrice  (same-as/tests/questions/which-dishes-share-a-price.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
