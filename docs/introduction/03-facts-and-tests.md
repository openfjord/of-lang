---
title: Facts and tests — Introduction to the OpenFjord modeling language
---

# Facts and tests

> Part of [the introduction](README.md).

The package defines kinds and links. These are definitions: they tell what
can exist. The kitchen also wants to see that the definitions fit its real
dishes. It writes down some real things, and what the package must tell
about them.

## Facts

A fact tells what one thing is: its kind, its fields and its links. Facts
are in the `tests/` folder of the package, next to `defs/`. The kitchen
writes an example about the pizza margherita:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Station]
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish, madeAtStation, usesRecipe, usesIngredient, usesGarnish]
  - from: file:../../defs/ingredient.of.yaml
    use: [Ingredient, Allergen, containsAllergen]
  - from: file:../../defs/allergens.of.yaml
    use: [listsAllergen]

allergen/nuts:
  is: Fact
  kind: Allergen
  fields:
    name: nuts

ingredient/pine_nuts:
  is: Fact
  kind: Ingredient
  fields:
    name: pine_nuts
  links:
    containsAllergen: allergen/nuts

recipe/pesto:
  is: Fact
  kind: Recipe
  fields:
    name: pesto
    steps: https://recipes.example.com/pesto
  links:
    madeAtStation: station/sauce
    usesIngredient: ingredient/pine_nuts

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    steps: https://recipes.example.com/margherita
    price: 12
  links:
    madeAtStation: station/pizza
    usesRecipe: [recipe/tomato_sauce, recipe/pesto]
    usesGarnish: ingredient/parmesan
    listsAllergen: [allergen/nuts, allergen/milk]
```

The name of a fact is its id. An id is any text. Usually it is a path, such
as `dish/margherita`.

- `kind` is the kind of the thing.
- `fields` gives a value to each field of the kind. The margherita is a
  dish, and a dish is a recipe, so it has a `name` and `steps`. It gives no
  `status`, so its status is the default, `active`.
- `links` gives, for each link, the id of the thing at the other end, or a
  list of ids. Each id is a fact in the same file.

A test file imports each kind and each link that it uses. It does not import
the names of fields.

The file holds more facts in the same form: the two stations, the allergen
milk, the parmesan, the tomatoes, the tomato sauce and a second dish, the
marinara. [The file in `steps/step-03`](steps/step-03/tests/examples/allergen-through-a-sauce.of.yaml)
shows them all.

## A test

Facts alone prove nothing. A file of facts also holds exactly one `Test`
statement, which tells what must follow from the facts:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
allergenThroughASauce:
  is: Test
  says: >
    Parmesan is on the margherita as a garnish. A garnish is an ingredient,
    so the margherita uses parmesan. The marinara has no garnish, so it does
    not use parmesan.
  expect:
    derived:
      - link: usesIngredient
        from: dish/margherita
        to:   ingredient/parmesan
    notDerived:
      - link: usesIngredient
        from: dish/marinara
        to:   ingredient/parmesan
```

`says` tells in plain words what the example shows. `expect` is what the
tool compares with its results.

No fact states that the margherita uses parmesan. The fact of the margherita
states only a `usesGarnish` link. But `usesGarnish` is under
`usesIngredient`, so the tool adds the broader link. A fact that the tool
adds is a derived fact. `derived` lists the facts that the tool must add,
and `notDerived` lists facts that it must not add.

`test` runs each test of the package:

```console
$ of-lang test step-03
pass  test allergenThroughASauce  (step-03/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test garnishIsAnIngredient  (step-03/tests/terms/garnish.of.yaml)
2 passed, 0 failed
```

`run` shows what the tool derives from a file of facts. Each line is the
`reads` sentence of a link, with the reason for the fact:

```console
$ of-lang run step-03 step-03/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
rejected:
  nothing
```

The package has nothing that rejects facts yet.
[Chapter 5](05-constraints.md) adds that.

A test protects the kitchen against mistakes. If somebody deletes the
garnish from the margherita, the test fails and tells what is missing:

```console
$ of-lang test step-03
FAIL  test allergenThroughASauce  (step-03/tests/examples/allergen-through-a-sauce.of.yaml)
      expected derived, but not: dish/margherita uses ingredient/parmesan
pass  test garnishIsAnIngredient  (step-03/tests/terms/garnish.of.yaml)
1 passed, 1 failed
```

## A test without facts

The second test in the output has no facts. It tells what the definitions
derive about kinds and links, before there are any things:

```yaml
# tests/terms/garnish.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish, madeAtStation, usesRecipe, usesIngredient, usesGarnish]

garnishIsAnIngredient:
  is: Test
  says: >
    A garnish is one of the ingredients that a recipe uses. A dish is a
    recipe, so a dish can have a station. A recipe can use a recipe that uses
    another recipe.
  expect:
    derived:
      - kind: Dish
        is a: Recipe
      - kind: Dish
        can have: madeAtStation
      - link: usesGarnish
        implies: usesIngredient
      - link: usesRecipe
        chains: true
    notDerived:
      - link: usesIngredient
        implies: usesGarnish
```

Each item is a derived term. The four items under `derived` follow from
`extends`, from the `from` of a link, from `under` and from `chains`. The
item under `notDerived` tells that `under` goes in one direction only: an
ingredient of a recipe is not always a garnish.

The author of a package chooses the folders in `tests/`. The kitchen puts
the files with facts in `examples/` and the tests without facts in `terms/`.

The tool derives little so far: only what `under` and `chains` imply. The
kitchen wants more. The pesto contains pine nuts, and pine nuts contain
nuts, so the margherita contains nuts. The next chapter teaches the package
to conclude that: [Derivatives](04-derivatives.md).

More about this: [facts](../../language/reference/fact.md),
[tests](../../language/reference/test.md) and
[what links imply](../../language/reference/link.md#what-links-imply) in the
reference, and [the output of the tool](../../cli/README.md#output). In the
language by example: [Facts](../by-example/facts.md),
[Expect derived facts](../by-example/expect-derived.md) and
[Derived terms](../by-example/derived-terms.md).
