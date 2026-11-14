---
title: What it outputs — A modeling language for software people
---

# What it outputs

> Part of [a modeling language for software people](README.md), a first
> draft. Reason numbers refer to the table there.

## The change report

Every change to a package returns one report, for the facts in its tests and
examples, before it lands (reasons 12 and 18):

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
    compute-revenue reads order_lines from another application without a contract.
```

Large results are grouped by link, with new kinds of conclusion and
rejections listed first.

## Product's reading

The same files generate a page per term for product: the definition,
examples and non-examples, every count as a sentence, and every constraint
and rule as its `says`, shown as written.

> **Table** — A set of rows with one schema, stored in a database.
> - Every table is in exactly one database.
> - Every table has at least one column.
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

## What it generates

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
| counts | `sh:minCount`, `sh:maxCount` — never an OWL cardinality restriction |
| constraint | a SHACL shape with a SPARQL (SPARQL Protocol and RDF Query Language) check |
| rule | a SHACL rule |
| definition, `reads` | `rdfs:comment`, and the generated sentence as a label for readers |
| facts | RDF triples in the package's named graph |
| derived facts | RDF triples in a derived graph, with PROV-O (the W3C provenance vocabulary) linking each to its rule |
| mappings | `owl:equivalentClass`, `owl:equivalentProperty`, `skos:closeMatch` |
