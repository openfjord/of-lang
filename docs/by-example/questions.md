---
title: Questions — The OpenFjord modeling language by example
---

# Questions

> Part of [the language by example](README.md).

A question asks something about the facts of an example, and tells what the
answer must be.

<table>
<tr>
<td>

A team owns a service.

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

An example gives the facts: one team and two services. The search service
has no owner.

</td>
<td>

```yaml
# tests/examples/two-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, Team, ownedBy]

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

twoServices:
  is: Test
  says: The payments team owns the checkout service, and the search service has no owner.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question goes in a file in `tests/questions/`. `asks` is the question in
plain words.

</td>
<td>

```yaml
# tests/questions/which-services-have-an-owner.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, Team, ownedBy]

whichServicesHaveAnOwner:
  is: Question
  asks: Which services have an owner?
```

</td>
</tr>
<tr>
<td>

`given` names the example file with the facts to ask about.

</td>
<td>

```yaml
  given: file:../examples/two-services.of.yaml
```

</td>
</tr>
<tr>
<td>

`match` tells what to find. A kind item gives a variable to the things of a
kind: `$service` represents each service in turn.

</td>
<td>

```yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

A link item is one link that must exist. Its `to` end is a kind alone, so
any team satisfies it.

</td>
<td>

```yaml
    - link: ownedBy
      from: $service
      to:   Team
```

</td>
</tr>
<tr>
<td>

`expect` lists the ids that the variable must represent. The answer is the
checkout service only.

</td>
<td>

```yaml
  expect:
    $service: [service/checkout]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test questions
pass  test twoServices  (questions/tests/examples/two-services.of.yaml)
pass  question whichServicesHaveAnOwner  (questions/tests/questions/which-services-have-an-owner.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions).
