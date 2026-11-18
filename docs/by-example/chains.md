---
title: Links that chain — The OpenFjord modeling language by example
---

# Links that chain

> Part of [the language by example](README.md).

A link with `chains: true` follows through: when a first thing has the link
to a second thing, and the second thing has it to a third thing, the first
thing has it to the third thing.

<table>
<tr>
<td>

A service can depend on another service. The link starts and ends at the
same kind.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

dependsOn:
  is: Link
  from: Service
  to: Service
  reads: "{from} depends on {to}"
  definition: The first service cannot do its work without the second service, directly or through other services.
```

</td>
</tr>
<tr>
<td>

`chains: true` tells that a chain of the link implies the link. In a link
that chains, `to` must be the same kind as `from`, or a kind that extends
it.

</td>
<td>

```yaml
  chains: true
```

</td>
</tr>
<tr>
<td>

A test file states three services. The checkout service depends on the
orders service, and the orders service depends on the database.

</td>
<td>

```yaml
# tests/examples/checkout-needs-the-database.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, dependsOn]

service/database:
  is: Fact
  kind: Service

service/orders:
  is: Fact
  kind: Service
  links:
    dependsOn: service/database

service/checkout:
  is: Fact
  kind: Service
  links:
    dependsOn: service/orders
```

</td>
</tr>
<tr>
<td>

No fact states that the checkout service depends on the database. The test
expects it as a derived fact.

</td>
<td>

```yaml
checkoutNeedsTheDatabase:
  is: Test
  says: Checkout depends on orders, which depends on the database. Thus, checkout depends on the database.
  expect:
    derived:
      - link: dependsOn
        from: service/checkout
        to:   service/database
```

</td>
</tr>
<tr>
<td>

The link stays in one direction: the database does not depend on the
checkout service.

</td>
<td>

```yaml
    notDerived:
      - link: dependsOn
        from: service/database
        to:   service/checkout
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test chains
pass  test checkoutNeedsTheDatabase  (chains/tests/examples/checkout-needs-the-database.of.yaml)
1 passed, 0 failed
```

More about this: [what links imply](../../language/reference/link.md#what-links-imply) and [the rules for links](../../language/reference/link.md#rules).
