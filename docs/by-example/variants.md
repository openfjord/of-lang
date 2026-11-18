---
title: Variants — The OpenFjord modeling language by example
---

# Variants

> Part of [the language by example](README.md).

Some values have one of a small number of shapes. `oneOf` on a field lists
the shapes, and each value matches exactly one of them.

<table>
<tr>
<td>

A service tells how many replicas it runs. Some services always run the
same number, and other services run more replicas when they get more
requests. Thus, `scaling` has two shapes.

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
    scaling:
      oneOf:
```

</td>
</tr>
<tr>
<td>

Each choice is a schema. The first choice is an object with one number of
replicas. Its `variant` part has a fixed value, which `const` gives.

</td>
<td>

```yaml
        - type: object
          properties:
            variant:
              const: fixed
            replicas:
              type: integer
```

</td>
</tr>
<tr>
<td>

The second choice is an object with the smallest and the largest number
of replicas. Its `variant` part has a different fixed value, so a reader
can see which shape a value has.

</td>
<td>

```yaml
        - type: object
          properties:
            variant:
              const: automatic
            minReplicas:
              type: integer
            maxReplicas:
              type: integer
```

</td>
</tr>
<tr>
<td>

A test file states two services. The database always runs three replicas,
and the checkout service runs from two to ten.

</td>
<td>

```yaml
# tests/examples/two-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/database:
  is: Fact
  kind: Service
  fields:
    name: database
    scaling:
      variant: fixed
      replicas: 3

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    scaling:
      variant: automatic
      minReplicas: 2
      maxReplicas: 10
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so each
value of `scaling` matches exactly one choice.

</td>
<td>

```yaml
twoServices:
  is: Test
  says: A fixed number of replicas and an automatic number of replicas are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A condition reads a part of the value with `part`. This question uses the
`variant` part to find the services with a fixed number of replicas.

</td>
<td>

```yaml
# tests/questions/which-services-have-fixed-replicas.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesHaveFixedReplicas:
  is: Question
  asks: Which services have a fixed number of replicas?
  given: file:../examples/two-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      field: scaling
      part: variant
      equals:
        value: fixed
  expect:
    $service: [service/database]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test variants
pass  test twoServices  (variants/tests/examples/two-services.of.yaml)
pass  question whichServicesHaveFixedReplicas  (variants/tests/questions/which-services-have-fixed-replicas.of.yaml)
2 passed, 0 failed
```

More about this: [variants](../../language/reference/kind.md#variants) and
[the field rules](../../language/reference/kind.md#field-rules).
