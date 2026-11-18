---
title: A modeling language for software people
---

# A modeling language for software people

The OpenFjord modeling language describes the things of a business, the
links between them, and the rules about them, in YAML files next to the code.
Its tool, `of-lang`, checks those files, derives the facts that nobody wrote
down, explains each one, and runs the model's tests.

<p align="center">
  <img src="docs/images/openfjord.png" alt="The OpenFjord logo: a fjord between mountains, above circuit lines" width="220">
</p>

> **Version 0.1.0.** Before version 1.0, the language can change in ways
> that break a package.

<!-- not on the site -->
The documentation is also a site: <https://lang.openfjord.io/>.
<!-- end -->

## From nothing to an answer

A guest asks: "Does the margherita contain nuts?" Nobody wrote the answer
down. The margherita uses the pesto, and the pesto uses pine nuts. This page
builds a small model of a restaurant's menu that finds the answer, shows the
steps, and rejects a menu that does not warn the guest.

### 1. Install the tool

With [Homebrew](https://brew.sh):

```bash
brew install openfjord/tap/of-lang
```

Or download a binary for macOS or Linux from the
[releases](https://github.com/openfjord/of-lang/releases).

### 2. Make a package

A model is in a package: a folder with a manifest, definitions and tests.

```
menu/                      # the root of the package
  package.yaml             # the manifest
  defs/
    menu.of.yaml           # the definitions
  tests/
    examples/
      margherita.of.yaml   # facts, and what must follow from them
```

The manifest gives the package its name, its version and its owner. Only
the export to RDF (Resource Description Framework) uses `rdfNamespace`:

```yaml
# package.yaml
package: "@bistro/menu"
version: 0.1.0
owner: kitchen
rdfNamespace: menu
```

### 3. Name the things

A kind is a category of things. Each one has a definition in plain words. A
dish is a recipe that the restaurant sells, so `Dish` extends `Recipe`.

```yaml
# defs/menu.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.

Dish:
  is: Kind
  extends: Recipe
  definition: A recipe that the restaurant sells to guests.

Ingredient:
  is: Kind
  definition: A food that the kitchen buys and does not make.

Allergen:
  is: Kind
  definition: A substance in food that can cause a dangerous reaction in some persons.
```

A link is a named relationship from one kind to a different kind, or to the
same kind. `reads` is the sentence that people read for one such link.

```yaml
usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe.

usesIngredient:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} uses {to}"
  definition: The recipe needs the ingredient.

containsAllergen:
  is: Link
  from: Ingredient
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is a part of the ingredient.

hasAllergen:
  is: Link
  from: Recipe
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is in the recipe, through one of its ingredients or one of its recipes.
```

### 4. Say what follows

Nobody states the last link, `hasAllergen`. Two derivatives add it. The
first one finds the allergens of the ingredients of a recipe. The second one
passes them on to each recipe that uses that recipe.

```yaml
allergenFromIngredient:
  is: Derivative
  says: A recipe contains each allergen that one of its ingredients contains.
  match:
    - kind: { name: Recipe, as: $recipe }
    - kind: { name: Ingredient, as: $ingredient }
    - kind: { name: Allergen, as: $allergen }
    - link: usesIngredient
      from: $recipe
      to:   $ingredient
    - link: containsAllergen
      from: $ingredient
      to:   $allergen
  derive:
    - link: hasAllergen
      from: $recipe
      to:   $allergen

allergenFromRecipe:
  is: Derivative
  says: A recipe contains each allergen that one of its recipes contains.
  match:
    - kind: { name: Recipe, as: $recipe }
    - kind: { name: Recipe, as: $part }
    - kind: { name: Allergen, as: $allergen }
    - link: usesRecipe
      from: $recipe
      to:   $part
    - link: hasAllergen
      from: $part
      to:   $allergen
  derive:
    - link: hasAllergen
      from: $recipe
      to:   $allergen
```

### 5. Say what is not allowed

A constraint rejects the facts that break a rule. Here, each match is a
violation, unless the menu gives a warning about the allergen.

```yaml
listsAllergen:
  is: Link
  from: Dish
  to: Allergen
  reads: "the menu tells guests that {from} contains {to}"
  definition: The menu gives a warning about the allergen next to the dish.

dishListsItsAllergens:
  is: Constraint
  says: The menu must give a warning about each allergen in a dish.
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Allergen, as: $allergen }
    - link: hasAllergen
      from: $dish
      to:   $allergen
  unless:
    - at least one:
        match:
          - link: listsAllergen
            from: $dish
            to:   $allergen
  message: "The menu must warn that {$dish} contains {$allergen}."
```

### 6. State the facts

The facts tell about one kitchen. No fact says that the margherita contains
nuts.

```yaml
# tests/examples/margherita.of.yaml
import:
  - from: file:../../defs/menu.of.yaml
    use: [Recipe, Dish, Ingredient, Allergen, usesRecipe, usesIngredient, containsAllergen, hasAllergen, dishListsItsAllergens]

allergen/nuts:
  is: Fact
  kind: Allergen

ingredient/pine_nuts:
  is: Fact
  kind: Ingredient
  links:
    containsAllergen: allergen/nuts

recipe/pesto:
  is: Fact
  kind: Recipe
  links:
    usesIngredient: ingredient/pine_nuts

dish/margherita:
  is: Fact
  kind: Dish
  links:
    usesRecipe: recipe/pesto
```

A test tells what must follow from these facts:

```yaml
margheritaContainsNuts:
  is: Test
  says: >
    The margherita uses the pesto, and the pesto contains pine nuts. Thus,
    the margherita contains nuts. The menu does not say so, so the
    constraint rejects the dish.
  expect:
    derived:
      - link: hasAllergen
        from: dish/margherita
        to:   allergen/nuts
    rejected: [dishListsItsAllergens]
```

### 7. Get the answer

`check` makes sure that the package obeys the rules of the language:

```console
$ of-lang check menu
ok
```

`run` derives the facts that follow, and then applies the constraints. The
margherita contains nuts, and the menu does not say so:

```console
$ of-lang run menu menu/tests/examples/margherita.of.yaml
derived:
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
rejected:
  dishListsItsAllergens: The menu must warn that dish/margherita contains allergen/nuts.
```

`explain` shows the steps to one derived fact, down to the lines of the
file that state the facts:

```console
$ of-lang explain menu menu/tests/examples/margherita.of.yaml dish/margherita hasAllergen allergen/nuts
dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  A recipe contains each allergen that one of its recipes contains.
  with $allergen = allergen/nuts, $part = recipe/pesto, $recipe = dish/margherita
  because:
    dish/margherita uses recipe/pesto  (written at menu/tests/examples/margherita.of.yaml:25:5)
    recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
      A recipe contains each allergen that one of its ingredients contains.
      with $allergen = allergen/nuts, $ingredient = ingredient/pine_nuts, $recipe = recipe/pesto
      because:
        recipe/pesto uses ingredient/pine_nuts  (written at menu/tests/examples/margherita.of.yaml:19:5)
        ingredient/pine_nuts contains allergen/nuts  (written at menu/tests/examples/margherita.of.yaml:13:5)
```

`test` runs the test, so that the answer stays correct when the model
changes:

```console
$ of-lang test menu
pass  test margheritaContainsNuts  (menu/tests/examples/margherita.of.yaml)
1 passed, 0 failed
```

The same package is in [`examples/menu`](examples/menu/).

## Where to go next

The same restaurant is the example in all of the documentation.

- [The introduction](docs/introduction/README.md) grows this model into the
  full kitchen in eight chapters. Each chapter is a question of the kitchen,
  such as "Which dishes contain nuts?", and ends with the answer of the tool.
- [The language by example](docs/by-example/README.md) has one short page
  for each thing in the language, with the text next to the code.
- [The language specification](language/README.md) tells what the language
  is, with its principles and its vocabulary.
  - [The cheat sheet](language/cheatsheet.md) has the whole language on one
    page.
  - [The reference](language/README.md#reference) has every rule.
- [The `of-lang` tool](cli/README.md) tells about each command and its
  output, and what is planned.
- [The full example package](examples/restaurant/) is the kitchen that the
  introduction builds. It depends on
  [a second package](examples/suppliers/), the restaurant's suppliers.
- [Why a model](docs/why.md) tells which problems a model solves, who it is
  for, and why this is a language of its own.
- [The design of the language](docs/design.md) tells what the language
  leaves out, the roads not taken, and the open questions.
- [The changelog](CHANGELOG.md) and [how to contribute](CONTRIBUTING.md).

## License

Licensed under the [Apache License 2.0](LICENSE).
