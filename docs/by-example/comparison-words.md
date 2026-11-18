---
title: Comparison words — The OpenFjord modeling language by example
---

# Comparison words

> Part of [the language by example](README.md).

Some comparison words do not take `value:` or `same as:`. `is one of` takes
a list, `is a` takes a kind, and `is set` and `is not set` take only a
field.

<table>
<tr>
<td>

A service has a stage. It can have `onCall`, the team to call when the
service fails.

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
    stage:
      enum: [live, beta, retired]
    onCall:
      type: string
      optional: true
```

</td>
</tr>
<tr>
<td>

A public service is a service.

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

A test file states two services that are not public. The reports service
gives no `onCall`.

</td>
<td>

```yaml
# tests/examples/four-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService]

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
    stage: live
    onCall: payments

service/reports:
  is: Fact
  kind: Service
  fields:
    name: reports
    stage: retired
```

</td>
</tr>
<tr>
<td>

It also states two public services. The search service gives no `onCall`.

</td>
<td>

```yaml
service/checkout:
  is: Fact
  kind: PublicService
  fields:
    name: checkout
    stage: beta
    onCall: payments

service/search:
  is: Fact
  kind: PublicService
  fields:
    name: search
    stage: live
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
fourServices:
  is: Test
  says: Two services and two public services are valid, with or without a team to call.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

`is one of` takes `values:` with a list. It holds when the field equals one
of the values. Only the reports service has a different stage.

</td>
<td>

```yaml
# tests/questions/which-services-are-live-or-in-beta.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesAreLiveOrInBeta:
  is: Question
  asks: Which services are live or in beta?
  given: file:../examples/four-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      field: stage
      is one of:
        values: [live, beta]
  expect:
    $service: [service/orders, service/checkout, service/search]
```

</td>
</tr>
<tr>
<td>

`is a` takes a kind. It applies to the thing, so the condition has no
`field`. It holds when the thing has that kind.

</td>
<td>

```yaml
# tests/questions/which-services-are-public.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService]

whichServicesArePublic:
  is: Question
  asks: Which services are public services?
  given: file:../examples/four-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      is a: PublicService
  expect:
    $service: [service/checkout, service/search]
```

</td>
</tr>
<tr>
<td>

`is not set` holds when the field has no value. `is set` holds when the
field has a value. The two words compare with nothing.

</td>
<td>

```yaml
# tests/questions/which-services-have-no-team-to-call.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesHaveNoTeamToCall:
  is: Question
  asks: Which services do not give a team to call?
  given: file:../examples/four-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      field: onCall
      is not set: true
  expect:
    $service: [service/reports, service/search]
```

</td>
</tr>
</table>

`test` runs the test and the three questions:

```console
$ of-lang test comparison-words
pass  test fourServices  (comparison-words/tests/examples/four-services.of.yaml)
pass  question whichServicesAreLiveOrInBeta  (comparison-words/tests/questions/which-services-are-live-or-in-beta.of.yaml)
pass  question whichServicesArePublic  (comparison-words/tests/questions/which-services-are-public.of.yaml)
pass  question whichServicesHaveNoTeamToCall  (comparison-words/tests/questions/which-services-have-no-team-to-call.of.yaml)
4 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
