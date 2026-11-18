---
title: Links — The OpenFjord modeling language by example
---

# Links

> Part of [the language by example](README.md).

A link is a named relationship from a kind to a kind, in one direction.

<table>
<tr>
<td>

A link connects two kinds. This package has services and teams.

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
```

</td>
</tr>
<tr>
<td>

A link is a statement of its own. Its name is lowerCamelCase and has two
words or more. `from` is the kind where the link starts. `to` is the kind
where the link ends.

</td>
<td>

```yaml
ownedBy:
  is: Link
  from: Service
  to: Team
```

</td>
</tr>
<tr>
<td>

`reads` is the sentence that product people see. It must contain `{from}`
and `{to}`.

</td>
<td>

```yaml
  reads: "{from} is owned by {to}"
```

</td>
</tr>
<tr>
<td>

Every link needs a `definition`: one sentence that tells what the link
means.

</td>
<td>

```yaml
  definition: The team that builds and runs the service.
```

</td>
</tr>
<tr>
<td>

To see the link at work, a test file states a team and a service. The
service gives the link under `links`, with the id of the team.

</td>
<td>

```yaml
# tests/examples/checkout-has-an-owner.of.yaml
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
```

</td>
</tr>
<tr>
<td>

The test expects that the package rejects nothing: the link starts at a
service and ends at a team.

</td>
<td>

```yaml
checkoutHasAnOwner:
  is: Test
  says: The payments team owns the checkout service.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test links
pass  test checkoutHasAnOwner  (links/tests/examples/checkout-has-an-owner.of.yaml)
1 passed, 0 failed
```

More about this: [links](../../language/reference/link.md) and [the rules for links](../../language/reference/link.md#rules).
