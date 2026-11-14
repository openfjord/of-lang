---
title: of-lang — The OpenFjord modeling language
---

# `of-lang`

A command-line tool for [the OpenFjord modeling language](../language/README.md):
it checks packages, runs their rules over facts, explains what it derived,
compares two sets of facts, and runs a package's tests.

> **It lags the draft.** The tool implements an earlier version of the
> language: `links` and `where` instead of `match` and `filter`, variables
> named in `says`, kinds given only on a variable's first mention, and
> `definitions/` instead of `defs/`. It does not run the
> [example package](../examples/data-model/) as it stands. It will be brought
> up to date once the draft settles.

## Build

```bash
cd cli && go build -o of-lang ./cmd/of-lang
```

## Commands

```
of-lang check <package-dir>...
of-lang run <package-dir> <facts.yaml>
of-lang explain <package-dir> <facts.yaml> <subject> <link-or-kind> [object]
of-lang diff <package-dir> <before.yaml> <after.yaml>
of-lang test <package-dir>
```

| Command | Does |
|---|---|
| `check` | checks definitions and docs of one or more packages |
| `run` | derives what follows from a facts file, then checks constraints |
| `explain` | shows why one fact holds: a link (`subject link object`) or a kind (`subject Kind`) |
| `diff` | shows what changes between two facts files: stated, derived and rejected |
| `test` | runs every example and question in the package's `tests/` |

| Flag | Means |
|---|---|
| `--packages <dir>` | where to look for dependencies; repeatable; by default the folder that holds the package |
| `--format text\|json` | the output format; `text` by default |

| Exit code | Means |
|---|---|
| 0 | everything passes |
| 1 | problems found: a check failed, a constraint rejected something, or a test failed |
| 2 | the tool could not run: a bad command line, or a file it cannot read |

## Output

All examples below run in the folder that holds the package.

### `check`

```
$ of-lang check data-model
ok
```

Each problem is one line, with the file, line and column:

```
$ of-lang check data-model
data-model/defs/table.yaml:39:3: reads: of inDatabase must use both {from} and {to}
```

The draft has each error name the rule it enforces, such as `(LK2)`; the tool
does not yet.

### `run`

Every derived fact as its `reads` sentence, with what derived it: a rule, or
`under` or `chains` on a link. Then every rejection, with its message.

```
$ of-lang run data-model data-model/tests/examples/two-apps-sharing-a-table.yaml
derived:
  app/order_intake/job/ingest_orders uses db/sales/table/order_lines  (under)
  app/revenue_report/job/compute_revenue uses db/sales/table/order_lines  (under)
  app/revenue_report/job/compute_revenue uses db/sales/table/revenue_daily  (under)
  app/revenue_report depends on app/order_intake  (rule appDependsOnWriter)
  db/sales/table/order_lines belongs to team/commerce  (rule tableTeamFromDatabase)
  db/sales/table/revenue_daily belongs to team/commerce  (rule tableTeamFromDatabase)
  db/sales/table/order_lines feeds db/sales/table/revenue_daily  (rule tableFeedsTable)
rejected:
  nothing
```

```
$ of-lang run data-model data-model/tests/examples/cross-app-read-without-contract.yaml
derived:
  ...
rejected:
  crossAppReadNeedsContract: app/revenue_report/job/compute_revenue reads db/sales/table/order_lines from app/order_intake without a contract.
```

A rejection makes `run` exit with 1.

### `explain`

The fact, the statement's `says`, what each variable stood for, and the facts
the match used. A derived fact among them unfolds into its own explanation,
indented below it.

```
$ of-lang explain data-model data-model/tests/examples/two-apps-sharing-a-table.yaml app/revenue_report dependsOnApp app/order_intake
app/revenue_report depends on app/order_intake  (rule appDependsOnWriter)
  An application depends on every other application whose jobs write a table that one of its own jobs reads.
  with $app = app/revenue_report, $other = app/order_intake, $reader = app/revenue_report/job/compute_revenue, $table = db/sales/table/order_lines, $writer = app/order_intake/job/ingest_orders
  because:
    app/revenue_report/job/compute_revenue is a job of app/revenue_report  (written at data-model/tests/examples/two-apps-sharing-a-table.yaml:87:3)
    app/revenue_report/job/compute_revenue reads db/sales/table/order_lines  (written at data-model/tests/examples/two-apps-sharing-a-table.yaml:88:3)
    app/order_intake/job/ingest_orders writes db/sales/table/order_lines  (written at data-model/tests/examples/two-apps-sharing-a-table.yaml:81:3)
    app/order_intake/job/ingest_orders is a job of app/order_intake  (written at data-model/tests/examples/two-apps-sharing-a-table.yaml:80:3)
```

```
$ of-lang explain data-model data-model/tests/examples/critical-table.yaml db/sales/table/customers CriticalTable
db/sales/table/customers is a CriticalTable  (rule criticalTable)
  A table read by more than two jobs is critical.
  with $table = db/sales/table/customers
```

### `diff`

What changes from one facts file to the other, in three parts: facts stated,
facts derived, and rejections. `+` is added, `-` is gone. A field value is a
fact of its own, so a changed value shows as one `-` and one `+`.

```
$ of-lang diff data-model data-model/tests/examples/cross-app-read-without-contract.yaml data-model/tests/examples/two-apps-sharing-a-table.yaml
asserted:
  + app/revenue_report consumes contract/orders
  + app/revenue_report/job/compute_revenue writes db/sales/table/revenue_daily
  + contract/orders covers db/sales/table/order_lines
  + contract/orders has name orders
  + contract/orders is a Contract
  + db/sales/table/order_lines has visibility public
  + db/sales/table/revenue_daily has name revenue_daily
  + db/sales/table/revenue_daily has visibility internal
  + db/sales/table/revenue_daily is a Table
  + db/sales/table/revenue_daily is in db/sales
  - db/sales/table/order_lines has visibility internal
derived:
  + app/revenue_report/job/compute_revenue uses db/sales/table/revenue_daily
  + db/sales/table/order_lines feeds db/sales/table/revenue_daily
  + db/sales/table/revenue_daily belongs to team/commerce
rejected:
  - crossAppReadNeedsContract: app/revenue_report/job/compute_revenue reads db/sales/table/order_lines from app/order_intake without a contract.
```

`diff` always exits with 0: a difference is not a problem.

### `test`

One line per test or question, then a total. A failure says what was
expected and did not happen.

```
$ of-lang test data-model
pass  test tableReadByThreeJobs  (data-model/tests/examples/critical-table.yaml)
pass  test crossAppReadWithoutContract  (data-model/tests/examples/cross-app-read-without-contract.yaml)
FAIL  test twoAppsSharingATable  (data-model/tests/examples/two-apps-sharing-a-table.yaml)
      expected derived, but not: app/revenue_report depends on app/crm
pass  question whoBreaksIfOrderIntakeGoes  (data-model/tests/questions/who-breaks-if-order-intake-goes.yaml)
3 passed, 1 failed
```

### JSON

`--format json` prints the same results as one JSON object, for other tools
to read. A derived link, for example:

```json
{
  "type": "link",
  "sentence": "app/revenue_report depends on app/order_intake",
  "from": "app/revenue_report",
  "link": "dependsOnApp",
  "to": "app/order_intake"
}
```

and a rejection carries the constraint, its message, and what each variable
stood for:

```json
{
  "constraint": "crossAppReadNeedsContract",
  "message": "app/revenue_report/job/compute_revenue reads db/sales/table/order_lines from app/order_intake without a contract.",
  "bindings": [
    { "var": "$job", "id": "app/revenue_report/job/compute_revenue" }
  ]
}
```

## Planned outputs

Not built yet; the draft describes them so the language can be judged by what
it produces.

### The change report

Every change to a package returns one report, for the facts in its tests and
examples, before it lands. It joins `diff` and `explain`: what was stated,
what newly follows and why, what no longer follows, and what is rejected.

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
      filter  app/order-intake does not equal app/revenue-report
  + order_lines feedsTable revenue_daily
      rule    tableFeedsTable
  + (1 more, grouped: feedsTable)

No longer derived (0)

Rejected (1)
  ✗ crossAppReadNeedsContract
    compute-revenue must not read order_lines from app/order-intake without a contract.
```

Large results are grouped by link, with new kinds of conclusion and
rejections listed first.

### Product's pages

A page per definition for product: the definition, examples and
non-examples, every count as a sentence, and every constraint and rule as its
`says`, shown as written.

> **Table** — A set of rows with one schema, stored in a database.
> - Every table is in exactly one database.
> - A table belongs to the team that owns its database. *(rule)*
> - A job of one application may read a table that a job of another
>   application writes, only if the table is public and the reading
>   application consumes a contract that covers the table.
>   *(constraint)*
>
> Example: the order_lines table in the sales database. Not a table: a
> temporary copy a test run makes and throws away.

A change to a product-owned package shows the sentences that changed, for
approval on that page.

### The RDF export

RDF (Resource Description Framework), OWL (Web Ontology Language) and SHACL
(Shapes Constraint Language) are generated for exchange and for existing
tools. Nobody writes them.

| In the language | Exported as |
|---|---|
| package | an ontology with a version and `owl:imports`; its statements in their own named graph |
| kind | `owl:Class`, plus a closed `sh:NodeShape` for its fields and links |
| `extends` | `rdfs:subClassOf` |
| field | `owl:DatatypeProperty`, checked by a `sh:property` built from its JSON Schema keywords |
| value type | `rdfs:Datatype` with restrictions (`xsd:pattern`, `xsd:maxLength`, `xsd:minInclusive` and the like) |
| `string`, `integer`, `number`, `boolean`, `date`, `date-time`, `duration` | `xsd:string`, `xsd:integer`, `xsd:decimal`, `xsd:boolean`, `xsd:date`, `xsd:dateTime`, `xsd:duration` |
| link | `owl:ObjectProperty`; `from` and `to` as SHACL checks, never `rdfs:domain` or `rdfs:range` |
| `under` | `rdfs:subPropertyOf` |
| `inverse` | `owl:inverseOf` |
| `chains` | `owl:TransitiveProperty` |
| counts | `sh:minCount`, `sh:maxCount`, never an OWL cardinality restriction |
| constraint | a SHACL shape with a SPARQL (SPARQL Protocol and RDF Query Language) check |
| rule | a SHACL rule |
| definition, `reads` | `rdfs:comment`, and the generated sentence as a label for readers |
| facts | RDF triples in the package's named graph |
| derived facts | RDF triples in a derived graph, with PROV-O (the W3C provenance vocabulary) linking each to its rule |

Each IRI (Internationalized Resource Identifier) is built from the package's
`rdfNamespace` and the definition's name.
