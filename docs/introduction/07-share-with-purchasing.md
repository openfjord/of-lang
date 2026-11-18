---
title: Share a model — Introduction to the OpenFjord modeling language
---

# Share a model

> Part of [the introduction](README.md).

The kitchen is one team in the restaurant. The purchasing team has a model
of its own, with the suppliers of the restaurant. The kitchen wants to tell
which supplier sells each of its ingredients, and it does not want a second
list of suppliers that nobody keeps correct. The two teams want two models
that fit together, where each team owns its own part.

## Buy from the suppliers of purchasing

The package of purchasing is
[`@bistro/suppliers`](../../examples/suppliers/), with the kinds `Supplier`
and `Ingredient`. The kitchen does not define suppliers again. It depends on
the package that defines them:

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

The links of the kitchen end at the kinds of purchasing. The kitchen cannot
change those kinds: they belong to purchasing.

## The two packages check together

The tool must find each package in `dependsOn`. By default, it looks in the
folder that holds the package. The suppliers package is not there, so a
plain `of-lang check step-07` reports that it did not find
`@bistro/suppliers`. `--packages` names the folder to look in. In the
repository, the suppliers package is in `examples/`:

```console
$ of-lang check --packages ../../../examples step-07
ok
```

The tool checks the kitchen against the real definitions of purchasing. If
purchasing removes `Supplier`, or the kitchen misspells it, `check` tells.
Each team works on its own package, and the tool keeps the two together.

## What other packages can use

The kitchen can import `Supplier` because the manifest of
`@bistro/suppliers` exports it. A package decides what other packages can
import from it. The kitchen has a list of its own:

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

## A part of the model for the allergen app

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
```

A set of tagged definitions must be complete, or the app gets a model with
holes. The link `containsAllergen` has the tag `allergens`, so `Ingredient`
and `Allergen` must have it too, and so must the value type `Name` of their
fields. When somebody removes the tag `allergens` from `Name`, `check`
tells what is missing:

```console
$ of-lang check --packages ../../../examples step-07
step-07/defs/ingredient.of.yaml:9:10: Ingredient has the tag allergens and refers to Name, which does not have it (TG5)
step-07/defs/ingredient.of.yaml:22:10: Allergen has the tag allergens and refers to Name, which does not have it (TG5)
step-07/defs/recipe.of.yaml:11:10: Recipe has the tag allergens and refers to Name, which does not have it (TG5)
```

A tag that the manifest does not list is also an error, so a misspelled tag
does not start a different set.

Tags do not change how derivatives, constraints and tests run. They select
what an export to RDF (Resource Description Framework) includes. The export
is [a planned output of the tool](../../cli/README.md#the-rdf-export).

## The reasons, next to the rules

A reader from a different team sees the rules, but not the story of the
guest who became ill. A package keeps such text in Markdown files in its
`docs/` folder. The kitchen writes
[`docs/allergens.md`](steps/step-07/docs/allergens.md). A doc has Markdown
links to definitions, or embeds them with the syntax of an image:

```markdown
![](../defs/allergens.of.yaml#allergenFromIngredient)

Refer to
[dishListsItsAllergens](../defs/allergens.of.yaml#dishListsItsAllergens).
```

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

`defs/ingredient.of.yaml` imports the name `AllergenSafety`, as it imports
each name from a different file.

Each Markdown file in `docs/` must have a `Doc`, unless a different doc
refers to it. The package has a second doc, `docs/spoilage.md`, with the
`Doc` definition `Spoilage`, which the kind `Batch` lists.

## The finished package

The package holds the full model of the kitchen, and the same files as
[the example package](../../examples/restaurant/). All its tests pass, next
to the package of purchasing:

```console
$ of-lang test --packages ../../../examples step-07
pass  test allergenThroughASauce  (step-07/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (step-07/tests/examples/base-recipe.of.yaml)
pass  test batchKeptTooWarm  (step-07/tests/examples/batch-kept-too-warm.of.yaml)
pass  test dishWithoutAllergenWarning  (step-07/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-07/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (step-07/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (step-07/tests/terms/garnish.of.yaml)
7 passed, 0 failed
```

Purchasing owns the suppliers, the kitchen owns the recipes, and the
allergen app gets the part that it needs.

A model that other teams depend on still has to change:
[Change a rule](08-change-a-rule.md).

More about this: [packages](../../language/reference/package.md),
[tags](../../language/reference/tag.md) and
[imports](../../language/reference/files.md#imports) in the reference. In the
language by example: [Dependencies](../by-example/dependencies.md),
[Import as](../by-example/import-as.md),
[Open for](../by-example/open-for.md), [Tags](../by-example/tags.md) and
[Docs](../by-example/docs.md).
