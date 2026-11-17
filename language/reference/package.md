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
  package.lock.yaml                # the exact versions of the dependencies; the tool writes it
  defs/                            # definitions
    names.of.yaml                     # value types
    station.of.yaml                   # Station, its links, the derivative about stations
    recipe.of.yaml                    # Recipe and its kinds, their links, the derivative about recipes
    ingredient.of.yaml                # Ingredient, Allergen, and the link between them
    allergens.of.yaml                 # the allergen topic: its links, derivatives and constraint
  docs/                            # Markdown that explains the definitions
    allergens.md
  tests/                           # examples and questions
    examples/
    questions/
    terms/                            # examples with no facts: derived terms
```

- **PK1.** The top level of a package MUST hold the manifest and `defs/`, MAY
  hold the lock file, `docs/` and `tests/`, and MUST NOT hold anything else. The manifest
  is a YAML file named `package.yaml` or `package.yml`, in every form of the
  language. A package MUST NOT hold both.
- **PK2.** In `defs/`, `docs/` and `tests/`, the author chooses the layout.
  Files in `defs/` and `tests/` have names that end in `.of.`, then the
  name of the form of the package, such as `.of.yaml`
  ([FL15](files.md#files-and-statements)). A file in `defs/` MAY hold
  definitions of any sort, in any combination.
- **PK3.** In `defs/`, a definition SHOULD be in the file of its topic:
  - a link, constraint or derivative that concerns one kind, with that kind
  - definitions that serve one topic across several kinds, such as
    allergens, together in one file for that topic.

## The manifest

```yaml
# package.yaml
package: "@bistro/kitchen"        # namespace, then name; quoted, since YAML reserves a leading @
version: 0.1.0
owner: kitchen                    # who must approve any change to the package
rdfNamespace: kitchen             # used only in exported RDF, never in the files
dependsOn:
  "@bistro/suppliers": ^1.2       # semantic-version range; package.lock.yaml pins the exact version
exports:                          # what other packages may use; the rest is internal
  kinds: [Station, Recipe, Dish, BaseRecipe, Ingredient, Allergen]
  links: [madeAtStation, usesRecipe, usesIngredient, usesGarnish, containsAllergen, hasAllergen, listsAllergen, dependsOnStation]
  valueTypes: [Name, Minutes, Price, URL]
openFor:                          # extension points: what other packages may add to
  - extend: Recipe                #   another package may declare a kind of recipe
  - under: usesIngredient         #   another package may declare a narrower link under usesIngredient
tags:                             # the labels that definitions of this package can carry
  - name: allergens
    description: What an allergen app needs to tell guests which dishes contain which allergens.
  - name: spoilage
    description: What the kitchen needs to keep food safe until the cooks use it. Only the kitchen uses it, so the package does not export it.
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
  no other key. Two items MUST NOT have the same `name`.
- **PK16.** The lock file is a YAML file named `package.lock.yaml` or
  `package.lock.yml`. It gives the exact version of each package that the
  package depends on, directly or through other packages. The tool writes
  it, and a person MUST NOT edit it. The tool uses the lock file only when it
  checks, tests or exports this package itself. When a different package
  depends on this package, the lock file of that different package decides
  the versions.

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
(reasons 7 and 8 in [the design notes](../../docs/design.md#why-ontologies-are-hard-for-software-people)). The users choose which
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
- **PK12.** A package MAY use the definitions of its dependencies: as the
  ends of its own links, in its own unions, and in the `match` and
  conditions of its own constraints and derivatives. It MAY extend a kind,
  or declare a link `under` a link, only where the dependency opens it in
  `openFor`. A derivative MUST NOT add a link, a kind or a field value of a
  different package.
- **PK13.** Every name in `openFor` MUST be a definition of the package:
  `extend` identifies a kind, and `under` identifies a link. `extend` MUST
  NOT identify a kind with `oneOf` ([KD9](kind.md#unions)).
- **PK14.** A new version MUST be at least the step that the change needs.
  The tool reads the [pending change](pending.md) and tells the smallest
  step, by what the change can do to the facts of a package that depends on
  this one:

  | Step | The change can | Examples |
  |---|---|---|
  | major | reject facts that passed, change what follows, or remove a name that other packages use | remove, rename or stop exporting a definition; add a constraint; change the `match`, `filter` or `unless` of a constraint; add or change a derivative; add a required field; narrow a field, a value type, or the `from` or `to` of a link; add `under` or `chains` to a link |
  | minor | add names, and nothing more | add a kind, a link, a value type or an optional field; add a tag; remove a constraint; export a definition |
  | patch | change only text | `definition`, `reads`, `says`, `message`, `examples`, `notExamples`, docs, tests |

  A derivative is major because the facts that it adds can make the
  constraints of a different package reject. Removing a constraint is minor,
  because then the package only rejects less.

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

### Doc definitions

A `Doc` gives a doc a name among the definitions of the package. A kind
lists the docs that tell more about it in `docs`. Thus, the generated page of
the kind can show its docs, and a reader of the kind can find them:

```yaml
# defs/allergens.of.yaml
AllergenSafety:
  is: Doc
  path: file:../docs/allergens.md

# defs/ingredient.of.yaml
Allergen:
  is: Kind
  tags: [allergens]
  definition: A substance in food that can cause a dangerous reaction in some persons.
  examples:
    - nuts
    - milk
  fields:
    name:
      type: Name
  docs: [AllergenSafety]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Doc` |
| `path` | yes | `file:` and the path of a Markdown file in `docs/`, relative to this file |

A doc that only other docs refer to is internal: it needs no `Doc`. For
example, a long doc can move one of its parts to a file of its own, and link
to it.

- **DC1.** Paths MUST be relative to the doc.
- **DC2.** A Markdown link or an embed MUST point to a file that exists. If
  there is a part after `#`, it MUST be the name of a statement in that file.
- **DC3.** A `Doc` MUST have `path`, with `file:` and the path of a
  Markdown file (`.md`) in the `docs/` folder of its package. The path is
  relative to the file of the `Doc`, and the file MUST exist.
- **DC4.** Two `Doc` definitions MUST NOT have the same file.
- **DC5.** Each Markdown file in `docs/` MUST be the file of a `Doc`, or the
  target of a Markdown link or an embed in a different doc in `docs/`.
- **DC6.** A `Doc` is internal to its package. `exports` does not list it.
