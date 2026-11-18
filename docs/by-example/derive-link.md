---
title: Derive a link — The OpenFjord modeling language by example
---

# Derive a link

> Part of [the language by example](README.md).

A derivative adds facts that follow from other facts. This one adds a link
between two things.

<table>
<tr>
<td>

The package has three kinds: a dish, an ingredient and an allergen.

</td>
<td>

```yaml
# defs/allergens.of.yaml
Dish:
  is: Kind
  definition: A food that the restaurant sells to guests.

Ingredient:
  is: Kind
  definition: A food that the kitchen buys and uses in dishes.

Allergen:
  is: Kind
  definition: A substance in food that can make a guest ill.
```

</td>
</tr>
<tr>
<td>

Facts state two links: a dish uses an ingredient, and an ingredient
contains an allergen.

</td>
<td>

```yaml
usesIngredient:
  is: Link
  from: Dish
  to: Ingredient
  reads: "{from} uses {to}"
  definition: The cooks put the ingredient in the dish.

containsAllergen:
  is: Link
  from: Ingredient
  to: Allergen
  reads: "{from} contains {to}"
  definition: The ingredient has the allergen in it.
```

</td>
</tr>
<tr>
<td>

No fact states the third link. A derivative adds it.

</td>
<td>

```yaml
hasAllergen:
  is: Link
  from: Dish
  to: Allergen
  reads: "{from} has {to}"
  definition: The dish has the allergen in it.
```

</td>
</tr>
<tr>
<td>

A derivative has a `says` sentence that tells what it concludes, in plain
words.

</td>
<td>

```yaml
allergenFromIngredient:
  is: Derivative
  says: A dish has each allergen that one of its ingredients contains.
```

</td>
</tr>
<tr>
<td>

`match` finds each dish that uses an ingredient that contains an allergen.
A variable, such as `$dish`, represents each thing of its kind in turn.

</td>
<td>

```yaml
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Ingredient, as: $ingredient }
    - kind: { name: Allergen, as: $allergen }
    - link: usesIngredient
      from: $dish
      to:   $ingredient
    - link: containsAllergen
      from: $ingredient
      to:   $allergen
```

</td>
</tr>
<tr>
<td>

`derive` tells what to add for each match. An item with `link`, `from`
and `to` adds a link.

</td>
<td>

```yaml
  derive:
    - link: hasAllergen
      from: $dish
      to:   $allergen
```

</td>
</tr>
<tr>
<td>

A test file states the facts. The margherita uses parmesan, and parmesan
contains milk.

</td>
<td>

```yaml
# tests/examples/margherita.of.yaml
import:
  - from: file:../../defs/allergens.of.yaml
    use: [Dish, Ingredient, Allergen, usesIngredient, containsAllergen, hasAllergen]

allergen/milk:
  is: Fact
  kind: Allergen

ingredient/parmesan:
  is: Fact
  kind: Ingredient
  links:
    containsAllergen: allergen/milk

dish/margherita:
  is: Fact
  kind: Dish
  links:
    usesIngredient: ingredient/parmesan
```

</td>
</tr>
<tr>
<td>

The test expects the link that the derivative adds. A fact under `derived`
has the same form as an item under `derive`, with ids in place of
variables.

</td>
<td>

```yaml
margherita:
  is: Test
  says: The margherita uses parmesan, which contains milk. Thus, the margherita has milk.
  expect:
    derived:
      - link: hasAllergen
        from: dish/margherita
        to:   allergen/milk
```

</td>
</tr>
</table>

`run` shows each derived fact as the `reads` sentence of its link, with the
derivative that added it:

```console
$ of-lang run derive-link derive-link/tests/examples/margherita.of.yaml
derived:
  dish/margherita has allergen/milk  (derivative allergenFromIngredient)
rejected:
  nothing
```

More about this: [derivatives](../../language/reference/derivative.md) and
[derive](../../language/reference/derivative.md#derive).
