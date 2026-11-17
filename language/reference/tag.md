---
title: Tags — The OpenFjord modeling language
---

# Tags

> Part of [the reference](../README.md#reference).

A tag is a label that marks definitions that belong together for one
purpose. By default, the export to RDF (Resource Description Framework)
includes the full package. An optional list of tags limits the export to the
definitions with those tags. For example, an allergen app needs recipes,
ingredients, allergens and the links between them. The menu app needs dishes
and their prices.

```yaml
# defs/recipe.of.yaml
Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
  tags: [allergens]
```

The manifest of the package lists its tags, each with a name and a description
([the manifest](package.md#the-manifest)):

```yaml
# package.yaml
tags:
  - name: allergens
    description: What an allergen app needs to tell guests which dishes contain which allergens.
```

Thus, a misspelled tag is an error, not a new subset. Each package has its
own tags. An export with a list of tags selects, in each package of the
export, the definitions that carry a tag with one of those names.

## Tagging a definition

Every definition can have an optional `tags` key, which is a list of tags:
kinds, links, value types, constraints and derivatives. Single fields have no
tags. They are exported with their kind.

```yaml
allergenFromRecipe:
  is: Derivative
  tags: [allergens]
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

A tagged derivative or constraint is in the subset, as a rule of the
export.

Tags select what is exported. They do not change how derivatives,
constraints and tests run.

Tags do not pass from one definition to a different definition. A kind that
extends a tagged kind does not have the tag, unless it lists the tag itself.
The same applies to a link `under` a tagged link.

A package tags only its own definitions. `tags` is a key on the definition
itself, so a package cannot tag the definitions of a different package.

## Exporting a subset

By default, an export includes the full package and its dependencies. An
optional list of tags limits the export: then it includes, in each package,
each definition that has one or more of these tags.

A subset is complete in its references: each definition that a tagged
definition refers to has the same tag ([TG5](#rules)). Without this rule, an
export must drop the reference, which gives wrong answers, or export the
untagged definition, which adds knowledge that nobody tagged. TG5 does not
cover an untagged derivative that adds facts that a tagged definition uses.
How the tool checks that a subset answers its questions as the full package
does is still open (see [the overview](../../README.md#open-questions)).

For example, `Dish` extends `Recipe`. When `Dish` has the tag `allergens`,
`Recipe` must also have it. Otherwise the tool rejects the package and names
`Recipe`. In the same way, when `allergenFromRecipe` has the tag
`allergens`, so must `Recipe`, `Allergen`, `usesRecipe` and `hasAllergen`.

## Rules

- **TG2.** A tag MUST be lowerCamelCase.
- **TG3.** `tags` MAY appear only on a definition: a kind, a link, a value
  type, a constraint or a derivative. Each
  tag in the list MUST be the `name` of a tag in the `tags` of the package
  manifest.
- **TG4.** A definition MUST NOT list the same tag twice.
- **TG5.** Each definition that a tagged definition refers to MUST have a
  tag with the same name. This applies to each tag of the definition, and to
  definitions in other packages. The references are:
  - the `from` and `to` kinds of a link
  - the kind that a kind extends, and the kinds that a union lists
  - the value type of each field
  - the link that a link is `under`
  - in a constraint or a derivative, each kind and each link that it names:
    in `match`, in conditions, in nested blocks and in `derive`.
