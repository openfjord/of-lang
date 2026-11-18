---
title: Open for — The OpenFjord modeling language by example
---

# Open for

> Part of [the language by example](README.md).

A different package can extend a kind of this package, or declare a link
under a link of this package, only where the manifest opens it in `openFor`.

<table>
<tr>
<td>

The package defines a kind for recipes and a kind for ingredients.

</td>
<td>

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.

Ingredient:
  is: Kind
  definition: A food that the kitchen buys and does not make.
```

</td>
</tr>
<tr>
<td>

It also defines a link between them.

</td>
<td>

```yaml
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

The manifest starts with the four required keys.

</td>
<td>

```yaml
# package.yaml
package: "@by-example/open-for"
version: 0.1.0
owner: kitchen
rdfNamespace: open-for
```

</td>
</tr>
<tr>
<td>

`exports` lists the definitions that other packages can import.

</td>
<td>

```yaml
exports:
  kinds: [Recipe, Ingredient]
  links: [usesIngredient]
```

</td>
</tr>
<tr>
<td>

`openFor` lists the extension points. `extend` identifies a kind: a
different package can declare a kind that extends `Recipe`.

</td>
<td>

```yaml
openFor:
  - extend: Recipe
```

</td>
</tr>
<tr>
<td>

`under` identifies a link: a different package can declare a link under
`usesIngredient`. `Ingredient` is not in the list, so no other package can
extend it.

</td>
<td>

```yaml
  - under: usesIngredient
```

</td>
</tr>
</table>

The package of the page [Dependencies](dependencies.md) depends on this
package and declares `Dish`, a kind that extends `Recipe`. `check` accepts
the two packages:

```console
$ of-lang check open-for dependencies
ok
```

More about this: [the manifest](../../language/reference/package.md#the-manifest)
and [exports](../../language/reference/package.md#exports).
