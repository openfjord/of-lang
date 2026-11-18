---
title: Any of — The OpenFjord modeling language by example
---

# Any of

> Part of [the language by example](README.md).

Each condition in `filter` must hold. `any of` lists conditions, and holds
when one or more of them hold.

<table>
<tr>
<td>

Each service has a tier.

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
    tier:
      enum: [critical, standard]
```

</td>
</tr>
<tr>
<td>

A public service is a service that customers call.

</td>
<td>

```yaml
PublicService:
  is: Kind
  extends: Service
  definition: A service that customers call.
```

</td>
</tr>
<tr>
<td>

A test file states three services. The database is critical, and the search
service is public.

</td>
<td>

```yaml
# tests/examples/three-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService]

service/database:
  is: Fact
  kind: Service
  fields:
    name: database
    tier: critical

service/search:
  is: Fact
  kind: PublicService
  fields:
    name: search
    tier: standard

service/reports:
  is: Fact
  kind: Service
  fields:
    name: reports
    tier: standard
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
  says: A critical service, a public service and a standard service are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services are critical or public. Its `match` finds
each service.

</td>
<td>

```yaml
# tests/questions/which-services-are-critical-or-public.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService]

whichServicesAreCriticalOrPublic:
  is: Question
  asks: Which services are critical or public?
  given: file:../examples/three-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

`any of` is one condition of the `filter`. Its first condition holds for a
critical service.

</td>
<td>

```yaml
  filter:
    - any of:
        - node: $service
          field: tier
          equals:
            value: critical
```

</td>
</tr>
<tr>
<td>

Its second condition holds for a public service.

</td>
<td>

```yaml
        - node: $service
          is a: PublicService
```

</td>
</tr>
<tr>
<td>

One condition holds for the database, and the other holds for the search
service. No condition holds for the reports service.

</td>
<td>

```yaml
  expect:
    $service: [service/database, service/search]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test any-of
pass  test threeServices  (any-of/tests/examples/three-services.of.yaml)
pass  question whichServicesAreCriticalOrPublic  (any-of/tests/questions/which-services-are-critical-or-public.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
