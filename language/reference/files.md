---
title: Files and statements — The OpenFjord modeling language
---

# Files and statements

> Part of [the reference](README.md).

This page describes the YAML form of the language, which is the default
form ([principle 7](../README.md#principles)). In this form, the statements
are in YAML files with names that end in `.of.yaml`,
such as `recipe.of.yaml`. The extension tells editors and tools which files are
in the language. The only exception is the manifest of a package. Its name is
always `package.yaml`. A file is a mapping: an optional `import` block, and
then statements. The key of each statement is its name.

- **FL15.** Every file in `defs/` and `tests/` MUST have a name that ends in
  `.of.yaml`. Any other file there is an error. The tool does not ignore it.

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

| `is:` | Where | Described in |
|---|---|---|
| `Kind` | `defs/` | [Kinds](kind.md) |
| `Link` | `defs/` | [Links](link.md) |
| `ValueType` | `defs/` | [Fields and value types](field.md#value-types) |
| `Constraint` | `defs/` | [Constraints](constraint.md) |
| `Derivative` | `defs/` | [Derivatives](derivative.md) |
| `Tag` | `defs/` | [Tags](tag.md) |
| `Fact` | `tests/` | [Facts](fact.md) |
| `Test` | `tests/` | [Tests](test.md) |
| `Question` | `tests/` | [Tests](test.md#questions) |

- **FL1.** Every top-level key other than `import` MUST be a statement: a
  mapping with an `is:` key that has a value from the table above.
- **FL2.** A statement MUST NOT have a key that its page does not list. The
  tool rejects all keys that are not listed, so a misspelled key is an error.
  The tool does not ignore it.
- **FL3.** A YAML file MUST NOT repeat a key in the same mapping.
- **FL4.** Statements in `defs/` MUST be `Kind`, `Link`, `ValueType`,
  `Constraint`, `Derivative` or `Tag`. Statements in `tests/` MUST be `Fact`, `Test` or
  `Question`.
- **FL5.** No statement MAY be named `import`.

## Names

| Statement | Written as | Examples |
|---|---|---|
| kind, value type | UpperCamelCase | `Recipe`, `BaseRecipe`, `Minutes` |
| link, link inverse | lowerCamelCase, two words or more | `madeAtStation`, `usesRecipe`, `makesRecipe` |
| constraint, derivative, tag, test, question | lowerCamelCase | `dishListsItsAllergens`, `allergens` |
| fact | an id: any text, usually a path | `dish/margherita` |

- **FL6.** The name of a definition MUST be unique in its package. The
  inverse name that a link declares is also a name of that package.
- **FL7.** Kind and value type names MUST be UpperCamelCase; constraint,
  derivative, tag, test and question names MUST be lowerCamelCase.
  [LK1](link.md#rules) gives the rule for link names.
- **FL8.** The name of a statement in a test file MUST NOT be the name of a
  definition in the package. Statements in a test file are local to that file.

## Imports

A file can use a name only if the file defines or imports that name. The
`import` block is at the top of the file. It shows the source of each name:

```yaml
import:
  - from: file:recipe.of.yaml             # a file in this package, relative to this one
    use: [Recipe, madeAtStation]
  - from: file:../shared/station.of.yaml  # a file in the parent folder
    use: [Station]
  - from: pkg:@bistro/suppliers           # another package
    use: [Supplier, { name: Ingredient, as: SupplierIngredient }]
```

| Key | Value |
|---|---|
| `from` | `file:` and a path relative to the importing file, or `pkg:` and the full name of a package |
| `use` | a list of names; a name can be `{ name: …, as: … }` to rename it in this file |

- **FL9.** A file MUST import every name that it uses but does not define.
  This includes names from other files of the same package. Thus, the top of
  a file shows all that the file depends on.
- **FL10.** A `file:` import MUST point to a file in the same package, and
  that file MUST define each name that the import uses. `./` is implied. Use
  `../` only to go to the parent folder.
- **FL11.** A `pkg:` import MUST identify a package that is in the `dependsOn`
  of the importing package. It MUST use only names that the package
  [exports](package.md#exports).
- **FL12.** Two names in one file MUST NOT be the same; use `as` to rename one
  of them.
- **FL13.** A file MAY import the inverse name of a link alone, as for any
  other name. For example, `makesRecipe` is the inverse of `madeAtStation`, and a
  file can import `makesRecipe` without `madeAtStation`.
- **FL14.** Field names are not imported: a field belongs to its kind.

## Free-form text

Some keys take text that people write in their own words. This text follows
ASD-STE100 Simplified Technical English (STE): short sentences, one meaning
for each word, active voice and simple verb tenses.

| Text | Where |
|---|---|
| `definition` | kinds, links, value types, tags |
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
An id is a name, not a web address (reason 10 in
[the research](../../README.md)). The tool builds the full IRI
(Internationalized Resource Identifier) from the `rdfNamespace` of the
package. The full IRI is only in the exported RDF.
