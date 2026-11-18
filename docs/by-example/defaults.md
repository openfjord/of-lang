---
title: Defaults — The OpenFjord modeling language by example
---

# Defaults

> Part of [the language by example](README.md).

A field with a `default` has that value when a fact does not give one.

<table>
<tr>
<td>

A field is required, unless it says something else. Each service must have
a `name`.

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
```

</td>
</tr>
<tr>
<td>

`tier` has a default. A service that gives no tier is `standard`.

</td>
<td>

```yaml
    tier:
      enum: [critical, standard]
      default: standard
```

</td>
</tr>
<tr>
<td>

To see the default at work, a test file states two services. The first one
gives no tier.

</td>
<td>

```yaml
# tests/examples/two-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/search:
  is: Fact
  kind: Service
  fields:
    name: search
```

</td>
</tr>
<tr>
<td>

The second one gives a tier of its own.

</td>
<td>

```yaml
service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    tier: critical
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test. This one expects that the package
rejects nothing: a service without a tier is valid.

</td>
<td>

```yaml
twoServices:
  is: Test
  says: A service without a tier and a critical service are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services are standard, given those facts.

</td>
<td>

```yaml
# tests/questions/which-services-are-standard.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesAreStandard:
  is: Question
  asks: Which services are standard?
  given: file:../examples/two-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      field: tier
      equals:
        value: standard
```

</td>
</tr>
<tr>
<td>

The answer is the search service, although its fact does not state a tier.

</td>
<td>

```yaml
  expect:
    $service: [service/search]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test defaults
pass  test twoServices  (defaults/tests/examples/two-services.of.yaml)
pass  question whichServicesAreStandard  (defaults/tests/questions/which-services-are-standard.of.yaml)
2 passed, 0 failed
```

More about this: [fields](../../language/reference/kind.md#fields).
