---
title: Narrowing a value — The OpenFjord modeling language by example
---

# Narrowing a value

> Part of [the language by example](README.md).

A type allows many values. Keywords next to the type narrow the field to
the values that are correct for it.

<table>
<tr>
<td>

`pattern` is a regular expression that the text must match. The name of a
service starts with a lowercase letter. Then it has lowercase letters,
digits and hyphens.

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
      pattern: "^[a-z][a-z0-9-]*$"
```

</td>
</tr>
<tr>
<td>

`minLength` and `maxLength` set how many characters the text can have.

</td>
<td>

```yaml
    summary:
      type: string
      minLength: 10
      maxLength: 80
```

</td>
</tr>
<tr>
<td>

`minimum` and `maximum` set the smallest number and the largest number.

</td>
<td>

```yaml
    port:
      type: integer
      minimum: 1
      maximum: 65535
```

</td>
</tr>
<tr>
<td>

`format` tells what the text holds. With `uri`, it is an absolute URI
(Uniform Resource Identifier), such as an address that starts with
`https:`. The other formats are `date`, `date-time` and `duration`.

</td>
<td>

```yaml
    runbook:
      type: string
      format: uri
```

</td>
</tr>
<tr>
<td>

`const` allows exactly one value. Each service runs in production.

</td>
<td>

```yaml
    environment:
      const: production
```

</td>
</tr>
<tr>
<td>

A test file states one service. Each value satisfies the keywords of its
field.

</td>
<td>

```yaml
# tests/examples/one-service.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    summary: Takes the payment for an order.
    port: 8443
    runbook: https://runbooks.example.com/checkout
    environment: production
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so the
fields allow all five values.

</td>
<td>

```yaml
oneService:
  is: Test
  says: Each value of the checkout service is a value that its field allows.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test narrowing
pass  test oneService  (narrowing/tests/examples/one-service.of.yaml)
1 passed, 0 failed
```

More about this: [the keywords of a field](../../language/reference/kind.md#keywords)
and [the field rules](../../language/reference/kind.md#field-rules).
