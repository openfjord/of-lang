---
title: Questions — The OpenFjord modeling language by example
---

# Questions

> Part of [the language by example](README.md).

A question asks something about the facts of an example, and tells what the
answer must be.

<table>
<tr>
<td>

Cooks make each recipe at a station.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food.

Station:
  is: Kind
  definition: A place in the kitchen where cooks make recipes.

madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.
```

</td>
</tr>
<tr>
<td>

An example gives the facts: two stations and one recipe. No recipe is made
at the grill station.

</td>
<td>

```yaml
# tests/examples/two-stations.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Station, madeAtStation]

station/sauce:
  is: Fact
  kind: Station

station/grill:
  is: Fact
  kind: Station

recipe/pesto:
  is: Fact
  kind: Recipe
  links:
    madeAtStation: station/sauce

twoStations:
  is: Test
  says: Cooks make the pesto at the sauce station, and nothing at the grill station.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question goes in a file in `tests/questions/`. `asks` is the question in
plain words.

</td>
<td>

```yaml
# tests/questions/which-stations-make-a-recipe.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Station, madeAtStation]

whichStationsMakeARecipe:
  is: Question
  asks: Which stations make a recipe?
```

</td>
</tr>
<tr>
<td>

`given` names the example file with the facts to ask about.

</td>
<td>

```yaml
  given: file:../examples/two-stations.of.yaml
```

</td>
</tr>
<tr>
<td>

`match` tells what to find. A kind item gives a variable to the things of a
kind: `$station` represents each station in turn.

</td>
<td>

```yaml
  match:
    - kind: { name: Station, as: $station }
```

</td>
</tr>
<tr>
<td>

A link item is one link that must exist. Its `from` end is a kind alone,
so any recipe satisfies it.

</td>
<td>

```yaml
    - link: madeAtStation
      from: Recipe
      to:   $station
```

</td>
</tr>
<tr>
<td>

`expect` lists the ids that the variable must represent. The answer is the
sauce station only.

</td>
<td>

```yaml
  expect:
    $station: [station/sauce]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test questions
pass  test twoStations  (questions/tests/examples/two-stations.of.yaml)
pass  question whichStationsMakeARecipe  (questions/tests/questions/which-stations-make-a-recipe.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions).
