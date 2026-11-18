---
title: Kinds — Introduction to the OpenFjord modeling language
---

# Kinds

> Part of [the introduction](README.md).

A small restaurant has a kitchen. The kitchen has stations, where cooks make
recipes. Some recipes are dishes, which the restaurant sells. Recipes use
ingredients, and some ingredients contain allergens.

Each cook knows these words, but nobody wrote down what they mean. The
kitchen starts there: it gives a name and a definition to each kind of thing.

## A package

A model is in a package. A package is a folder with a manifest and a `defs/`
folder:

```
step-01/             # the root of the package
  package.yaml       # the manifest
  defs/              # the definitions
    names.of.yaml
    station.of.yaml
    recipe.of.yaml
    ingredient.of.yaml
```

The manifest is the file `package.yaml`:

```yaml
# package.yaml
package: "@bistro/kitchen"
version: 0.1.0
owner: kitchen
rdfNamespace: kitchen
```

The manifest must have these four keys:

- `package` is the name of the package. `@bistro` is a namespace, and
  `kitchen` is the name in that namespace.
- `version` is the version of the package.
- `owner` is who must approve each change to the package.
- `rdfNamespace` is a prefix that the tool uses only when it exports the
  package to RDF (Resource Description Framework), a standard format for
  models.

[Chapter 6](06-packages.md) tells more about the manifest.

## The first kind

The definitions of the package are in the files in `defs/`. The name of each
file ends in `.of.yaml`. The kitchen chooses how many files it has and what
goes in each one. It writes the file for stations first:

```yaml
# defs/station.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name]

Station:
  is: Kind
  definition: A part of the kitchen where a group of cooks makes some of the recipes.
  examples:
    - the pizza station
    - the pastry station
  fields:
    name:
      type: Name
```

`Station` is a statement. A file holds statements, and the key of each
statement is its name. The `is` key tells what the statement is. This one is
a kind: a category of things.

- `definition` is one sentence that tells what a station is. Each kind must
  have one.
- `examples` lists some stations, in words. The examples help a reader to
  understand the definition.
- `fields` lists the values that each station has. A station has one field,
  its `name`.

The name of a kind starts with a capital letter, and so does each word in
it, such as `Station`.

## Value types and imports

The `name` field has the type `Name`. `Name` is not a part of the language.
It is a value type: a description of a value that has a name of its own, so
that many fields can use it. The kitchen defines its value types in one
file:

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

URL:
  is: ValueType
  definition: The address of a page on the web.
  type: string
  format: uri

Price:
  is: ValueType
  definition: An amount of money in euros.
  type: number
  minimum: 0
```

A value type has a built-in `type`, such as `string`, `integer` or `number`,
and keywords that narrow it. A `Name` is text that matches a pattern and has
63 characters at most. A value of `Minutes` is a whole number from 1 to
1440. `format: uri` allows only an absolute URI (Uniform Resource
Identifier), such as `https://recipes.example.com/pesto`.

A file can use a name only if the file defines the name or imports it. Thus,
`station.of.yaml` starts with an `import` block. `from` gives the file that
defines the names, relative to the file that imports them. `use` lists the
names. The top of a file shows all that the file depends on.

## A kind with more fields

A recipe has a name, the address of the page with its steps, and the number
of minutes that the work takes:

```yaml
# defs/recipe.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name, Minutes, Price, URL]

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
    steps:
      type: URL
    prepMinutes:
      type: Minutes
      optional: true
```

`notExamples` lists things that look like recipes but are not. Together,
`examples` and `notExamples` show where the kind starts and where it stops.

A field is required by default: each recipe must have a `name` and `steps`.
`optional: true` makes a field optional, so a recipe can be without
`prepMinutes`.

## A narrower kind

A dish is a recipe that the restaurant sells. The kitchen does not write the
fields of a recipe a second time. `Dish` extends `Recipe`:

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

A dish has each field of a recipe, and two fields of its own. Each dish is
also a recipe.

The `status` field does not name a type. `enum` lists the values that it can
have. A dish that gives no status has the `default`, `active`.

The last file defines the two kinds that are left. An ingredient is a food
that the restaurant buys, and an allergen is a substance in food:

```yaml
# defs/ingredient.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name]

Ingredient:
  is: Kind
  definition: A food that the restaurant buys and does not make.
  examples:
    - flour
    - pine nuts
  notExamples:
    - the tomato sauce, which the kitchen makes
  fields:
    name:
      type: Name

Allergen:
  is: Kind
  definition: A substance in food that can cause a dangerous reaction in some persons.
  examples:
    - nuts
    - milk
  fields:
    name:
      type: Name
```

The package has five kinds, but it does not tell how they relate to each
other. That is the subject of the next chapter: [Links](02-links.md).

More about this:
[kinds, fields and value types](../../language/reference/kind.md),
[files and statements](../../language/reference/files.md) and
[the manifest](../../language/reference/package.md#the-manifest) in the
reference. In the language by example:
[Hello, kind](../by-example/hello-kind.md), [Fields](../by-example/fields.md),
[Value types](../by-example/value-types.md),
[Defaults](../by-example/defaults.md), [Enums](../by-example/enums.md),
[Extends](../by-example/extends.md) and [Imports](../by-example/imports.md).
