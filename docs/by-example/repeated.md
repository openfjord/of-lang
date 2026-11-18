---
title: Repeated — The OpenFjord modeling language by example
---

# Repeated

> Part of [the language by example](README.md).

`repeated: true` on a link item follows the link for one or more steps. It
is for a link that does not chain.

<table>
<tr>
<td>

The package has services.

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

A service can depend on another service. The link has no `chains: true`, so
a chain of the link does not imply the link.

</td>
<td>

```yaml
dependsOn:
  is: Link
  from: Service
  to: Service
  reads: "{from} depends on {to}"
  definition: The first service calls the second service to do its work.
```

</td>
</tr>
<tr>
<td>

A test file states three services in a sequence. The checkout service
depends on the orders service, and the orders service depends on the
database.

</td>
<td>

```yaml
# tests/examples/four-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, dependsOn]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
  links:
    dependsOn: service/orders

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
  links:
    dependsOn: service/database

service/database:
  is: Fact
  kind: Service
  fields:
    name: database
```

</td>
</tr>
<tr>
<td>

The search service depends on nothing.

</td>
<td>

```yaml
service/search:
  is: Fact
  kind: Service
  fields:
    name: search
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test.

</td>
<td>

```yaml
fourServices:
  is: Test
  says: Three services that depend on each other in a sequence, and a fourth service, are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services depend on the database, directly or through
other services.

</td>
<td>

```yaml
# tests/questions/which-services-depend-on-the-database.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, dependsOn]

whichServicesDependOnTheDatabase:
  is: Question
  asks: Which services depend on the database, directly or through other services?
  given: file:../examples/four-services.of.yaml
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Service, as: $other }
```

</td>
</tr>
<tr>
<td>

With `repeated: true`, the link item matches the link one or more times in a
sequence. The services between the two ends have no names.

</td>
<td>

```yaml
    - link: dependsOn
      from: $service
      to:   $other
      repeated: true
```

</td>
</tr>
<tr>
<td>

A condition keeps the matches that end at the database.

</td>
<td>

```yaml
  filter:
    - node: $other
      equals:
        value: service/database
```

</td>
</tr>
<tr>
<td>

The answer has the orders service, and also the checkout service, which is
two steps away. Without `repeated: true`, the answer has only the orders
service.

</td>
<td>

```yaml
  expect:
    $service: [service/checkout, service/orders]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test repeated
pass  test fourServices  (repeated/tests/examples/four-services.of.yaml)
pass  question whichServicesDependOnTheDatabase  (repeated/tests/questions/which-services-depend-on-the-database.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[match](../../language/reference/constraint.md#match), [what links
imply](../../language/reference/link.md#what-links-imply).
