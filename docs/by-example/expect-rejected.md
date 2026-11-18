---
title: Expect rejected facts — The OpenFjord modeling language by example
---

# Expect rejected facts

> Part of [the language by example](README.md).

A constraint rejects facts. A test lists under `rejected` all constraints
that its facts must break.

<table>
<tr>
<td>

A dish has a price.

</td>
<td>

```yaml
# defs/dish.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    price:
      type: number
```

</td>
</tr>
<tr>
<td>

A constraint finds things and reports each match as a violation. This one
finds each dish with a price of 0.

</td>
<td>

```yaml
dishIsNotFree:
  is: Constraint
  says: A dish must not be free.
  match:
    - kind: { name: Dish, as: $dish }
  filter:
    - node: $dish
      field: price
      equals:
        value: 0
  message: "{$dish} must have a price above 0."
```

</td>
</tr>
<tr>
<td>

A test file imports each constraint that its test names. It states two
dishes, and the bread is free.

</td>
<td>

```yaml
# tests/examples/free-bread.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish, dishIsNotFree]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    price: 12

dish/bread:
  is: Fact
  kind: Dish
  fields:
    price: 0
```

</td>
</tr>
<tr>
<td>

`rejected` lists the constraints that the facts must break. The test
passes only if the facts break exactly these constraints. Thus, a test
with `rejected: []` expects that the facts break none.

</td>
<td>

```yaml
freeBread:
  is: Test
  says: The bread has a price of 0, so the package rejects it.
  expect:
    rejected: [dishIsNotFree]
```

</td>
</tr>
</table>

`test` runs the test. The test passes because the constraint rejects the bread:

```console
$ of-lang test expect-rejected
pass  test freeBread  (expect-rejected/tests/examples/free-bread.of.yaml)
1 passed, 0 failed
```

More about this: [examples](../../language/reference/test.md#examples) and [constraints](../../language/reference/constraint.md).
