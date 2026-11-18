---
title: Facts — The OpenFjord modeling language by example
---

# Facts

> Part of [the language by example](README.md).

A fact tells what a thing is: its kind, its fields and its links. Facts
occur only in the tests of a package, where they are the data for examples
and questions.

<table>
<tr>
<td>

The definitions tell what a fact can state. A dish has a name and a price,
and it uses ingredients.

</td>
<td>

```yaml
# defs/dish.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.
  fields:
    name:
      type: string
    price:
      type: number

Ingredient:
  is: Kind
  definition: A raw food that the restaurant buys.
  fields:
    name:
      type: string

usesIngredient:
  is: Link
  from: Dish
  to: Ingredient
  reads: "{from} uses {to}"
  definition: The dish needs the ingredient.
```

</td>
</tr>
<tr>
<td>

Facts go in files in `tests/`. A file of facts imports each kind and each
link that its facts use. It does not import field names.

</td>
<td>

```yaml
# tests/examples/one-dish.of.yaml
import:
  - from: file:../../defs/dish.of.yaml
    use: [Dish, Ingredient, usesIngredient]
```

</td>
</tr>
<tr>
<td>

The name of the statement is the id of the fact. An id is any text, and
usually a path. `kind` names the kind of the thing.

</td>
<td>

```yaml
ingredient/tomatoes:
  is: Fact
  kind: Ingredient
```

</td>
</tr>
<tr>
<td>

`fields` gives the value of each field, by the name of the field.

</td>
<td>

```yaml
  fields:
    name: tomatoes
```

</td>
</tr>
<tr>
<td>

A second ingredient has the same form.

</td>
<td>

```yaml
ingredient/mozzarella:
  is: Fact
  kind: Ingredient
  fields:
    name: mozzarella
```

</td>
</tr>
<tr>
<td>

The dish must give its two fields, because a field is required by default.

</td>
<td>

```yaml
dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12
```

</td>
</tr>
<tr>
<td>

`links` gives one id or a list of ids for each link. Each id is the id of
a fact in the same file, and that fact must be of the kind where the link
ends.

</td>
<td>

```yaml
  links:
    usesIngredient: [ingredient/tomatoes, ingredient/mozzarella]
```

</td>
</tr>
<tr>
<td>

A file of facts holds exactly one test. This one expects that the package
rejects nothing: each fact obeys its kind.

</td>
<td>

```yaml
oneDish:
  is: Test
  says: The margherita uses tomatoes and mozzarella.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` runs the test:

```console
$ of-lang test facts
pass  test oneDish  (facts/tests/examples/one-dish.of.yaml)
1 passed, 0 failed
```

More about this: [facts](../../language/reference/fact.md) and [ids](../../language/reference/files.md#ids).
