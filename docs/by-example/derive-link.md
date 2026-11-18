---
title: Derive a link — The OpenFjord modeling language by example
---

# Derive a link

> Part of [the language by example](README.md).

A derivative adds facts that follow from other facts. This one adds a link
between two things.

<table>
<tr>
<td>

The package has three kinds: a service, a table and a data class.

</td>
<td>

```yaml
# defs/data.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

Table:
  is: Kind
  definition: A table in a database.

DataClass:
  is: Kind
  definition: A sort of data that needs care, such as personal data or payment data.
```

</td>
</tr>
<tr>
<td>

Facts state two links: a service reads a table, and a table holds a data
class.

</td>
<td>

```yaml
readsTable:
  is: Link
  from: Service
  to: Table
  reads: "{from} reads {to}"
  definition: The service gets rows from the table.

holdsData:
  is: Link
  from: Table
  to: DataClass
  reads: "{from} holds {to}"
  definition: The rows of the table contain data of the data class.
```

</td>
</tr>
<tr>
<td>

No fact states the third link. A derivative adds it.

</td>
<td>

```yaml
handlesData:
  is: Link
  from: Service
  to: DataClass
  reads: "{from} handles {to}"
  definition: Data of the data class passes through the service.
```

</td>
</tr>
<tr>
<td>

A derivative has a `says` sentence that tells what it concludes, in plain
words.

</td>
<td>

```yaml
dataFromTable:
  is: Derivative
  says: A service handles each data class that one of the tables it reads holds.
```

</td>
</tr>
<tr>
<td>

`match` finds each service that reads a table that holds a data class.
A variable, such as `$service`, represents each thing of its kind in turn.

</td>
<td>

```yaml
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Table, as: $table }
    - kind: { name: DataClass, as: $data }
    - link: readsTable
      from: $service
      to:   $table
    - link: holdsData
      from: $table
      to:   $data
```

</td>
</tr>
<tr>
<td>

`derive` tells what to add for each match. An item with `link`, `from`
and `to` adds a link.

</td>
<td>

```yaml
  derive:
    - link: handlesData
      from: $service
      to:   $data
```

</td>
</tr>
<tr>
<td>

A test file states the facts. The checkout service reads the customers
table, and the customers table holds personal data.

</td>
<td>

```yaml
# tests/examples/checkout.of.yaml
import:
  - from: file:../../defs/data.of.yaml
    use: [Service, Table, DataClass, readsTable, holdsData, handlesData]

data/personal:
  is: Fact
  kind: DataClass

table/customers:
  is: Fact
  kind: Table
  links:
    holdsData: data/personal

service/checkout:
  is: Fact
  kind: Service
  links:
    readsTable: table/customers
```

</td>
</tr>
<tr>
<td>

The test expects the link that the derivative adds. A fact under `derived`
has the same form as an item under `derive`, with ids in place of
variables.

</td>
<td>

```yaml
checkout:
  is: Test
  says: The checkout service reads the customers table, which holds personal data. Thus, the checkout service handles personal data.
  expect:
    derived:
      - link: handlesData
        from: service/checkout
        to:   data/personal
```

</td>
</tr>
</table>

`run` shows each derived fact as the `reads` sentence of its link, with the
derivative that added it:

```console
$ of-lang run derive-link derive-link/tests/examples/checkout.of.yaml
derived:
  service/checkout handles data/personal  (derivative dataFromTable)
rejected:
  nothing
```

More about this: [derivatives](../../language/reference/derivative.md) and
[derive](../../language/reference/derivative.md#derive).
