---
title: Open for — The OpenFjord modeling language by example
---

# Open for

> Part of [the language by example](README.md).

A different package can extend a kind of this package, or declare a link
under a link of this package, only where the manifest opens it in `openFor`.

<table>
<tr>
<td>

The package defines a kind for services and a kind for teams.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

Team:
  is: Kind
  definition: A group of people who build and run services.
```

</td>
</tr>
<tr>
<td>

It also defines a link between them.

</td>
<td>

```yaml
ownedBy:
  is: Link
  from: Service
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that builds and runs the service.
```

</td>
</tr>
<tr>
<td>

The manifest starts with the four required keys.

</td>
<td>

```yaml
# package.yaml
package: "@by-example/open-for"
version: 0.1.0
owner: platform
rdfNamespace: open-for
```

</td>
</tr>
<tr>
<td>

`exports` lists the definitions that other packages can import.

</td>
<td>

```yaml
exports:
  kinds: [Service, Team]
  links: [ownedBy]
```

</td>
</tr>
<tr>
<td>

`openFor` lists the extension points. `extend` identifies a kind: a
different package can declare a kind that extends `Service`.

</td>
<td>

```yaml
openFor:
  - extend: Service
```

</td>
</tr>
<tr>
<td>

`under` identifies a link: a different package can declare a link under
`ownedBy`. `Team` is not in the list, so no other package can extend it.

</td>
<td>

```yaml
  - under: ownedBy
```

</td>
</tr>
</table>

The package of the page [Dependencies](dependencies.md) depends on this
package and declares `PublicService`, a kind that extends `Service`. `check`
accepts the two packages:

```console
$ of-lang check open-for dependencies
ok
```

More about this: [the manifest](../../language/reference/package.md#the-manifest)
and [exports](../../language/reference/package.md#exports).
