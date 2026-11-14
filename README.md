---
title: A modeling language for software people — Research
---

# A modeling language for software people — Research

> **Status: a first draft, explored, not decided.** It lists why ontologies are
> hard for software product and development people to read and write, and
> sketches a language that answers each reason. Nothing here is binding. The
> examples use a small made-up domain — teams, databases, tables, applications
> and the jobs that read and write tables. Written 2026-11-13.
>
> A command-line tool, `of-lang`, implements the draft: [cli/](cli/), with an
> example package in [examples/data-model/](examples/data-model/).

## The question

The W3C semantic-web standards — RDF (Resource Description Framework), RDFS
(RDF Schema), OWL (Web Ontology Language) and SHACL (Shapes Constraint
Language) — can describe a software product's domain exactly. Few product or
development people will read or write them. Can a language keep what the
standards are good at, inference above all, and drop what makes them hard?

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
| 16 | **No separation of concerns,** and no tests. | Each concern a separate entry — kind, link, constraint, rule — placed next to what it is about; competency questions and example data as tests. |
| 17 | **Trust and freshness are not modeled.** Every fact looks equally sure and equally current. | Source, confirmation and expiry per fact; trust levels; views filtered by trust. |
| 18 | **Inference must be visible.** Its value is the knowledge nobody wrote down, which only counts if people see it when they add something. | Every change reports what is newly derived and what is no longer derived, each derived fact with its explanation. |

### What others do

- **OKF (Open Knowledge Format, Google Cloud, version 0.2)** stores knowledge
  as one Markdown page per concept with YAML frontmatter. It fits reasons 5, 7
  and 17 well: prose first, readable by product, built for git, with fields for
  sources, confirmation and expiry. Its links carry no type, its types are free
  text, and readers must never reject anything, so it gives up reasons 2, 13
  and 14 entirely.
- **Palantir Foundry** models object types, link types and action types over
  tables. Every write goes through a governed action — the clearest answer to
  reason 15 found. It does not infer, and its model is managed in its own
  interface rather than as code.
- **Databricks Unity Catalog** is a catalog with a three-level name, lineage
  and metric definitions. It answers namespaces and provenance, and has no
  typed links.

Both platforms chose closed-world checking and no inference. This draft keeps
closed-world checking and adds inference by rules.

## The draft, file by file

- [The language](language.md): principles, vocabulary, entries and imports, ids.
- Terms, one file each:
  - [Packages](terms/package.md): ownership, manifests, exports, docs.
  - [Kinds](terms/kind.md): fields, link counts, `extends`.
  - [Fields and value types](terms/field.md): the JSON Schema subset.
  - [Links](terms/link.md): `under`, `inverse`, `chains`.
  - [Patterns](terms/pattern.md): what constraints, rules and questions share.
  - [Constraints](terms/constraint.md): `require` and messages.
  - [Rules](terms/rule.md): `derive`.
  - [Facts](terms/fact.md): stated facts in tests, and derived facts.
  - [Tests](terms/test.md): examples and questions.
- [What it outputs](outputs.md): the change report, product's pages, the RDF
  export.
- [Scope and open questions](scope.md): what it leaves out, roads not taken,
  what is deferred, and what is still open.
