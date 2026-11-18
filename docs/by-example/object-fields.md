---
title: Objects and arrays — The OpenFjord modeling language by example
---

# Objects and arrays

> Part of [the language by example](README.md).

A field can hold more than one value. An `object` has named parts, and an
`array` is a list of values.

<table>
<tr>
<td>

Each service has a `name`, which is one value.

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
```

</td>
</tr>
<tr>
<td>

`properties` lists the parts of an `object`, by name. Each part has a
schema, like a field. A part is required, unless it has `optional: true`.

</td>
<td>

```yaml
    onCall:
      type: object
      properties:
        team:
          type: string
        phone:
          type: string
          optional: true
```

</td>
</tr>
<tr>
<td>

`items` gives the schema of each value of an `array`. Here, each value is
a whole number.

</td>
<td>

```yaml
    ports:
      type: array
      items:
        type: integer
```

</td>
</tr>
<tr>
<td>

A test file states two services. The first one gives the two parts of
`onCall`, and a list of two ports.

</td>
<td>

```yaml
# tests/examples/two-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    onCall:
      team: payments
      phone: "+31 20 555 0100"
    ports: [8443, 9090]
```

</td>
</tr>
<tr>
<td>

The second one omits the optional part, and its list has one port.

</td>
<td>

```yaml
service/search:
  is: Fact
  kind: Service
  fields:
    name: search
    onCall:
      team: discovery
    ports: [8080]

twoServices:
  is: Test
  says: A service with a phone number for the team on call and a service without one are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A condition names one part of an `object` field with `part`. This question
asks which services have the payments team on call.

</td>
<td>

```yaml
# tests/questions/which-services-have-payments-on-call.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesHavePaymentsOnCall:
  is: Question
  asks: Which services have the payments team on call?
  given: file:../examples/two-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
  filter:
    - node: $service
      field: onCall
      part: team
      equals:
        value: payments
  expect:
    $service: [service/checkout]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test object-fields
pass  test twoServices  (object-fields/tests/examples/two-services.of.yaml)
pass  question whichServicesHavePaymentsOnCall  (object-fields/tests/questions/which-services-have-payments-on-call.of.yaml)
2 passed, 0 failed
```

More about this: [the keywords of a field](../../language/reference/kind.md#keywords)
and [conditions on a part](../../language/reference/constraint.md#filter-and-unless).
