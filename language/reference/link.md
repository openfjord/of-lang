---
title: Links — The OpenFjord modeling language
---

# Links

> Part of [the reference](../README.md#reference).

A link is a named relationship from a kind to a kind, in one direction. A link
is declared one time, as its own statement. Usually, it is in the same file as
the kind where it starts. A thing can have each link that starts at its kind
([links of a thing](kind.md#links-of-a-thing)).

```yaml
madeAtStation:
  is: Link
  from: Recipe
  to: Station
  reads: "{from} is made at {to}"
  definition: The station where cooks make the recipe.
  inverse: makesRecipe

usesGarnish:
  is: Link
  from: Recipe
  to: Ingredient
  reads: "{from} is garnished with {to}"
  definition: The cook puts the ingredient on the food at the end, to make it look or taste better.
  under: usesIngredient              # every usesGarnish is also a usesIngredient

usesRecipe:
  is: Link
  from: Recipe
  to: Recipe
  reads: "{from} uses {to}"
  definition: The first recipe needs the second recipe, directly or through other recipes.
  chains: true                       # a uses b and b uses c means a uses c
  inverse: usedByRecipe
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Link` |
| `from` | yes | the kind where the link starts |
| `to` | yes | the kind where the link ends |
| `reads` | yes | the sentence that product sees, with `{from}` and `{to}` |
| `definition` | yes | one sentence that tells what the link means |
| `under` | no | the broader link; this link is a narrower sort of it |
| `inverse` | no | a name for the link in the opposite direction |
| `chains` | no | `true` when a chain of the link implies the link |
| `tags` | no | the [tags](tag.md) that the link belongs to |

A link has only its two ends. It has no fields. Any data about a link becomes
a thing in the middle, for example a `Portion` between a recipe and an ingredient
that it uses.

## Rules

- **LK1.** A link's name, and its inverse name, MUST be lowerCamelCase, with
  two or more words that tell what the link connects and how: `usesRecipe`,
  never `uses` (reason 2 in [the research](../../README.md)).
- **LK2.** A link MUST have a `definition` and a `reads` sentence. `reads`
  MUST contain both `{from}` and `{to}` (reason 7).
- **LK3.** `from` and `to` MUST each name a kind that the file defines or
  imports. A name that is not a kind, or that names nothing, is an error.
  `from` and `to` are checks: the language
  rejects a fact that links a station with `usesRecipe`. A link never gives a
  thing a kind.
- **LK4.** `under` MUST name a declared link, not an inverse name. The `from`
  of the narrower link MUST be the `from` of the broader link, or a kind that
  extends it. The same rule applies to `to`. A link MUST NOT be under itself,
  directly or through other links.
- **LK5.** A link with `chains: true` MUST connect a kind to the same kind, or
  to a kind that extends it.

## What links imply

These are the only inferences that the language makes without a derivative.
Each fact that they add is a [derived fact](fact.md#derived-facts), and
records its reason.

| Declared | Adds |
|---|---|
| `inverse: makesRecipe` on `madeAtStation` | nothing new: `makesRecipe` is every `madeAtStation` in the opposite direction, in all derivatives, constraints, questions and facts |
| `under: usesIngredient` on `usesGarnish` | a `usesIngredient` for every `usesGarnish` |
| `chains: true` on `usesRecipe` | `a usesRecipe c` each time that `a usesRecipe b` and `b usesRecipe c` |
