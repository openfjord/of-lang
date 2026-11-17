---
title: Packages — The OpenFjord modeling language
---

# Packages

> Part of [the reference](../README.md#reference).

A package is the unit that people version, own, review and depend on. It
holds definitions: kinds, links, value types, constraints and derivatives. It
also holds the docs that explain these definitions and the tests that prove
them.

## Layout

```
restaurant/
  package.yaml                     # the manifest
  defs/                            # definitions
    names.of.yaml                     # value types
    station.of.yaml                   # Station, its links, the derivative about stations
    recipe.of.yaml                    # Recipe and its kinds, their links, the derivative about recipes
    ingredient.of.yaml                # Ingredient, Allergen, and the link between them
    allergens.of.yaml                 # a constraint and derivatives that concern several kinds
  docs/                            # Markdown that explains the definitions
    allergens.md
  tests/                           # examples and questions
    examples/
    questions/
    terms/                            # examples with no facts: derived terms
```

- **PK1.** The top level of a package MUST hold the manifest and `defs/`, MAY
  hold `docs/` and `tests/`, and MUST NOT hold anything else. The manifest
  is a YAML file named `package.yaml` or `package.yml`, in every form of the
  language. A package MUST NOT hold both.
- **PK2.** In `defs/`, `docs/` and `tests/`, the author chooses the layout.
  Files in `defs/` and `tests/` have names that end in `.of.`, then the
  name of the form of the package, such as `.of.yaml`
  ([FL15](files.md#files-and-statements)). A file in `defs/` MAY hold
  definitions of any sort, in any combination.
- **PK3.** In `defs/`, a definition SHOULD be next to the item that it is
  about:
  - a link with the kind where the link starts
  - a constraint or derivative with the kind that it concerns
  - a constraint or derivative that concerns several kinds one level higher,
    next to all of those kinds.

## The manifest

```yaml
package: "@bistro/kitchen"        # namespace, then name; quoted, since YAML reserves a leading @
version: 0.1.0
owner: kitchen                    # who must approve any change to the package
rdfNamespace: kitchen             # used only in exported RDF, never in the files
dependsOn:
  "@bistro/suppliers": ^1.2       # semantic-version range, pinned in a lock file
exports:                          # what other packages may use; the rest is internal
  kinds: [Station, Recipe, Dish, BaseRecipe, Ingredient, Allergen]
  links: [madeAtStation, usesRecipe, usesIngredient, usesGarnish, containsAllergen, hasAllergen, listsAllergen, dependsOnStation]
  valueTypes: [Name, Minutes, Price]
openFor:                          # extension points: what other packages may add to
  - extend: Recipe                #   another package may declare a kind of recipe
  - under: usesIngredient         #   another package may declare a narrower link under usesIngredient
tags:                             # the labels that statements of this package can carry
  - name: allergens
    description: What an allergen app needs to tell guests which dishes contain which allergens.
  - name: menu
    description: What the menu app needs to show the dishes and their prices.
```

| Key | Required | Value |
|---|---|---|
| `package` | yes | the full name of the package |
| `version` | yes | a semantic version, `MAJOR.MINOR.PATCH` |
| `owner` | yes | who must approve any change to the package |
| `rdfNamespace` | yes | the prefix that the tool puts in front of each name in the exported RDF |
| `dependsOn` | no | a mapping from package names to version ranges |
| `exports` | no | the `kinds`, `links` and `valueTypes` that other packages can import |
| `openFor` | no | a list of extension points, each `extend: <kind>` or `under: <link>` |
| `tags` | no | a list of [tags](tag.md), each with a `name` and a `description` |

- **PK4.** The manifest MUST have `package`, `version`, `owner` and
  `rdfNamespace`, and MUST NOT have a key that is not in the table.
- **PK5.** `version` MUST be `MAJOR.MINOR.PATCH`, three whole numbers.
- **PK15.** Each item in `tags` MUST have a `name` and a `description`, and
  no other key. No two items MAY have the same `name`.

## Names and namespaces

A package name can start with a namespace, as in npm: `@bistro/kitchen`
is the package `kitchen` in the namespace `@bistro`.

- **PK6.** A package name MUST be lowercase letters, digits and hyphens,
  with a letter first. An optional namespace can come before the name: `@`,
  then the same characters, then `/`.
- **PK7.** `dependsOn` and `pkg:` imports MUST use the full name, with the
  namespace.

## Owners

A package has one owner. That owner owns everything in the package:
definitions, docs and tests. The owner must approve each change to the package
(reasons 7 and 8 in [the research](../../README.md)). The users choose which
packages hold which terms.

- **PK8.** No item in a package has its own owner. Thus, terms with different
  owners MUST be in different packages.

## Dependencies

A version range is one of:

| Range | Means |
|---|---|
| `1.2.3` | exactly 1.2.3 |
| `^1.2`, `^1.2.3` | at least that version, below the next major version (below the next minor version when the major version is 0) |
| `~1.2`, `~1.2.3` | at least that version, below the next minor version |

- **PK9.** The tool MUST find every package in `dependsOn`, and the version
  of each package MUST be in the range.
- **PK10.** Packages MUST NOT depend on each other in a circle.

## Exports

The definitions of a package are internal if the package does not export
them.

- **PK11.** Every name in `exports` MUST be a definition of the package, and
  of the sort that its list shows. The inverse name of a link MUST NOT be in
  `exports`. The package exports it with its link.
- **PK12.** A package MUST state things only about its own definitions,
  except through the extension points that its dependencies open in
  `openFor`.
- **PK13.** Every name in `openFor` MUST be a definition of the package:
  `extend` identifies a kind, and `under` identifies a link. `extend` MUST
  NOT identify a kind with `oneOf` ([KD9](kind.md#unions)).
- **PK14.** A new version MUST match the change:
  - A change that removes or renames a definition, or makes a constraint
    stricter, needs a new major version.
  - A change that adds a definition or makes a constraint less strict does not
    need a new major version.

## Docs

Markdown files in `docs/` give the reasons, the methods and the context
(reason 5). A doc has a Markdown link to a definition or embeds it. Thus, a
reader sees the text and the definitions together:

````markdown
# Allergens

An allergen can cause a dangerous reaction in some persons. A recipe contains
each allergen of its ingredients:

![](../defs/allergens.of.yaml#allergenFromIngredient)

In the past, a guest became ill. The pesto in a dish contained nuts, and the
menu did not tell the guest. Thus, the menu must give a warning about each
allergen in a dish. Refer to
[dishListsItsAllergens](../defs/allergens.of.yaml#dishListsItsAllergens).
````

| Written | Means |
|---|---|
| `[text](path/file.of.yaml#Name)` | a Markdown link to the statement `Name` |
| `![](path/file.of.yaml#Name)` | the statement `Name`, embedded at that position |
| `![](path/file.of.yaml)` | the full file, embedded |

Both forms are standard Markdown. Thus, a viewer that does not know this
convention still shows a Markdown link that works. The part after `#` is the
name of a statement, not a line number. Thus, an embed always shows its
definition, also when the definition moves in the file. The generated page of
each definition lists the docs that have a Markdown link to it or embed it.

- **DC1.** Paths MUST be relative to the doc.
- **DC2.** A Markdown link or an embed MUST point to a file that exists. If
  there is a part after `#`, it MUST be the name of a statement in that file.
