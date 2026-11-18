---
title: Unless — The OpenFjord modeling language by example
---

# Unless

> Part of [the language by example](README.md).

`unless` gives the conditions that excuse a match of a constraint. A match
that satisfies them is not a violation.

<table>
<tr>
<td>

Each dish states how many minutes the cooks need for it, and whether
guests order it in advance.

</td>
<td>

```yaml
# defs/dish.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    prepMinutes:
      type: integer
    orderedInAdvance:
      type: boolean
```

</td>
</tr>
<tr>
<td>

`match` and `filter` select the matches that the constraint applies to:
each dish that takes more than 30 minutes.

</td>
<td>

```yaml
dishIsReadyInTime:
  is: Constraint
  says: >
    Cooks make a dish in 30 minutes or less, unless guests order it in
    advance.
  match:
    - kind: { name: Dish, as: $dish }
  filter:
    - node:  $dish
      field: prepMinutes
      is above:
        value: 30
```

</td>
</tr>
<tr>
<td>

`unless` tells which of these matches to excuse. Its conditions have the
same form as the conditions under `filter`. It excuses a match only when all
its conditions hold.

</td>
<td>

```yaml
  unless:
    - node:  $dish
      field: orderedInAdvance
      equals:
        value: true
```

</td>
</tr>
<tr>
<td>

The message is for the matches that `unless` does not excuse.

</td>
<td>

```yaml
  message: "{$dish} must take 30 minutes or less, or guests must order it in advance."
```

</td>
</tr>
<tr>
<td>

A test file states two dishes that take more than 30 minutes. Guests
order the roast in advance.

</td>
<td>

```yaml
# tests/examples/two-slow-dishes.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish, dishIsReadyInTime]

dish/roast:
  is: Fact
  kind: Dish
  fields:
    prepMinutes: 90
    orderedInAdvance: true
```

</td>
</tr>
<tr>
<td>

Guests do not order the risotto in advance.

</td>
<td>

```yaml
dish/risotto:
  is: Fact
  kind: Dish
  fields:
    prepMinutes: 40
    orderedInAdvance: false
```

</td>
</tr>
<tr>
<td>

The facts break the constraint, because the risotto is a violation.

</td>
<td>

```yaml
twoSlowDishes:
  is: Test
  says: >
    The roast and the risotto both take more than 30 minutes. Guests order
    the roast in advance, so the constraint rejects only the risotto.
  expect:
    rejected: [dishIsReadyInTime]
```

</td>
</tr>
</table>

`run` reports one violation. The roast is a match, but `unless` excuses it:

```console
$ of-lang run unless unless/tests/examples/two-slow-dishes.of.yaml
derived:
  nothing
rejected:
  dishIsReadyInTime: dish/risotto must take 30 minutes or less, or guests must order it in advance.
```

More about this: [filter and unless](../../language/reference/constraint.md#filter-and-unless).
