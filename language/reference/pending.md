---
title: Pending changes — The OpenFjord modeling language
---

# Pending changes

> Part of [the reference](README.md).

A pending change is a change to a package that nobody has accepted yet. The
change is in the files of the package, as marks on the statements that it
adds, removes or updates. Thus, the tool and the web app can show what the
change does without git.

```yaml
# defs/recipe.of.yaml
Garnish:                     # the change adds this statement
  is: Kind
  pending: added
  definition: An ingredient that the cook puts on the food at the end.

usesGarnish:                 # the change removes this statement
  is: Link
  pending: removed
  from: Recipe
  to: Ingredient
  reads: "{from} is garnished with {to}"
  definition: The cook puts the ingredient on the food at the end, to make it look or taste better.
  under: usesIngredient

Recipe:                      # the change updates this statement
  is: Kind
  pending: updated
  definition: The steps that make one dish or one part of a dish.
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
    hasAllergen: any number
  was:
    definition: The steps that make one food from ingredients and other recipes.
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
```

| Key | Required | Value |
|---|---|---|
| `pending` | no | `added`, `removed` or `updated` |
| `was` | only with `pending: updated` | the whole statement as it was before the change, without `is` |

An updated statement shows its new version. `was` keeps the old version, so
that the tool can show the difference and can reject the change. In this
example, the update changes the definition of `Recipe`. It also removes the
count for `usesGarnish`, because the change removes that link.

## One change at a time

All the marks in a package together are its pending change. A package has
one pending change or none. The change has no name.

The package has two versions:

| Version | Added statements | Removed statements | Updated statements |
|---|---|---|---|
| without the change | not included | included, as written | the `was` version |
| with the change | included | not included | the new version |

The tool checks the two versions. The [change report](../../cli/README.md#the-change-report)
compares them: what the change states, what newly follows, what no longer
follows, and what is rejected.

Other packages see only the version without the change. A pending change has
no effect outside its package until somebody accepts it.

## Accepting and rejecting

To accept the change, the tool does these steps as one operation:

1. It deletes each removed statement.
2. It deletes `pending` and `was` from each added and updated statement.

To reject the change, the tool does these steps as one operation:

1. It deletes each added statement.
2. It replaces each updated statement with its `was` version.
3. It deletes `pending` from each removed statement.

## Rules

- **PD1.** Any statement in `defs/` or `tests/` MAY have `pending`, with the
  value `added`, `removed` or `updated`.
- **PD2.** A statement with `pending: updated` MUST have `was`. `was` MUST
  hold the whole statement as it was before the change, without `is`,
  `pending` and `was`. A statement with a different `pending` value MUST NOT
  have `was`.
- **PD3.** An update MUST NOT change `is`. To change what a statement is,
  the change removes the statement and adds a different one.
- **PD4.** The version without the change and the version with the change
  MUST each pass every rule in this specification.

The manifest, `package.yaml`, and the Markdown files in `docs/` have no
pending marks. A pending change to them is deferred.
