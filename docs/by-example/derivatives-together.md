---
title: Derivatives together — The OpenFjord modeling language by example
---

# Derivatives together

> Part of [the language by example](README.md).

Derivatives run again and again, until a pass adds nothing. Thus a
derivative can match a link that a derivative adds.

<table>
<tr>
<td>

A recipe is made at a station, and a recipe can use another recipe. Facts
state these two links.

</td>
<td>

```yaml
# defs/station.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.

Station:
  is: Kind
  definition: A place in the kitchen where cooks make food.

madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.

usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe directly.
```

</td>
</tr>
<tr>
<td>

No fact states `needsStation`. Two derivatives add it.

</td>
<td>

```yaml
needsStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} needs {to}"
  definition: The kitchen cannot make the recipe when the station stops.
```

</td>
</tr>
<tr>
<td>

The first derivative matches only stated facts.

</td>
<td>

```yaml
stationFromRecipe:
  is: Derivative
  says: A recipe needs the station where cooks make it.
  match:
    - kind: { name: Recipe, as: $recipe }
    - kind: { name: Station, as: $station }
    - link: madeAtStation
      from: $recipe
      to:   $station
  derive:
    - link: needsStation
      from: $recipe
      to:   $station
```

</td>
</tr>
<tr>
<td>

The second derivative matches `needsStation`, the link that the first
derivative adds.

</td>
<td>

```yaml
stationFromPart:
  is: Derivative
  says: A recipe needs each station that one of its recipes needs.
  match:
    - kind: { name: Recipe, as: $recipe }
    - kind: { name: Recipe, as: $part }
    - kind: { name: Station, as: $station }
    - link: usesRecipe
      from: $recipe
      to:   $part
    - link: needsStation
      from: $part
      to:   $station
```

</td>
</tr>
<tr>
<td>

It also adds `needsStation`. Thus one pass can add a fact that the same
derivative matches in the pass after it.

</td>
<td>

```yaml
  derive:
    - link: needsStation
      from: $recipe
      to:   $station
```

</td>
</tr>
<tr>
<td>

A test file states a chain of three recipes. The margherita uses the
tomato sauce, the tomato sauce uses the stock, and cooks make the stock at the
stove.

</td>
<td>

```yaml
# tests/examples/stock-on-the-stove.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Recipe, Station, madeAtStation, usesRecipe, needsStation]

station/stove:
  is: Fact
  kind: Station

recipe/stock:
  is: Fact
  kind: Recipe
  links:
    madeAtStation: station/stove

recipe/tomato_sauce:
  is: Fact
  kind: Recipe
  links:
    usesRecipe: recipe/stock

recipe/margherita:
  is: Fact
  kind: Recipe
  links:
    usesRecipe: recipe/tomato_sauce
```

</td>
</tr>
<tr>
<td>

The test expects a link that no single pass can add. The tomato sauce must
need the stove before the derivative can find that the margherita needs it.

</td>
<td>

```yaml
stockOnTheStove:
  is: Test
  says: >
    Cooks make the stock at the stove. The margherita uses the tomato sauce,
    which uses the stock. Thus, the margherita needs the stove.
  expect:
    derived:
      - link: needsStation
        from: recipe/margherita
        to:   station/stove
```

</td>
</tr>
</table>

`explain` shows why a fact holds, down to the facts that the file states:

```console
$ of-lang explain derivatives-together derivatives-together/tests/examples/stock-on-the-stove.of.yaml recipe/margherita needsStation station/stove
recipe/margherita needs station/stove  (derivative stationFromPart)
  A recipe needs each station that one of its recipes needs.
  with $part = recipe/tomato_sauce, $recipe = recipe/margherita, $station = station/stove
  because:
    recipe/margherita uses recipe/tomato_sauce  (written at derivatives-together/tests/examples/stock-on-the-stove.of.yaml:25:5)
    recipe/tomato_sauce needs station/stove  (derivative stationFromPart)
      A recipe needs each station that one of its recipes needs.
      with $part = recipe/stock, $recipe = recipe/tomato_sauce, $station = station/stove
      because:
        recipe/tomato_sauce uses recipe/stock  (written at derivatives-together/tests/examples/stock-on-the-stove.of.yaml:19:5)
        recipe/stock needs station/stove  (derivative stationFromRecipe)
          A recipe needs the station where cooks make it.
          with $recipe = recipe/stock, $station = station/stove
          because:
            recipe/stock is made at station/stove  (written at derivatives-together/tests/examples/stock-on-the-stove.of.yaml:13:5)
```

More about this: [how derivatives run](../../language/reference/derivative.md#how-derivatives-run)
and [explanations](../../language/reference/derivative.md#explanations).
