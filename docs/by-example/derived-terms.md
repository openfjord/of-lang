---
title: Derived terms — The OpenFjord modeling language by example
---

# Derived terms

> Part of [the language by example](README.md).

The definitions derive information before there are any facts: what
follows for kinds and links. A test with no facts can expect this
information.

<table>
<tr>
<td>

A dish extends a recipe.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food.

Dish:
  is: Kind
  extends: Recipe
  definition: A recipe that the restaurant sells to guests.

Station:
  is: Kind
  definition: A place in the kitchen where cooks make recipes.
```

</td>
</tr>
<tr>
<td>

One link has an inverse name, and one link chains.

</td>
<td>

```yaml
madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.
  inverse: makesRecipe

usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe, directly or through other recipes.
  chains: true
```

</td>
</tr>
<tr>
<td>

One link is under another link.

</td>
<td>

```yaml
usesSauce:
  is: Link
  from: Dish
  to: Recipe
  reads: "{from} is served with {to}"
  definition: The cook puts the second recipe on the dish as its sauce.
  under: usesRecipe
```

</td>
</tr>
<tr>
<td>

A test in `tests/terms/` has no facts. It imports each kind and each link
that it names.

</td>
<td>

```yaml
# tests/terms/dish.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish, madeAtStation, makesRecipe, usesRecipe, usesSauce]

dishIsARecipe:
  is: Test
  says: >
    A dish is a recipe, so a dish can have a station. A sauce is one of the
    recipes that a dish uses.
  expect:
    derived:
```

</td>
</tr>
<tr>
<td>

A term about a kind starts with `kind`. `is a` comes from `extends`. `can
have` tells that a thing of the kind can have the link: a dish gets
`madeAtStation` from `Recipe`.

</td>
<td>

```yaml
      - kind: Dish
        is a: Recipe
      - kind: Dish
        can have: madeAtStation
```

</td>
</tr>
<tr>
<td>

A term about a link starts with `link`. `implies` comes from `under`,
`chains` comes from `chains`, and `inverse of` comes from the inverse
name.

</td>
<td>

```yaml
      - link: usesSauce
        implies: usesRecipe
      - link: usesRecipe
        chains: true
      - link: makesRecipe
        inverse of: madeAtStation
```

</td>
</tr>
<tr>
<td>

`notDerived` lists terms that must not hold. A recipe is not always a
dish, and a recipe that a dish uses is not always its sauce.

</td>
<td>

```yaml
    notDerived:
      - kind: Recipe
        is a: Dish
      - link: usesRecipe
        implies: usesSauce
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test derived-terms
pass  test dishIsARecipe  (derived-terms/tests/terms/dish.of.yaml)
1 passed, 0 failed
```

More about this: [derived terms](../../language/reference/test.md#derived-terms).
