---
title: Inverse names — The OpenFjord modeling language by example
---

# Inverse names

> Part of [the language by example](README.md).

An inverse name is a name for a link in the opposite direction. It adds no
link: it is the same link, read the other way.

<table>
<tr>
<td>

The link `ownedBy` goes from a service to a team.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

Team:
  is: Kind
  definition: A group of people who build and run services.

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

`inverse` gives the name `ownsService` to the opposite direction, from a
team to a service. An inverse name follows the same rules as the name of
a link.

</td>
<td>

```yaml
  inverse: ownsService
```

</td>
</tr>
<tr>
<td>

A test file states two teams and two services. The checkout service gives
the link by its own name.

</td>
<td>

```yaml
# tests/examples/two-teams.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, Team, ownedBy, ownsService]

team/payments:
  is: Fact
  kind: Team

service/checkout:
  is: Fact
  kind: Service
  links:
    ownedBy: team/payments

service/search:
  is: Fact
  kind: Service
```

</td>
</tr>
<tr>
<td>

A fact can use the inverse name. This link on the catalog team states the
same fact as `ownedBy: team/catalog` on the search service.

</td>
<td>

```yaml
team/catalog:
  is: Fact
  kind: Team
  links:
    ownsService: service/search
```

</td>
</tr>
<tr>
<td>

The test expects that the link holds from the search service to the
catalog team, although no fact states it in that direction.

</td>
<td>

```yaml
twoTeams:
  is: Test
  says: Each of the two teams owns one service.
  expect:
    derived:
      - link: ownedBy
        from: service/search
        to:   team/catalog
```

</td>
</tr>
<tr>
<td>

A question can also use the inverse name. This one starts at the team.

</td>
<td>

```yaml
# tests/questions/which-teams-own-a-service.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, Team, ownsService]

whichTeamsOwnAService:
  is: Question
  asks: Which teams own a service?
  given: file:../examples/two-teams.of.yaml
  match:
    - kind: { name: Team, as: $team }
    - link: ownsService
      from: $team
      to:   Service
```

</td>
</tr>
<tr>
<td>

The answer has the two teams. It does not matter which name the fact
uses.

</td>
<td>

```yaml
  expect:
    $team: [team/payments, team/catalog]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test inverse-names
pass  test twoTeams  (inverse-names/tests/examples/two-teams.of.yaml)
pass  question whichTeamsOwnAService  (inverse-names/tests/questions/which-teams-own-a-service.of.yaml)
2 passed, 0 failed
```

More about this: [what links imply](../../language/reference/link.md#what-links-imply) and [the rules for facts](../../language/reference/fact.md#rules).
