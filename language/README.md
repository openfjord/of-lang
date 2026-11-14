---
title: The OpenFjord modeling language
---

# The OpenFjord modeling language

> **Status: a draft, explored, not decided.** This is the specification of
> the language described in [the research](../README.md). Nothing here is
> binding yet.

The language describes what a software product is made of — the kinds of
things in it, how they link to each other, what follows from those links, and
what is not allowed — in YAML files that live in git beside the code.
Developers read and write the files. Product reads sentences and pages
generated from the same files. A tool checks the files, derives what follows,
and rejects what breaks a rule.

It keeps what the semantic-web standards (RDF, OWL, SHACL) are good at —
exact meaning, and inference above all — and drops what makes them hard for
software people: the notation, the jargon, the open world, and the lack of
packages, tests and tooling.

## How to read this

| If you want to | Read |
|---|---|
| learn the language from nothing | [the tour](tour/README.md), seven short steps that build one package |
| look up exactly how something works | [the reference](reference/README.md), one page per part of the language |
| see the whole language at a glance | [the cheat sheet](cheatsheet.md) |
| run it | [the `of-lang` tool](../cli/README.md) |

Every example uses one small made-up domain: teams own databases, databases
hold tables, applications run jobs, and jobs read and write tables. The tour
builds that domain step by step; the finished package is in
[examples/data-model](../examples/data-model/).

## Conventions

The words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT** and **MAY** in
the reference are used as RFC 2119 defines them: MUST and MUST NOT are
requirements the tool checks, SHOULD and SHOULD NOT are strong advice, and MAY
is a permission.

Each requirement has a code, such as **KD3** (the third rule about kinds), so
an error from the tool can name the rule it enforces.

## Principles

1. **One model, two readings.** Developers write and read objects; product
   reads sentences and examples generated from the same files.
2. **Check closed, infer by rules.** Constraints reject. Rules derive. Nothing
   does both, and no construct ever concludes a thing's kind from a link it
   has.
3. **Engineered like code.** Packages, versions, namespaces, review, CI
   (continuous integration) and tests.
4. **Every change says what it means.** What was stated, what follows, what no
   longer follows, what is rejected, and why.
5. **The standards are the export, not the source.** RDF, OWL and SHACL are
   generated for exchange and for existing tools; nobody writes them.

## Vocabulary

| Word | Meaning | Closest standard construct |
|---|---|---|
| **[package](reference/package.md)** | a versioned unit with one manifest and one owner, who owns everything in it; its name may start with a namespace, as in npm | an ontology with `owl:imports` |
| **[entry](reference/files.md)** | one named item in a file, saying what it is with `is:` | — |
| **[kind](reference/kind.md)** | a kind of thing (database, table, application) | `owl:Class` |
| **[field](reference/field.md)** | a value of the thing itself (a name, a format), described with JSON Schema | an attribute (`owl:DatatypeProperty`) |
| **[value type](reference/field.md#value-types)** | a named, reusable description of a field's value | `rdfs:Datatype` |
| **[link](reference/link.md)** | a named, directed relationship between two things | `owl:ObjectProperty` |
| **[pattern](reference/pattern.md)** | the links between things, each with its kind, and the conditions a constraint, rule or question looks for | a SPARQL graph pattern |
| **[constraint](reference/constraint.md)** | a pattern that rejects facts | a SHACL shape |
| **[rule](reference/rule.md)** | a pattern that derives new links, kinds or field values | a SHACL rule |
| **[fact](reference/fact.md)** | a stated thing with its fields and links; used in tests and examples | RDF triples |
| **[derived fact](reference/fact.md#derived-facts)** | a fact a rule or a link's declaration adds | RDF triples in a separate graph |
| **[test](reference/test.md)** | facts with what they must derive and reject, or a question with the answer it expects | a competency question |

Fields and links stay apart: a field is a value of the thing, a link points
at another thing.
