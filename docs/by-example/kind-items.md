---
title: Kind items — The OpenFjord modeling language by example
---

# Kind items

> Part of [the language by example](README.md).

A kind item in `match` gives a variable to the things of a kind.

<table>
<tr>
<td>

The package has two kinds. The first kind is `Dish`.

</td>
<td>

```yaml
# defs/kitchen.of.yaml
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

The second kind is `Station`.

</td>
<td>

```yaml
Station:
  is: Kind
  definition: A place in the kitchen where cooks make food.
  fields:
    name:
      type: string
```

</td>
</tr>
<tr>
<td>

A test file states two dishes and one station.

</td>
<td>

```yaml
# tests/examples/two-dishes-and-a-station.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Station]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue

station/pizza:
  is: Fact
  kind: Station
  fields:
    name: pizza
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test. This one expects that the package
rejects nothing.

</td>
<td>

```yaml
twoDishesAndAStation:
  is: Test
  says: Two dishes and a station are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes there are, given those facts.

</td>
<td>

```yaml
# tests/questions/which-dishes-are-there.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish]

whichDishesAreThere:
  is: Question
  asks: Which dishes are there?
  given: file:../examples/two-dishes-and-a-station.of.yaml
```

</td>
</tr>
<tr>
<td>

The kind item names the kind `Dish` and gives it the variable `$dish`. The
variable represents each dish in turn. Each variable has exactly one kind
item.

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

The answer is the two dishes. The station is not a dish, so `$dish` never
represents it.

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
$ of-lang test kind-items
pass  test twoDishesAndAStation  (kind-items/tests/examples/two-dishes-and-a-station.of.yaml)
pass  question whichDishesAreThere  (kind-items/tests/questions/which-dishes-are-there.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match).
