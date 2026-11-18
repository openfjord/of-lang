---
title: Expect derived facts — The OpenFjord modeling language by example
---

# Expect derived facts

> Part of [the language by example](README.md).

A test tells what a set of facts must derive. `derived` lists facts that
must hold, and `notDerived` lists facts that must not hold.

<table>
<tr>
<td>

The package needs something that derives a fact. Here, each `usesGarnish`
is also a `usesIngredient`.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food.

Ingredient:
  is: Kind
  definition: A raw food that the restaurant buys.

usesIngredient:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} uses {to}"
  definition: The recipe needs the ingredient.

usesGarnish:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} is garnished with {to}"
  definition: The cook puts the ingredient on the food at the end.
  under: usesIngredient
```

</td>
</tr>
<tr>
<td>

An example is a file of facts with one `Test` statement. The margherita
has basil as a garnish. The marinara uses oregano, but not as a garnish.

</td>
<td>

```yaml
# tests/examples/garnish-and-ingredient.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Ingredient, usesIngredient, usesGarnish]

ingredient/basil:
  is: Fact
  kind: Ingredient

ingredient/oregano:
  is: Fact
  kind: Ingredient

recipe/margherita:
  is: Fact
  kind: Recipe
  links:
    usesGarnish: ingredient/basil

recipe/marinara:
  is: Fact
  kind: Recipe
  links:
    usesIngredient: ingredient/oregano
```

</td>
</tr>
<tr>
<td>

A test must have `says` and `expect`. `says` tells in plain words what the
example shows.

</td>
<td>

```yaml
garnishAndIngredient:
  is: Test
  says: >
    Basil is a garnish on the margherita, so the margherita uses basil. The
    marinara uses oregano, but oregano is not a garnish on it.
  expect:
```

</td>
</tr>
<tr>
<td>

`derived` lists the facts that must hold after the package derives all that
it can. A link is written with `link`, `from` and `to`, and the two ends are
ids of facts.

</td>
<td>

```yaml
    derived:
      - link: usesIngredient
        from: recipe/margherita
        to:   ingredient/basil
```

</td>
</tr>
<tr>
<td>

`notDerived` lists the facts that must not hold. The test fails if one of
them holds.

</td>
<td>

```yaml
    notDerived:
      - link: usesGarnish
        from: recipe/marinara
        to:   ingredient/oregano
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test expect-derived
pass  test garnishAndIngredient  (expect-derived/tests/examples/garnish-and-ingredient.of.yaml)
1 passed, 0 failed
```

More about this: [examples](../../language/reference/test.md#examples) and [derived facts](../../language/reference/fact.md#derived-facts).
