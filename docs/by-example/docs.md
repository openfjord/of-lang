---
title: Docs — The OpenFjord modeling language by example
---

# Docs

> Part of [the language by example](README.md).

Markdown files in the `docs/` folder of a package explain its definitions. A
`Doc` gives such a file a name, and a kind lists the docs that tell more
about it.

<table>
<tr>
<td>

The package defines the kind `Allergen`.

</td>
<td>

```yaml
# defs/allergen.of.yaml
Allergen:
  is: Kind
  definition: A substance in food that can cause a dangerous reaction in some persons.
```

</td>
</tr>
<tr>
<td>

`docs` lists the docs that tell more about the kind. Each name is a `Doc` of
this package.

</td>
<td>

```yaml
  docs: [AllergenSafety]
```

</td>
</tr>
<tr>
<td>

A `Doc` is a definition with one key, `path`: `file:`, then the path of a
Markdown file in `docs/`, relative to this file.

</td>
<td>

```yaml
AllergenSafety:
  is: Doc
  path: file:../docs/allergens.md
```

</td>
</tr>
</table>

The file `docs/allergens.md` is standard Markdown. It tells why the kitchen
keeps a list of allergens. It embeds the definition of the kind with
`![](../defs/allergen.of.yaml#Allergen)`, and it has a Markdown link to the
kind: `[Allergen](../defs/allergen.of.yaml#Allergen)`. The part after `#` is
the name of a statement. `check` makes sure that the Markdown file exists and
that each Markdown link and each embed points to a statement that exists:

```console
$ of-lang check docs
ok
```

More about this: [docs](../../language/reference/package.md#docs) and
[doc definitions](../../language/reference/package.md#doc-definitions).
