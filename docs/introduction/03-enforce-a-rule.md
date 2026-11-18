---
title: Enforce a rule — Introduction to the OpenFjord modeling language
---

# Enforce a rule

> Part of [the introduction](README.md).

In the past, a guest of the restaurant became ill. The pesto in a dish
contained nuts, and the menu did not tell the guest. The model can conclude
that a dish contains nuts. But it accepts a menu without the warning. The
kitchen wants the tool to reject such a menu, and to tell which warning is
missing.

## What the menu tells

The menu shows each dish with its price, and with a warning for each
allergen. A dish that the restaurant does not sell at this time is not on
the menu. The kitchen adds this to the model. A dish gets two fields:

```yaml
# defs/recipe.of.yaml
Dish:
  is: Kind
  extends: Recipe
  definition: A recipe that the restaurant sells to guests.
  examples:
    - the pizza margherita
  notExamples:
    - the tomato sauce, which is only a part of a dish
  fields:
    price:
      type: Price
    status:
      enum: [active, retired]
      default: active
```

`price` has a value type, `Price`, which
[`defs/names.of.yaml`](steps/step-03/defs/names.of.yaml) defines as a number
that is 0 or more. `status` does not name a type. `enum` lists the values
that it can have. A dish that gives no status has the `default`, `active`.
A field is required by default, so each dish in the facts gets a price.

A warning on the menu is a link from a dish to an allergen:

```yaml
# defs/allergens.of.yaml
listsAllergen:
  is: Link
  from: Dish
  to: Allergen
  reads: "the menu tells guests that {from} contains {to}"
  definition: The menu gives a warning about the allergen next to the dish.
```

This link starts at `Dish`, not at `Recipe`. Only a dish is on the menu, so
only a dish can have the link.

## The rule

A constraint is a rule about what must not occur. It finds things by their
links and conditions, and rejects each match that breaks it:

```yaml
# defs/allergens.of.yaml
dishListsItsAllergens:
  is: Constraint
  says: >
    The menu must give a warning about each allergen in an active dish.
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Allergen, as: $allergen }
    - link: hasAllergen
      from: $dish
      to:   $allergen
  filter:
    - node:  $dish
      field: status
      equals:
        value: active
  unless:
    - at least one:
        match:
          - link: listsAllergen
            from: $dish
            to:   $allergen
  message: "The menu must warn that {$dish} contains {$allergen}."
```

`says` and `match` are the same as in a derivative. The match finds each
dish with each allergen that it contains. The filter keeps the dishes with
the status `active`: a condition can name a `field` of a thing and compare
its value.

A constraint has no `derive`. It has two keys of its own:

- `unless` lists conditions. A match is a violation unless all of them hold.
  Here, the match is correct when the menu lists the allergen for the dish.
  `at least one` holds when its nested block matches one time or more.
- `message` is the text that the tool shows for each violation. It can name
  variables.

The match uses `hasAllergen`, a link that only derivatives add. All the
derivatives finish first. Then the constraints check the facts that the
files state and the derived facts together.

## The menu that made a guest ill

The kitchen writes down the case of the guest, in
[a second file of facts](steps/step-03/tests/examples/dish-without-allergen-warning.of.yaml):
a margherita that uses the pesto, and no `listsAllergen` link. `run` shows
what the tool does with it:

```console
$ of-lang run step-03 step-03/tests/examples/dish-without-allergen-warning.of.yaml
derived:
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
rejected:
  dishListsItsAllergens: The menu must warn that dish/margherita contains allergen/nuts.
```

The tool rejects the menu, and the message tells what the menu must do. A
rejection also makes `run` exit with the code 1, so a script can stop the
menu before it goes to the printer.

The file has a test that keeps this result:

```yaml
# tests/examples/dish-without-allergen-warning.of.yaml
dishWithoutAllergenWarning:
  is: Test
  says: >
    The margherita uses the pesto, and the pesto contains nuts. The menu does
    not give a warning about nuts, so the constraint rejects the dish.
  expect:
    rejected: [dishListsItsAllergens]
```

Each file of facts holds exactly one test. `says` tells in plain words what
the example shows. `expect` is what the tool compares with its results, and
`rejected` lists all the constraints that the facts must break. The test
fails if the facts break a different constraint, or none. The file imports
the constraint by its name, as it imports kinds and links.

## The menu with the warnings

In the example from the chapters before this one, the margherita gets a
price and its two warnings, `listsAllergen: [allergen/nuts, allergen/milk]`.
The marinara gets a price. It contains no allergen, so it needs no warning:

```console
$ of-lang run step-03 step-03/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/milk  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
rejected:
  nothing
```

All the tests pass. The test that expects a rejection passes because the
rejection occurs:

```console
$ of-lang test step-03
pass  test allergenThroughASauce  (step-03/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test dishWithoutAllergenWarning  (step-03/tests/examples/dish-without-allergen-warning.of.yaml)
pass  question whichDishesContainNuts  (step-03/tests/questions/which-dishes-contain-nuts.of.yaml)
3 passed, 0 failed
```

The mistake that made a guest ill cannot come back without a message from
the tool. When a cook adds an ingredient with an allergen to a recipe, each
active dish that uses the recipe needs the warning, or the tool rejects the
menu.

Allergens are not the only thing that goes through the kitchen from recipe
to recipe. Work does too:
[Trace dependencies](04-trace-dependencies.md)

More about this: [constraints](../../language/reference/constraint.md),
[fields](../../language/reference/kind.md#fields) and
[tests](../../language/reference/test.md#examples) in the reference. In the
language by example: [Constraints](../by-example/constraints.md),
[Unless](../by-example/unless.md),
[At least one](../by-example/at-least-one.md),
[Messages](../by-example/messages.md),
[Expect rejected facts](../by-example/expect-rejected.md),
[Enums](../by-example/enums.md), [Defaults](../by-example/defaults.md) and
[Value types](../by-example/value-types.md).
