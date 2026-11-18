---
title: Cheat sheet — The OpenFjord modeling language
---

# Cheat sheet

> Part of [the OpenFjord modeling language](README.md). All of the syntax is
> on one page, in the default YAML form. [The reference](README.md#reference) gives the exact rules.

## A package

The name of every file in `defs/` and `tests/` ends in `.of.yaml` or `.of.yml`
(in the JSON form, `.of.json`; a package uses one form). The manifest is always YAML:
`package.yaml` or `package.yml`.

```
my-package/
  package.yaml       # required: the manifest
  defs/              # required: Kind, Link, ValueType, Constraint, Derivative, Doc
  docs/              # optional: Markdown that links to or embeds definitions; each file has a Doc, or other docs refer to it
  tests/             # optional: Fact, Test, Question
```

```yaml
# package.yaml
package: "@bistro/kitchen"         # required; quoted because of the leading @
version: 0.1.0                     # required
owner: kitchen                     # required; owns everything in the package
rdfNamespace: kitchen              # required; used only in exported RDF
dependsOn:
  "@bistro/suppliers": ^1.2
exports:
  kinds: [Recipe]
  links: [madeAtStation]
  valueTypes: [Name]
openFor:
  - extend: Recipe
  - under: usesIngredient
tags:                                # each label that definitions can carry
  - name: allergens
    description: What an allergen app needs to tell guests which dishes contain which allergens.
```

## A file

```yaml
import:
  - from: file:names.of.yaml                     # relative to this file
    use: [Name]
  - from: pkg:@bistro/suppliers                  # must be in dependsOn and exported there
    use: [Supplier, { name: Ingredient, as: SupplierIngredient }]

SomeStatement:
  is: Kind        # or Link, ValueType, Constraint, Derivative, Doc — Fact, Test, Question in tests/
```

The import schemes are `file:` and `pkg:`. A later version can add more.
A file imports each name that it uses but does not define, and no other
name. A file does not import field names.

All free-form text (`definition`, `reads`, `examples`, `says`, `asks`,
`message`, docs) follows ASD-STE100 Simplified Technical English.

## Definitions

```yaml
Recipe:                                  # UpperCamelCase
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.   # required
  tags: [allergens]                      # optional; names from package.yaml; on any definition but a Doc
  examples:    [the tomato sauce]
  notExamples: [a new dish that a cook tries one time and does not keep]
  fields:
    name:        { type: Name }                       # required by default
    steps:       { type: URL }                        # where the steps are written
    prepMinutes: { type: Minutes, optional: true }
  docs: [KitchenGuide]                   # optional; Docs of this package

KitchenGuide:                            # UpperCamelCase
  is: Doc
  path: file:../docs/kitchen.md          # a Markdown file in docs/, relative to this file

Dish:
  is: Kind
  extends: Recipe                        # optional
  definition: A recipe that the restaurant sells to guests.
  fields:
    price:  { type: Price }
    status: { enum: [active, retired], default: active }

Food:
  is: Kind
  definition: Something that a recipe uses to make its food.
  oneOf: [Recipe, Ingredient]            # a union: each food is exactly one; closed; no fields or links

usesRecipe:                              # lowerCamelCase, two words or more
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"              # required, with both {from} and {to}
  definition: The first recipe needs the second recipe, directly or through other recipes.
  chains: true                           # optional: a→b and b→c means a→c
  inverse: usedByRecipe                  # optional: the name read the other way

usesGarnish:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} is garnished with {to}"
  definition: The cook puts the ingredient on the food at the end, to make it look or taste better.
  under: usesIngredient                  # optional: every usesGarnish is a usesIngredient

Name:
  is: ValueType
  definition: A lowercase identifier for stations, recipes, ingredients and allergens.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"

Minutes:
  is: ValueType
  definition: A number of minutes of work.
  type: integer                          # a type and its keywords
  minimum: 1
  maximum: 1440

URL:
  is: ValueType
  definition: The address of a page on the web.
  type: string
  format: uri                            # an absolute URI, such as https://…
```

By default, the RDF export holds everything. An optional list of tags limits it. Tags do not pass through
`extends` or `under`. The export does not include a reference to an untagged
definition. It also does not include a link whose `from` or `to` kind is
untagged.

Field keywords: `type` (`string` `integer` `number` `boolean` `array` `object`
or a value type), `enum` `const` `default` `pattern` `minLength` `maxLength`
`format` (`date` `date-time` `duration` `uri`) `minimum` `maximum` `items`
`properties` `optional` `oneOf` (variants: exactly one choice matches; object
choices carry a `const` part such as `variant`). The language does not have
`null`, `$ref`, `required`, `anyOf`, `allOf`, `not` or `if`.

## Match and filter

Derivatives, constraints and questions each define their own `match` and
`filter`. All three use the same syntax, but each can change apart from the
others.

```yaml
match:                                   # things and links that must exist
  - kind: { name: Dish, as: $dish }      # each variable has exactly one kind item
  - kind: { name: Recipe, as: $part }
  - link: makesRecipe                    # a link or an inverse name
    from: Station                        # a kind alone: nothing else refers to this end
    to:   $dish
  - link: usesRecipe
    from: $dish
    to:   $part                          # usesRecipe chains, so this also finds longer chains
filter:                                  # conditions; all must hold
  - node: $part
    does not equal: { same as: $dish }
  - node: $dish
    field: status
    equals: { value: active }
```

`repeated: true` on a link item follows a link that does not chain for one or
more steps.

A match can hold only a kind item. The variable then represents each thing
of the kind:

```yaml
match:
  - kind: { name: Dish, as: $dish }      # each dish in turn
filter:
  - node: $dish
    field: price
    is above: { value: 30 }
```

| Comparison | Takes |
|---|---|
| `equals` `does not equal` `is above` `is below` `is at least` `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set` `is not set` | (only a field) |

A condition on an object field can add `part: variant`, or `part: [a, b]`
for a part in a part.

| Nested | Holds when |
|---|---|
| `at least one: { match, filter }` | it matches at least once |
| `none: { match, filter }` | it never matches |
| `count: { match, filter }` + a comparison | the number of matches satisfies the comparison |
| `any of: [conditions]` | one of them holds |

## Constraints and derivatives

```yaml
dishListsItsAllergens:
  is: Constraint
  says: The menu must give a warning about each allergen in an active dish.   # no $ words
  match: [ ... ]
  filter: [ ... ]
  unless: [ ... ]                        # optional; a match is a violation unless these hold
  message: "The menu must warn that {$dish} contains {$allergen}."

allergenFromIngredient:
  is: Derivative
  says: A recipe contains each allergen that one of its ingredients contains.
  match: [ ... ]
  derive:
    - { link: hasAllergen, from: $recipe, to: $allergen }          # a link
    - { node: $recipe, becomes: BaseRecipe }                       # a kind
    - { node: $recipe, field: prepMinutes, set to: { value: 30 } } # a field value
```

## Tests

```yaml
# tests/examples/some-example.of.yaml — facts plus exactly one Test
dish/margherita:                         # a fact's id
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12
  links:
    madeAtStation: station/pizza
    usesRecipe: [recipe/tomato_sauce, recipe/pesto]   # a link: an id, or a list of ids

someExample:
  is: Test
  says: What this example shows.
  expect:
    derived:    [ { link: hasAllergen, from: dish/margherita, to: allergen/nuts } ]
    notDerived: [ { node: recipe/pesto, becomes: BaseRecipe } ]
    rejected:   [dishListsItsAllergens]

# tests/terms/garnish.of.yaml — no facts: what the definitions derive about kinds and links
garnishIsAnIngredient:
  is: Test
  says: A garnish is one of the ingredients that a recipe uses.
  expect:
    derived:
      - { kind: Dish, is a: Recipe }                  # also: never is, can have
      - { link: usesGarnish, implies: usesIngredient } # also: chains, inverse of
    notDerived:
      - { link: usesIngredient, implies: usesGarnish }

# tests/questions/which-dishes-contain-nuts.of.yaml
whichDishesContainNuts:
  is: Question
  asks: Which dishes contain nuts?
  given: file:../examples/allergen-through-a-sauce.of.yaml
  match: [ ... ]
  filter: [ ... ]
  expect:
    $dish: [dish/margherita]
```

## Pending changes

```yaml
Garnish:
  is: Kind
  pending: added                         # or removed, or updated
  definition: An ingredient that the cook puts on the food at the end.

Recipe:
  is: Kind
  pending: updated
  definition: The steps that make one dish or one part of a dish.   # the new version
  was:                                   # the whole old version, without is
    definition: The steps that make one food from ingredients and other recipes.
```

A package has one pending change or none. Accept applies all marks in one
step, and reject undoes them. Other packages see the package without the
change.

## How it runs

The derivatives run repeatedly until no new fact follows. `under` and `chains`
also add their facts. Then the constraints check everything. Anything that is
not stated or derived is false. Every derived fact and every rejection has
its explanation.
