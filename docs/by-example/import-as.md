---
title: Import as — The OpenFjord modeling language by example
---

# Import as

> Part of [the language by example](README.md).

An import can give a name a different name in the file that imports it.

<table>
<tr>
<td>

This file defines the kind `Ingredient`.

</td>
<td>

```yaml
# defs/ingredient.of.yaml
Ingredient:
  is: Kind
  definition: A food that the kitchen buys and does not make.
```

</td>
</tr>
<tr>
<td>

A second file imports `Ingredient`, but calls it `Product`. An item in `use`
can be a name, or `name` and `as` together.

</td>
<td>

```yaml
# defs/buying.of.yaml
import:
  - from: file:ingredient.of.yaml
    use: [{ name: Ingredient, as: Product }]
```

</td>
</tr>
<tr>
<td>

The rest of the file uses the name `Product`. In this file, the name
`Ingredient` does not exist.

</td>
<td>

```yaml
Supplier:
  is: Kind
  definition: A company that sells food to the restaurant.

soldBy:
  is: Link
  from: Product
  to: Supplier
  reads: "{from} is sold by {to}"
  definition: The supplier that sells the ingredient to the kitchen.
```

</td>
</tr>
</table>

The kind keeps its name `Ingredient` in all other files. Two names in one
file must not be the same. Thus, `as` is necessary when a file imports the
same name from two sources. `check` accepts the package:

```console
$ of-lang check import-as
ok
```

More about this: [imports](../../language/reference/files.md#imports).
