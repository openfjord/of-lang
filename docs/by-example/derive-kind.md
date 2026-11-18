---
title: Derive a kind — The OpenFjord modeling language by example
---

# Derive a kind

> Part of [the language by example](README.md).

A derivative can give a thing a narrower kind. The thing keeps its own kind
and also becomes a thing of the narrower kind.

<table>
<tr>
<td>

A service can depend on other services.

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
  definition: The first service needs the second service to do its work.
```

</td>
</tr>
<tr>
<td>

A shared service is a narrower sort of service: `SharedService` extends
`Service`. No fact states that a service is a shared service.

</td>
<td>

```yaml
SharedService:
  is: Kind
  extends: Service
  definition: A service that other services depend on.
```

</td>
</tr>
<tr>
<td>

The derivative finds each service that some service depends on. The end
`from: Service` gives a kind alone, because nothing else refers to the
service at that end.

</td>
<td>

```yaml
sharedServiceFromUse:
  is: Derivative
  says: A service that another service depends on is a shared service.
  match:
    - kind: { name: Service, as: $service }
    - link: dependsOn
      from: Service
      to:   $service
```

</td>
</tr>
<tr>
<td>

An item with `node` and `becomes` adds a kind to the thing. The kind must
extend the kind of the variable.

</td>
<td>

```yaml
  derive:
    - node: $service
      becomes: SharedService
```

</td>
</tr>
<tr>
<td>

A test file states two services. The checkout service depends on the
database.

</td>
<td>

```yaml
# tests/examples/database.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, SharedService, dependsOn]

service/database:
  is: Fact
  kind: Service

service/checkout:
  is: Fact
  kind: Service
  links:
    dependsOn: service/database
```

</td>
</tr>
<tr>
<td>

The test expects that the database becomes a shared service. No service
depends on the checkout service, so `notDerived` expects that the checkout
service does not.

</td>
<td>

```yaml
database:
  is: Test
  says: The checkout service depends on the database, so only the database is a shared service.
  expect:
    derived:
      - node: service/database
        becomes: SharedService
    notDerived:
      - node: service/checkout
        becomes: SharedService
```

</td>
</tr>
</table>

`run` shows the kind that the derivative adds:

```console
$ of-lang run derive-kind derive-kind/tests/examples/database.of.yaml
derived:
  service/database is a SharedService  (derivative sharedServiceFromUse)
rejected:
  nothing
```

More about this: [derive](../../language/reference/derivative.md#derive) and
[kinds](../../language/reference/kind.md).
