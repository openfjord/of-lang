---
title: The language — A modeling language for software people
---

# The language

> Part of [a modeling language for software people](README.md), a first
> draft. Reason numbers refer to the table there.

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
5. **The standards are the export, not the source.** RDF, OWL and SHACL are
   generated for exchange and for existing tools; nobody writes them.

## Vocabulary

| Word | Meaning | Closest standard construct |
|---|---|---|
| **[package](terms/package.md)** | a versioned unit with one manifest and one owner, who owns everything in it; its name may start with a namespace, as in npm | an ontology with `owl:imports` |
| **[kind](terms/kind.md)** | a kind of thing (database, table, application) | `owl:Class` |
| **[field](terms/field.md)** | a value of the thing itself (a name, a format), described with JSON Schema | an attribute (`owl:DatatypeProperty`) |
| **[value type](terms/field.md#value-types)** | a named, reusable description of a field's value | `rdfs:Datatype` |
| **[link](terms/link.md)** | a named, directed relationship between two things | `owl:ObjectProperty` |
| **[constraint](terms/constraint.md)** | a rule that rejects data | a SHACL shape |
| **[rule](terms/rule.md)** | a pattern that derives new links, kinds or field values | a SHACL rule |
| **[pattern](terms/pattern.md)** | the links between things, each with its kind, and the conditions a constraint, rule or question looks for | a SPARQL graph pattern |
| **[fact](terms/fact.md)** | a stated thing with its fields and links; used in tests and examples | RDF triples |
| **[derived fact](terms/fact.md#derived-facts)** | a link produced by a rule | RDF triples in a separate graph |
| **[test](terms/test.md)** | facts with what they must derive and reject, or a question with the answer it expects | a competency question |

Fields and links stay apart: a field is a value of the thing, a link points
at another thing.

## Entries and imports

**Every entry says what it is.** Each top-level entry in a file is keyed by
its name and states what it is with `is:` — `Kind`, `Link`, `ValueType`,
`Constraint`, `Rule` or `Mapping` for definitions; `Fact`, `Test` or
`Question` in tests. Names are unique within a package.

```yaml
Table:
  is: Kind
  …

inDatabase:
  is: Link
  …

db/sales/table/order_lines:
  is: Fact
  kind: Table
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
  - from: pkg:@openfjord/org-model   # another package
    use: [Contract, { name: Team, as: OrgTeam }]
```

- `file:` names a file in the same package, relative to the importing file;
  `../` goes up a folder.
- `pkg:` names another package by its full name, namespace included. It must
  be in `dependsOn`, and only names it exports may be imported.
- Imports from the same package are required too, so the top of a file shows
  everything it depends on, and every name leads to its definition.
- `as` renames an import when two names clash.
- A link's inverse is a name like any other: `hasJob`, the inverse declared on
  `jobOf`, can be imported on its own.
- `import` is reserved: no entry may be named `import`.

## Ids

Ids are names, not links (reason 10). Files use plain names — `Table` for
a kind, `db/sales/table/order_lines` for a fact — and imports say where each
comes from. The full IRI, built from the package's `rdfNamespace`, appears only in
the exported RDF.
