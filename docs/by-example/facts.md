---
title: Facts — The OpenFjord modeling language by example
---

# Facts

> Part of [the language by example](README.md).

A fact tells what a thing is: its kind, its fields and its links. Facts
occur only in the tests of a package, where they are the data for examples
and questions.

<table>
<tr>
<td>

The definitions tell what a fact can state. A service has a name and a
number of replicas, and it reads tables.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    name:
      type: string
    replicas:
      type: integer

Table:
  is: Kind
  definition: A table in a database.
  fields:
    name:
      type: string

readsTable:
  is: Link
  from: Service
  to: Table
  reads: "{from} reads {to}"
  definition: The service gets rows from the table.
```

</td>
</tr>
<tr>
<td>

Facts go in files in `tests/`. A file of facts imports each kind and each
link that its facts use. It does not import field names.

</td>
<td>

```yaml
# tests/examples/one-service.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, Table, readsTable]
```

</td>
</tr>
<tr>
<td>

The name of the statement is the id of the fact. An id is any text, and
usually a path. `kind` names the kind of the thing.

</td>
<td>

```yaml
table/customers:
  is: Fact
  kind: Table
```

</td>
</tr>
<tr>
<td>

`fields` gives the value of each field, by the name of the field.

</td>
<td>

```yaml
  fields:
    name: customers
```

</td>
</tr>
<tr>
<td>

A second table has the same form.

</td>
<td>

```yaml
table/orders:
  is: Fact
  kind: Table
  fields:
    name: orders
```

</td>
</tr>
<tr>
<td>

The service must give its two fields, because a field is required by
default.

</td>
<td>

```yaml
service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    replicas: 3
```

</td>
</tr>
<tr>
<td>

`links` gives one id or a list of ids for each link. Each id is the id of
a fact in the same file, and that fact must be of the kind where the link
ends.

</td>
<td>

```yaml
  links:
    readsTable: [table/customers, table/orders]
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test. This one expects that the package
rejects nothing: each fact obeys its kind.

</td>
<td>

```yaml
oneService:
  is: Test
  says: The checkout service reads the customers table and the orders table.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test facts
pass  test oneService  (facts/tests/examples/one-service.of.yaml)
1 passed, 0 failed
```

More about this: [facts](../../language/reference/fact.md) and [ids](../../language/reference/files.md#ids).
