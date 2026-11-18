---
title: Links under links — The OpenFjord modeling language by example
---

# Links under links

> Part of [the language by example](README.md).

A link can be `under` a broader link. Then every such link is also the
broader link.

<table>
<tr>
<td>

A service uses tables. `usesTable` is the broader link.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

Table:
  is: Kind
  definition: A table in a database.

usesTable:
  is: Link
  from: Service
  to: Table
  reads: "{from} uses {to}"
  definition: The service reads the table, writes it, or does both.
```

</td>
</tr>
<tr>
<td>

To write a table is one way to use it. `writesTable` is a narrower sort of
`usesTable`.

</td>
<td>

```yaml
writesTable:
  is: Link
  from: Service
  to: Table
  reads: "{from} writes {to}"
  definition: The service adds rows to the table or changes its rows.
```

</td>
</tr>
<tr>
<td>

`under` names the broader link. The two ends of the narrower link must be
the kinds of the broader link, or kinds that extend them.

</td>
<td>

```yaml
  under: usesTable
```

</td>
</tr>
<tr>
<td>

A test file states that the checkout service writes the customers table.
It states no `usesTable`.

</td>
<td>

```yaml
# tests/examples/checkout-writes-customers.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, Table, usesTable, writesTable]

table/customers:
  is: Fact
  kind: Table

service/checkout:
  is: Fact
  kind: Service
  links:
    writesTable: table/customers
```

</td>
</tr>
<tr>
<td>

The test expects the broader link as a derived fact: the checkout service
uses the customers table.

</td>
<td>

```yaml
checkoutWritesCustomers:
  is: Test
  says: The checkout service writes the customers table, so it uses that table.
  expect:
    derived:
      - link: usesTable
        from: service/checkout
        to:   table/customers
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test under
pass  test checkoutWritesCustomers  (under/tests/examples/checkout-writes-customers.of.yaml)
1 passed, 0 failed
```

More about this: [what links imply](../../language/reference/link.md#what-links-imply) and [the rules for links](../../language/reference/link.md#rules).
