---
title: A kind alone — The OpenFjord modeling language by example
---

# A kind alone

> Part of [the language by example](README.md).

An end of a link item can be a kind alone, without a variable. This is for a
thing that nothing else refers to.

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

A test file states two teams and one service. The payments team owns the
service.

</td>
<td>

```yaml
# tests/examples/one-service-two-teams.of.yaml
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

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
  links:
    ownedBy: team/payments
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
oneServiceTwoTeams:
  is: Test
  says: A team that owns a service and a team that owns no service are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which teams own a service. The teams have the variable
`$team`.

</td>
<td>

```yaml
# tests/questions/which-teams-own-a-service.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Team, ownedBy]

whichTeamsOwnAService:
  is: Question
  asks: Which teams own a service?
  given: file:../examples/one-service-two-teams.of.yaml
  match:
    - kind: { name: Team, as: $team }
```

</td>
</tr>
<tr>
<td>

The question does not refer to the service again, so the `from` end gives
the kind alone. This end only has to exist. A variable that only one link
end uses is an error.

</td>
<td>

```yaml
    - link: ownedBy
      from: Service
      to:   $team
```

</td>
</tr>
<tr>
<td>

The answer is the payments team. No service has a link to the discovery
team.

</td>
<td>

```yaml
  expect:
    $team: [team/payments]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test kind-alone
pass  test oneServiceTwoTeams  (kind-alone/tests/examples/one-service-two-teams.of.yaml)
pass  question whichTeamsOwnAService  (kind-alone/tests/questions/which-teams-own-a-service.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match).
