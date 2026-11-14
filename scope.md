---
title: Scope and open questions — A modeling language for software people
---

# Scope and open questions

> Part of [a modeling language for software people](README.md), a first
> draft. Reason numbers refer to the table there.

## What it leaves out, and why

- **Open-world reasoning.** No "either this or that", no "something exists but
  is unnamed". It cannot be explained in a change report.
- **`owl:sameAs` and anything that merges things.** One mistaken statement
  would silently join two records.
- **Domain and range as inference.** They become checks.
- **OWL restrictions as constraints.** Counts and allowed values are SHACL.
- **Facts as definitions.** Facts appear only in tests and examples, and
  define nothing.

## Roads not taken

- **Write RDF, OWL and SHACL directly.** Exact and standard, and every reason
  in [the table](README.md#why-ontologies-are-hard-for-software-people)
  applies in full.
- **Prose pages with untyped links, as in OKF.** Easy to read and write;
  nothing can be checked or concluded.
- **A platform-owned model, as in Palantir Foundry.** Typed and governed, but
  managed outside the codebase and without inference.
- **An expression language for constraints and rules,** such as CEL (Common
  Expression Language) or a string syntax of our own. Shorter, but a reader
  must parse operators and paths in their head, and dotted paths hide the
  things they pass through.
- **Datalog as written.** The semantics wanted, but variables threaded through
  terse lines few developers read comfortably.
- **MongoDB-style operators, as in Mingo.** YAML-native, but one document at a
  time, with no derivation or explanations.
- **A graph query language,** such as GQL (Graph Query Language) or Cypher.
  Readable patterns, but as strings, and for querying rather than deriving.
- **A separate specification for the logic.** It would let other languages
  reuse it, at the cost of a second schema interface to keep in step; the
  logic is written against this language's own kinds and links instead.
- **Fields on links, as in a property graph.** It lets a link say when it
  started or how much it carries. But derived and chained links would have
  nothing to fill the fields with, and patterns would need a way to reach
  them. A link stays only its link, from and to; anything said about it
  becomes a thing in the middle, such as a `Read` between a job and a table.

## Deferred: facts at runtime

Facts written by running systems are set aside for now. Two answers in the
reasons table wait on them:

- **Trust and freshness** (reason 17): source, confirmation and expiry per
  fact, trust levels, and whether a single link carries its own source.
- **Runtime writes and types** (reason 15): how facts written at runtime are
  kept from changing kinds, links, constraints or rules.

## Open questions

1. **Can a technical package add rules about product kinds?** `Table` lives
   in a product package. The platform team wants a rule such as "a table with
   more than 1,000 columns is oversized". Everything in a package is owned by
   its owner, so the rule cannot go in the product package without product
   approving it. May the platform team put that rule in their own package,
   through an extension point, or must every rule about `Table` go through
   product?
2. **Optional links in a pattern.** Every link in a pattern must exist, or
   there is no match. A rule that wants "every table, and its team if it has
   one" loses the tables without a team. Should a link in a pattern be markable
   as optional, or should authors work around it with two rules or `any of`?
3. **Limits on chains.** `repeated: true` follows a link any number of steps:
   table to table to table. Sometimes only nearby things matter, such as
   "tables at most two steps upstream". Should chains take a limit
   (`at most: 2`), or stay unlimited?
4. **Dates and durations.** Conditions compare a value with another value.
   They cannot say "a job that has not run for 30 days", because that needs
   "now minus 30 days". Should conditions gain date arithmetic, and if so, in
   what YAML form?
