---
title: Packages — Introduction to the OpenFjord modeling language
---

# Packages

> Part of [the introduction](README.md).

The package of the kitchen is a model that works. But a package is also the
unit that people own, version and depend on. The kitchen is one team in the
restaurant. The purchasing team has a package of its own, an app for guests
wants a part of the kitchen's model, and a reader wants to know the reasons
for the rules. This chapter adds what the package needs for them. All of it
is in the manifest, in one more file of definitions, and in a `docs/`
folder.

## A dependency

The purchasing team owns the suppliers of the restaurant. Its package is
[`@bistro/suppliers`](../../examples/suppliers/), with the kinds `Supplier`
and `Ingredient`. The kitchen wants to tell which supplier sells each of its
ingredients. It does not define suppliers again. It depends on the package
that defines them:

```yaml
# package.yaml
dependsOn:
  "@bistro/suppliers": ^1.2
```

`dependsOn` gives, for each package, the range of versions that the kitchen
accepts. `^1.2` is version 1.2 or later, but below 2.0.

A file gets names from a different package with a `pkg:` import:

```yaml
# defs/buying.of.yaml
import:
  - from: file:ingredient.of.yaml
    use: [Ingredient]
  - from: pkg:@bistro/suppliers
    use: [Supplier, { name: Ingredient, as: SupplierIngredient }]

boughtFrom:
  is: Link
  from: Ingredient
  to: Supplier
  reads: "the kitchen buys {from} from {to}"
  definition: The supplier that sells the ingredient to the kitchen.

boughtAs:
  is: Link
  from: Ingredient
  to: SupplierIngredient
  reads: "the kitchen buys {from} as {to}"
  definition: The product of a supplier that the kitchen buys to have the ingredient.
```

The two packages each have a kind with the name `Ingredient`. For the
kitchen, it is a food. For purchasing, it is a product in the catalog of a
supplier. Two names in one file must be different, so the import gives the
second one a different name in this file, with `as`.

The links of the kitchen can end at the kinds of the other package. The
kitchen cannot change those kinds: they belong to purchasing.

The tool must find each package in `dependsOn`. By default, it looks in the
folder that holds the package. `--packages` names a different folder. In the
repository, the suppliers package is in `examples/`:

```console
$ of-lang check --packages ../../../examples step-06
ok
```

## Exports

A package decides what other packages can import from it. The manifest of
`@bistro/suppliers` exports its two kinds, and that is why the kitchen can
import them. The kitchen has a list of its own:

```yaml
# package.yaml
exports:
  kinds: [Station, Recipe, Dish, BaseRecipe, Ingredient, Allergen]
  links: [madeAtStation, usesRecipe, usesIngredient, usesGarnish, containsAllergen, hasAllergen, listsAllergen, dependsOnStation]
  valueTypes: [Name, Minutes, Price, URL]
```

A definition that is not in `exports` is internal. The definitions about
spoilage and about buying are not in the lists, so no other package can use
them, and the kitchen can change them freely. An inverse name, such as
`makesRecipe`, is not in the list. The package exports it with its link.

## Extension points

A different package can use an exported kind as the end of its own links.
It cannot extend the kind, or put a link under an exported link, unless the
kitchen opens that point:

```yaml
# package.yaml
openFor:
  - extend: Recipe
  - under: usesIngredient
```

With these two lines, a package of the bar can declare a kind `Cocktail`
that extends `Recipe`, and a link under `usesIngredient`.

## Tags

An app for guests tells which dishes contain which allergens. It reads an
export of the package, but it does not need all of it: it needs the recipes,
the ingredients and the allergens, not the stations or the batches. A tag is
a label that marks the definitions that belong together for one purpose.
The manifest lists the tags of the package:

```yaml
# package.yaml
tags:
  - name: allergens
    description: What an allergen app needs to tell guests which dishes contain which allergens.
  - name: spoilage
    description: What the kitchen needs to keep food safe until the cooks use it. Only the kitchen uses it, so the package does not export it.
```

Each definition lists its tags in a `tags` key. Most definitions of the
kitchen get one tag. `Ingredient` is necessary for the two purposes, so it
gets the two tags:

```yaml
# defs/ingredient.of.yaml
Ingredient:
  is: Kind
  tags: [allergens, spoilage]
  definition: A food that the restaurant buys and does not make.
  examples:
    - flour
    - pine nuts
  notExamples:
    - the tomato sauce, which the kitchen makes
  fields:
    name:
      type: Name

containsAllergen:
  is: Link
  tags: [allergens]
  from: Ingredient
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is a part of the ingredient.
```

A set of tagged definitions must be complete. `containsAllergen` has the tag
`allergens`, so `Ingredient` and `Allergen` must have it too, and so must the
value type `Name` of their fields. `check` rejects a package where a tagged
definition refers to a definition without the tag. A tag that the manifest
does not list is also an error, so a misspelled tag does not start a
different set.

Tags do not change how derivatives, constraints and tests run. They select
what an export to RDF (Resource Description Framework) includes. The export
is [a planned output of the tool](../../cli/README.md#the-rdf-export).

## Docs

The definitions tell what the rules are. They do not tell the story of the
guest who became ill. A package keeps such text in Markdown files in its
`docs/` folder. A doc has Markdown links to definitions, or embeds them with
the syntax of an image:

````markdown
# Allergens

An allergen can cause a dangerous reaction in some persons. A recipe contains
each allergen of its ingredients, and each allergen of the recipes that it
uses:

![](../defs/allergens.of.yaml#allergenFromIngredient)

![](../defs/allergens.of.yaml#allergenFromRecipe)

In the past, a guest became ill. The pesto in a dish contained nuts, and the
menu did not tell the guest. Thus, the menu must give a warning about each
allergen in a dish. Refer to
[dishListsItsAllergens](../defs/allergens.of.yaml#dishListsItsAllergens).
````

The part after `#` is the name of a statement. `check` makes sure that the
file and the statement exist.

A `Doc` definition gives the file a name among the definitions. A kind then
lists the docs that tell more about it:

```yaml
# defs/allergens.of.yaml
AllergenSafety:
  is: Doc
  path: file:../docs/allergens.md

# defs/ingredient.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name]
  - from: file:allergens.of.yaml
    use: [AllergenSafety]

Allergen:
  is: Kind
  tags: [allergens]
  definition: A substance in food that can cause a dangerous reaction in some persons.
  examples:
    - nuts
    - milk
  fields:
    name:
      type: Name
  docs: [AllergenSafety]
```

Each Markdown file in `docs/` must have a `Doc`, unless a different doc
refers to it. The package has a second doc, `docs/spoilage.md`, with the
`Doc` definition `Spoilage`, which the kind `Batch` lists.

## The finished package

The package holds the full model of the kitchen, and all its tests pass:

```console
$ of-lang test --packages ../../../examples step-06
pass  test allergenThroughASauce  (step-06/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (step-06/tests/examples/base-recipe.of.yaml)
pass  test batchKeptTooWarm  (step-06/tests/examples/batch-kept-too-warm.of.yaml)
pass  test dishWithoutAllergenWarning  (step-06/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-06/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-06/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (step-06/tests/terms/garnish.of.yaml)
7 passed, 0 failed
```

It holds the same files as [the example package](../../examples/restaurant/).
One idea is left: how the package shows a change before somebody accepts it.
That is the last chapter: [Pending changes](07-pending-changes.md).

More about this: [packages](../../language/reference/package.md),
[tags](../../language/reference/tag.md) and
[imports](../../language/reference/files.md#imports) in the reference. In the
language by example: [Dependencies](../by-example/dependencies.md),
[Import as](../by-example/import-as.md),
[Open for](../by-example/open-for.md), [Tags](../by-example/tags.md) and
[Docs](../by-example/docs.md).
