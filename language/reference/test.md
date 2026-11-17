---
title: Tests — The OpenFjord modeling language
---

# Tests

> Part of [the reference](../README.md#reference).

A package uses tests in its `tests/` folder to prove that it is correct
(reason 16 in [the research](../../README.md#why-ontologies-are-hard-for-software-people)). A test is an example or a
question:

- An **example** tells what a set of facts must derive, and which constraints
  the facts must break. An example with no facts tells what the definitions
  derive about kinds and links ([derived terms](#derived-terms)).
- A **question** asks something about a set of facts, and tells what the
  answer must be.

## Examples

An example is a file of facts with one `Test` statement. All the
derivatives and constraints of the package run on the facts of the
example. A question asks about the stated facts and the derived facts
together.

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
| `derived` | facts and terms that must hold after the derivatives run |
| `notDerived` | facts and terms that must not hold |
| `rejected` | all constraints that the facts must break |

A fact under `derived` or `notDerived` has the same form as an item under `derive`,
with ids in place of variables:

| Fact | Written |
|---|---|
| a link | `link`, `from`, `to` |
| a kind | `node`, `becomes` |
| a field value | `node`, `field`, `set to` with `value:` |

- **TS1.** An example file MUST hold exactly one `Test` statement, which MUST
  have `says` and `expect`.
- **TS2.** The test passes only if all of these are true:
  - every fact and term under `derived` holds
  - no fact or term under `notDerived` holds
  - the facts break exactly the constraints under `rejected`
  - every fact obeys its kind ([FC3](fact.md#rules)).

  A test with no `rejected` is the same as a test that expects no broken
  constraints.
- **TS3.** The file MUST import each link, kind and constraint that the test
  names.
- **TS15.** A term under `derived` or `notDerived` MUST have one of the
  forms in [derived terms](#derived-terms). A test MAY list only terms. A
  test in a file with no facts MUST list only terms under `expect`, and no
  `rejected`.

## Derived terms

The definitions derive information that is not a fact: what follows for
kinds and links before there are any facts. A test can expect this
information under `derived` and `notDerived`, next to facts or alone:

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

| Term | Written | Derived from |
|---|---|---|
| a kind is another kind | `kind`, `is a` | `extends`, in one or more steps, and the kinds that a [union](kind.md#unions) lists |
| two kinds never share a thing | `kind`, `never is` | a union |
| a kind can have a link | `kind`, `can have` | the `from` of the link, or the `to` for an inverse name, with `extends` and unions |
| a link implies another link | `link`, `implies` | `under`, in one or more steps |
| a link chains | `link`, `chains: true` | `chains` |
| a link is another link in the opposite direction | `link`, `inverse of` | the inverse name |

A term names kinds and links only, never ids of facts.

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

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Question` |
| `asks` | yes | the question, in plain words |
| `given` | no | `file:` and an example file with the facts to ask about; if not given, the facts in this file |
| `match` | yes | the things and the links to find |
| `filter` | no | conditions that narrow the matches |
| `expect` | yes | for each listed variable, the ids that it must represent |

Each item in `match` is a kind item or a link item. A kind item, such as
`kind: { name: Dish, as: $dish }`, gives a variable to the things of a kind:
the variable represents each thing of that kind in turn. Each variable has
exactly one kind item. A link item is one link that must exist, with
`link`, `from` and `to`. `from` and `to` are each a variable, or a kind
alone for a thing that nothing else refers to. `repeated: true` on a link
item follows the link for one or more steps. `match` finds each combination
of things for which all its items hold. Each combination is one match, which
gives one thing to each variable. An end with a kind alone only has to
exist.

`filter` keeps the matches for which each condition holds. A condition names
a `node`, an optional `field`, a comparison word and what to compare with:

| Word | Takes |
|---|---|
| `equals`, `does not equal`, `is above`, `is below`, `is at least`, `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set`, `is not set` | a `field`, and nothing more |

A condition can also contain a nested block with its own `match` and
`filter`. The words for this are `at least one`, `none`, or `count` with a
comparison word. `any of` holds when one of its listed conditions holds.

A condition can also name a `part` of an `object` field, as in
[constraints](constraint.md#filter-and-unless).

- **TS4.** A question MUST have an `asks` sentence, and the sentence MUST NOT
  contain a `$` word.
- **TS5.** Every item in `match` MUST be a kind item or a link item. A kind
  item MUST have `kind` with `name` and `as`. `name` MUST be a kind that the
  file imports, and `as` MUST be a variable. A link item MUST have `link`,
  `from` and `to`, and MAY have `repeated: true`. `link` MUST be a link, or
  an inverse name, that the file imports.
- **TS6.** Each `from` and each `to` MUST be a variable or a kind. A
  variable there MUST have a kind item, in the same `match` or in a
  `match` that encloses its nested block. The kind at each end MUST be the kind of
  the link at that end, or a kind that extends it. For a variable, this is
  the kind in its kind item.
- **TS7.** Conditions MUST write a variable alone (`$allergen`), and every
  variable that they use MUST have a kind item. A variable of a nested block
  MUST NOT be used outside that block.
- **TS8.** A comparison MUST compare values of the same type: two fields that
  name the same value type, or two fields that name no value type and have
  the same built-in type. A literal MUST be a value that the field allows. A
  literal compared with a node without a field MUST be the id of a fact. A
  field MUST be a field of the kind of its node, or of a kind that it
  extends. Two nodes that a condition
  compares without fields MUST have the same kind, or kinds where one
  extends the other.
- **TS9.** Each key under `expect` MUST be a variable of the question's
  `match`.
- **TS10.** The question passes only if, for each listed variable, the
  distinct things that it represents in all matches are exactly the listed
  ids.
- **TS11.** `part` MAY appear only next to `field`. Each name in `part`
  MUST be a part that the schema of the field defines. For a `oneOf`, one
  or more of the choices MUST define it. The comparison uses the type of
  that part.
- **TS12.** A question MUST have `match`, with at least one item. Each
  nested block MUST also have `match`, with at least one item.
- **TS13.** Each variable MUST have exactly one kind item. Two kind items
  in one question MUST NOT give the same variable.
- **TS14.** Each variable MUST be used outside its kind item: at two or
  more places, or in a condition or `expect`. A variable at one link end
  and nowhere else refers to nothing, so that end gives the kind alone.
- **TS16.** A condition MUST be one of the forms above. `is set` and `is
  not set` MUST name a field. `is a` MUST name a kind.
- **TS17.** `count` MUST have one comparison word: `equals`, `does not
  equal`, `is above`, `is below`, `is at least` or `is at most`. Its value
  MUST be `value:` with a whole number, 0 or more.