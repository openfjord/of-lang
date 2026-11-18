---
title: At least one — The OpenFjord modeling language by example
---

# At least one

> Part of [the language by example](README.md).

A condition cannot see links directly. `at least one` holds a nested block
with its own `match` and `filter`, and holds when the block has one or more
matches.

<table>
<tr>
<td>

A table holds personal data or does not.

</td>
<td>

```yaml
# defs/product.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    name:
      type: string

Table:
  is: Kind
  definition: A table in a database.
  fields:
    name:
      type: string
    personalData:
      type: boolean
```

</td>
</tr>
<tr>
<td>

A service reads tables.

</td>
<td>

```yaml
readsTable:
  is: Link
  from: Service
  to: Table
  reads: "{from} reads {to}"
  definition: The service reads rows of the table.
```

</td>
</tr>
<tr>
<td>

A test file states two tables. Only the customers table holds personal data.

</td>
<td>

```yaml
# tests/examples/two-tables.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Table, readsTable]

table/products:
  is: Fact
  kind: Table
  fields:
    name: products
    personalData: false

table/customers:
  is: Fact
  kind: Table
  fields:
    name: customers
    personalData: true
```

</td>
</tr>
<tr>
<td>

The search service reads only the products table. The checkout service also
reads the customers table.

</td>
<td>

```yaml
service/search:
  is: Fact
  kind: Service
  fields:
    name: search
  links:
    readsTable: table/products

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
  links:
    readsTable: [table/products, table/customers]
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
twoTables:
  is: Test
  says: A service that reads personal data and a service that reads none are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services read a table that holds personal data. Its
`match` finds each service.

</td>
<td>

```yaml
# tests/questions/which-services-read-personal-data.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Table, readsTable]

whichServicesReadPersonalData:
  is: Question
  asks: Which services read a table that holds personal data?
  given: file:../examples/two-tables.of.yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

The nested block runs one time for each service. Its `match` finds the
tables that `$service` reads. The variable `$table` belongs to the nested
block, and nothing outside the block can use it.

</td>
<td>

```yaml
  filter:
    - at least one:
        match:
          - kind: { name: Table, as: $table }
          - link: readsTable
            from: $service
            to:   $table
```

</td>
</tr>
<tr>
<td>

The `filter` of the nested block keeps the tables that hold personal data.

</td>
<td>

```yaml
        filter:
          - node: $table
            field: personalData
            equals:
              value: true
```

</td>
</tr>
<tr>
<td>

For the checkout service, the nested block has one match. For the search
service, it has no matches.

</td>
<td>

```yaml
  expect:
    $service: [service/checkout]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test at-least-one
pass  test twoTables  (at-least-one/tests/examples/two-tables.of.yaml)
pass  question whichServicesReadPersonalData  (at-least-one/tests/questions/which-services-read-personal-data.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
