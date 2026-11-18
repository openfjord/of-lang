---
title: Hello, kind — The OpenFjord modeling language by example
---

# Hello, kind

> Part of [the language by example](README.md).

The smallest package has a manifest and one definition. This one defines a
kind: a category of things.

<table>
<tr>
<td>

The manifest gives the package its name and its version. The name starts
with a scope, here `@by-example`.

</td>
<td>

```yaml
# package.yaml
package: "@by-example/hello-kind"
version: 0.1.0
```

</td>
</tr>
<tr>
<td>

The owner owns everything in the package. The export to RDF (Resource
Description Framework) uses `rdfNamespace`, and nothing else does.

</td>
<td>

```yaml
owner: platform
rdfNamespace: hello-kind
```

</td>
</tr>
<tr>
<td>

Definitions go in files in `defs/`. A file holds statements, and the key
of each statement is its name. `is` tells what the statement is.

</td>
<td>

```yaml
# defs/service.of.yaml
Service:
  is: Kind
```

</td>
</tr>
<tr>
<td>

Every kind needs a `definition`: one sentence that tells what the kind is.

</td>
<td>

```yaml
  definition: A program that runs in production and answers requests.
```

</td>
</tr>
<tr>
<td>

`examples` and `notExamples` are optional. They show where the kind starts
and where it stops.

</td>
<td>

```yaml
  examples:
    - the checkout service
  notExamples:
    - a script that a developer runs one time on a laptop
```

</td>
</tr>
</table>

`check` reads the package and reports each rule that it breaks:

```console
$ of-lang check hello-kind
ok
```

More about this: [kinds](../../language/reference/kind.md) and
[the manifest](../../language/reference/package.md#the-manifest).
