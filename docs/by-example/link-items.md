---
title: Link items — The OpenFjord modeling language by example
---

# Link items

> Part of [the language by example](README.md).

A link item in `match` is one link that must exist between two things.

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

The link `ownedBy` goes from a service to a team.

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

A test file states two teams.

</td>
<td>

```yaml
# tests/examples/services-and-owners.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Team, ownedBy]

team/payments:
  is: Fact
  kind: Team
  fields:
    name: payments

team/discovery:
  is: Fact
  kind: Team
  fields:
    name: discovery
```

</td>
</tr>
<tr>
<td>

The payments team owns the checkout service. The search service has no
owner.

</td>
<td>

```yaml
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
servicesAndOwners:
  is: Test
  says: A service with an owner and a service without an owner are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which teams own which services. Each variable has a kind
item.

</td>
<td>

```yaml
# tests/questions/which-teams-own-which-services.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Team, ownedBy]

whichTeamsOwnWhichServices:
  is: Question
  asks: Which teams own which services?
  given: file:../examples/services-and-owners.of.yaml
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Team, as: $team }
```

</td>
</tr>
<tr>
<td>

The link item has `link`, `from` and `to`. `match` finds each combination of
a service and a team for which the `ownedBy` link exists.

</td>
<td>

```yaml
    - link: ownedBy
      from: $service
      to:   $team
```

</td>
</tr>
<tr>
<td>

There is one match: the checkout service and the payments team. The search
service and the discovery team have no such link, so they are not in the
answer.

</td>
<td>

```yaml
  expect:
    $service: [service/checkout]
    $team: [team/payments]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test link-items
pass  test servicesAndOwners  (link-items/tests/examples/services-and-owners.of.yaml)
pass  question whichTeamsOwnWhichServices  (link-items/tests/questions/which-teams-own-which-services.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match).
