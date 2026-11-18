---
title: Tags — The OpenFjord modeling language by example
---

# Tags

> Part of [the language by example](README.md).

A tag is a label that marks the definitions that belong together for one
purpose. An export to RDF (Resource Description Framework) can include only
the definitions with given tags.

<table>
<tr>
<td>

The manifest starts with the four required keys.

</td>
<td>

```yaml
# package.yaml
package: "@by-example/tags"
version: 0.1.0
owner: platform
rdfNamespace: tags
```

</td>
</tr>
<tr>
<td>

`tags` lists each tag that the definitions of the package can carry. Each
tag has a `name` and a `description`.

</td>
<td>

```yaml
tags:
  - name: privacy
    description: What a privacy review needs to see which tables hold which sorts of data.
```

</td>
</tr>
<tr>
<td>

A definition lists its tags in `tags`. Each tag must be in the manifest, so
a misspelled tag is an error.

</td>
<td>

```yaml
# defs/table.of.yaml
Table:
  is: Kind
  tags: [privacy]
  definition: A set of rows in a database that all have the same columns.
```

</td>
</tr>
<tr>
<td>

Each definition lists its own tags. A tag does not pass from one definition
to a different definition.

</td>
<td>

```yaml
DataClass:
  is: Kind
  tags: [privacy]
  definition: A sort of data that needs care, such as personal data or payment data.
```

</td>
</tr>
<tr>
<td>

A link can also have tags. Each definition that a tagged definition refers
to must have the same tag. Here, these are the `from` kind and the `to` kind.

</td>
<td>

```yaml
holdsData:
  is: Link
  tags: [privacy]
  from: Table
  to: DataClass
  reads: "{from} holds {to}"
  definition: Some columns of the table contain data of this sort.
```

</td>
</tr>
<tr>
<td>

`tags` is optional. `Team` has no tag, so an export of the tag `privacy`
does not include it.

</td>
<td>

```yaml
Team:
  is: Kind
  definition: A group of people who build and run services.
```

</td>
</tr>
</table>

Tags only select what is exported. They do not change how the package runs.
`check` makes sure that each tag is in the manifest and that the tagged
definitions refer only to definitions with the same tag:

```console
$ of-lang check tags
ok
```

More about this: [tags](../../language/reference/tag.md) and
[the manifest](../../language/reference/package.md#the-manifest).
