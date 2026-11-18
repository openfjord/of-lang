---
title: Imports — The OpenFjord modeling language by example
---

# Imports

> Part of [the language by example](README.md).

A file can use a name only if the file defines that name or imports it. The
`import` block at the top of a file shows the source of each name.

<table>
<tr>
<td>

This file defines the kind `Station`. It uses no other name, so it has no
`import` block.

</td>
<td>

```yaml
# defs/station.of.yaml
Station:
  is: Kind
  definition: A part of the kitchen where a group of cooks makes some of the recipes.
```

</td>
</tr>
<tr>
<td>

A second file uses `Station`, so it imports the name. `from` gives the
source: `file:`, then a path relative to this file. `use` lists the names.

</td>
<td>

```yaml
# defs/recipe.of.yaml
import:
  - from: file:station.of.yaml
    use: [Station]
```

</td>
</tr>
<tr>
<td>

The file defines `Recipe` itself, so it does not import that name.

</td>
<td>

```yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
```

</td>
</tr>
<tr>
<td>

The link uses the two names: one that the file defines and one that it
imports.

</td>
<td>

```yaml
madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where the cooks make the recipe.
```

</td>
</tr>
<tr>
<td>

A test file imports each name that it uses, with one item for each source
file. `../` goes to the parent folder.

</td>
<td>

```yaml
# tests/examples/sauce-at-a-station.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Station]
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, madeAtStation]
```

</td>
</tr>
<tr>
<td>

The facts use the two kinds and the link, so the file imports all three
names.

</td>
<td>

```yaml
station/sauce:
  is: Fact
  kind: Station

recipe/tomato_sauce:
  is: Fact
  kind: Recipe
  links:
    madeAtStation: station/sauce

sauceAtAStation:
  is: Test
  says: A recipe that cooks make at a station is valid.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` reads the files with their imports and runs the test:

```console
$ of-lang test imports
pass  test sauceAtAStation  (imports/tests/examples/sauce-at-a-station.of.yaml)
1 passed, 0 failed
```

More about this: [imports](../../language/reference/files.md#imports).
