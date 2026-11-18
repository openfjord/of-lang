---
title: Filter — The OpenFjord modeling language by example
---

# Filter

> Part of [the language by example](README.md).

`filter` holds conditions. It keeps the matches for which each condition
holds.

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

A test file states three dishes with different prices.

</td>
<td>

```yaml
# tests/examples/three-prices.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
    price: 28

dish/lobster:
  is: Fact
  kind: Dish
  fields:
    name: lobster
    price: 45
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threePrices:
  is: Test
  says: Three dishes with different prices are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks about the prices. Its `match` finds each dish.

</td>
<td>

```yaml
# tests/questions/which-dishes-cost-12-or-more-but-less-than-45.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

whichDishesCost12OrMoreButLessThan45:
  is: Question
  asks: Which dishes cost 12 or more, but less than 45?
  given: file:../examples/three-prices.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

A condition names a `node` and a `field`. Then it gives a comparison word,
and `value:` with a literal to compare with. `is at least` holds when the
price is 12 or more.

</td>
<td>

```yaml
  filter:
    - node: $dish
      field: price
      is at least:
        value: 12
```

</td>
</tr>
<tr>
<td>

`is below` holds when the price is less than 45. The other comparison words
that take `value:` are `equals`, `does not equal`, `is above` and `is at
most`.

</td>
<td>

```yaml
    - node: $dish
      field: price
      is below:
        value: 45
```

</td>
</tr>
<tr>
<td>

The two conditions hold for the margherita and the fondue. The lobster costs
45, so the second condition does not hold for it.

</td>
<td>

```yaml
  expect:
    $dish: [dish/margherita, dish/fondue]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test filter
pass  test threePrices  (filter/tests/examples/three-prices.of.yaml)
pass  question whichDishesCost12OrMoreButLessThan45  (filter/tests/questions/which-dishes-cost-12-or-more-but-less-than-45.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
