---
title: Links that chain — The OpenFjord modeling language by example
---

# Links that chain

> Part of [the language by example](README.md).

A link with `chains: true` follows through: when a first thing has the link
to a second thing, and the second thing has it to a third thing, the first
thing has it to the third thing.

<table>
<tr>
<td>

A recipe can use another recipe. The link starts and ends at the same
kind.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food.

usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe, directly or through other recipes.
```

</td>
</tr>
<tr>
<td>

`chains: true` tells that a chain of the link implies the link. In a link
that chains, `to` must be the same kind as `from`, or a kind that extends
it.

</td>
<td>

```yaml
  chains: true
```

</td>
</tr>
<tr>
<td>

A test file states three recipes. The lasagne uses the ragu, and the ragu
uses the stock.

</td>
<td>

```yaml
# tests/examples/stock-in-the-lasagne.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, usesRecipe]

recipe/stock:
  is: Fact
  kind: Recipe

recipe/ragu:
  is: Fact
  kind: Recipe
  links:
    usesRecipe: recipe/stock

recipe/lasagne:
  is: Fact
  kind: Recipe
  links:
    usesRecipe: recipe/ragu
```

</td>
</tr>
<tr>
<td>

No fact states that the lasagne uses the stock. The test expects it as a
derived fact.

</td>
<td>

```yaml
stockInTheLasagne:
  is: Test
  says: The lasagne uses the ragu, which uses the stock. Thus, the lasagne uses the stock.
  expect:
    derived:
      - link: usesRecipe
        from: recipe/lasagne
        to:   recipe/stock
```

</td>
</tr>
<tr>
<td>

The link stays in one direction: the stock does not use the lasagne.

</td>
<td>

```yaml
    notDerived:
      - link: usesRecipe
        from: recipe/stock
        to:   recipe/lasagne
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test chains
pass  test stockInTheLasagne  (chains/tests/examples/stock-in-the-lasagne.of.yaml)
1 passed, 0 failed
```

More about this: [what links imply](../../language/reference/link.md#what-links-imply) and [the rules for links](../../language/reference/link.md#rules).
