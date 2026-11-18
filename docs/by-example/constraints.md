---
title: Constraints — The OpenFjord modeling language by example
---

# Constraints

> Part of [the language by example](README.md).

A constraint rejects facts. It finds things, and it reports each match as a
violation.

<table>
<tr>
<td>

Each dish has a price and a cost.

</td>
<td>

```yaml
# defs/dish.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    price:
      type: integer
    cost:
      type: integer
```

</td>
</tr>
<tr>
<td>

A constraint has a `says` sentence that gives the rule, in plain words.

</td>
<td>

```yaml
dishSellsAboveCost:
  is: Constraint
  says: The price of a dish is above its cost.
```

</td>
</tr>
<tr>
<td>

`match` lists the things to find. Here the variable `$dish` represents
each dish in turn.

</td>
<td>

```yaml
  match:
    - kind: { name: Dish, as: $dish }
```

</td>
</tr>
<tr>
<td>

`filter` keeps the matches for which each condition holds. This condition
holds when the price of the dish is its cost or less. Each match that stays
is a violation, so `match` and `filter` describe what breaks the rule.

</td>
<td>

```yaml
  filter:
    - node:  $dish
      field: price
      is at most:
        same as:
          node:  $dish
          field: cost
```

</td>
</tr>
<tr>
<td>

`message` is what a violation says. `{$dish}` gives the id of the dish.

</td>
<td>

```yaml
  message: "{$dish} must have a price above its cost."
```

</td>
</tr>
<tr>
<td>

A test file states two dishes. The price of the fondue is below its
cost.

</td>
<td>

```yaml
# tests/examples/cheap-fondue.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish, dishSellsAboveCost]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    price: 12
    cost: 4

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    price: 9
    cost: 11
```

</td>
</tr>
<tr>
<td>

`rejected` lists all constraints that the facts must break. The file
imports each constraint that the test names.

</td>
<td>

```yaml
cheapFondue:
  is: Test
  says: The fondue costs more than its price, so the constraint rejects it.
  expect:
    rejected: [dishSellsAboveCost]
```

</td>
</tr>
</table>

`run` reports each violation with the name of the constraint and its
message:

```console
$ of-lang run constraints constraints/tests/examples/cheap-fondue.of.yaml
derived:
  nothing
rejected:
  dishSellsAboveCost: dish/fondue must have a price above its cost.
```

More about this: [constraints](../../language/reference/constraint.md) and
[how constraints run](../../language/reference/constraint.md#how-constraints-run).
