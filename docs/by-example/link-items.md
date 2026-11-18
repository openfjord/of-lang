---
title: Link items — The OpenFjord modeling language by example
---

# Link items

> Part of [the language by example](README.md).

A link item in `match` is one link that must exist between two things.

<table>
<tr>
<td>

The package has dishes and stations.

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

The link `madeAtStation` goes from a dish to a station.

</td>
<td>

```yaml
madeAtStation:
  is: Link
  from: Dish
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the dish.
```

</td>
</tr>
<tr>
<td>

A test file states two stations.

</td>
<td>

```yaml
# tests/examples/dishes-at-stations.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Station, madeAtStation]

station/pizza:
  is: Fact
  kind: Station
  fields:
    name: pizza

station/grill:
  is: Fact
  kind: Station
  fields:
    name: grill
```

</td>
</tr>
<tr>
<td>

Cooks make the margherita at the pizza station. The fondue has no station.

</td>
<td>

```yaml
dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
  links:
    madeAtStation: station/pizza

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    name: fondue
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
dishesAtStations:
  is: Test
  says: A dish with a station and a dish without a station are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which dishes cooks make at which stations. Each variable has
a kind item.

</td>
<td>

```yaml
# tests/questions/which-dishes-are-made-where.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Station, madeAtStation]

whichDishesAreMadeWhere:
  is: Question
  asks: Which dishes do cooks make at which stations?
  given: file:../examples/dishes-at-stations.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Station, as: $station }
```

</td>
</tr>
<tr>
<td>

The link item has `link`, `from` and `to`. `match` finds each combination of
a dish and a station for which the `madeAtStation` link exists.

</td>
<td>

```yaml
    - link: madeAtStation
      from: $dish
      to:   $station
```

</td>
</tr>
<tr>
<td>

There is one match: the margherita and the pizza station. The fondue and the
grill station have no such link, so they are not in the answer.

</td>
<td>

```yaml
  expect:
    $dish: [dish/margherita]
    $station: [station/pizza]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test link-items
pass  test dishesAtStations  (link-items/tests/examples/dishes-at-stations.of.yaml)
pass  question whichDishesAreMadeWhere  (link-items/tests/questions/which-dishes-are-made-where.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match).
