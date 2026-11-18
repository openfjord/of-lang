---
title: Fields — The OpenFjord modeling language by example
---

# Fields

> Part of [the language by example](README.md).

A field is a value of the thing itself, for example a name or a port. A
kind lists its fields, and each field has a type.

<table>
<tr>
<td>

`fields` lists the fields of the kind, by name. `type` gives the type of
the value. A `string` is text.

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

A `number` can have decimals, and an `integer` is a whole number. A
`boolean` is `true` or `false`.

</td>
<td>

```yaml
    monthlyCost:
      type: number
    replicas:
      type: integer
    public:
      type: boolean
```

</td>
</tr>
<tr>
<td>

A field is required, unless it says something else. With `optional: true`,
a fact can omit the field.

</td>
<td>

```yaml
    onCall:
      type: string
      optional: true
```

</td>
</tr>
<tr>
<td>

A test file states two services. The first one gives a value for each
field.

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
    monthlyCost: 1250.5
    replicas: 3
    public: true
    onCall: payments
```

</td>
</tr>
<tr>
<td>

The second one omits `onCall`, which is the only optional field.

</td>
<td>

```yaml
service/search:
  is: Fact
  kind: Service
  fields:
    name: search
    monthlyCost: 400
    replicas: 2
    public: false
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so each
service has every required field, and each value has the type of its
field.

</td>
<td>

```yaml
twoServices:
  is: Test
  says: A service with a team on call and a service without one are both valid.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test fields
pass  test twoServices  (fields/tests/examples/two-services.of.yaml)
1 passed, 0 failed
```

More about this: [fields](../../language/reference/kind.md#fields) and
[the rules for facts](../../language/reference/fact.md#rules).
