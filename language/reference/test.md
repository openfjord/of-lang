---
title: Tests — The OpenFjord modeling language
---

# Tests

> Part of [the reference](README.md).

A package uses tests in its `tests/` folder to prove that it is correct
(reason 16 in [the research](../../README.md)). A test is an example or a
question:

- An **example** tells what a set of facts must derive, and which constraints
  the facts must break.
- A **question** asks something about a set of facts, and tells what the
  answer must be.

## Examples

An example is a file of facts with one `Test` statement:

```yaml
# tests/examples/allergen-through-a-sauce.of.yaml
import:
  - from: file:../../defs/station.of.yaml
    use: [Station, dependsOnStation]
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish, madeAtStation, usesRecipe, usesIngredient, usesGarnish]
  - from: file:../../defs/ingredient.of.yaml
    use: [Ingredient, Allergen, containsAllergen]
  - from: file:../../defs/allergens.of.yaml
    use: [hasAllergen, listsAllergen]

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

station/pizza:
  is: Fact
  kind: Station
  fields:
    name: pizza

station/sauce:
  is: Fact
  kind: Station
  fields:
    name: sauce

allergen/nuts:
  is: Fact
  kind: Allergen
  fields:
    name: nuts

allergen/milk:
  is: Fact
  kind: Allergen
  fields:
    name: milk

ingredient/pine_nuts:
  is: Fact
  kind: Ingredient
  fields:
    name: pine_nuts
  links:
    containsAllergen: allergen/nuts

ingredient/parmesan:
  is: Fact
  kind: Ingredient
  fields:
    name: parmesan
  links:
    containsAllergen: allergen/milk

ingredient/tomatoes:
  is: Fact
  kind: Ingredient
  fields:
    name: tomatoes

recipe/pesto:
  is: Fact
  kind: Recipe
  fields:
    name: pesto
  links:
    madeAtStation: station/sauce
    usesIngredient: ingredient/pine_nuts

recipe/tomato_sauce:
  is: Fact
  kind: Recipe
  fields:
    name: tomato_sauce
  links:
    madeAtStation: station/sauce
    usesIngredient: ingredient/tomatoes

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12
  links:
    madeAtStation: station/pizza
    usesRecipe: [recipe/tomato_sauce, recipe/pesto]
    usesGarnish: ingredient/parmesan
    listsAllergen: [allergen/nuts, allergen/milk]

dish/marinara:
  is: Fact
  kind: Dish
  fields:
    name: marinara
    price: 10
  links:
    madeAtStation: station/pizza
    usesRecipe: recipe/tomato_sauce
```

| Key under `expect` | Value |
|---|---|
| `derived` | facts that must hold after the derivatives run |
| `notDerived` | facts that must not hold |
| `rejected` | all constraints that the facts must break |

A fact under `derived` or `notDerived` has the same form as a derivation,
with ids in place of variables:

| Fact | Written |
|---|---|
| a link | `link`, `from`, `to` |
| a kind | `node`, `becomes` |
| a field value | `node`, `field`, `set to` with `value:` |

- **TS1.** An example file MUST hold exactly one `Test` statement, which MUST
  have `says` and `expect`.
- **TS2.** The test passes only if all of these are true:
  - every fact under `derived` holds
  - no fact under `notDerived` holds
  - the facts break exactly the constraints under `rejected`
  - every fact obeys its kind ([FC3](fact.md#rules)).

  A test with no `rejected` is the same as a test that expects no broken
  constraints.
- **TS3.** The file MUST import each link, kind and constraint that the test
  names.

## Questions

A question asks something about the facts of an example, and tells what the
answer must be. It finds things and the links between them, and it uses
conditions to narrow the matches. Then it compares the things that a variable
represented with the expected list.

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
    - link: hasAllergen
      from: {$dish: Dish}
      to:   {$allergen: Allergen}
  filter:
    - node: $allergen
      equals:
        value: allergen/nuts
  expect:
    $dish: [dish/margherita]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Question` |
| `asks` | yes | the question, in plain words |
| `given` | no | `file:` and an example file with the facts to ask about; if not given, the facts in this file |
| `match` | yes | the links to find |
| `filter` | no | conditions that narrow the matches |
| `expect` | yes | for each listed variable, the ids that it must represent |

Each item in `match` is one link that must exist, with `link`, `from` and
`to`. `from` and `to` each give a variable and its kind, as
`{$variable: Kind}`. `match` finds each combination of things for which all
its links hold. Each combination is one match, which gives one thing to each
variable. `repeated: true` on an item follows the link for one or more steps.

`filter` keeps the matches for which each condition holds. A condition names
a `node`, an optional `field`, a comparison word and what to compare with:

| Word | Takes |
|---|---|
| `equals`, `does not equal`, `is above`, `is below`, `is at least`, `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set`, `is not set` | nothing more |

A condition can also contain a nested block with its own `match` and
`filter`. The words for this are `at least one`, `none`, or `count` with a
comparison word. `any of` holds when one of its listed conditions holds.

- **TS4.** A question MUST have an `asks` sentence, and the sentence MUST NOT
  contain a `$` word.
- **TS5.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file imports.
- **TS6.** Each `from` and each `to` MUST use the form `{$variable: Kind}`,
  every time and also in nested blocks. Every mention of a variable MUST give
  the same kind. The kind MUST be the kind of the link at that end, or a kind
  that extends it.
- **TS7.** Conditions MUST write a variable without its kind (`$allergen`), and
  every variable that they use MUST appear in a `match` item. A question
  MUST NOT use a variable outside a nested block if the variable appears only
  in that block.
- **TS8.** A comparison MUST compare values of the same type: two fields that
  name the same value type, or two fields that name no value type and have
  the same built-in type. A literal MUST be a value that the field allows. A
  field MUST be a field of the kind of its node. Two nodes that a condition
  compares without fields MUST have the same kind, or kinds where one
  extends the other.
- **TS9.** Each key under `expect` MUST be a variable of the question's
  `match`.
- **TS10.** The question passes only if, for each listed variable, the
  distinct things that it represents in all matches are exactly the listed
  ids.
