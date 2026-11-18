---
title: Ask a question — Introduction to the OpenFjord modeling language
---

# Ask a question

> Part of [the introduction](README.md).

The waiter cannot go to the kitchen for each guest and each dish. The
restaurant wants one list: all the dishes that contain nuts. The kitchen
does not make this list by hand, because a list made by hand is wrong after
the next change to a recipe. It asks the model.

## A second dish

A list of one dish proves little. The kitchen writes down a second dish, the
marinara, in the same file of facts. The marinara uses the tomato sauce, but
not the pesto:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
dish/marinara:
  is: Fact
  kind: Dish
  fields:
    name: marinara
  links:
    usesRecipe: recipe/tomato_sauce
```

## Ask the model

A question asks something about the facts of an example, and tells what the
answer must be. The kitchen puts its questions in a folder of their own:

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

- `asks` is the question in plain words.
- `given` names the file with the facts.
- `match` has the same form as in a derivative. It finds each dish with each
  allergen that the dish contains. A question sees the facts that the file
  states and the derived facts in the same way. Here, each `hasAllergen`
  link is a derived fact.
- `filter` lists conditions that narrow the matches. This condition keeps
  the matches where the allergen is `allergen/nuts`.
- `expect` gives, for a variable, the ids that it must represent in all the
  matches together.

`test` asks each question of the package, and compares the answer with
`expect`:

```console
$ of-lang test step-02
pass  test allergenThroughASauce  (step-02/tests/examples/allergen-through-a-sauce.of.yaml)
pass  question whichDishesContainNuts  (step-02/tests/questions/which-dishes-contain-nuts.of.yaml)
2 passed, 0 failed
```

The answer is the margherita, and not the marinara.

## A change to a recipe changes the answer

A cook decides to put pesto on the marinara too. The cook changes one line
of the fact `dish/marinara`: `usesRecipe: [recipe/tomato_sauce, recipe/pesto]`.
Nobody thinks about the list of dishes with nuts. The tool does:

```console
$ of-lang test step-02
pass  test allergenThroughASauce  (step-02/tests/examples/allergen-through-a-sauce.of.yaml)
FAIL  question whichDishesContainNuts  (step-02/tests/questions/which-dishes-contain-nuts.of.yaml)
      $dish: expected [dish/margherita], got [dish/margherita, dish/marinara]
1 passed, 1 failed
```

The line below `FAIL` gives the answer of the model after the change. The
restaurant then knows that the list must have the two dishes. The package in
[`steps/step-02`](steps/step-02/) does not have this change.

## The garnish counts too

The cook at the pizza station puts parmesan on the margherita at the end.
Parmesan contains milk. A cook who lists the ingredients of a recipe can
forget the garnish, so the kitchen gives it a link of its own:

```yaml
# defs/recipe.of.yaml
usesGarnish:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} is garnished with {to}"
  definition: The cook puts the ingredient on the food at the end, to make it look or taste better.
  under: usesIngredient
```

`under` names a broader link. Each garnish of a recipe is also an ingredient
of that recipe. `under` is for links what `extends` is for kinds.

The file of facts gets the allergen milk, and the parmesan, which contains
it. The margherita gets its garnish:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
  links:
    usesRecipe: [recipe/tomato_sauce, recipe/pesto]
    usesGarnish: ingredient/parmesan
```

No fact states that the margherita uses parmesan as an ingredient. The tool
adds that fact because of `under`. Then the derivative from
[chapter 1](01-margherita-and-nuts.md), which finds `usesIngredient` links,
finds the parmesan too:

```console
$ of-lang run step-02 step-02/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/milk  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
rejected:
  nothing
```

The first line has the reason `(under)`, where the other lines name a
derivative. The margherita contains nuts through a sauce, and milk through a
garnish. The kitchen wrote no rule about garnishes and allergens.

The restaurant can ask the model which dishes contain an allergen, and the
answer follows each change to a recipe. But nothing makes the menu tell the
guests. That is the next chapter:
[Enforce a rule](03-menu-that-hides-an-allergen.md).

More about this: [questions](../../language/reference/test.md#questions) and
[what links imply](../../language/reference/link.md#what-links-imply) in the
reference. In the language by example:
[Questions](../by-example/questions.md), [Filter](../by-example/filter.md),
[Kind items](../by-example/kind-items.md),
[Link items](../by-example/link-items.md) and
[Links under links](../by-example/under.md).
