---
title: Dependencies — The OpenFjord modeling language by example
---

# Dependencies

> Part of [the language by example](README.md).

A package can use the definitions of a different package. The manifest
names the other package in `dependsOn`, and a file imports a name from it
with `pkg:`.

<table>
<tr>
<td>

The manifest starts with the four required keys.

</td>
<td>

```yaml
# package.yaml
package: "@by-example/dependencies"
version: 0.1.0
owner: platform
rdfNamespace: dependencies
```

</td>
</tr>
<tr>
<td>

`dependsOn` gives the full name of each package that this package depends
on, and a range of versions. `^0.1` accepts each version 0.1 of
`@by-example/open-for`, the package of the page [Open for](open-for.md).

</td>
<td>

```yaml
dependsOn:
  "@by-example/open-for": ^0.1
```

</td>
</tr>
<tr>
<td>

A `pkg:` import gives the full name of the other package. The file can use
only the names that the other package exports. `Service` is one of them.

</td>
<td>

```yaml
# defs/public-service.of.yaml
import:
  - from: pkg:@by-example/open-for
    use: [Service]
```

</td>
</tr>
<tr>
<td>

The file then uses `Service` as it uses a name of its own package. The
other package lets this package extend `Service`.

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

A test file also imports the name from the other package. The test
expects that a public service is a service.

</td>
<td>

```yaml
# tests/terms/public-service-is-a-service.of.yaml
import:
  - from: pkg:@by-example/open-for
    use: [Service]
  - from: file:../../defs/public-service.of.yaml
    use: [PublicService]

publicServiceIsAService:
  is: Test
  says: A public service is a service.
  expect:
    derived:
      - kind: PublicService
        is a: Service
```

</td>
</tr>
</table>

The tool looks for each dependency in the folder that holds the package. The
folder `open-for` is next to the folder `dependencies`, so `test` finds it:

```console
$ of-lang test dependencies
pass  test publicServiceIsAService  (dependencies/tests/terms/public-service-is-a-service.of.yaml)
1 passed, 0 failed
```

More about this: [dependencies](../../language/reference/package.md#dependencies),
[exports](../../language/reference/package.md#exports) and
[imports](../../language/reference/files.md#imports).
