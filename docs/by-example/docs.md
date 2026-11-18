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

The package defines the kind `DataClass`.

</td>
<td>

```yaml
# defs/data-class.of.yaml
DataClass:
  is: Kind
  definition: A sort of data that needs care, such as personal data or payment data.
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
  docs: [PersonalDataCare]
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
PersonalDataCare:
  is: Doc
  path: file:../docs/personal-data.md
```

</td>
</tr>
</table>

The file `docs/personal-data.md` is standard Markdown. It tells why each
team must know which tables hold personal data. It embeds the definition of
the kind with `![](../defs/data-class.of.yaml#DataClass)`, and it has a
Markdown link to the kind: `[DataClass](../defs/data-class.of.yaml#DataClass)`.
The part after `#` is the name of a statement. Each Markdown file in `docs/`
needs a `Doc`, unless a different doc refers to it. `check` makes sure that
the Markdown file exists and that each Markdown link and each embed points
to a statement that exists:

```console
$ of-lang check docs
ok
```

More about this: [docs](../../language/reference/package.md#docs) and
[doc definitions](../../language/reference/package.md#doc-definitions).
