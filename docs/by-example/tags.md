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
owner: kitchen
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
  - name: allergens
    description: What an allergen app needs to tell guests which ingredients contain which allergens.
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
# defs/ingredient.of.yaml
Ingredient:
  is: Kind
  tags: [allergens]
  definition: A food that the kitchen buys and does not make.
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
Allergen:
  is: Kind
  tags: [allergens]
  definition: A substance in food that can cause a dangerous reaction in some persons.
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
containsAllergen:
  is: Link
  tags: [allergens]
  from: Ingredient
  to: Allergen
  reads: "{from} contains {to}"
  definition: The allergen is a part of the ingredient.
```

</td>
</tr>
<tr>
<td>

`tags` is optional. `Station` has no tag, so an export of the tag
`allergens` does not include it.

</td>
<td>

```yaml
Station:
  is: Kind
  definition: A part of the kitchen where a group of cooks makes some of the recipes.
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
