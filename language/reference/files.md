---
title: Files and statements — The OpenFjord modeling language
---

# Files and statements

> Part of [the reference](../README.md#reference).

Packages, the `defs/` and `tests/` folders, files and imports are the same
in every form of the language ([principle 7](../README.md#principles)). Only
the syntax inside a file depends on its form. This page shows the YAML form,
which is the default.

The name of a file ends in `.of.`, then the name of its form: `.of.yaml`
or `.of.yml` for YAML, `.of.json` for JSON. For example, `recipe.of.yaml`. The ending
tells editors and tools which files are in the language and which form they
use. All the files of one package use the same form.

The manifest of a package is different. It is always YAML, whatever the form
of the package, and its name is `package.yaml` or `package.yml`
([the manifest](package.md#the-manifest)).

In the YAML form, a file is a mapping: an optional `import` block, and then
statements. The key of each statement is its name.

- **FL15.** The name of every file in `defs/` and `tests/` MUST end in
  `.of.`, then an ending of a form that the tool reads: `yaml` or `yml` for
  YAML, `json` for JSON. Any other file there
  is an error. The tool does not ignore it.
- **FL18.** All the files in `defs/` and `tests/` of a package MUST use the
  same form.

```yaml
import:
  - from: file:names.of.yaml
    use: [Name]
  - from: file:station.of.yaml
    use: [Station]

Recipe:
  is: Kind
  definition: The steps that make one food from ingredients and other recipes.
  fields:
    name:
      type: Name

madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.
```

## Statements

The `is:` key of every statement tells what the statement is.

At the top level of a statement, every key belongs to the language. The
names that users choose, such as the names of fields and links, are always
one level lower: under `fields` in a kind, and under `fields` and `links`
in a fact. Thus, a
name that a user chooses never clashes with a key of the language. This
includes keys that a later version of the language adds.

| `is:` | Where | Described in |
|---|---|---|
| `Kind` | `defs/` | [Kinds](kind.md) |
| `Link` | `defs/` | [Links](link.md) |
| `ValueType` | `defs/` | [Value types](kind.md#value-types) |
| `Constraint` | `defs/` | [Constraints](constraint.md) |
| `Derivative` | `defs/` | [Derivatives](derivative.md) |
| `Doc` | `defs/` | [Docs](package.md#docs) |
| `Fact` | `tests/` | [Facts](fact.md) |
| `Test` | `tests/` | [Tests](test.md) |
| `Question` | `tests/` | [Tests](test.md#questions) |

- **FL1.** Every top-level key other than `import` MUST be a statement: a
  mapping with an `is:` key that has a value from the table above.
- **FL2.** A statement MUST NOT have a key that its page does not list,
  other than `pending` and `was` ([pending changes](pending.md)). The
  tool rejects all keys that are not listed, so a misspelled key is an error.
  The tool does not ignore it.
- **FL3.** A file MUST NOT repeat a key in the same mapping.
- **FL4.** Statements in `defs/` MUST be `Kind`, `Link`, `ValueType`,
  `Constraint`, `Derivative` or `Doc`. Statements in `tests/` MUST be `Fact`, `Test` or
  `Question`.
- **FL5.** A statement MUST NOT be named `import`.

## Names

| Statement | Written as | Examples |
|---|---|---|
| kind, value type, doc | UpperCamelCase | `Recipe`, `BaseRecipe`, `Minutes`, `AllergenSafety` |
| link, link inverse | lowerCamelCase, two words or more | `madeAtStation`, `usesRecipe`, `makesRecipe` |
| constraint, derivative, test, question | lowerCamelCase | `dishListsItsAllergens`, `whichDishesContainNuts` |
| fact | an id: any text, usually a path | `dish/margherita` |

- **FL6.** The name of a definition MUST be unique in its package, also
  when case is ignored: a package cannot have both `BaseRecipe` and
  `baseRecipe`. The inverse name that a link declares is also a name of that
  package.
- **FL7.** Kind, value type and doc names MUST be UpperCamelCase; constraint,
  derivative, test and question names MUST be lowerCamelCase.
  [LK1](link.md#rules) gives the rule for link names.
- **FL8.** The name of a statement in a test file MUST NOT be the name of a
  definition in the package, also when case is ignored. Statements in a test file are local to that file.

## Imports

A file can use a name only if the file defines or imports that name. The
`import` block is at the top of the file. It shows the source of each name:

```yaml
import:
  - from: file:recipe.of.yaml             # a file in this package, relative to this one
    use: [Recipe, madeAtStation]
  - from: file:../station.of.yaml         # from defs/kitchen/: a file in defs/
    use: [Station]
  - from: pkg:@bistro/suppliers           # another package
    use: [Supplier, { name: Ingredient, as: SupplierIngredient }]
```

| Key | Value |
|---|---|
| `from` | the source: a scheme, then a colon, then a value that the scheme defines |
| `use` | a list of names; a name can be `{ name: …, as: … }` to rename it in this file |

The scheme tells where the names come from. The language has two schemes:

| Scheme | Value | Source |
|---|---|---|
| `file:` | a path relative to the importing file | a file in the same package |
| `pkg:` | the full name of a package | a different package |

The list of schemes can grow. A later version of the language can add a
scheme for a different source. The other import rules apply to all schemes.

- **FL9.** A file MUST import every name that it uses but does not define.
  This includes names from other files of the same package. Thus, the top of
  a file shows all that the file depends on.
- **FL10.** A `file:` import MUST point to a file in the same package, and
  that file MUST define each name that the import uses. `./` is implied. Use
  `../` only to go to the parent folder.
- **FL11.** A `pkg:` import MUST identify a package that is in the `dependsOn`
  of the importing package. It MUST use only names that the package
  [exports](package.md#exports).
- **FL17.** `from` MUST start with a scheme that the language defines. The
  tool rejects an unknown scheme.
- **FL12.** Two names in one file MUST NOT be the same; use `as` to rename one
  of them.
- **FL13.** A file MAY import the inverse name of a link alone, as for any
  other name. For example, `makesRecipe` is the inverse of `madeAtStation`, and a
  file can import `makesRecipe` without `madeAtStation`.
- **FL14.** Field names are not imported: a field belongs to its kind.
- **FL19.** Files of one package MAY import each other, also in a circle.

## Free-form text

Some keys take text that people write in their own words. This text follows
ASD-STE100 Simplified Technical English (STE): short sentences, one meaning
for each word, active voice and simple verb tenses.

| Text | Where |
|---|---|
| `definition` | kinds, links, value types |
| `description` | tags in the manifest |
| `reads` | links |
| `examples`, `notExamples` | kinds |
| `says` | constraints, derivatives, tests |
| `asks` | questions |
| `message` | constraints |
| Markdown | files in `docs/` |

- **FL16.** All free-form text SHOULD follow ASD-STE100. The tool MAY warn
  about the rules that it can check, such as sentence length and noun
  clusters.

## Ids

Files use plain names, for example `Recipe` for a kind and
`dish/margherita` for a fact. Imports show the source of each name.
An id is a name, not a web address. In the RDF export only, the tool puts
the `rdfNamespace` of the package in front of each name.
