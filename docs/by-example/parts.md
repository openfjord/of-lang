---
title: Parts — The OpenFjord modeling language by example
---

# Parts

> Part of [the language by example](README.md).

A condition on an `object` field names one part of the value with `part`.

<table>
<tr>
<td>

The field `runtime` is an object with two parts. Each service must have a
`language`. The `version` is optional.

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
    runtime:
      type: object
      properties:
        language: { type: string }
        version:  { type: string, optional: true }
```

</td>
</tr>
<tr>
<td>

A test file states a service that is written in Go.

</td>
<td>

```yaml
# tests/examples/three-runtimes.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    runtime:
      language: go
      version: "1.24"
```

</td>
</tr>
<tr>
<td>

The orders service is written in Go too, and the search service in Python.

</td>
<td>

```yaml
service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
    runtime:
      language: go

service/search:
  is: Fact
  kind: Service
  fields:
    name: search
    runtime:
      language: python
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threeRuntimes:
  is: Test
  says: Services in Go and in Python are valid, with or without a version.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services are written in Go.

</td>
<td>

```yaml
# tests/questions/which-services-are-written-in-go.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesAreWrittenInGo:
  is: Question
  asks: Which services are written in Go?
  given: file:../examples/three-runtimes.of.yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

`part` goes next to `field`. The condition compares the `language` part of
`runtime`, not the whole object. For a part in a part, `part` is a list of
the names.

</td>
<td>

```yaml
  filter:
    - node: $service
      field: runtime
      part: language
      equals:
        value: go
```

</td>
</tr>
<tr>
<td>

The answer is the two services that are written in Go.

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
$ of-lang test parts
pass  test threeRuntimes  (parts/tests/examples/three-runtimes.of.yaml)
pass  question whichServicesAreWrittenInGo  (parts/tests/questions/which-services-are-written-in-go.of.yaml)
2 passed, 0 failed
```

More about this: [filter and
unless](../../language/reference/constraint.md#filter-and-unless),
[fields](../../language/reference/kind.md#fields).
