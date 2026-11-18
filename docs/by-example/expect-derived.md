---
title: Expect derived facts — The OpenFjord modeling language by example
---

# Expect derived facts

> Part of [the language by example](README.md).

A test tells what a set of facts must derive. `derived` lists facts that
must hold, and `notDerived` lists facts that must not hold.

<table>
<tr>
<td>

The package needs something that derives a fact. Here, `dependsOn` chains:
a service also depends on each service that its dependencies depend on.

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
  chains: true
```

</td>
</tr>
<tr>
<td>

An example is a file of facts with one `Test` statement. The checkout
service depends on the orders service, and the orders service depends on
the database. The search service depends on nothing.

</td>
<td>

```yaml
# tests/examples/what-checkout-depends-on.of.yaml
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

service/search:
  is: Fact
  kind: Service
```

</td>
</tr>
<tr>
<td>

A test must have `says` and `expect`. `says` tells in plain words what the
example shows.

</td>
<td>

```yaml
whatCheckoutDependsOn:
  is: Test
  says: >
    Checkout depends on orders, which depends on the database. Thus,
    checkout depends on the database. Search does not depend on the
    database.
  expect:
```

</td>
</tr>
<tr>
<td>

`derived` lists the facts that must hold after the package derives all that
it can. A link is written with `link`, `from` and `to`, and the two ends are
ids of facts.

</td>
<td>

```yaml
    derived:
      - link: dependsOn
        from: service/checkout
        to:   service/database
```

</td>
</tr>
<tr>
<td>

`notDerived` lists the facts that must not hold. The test fails if one of
them holds.

</td>
<td>

```yaml
    notDerived:
      - link: dependsOn
        from: service/search
        to:   service/database
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test expect-derived
pass  test whatCheckoutDependsOn  (expect-derived/tests/examples/what-checkout-depends-on.of.yaml)
1 passed, 0 failed
```

More about this: [examples](../../language/reference/test.md#examples) and [derived facts](../../language/reference/fact.md#derived-facts).
