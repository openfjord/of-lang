---
title: Filter — The OpenFjord modeling language by example
---

# Filter

> Part of [the language by example](README.md).

`filter` holds conditions. It keeps the matches for which each condition
holds.

<table>
<tr>
<td>

Each service has a name and a number of replicas. A replica is one copy of
the service that runs.

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
    replicas:
      type: integer
```

</td>
</tr>
<tr>
<td>

A test file states three services with different numbers of replicas.

</td>
<td>

```yaml
# tests/examples/three-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
    replicas: 2

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    replicas: 4

service/search:
  is: Fact
  kind: Service
  fields:
    name: search
    replicas: 12
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
threeServices:
  is: Test
  says: Three services with different numbers of replicas are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks about the numbers of replicas. Its `match` finds each
service.

</td>
<td>

```yaml
# tests/questions/which-services-have-2-or-more-replicas-but-fewer-than-12.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesHave2OrMoreReplicasButFewerThan12:
  is: Question
  asks: Which services have 2 or more replicas, but fewer than 12?
  given: file:../examples/three-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

A condition names a `node` and a `field`. Then it gives a comparison word,
and `value:` with a literal to compare with. `is at least` holds when the
service has 2 or more replicas.

</td>
<td>

```yaml
  filter:
    - node: $service
      field: replicas
      is at least:
        value: 2
```

</td>
</tr>
<tr>
<td>

`is below` holds when the service has fewer than 12 replicas. The other
comparison words that take `value:` are `equals`, `does not equal`, `is
above` and `is at most`.

</td>
<td>

```yaml
    - node: $service
      field: replicas
      is below:
        value: 12
```

</td>
</tr>
<tr>
<td>

The two conditions hold for the orders service and the checkout service. The
search service has 12 replicas, so the second condition does not hold for
it.

</td>
<td>

```yaml
  expect:
    $service: [service/orders, service/checkout]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test filter
pass  test threeServices  (filter/tests/examples/three-services.of.yaml)
pass  question whichServicesHave2OrMoreReplicasButFewerThan12  (filter/tests/questions/which-services-have-2-or-more-replicas-but-fewer-than-12.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
