---
title: Does the margherita contain nuts? — Introduction to the OpenFjord modeling language
---

# Does the margherita contain nuts?

> Part of [the introduction](README.md).

A guest with a nut allergy asks the waiter: "Does the pizza margherita
contain nuts?" A wrong answer can make the guest ill. The cooks know the
parts of the answer: the margherita uses the pesto, the pesto uses pine
nuts, and pine nuts contain nuts. But nobody wrote these parts down, and
nobody put them together.

In this chapter, the kitchen writes down what it knows. Then the tool,
`of-lang`, puts the parts together. This is the answer that the chapter
ends with:

```console
$ of-lang run step-01 step-01/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
rejected:
  nothing
```

The second line under `derived` is the answer: the margherita contains nuts.
No person wrote that line. The tool concluded it.

## A package for the kitchen

A model is in a package. A package is a folder with a manifest, a `defs/`
folder for the definitions, and a `tests/` folder for real things and what
must follow from them. The manifest is the file `package.yaml`:

```yaml
# package.yaml
package: "@bistro/kitchen"
version: 0.1.0
owner: kitchen
rdfNamespace: kitchen
```

`package` is the name of the package: `kitchen`, in the namespace `@bistro`.
`owner` is who must approve each change to it. `rdfNamespace` is a prefix
that the tool uses only when it exports the package to RDF (Resource
Description Framework), a standard format for models.

## The words of the question

The question is about a dish and an allergen. The answer also needs a recipe
that is not a dish, the pesto, and an ingredient, the pine nuts. The kitchen
gives a name and a definition to each of these kinds of things. Recipes and
dishes are in one file:

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
  examples:
    - the tomato sauce
    - the pizza dough
  notExamples:
    - a food that the restaurant buys and does not make, such as flour
    - a new dish that a cook tries one time and does not keep
  fields:
    name:
      type: Name

Dish:
  is: Kind
  extends: Recipe
  definition: A recipe that the restaurant sells to guests.
  examples:
    - the pizza margherita
  notExamples:
    - the tomato sauce, which is only a part of a dish
```

`Recipe` and `Dish` are statements. A file holds statements, and the key of
each statement is its name. The `is` key tells what the statement is. These
two are kinds: categories of things.

- `definition` is one sentence that tells what the kind is. `examples` and
  `notExamples` show where the kind starts and where it stops.
- `fields` lists the values that each thing of the kind has. A recipe has
  one field, its `name`.
- `extends` makes a narrower kind. Each dish is also a recipe, so a dish has
  the fields of a recipe, and the tool can use a dish where it looks for a
  recipe.

The `name` field has the type `Name`. `Name` is a value type: a description
of a value, with a name of its own, so that many fields can use it. The file
[`defs/names.of.yaml`](steps/step-01/defs/names.of.yaml) defines it as text
that matches a pattern.

The file [`defs/ingredient.of.yaml`](steps/step-01/defs/ingredient.of.yaml)
defines the kinds `Ingredient` and `Allergen` in the same form.

## What uses what

A kind tells that recipes exist. It does not tell that a recipe uses an
ingredient. A link does this. A link is a named relationship from one kind
to a different kind, or to the same kind, in one direction:

```yaml
# defs/recipe.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name]
  - from: file:ingredient.of.yaml
    use: [Ingredient]

usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe, directly or through other recipes.
  chains: true
  inverse: usedByRecipe
```

`from` is the kind where the link starts, and `to` is the kind where it
ends. `reads` is the sentence that people read for one such link: "the
margherita uses the pesto". The tool uses this sentence in its answers.

The last two keys are optional. `chains: true` tells that a chain of the
link is also the link: when the lasagne uses the ragù, and the ragù uses the
tomato sauce, then the lasagne uses the tomato sauce. `inverse` gives a name
to the link in the opposite direction.

A file can use a name only if the file defines the name or imports it. Thus,
the file starts with an `import` block. `from` gives the file that defines
the names, and `use` lists the names.

Two more links have the same form, without the optional keys.
`usesIngredient` goes from `Recipe` to `Ingredient`, in the same file.
`containsAllergen` goes from `Ingredient` to `Allergen`, in
`defs/ingredient.of.yaml`.

## The margherita, written down

The definitions tell what can exist. A fact tells what one real thing is:
its kind, its fields and its links. Facts are in the `tests/` folder:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
recipe/pesto:
  is: Fact
  kind: Recipe
  fields:
    name: pesto
  links:
    usesIngredient: ingredient/pine_nuts

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
  links:
    usesRecipe: [recipe/tomato_sauce, recipe/pesto]
```

The name of a fact is its id, such as `dish/margherita`. `links` gives, for
each link, the id of the thing at the other end, or a list of ids.
[The file](steps/step-01/tests/examples/allergen-through-a-sauce.of.yaml)
holds the other things in the same form: the pine nuts, which contain the
allergen nuts, and the tomato sauce with its tomatoes.

These facts are the three parts that the cooks know. No fact tells that the
margherita contains nuts.

## The rule that each cook knows

A cook puts the parts together with a rule: a recipe contains each allergen
of its ingredients, and each allergen of the recipes that it uses. The
kitchen writes the rule in a file for the topic, `defs/allergens.of.yaml`.
The result of the rule is a link, and a derivative adds it. A derivative
finds things by their links, and adds a fact for each match:

```yaml
# defs/allergens.of.yaml
hasAllergen:
  is: Link
  from: Recipe
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is in the recipe, through one of its ingredients or one of its recipes.

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
```

No fact states a `hasAllergen` link. Only derivatives add it.

`says` is the rule in plain words. `match` lists what to find. A kind item,
such as `kind: { name: Recipe, as: $recipe }`, gives a variable to the
things of a kind. A link item is one link that must exist between two
things. `derive` tells what to add for each match. Read it as follows: for
each recipe that uses an ingredient that contains an allergen, add "the
recipe contains the allergen".

This derivative gives nuts to the pesto. A second derivative in
[the file](steps/step-01/defs/allergens.of.yaml), `allergenFromRecipe`, has
the same form. It gives each allergen of a recipe to each recipe that uses
it. It finds `hasAllergen` links and adds `hasAllergen` links. The
derivatives run repeatedly until no new fact follows.

## The answer, with its steps

`run` shows what the tool derives from a file of facts. Its output is at the
top of this chapter: the pesto contains nuts, and thus the margherita
contains nuts. A fact that the tool adds is a derived fact.

A cook does not have to believe the tool. `explain` shows why one fact
holds, down to the facts that the file states:

```console
$ of-lang explain step-01 step-01/tests/examples/allergen-through-a-sauce.of.yaml dish/margherita hasAllergen allergen/nuts
dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  A recipe contains each allergen that one of its recipes contains.
  with $allergen = allergen/nuts, $part = recipe/pesto, $recipe = dish/margherita
  because:
    dish/margherita uses recipe/pesto  (written at step-01/tests/examples/allergen-through-a-sauce.of.yaml:62:5)
    recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
      A recipe contains each allergen that one of its ingredients contains.
      with $allergen = allergen/nuts, $ingredient = ingredient/pine_nuts, $recipe = recipe/pesto
      because:
        recipe/pesto uses ingredient/pine_nuts  (written at step-01/tests/examples/allergen-through-a-sauce.of.yaml:46:5)
        ingredient/pine_nuts contains allergen/nuts  (written at step-01/tests/examples/allergen-through-a-sauce.of.yaml:32:5)
```

Each step names its derivative, gives the rule in plain words, and tells
what each variable represents. The last lines point to the places in the
file where a person wrote a fact.

## Keep the answer

The answer must stay correct when the package changes. Thus, the file of
facts also holds one test, `allergenThroughASauce`. A test tells what must
follow from the facts. This one expects one derived fact: the margherita
contains nuts. `test` runs each test of the package:

```console
$ of-lang test step-01
pass  test allergenThroughASauce  (step-01/tests/examples/allergen-through-a-sauce.of.yaml)
1 passed, 0 failed
```

If a change to the package loses this answer, the test fails. The waiter can
tell the guest that the margherita contains nuts, and the kitchen can show
why.

One guest has an answer. The restaurant wants the answer for all its dishes
at one time: [Which dishes contain nuts?](02-dishes-with-nuts.md)

More about this: [the manifest](../../language/reference/package.md#the-manifest),
[kinds](../../language/reference/kind.md),
[links](../../language/reference/link.md),
[facts](../../language/reference/fact.md) and
[derivatives](../../language/reference/derivative.md) in the reference, and
[the output of the tool](../../cli/README.md#output). In the language by
example: [Hello, kind](../by-example/hello-kind.md),
[Extends](../by-example/extends.md), [Imports](../by-example/imports.md),
[Links](../by-example/links.md), [Facts](../by-example/facts.md),
[Derive a link](../by-example/derive-link.md) and
[Derivatives together](../by-example/derivatives-together.md).
