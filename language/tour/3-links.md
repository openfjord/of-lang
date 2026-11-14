---
title: 3. Links — Tour
---

# 3. Links

> Part of [the tour](README.md). Previous: [2. Facts and a test](2-facts-and-a-test.md).
> Next: [4. Rules](4-rules.md).

A field is a value of the thing itself. A link points at another thing.

## A link

A database is owned by a team. Write `defs/database.yaml`:

```yaml
import:
  - from: file:names.yaml
    use: [Name, RetentionDays]
  - from: file:team.yaml
    use: [Team]

Database:
  is: Kind
  definition: A named store of tables, owned by one team.
  fields:
    name:
      type: Name
    retentionDays:
      type: RetentionDays
      optional: true
  links:
    ownedByTeam: exactly 1

ownedByTeam:
  is: Link
  from: Database
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that answers for the database.
  inverse: ownsDatabase
```

(`RetentionDays` is a second value type in `names.yaml`: an integer from 1 to
3650.)

The link is its own entry, next to the kind it comes out of:

- `from` and `to` are the kinds it joins. They are checks: linking a table
  with `ownedByTeam` is rejected.
- `reads` is the sentence product sees: "db/sales is owned by team/commerce".
- `inverse` names the link read the other way. A team `ownsDatabase` exactly
  what is `ownedByTeam` it, with nothing extra to state.

The kind says how many each thing carries: `exactly 1`, `at least N`,
`at most N` or `any number`. A count is a check. A database with two owning
teams is rejected; it is never read as "the two teams are the same".

Name links with two words or more, saying what they link: `ownedByTeam`,
never `owner` or `relatedTo`. A vague link adds lines and nothing can be
concluded from it.

## Tables

`defs/table.yaml` adds tables, which sit in a database:

```yaml
import:
  - from: file:names.yaml
    use: [Name, RetentionDays]
  - from: file:database.yaml
    use: [Database]

Table:
  is: Kind
  definition: A set of rows with one schema, stored in a database.
  examples:
    - the order_lines table in the sales database
  notExamples:
    - a temporary copy a test run makes and throws away
    - a view computed on read
  fields:
    name:
      type: Name
    visibility:
      enum: [public, internal]
      default: internal
  links:
    inDatabase: exactly 1

inDatabase:
  is: Link
  from: Table
  to: Database
  reads: "{from} is in {to}"
  definition: The database that stores the table.
  inverse: hasTable
```

`examples` and `notExamples` show up on product's page for the kind. A good
non-example is something people would mistake for a table.

## Links in facts

A fact states a link like a field, by the id of the thing it points at:

```yaml
db/sales:
  is: Fact
  kind: Database
  name: sales
  ownedByTeam: team/commerce

db/sales/table/order_lines:
  is: Fact
  kind: Table
  name: order_lines
  inDatabase: db/sales
```

Leave out `ownedByTeam` and the test fails: a database needs exactly one
owning team.

```
FAIL  test aDatabase  (data-model/tests/examples/a-database.yaml)
      db/sales: ownedByTeam needs exactly 1, has 0 (FC3)
```
