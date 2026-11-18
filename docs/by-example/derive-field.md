---
title: Derive a field value — The OpenFjord modeling language by example
---

# Derive a field value

> Part of [the language by example](README.md).

A derivative can set the value of a field. The fact about the thing does not
state that value.

<table>
<tr>
<td>

Each dish must state how many portions the kitchen has.

</td>
<td>

```yaml
# defs/dish.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    portionsLeft:
      type: integer
      minimum: 0
```

</td>
</tr>
<tr>
<td>

`soldOut` is optional, so a fact can omit it.

</td>
<td>

```yaml
    soldOut:
      type: boolean
      optional: true
```

</td>
</tr>
<tr>
<td>

The derivative finds each dish, and `filter` keeps the dishes with no
portions.

</td>
<td>

```yaml
soldOutFromPortions:
  is: Derivative
  says: A dish with no portions left is sold out.
  match:
    - kind: { name: Dish, as: $dish }
  filter:
    - node:  $dish
      field: portionsLeft
      equals:
        value: 0
```

</td>
</tr>
<tr>
<td>

An item with `node`, `field` and `set to` sets a field value. `value:`
gives a literal, which must be a value that the field allows.

</td>
<td>

```yaml
  derive:
    - node:  $dish
      field: soldOut
      set to:
        value: true
```

</td>
</tr>
<tr>
<td>

A test file states two dishes. Neither fact states `soldOut`.

</td>
<td>

```yaml
# tests/examples/no-fondue-left.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    portionsLeft: 12

dish/fondue:
  is: Fact
  kind: Dish
  fields:
    portionsLeft: 0
```

</td>
</tr>
<tr>
<td>

The test expects that the fondue is sold out, and that the margherita is
not.

</td>
<td>

```yaml
noFondueLeft:
  is: Test
  says: The kitchen has no fondue and twelve portions of margherita, so only the fondue is sold out.
  expect:
    derived:
      - node:  dish/fondue
        field: soldOut
        set to:
          value: true
    notDerived:
      - node:  dish/margherita
        field: soldOut
        set to:
          value: true
```

</td>
</tr>
</table>

`run` shows the field value that the derivative sets:

```console
$ of-lang run derive-field derive-field/tests/examples/no-fondue-left.of.yaml
derived:
  dish/fondue has soldOut true  (derivative soldOutFromPortions)
rejected:
  nothing
```

More about this: [derive](../../language/reference/derivative.md#derive) and
[fields](../../language/reference/kind.md#fields).
