---
title: Unless — The OpenFjord modeling language by example
---

# Unless

> Part of [the language by example](README.md).

`unless` gives the conditions that excuse a match of a constraint. A match
that satisfies them is not a violation.

<table>
<tr>
<td>

Each table states how many days it keeps a row, and whether the law
requires the company to keep its rows.

</td>
<td>

```yaml
# defs/table.of.yaml
Table:
  is: Kind
  definition: A table in a database.
  fields:
    retentionDays:
      type: integer
    requiredByLaw:
      type: boolean
```

</td>
</tr>
<tr>
<td>

`match` and `filter` select the matches that the constraint applies to:
each table that keeps a row for more than 90 days.

</td>
<td>

```yaml
tableDeletesOldRows:
  is: Constraint
  says: >
    A table keeps a row for 90 days or less, unless the law requires the
    company to keep its rows.
  match:
    - kind: { name: Table, as: $table }
  filter:
    - node:  $table
      field: retentionDays
      is above:
        value: 90
```

</td>
</tr>
<tr>
<td>

`unless` tells which of these matches to excuse. Its conditions have the
same form as the conditions under `filter`. It excuses a match only when all
its conditions hold.

</td>
<td>

```yaml
  unless:
    - node:  $table
      field: requiredByLaw
      equals:
        value: true
```

</td>
</tr>
<tr>
<td>

The message is for the matches that `unless` does not excuse.

</td>
<td>

```yaml
  message: "{$table} must keep a row for 90 days or less, unless the law requires its rows."
```

</td>
</tr>
<tr>
<td>

A test file states two tables that keep a row for more than 90 days. The
law requires the company to keep the rows of the invoices table.

</td>
<td>

```yaml
# tests/examples/two-old-tables.of.yaml
import:
  - from: file:../../defs/table.of.yaml
    use: [Table, tableDeletesOldRows]

table/invoices:
  is: Fact
  kind: Table
  fields:
    retentionDays: 3650
    requiredByLaw: true
```

</td>
</tr>
<tr>
<td>

The law does not require the rows of the customers table.

</td>
<td>

```yaml
table/customers:
  is: Fact
  kind: Table
  fields:
    retentionDays: 365
    requiredByLaw: false
```

</td>
</tr>
<tr>
<td>

The facts break the constraint, because the customers table is a
violation.

</td>
<td>

```yaml
twoOldTables:
  is: Test
  says: >
    The invoices table and the customers table both keep a row for more
    than 90 days. The law requires the rows of the invoices table, so the
    constraint rejects only the customers table.
  expect:
    rejected: [tableDeletesOldRows]
```

</td>
</tr>
</table>

`run` reports one violation. The invoices table is a match, but `unless`
excuses it:

```console
$ of-lang run unless unless/tests/examples/two-old-tables.of.yaml
derived:
  nothing
rejected:
  tableDeletesOldRows: table/customers must keep a row for 90 days or less, unless the law requires its rows.
```

More about this: [filter and unless](../../language/reference/constraint.md#filter-and-unless).
