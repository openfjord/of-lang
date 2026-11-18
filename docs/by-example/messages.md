---
title: Messages — The OpenFjord modeling language by example
---

# Messages

> Part of [the language by example](README.md).

The message of a constraint is what a violation says. It quotes variables in
braces, and each match gives their values.

<table>
<tr>
<td>

A service states how many replicas it has, and the smallest number that it
needs. A replica is one running copy of the service. The person on call is
optional.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    replicas:
      type: integer
    minReplicas:
      type: integer
    onCall:
      type: string
      optional: true
```

</td>
</tr>
<tr>
<td>

The constraint matches each service that has too few replicas.

</td>
<td>

```yaml
serviceHasEnoughReplicas:
  is: Constraint
  says: Each service has the number of replicas that it needs, or more.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node:  $service
      field: replicas
      is below:
        same as:
          node:  $service
          field: minReplicas
```

</td>
</tr>
<tr>
<td>

The message states the rule that the match breaks, with the word "must".
`{$service}` gives the id of the thing. `{$service: replicas}` gives the
value of a field.

</td>
<td>

```yaml
  message: >-
    {$service} must have {$service: minReplicas} replicas or more, but it
    has {$service: replicas}. On call: {$service: onCall}.
```

</td>
</tr>
<tr>
<td>

A test file states two services that have too few replicas. The first one
states the person on call.

</td>
<td>

```yaml
# tests/examples/too-few-replicas.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, serviceHasEnoughReplicas]

service/checkout:
  is: Fact
  kind: Service
  fields:
    replicas: 2
    minReplicas: 3
    onCall: dana
```

</td>
</tr>
<tr>
<td>

The second one states no person on call.

</td>
<td>

```yaml
service/search:
  is: Fact
  kind: Service
  fields:
    replicas: 1
    minReplicas: 2
```

</td>
</tr>
<tr>
<td>

The test expects that the facts break the constraint.

</td>
<td>

```yaml
tooFewReplicas:
  is: Test
  says: The two services have fewer replicas than they need, so the constraint rejects them.
  expect:
    rejected: [serviceHasEnoughReplicas]
```

</td>
</tr>
</table>

`run` prints one line for each violation: the name of the constraint, then
the message with the values of that match. A field with no value gives the
words "not set":

```console
$ of-lang run messages messages/tests/examples/too-few-replicas.of.yaml
derived:
  nothing
rejected:
  serviceHasEnoughReplicas: service/checkout must have 3 replicas or more, but it has 2. On call: dana.
  serviceHasEnoughReplicas: service/search must have 2 replicas or more, but it has 1. On call: not set.
```

More about this: [messages](../../language/reference/constraint.md#messages).
