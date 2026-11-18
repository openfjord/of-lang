---
title: Constraints — The OpenFjord modeling language by example
---

# Constraints

> Part of [the language by example](README.md).

A constraint rejects facts. It finds things, and it reports each match as a
violation.

<table>
<tr>
<td>

Each service has a monthly cost and a monthly budget.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    monthlyCost:
      type: integer
    monthlyBudget:
      type: integer
```

</td>
</tr>
<tr>
<td>

A constraint has a `says` sentence that gives the rule, in plain words.

</td>
<td>

```yaml
serviceStaysInBudget:
  is: Constraint
  says: A service does not cost more than its budget.
```

</td>
</tr>
<tr>
<td>

`match` lists the things to find. Here the variable `$service` represents
each service in turn.

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

`filter` keeps the matches for which each condition holds. This condition
holds when the cost of the service is above its budget. Each match that
stays is a violation, so `match` and `filter` describe what breaks the rule.

</td>
<td>

```yaml
  filter:
    - node:  $service
      field: monthlyCost
      is above:
        same as:
          node:  $service
          field: monthlyBudget
```

</td>
</tr>
<tr>
<td>

`message` is what a violation says. `{$service}` gives the id of the
service.

</td>
<td>

```yaml
  message: "{$service} must not cost more than its budget."
```

</td>
</tr>
<tr>
<td>

A test file states two services. The cost of the search service is above
its budget.

</td>
<td>

```yaml
# tests/examples/search-over-budget.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, serviceStaysInBudget]

service/checkout:
  is: Fact
  kind: Service
  fields:
    monthlyCost: 400
    monthlyBudget: 500

service/search:
  is: Fact
  kind: Service
  fields:
    monthlyCost: 900
    monthlyBudget: 600
```

</td>
</tr>
<tr>
<td>

`rejected` lists all constraints that the facts must break. The file
imports each constraint that the test names.

</td>
<td>

```yaml
searchOverBudget:
  is: Test
  says: The search service costs more than its budget, so the constraint rejects it.
  expect:
    rejected: [serviceStaysInBudget]
```

</td>
</tr>
</table>

`run` reports each violation with the name of the constraint and its
message:

```console
$ of-lang run constraints constraints/tests/examples/search-over-budget.of.yaml
derived:
  nothing
rejected:
  serviceStaysInBudget: service/search must not cost more than its budget.
```

More about this: [constraints](../../language/reference/constraint.md) and
[how constraints run](../../language/reference/constraint.md#how-constraints-run).
