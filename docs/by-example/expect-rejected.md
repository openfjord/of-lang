---
title: Expect rejected facts — The OpenFjord modeling language by example
---

# Expect rejected facts

> Part of [the language by example](README.md).

A constraint rejects facts. A test lists under `rejected` all constraints
that its facts must break.

<table>
<tr>
<td>

A service has a number of replicas: the copies of the service that run at
the same time.

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
```

</td>
</tr>
<tr>
<td>

A constraint finds things and reports each match as a violation. This one
finds each service with 0 replicas.

</td>
<td>

```yaml
serviceRuns:
  is: Constraint
  says: A service must have at least one replica.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      field: replicas
      equals:
        value: 0
  message: "{$service} must have at least one replica."
```

</td>
</tr>
<tr>
<td>

A test file imports each constraint that its test names. It states two
services, and the search service has 0 replicas.

</td>
<td>

```yaml
# tests/examples/search-does-not-run.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, serviceRuns]

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

`rejected` lists the constraints that the facts must break. The test
passes only if the facts break exactly these constraints. Thus, a test
with `rejected: []` expects that the facts break none.

</td>
<td>

```yaml
searchDoesNotRun:
  is: Test
  says: The search service has 0 replicas, so the package rejects it.
  expect:
    rejected: [serviceRuns]
```

</td>
</tr>
</table>

`test` runs the test. The test passes because the constraint rejects the search service:

```console
$ of-lang test expect-rejected
pass  test searchDoesNotRun  (expect-rejected/tests/examples/search-does-not-run.of.yaml)
1 passed, 0 failed
```

More about this: [examples](../../language/reference/test.md#examples) and [constraints](../../language/reference/constraint.md).
