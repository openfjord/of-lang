---
title: Links under links — The OpenFjord modeling language by example
---

# Links under links

> Part of [the language by example](README.md).

A link can be `under` a broader link. Then every such link is also the
broader link.

<table>
<tr>
<td>

A recipe uses ingredients. `usesIngredient` is the broader link.

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
```

</td>
</tr>
<tr>
<td>

A garnish is one way to use an ingredient. `usesGarnish` is a narrower sort
of `usesIngredient`.

</td>
<td>

```yaml
usesGarnish:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} is garnished with {to}"
  definition: The cook puts the ingredient on the food at the end.
```

</td>
</tr>
<tr>
<td>

`under` names the broader link. The two ends of the narrower link must be
the kinds of the broader link, or kinds that extend them.

</td>
<td>

```yaml
  under: usesIngredient
```

</td>
</tr>
<tr>
<td>

A test file states that the margherita has basil as a garnish. It states
no `usesIngredient`.

</td>
<td>

```yaml
# tests/examples/basil-on-the-margherita.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Ingredient, usesIngredient, usesGarnish]

ingredient/basil:
  is: Fact
  kind: Ingredient

recipe/margherita:
  is: Fact
  kind: Recipe
  links:
    usesGarnish: ingredient/basil
```

</td>
</tr>
<tr>
<td>

The test expects the broader link as a derived fact: the margherita uses
basil.

</td>
<td>

```yaml
basilOnTheMargherita:
  is: Test
  says: Basil is a garnish on the margherita, so the margherita uses basil.
  expect:
    derived:
      - link: usesIngredient
        from: recipe/margherita
        to:   ingredient/basil
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test under
pass  test basilOnTheMargherita  (under/tests/examples/basil-on-the-margherita.of.yaml)
1 passed, 0 failed
```

More about this: [what links imply](../../language/reference/link.md#what-links-imply) and [the rules for links](../../language/reference/link.md#rules).
