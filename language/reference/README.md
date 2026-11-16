---
title: Reference — The OpenFjord modeling language
---

# Reference

> Part of [the OpenFjord modeling language](../README.md). This reference
> uses MUST, MUST NOT, SHOULD and MAY as RFC (Request for Comments) 2119 defines them. Each
> requirement has a code, and an error from the tool can show that code.

Each page describes one part of the language: its purpose, every key that it
takes, and the requirements that the tool checks.

| Page | Covers | Codes |
|---|---|---|
| [Files and statements](files.md) | files, statements, `is:`, names, imports, ids | FL |
| [Packages](package.md) | the manifest, owners, namespaces, the fixed layout, exports, docs | PK, DC |
| [Kinds, fields and value types](kind.md) | kinds, `extends`, link counts, fields as a JSON Schema subset, value types | KD, FD, VT |
| [Links](link.md) | links, `under`, `inverse`, `chains` | LK |
| [Tags](tag.md) | tags, and the export of a subset to RDF | TG |
| [Constraints](constraint.md) | `match`, `filter`, `unless`, messages, how they run | CN |
| [Derivatives](derivative.md) | `match`, `filter`, `derive`, how they run, explanations | DV |
| [Facts](fact.md) | facts in tests, derived facts | FC |
| [Tests](test.md) | examples, and questions with their own `match` and `filter` | TS |
| [Pending changes](pending.md) | `pending` and `was`: a change that nobody has accepted yet | PD |
