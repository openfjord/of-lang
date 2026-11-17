---
title: Facts — The OpenFjord modeling language
---

# Facts

> Part of [the reference](../README.md#reference).

A fact tells what a thing is: its kind, its fields and its links. Facts that
files contain are stated facts. They occur only in the tests of a package,
where they are the data for examples and questions. Facts that derivatives add
are [derived](#derived-facts). Facts that running systems write are deferred
(see [the overview](../../README.md#deferred-facts-at-runtime)).

```yaml
import:
  - from: file:../../defs/recipe.of.yaml
    use: [Recipe, Dish, madeAtStation, usesRecipe]

dish/margherita:
  is: Fact
  kind: Dish
  fields:
    name: margherita
    price: 12
  links:
    madeAtStation: station/pizza
    usesRecipe: [recipe/tomato_sauce, recipe/pesto]
```

The fact omits `status`, so the dish takes the default, `active`. The kind
`Dish` extends `Recipe`, so the fact can state `name` and `madeAtStation`.

| Key | Required | Value |
|---|---|---|
| the fact's id | — | the name of the statement: any text, usually a path such as `dish/margherita` |
| `is` | yes | `Fact` |
| `kind` | yes | the kind of the thing |
| `fields` | as the kind tells | the value of each field, by the name of the field |
| `links` | no | for each link, one id or a list of ids |

## Rules

- **FC1.** A fact MUST name its `kind`, and the kind MUST be imported.
- **FC2.** Each key under `fields` MUST be a field of the kind (inherited
  fields included). Each key under `links` MUST be a link that the file
  imports, and the thing MUST be able to have the link
  ([links of a thing](kind.md#links-of-a-thing)). Fields do not need an
  import, but links do.
- **FC3.** The facts about a thing, stated and derived together, MUST have
  every field that its kind requires, and no field that its kind does not
  list.

  An omitted field takes its `default`, if it has one.
- **FC4.** Each field value MUST be a value that its field allows.
- **FC5.** The value of a link MUST be the id of a fact in the same file, or a
  list of such ids. Each of these facts MUST be of the kind where the link
  ends, of a kind that extends it, or of a kind that a union there lists. A fact MAY use an inverse name under `links`: `makesRecipe: recipe/pesto` on a
  station states the same fact as `madeAtStation: station/sauce` on the
  recipe.
- **FC6.** Facts define nothing: a fact never declares a kind, a link or a
  field.

## Derived facts

A [derivative](derivative.md) adds a derived fact. The `under` or `chains`
key of a link can also add one ([what links imply](link.md#what-links-imply)).
No file declares a derived fact, because then it is not derived. A derived
fact exists only after the derivatives run, and the tool shows it in its
output ([how the tool shows derived facts](../../cli/README.md#derived-facts)).
A test can expect a derived fact, or expect that it is absent, under `expect`
([tests](test.md#examples)). Constraints and questions see stated facts and
derived facts in the same way.

Each derived fact records why it is true. The record contains the derivative
or declaration that added the fact, what each variable stands for, and the
facts that the derived fact comes from. An explanation shows that record.
