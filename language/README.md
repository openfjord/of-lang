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

The language keeps the strengths of the semantic-web standards (RDF, OWL,
SHACL): exact meaning and, most of all, inference. It removes the parts that
make these standards difficult for software people:

- the notation
- the jargon
- the open world
- the absence of packages, tests and tools.

## How to read this

| If you want to | Read |
|---|---|
| find exactly how something works | [the reference](reference/README.md), one page for each part of the language |
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
2. **Check closed, infer by derivatives.** Constraints reject. Derivatives
   derive. No construct does both. No construct concludes the kind of a thing
   from a link that the thing has.
3. **Engineered like code.** The language uses the practices of code:
   packages, versions, namespaces, review, CI (continuous integration) and
   tests.
4. **Every change shows what it means.** A change shows what was stated, what
   follows, what no longer follows, what is rejected, and why.
5. **The standards are the export, not the source.** The tool generates RDF,
   OWL and SHACL for exchange and for existing tools. Nobody writes them.
6. **All text is Simplified Technical English.** This specification and all
   free-form text in a package follow ASD-STE100 Simplified Technical
   English (STE). See [FL16](reference/files.md#free-form-text).
7. **YAML is the default form, not the language.** The language is the
   statements, their keys and their values. YAML is the default way to write
   them, and this specification shows them in YAML. A different form, such as
   JSON, can hold the same statements. Packages, the `defs/` and `tests/`
   folders, files and imports are the same in every form. Only the syntax
   inside a file changes. The tool reads each form into the same model. A
   user interface is not a form. It shows and changes the statements in the
   files of a package, in whichever form those files are.

## Vocabulary

| Word | Meaning | Closest standard construct |
|---|---|---|
| **[package](reference/package.md)** | a versioned unit with one manifest and one owner; the owner owns everything in the package; its name can start with a namespace, as in npm | an ontology with `owl:imports` |
| **[statement](reference/files.md)** | one named item in a package; its `is:` key tells what it is | — |
| **[kind](reference/kind.md)** | a kind of thing (station, recipe, dish) | `owl:Class` |
| **[field](reference/field.md)** | a value of the thing itself (a name, a price); JSON Schema describes it | an attribute (`owl:DatatypeProperty`) |
| **[value type](reference/field.md#value-types)** | a named, reusable description of a field's value | `rdfs:Datatype` |
| **[link](reference/link.md)** | a named, directed relationship between two things | `owl:ObjectProperty` |
| **[tag](reference/tag.md)** | a named purpose that kinds, links and value types can carry; only the items that carry it are exported to RDF | — |
| **[constraint](reference/constraint.md)** | finds things by their links and conditions, and rejects the matches that break it | a SHACL shape |
| **[derivative](reference/derivative.md)** | finds things by their links and conditions, and adds links, kinds or field values for every match | a SHACL rule |
| **[fact](reference/fact.md)** | a stated thing with its fields and links; tests and examples use facts | RDF triples |
| **[derived fact](reference/fact.md#derived-facts)** | a fact that a derivative or a link declaration adds; tests can expect one, but no file declares one | RDF triples |
| **[test](reference/test.md)** | facts with what they must derive and reject, or a question with the answer it expects | a competency question |

Fields and links are different: a field is a value of the thing, and a link
points to another thing.
