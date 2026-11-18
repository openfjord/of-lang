---
title: Links — The OpenFjord modeling language by example
---

# Links

> Part of [the language by example](README.md).

A link is a named relationship from a kind to a kind, in one direction.

<table>
<tr>
<td>

A link connects two kinds. This package has recipes and stations.

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
```

</td>
</tr>
<tr>
<td>

A link is a statement of its own. Its name is lowerCamelCase and has two
words or more. `from` is the kind where the link starts. `to` is the kind
where the link ends.

</td>
<td>

```yaml
madeAtStation:
  is: Link
  from: Recipe
  to: Station
```

</td>
</tr>
<tr>
<td>

`reads` is the sentence that product people see. It must contain `{from}`
and `{to}`.

</td>
<td>

```yaml
  reads: "{from} is made at {to}"
```

</td>
</tr>
<tr>
<td>

Every link needs a `definition`: one sentence that tells what the link
means.

</td>
<td>

```yaml
  definition: The station where cooks make the recipe.
```

</td>
</tr>
<tr>
<td>

To see the link at work, a test file states a station and a recipe. The
recipe gives the link under `links`, with the id of the station.

</td>
<td>

```yaml
# tests/examples/pesto-at-the-sauce-station.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Station, madeAtStation]

station/sauce:
  is: Fact
  kind: Station

recipe/pesto:
  is: Fact
  kind: Recipe
  links:
    madeAtStation: station/sauce
```

</td>
</tr>
<tr>
<td>

The test expects that the package rejects nothing: the link starts at a
recipe and ends at a station.

</td>
<td>

```yaml
pestoAtTheSauceStation:
  is: Test
  says: Cooks make the pesto at the sauce station.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test links
pass  test pestoAtTheSauceStation  (links/tests/examples/pesto-at-the-sauce-station.of.yaml)
1 passed, 0 failed
```

More about this: [links](../../language/reference/link.md) and [the rules for links](../../language/reference/link.md#rules).
