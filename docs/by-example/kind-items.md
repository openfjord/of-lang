---
title: Kind items — The OpenFjord modeling language by example
---

# Kind items

> Part of [the language by example](README.md).

A kind item in `match` gives a variable to the things of a kind.

<table>
<tr>
<td>

The package has two kinds. The first kind is `Service`.

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
```

</td>
</tr>
<tr>
<td>

The second kind is `Team`.

</td>
<td>

```yaml
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

A test file states two services and one team.

</td>
<td>

```yaml
# tests/examples/two-services-and-a-team.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service, Team]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout

service/search:
  is: Fact
  kind: Service
  fields:
    name: search

team/payments:
  is: Fact
  kind: Team
  fields:
    name: payments
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test. This one expects that the package
rejects nothing.

</td>
<td>

```yaml
twoServicesAndATeam:
  is: Test
  says: Two services and a team are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services there are, given those facts.

</td>
<td>

```yaml
# tests/questions/which-services-are-there.of.yaml
import:
  - from: file:../../defs/product.of.yaml
    use: [Service]

whichServicesAreThere:
  is: Question
  asks: Which services are there?
  given: file:../examples/two-services-and-a-team.of.yaml
```

</td>
</tr>
<tr>
<td>

The kind item names the kind `Service` and gives it the variable `$service`.
The variable represents each service in turn. Each variable has exactly one
kind item.

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

The answer is the two services. The team is not a service, so `$service`
never represents it.

</td>
<td>

```yaml
  expect:
    $service: [service/checkout, service/search]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test kind-items
pass  test twoServicesAndATeam  (kind-items/tests/examples/two-services-and-a-team.of.yaml)
pass  question whichServicesAreThere  (kind-items/tests/questions/which-services-are-there.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match).
