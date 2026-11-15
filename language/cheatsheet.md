---
title: Cheat sheet — The OpenFjord modeling language
---

# Cheat sheet

> Part of [the OpenFjord modeling language](README.md). All of the syntax is
> on one page. [The reference](reference/README.md) gives the exact rules.

## A package

The name of every file in `defs/` and `tests/` ends in `.of.yaml`. The
manifest is always `package.yaml`.

```
my-package/
  package.yaml       # required: the manifest
  defs/              # required: Kind, Link, ValueType, Tag, Constraint, Derivative
  docs/              # optional: Markdown that links to or embeds definitions
  tests/             # optional: Fact, Test, Question
```

```yaml
# package.yaml
package: "@openfjord/data-model"   # required; quoted because of the leading @
version: 0.3.0                     # required
owner: product                     # required; owns everything in the package
rdfNamespace: data                 # required; used only in exported RDF
dependsOn:
  "@openfjord/org-model": ^1.2
exports:
  kinds: [Table]
  links: [inDatabase]
  valueTypes: [Name]
  tags: [lineage]
openFor:
  - extend: Table
  - under: readsTable
```

## A file

```yaml
import:
  - from: file:names.of.yaml                     # relative to this file
    use: [Name]
  - from: pkg:@openfjord/org-model            # must be in dependsOn and exported there
    use: [Contract, { name: Team, as: OrgTeam }]

SomeStatement:
  is: Kind        # or Link, ValueType, Tag, Constraint, Derivative — Fact, Test, Question in tests/
```

A file imports each name that it uses but does not define. A file does not
import field names.

All free-form text (`definition`, `reads`, `examples`, `says`, `asks`,
`message`, docs) follows ASD-STE100 Simplified Technical English.

## Definitions

```yaml
Table:                                   # UpperCamelCase
  is: Kind
  definition: A set of rows with one schema, stored in a database.   # required
  extends: Dataset                       # optional
  tags: [lineage]                        # optional; kinds, links and value types only
  examples:    [the order_lines table]
  notExamples: [a view computed on read]
  fields:
    name:          { type: Name }                          # required by default
    visibility:    { enum: [public, internal], default: internal }
    retentionDays: { type: integer, minimum: 1, optional: true }
  links:
    inDatabase:    exactly 1             # exactly N | at least N | at most N | any number
    feedsTable:    any number

inDatabase:                              # lowerCamelCase, two words or more
  is: Link
  from: Table
  to: Database
  reads: "{from} is in {to}"             # required, with both {from} and {to}
  definition: The database that stores the table.
  inverse: hasTable                      # optional: the name read the other way
  under: locatedIn                       # optional: every inDatabase is a locatedIn
  chains: true                           # optional: a→b and b→c means a→c

Name:
  is: ValueType
  definition: A lowercase identifier.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"

lineage:                                 # lowerCamelCase
  is: Tag
  definition: What a lineage tool needs to trace data from table to table.
```

Tags only select what an RDF export holds. Tags do not pass through
`extends` or `under`. The export does not include a reference to an untagged
definition. It also does not include a link whose `from` or `to` kind is
untagged.

Field keywords: `type` (`string` `integer` `number` `boolean` `array` `object`
or a value type), `enum` `const` `default` `pattern` `minLength` `maxLength`
`format` (`date` `date-time` `duration`) `minimum` `maximum` `items`
`properties` `optional`. The language does not have `null`, `$ref`,
`required`, `oneOf` or similar keywords.

## Match and filter

Derivatives, constraints and questions each define their own `match` and
`filter`. All three use the same syntax, but each can change apart from the
others.

```yaml
match:                                   # links that must exist
  - link: readsTable                     # a link or an inverse name
    from: {$job: Job}                    # every mention gives the variable's kind
    to:   {$table: Table}
  - link: feedsTable
    from: {$table: Table}
    to:   {$next: Table}
    repeated: true                       # one or more steps
filter:                                  # conditions; all must hold
  - node: $next                          # bare variable outside match
    does not equal: { same as: $table }
  - node: $table
    field: visibility
    equals: { value: public }
```

| Comparison | Takes |
|---|---|
| `equals` `does not equal` `is above` `is below` `is at least` `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set` `is not set` | (only a field) |

| Nested | Holds when |
|---|---|
| `at least one: { match, filter }` | it matches at least once |
| `none: { match, filter }` | it never matches |
| `count: { match, filter }` + a comparison | the number of matches satisfies the comparison |
| `any of: [conditions]` | one of them holds |

## Constraints and derivatives

```yaml
crossAppReadNeedsContract:
  is: Constraint
  says: A job may read another application's table only under a contract.   # no $ words
  match: [ ... ]
  filter: [ ... ]
  unless: [ ... ]                        # optional; a match is a violation unless these hold
  message: "{$job} must not read {$table} without a contract."

tableTeamFromDatabase:
  is: Derivative
  says: A table belongs to the team that owns its database.
  match: [ ... ]
  derive:
    - { link: belongsToTeam, from: $table, to: $team }             # a link
    - { node: $table, becomes: CriticalTable }                     # a kind
    - { node: $table, field: visibility, set to: { value: internal } }   # a field value
```

## Tests

```yaml
# tests/examples/some-example.of.yaml — facts plus exactly one Test
db/sales:                                # a fact's id
  is: Fact
  kind: Database
  name: sales
  ownedByTeam: team/commerce             # a link: an id, or a list of ids

someExample:
  is: Test
  says: What this example shows.
  expect:
    derived:    [ { link: belongsToTeam, from: db/sales/table/x, to: team/commerce } ]
    notDerived: [ { node: db/sales/table/x, becomes: CriticalTable } ]
    rejected:   [crossAppReadNeedsContract]

# tests/questions/some-question.of.yaml
someQuestion:
  is: Question
  asks: Which applications depend on order intake?
  given: file:../examples/some-example.of.yaml
  match: [ ... ]
  filter: [ ... ]
  expect:
    $app: [app/revenue_report]
```

## How it runs

The derivatives run repeatedly until no new fact follows. `under` and `chains`
also add their facts. Then the constraints check everything. Anything that is
not stated or derived is false. Every derived fact and every rejection has
its explanation.
