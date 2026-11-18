---
title: Derive a kind — The OpenFjord modeling language by example
---

# Derive a kind

> Part of [the language by example](README.md).

A derivative can give a thing a narrower kind. The thing keeps its own kind
and also becomes a thing of the narrower kind.

<table>
<tr>
<td>

A recipe can use other recipes.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.

usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe.
```

</td>
</tr>
<tr>
<td>

A base recipe is a narrower sort of recipe: `BaseRecipe` extends `Recipe`.
No fact states that a recipe is a base recipe.

</td>
<td>

```yaml
BaseRecipe:
  is: Kind
  extends: Recipe
  definition: A recipe that the kitchen makes as a part of other recipes.
```

</td>
</tr>
<tr>
<td>

The derivative finds each recipe that some recipe uses. The end
`from: Recipe` gives a kind alone, because nothing else refers to the recipe
at that end.

</td>
<td>

```yaml
baseRecipeFromUse:
  is: Derivative
  says: A recipe that another recipe uses is a base recipe.
  match:
    - kind: { name: Recipe, as: $recipe }
    - link: usesRecipe
      from: Recipe
      to:   $recipe
```

</td>
</tr>
<tr>
<td>

An item with `node` and `becomes` adds a kind to the thing. The kind must
extend the kind of the variable.

</td>
<td>

```yaml
  derive:
    - node: $recipe
      becomes: BaseRecipe
```

</td>
</tr>
<tr>
<td>

A test file states two recipes. The margherita uses the tomato sauce.

</td>
<td>

```yaml
# tests/examples/tomato-sauce.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, BaseRecipe, usesRecipe]

recipe/tomato_sauce:
  is: Fact
  kind: Recipe

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

The test expects that the tomato sauce becomes a base recipe. No recipe
uses the margherita, so `notDerived` expects that the margherita does not.

</td>
<td>

```yaml
tomatoSauce:
  is: Test
  says: The margherita uses the tomato sauce, so only the tomato sauce is a base recipe.
  expect:
    derived:
      - node: recipe/tomato_sauce
        becomes: BaseRecipe
    notDerived:
      - node: recipe/margherita
        becomes: BaseRecipe
```

</td>
</tr>
</table>

`run` shows the kind that the derivative adds:

```console
$ of-lang run derive-kind derive-kind/tests/examples/tomato-sauce.of.yaml
derived:
  recipe/tomato_sauce is a BaseRecipe  (derivative baseRecipeFromUse)
rejected:
  nothing
```

More about this: [derive](../../language/reference/derivative.md#derive) and
[kinds](../../language/reference/kind.md).
