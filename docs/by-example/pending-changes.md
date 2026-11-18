---
title: Pending changes — The OpenFjord modeling language by example
---

# Pending changes

> Part of [the language by example](README.md).

A pending change is a change to a package that nobody has accepted yet. The
change is in the files of the package, as marks on the statements that it
adds, removes or updates.

<table>
<tr>
<td>

`pending: added` marks a statement that the change adds. The package
without the change does not include it.

</td>
<td>

```yaml
# defs/menu.of.yaml
Drink:
  is: Kind
  pending: added
  definition: A liquid that the restaurant sells to guests.
```

</td>
</tr>
<tr>
<td>

`pending: removed` marks a statement that the change removes. The statement
stays in the file, as written, until somebody accepts the change.

</td>
<td>

```yaml
Snack:
  is: Kind
  pending: removed
  definition: A small food that the restaurant sells at the bar.
```

</td>
</tr>
<tr>
<td>

`pending: updated` marks a statement that the change updates. The statement
shows its version with the change: here, a dish has a `course`.

</td>
<td>

```yaml
Dish:
  is: Kind
  pending: updated
  definition: A food that the restaurant sells to guests.
  fields:
    name:
      type: string
    course:
      enum: [starter, main, dessert]
      default: main
```

</td>
</tr>
<tr>
<td>

`was` keeps the whole statement as it was before the change, without `is`.
Only an updated statement has `was`.

</td>
<td>

```yaml
  was:
    definition: A food that the restaurant sells to guests.
    fields:
      name:
        type: string
```

</td>
</tr>
<tr>
<td>

The tests use the package with the change. Thus, a test file can state a
drink, and a dish with a course. It cannot state a snack.

</td>
<td>

```yaml
# tests/examples/menu-with-the-change.of.yaml
import:
  - from: file:../../defs/menu.of.yaml
    use: [Drink, Dish]

drink/lemonade:
  is: Fact
  kind: Drink

dish/tiramisu:
  is: Fact
  kind: Dish
  fields:
    name: tiramisu
    course: dessert
```

</td>
</tr>
<tr>
<td>

The test expects that the package with the change rejects nothing.

</td>
<td>

```yaml
menuWithTheChange:
  is: Test
  says: A drink and a dish with a course are valid.
  expect:
    rejected: []
```

</td>
</tr>
</table>

All the marks in a package together are its pending change: a package has
one pending change or none. To accept the change is to apply all the marks
in one step, and to reject it is to undo them all. `test` runs the package
with the change:

```console
$ of-lang test pending-changes
pass  test menuWithTheChange  (pending-changes/tests/examples/menu-with-the-change.of.yaml)
1 passed, 0 failed
```

More about this: [pending changes](../../language/reference/pending.md).
