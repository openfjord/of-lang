---
title: Extends — The OpenFjord modeling language by example
---

# Extends

> Part of [the language by example](README.md).

A kind can extend a broader kind. It has every field of the broader kind,
and each of its things is also a thing of the broader kind.

<table>
<tr>
<td>

`Service` is the broader kind. Each service has a name and a port.

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
    port:
      type: integer
```

</td>
</tr>
<tr>
<td>

`extends` names the broader kind. A public service is a narrower sort of
service.

</td>
<td>

```yaml
PublicService:
  is: Kind
  extends: Service
  definition: A service that customers call.
```

</td>
</tr>
<tr>
<td>

A public service has the fields of a service, and it adds a field of its
own. It cannot declare `name` or `port` again.

</td>
<td>

```yaml
  fields:
    domain:
      type: string
```

</td>
</tr>
<tr>
<td>

A test file states one service and one public service. The public service
gives the two fields of a service, and also its domain.

</td>
<td>

```yaml
# tests/examples/a-service-and-a-public-service.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService]

service/orders:
  is: Fact
  kind: Service
  fields:
    name: orders
    port: 8080

service/checkout:
  is: Fact
  kind: PublicService
  fields:
    name: checkout
    port: 8443
    domain: checkout.example.com
```

</td>
</tr>
<tr>
<td>

The test expects what the definitions derive about the two kinds: a
public service is a service, but a service is not a public service.

</td>
<td>

```yaml
aServiceAndAPublicService:
  is: Test
  says: A public service is a service, and it has the fields of a service.
  expect:
    derived:
      - kind: PublicService
        is a: Service
    notDerived:
      - kind: Service
        is a: PublicService
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks for each service, given those facts. The answer includes
the checkout service, because each public service is also a service.

</td>
<td>

```yaml
# tests/questions/which-services-are-there.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service]

whichServicesAreThere:
  is: Question
  asks: Which services are there?
  given: file:../examples/a-service-and-a-public-service.of.yaml
  match:
    - kind: { name: Service, as: $service }
  expect:
    $service: [service/orders, service/checkout]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test extends
pass  test aServiceAndAPublicService  (extends/tests/examples/a-service-and-a-public-service.of.yaml)
pass  question whichServicesAreThere  (extends/tests/questions/which-services-are-there.of.yaml)
2 passed, 0 failed
```

More about this: [the kind rules](../../language/reference/kind.md#kind-rules)
and [derived terms](../../language/reference/test.md#derived-terms).
