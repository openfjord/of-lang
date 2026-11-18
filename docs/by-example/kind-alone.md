---
title: A kind alone — The OpenFjord modeling language by example
---

# A kind alone

> Part of [the language by example](README.md).

An end of a link item can be a kind alone, without a variable. This is for a
thing that nothing else refers to.

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

A test file states two stations and one dish. Cooks make the dish at the
pizza station.

</td>
<td>

```yaml
# tests/examples/one-dish-two-stations.of.yaml
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

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
  links:
    madeAtStation: station/pizza
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
oneDishTwoStations:
  is: Test
  says: A station that makes a dish and a station that makes no dish are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which stations make a dish. The stations have the variable
`$station`.

</td>
<td>

```yaml
# tests/questions/which-stations-make-a-dish.of.yaml
import:
  - from: file:../../defs/kitchen.of.yaml
    use: [Dish, Station, madeAtStation]

whichStationsMakeADish:
  is: Question
  asks: Which stations make a dish?
  given: file:../examples/one-dish-two-stations.of.yaml
  match:
    - kind: { name: Station, as: $station }
```

</td>
</tr>
<tr>
<td>

The question does not refer to the dish again, so the `from` end gives the
kind alone. This end only has to exist. A variable that only one link end
uses is an error.

</td>
<td>

```yaml
    - link: madeAtStation
      from: Dish
      to:   $station
```

</td>
</tr>
<tr>
<td>

The answer is the pizza station. No dish has a link to the grill station.

</td>
<td>

```yaml
  expect:
    $station: [station/pizza]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test kind-alone
pass  test oneDishTwoStations  (kind-alone/tests/examples/one-dish-two-stations.of.yaml)
pass  question whichStationsMakeADish  (kind-alone/tests/questions/which-stations-make-a-dish.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match).
