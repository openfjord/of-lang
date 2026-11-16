---
title: The OpenFjord modeling language
---

# The OpenFjord modeling language

> **Status: a draft, explored, not decided.** This is the specification of
> the language that [the research](../README.md) describes. Nothing here is
> mandatory yet. For why the language is necessary, read
> [Why](../README.md#why) in the research.

The language describes the parts of a software product. The description is
a set of statements. By default, the statements are in YAML files, which are
in git next to the code. The statements describe:

- the kinds of things in the product
- how these things link to each other
- what follows from those links
- what is not allowed.

Developers read and write the statements. Product reads the sentences and
pages that a tool generates from the same statements. The tool checks the
statements, derives what follows, and rejects what breaks a rule.

## How to read this

| If you want to | Read |
|---|---|
| find exactly how something works | [the reference](#reference), one page for each part of the language |
| see all of the language on one page | [the cheat sheet](cheatsheet.md) |
| run it | [the `of-lang` tool](../cli/README.md) |

Every example uses one small, fictional restaurant. Stations in the kitchen
make recipes, dishes use recipes and ingredients, and ingredients contain
allergens. [examples/restaurant](../examples/restaurant/) holds the full
domain as one package.

## Conventions

The reference uses the words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD
NOT** and **MAY** as RFC (Request for Comments) 2119 defines them:

- MUST and MUST NOT are requirements that the tool checks.
- SHOULD and SHOULD NOT are strong advice.
- MAY is a permission.

Each requirement has a code, such as **KD3** (the third rule about kinds).
Thus, an error from the tool can identify the rule that it enforces.

## Principles

1. **One model, two readings.** Developers write and read objects. Product
   reads the sentences and examples that the tool generates from the same
   files.
2. **Constraints reject, derivatives derive.** No construct does both. Only
   a derivative can give a thing a narrower kind. A link never changes the
   kind of a thing.
3. **Engineered like code.** The language uses the practices of code:
   packages, versions, namespaces, review, CI (continuous integration) and
   tests.
4. **Every change shows what it means.** A change shows what was stated, what
   follows, what no longer follows, what is rejected, and why.
5. **RDF is an export, not the source.** The tool generates RDF (Resource
   Description Framework), OWL (Web Ontology Language) and SHACL (Shapes
   Constraint Language) for exchange and for existing tools. Nobody writes
   them. The export does not translate each statement into one RDF item. It
   holds the same knowledge as the package and its dependencies, so it
   answers the same questions.
6. **All text is Simplified Technical English.** This specification and all
   free-form text in a package follow ASD-STE100 Simplified Technical
   English (STE). See [FL16](reference/files.md#free-form-text).
7. **YAML is the default form, not the language.** The language is the
   statements, their keys and their values. YAML is the default way to write
   them, and this specification shows them in YAML. A different form, such as
   JSON, can hold the same statements. Packages, the `defs/` and `tests/`
   folders, files and imports are the same in every form. Only the syntax
   inside a file and the ending of its name change, and a package uses one
   form. The manifest is always YAML. The tool reads each form into the
   same model. A
   user interface is not a form. It edits packages: it changes the files of a
   package, in whichever form those files are.
8. **No prior knowledge.** A reader needs no knowledge of RDF or of
   ontologies to understand the language. This specification explains each
   idea in its own words, and does not compare the language with other
   standards.

## Vocabulary

| Word | Meaning |
|---|---|
| **[package](reference/package.md)** | a versioned unit with one manifest and one owner; the owner owns everything in the package; its name can start with a namespace, as in npm |
| **[statement](reference/files.md)** | one named item in a package; its `is:` key tells what it is |
| **[kind](reference/kind.md)** | a kind of thing (station, recipe, dish); a kind can also be exactly one of a fixed list of kinds, with `oneOf` |
| **[field](reference/kind.md#fields)** | a value of the thing itself (a name, a price); JSON Schema describes it |
| **[value type](reference/kind.md#value-types)** | a named, reusable description of a field's value |
| **[link](reference/link.md)** | a named, directed relationship between two things |
| **[tag](reference/tag.md)** | a label for a purpose that kinds, links and value types can carry; the package manifest lists each tag with a name and a description; by default the RDF export includes everything, and an optional list of tags limits it to the items that carry them |
| **[constraint](reference/constraint.md)** | finds things by their links and conditions, and rejects the matches that break it |
| **[derivative](reference/derivative.md)** | finds things by their links and conditions, and adds links, kinds or field values for every match |
| **[fact](reference/fact.md)** | a stated thing with its fields and links; tests and examples use facts |
| **[derived fact](reference/fact.md#derived-facts)** | a fact that nobody stated: a derivative adds it, or the `under` or `chains` key of a link implies it; tests can expect one, but no file declares one |
| **[test](reference/test.md)** | facts with what they must derive and reject, or a question with the answer it expects |
| **[pending change](reference/pending.md)** | the statements that a change adds, removes or updates, marked in the files until somebody accepts or rejects the change; one for each package at most |

Fields and links are different: a field is a value of the thing, and a link
points to another thing.

## Reference

Each page describes one part of the language: its purpose, every key that it
takes, and the requirements that the tool checks.

| Page | Covers | Codes |
|---|---|---|
| [Files and statements](reference/files.md) | files, statements, `is:`, names, imports, ids | FL |
| [Packages](reference/package.md) | the manifest, owners, namespaces, the fixed layout, exports, docs | PK, DC |
| [Kinds, fields and value types](reference/kind.md) | kinds, `extends`, unions, fields as a JSON Schema subset, value types | KD, FD, VT |
| [Links](reference/link.md) | links, `under`, `inverse`, `chains` | LK |
| [Tags](reference/tag.md) | tags, and the export of a subset to RDF | TG |
| [Constraints](reference/constraint.md) | `match`, `filter`, `unless`, messages, how they run | CN |
| [Derivatives](reference/derivative.md) | `match`, `filter`, `derive`, how they run, explanations | DV |
| [Facts](reference/fact.md) | facts in tests, derived facts | FC |
| [Tests](reference/test.md) | examples, and questions with their own `match` and `filter` | TS |
| [Pending changes](reference/pending.md) | `pending` and `was`: a change that nobody has accepted yet | PD |
