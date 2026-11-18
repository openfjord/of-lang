---
title: Derivatives together — The OpenFjord modeling language by example
---

# Derivatives together

> Part of [the language by example](README.md).

Derivatives run again and again, until a pass adds nothing. Thus a
derivative can match a link that a derivative adds.

<table>
<tr>
<td>

A service can store a data class, and a service can depend on another
service. Facts state these two links.

</td>
<td>

```yaml
# defs/data.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

DataClass:
  is: Kind
  definition: A sort of data that needs care, such as personal data or payment data.

storesData:
  is: Link
  from: Service
  to: DataClass
  reads: "{from} stores {to}"
  definition: The service keeps data of the data class.

dependsOn:
  is: Link
  from: Service
  to: Service
  reads: "{from} depends on {to}"
  definition: The first service calls the second service directly.
```

</td>
</tr>
<tr>
<td>

No fact states `handlesData`. Two derivatives add it.

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

The first derivative matches only stated facts.

</td>
<td>

```yaml
dataFromStore:
  is: Derivative
  says: A service handles each data class that it stores.
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: DataClass, as: $data }
    - link: storesData
      from: $service
      to:   $data
  derive:
    - link: handlesData
      from: $service
      to:   $data
```

</td>
</tr>
<tr>
<td>

The second derivative matches `handlesData`, the link that the first
derivative adds.

</td>
<td>

```yaml
dataFromDependency:
  is: Derivative
  says: A service handles each data class that one of the services it depends on handles.
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Service, as: $other }
    - kind: { name: DataClass, as: $data }
    - link: dependsOn
      from: $service
      to:   $other
    - link: handlesData
      from: $other
      to:   $data
```

</td>
</tr>
<tr>
<td>

It also adds `handlesData`. Thus one pass can add a fact that the same
derivative matches in the pass after it.

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

A test file states a chain of three services. The checkout service depends
on the orders service, the orders service depends on the database, and the
database stores personal data.

</td>
<td>

```yaml
# tests/examples/data-in-the-database.of.yaml
import:
  - from: file:../../defs/data.of.yaml
    use: [Service, DataClass, storesData, dependsOn, handlesData]

data/personal:
  is: Fact
  kind: DataClass

service/database:
  is: Fact
  kind: Service
  links:
    storesData: data/personal

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

The test expects a link that no single pass can add. The orders service
must handle personal data before the derivative can find that the checkout
service handles it.

</td>
<td>

```yaml
dataInTheDatabase:
  is: Test
  says: >
    The database stores personal data. The checkout service depends on the
    orders service, which depends on the database. Thus, the checkout
    service handles personal data.
  expect:
    derived:
      - link: handlesData
        from: service/checkout
        to:   data/personal
```

</td>
</tr>
</table>

`explain` shows why a fact holds, down to the facts that the file states:

```console
$ of-lang explain derivatives-together derivatives-together/tests/examples/data-in-the-database.of.yaml service/checkout handlesData data/personal
service/checkout handles data/personal  (derivative dataFromDependency)
  A service handles each data class that one of the services it depends on handles.
  with $data = data/personal, $other = service/orders, $service = service/checkout
  because:
    service/checkout depends on service/orders  (written at derivatives-together/tests/examples/data-in-the-database.of.yaml:25:5)
    service/orders handles data/personal  (derivative dataFromDependency)
      A service handles each data class that one of the services it depends on handles.
      with $data = data/personal, $other = service/database, $service = service/orders
      because:
        service/orders depends on service/database  (written at derivatives-together/tests/examples/data-in-the-database.of.yaml:19:5)
        service/database handles data/personal  (derivative dataFromStore)
          A service handles each data class that it stores.
          with $data = data/personal, $service = service/database
          because:
            service/database stores data/personal  (written at derivatives-together/tests/examples/data-in-the-database.of.yaml:13:5)
```

More about this: [how derivatives run](../../language/reference/derivative.md#how-derivatives-run)
and [explanations](../../language/reference/derivative.md#explanations).
