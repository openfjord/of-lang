---
title: Count — The OpenFjord modeling language by example
---

# Count

> Part of [the language by example](README.md).

`count` holds a nested block, and compares the number of its matches with a
number.

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

A service depends on other services.

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

A test file states four services. The checkout service depends on three
services, and the orders service depends on two.

</td>
<td>

```yaml
# tests/examples/three-and-two-dependencies.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, dependsOn]

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
  links:
    dependsOn: [service/orders, service/search, service/database]

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
  links:
    dependsOn: [service/search, service/database]
```

</td>
</tr>
<tr>
<td>

The search service and the database depend on nothing.

</td>
<td>

```yaml
service/search:
  is: Fact
  kind: Service
  fields:
    name: search

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

A file of facts holds exactly one test.

</td>
<td>

```yaml
threeAndTwoDependencies:
  is: Test
  says: Services that depend on three, on two and on no services are valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which services depend on more than two services. Its `match`
finds each service.

</td>
<td>

```yaml
# tests/questions/which-services-depend-on-more-than-two-services.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, dependsOn]

whichServicesDependOnMoreThanTwoServices:
  is: Question
  asks: Which services depend on more than two services?
  given: file:../examples/three-and-two-dependencies.of.yaml
  match:
    - kind: { name: Service, as: $service }
```

</td>
</tr>
<tr>
<td>

The nested block runs one time for each service. It has one match for each
service that `$service` depends on.

</td>
<td>

```yaml
  filter:
    - count:
        match:
          - link: dependsOn
            from: $service
            to:   Service
```

</td>
</tr>
<tr>
<td>

The comparison word goes next to `count`. It takes `value:` with a whole
number. The condition holds when the number of matches is above 2.

</td>
<td>

```yaml
      is above:
        value: 2
```

</td>
</tr>
<tr>
<td>

The answer is the checkout service. The count for the orders service is 2.

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
$ of-lang test count
pass  test threeAndTwoDependencies  (count/tests/examples/three-and-two-dependencies.of.yaml)
pass  question whichServicesDependOnMoreThanTwoServices  (count/tests/questions/which-services-depend-on-more-than-two-services.of.yaml)
2 passed, 0 failed
```

More about this: [questions](../../language/reference/test.md#questions),
[filter and
unless](../../language/reference/constraint.md#filter-and-unless).
