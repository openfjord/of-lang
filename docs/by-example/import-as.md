---
title: Import as — The OpenFjord modeling language by example
---

# Import as

> Part of [the language by example](README.md).

An import can give a name a different name in the file that imports it.

<table>
<tr>
<td>

This file defines the kind `Table`.

</td>
<td>

```yaml
# defs/table.of.yaml
Table:
  is: Kind
  definition: A set of rows in a database that all have the same columns.
```

</td>
</tr>
<tr>
<td>

A second file imports `Table`, but calls it `DatabaseTable`. An item in
`use` can be a name, or `name` and `as` together.

</td>
<td>

```yaml
# defs/service.of.yaml
import:
  - from: file:table.of.yaml
    use: [{ name: Table, as: DatabaseTable }]
```

</td>
</tr>
<tr>
<td>

The rest of the file uses the name `DatabaseTable`. In this file, the name
`Table` does not exist.

</td>
<td>

```yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

readsTable:
  is: Link
  from: Service
  to: DatabaseTable
  reads: "{from} reads {to}"
  definition: The service gets rows from the table.
```

</td>
</tr>
</table>

The kind keeps its name `Table` in all other files. Two names in one file
must not be the same. Thus, `as` is necessary when a file imports the same
name from two sources. `check` accepts the package:

```console
$ of-lang check import-as
ok
```

More about this: [imports](../../language/reference/files.md#imports).
