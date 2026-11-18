---
title: Derived terms — The OpenFjord modeling language by example
---

# Derived terms

> Part of [the language by example](README.md).

The definitions derive information before there are any facts: what
follows for kinds and links. A test with no facts can expect this
information.

<table>
<tr>
<td>

A public service extends a service.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

PublicService:
  is: Kind
  extends: Service
  definition: A service that customers call.

Team:
  is: Kind
  definition: A group of people who build and run services.
```

</td>
</tr>
<tr>
<td>

One link has an inverse name, and one link chains.

</td>
<td>

```yaml
ownedBy:
  is: Link
  from: Service
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that builds and runs the service.
  inverse: ownsService

dependsOn:
  is: Link
  from: Service
  to: Service
  reads: "{from} depends on {to}"
  definition: The first service cannot do its work without the second service, directly or through other services.
  chains: true
```

</td>
</tr>
<tr>
<td>

One link is under another link.

</td>
<td>

```yaml
callsService:
  is: Link
  from: Service
  to: Service
  reads: "{from} calls {to}"
  definition: The first service sends requests to the second service.
  under: dependsOn
```

</td>
</tr>
<tr>
<td>

A test in `tests/terms/` has no facts. It imports each kind and each link
that it names.

</td>
<td>

```yaml
# tests/terms/public-service.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService, ownedBy, ownsService, dependsOn, callsService]

publicServiceIsAService:
  is: Test
  says: >
    A public service is a service, so a public service can have an owner. A
    service depends on each service that it calls.
  expect:
    derived:
```

</td>
</tr>
<tr>
<td>

A term about a kind starts with `kind`. `is a` comes from `extends`. `can
have` tells that a thing of the kind can have the link: a public service
gets `ownedBy` from `Service`.

</td>
<td>

```yaml
      - kind: PublicService
        is a: Service
      - kind: PublicService
        can have: ownedBy
```

</td>
</tr>
<tr>
<td>

A term about a link starts with `link`. `implies` comes from `under`,
`chains` comes from `chains`, and `inverse of` comes from the inverse
name.

</td>
<td>

```yaml
      - link: callsService
        implies: dependsOn
      - link: dependsOn
        chains: true
      - link: ownsService
        inverse of: ownedBy
```

</td>
</tr>
<tr>
<td>

`notDerived` lists terms that must not hold. A service is not always a
public service, and a service does not always call a service that it
depends on.

</td>
<td>

```yaml
    notDerived:
      - kind: Service
        is a: PublicService
      - link: dependsOn
        implies: callsService
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test derived-terms
pass  test publicServiceIsAService  (derived-terms/tests/terms/public-service.of.yaml)
1 passed, 0 failed
```

More about this: [derived terms](../../language/reference/test.md#derived-terms).
