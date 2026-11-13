---
title: A modeling language for software people — Research
---

# A modeling language for software people — Research

> **Status: a first draft, explored, not decided.** It lists why ontologies are
> hard for software product and development people to read and write, and
> sketches a language that answers each reason. Nothing here is binding. The
> examples use a small made-up domain — teams, databases, tables, applications
> and the jobs that read and write tables. Written 2026-11-13.

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
| 5 | **An ontology cannot hold a whole product.** It is weak at why, at workflows, and at context. | Checkable facts in the model; reasoning and sequences in Markdown. Docs link to terms by hand; the reverse links are generated. |
| 6 | **No tools fit a software factory.** Existing editors are built for logicians or draw unreadable whole-graph pictures. | Text files in git, checks in CI (continuous integration), editor support, one page per term, impact views, generated code and docs. |
| 7 | **Product is the hardest reader.** Product owns many meanings but cannot read the format. | Generated sentences, examples and non-examples, ownership that follows from where a term lives, review on a rendered page. |
| 8 | **The filesystem has one dimension**; the graph has many groupings. | Folders follow who reviews a change; every other grouping is a property of the term; views are generated. |
| 9 | **No packages or dependencies.** Imports have no versions, no exports and no rules about direction. | A manifest per package, public and internal terms, allowed directions checked in CI, clear breaking-change rules. |
| 10 | **IRIs (Internationalized Resource Identifiers) look like links.** Developers click them, try to host them, or treat them as locations. | Show short names; an id is a name, not a link; if ids resolve, they open the term's page. |
| 11 | **The graph is flat.** Anyone can say anything about any term, and merged facts forget where they came from. | A namespace per package; a package only makes statements about its own terms; explicit extension points; each package's statements kept in their own graph. |
| 12 | **Adding data has invisible effects.** Hidden conclusions, and nothing ever fails. | A report of what a change concludes, before it lands; checks that reject; a small set of inference rules; guided entry. |
| 13 | **Mechanics, lint and business rules look alike.** An OWL "exactly one owner" merges two owners into one instead of rejecting. | One construct per kind of rule, each with its own owner and its own behavior on failure. |
| 14 | **The type system and the knowledge are mixed** (the TBox, terminology, and the ABox, assertions). | Separate packages for types and for facts, reviewed differently. |
| 15 | **Schema changes and data changes are the same language** — there is no DDL/DML (data definition / data manipulation language) split. | Types arrive only by release; runtime writes go through a gateway that accepts facts only. |
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

## Principles

1. **One model, two readings.** Developers write and read objects; product
   reads sentences and examples generated from the same files.
2. **Check closed, infer by rules.** Constraints reject. Rules derive. Nothing
   does both, and no construct ever concludes a thing's kind from a link it
   has.
3. **Engineered like code.** Packages, versions, namespaces, review, CI and
   tests.
4. **Every change says what it means.** What was stated, what follows, what no
   longer follows, what is rejected, and why.
5. **Every fact says how far to trust it.** Where it came from, who confirmed
   it, when to check it again.
6. **The standards are the export, not the source.** RDF, OWL and SHACL are
   generated for exchange and for existing tools; nobody writes them.

## The language

### Vocabulary

| Word | Meaning | Closest standard construct |
|---|---|---|
| **package** | a versioned unit with one namespace and one manifest; who owns it follows from where it sits | an ontology with `owl:imports` |
| **kind** | a kind of thing (database, table, application) | `owl:Class` |
| **field** | a value of the thing itself (a name, a format) | an attribute (`owl:DatatypeProperty`) |
| **link** | a named, directed relationship between two things | `owl:ObjectProperty` |
| **constraint** | a rule that rejects data | a SHACL shape |
| **rule** | a pattern that derives new links, kinds or field values | a SHACL rule |
| **fact** | a stated thing with its fields and links | RDF triples |
| **derived fact** | a link produced by a rule | RDF triples in a separate graph |

Fields and links stay apart: a field is a value of the thing, a link points
at another thing.

### Packages

A package is a folder. It holds either **types** (kinds, links, constraints,
rules) or **facts**, never both (reasons 14 and 15).

**Ownership is decided by location.** Packages sit in folders, and a folder
says who must approve a change to everything in it — product or not. Nothing
in the language names an owner. A package is owned as a whole, so a single
file never holds both product-owned and other things: a product term and a
technical term live in different packages (reasons 7 and 8).

Only `package.yaml`, at the package's root, is required. Every other YAML
file in the package may hold kinds, links, constraints, rules and mappings in
any mix, in any folders the author chooses. One way to lay out a package:

```
data-model/                        # a types package
  package.yaml                     # name, version, dependencies, exports
  definitions/
    job.yaml                       # Job, and the links that come out of it
    table.yaml                     # Table, its links, the rules about tables
    database.yaml
    application.yaml               # Application, its links, the rules about applications
    reading-across-applications.yaml
                                   # a constraint spanning jobs, tables and applications
  docs/                            # Markdown: why and how, linking to terms
  tests/
    questions/                     # questions with expected answers
    examples/                      # data that must pass, data that must fail

data-facts/                        # a facts package
  package.yaml
  databases.yaml
  tables.yaml
  applications.yaml
```

**Place things next to what they are about.** A link sits with the kind it
comes out of. A constraint, rule or mapping sits with the kind it concerns;
one that concerns several sits a level up, where it can be seen beside all of
them. This is guidance for the author, not a check: the language enforces no
folder structure.

**Every entry says what it is.** Each top-level entry in a file is keyed by
its name and states its type with `is:` — `Kind`, `Link`, `Constraint`,
`Rule` or `Mapping` for definitions, and one of the package's kinds for a
fact. Names are unique within a package.

```yaml
Table:
  is: Kind
  …

inDatabase:
  is: Link
  …

db/sales/table/order_lines:
  is: Table
  …
```

**Every file imports what it uses.** A file may use a name only if it defines
it or imports it, at the top, in an `import` block:

```yaml
import:
  - from: file:table.yaml            # a file in this package, relative to this one
    use: [Table, inDatabase]
  - from: file:../shared/team.yaml   # going up a folder
    use: [Team]
  - from: pkg:org-model              # another package
    use: [Contract, { name: Team, as: OrgTeam }]
```

- `file:` names a file in the same package, relative to the importing file;
  `../` goes up a folder.
- `pkg:` names another package, which must be in `dependsOn`, and only names
  it exports may be imported.
- Imports from the same package are required too, so the top of a file shows
  everything it depends on, and every name leads to its definition.
- `as` renames an import when two names clash.
- `import` is reserved: no entry may be named `import`.

The manifest:

```yaml
package: data-model
version: 0.3.0
namespace: data                   # used only in exported RDF, never in the files
holds: types                      # types | facts
dependsOn:
  org-model: ^1.2                 # semantic-version range, pinned in a lock file
exports:                          # what other packages may use; the rest is internal
  kinds: [Database, Table, Application, Job]
  links: [inDatabase, hasTable, readsTable, writesTable, jobOf]
openFor:                          # extension points: what other packages may add to
  - extend: Table                 #   another package may declare a kind of table
  - under: readsTable             #   another package may declare a narrower link under readsTable
```

Checks on packages, all in CI:

- A package uses only terms it defines or that a dependency exports.
- A package states things only about its own terms, except through the
  extension points its dependencies open (reason 11).
- A facts package defines no kinds, links, constraints or rules.
- Dependencies follow allowed directions; a cycle is an error.
- Version changes match what changed: removing or renaming a term, or
  tightening a constraint, is a breaking change; adding a term or loosening a
  constraint is not.

Ids are names, not links (reason 10). Files use plain names — `Table` in
types, `db/sales/table/order_lines` for facts — and imports say where each
comes from. The full IRI, built from the package's namespace, appears only in
the exported RDF.

### Kinds

A kind reads like an object. It says what it is, gives examples, lists its
fields, and lists the links it carries with how many of each.

```yaml
Table:
  is: Kind
  definition: A set of rows with one schema, stored in a database.
  examples:
    - the order_lines table in the sales database
  notExamples:
    - a temporary copy a test run makes and throws away
    - a view computed on read
  fields:
    name:       { type: name }
    format:     { oneOf: [row, columnar] }
    visibility: { oneOf: [public, internal], default: internal }
  links:
    inDatabase: exactly 1
    hasColumn:  at least 1
  docs: docs/tables.md

Database:
  is: Kind
  definition: A named store of tables, owned by one team.
  fields:
    name: { type: name }
  links:
    ownedByTeam: exactly 1

Job:
  is: Kind
  definition: A unit of work an application runs, reading and writing tables.
  fields:
    name: { type: name }
  links:
    jobOf: exactly 1
```

A sub-kind names its parent: `PartitionedTable: { is: Kind, extends: Table, … }`.

Every field and link a kind lists is required unless a count or a default
says otherwise, and anything not listed is rejected. The counts are
constraints, not inference: a database naming two teams is rejected, never
read as "the two teams are the same".

### Links

A link is declared once, as its own entry, and the kind it comes out of says
how many it carries. The link usually sits in the same file as that kind, so
an object's fields, the links that come out of it and the limits on them read
together.

```yaml
inDatabase:
  is: Link
  from: Table
  to: Database
  reads: "{from} is in {to}"
  definition: The database that stores the table.
  inverse: hasTable

jobOf:
  is: Link
  from: Job
  to: Application
  reads: "{from} is a job of {to}"
  definition: The application that runs the job.
  inverse: hasJob

readsTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} reads {to}"
  definition: The job reads rows from the table when it runs.
  under: usesTable                   # every readsTable is also a usesTable
  inverse: readByJob

writesTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} writes {to}"
  definition: The job writes rows to the table when it runs.
  under: usesTable
  inverse: writtenByJob

feedsTable:
  is: Link
  from: Table
  to: Table
  reads: "{from} feeds {to}"
  definition: >
    Rows of the first table flow, directly or through other tables, into the
    second.
  chains: true                       # a feeds b and b feeds c means a feeds c
```

Rules on links:

- **The name says what it links and how** — two words at least, never a bare
  verb: `readsTable`, never `reads` (reason 2).
- **A definition is required**, and so is `reads`, the sentence template that
  product sees (reason 7).
- **`from` and `to` are checks.** A team that `readsTable` is rejected. No
  link ever makes something into a kind.
- **`under`** places a link beneath a broader one, so specific links can be
  written and general questions still answered.
- **`inverse`** and **`chains`** are the only built-in inference on links.

### Constraints and rules

Counts and allowed values sit on kinds. Anything that looks across more than
one thing is a **constraint** (it rejects) or a **rule** (it derives). Both
are written the same way: match a
pattern of links, narrow it with conditions, then require something of every
match or derive new facts from it.

Four principles shape how they are written:

1. **Readability comes before brevity.** Every part is spelled out in YAML
   keys and words. There are no expression strings to parse.
2. **No implicit variables.** Every thing a statement talks about has a name.
   Nothing is reached through a dotted path that skips over an unnamed thing.
3. **Variables start with `$`,** so a variable is never mistaken for a kind, a
   link, a field or a value.
4. **The sentence declares the variables.** Every variable appears in the
   statement's `says` sentence, and the logic may use only those. The sentence
   product reads and the logic that runs cannot drift apart.

#### The shape of a statement

| Key | Holds | In |
|---|---|---|
| `says` | the sentence, with every variable in it | both |
| `links` | the pattern to match | both (may be left out when there is one variable) |
| `where` | conditions that narrow the matches | both, optional |
| `require` | conditions every match must satisfy | constraints |
| `message` | what a failure says | constraints |
| `derive` | facts to add for every match | rules |

```yaml
crossAppReadNeedsContract:
  is: Constraint
  says: >
    A job $job of application $jobApp may read a table $table that job
    $writer of another application $writerApp writes, only if $table is
    public and $jobApp consumes a contract $contract that covers it.

  links:
    - from: $job
      link: jobOf
      to:   $jobApp
    - from: $job
      link: readsTable
      to:   $table
    - from: $writer
      link: writesTable
      to:   $table
    - from: $writer
      link: jobOf
      to:   $writerApp

  where:
    - node: $writerApp
      does not equal:
        same as: $jobApp

  require:
    - node:  $table
      field: visibility
      equals:
        value: public
    - at least one:
        links:
          - from: $jobApp
            link: consumesContract
            to:   $contract
          - from: $contract
            link: coversTable
            to:   $table

  message: "{$job} reads {$table} from {$writerApp} without a contract."
```

Read it as: for every `$job`, `$jobApp`, `$table`, `$writer` and `$writerApp`
joined by these links, where the two applications differ, the table must be
public and at least one contract must link them.

#### Variables

- **Declared in `says`.** A word starting with `$` in the sentence declares a
  variable. A variable used anywhere else must appear in `says`; a variable in
  `says` that nothing uses is also an error.
- **Named** `$` then a lowercase letter, then letters and digits: `$job`,
  `$writerApp`.
- **Each variable stands for one thing** in a match. Two variables may match
  the same thing unless a condition says otherwise.
- **Kinds come from the links.** `$job` is a `Job` because `readsTable` goes
  from `Job`. If two links disagree — one says `Table`, the other
  `PartitionedTable` — the narrower kind wins; if neither kind extends the
  other, it is an error.
- **To narrow a kind** beyond what the links say, use a condition: `node:
  $table` with `is a: PartitionedTable`.
- **Scope.** A variable first used inside a nested block (`at least one`,
  `none`, `count`) belongs to that block and is not visible outside it.

#### Links in a pattern

Each entry in `links` is one link that must exist. The link must be declared
in the package or exported by a dependency, and `from` and `to` must be
variables. An inverse may be used
in either direction: `from: $table, link: writtenByJob, to: $writer` matches
exactly what `from: $writer, link: writesTable, to: $table` does.

`repeated: true` matches a link one or more times in a row. The things in
between are not named: they are the steps of the chain, and explanations list
them.

```yaml
- from: $source
  link: feedsTable
  to:   $target
  repeated: true
```

#### Conditions

Conditions appear in `where` and `require`. All entries in a list must hold.

A condition is about a field (`node` and `field`) or about a thing itself
(`node` alone), followed by a comparison word and then `value:` (a literal) or
`same as:` (a variable, or a variable and a field):

```yaml
- node:  $table
  field: retentionDays
  is at most:
    same as:
      node:  $database
      field: retentionDays
```

| Word | Meaning |
|---|---|
| `equals`, `does not equal` | equal, not equal |
| `is above`, `is below` | greater, smaller |
| `is at least`, `is at most` | greater or equal, smaller or equal |
| `is one of` | equal to one of a list (`values:` instead of `value:`) |
| `is a` | the thing is of this kind (takes a kind name directly) |
| `is set`, `is not set` | the field has a value, or has none |

A field with no value fails every comparison except `is not set`.

Patterns can nest inside conditions:

| Word | Meaning |
|---|---|
| `at least one` | the nested `links` (and its own `where`) match at least once |
| `none` | they never match |
| `count` | how many times they match, followed by a comparison word |
| `any of` | at least one of the listed conditions holds |

There is no general `not`: `does not equal`, `is not set` and `none` cover the
negative cases while keeping each condition readable on its own.

#### Derivations

A rule's `derive` lists facts to add for every match: a link (`from`, `link`,
`to`), a kind (`node` and `becomes`), or a field value (`node`, `field` and
`set to`). A derived field must end up with one value; two rules giving the
same thing two different values for one field is an error, reported with both
explanations.

```yaml
tableTeamFromDatabase:
  is: Rule
  says: A table $table belongs to the team $team that owns its database $database.

  links:
    - from: $table
      link: inDatabase
      to:   $database
    - from: $database
      link: ownedByTeam
      to:   $team

  derive:
    - from: $table
      link: belongsToTeam
      to:   $team

appDependsOnWriter:
  is: Rule
  says: >
    An application $app depends on every other application $other whose job
    $writer writes a table $table that a job $reader of $app reads.

  links:
    - from: $app
      link: hasJob
      to:   $reader
    - from: $reader
      link: readsTable
      to:   $table
    - from: $writer
      link: writesTable
      to:   $table
    - from: $writer
      link: jobOf
      to:   $other

  where:
    - node: $other
      does not equal:
        same as: $app

  derive:
    - from: $app
      link: dependsOnApp
      to:   $other

tableFeedsTable:
  is: Rule
  says: A table $source feeds every table $target that a job $job reading it writes.

  links:
    - from: $job
      link: readsTable
      to:   $source
    - from: $job
      link: writesTable
      to:   $target

  derive:
    - from: $source
      link: feedsTable
      to:   $target

criticalTable:
  is: Rule
  says: A table $table read by more than ten jobs $reader is critical.

  where:
    - count:
        links:
          - from: $reader
            link: readsTable
            to:   $table
      is above:
        value: 10

  derive:
    - node: $table
      becomes: CriticalTable
```

`feedsTable` is declared with `chains: true`, so the rule needs to derive only
the direct step; the chain follows. `criticalTable` has no top-level `links`:
with one variable outside the nested block, it matches every thing of that
variable's kind.

#### How they run

1. **Rules run over the stated facts,** each adding what it derives for every
   match, then again over the result, until a pass adds nothing. Rules never
   create things, merge things or remove facts, so this always ends.
2. **`none`, `count`, `does not equal` and `is not set` over derived facts wait.** A
   rule that uses them over something another rule derives runs only after
   that rule has finished. Rules that depend on each other through these words
   in a circle are rejected. (In Datalog terms, the rules must be stratified.)
3. **Constraints run last,** over stated and derived facts alike, and add
   nothing.
4. **Derived facts live in their own graph,** each recording the rule and the
   match that produced it.

Closed world: a link that is neither stated nor derived does not exist, and a
field with no value has none.

#### Explanations and messages

An explanation is the statement's match, filled in: each `links` entry becomes
a fact line and each condition a checked line. A line that is itself derived
unfolds into its own explanation, down to stated facts. A failed `at least
one` also lists what was searched.

`message` quotes variables in braces, filled in from the failing match:
`{$job}` gives the thing's id, `{$table: visibility}` a field's value. Only
variables of the top-level match may be quoted.

#### What is checked before anything runs

- A variable used but not in `says`, or in `says` but unused.
- A link, kind or field that is not declared.
- A link entry whose `from` or `to` cannot have that link's kind.
- A comparison between values of different types, or with a value a field does
  not allow.
- A derived field that two rules could set differently.
- Rules that depend on each other in a circle through negative words or
  `count`.
- A message quoting a variable that belongs to a nested block.

### Facts

Facts are records in a facts package. Each says what kind it is, its fields,
its links, and how far to trust it.

```yaml
import:
  - from: pkg:data-model
    use: [Table]

db/sales/table/order_lines:
  is: Table
  name: order_lines
  format: columnar
  visibility: public
  inDatabase: db/sales
  hasColumn: [order_id, sku, quantity]
  source: catalog-export/2026-11-10.json       # where this came from
  generated: { by: catalog-importer/2.1, at: 2026-11-10T09:12:00Z }
  verified:  { by: human:data-owner, at: 2026-11-11T14:00:00Z }
  staleAfter: 2027-02-01T00:00:00Z             # recheck after this
```

- **Trust levels** come from `verified`: unverified, confirmed by a machine,
  reviewed by a person (reason 17).
- **Stale facts** are flagged, not hidden.
- A derived fact is as trusted as the least trusted fact it came from.

### Tests

A types package proves itself with three kinds of test (reason 16):

```yaml
# tests/questions/who-breaks-if-order-intake-goes.yaml
asks: Which applications $app depend on the application $target?
given: tests/examples/two-apps-sharing-a-table.yaml
links:
  - from: $app
    link: dependsOnApp
    to:   $target
where:
  - node: $target
    equals:
      value: app/order-intake
expect:
  $app: [app/revenue-report]
```

```yaml
# tests/examples/cross-app-read-without-contract.yaml
expect:
  rejected:
    - crossAppReadNeedsContract
facts:
  - …
```

```yaml
# tests/examples/table-inherits-database-team.yaml
expect:
  derived:
    - db/sales/table/order_lines belongsToTeam team/commerce
facts:
  - …
```

### The change report

Every change — a pull request to a types or facts package, or a write through
the gateway — returns one report before it lands (reasons 12 and 18):

```
Stated (1)
  + job compute-revenue (app/revenue-report) readsTable order_lines

Derived (3)
  + app/revenue-report dependsOnApp app/order-intake
      rule    appDependsOnWriter
      match   $app = app/revenue-report, $reader = compute-revenue,
              $table = order_lines, $writer = ingest-orders,
              $other = app/order-intake
      links   app/revenue-report hasJob compute-revenue
              compute-revenue readsTable order_lines
              ingest-orders writesTable order_lines
              ingest-orders jobOf app/order-intake
      where   app/order-intake does not equal app/revenue-report
  + order_lines feedsTable revenue_daily
      rule    tableFeedsTable
  + (1 more, grouped: feedsTable)

No longer derived (0)

Rejected (1)
  ✗ crossAppReadNeedsContract
    compute-revenue reads order_lines from another application without a contract.
```

Large results are grouped by link, with new kinds of conclusion and
rejections listed first.

### The gateway

Running systems never write files. They write facts through a gateway that:

1. accepts facts only — a write that names a kind, link, constraint or rule
   as its subject is refused;
2. runs the constraints and rejects what fails;
3. runs the rules and returns the change report;
4. stamps `generated` and stores the facts in the writer's package graph.

Types arrive only by publishing a package version from CI.

### Product's reading

The same files generate a page per term for product: the definition,
examples and non-examples, every count as a sentence, and every constraint
and rule as its `says`, shown as written, variables included.

> **Table** — A set of rows with one schema, stored in a database.
> - Every table is in exactly one database.
> - Every table has at least one column.
> - A table $table belongs to the team $team that owns its database
>   $database. *(rule)*
> - A job $job of application $jobApp may read a table $table that job
>   $writer of another application $writerApp writes, only if $table is
>   public and $jobApp consumes a contract $contract that covers it.
>   *(constraint)*
>
> Example: the order_lines table in the sales database. Not a table: a
> temporary copy a test run makes and throws away.

A change to a package in a product-owned folder shows the sentences that
changed, for approval on that page.

### What it generates

| In the language | Exported as |
|---|---|
| package | an ontology with a version and `owl:imports`; its statements in their own named graph |
| kind | `owl:Class`, plus a closed `sh:NodeShape` for its fields and links |
| `extends` | `rdfs:subClassOf` |
| field | `owl:DatatypeProperty`, checked by a `sh:property` |
| link | `owl:ObjectProperty`; `from` and `to` as SHACL checks, never `rdfs:domain` or `rdfs:range` |
| `under` | `rdfs:subPropertyOf` |
| `inverse` | `owl:inverseOf` |
| `chains` | `owl:TransitiveProperty` |
| counts | `sh:minCount`, `sh:maxCount` — never an OWL cardinality restriction |
| constraint | a SHACL shape with a SPARQL (SPARQL Protocol and RDF Query Language) check |
| rule | a SHACL rule |
| definition, `reads` | `rdfs:comment`, and the generated sentence as a label for readers |
| facts | RDF triples in the package's named graph |
| derived facts | RDF triples in a derived graph, with PROV-O (the W3C provenance vocabulary) linking each to its rule |
| `source`, `generated`, `verified`, `staleAfter` | PROV-O statements on the fact's record |
| mappings | `owl:equivalentClass`, `owl:equivalentProperty`, `skos:closeMatch` |

## What it leaves out, and why

- **Open-world reasoning.** No "either this or that", no "something exists but
  is unnamed". It cannot be explained in a change report.
- **`owl:sameAs` and anything that merges things.** One mistaken statement
  would silently join two records.
- **Domain and range as inference.** They become checks.
- **OWL restrictions as constraints.** Counts and allowed values are SHACL.
- **Facts and types in one file.** Ruled out by the package split.

## Roads not taken

- **Write RDF, OWL and SHACL directly.** Exact and standard, and every reason
  in the table above applies in full.
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

## Open questions

1. **Trust per record, or per link?** A fact records where it came from and
   who confirmed it, once for the whole record. But one record can mix
   sources: `order_lines` came from the catalog export, and someone later added
   a link to it by hand. Is one source per record enough, or does each link
   need its own source and confirmation?
2. **Where do the Markdown docs live?** The draft puts `docs/` inside each
   package, so docs are versioned and reviewed with the model. The alternative
   is a separate set of pages beside the packages that link into them, the way
   OKF does it. Are docs part of the package or not?
3. **What does a facts package's version mean?** A types package has clear
   versions: removing a term is a breaking change. Facts change all the time —
   rows are added, jobs change, tables are dropped. Does a facts package have
   versions at all, and if so, what makes one version different from the
   next?
4. **Can a technical package add rules about product kinds?** `Table` lives
   in a product package. The platform team wants a rule such as "a table with
   more than 1,000 columns is oversized". Ownership comes from location, so
   the rule cannot go in the product package without product approving it. May
   the platform team put that rule in their own package, through an extension
   point, or must every rule about `Table` go through product?
5. **Optional links in a pattern.** Every link in a pattern must exist, or
   there is no match. A rule that wants "every table, and its team if it has
   one" loses the tables without a team. Should a link in a pattern be markable
   as optional, or should authors work around it with two rules or `any of`?
6. **Limits on chains.** `repeated: true` follows a link any number of steps:
   table to table to table. Sometimes only nearby things matter, such as
   "tables at most two steps upstream". Should chains take a limit
   (`at most: 2`), or stay unlimited?
7. **Variables only for the sentence.** An author might write "…owned by team
   $team…" in `says` because the sentence reads better, without using `$team`
   in the logic. The draft treats an unused variable as an error. Should it
   stay strict, or allow variables used only to make the sentence read well?
8. **Dates and durations.** Conditions compare a value with another value.
   They cannot say "a job that has not run for 30 days", because that needs
   "now minus 30 days". Should conditions gain date arithmetic, and if so, in
   what YAML form?
