---
title: Same as — The OpenFjord modeling language by example
---

# Same as

> Part of [the language by example](README.md).

`same as:` compares with a variable, or with a field of a variable, in place
of a literal.

<table>
<tr>
<td>

Each service has a name and a port.

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
    port:
      type: integer
```

</td>
</tr>
<tr>
<td>

A test file states three services. Two of them have the same port.

</td>
<td>

```yaml
# tests/examples/two-equal-ports.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    port: 8080

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
    port: 8080

service/search:
  is: Fact
  kind: Service
  fields:
    name: search
    port: 9090
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
twoEqualPorts:
  is: Test
  says: Two services with the same port and a third service are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services have the same port as another service. It has
two variables for services. `match` finds each combination of two services.

</td>
<td>

```yaml
# tests/questions/which-services-share-a-port.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesShareAPort:
  is: Question
  asks: Which services have the same port as another service?
  given: file:../examples/two-equal-ports.of.yaml
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Service, as: $other }
```

</td>
</tr>
<tr>
<td>

Two variables can represent the same thing. This condition prevents it:
`same as:` with a variable compares the two things.

</td>
<td>

```yaml
  filter:
    - node: $service
      does not equal:
        same as: $other
```

</td>
</tr>
<tr>
<td>

`same as:` with a `node` and a `field` compares with the value of that
field. The condition holds when the two ports are equal.

</td>
<td>

```yaml
    - node: $service
      field: port
      equals:
        same as:
          node: $other
          field: port
```

</td>
</tr>
<tr>
<td>

The checkout service and the orders service have the same port. No other
service has the port of the search service.

</td>
<td>

```yaml
  expect:
    $service: [service/checkout, service/orders]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test same-as
pass  test twoEqualPorts  (same-as/tests/examples/two-equal-ports.of.yaml)
pass  question whichServicesShareAPort  (same-as/tests/questions/which-services-share-a-port.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
