---
title: Derive a field value — The OpenFjord modeling language by example
---

# Derive a field value

> Part of [the language by example](README.md).

A derivative can set the value of a field. The fact about the thing does not
state that value.

<table>
<tr>
<td>

Each service must state how many replicas it has. A replica is one running
copy of the service.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    replicas:
      type: integer
      minimum: 0
```

</td>
</tr>
<tr>
<td>

`available` is optional, so a fact can omit it.

</td>
<td>

```yaml
    available:
      type: boolean
      optional: true
```

</td>
</tr>
<tr>
<td>

The derivative finds each service, and `filter` keeps the services with no
replicas.

</td>
<td>

```yaml
notAvailableFromReplicas:
  is: Derivative
  says: A service with no replicas is not available.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node:  $service
      field: replicas
      equals:
        value: 0
```

</td>
</tr>
<tr>
<td>

An item with `node`, `field` and `set to` sets a field value. `value:`
gives a literal, which must be a value that the field allows.

</td>
<td>

```yaml
  derive:
    - node:  $service
      field: available
      set to:
        value: false
```

</td>
</tr>
<tr>
<td>

A test file states two services. Neither fact states `available`.

</td>
<td>

```yaml
# tests/examples/search-is-down.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/checkout:
  is: Fact
  kind: Service
  fields:
    replicas: 3

service/search:
  is: Fact
  kind: Service
  fields:
    replicas: 0
```

</td>
</tr>
<tr>
<td>

The test expects that the search service is not available. It also expects
that the derivative does not set the field for the checkout service.

</td>
<td>

```yaml
searchIsDown:
  is: Test
  says: The search service has no replicas and the checkout service has three, so only the search service is not available.
  expect:
    derived:
      - node:  service/search
        field: available
        set to:
          value: false
    notDerived:
      - node:  service/checkout
        field: available
        set to:
          value: false
```

</td>
</tr>
</table>

`run` shows the field value that the derivative sets:

```console
$ of-lang run derive-field derive-field/tests/examples/search-is-down.of.yaml
derived:
  service/search has available false  (derivative notAvailableFromReplicas)
rejected:
  nothing
```

More about this: [derive](../../language/reference/derivative.md#derive) and
[fields](../../language/reference/kind.md#fields).
