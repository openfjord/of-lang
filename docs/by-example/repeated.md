---
title: Repeated — The OpenFjord modeling language by example
---

# Repeated

> Part of [the language by example](README.md).

`repeated: true` on a link item follows the link for one or more steps. It
is for a link that does not chain.

<table>
<tr>
<td>

The package has stations.

</td>
<td>

```yaml
# defs/station.of.yaml
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

A station can depend on another station. The link has no `chains: true`, so
a chain of the link does not imply the link.

</td>
<td>

```yaml
dependsOnStation:
  is: Link
  from: Station
  to: Station
  reads: "{from} depends on {to}"
  definition: The first station cannot work when the second station stops.
```

</td>
</tr>
<tr>
<td>

A test file states three stations in a sequence. The pizza station depends
on the sauce station, and the sauce station depends on the prep station.

</td>
<td>

```yaml
# tests/examples/three-stations.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Station, dependsOnStation]

station/pizza:
  is: Fact
  kind: Station
  fields:
    name: pizza
  links:
    dependsOnStation: station/sauce

station/sauce:
  is: Fact
  kind: Station
  fields:
    name: sauce
  links:
    dependsOnStation: station/prep

station/prep:
  is: Fact
  kind: Station
  fields:
    name: prep
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threeStations:
  is: Test
  says: Three stations that depend on each other in a sequence are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which stations the pizza station depends on, directly or
through other stations.

</td>
<td>

```yaml
# tests/questions/which-stations-does-pizza-depend-on.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Station, dependsOnStation]

whichStationsDoesPizzaDependOn:
  is: Question
  asks: Which stations does the pizza station depend on, directly or through other stations?
  given: file:../examples/three-stations.of.yaml
  match:
    - kind: { name: Station, as: $station }
    - kind: { name: Station, as: $other }
```

</td>
</tr>
<tr>
<td>

With `repeated: true`, the link item matches the link one or more times in a
sequence. The stations between the two ends have no names.

</td>
<td>

```yaml
    - link: dependsOnStation
      from: $station
      to:   $other
      repeated: true
```

</td>
</tr>
<tr>
<td>

A condition keeps the matches that start at the pizza station.

</td>
<td>

```yaml
  filter:
    - node: $station
      equals:
        value: station/pizza
```

</td>
</tr>
<tr>
<td>

The answer has the sauce station, and also the prep station, which is two
steps away. Without `repeated: true`, the answer has only the sauce station.

</td>
<td>

```yaml
  expect:
    $other: [station/sauce, station/prep]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test repeated
pass  test threeStations  (repeated/tests/examples/three-stations.of.yaml)
pass  question whichStationsDoesPizzaDependOn  (repeated/tests/questions/which-stations-does-pizza-depend-on.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match), [what links
imply](../../language/reference/link.md#what-links-imply).
