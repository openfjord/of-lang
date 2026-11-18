---
title: Design of the language — The OpenFjord modeling language
---

# Design of the language

Why the language looks the way it does: the reasons ontologies are hard for
software people, what the language leaves out, the roads not taken, and what
is still open. [The README](../README.md) gives the short version; [the
language specification](../language/README.md) gives the rules.

## The question

Can a language keep what the W3C standards are good at, inference above all,
and remove what makes them hard for software people?

## Why ontologies are hard for software people

| # | Reason | What would help |
|---|---|---|
| 1 | **The notation is odd and repetitive.** Triples instead of objects; links defined apart from the classes that use them, so nothing reads like `Foo { … }`. | An object-style way to write, generated into triples, that still shows what each link implies. |
| 2 | **Data without knowledge.** Vague, generic links reused for unrelated meanings (`uses`, `relatedTo`) add lines but nothing can be concluded from them. | Specific links, grouped under general ones; a one-sentence definition each; names that read as verb phrases with a direction. |
| 3 | **The jargon and perspective are unintuitive.** Familiar words mean other things: a class is a set, a property is not a field, domain and range draw conclusions instead of checking, a missing fact is unknown rather than false. | Teach the false friends first; checks that reject for "is this allowed", inference only for "what follows". |
| 4 | **The learning curve is steep** and most developers will not climb it: the payoff is late, hard to see, and rarely asked for. | A few people model, everyone else uses the result; a subset that fits on one page; a payoff on the first day. |
| 5 | **An ontology cannot hold a whole product.** It is weak at why, at workflows, and at context. | Checkable facts in the model; reasoning and sequences in Markdown, in a tree of its own. Docs link to or embed definitions by hand; the reverse links are generated. |
| 6 | **No tools fit a software factory.** Existing editors are built for logicians or draw unreadable whole-graph pictures. | Text files in git, checks in CI (continuous integration), editor support, one page per term, impact views, generated code and docs. |
| 7 | **Product is the hardest reader.** Product owns many meanings but cannot read the format. | Generated sentences, examples and non-examples, ownership that follows from the package a term lives in, review on a rendered page. |
| 8 | **The filesystem has one dimension**; the graph has many groupings. | Packages follow who owns them; every other grouping is a property of the term; views are generated. |
| 9 | **No packages or dependencies.** Imports have no versions, no exports and no rules about direction. | A manifest per package, public and internal terms, allowed directions checked in CI, clear breaking-change rules. |
| 10 | **IRIs (Internationalized Resource Identifiers) look like links.** Developers click them, try to host them, or treat them as locations. | Show short names; an id is a name, not a link; if ids resolve, they open the term's page. |
| 11 | **The graph is flat.** Anyone can say anything about any term, and merged facts forget where they came from. | Every term qualified by its package; a package only makes statements about its own terms; explicit extension points; each package's statements kept in their own graph. |
| 12 | **Adding data has invisible effects.** Hidden conclusions, and nothing ever fails. | A report of what a change concludes, before it lands; checks that reject; a small set of inference rules; guided entry. |
| 13 | **Mechanics, lint and business rules look alike.** An OWL "exactly one owner" merges two owners into one instead of rejecting. | One construct per kind of rule, each with its own owner and its own behavior on failure. |
| 14 | **The type system and the knowledge are mixed** (the TBox, terminology, and the ABox, assertions). | Separate packages for types and for facts, reviewed differently. |
| 15 | **Schema changes and data changes are the same language** — there is no DDL/DML (data definition / data manipulation language) split. | Facts can never define anything. How facts written at runtime are kept from changing the types is deferred, with runtime facts. |
| 16 | **No separation of concerns,** and no tests. | Each concern a separate statement — kind, link, value type, constraint, derivative — placed next to what it is about; competency questions and example data as tests. |
| 17 | **Trust and freshness are not modeled.** Every fact looks equally sure and equally current. | Source, confirmation and expiry per fact; trust levels; views filtered by trust. |
| 18 | **Inference must be visible.** Its value is the knowledge nobody wrote down, which only counts if people see it when they add something. | Every change reports what is newly derived and what is no longer derived, each derived fact with its explanation. |

## What it leaves out, and why

- **Open-world reasoning.** No "either this or that", no "something exists but
  is unnamed". It cannot be explained in a change report.
- **`owl:sameAs` and anything that merges things.** One mistaken statement
  would silently join two records.
- **Domain and range as inference.** They become checks.
- **OWL restrictions as constraints.** Constraints are not exported as OWL
  restrictions, which infer. Counts and allowed values become SHACL shapes,
  which reject.
- **Arithmetic on dates and the current time.** A rule about "now" gives a
  different answer tomorrow with the same facts. Time belongs to the system
  that is modeled, not to the definitions. Dates are plain field values.
- **Optional links in a match.** Nothing in the language returns rows with
  empty cells. "Every dish, and its garnish if it has one" is two
  statements, or a condition with `at least one` or `none`.
- **A limit on chains.** The links that chain, such as "uses" and "contains
  an allergen", are true at any distance. A rule that needs exactly two
  steps writes two link items.
- **Facts as definitions.** Facts appear only in tests and examples, and
  define nothing.

## Roads not taken

- **Write RDF, OWL and SHACL directly.** Exact and standard, and every reason
  in [the table](#why-ontologies-are-hard-for-software-people)
  applies in full.
- **Prose pages with untyped links, as in OKF.** Easy to read and write;
  nothing can be checked or concluded.
- **A platform-owned model, as in Palantir Foundry.** Typed and governed, but
  managed outside the codebase and without inference.
- **An expression language for constraints and derivatives,** such as CEL (Common
  Expression Language) or a string syntax of our own. Shorter, but a reader
  must parse operators and paths in their head, and dotted paths hide the
  things they pass through.
- **Datalog (a logic language of facts and rules) as written.** The semantics wanted, but variables threaded through
  terse lines few developers read comfortably.
- **MongoDB-style operators, as in Mingo (a JavaScript library that runs
  MongoDB queries).** YAML-native, but one document at a
  time, with no derivation or explanations.
- **A graph query language,** such as GQL (Graph Query Language) or Cypher.
  Readable patterns, but as strings, and for querying rather than deriving.
- **A separate specification for the logic.** It would let other languages
  reuse it, at the cost of a second schema interface to keep in step; the
  logic is written against this language's own kinds and links instead.
- **Fields on links, as in a property graph.** It lets a link say when it
  started, or how much of an ingredient a recipe uses. But derived and chained links would have
  nothing to fill the fields with, and matches would need a way to reach
  them. A link stays only its link, from and to; anything said about it
  becomes a thing in the middle, such as a `Portion` between a recipe and an ingredient.

## Deferred: mappings to outside vocabularies

Saying that a definition means the same as a term in another vocabulary —
`Dish` as a class in someone else's ontology — is left out of the
specification until it is designed.

## Deferred: facts at runtime

Facts written by running systems are set aside for now. Two answers in the
reasons table wait on them:

- **Trust and freshness** (reason 17): source, confirmation and expiry per
  fact, trust levels, and whether a single link carries its own source.
- **Runtime writes and types** (reason 15): how facts written at runtime are
  kept from changing kinds, links, constraints or derivatives.

## Open questions

1. **Checking a tagged export.** An export must hold the same knowledge as
   its package and answer the same questions. A tagged export must hold only
   the tagged knowledge: it must not answer a question that needs knowledge
   without the tag. Tests and questions carry no tags. Which of them belong
   to a subset, and how does the tool check that a tagged export passes
   those, and no others? An untagged derivative that adds facts to a tagged
   link is part of this question: the subset then misses those facts.

## What the tool may need to do

The tool may need to:

- **Validate packages.** Check each package against every rule in the
  specification.
- **Compile packages to RDF.** Make RDF, OWL and SHACL from a package and its
  dependencies, with an optional list of tags that selects a subset.
- **Support editors through the LSP** (Language Server Protocol): errors,
  completion and links to definitions while a person writes.
- **Give a user interface** that edits packages. It is a web app.
  `of-lang ui` starts a server that serves it. The interface always changes
  the files of a package, in whichever form those files are. The rest of its design is not yet discussed.
