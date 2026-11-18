---
title: Derivatives — Introduction to the OpenFjord modeling language
---

# Derivatives

> Part of [the introduction](README.md).

A guest asks: "Does the margherita contain nuts?" The facts do not tell it
directly. The margherita uses the pesto, the pesto uses pine nuts, and pine
nuts contain nuts. A cook knows the rule: a recipe contains each allergen of
its ingredients, and each allergen of the recipes that it uses. The kitchen
writes this rule down as derivatives. A derivative adds the facts that
follow from other facts.

## A link that nobody states

First, the rule needs a link for its result. No fact states this link. Only
derivatives add it:

```yaml
# defs/allergens.of.yaml
hasAllergen:
  is: Link
  from: Recipe
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is in the recipe, through one of its ingredients or one of its recipes.
```

## The first derivative

```yaml
# defs/allergens.of.yaml
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

`says` is the rule in plain words. The other two keys are the logic:

- `match` lists the things and the links to find. A kind item, such as
  `kind: { name: Recipe, as: $recipe }`, gives a variable to the things of a
  kind. A variable starts with `$`. A link item is one link that must exist
  between two things.
- `derive` tells what to add for each match.

Read it as follows: for each recipe that uses an ingredient that contains an
allergen, add "the recipe contains the allergen".

## A derivative that uses derived facts

The second half of the rule is about recipes that use other recipes:

```yaml
# defs/allergens.of.yaml
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

This derivative finds `hasAllergen` links and adds `hasAllergen` links. The
derivatives run repeatedly until no new fact follows. Thus, an allergen
goes from an ingredient to a sauce, and from the sauce to each dish that
uses it.

`run` shows each derived fact with the derivative that added it:

```console
$ of-lang run 04 04/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/milk  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  station/pizza depends on station/sauce  (derivative stationDependsOnStation)
rejected:
  nothing
```

The margherita contains milk because of its garnish: the parmesan is an
ingredient of the dish, as [chapter 3](03-facts-and-tests.md) shows.

`explain` shows why one fact holds, down to the facts that the file states:

```console
$ of-lang explain 04 04/tests/examples/allergen-through-a-sauce.of.yaml dish/margherita hasAllergen allergen/nuts
dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  A recipe contains each allergen that one of its recipes contains.
  with $allergen = allergen/nuts, $part = recipe/pesto, $recipe = dish/margherita
  because:
    dish/margherita uses recipe/pesto  (written at 04/tests/examples/allergen-through-a-sauce.of.yaml:117:5)
    recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
      A recipe contains each allergen that one of its ingredients contains.
      with $allergen = allergen/nuts, $ingredient = ingredient/pine_nuts, $recipe = recipe/pesto
      because:
        recipe/pesto uses ingredient/pine_nuts  (written at 04/tests/examples/allergen-through-a-sauce.of.yaml:96:5)
        ingredient/pine_nuts contains allergen/nuts  (written at 04/tests/examples/allergen-through-a-sauce.of.yaml:72:5)
```

## A derivative that gives a kind

Many dishes use the tomato sauce. A change to such a recipe needs extra
care, so the kitchen wants a name for it:

```yaml
# defs/recipe.of.yaml
BaseRecipe:
  is: Kind
  extends: Recipe
  definition: A recipe that many other recipes use, so a change to it needs extra care.

baseRecipeFromUses:
  is: Derivative
  says: >
    A recipe that more than two recipes use, directly or through other
    recipes, is a base recipe.
  match:
    - kind: { name: Recipe, as: $recipe }
  filter:
    - count:
        match:
          - link: usesRecipe
            from: Recipe
            to:   $recipe
      is above:
        value: 2
  derive:
    - node: $recipe
      becomes: BaseRecipe
```

No fact states that a recipe is a base recipe. The derivative gives the kind
to each recipe that enough recipes use.

- `filter` lists conditions that narrow the matches. This condition counts
  the matches of a nested block, and compares the number with 2.
- In the nested block, `from: Recipe` is a kind alone. The recipe at that
  end only has to exist, so it needs no variable.
- `node` and `becomes` add a kind to a thing. The kind must extend the kind
  of the variable.

`usesRecipe` chains, so the count includes the recipes that use `$recipe`
through other recipes.

```console
$ of-lang explain 04 04/tests/examples/base-recipe.of.yaml recipe/tomato_sauce BaseRecipe
recipe/tomato_sauce is a BaseRecipe  (derivative baseRecipeFromUses)
  A recipe that more than two recipes use, directly or through other recipes, is a base recipe.
  with $recipe = recipe/tomato_sauce
```

The file `defs/station.of.yaml` gets a link and a derivative in the same
way. The link `dependsOnStation` goes from a station to a station. The
derivative `stationDependsOnStation` adds it when a recipe of the first
station uses a recipe of the second station. Its match uses the inverse name
`makesRecipe`, and its filter keeps a station from a link to itself.
[The file in `steps/04`](steps/04/defs/station.of.yaml) shows the two
statements. That file imports from `recipe.of.yaml`, and `recipe.of.yaml`
imports from it. Files of one package can import each other in a circle.

## The test grows

The example from chapter 3 can expect more. It gets the facts that the
derivatives add, and its `says` tells the full story:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
allergenThroughASauce:
  is: Test
  says: >
    The margherita uses the pesto, which contains nuts. Parmesan, which
    contains milk, is on the margherita as a garnish. Thus, the margherita
    contains nuts and milk. The menu gives a warning about the two allergens,
    so the dish breaks no constraint. The pizza station depends on the sauce
    station.
  expect:
    derived:
      - link: hasAllergen
        from: dish/margherita
        to:   allergen/nuts
      - link: hasAllergen
        from: dish/margherita
        to:   allergen/milk
      - link: usesIngredient
        from: dish/margherita
        to:   ingredient/parmesan
      - link: dependsOnStation
        from: station/pizza
        to:   station/sauce
    notDerived:
      - link: hasAllergen
        from: dish/marinara
        to:   allergen/nuts
      - link: dependsOnStation
        from: station/sauce
        to:   station/pizza
    rejected: []
```

The last key, `rejected`, lists the constraints that the facts must break.
The list is empty, and the package has no constraints yet.
[Chapter 5](05-constraints.md) explains them.

## Questions

A test tells what must follow. A question asks something about the facts of
an example, and tells what the answer must be:

```yaml
# tests/questions/which-dishes-contain-nuts.of.yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Dish]
  - from: file:../../defs/ingredient.of.yaml
    use: [Allergen]
  - from: file:../../defs/allergens.of.yaml
    use: [hasAllergen]

whichDishesContainNuts:
  is: Question
  asks: Which dishes contain nuts?
  given: file:../examples/allergen-through-a-sauce.of.yaml
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Allergen, as: $allergen }
    - link: hasAllergen
      from: $dish
      to:   $allergen
  filter:
    - node: $allergen
      equals:
        value: allergen/nuts
  expect:
    $dish: [dish/margherita]
```

`given` names the example file with the facts. `match` and `filter` have the
same form as in a derivative. A question sees the stated facts and the
derived facts in the same way. `expect` gives, for a variable, the ids that
it must represent in all the matches together. The answer is the margherita,
and not the marinara.

A second question asks which stations stop when the sauce station stops. A
second example proves `baseRecipeFromUses` with three dishes that use the
tomato sauce.

```console
$ of-lang test 04
pass  test allergenThroughASauce  (04/tests/examples/allergen-through-a-sauce.of.yaml)
pass  test tomatoSauceInThreeDishes  (04/tests/examples/base-recipe.of.yaml)
pass  question whichDishesContainNuts  (04/tests/questions/which-dishes-contain-nuts.of.yaml)
pass  question whichStationsStopWithTheSauceStation  (04/tests/questions/which-stations-stop-with-the-sauce-station.of.yaml)
pass  test garnishIsAnIngredient  (04/tests/terms/garnish.of.yaml)
5 passed, 0 failed
```

A derivative only adds facts. It rejects nothing. The kitchen also has
rules about what must not occur, and the next chapter writes them down:
[Constraints](05-constraints.md).

More about this: [derivatives](../../language/reference/derivative.md),
[derived facts](../../language/reference/fact.md#derived-facts) and
[questions](../../language/reference/test.md#questions) in the reference. In
the language by example: [Derive a link](../by-example/derive-link.md),
[Derive a kind](../by-example/derive-kind.md),
[Derivatives together](../by-example/derivatives-together.md),
[Kind items](../by-example/kind-items.md),
[Link items](../by-example/link-items.md),
[A kind alone](../by-example/kind-alone.md),
[Filter](../by-example/filter.md), [Count](../by-example/count.md) and
[Questions](../by-example/questions.md).
