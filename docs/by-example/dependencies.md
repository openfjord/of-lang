---
title: Dependencies — The OpenFjord modeling language by example
---

# Dependencies

> Part of [the language by example](README.md).

A package can use the definitions of a different package. The manifest
names the other package in `dependsOn`, and a file imports a name from it
with `pkg:`.

<table>
<tr>
<td>

The manifest starts with the four required keys.

</td>
<td>

```yaml
# package.yaml
package: "@by-example/dependencies"
version: 0.1.0
owner: kitchen
rdfNamespace: dependencies
```

</td>
</tr>
<tr>
<td>

`dependsOn` gives the full name of each package that this package depends
on, and a range of versions. `^0.1` accepts each version 0.1 of
`@by-example/open-for`, the package of the page [Open for](open-for.md).

</td>
<td>

```yaml
dependsOn:
  "@by-example/open-for": ^0.1
```

</td>
</tr>
<tr>
<td>

A `pkg:` import gives the full name of the other package. The file can use
only the names that the other package exports. `Recipe` is one of them.

</td>
<td>

```yaml
# defs/dish.of.yaml
import:
  - from: pkg:@by-example/open-for
    use: [Recipe]
```

</td>
</tr>
<tr>
<td>

The file then uses `Recipe` as it uses a name of its own package. The other
package lets this package extend `Recipe`.

</td>
<td>

```yaml
Dish:
  is: Kind
  extends: Recipe
  definition: A recipe that the restaurant sells to guests.
```

</td>
</tr>
<tr>
<td>

A test file also imports the name from the other package. The test
expects that a dish is a recipe.

</td>
<td>

```yaml
# tests/terms/dish-is-a-recipe.of.yaml
import:
  - from: pkg:@by-example/open-for
    use: [Recipe]
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dishIsARecipe:
  is: Test
  says: A dish is a recipe.
  expect:
    derived:
      - kind: Dish
        is a: Recipe
```

</td>
</tr>
</table>

The tool looks for each dependency in the folder that holds the package. The
folder `open-for` is next to the folder `dependencies`, so `test` finds it:

```console
$ of-lang test dependencies
pass  test dishIsARecipe  (dependencies/tests/terms/dish-is-a-recipe.of.yaml)
1 passed, 0 failed
```

More about this: [dependencies](../../language/reference/package.md#dependencies),
[exports](../../language/reference/package.md#exports) and
[imports](../../language/reference/files.md#imports).
