---
title: Kinds, fields and value types — The OpenFjord modeling language
---

# Kinds, fields and value types

> Part of [the reference](README.md).

A kind is a category of things, for example a station, a recipe or a
dish. A kind tells what the thing is and gives examples of it. It also
lists the fields of the thing, and the links that the thing has, with a count
for each link.

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
  links:
    madeAtStation: exactly 1
    usesRecipe: any number
    usesIngredient: any number
    usesGarnish: any number
    hasAllergen: any number

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
  links:
    listsAllergen: any number
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Kind` |
| `definition` | yes | one sentence that tells what the kind is |
| `extends` | no | the broader kind; this kind is a narrower sort of it |
| `examples` | no | a list of things of this kind, in words |
| `notExamples` | no | a list of things that look like things of this kind, but are not |
| `fields` | no | the [fields](#fields) of the kind, by name |
| `links` | no | the links that a thing of this kind has, each with a count |
| `tags` | no | the [tags](tag.md) that the kind belongs to |

## Kind rules

- **KD1.** A kind MUST have a `definition`.
- **KD2.** `extends` MUST name a kind. A kind MUST NOT extend itself,
  directly or through other kinds.
- **KD3.** When kind B extends kind A:
  - B has every field and every link count of A.
  - B MAY add fields and link counts of its own.
  - Each thing of kind B is also a thing of kind A.

  For example, `Dish` extends `Recipe`. A dish has the fields and links of a
  recipe, and also a price. Each dish is also a recipe.
- **KD4.** Each key under `links` MUST be a link, or the inverse name of a
  link, that starts at this kind or at a kind that this kind extends. For an
  inverse name, the link must end at that kind.
- **KD5.** Each count MUST be one of `exactly N`, `at least N`, `at most N`
  or `any number`, where N is a whole number.

## What counts mean

A count is a check. A recipe that names two stations is an error. A fact cannot state a link that its kind does not list under
`links`.

Every field is required, unless it has `optional: true` or a `default`
([FD3](#field-rules)). These two rules together mean that a fact must have
all the items that its kind asks for, and nothing more
([FC3](fact.md#rules)).

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
| `optional: true` | a fact can omit the field |

There are five differences from JSON Schema. Each difference agrees with a
decision in a different part of the language:

| JSON Schema | Here | Why |
|---|---|---|
| `required: [ … ]` lists required parts | every field and every part of an `object` is required, unless it has `optional: true` | required is the default in all places |
| `$ref` refers to a different schema | `type: Name` names an imported value type | names come from imports |
| `null` is a type | there is no null type | a field with no value has no value; there is no third state |
| `oneOf`, `anyOf`, `allOf`, `not`, `if`/`then`/`else` combine schemas | not available | difficult to read; conditional logic goes in constraints |
| `additionalProperties` controls if parts that are not listed are allowed | always false | the language rejects all items that are not listed |

`number` is exact, so `0.1 + 0.2` is `0.3`. It is exported as
`xsd:decimal`, never as a floating-point type.

### Field rules

- **FD1.** A field MUST use only the keywords in the table above.
- **FD2.** `required`, `$ref`, `null`, `oneOf`, `anyOf`, `allOf`, `not`,
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

## Value types

A value type gives a name to a schema one time, so that many fields can use
the schema again. A value type is a statement like all other statements, and
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

  `Minutes` and a plain `integer` are different types. A field that
  narrows a value type ([VT3](#value-types)) keeps that value type.
