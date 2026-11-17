---
title: Kinds, fields and value types — The OpenFjord modeling language
---

# Kinds, fields and value types

> Part of [the reference](../README.md#reference).

A kind is a category of things, for example a station, a recipe or a
dish. A kind tells what the thing is and gives examples of it. It also
lists the fields of the thing.

```yaml
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
    prepMinutes:
      type: Minutes
      optional: true

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

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Kind` |
| `definition` | yes | one sentence that tells what the kind is |
| `extends` | no | the broader kind; this kind is a narrower sort of it |
| `oneOf` | no | the kinds that together make up this kind, with no overlap ([unions](#unions)) |
| `examples` | no | a list of things of this kind, in words |
| `notExamples` | no | a list of things that look like things of this kind, but are not |
| `fields` | no | the [fields](#fields) of the kind, by name |
| `tags` | no | the [tags](tag.md) that the kind belongs to |

## Kind rules

- **KD1.** A kind MUST have a `definition`.
- **KD2.** `extends` MUST name a kind. A kind MUST NOT extend itself,
  directly or through other kinds.
- **KD3.** When kind B extends kind A:
  - B has every field of A.
  - B MAY add fields of its own.
  - Each thing of kind B is also a thing of kind A.

  For example, `Dish` extends `Recipe`. A dish has the fields and links of a
  recipe, and also a price. Each dish is also a recipe.

## Links of a thing

Each [link](link.md) names the kind where it starts and the kind where it
ends. A thing can have each link that starts at its kind, at a kind that its
kind extends, or at a union that lists its kind.

A [constraint](constraint.md) with `count` sets how many links a thing has:

```yaml
recipeAtOneStation:
  is: Constraint
  says: Cooks make each recipe at exactly one station.
  match:
    - {$recipe: Recipe}
  filter:
    - count:
        match:
          - link: madeAtStation
            from: {$recipe: Recipe}
            to:   {$station: Station}
      does not equal:
        value: 1
  message: "{$recipe} must be made at exactly one station."
```

A recipe without a station is an error, and so is a recipe with two
stations. A link with no such constraint can occur any number of times,
or not at all.

Every field is required, unless it has `optional: true` or a `default`
([FD3](#field-rules)).

## Unions

Sometimes a kind is exactly one of a fixed list of kinds, and nothing else.
For example, the food that a recipe uses is a recipe or an ingredient. Each
food is one of the two, and never both. `oneOf` states this:

```yaml
Food:
  is: Kind
  definition: Something that a recipe uses to make its food.
  oneOf: [Recipe, Ingredient]

usesFood:
  is: Link
  from: Recipe
  to: Food
  reads: "{from} uses {to}"
  definition: The recipe needs the food.
```

`Recipe` and `Ingredient` do not extend `Food`, and they do not get anything
from it. `Food` gives a name to the two kinds together. Thus, a link, a
match or a condition can use `Food` where a recipe or an ingredient can
occur.

Each thing of a listed kind is also a thing of the union kind. Thus, where
a rule accepts a kind that extends a different kind, the rule also accepts a
listed kind in place of the union kind. For example, `usesFood` can link a
recipe to an ingredient.

A business menu shows a different use. It is a lunch with fixed courses at
one price. The guest picks one drink: a light drink or a glass of wine.

```yaml
LightDrink:
  is: Kind
  definition: A drink with little or no alcohol, such as juice or water.

Wine:
  is: Kind
  definition: A glass of wine.

Drink:
  is: Kind
  definition: The drink that comes with a business menu.
  oneOf: [LightDrink, Wine]

BusinessLunch:
  is: Kind
  definition: One business menu that a guest orders.

comesWithDrink:
  is: Link
  from: BusinessLunch
  to: Drink
  reads: "{from} comes with {to}"
  definition: The drink that the guest picks for the business lunch.
```

A business lunch can come with a light drink or with a wine. A drink is
never both, because each drink is exactly one listed kind. A business lunch
cannot come with a dish, because a dish is not a drink. A [constraint](constraint.md) with `count` can make sure that each
business lunch comes with exactly one drink
([links of a thing](#links-of-a-thing)).

A union is closed. Only the statement itself lists its kinds. No other kind,
and no other package, can add a kind to it.

- **KD7.** `oneOf` MUST list two or more kinds. No listed kind MAY extend a
  different listed kind, directly or through other kinds, because then each
  thing of the narrower kind is two listed kinds at once.
- **KD8.** A kind with `oneOf` MUST NOT have `fields`. Its listed kinds
  keep their own fields.
- **KD9.** No kind MAY extend a kind with `oneOf`.
- **KD10.** A fact MUST NOT name a kind with `oneOf` as its `kind`. It names
  one of the listed kinds.
- **KD11.** After the derivatives run, each thing of a union kind MUST be a
  thing of exactly one listed kind. The tool rejects a thing that a
  derivative makes a thing of two listed kinds.

## Fields

A field is a value of the thing itself, for example a name, a price or a
number of minutes. A subset of JSON Schema describes each field, and each keyword
has the meaning that JSON Schema gives it. Only field values use JSON Schema.
Kinds, links, constraints and derivatives use the syntax of the language
itself, because JSON Schema describes one document and has no links between
things.

```yaml
fields:
  name:
    type: Name                 # a value type, imported
  status:
    enum: [active, retired]
    default: active
  servings:
    type: integer
    minimum: 1
    optional: true
  plating:
    type: object
    properties:
      garnish: { type: string, optional: true }
      plate:   { type: string }
```

### Keywords

| Keyword | For |
|---|---|
| `type`: `string`, `integer`, `number`, `boolean`, `array`, `object`, or a value type's name | the type of the value |
| `enum`, `const`, `default` | fixed values |
| `pattern`, `minLength`, `maxLength`, `format` (`date`, `date-time`, `duration`) | text |
| `minimum`, `maximum` | numbers |
| `items` | the values of an `array` |
| `properties` | the parts of an `object` |
| `oneOf` | [variants](#variants): the value matches exactly one of the listed schemas |
| `optional: true` | a fact can omit the field |

There are five differences from JSON Schema. Each difference agrees with a
decision in a different part of the language:

| JSON Schema | Here | Why |
|---|---|---|
| `required: [ … ]` lists required parts | every field and every part of an `object` is required, unless it has `optional: true` | required is the default in all places |
| `$ref` refers to a different schema | `type: Name` names an imported value type | names come from imports |
| `null` is a type | there is no null type | a field with no value has no value; there is no third state |
| `oneOf`, `anyOf`, `allOf`, `not`, `if`/`then`/`else` combine schemas | only `oneOf` | `oneOf` gives variants; value types give reuse; conditional logic goes in constraints |
| `additionalProperties` controls if parts that are not listed are allowed | always false | the language rejects all items that are not listed |

`number` is exact, so `0.1 + 0.2` is `0.3`. It is exported as
`xsd:decimal`, never as a floating-point type.

### Variants

Some values have one of a small number of shapes. `oneOf` lists the shapes,
and each value matches exactly one of them. When the shapes are objects,
each shape SHOULD have a part with a fixed value, so that a reader, a
constraint and a derivative can see which shape a value has
([FD8](#field-rules)). Constraints and derivatives read a part with the
[`part` key](constraint.md#filter-and-unless).

A choice between sorts of things, such as a light drink or a wine, is not a
field. It is a [union of kinds](#unions).

### Field rules

- **FD1.** A field MUST use only the keywords in the table above.
- **FD2.** `required`, `$ref`, `null`, `anyOf`, `allOf`, `not`,
  `if`, `then`, `else` and `additionalProperties` MUST NOT appear. The error
  tells what to write instead.
- **FD3.** A field is required, unless it has `optional: true` or a
  `default`. A field SHOULD NOT have both, because a default already makes
  the field optional.
- **FD4.** `pattern` MUST be a valid regular expression. `format` MUST be
  `date`, `date-time` or `duration` (ISO 8601).
- **FD5.** The rest of the field's schema MUST allow each value in `enum`,
  the `default` and the `const`.
- **FD6.** There is no null. A field with no value is omitted, and never has
  the value `null`.
- **FD7.** Each choice in `oneOf` MUST be a schema that uses the keywords
  above, but not `optional` or `default`. A value MUST match exactly one
  choice.
- **FD8.** When the choices of a `oneOf` are objects, each choice SHOULD have
  a part with the same name and a different `const` value.

## Value types

A value type gives a name to a schema one time, so that many fields can use
the schema again. A value type can hold a `oneOf`. Then all the fields that
name it share the same variants. A value type is a statement like all other statements, and
is imported in the same way:

```yaml
# defs/names.of.yaml

Name:
  is: ValueType
  definition: A lowercase identifier for stations, recipes, ingredients and allergens.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"
  maxLength: 63

Minutes:
  is: ValueType
  definition: A number of minutes of work.
  type: integer
  minimum: 1
  maximum: 1440

Price:
  is: ValueType
  definition: An amount of money in euros.
  type: number
  minimum: 0
```

```yaml
# defs/recipe.of.yaml

import:
  - from: file:names.of.yaml
    use: [Name, Minutes, Price]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `ValueType` |
| `definition` | yes | one sentence that tells what the value holds; product's pages show it |
| `tags` | no | the [tags](tag.md) that the value type belongs to |
| any keyword above, but not `optional` | — | the schema |

- **VT1.** A value type MUST have a `definition`.
- **VT2.** A value type MUST NOT use the name of a built-in type (`string`,
  `integer`, `number`, `boolean`, `array`, `object`). Built-in types do not
  need an import.
- **VT3.** A field that names a value type MAY narrow it, but MUST NOT widen
  it. For example, `type: Name` with `maxLength: 30` is allowed, but with
  `maxLength: 100` it is an error. This rule also applies to `minLength`,
  `minimum`, `maximum` and `format`.
- **VT4.** The language checks a value against the value type and against
  the field's own keywords. The value must satisfy both.
- **VT5.** A value type is a type of its own, not only a narrower schema. Two
  fields have the same type only if one of these conditions is true:
  - Both fields name the same value type.
  - Neither field names a value type, and both fields have the same built-in
    type.
- **VT6.** A field that names a value type with `oneOf` MUST NOT add
  keywords to it, other than `optional` and `default`.

  `Minutes` and a plain `integer` are different types. A field that
  narrows a value type ([VT3](#value-types)) keeps that value type.
