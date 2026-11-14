---
title: Reference — The OpenFjord modeling language
---

# Reference

> Part of [the OpenFjord modeling language](../README.md). MUST, MUST NOT,
> SHOULD and MAY are used as RFC 2119 defines them; each requirement has a
> code the tool's errors can name.

Each page describes one part of the language: what it is for, every key it
takes, and the requirements the tool checks.

| Page | Covers | Codes |
|---|---|---|
| [Files and entries](files.md) | YAML files, entries, `is:`, names, imports, ids | FL |
| [Packages](package.md) | the manifest, owners, namespaces, the fixed layout, exports, docs | PK, DC |
| [Kinds](kind.md) | kinds, `extends`, link counts | KD |
| [Fields and value types](field.md) | the JSON Schema subset, value types | FD, VT |
| [Links](link.md) | links, `under`, `inverse`, `chains` | LK |
| [Patterns](pattern.md) | what constraints, rules and questions share: `match`, `filter`, variables, conditions, how they run | PT |
| [Constraints](constraint.md) | `require`, messages | CN |
| [Rules](rule.md) | `derive` | RL |
| [Facts](fact.md) | facts in tests, derived facts | FC |
| [Tests](test.md) | examples and questions | TS |
