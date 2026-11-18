---
title: Unions — The OpenFjord modeling language by example
---

# Unions

> Part of [the language by example](README.md).

Sometimes a kind is exactly one of a fixed list of kinds, and nothing else.
`oneOf` on a kind states this. Such a kind is a union.

<table>
<tr>
<td>

The product has services and tables. A service uses the two of them to do
its work.

</td>
<td>

```yaml
# defs/resource.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    name:
      type: string

Table:
  is: Kind
  definition: A set of rows with the same columns, in a database.
  fields:
    name:
      type: string
```

</td>
</tr>
<tr>
<td>

`oneOf` lists the kinds that together make up `Resource`. Each resource is
a service or a table, and never both. A union has no fields. The listed
kinds keep their own fields.

</td>
<td>

```yaml
Resource:
  is: Kind
  definition: Something that a service needs to do its work.
  oneOf: [Service, Table]
```

</td>
</tr>
<tr>
<td>

A link can end at the union. Then it can end at a service or at a table.

</td>
<td>

```yaml
usesResource:
  is: Link
  from: Service
  to: Resource
  reads: "{from} uses {to}"
  definition: The service needs the resource to do its work.
```

</td>
</tr>
<tr>
<td>

A test file states one table and two services. A fact names a listed kind
as its `kind`, and never the union.

</td>
<td>

```yaml
# tests/examples/a-checkout.of.yaml
import:
  - from: file:../../defs/resource.of.yaml
    use: [Service, Table, Resource, usesResource]

table/customers:
  is: Fact
  kind: Table
  fields:
    name: customers

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
```

</td>
</tr>
<tr>
<td>

The checkout service uses a service and a table through the same link.

</td>
<td>

```yaml
service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
  links:
    usesResource: [service/orders, table/customers]
```

</td>
</tr>
<tr>
<td>

The test expects what the definitions derive about the kinds: a service is
a resource, and a service is never a table.

</td>
<td>

```yaml
aCheckout:
  is: Test
  says: The checkout service uses a service and a table, and each one is a resource.
  expect:
    derived:
      - kind: Service
        is a: Resource
      - kind: Service
        never is: Table
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question can match the union. It asks for each resource that the
checkout service uses, and the answer has a service and a table.

</td>
<td>

```yaml
# tests/questions/which-resources-does-checkout-use.of.yaml
import:
  - from: file:../../defs/resource.of.yaml
    use: [Service, Resource, usesResource]

whichResourcesDoesCheckoutUse:
  is: Question
  asks: Which resources does the checkout service use?
  given: file:../examples/a-checkout.of.yaml
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Resource, as: $resource }
    - link: usesResource
      from: $service
      to:   $resource
  filter:
    - node: $service
      equals:
        value: service/checkout
  expect:
    $resource: [service/orders, table/customers]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test unions
pass  test aCheckout  (unions/tests/examples/a-checkout.of.yaml)
pass  question whichResourcesDoesCheckoutUse  (unions/tests/questions/which-resources-does-checkout-use.of.yaml)
2 passed, 0 failed
```

More about this: [unions](../../language/reference/kind.md#unions) and
[derived terms](../../language/reference/test.md#derived-terms).
