---
title: Tags — The OpenFjord modeling language
---

# Tags

> Part of [the reference](README.md).

A tag marks definitions that belong together for one purpose. By default,
the export to RDF (Resource Description Framework) includes the full package.
An optional list of tags limits the export to the definitions with those
tags. For
example, an allergen app needs recipes, ingredients, allergens and the links
between them. The menu app needs dishes and their prices. Each purpose has a
tag, and each definition lists the tags that it belongs to.

```yaml
# defs/tags.of.yaml
allergens:
  is: Tag
  definition: What an allergen app needs to tell guests which dishes contain which allergens.

menu:
  is: Tag
  definition: What the menu app needs to show the dishes and their prices.
```

```yaml
# defs/recipe.of.yaml
import:
  - from: file:tags.of.yaml
    use: [allergens]

Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
  tags: [allergens]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Tag` |
| `definition` | yes | one sentence that tells the purpose of the tagged definitions |

A tag is declared and imported like all other names. Thus, a misspelled tag
is an error, not a new subset. A package may [export](package.md#exports) its
tags, so that other packages can use them on their own definitions.

## Tagging a definition

Kinds, links and value types can have an optional `tags` key, which is a list
of tags. Nothing else has tags. Constraints, derivatives and single fields are
always exported with the full package, never in a subset.

Tags only select what is exported. They have no effect on checks,
derivatives, constraints or tests.

Tags do not pass from one definition to a different definition. A kind that
extends a tagged kind does not have the tag, unless it lists the tag itself.
The same applies to a link `under` a tagged link.

A package tags only its own definitions. `tags` is a key on the definition
itself, so a package cannot tag the definitions of a different package.

## Exporting a subset

By default, an export includes the full package. An optional list of tags
limits the export: then it includes each definition of the package that has
one or more of these tags.

A definition in the subset can refer to a definition that is not in the
subset. Examples of these references are:

- a link whose `from` or `to` kind is untagged
- a kind that extends an untagged kind
- a field whose value type is untagged
- a link `under` an untagged link.

Each such reference is not exported. A link whose `from` or `to` kind is not
exported is not exported at all. The export lists each item that it removes.

## Rules

- **TG1.** A tag MUST have a `definition`.
- **TG2.** A tag's name MUST be lowerCamelCase.
- **TG3.** `tags` MAY appear only on a kind, a link or a value type. Each tag
  in the list MUST be a tag that the file defines or imports.
- **TG4.** A definition MUST NOT list the same tag twice.
