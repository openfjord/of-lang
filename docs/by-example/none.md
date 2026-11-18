---
title: None — The OpenFjord modeling language by example
---

# None

> Part of [the language by example](README.md).

`none` holds a nested block, and holds when the block has no matches.

<table>
<tr>
<td>

The package has services and teams.

</td>
<td>

```yaml
# defs/product.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    name:
      type: string

Team:
  is: Kind
  definition: A group of people who build and run services.
  fields:
    name:
      type: string
```

</td>
</tr>
<tr>
<td>

A team owns a service.

</td>
<td>

```yaml
ownedBy:
  is: Link
  from: Service
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that builds and runs the service.
```

</td>
</tr>
<tr>
<td>

A test file states one team and two services. The payments team owns the
checkout service. The search service states no owner.

</td>
<td>

```yaml
# tests/examples/two-services.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Team, ownedBy]

team/payments:
  is: Fact
  kind: Team
  fields:
    name: payments

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
  links:
    ownedBy: team/payments

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

A file of facts holds exactly one test.

</td>
<td>

```yaml
twoServices:
  is: Test
  says: A service with an owner and a service without an owner are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services have no owner. Its `match` finds each
service.

</td>
<td>

```yaml
# tests/questions/which-services-have-no-owner.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Team, ownedBy]

whichServicesHaveNoOwner:
  is: Question
  asks: Which services have no owner?
  given: file:../examples/two-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

The nested block runs one time for each service. It finds the teams that own
`$service`. `none` holds when there are no such teams. A nested block can
also have a `filter`.

</td>
<td>

```yaml
  filter:
    - none:
        match:
          - link: ownedBy
            from: $service
            to:   Team
```

</td>
</tr>
<tr>
<td>

The answer is the search service. A link that is not stated does not exist,
so the nested block has no matches for it.

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
$ of-lang test none
pass  test twoServices  (none/tests/examples/two-services.of.yaml)
pass  question whichServicesHaveNoOwner  (none/tests/questions/which-services-have-no-owner.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
