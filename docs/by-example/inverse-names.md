---
title: Inverse names — The OpenFjord modeling language by example
---

# Inverse names

> Part of [the language by example](README.md).

An inverse name is a name for a link in the opposite direction. It adds no
link: it is the same link, read the other way.

<table>
<tr>
<td>

The link `madeAtStation` goes from a recipe to a station.

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

`inverse` gives the name `makesRecipe` to the opposite direction, from a
station to a recipe. An inverse name follows the same rules as the name of
a link.

</td>
<td>

```yaml
  inverse: makesRecipe
```

</td>
</tr>
<tr>
<td>

A test file states two stations and two recipes. The pesto gives the link
by its own name.

</td>
<td>

```yaml
# tests/examples/two-stations.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Station, madeAtStation, makesRecipe]

station/sauce:
  is: Fact
  kind: Station

recipe/pesto:
  is: Fact
  kind: Recipe
  links:
    madeAtStation: station/sauce

recipe/margherita:
  is: Fact
  kind: Recipe
```

</td>
</tr>
<tr>
<td>

A fact can use the inverse name. This link on the pizza station states the
same fact as `madeAtStation: station/pizza` on the margherita.

</td>
<td>

```yaml
station/pizza:
  is: Fact
  kind: Station
  links:
    makesRecipe: recipe/margherita
```

</td>
</tr>
<tr>
<td>

The test expects that the link holds from the margherita to the pizza
station, although no fact states it in that direction.

</td>
<td>

```yaml
twoStations:
  is: Test
  says: Each of the two stations makes one recipe.
  expect:
    derived:
      - link: madeAtStation
        from: recipe/margherita
        to:   station/pizza
```

</td>
</tr>
<tr>
<td>

A question can also use the inverse name. This one starts at the station.

</td>
<td>

```yaml
# tests/questions/which-stations-make-a-recipe.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Station, makesRecipe]

whichStationsMakeARecipe:
  is: Question
  asks: Which stations make a recipe?
  given: file:../examples/two-stations.of.yaml
  match:
    - kind: { name: Station, as: $station }
    - link: makesRecipe
      from: $station
      to:   Recipe
```

</td>
</tr>
<tr>
<td>

The answer has the two stations. It does not matter which name the fact
uses.

</td>
<td>

```yaml
  expect:
    $station: [station/sauce, station/pizza]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test inverse-names
pass  test twoStations  (inverse-names/tests/examples/two-stations.of.yaml)
pass  question whichStationsMakeARecipe  (inverse-names/tests/questions/which-stations-make-a-recipe.of.yaml)
2 passed, 0 failed
```

More about this: [what links imply](../../language/reference/link.md#what-links-imply) and [the rules for facts](../../language/reference/fact.md#rules).
