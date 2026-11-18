---
title: Links — Introduction to the OpenFjord modeling language
---

# Links

> Part of [the introduction](README.md).

The kitchen has five kinds: stations, recipes, dishes, ingredients and
allergens. A field holds a value of the thing itself, such as the name of a
recipe. But the kitchen also must tell how the things relate: which station
makes a recipe, and what a recipe uses. A link does this. A link is a named
relationship from one kind to a different kind, or to the same kind, in one
direction.

## A link between two kinds

Cooks make each recipe at a station:

```yaml
# defs/recipe.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name, Minutes, Price, URL]
  - from: file:station.of.yaml
    use: [Station]
  - from: file:ingredient.of.yaml
    use: [Ingredient]

madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.
  inverse: makesRecipe
```

A link is a statement of its own. It is not a part of a kind. Usually, it is
in the file of the kind where it starts.

- `from` is the kind where the link starts, and `to` is the kind where it
  ends. `Station` is in a different file, so the import block gets one more
  entry.
- `reads` is the sentence that people read for one such link: "the pesto is
  made at the sauce station". It must contain `{from}` and `{to}`.
- `definition` tells what the link means.
- `inverse` is optional. It gives a name to the link in the opposite
  direction. `makesRecipe` from a station to a recipe is the same
  relationship as `madeAtStation` from the recipe to the station.

The name of a link starts with a lowercase letter and has two words or more.
The words tell what the link connects and how: `madeAtStation`, not `at`.

A link has two ends and nothing more. It has no fields.

`Dish` extends `Recipe`, so a dish can have each link that starts at
`Recipe`. The kitchen does not declare the link a second time for dishes.

## A link that chains

A recipe can use a different recipe. The margherita uses the tomato sauce:

```yaml
# defs/recipe.of.yaml
usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe, directly or through other recipes.
  chains: true
  inverse: usedByRecipe
```

This link starts and ends at the same kind. `chains: true` tells that a
chain of the link is also the link. When the lasagne uses the ragù, and the
ragù uses the tomato sauce, then the lasagne uses the tomato sauce. Nobody
has to write the third link.

## A link under a link

A recipe also uses ingredients. Some ingredients are a garnish: the cook puts
them on the food at the end.

```yaml
# defs/recipe.of.yaml
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
  definition: The cook puts the ingredient on the food at the end, to make it look or taste better.
  under: usesIngredient
```

`under` names a broader link. `usesGarnish` is a narrower sort of
`usesIngredient`: each garnish of a recipe is also an ingredient of that
recipe. `under` is for links what `extends` is for kinds.

## Two more links

An ingredient can contain an allergen. The kitchen puts this link in the
file of `Ingredient`, where it starts:

```yaml
# defs/ingredient.of.yaml
containsAllergen:
  is: Link
  from: Ingredient
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is a part of the ingredient.
```

The menu gives warnings about allergens. This link is the first definition
about a topic, allergens, that concerns more than one kind. The kitchen
starts a file for the topic:

```yaml
# defs/allergens.of.yaml
import:
  - from: file:recipe.of.yaml
    use: [Dish]
  - from: file:ingredient.of.yaml
    use: [Allergen]

listsAllergen:
  is: Link
  from: Dish
  to: Allergen
  reads: "the menu tells guests that {from} contains {to}"
  definition: The menu gives a warning about the allergen next to the dish.
```

This link starts at `Dish`, not at `Recipe`. Only a dish is on the menu, so
only a dish can have the link.

## Check the package

```console
$ of-lang check 02
ok
```

`check` makes sure that each `from` and each `to` is a kind that the file
defines or imports. Without the import of `Station`, the package breaks a
rule:

```console
$ of-lang check 02
02/defs/recipe.of.yaml:43:7: Station is not defined in this file or imported (FL9)
```

The package tells which links can exist. It does not yet tell about one
real dish. The next chapter does: [Facts and tests](03-facts-and-tests.md).

More about this: [links](../../language/reference/link.md),
[the links of a thing](../../language/reference/kind.md#links-of-a-thing) and
[imports](../../language/reference/files.md#imports) in the reference. In the
language by example: [Links](../by-example/links.md),
[Inverse names](../by-example/inverse-names.md),
[Links that chain](../by-example/chains.md) and
[Links under links](../by-example/under.md).
