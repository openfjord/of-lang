---
title: Derivatives — The OpenFjord modeling language
---

# Derivatives

> Part of [the reference](README.md).

A derivative adds facts that follow from other facts. It finds things and the
links between them, and it uses conditions to narrow the matches. For each
match, it adds a link, a kind or a field value.

```yaml
allergenFromIngredient:
  is: Derivative
  says: A recipe contains each allergen that one of its ingredients contains.
  match:
    - link: usesIngredient
      from: {$recipe: Recipe}
      to:   {$ingredient: Ingredient}
    - link: containsAllergen
      from: {$ingredient: Ingredient}
      to:   {$allergen: Allergen}
  derive:
    - link: hasAllergen
      from: $recipe
      to:   $allergen
```

Read it as follows: for each recipe `$recipe` that uses an ingredient
`$ingredient` that contains an allergen `$allergen`, add "`$recipe` contains
`$allergen`".

Four principles apply to derivatives:

1. **Readability comes before brevity.** Keys and words state each part
   in full. There are no expression strings to parse.
2. **No implicit variables.** Each thing that a derivative refers to has a
   name and a kind. A derivative does not reach a thing through a dotted path
   that skips an unnamed thing.
3. **Variables start with `$`,** so a reader cannot confuse a variable with a
   kind, a link, a field or a value.
4. **The sentence is plain words.** Product reads `says`, so `says` names no
   variables. Variables are part of the logic.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Derivative` |
| `says` | yes | what the derivative concludes, in plain words |
| `match` | yes, unless [DV11](#rules) makes it optional | the links to find |
| `filter` | no | conditions that narrow the matches |
| `derive` | yes | what to add for each match |

## Match

Each item in `match` is one link that must exist between two things:

```yaml
match:
  - link: usesIngredient
    from: {$recipe: Recipe}
    to:   {$ingredient: Ingredient}
```

`from` and `to` each give a variable and the kind that the variable must
have. A variable is `$`, then a lowercase letter, then letters and digits:
`$recipe`, `$part`.

`match` finds each combination of things for which all its links hold. Each
combination is one match: a row that gives one thing to each variable. Two
matches that give the same things to the same variables are one match. Two
variables can represent the same thing, unless a condition prevents it.

The inverse name of a link reads the link in the opposite direction: `link: makesRecipe, from:
{$station: Station}, to: {$recipe: Recipe}` matches exactly the same things as `link: madeAtStation,
from: {$recipe: Recipe}, to: {$station: Station}`.

`repeated: true` matches the link one or more times in a sequence. The things
between the two ends have no names. Thus a longer chain between the same two
things is still one match. Explanations list the steps.

```yaml
- link: usesRecipe
  from: {$recipe: Recipe}
  to:   {$part: Recipe}
  repeated: true
```

A narrower kind limits the matches. For example, `{$part: BaseRecipe}`
on `usesRecipe` matches only recipes that use base recipes.

## Filter

`filter` receives the matches one at a time. It keeps the matches for which
each condition holds.

A condition applies to a thing (`node` alone) or to one of its fields (`node`
and `field`). Then it gives a comparison word and what to compare with:

- `value:` (a literal)
- `values:` (a list, for `is one of`)
- `same as:` (a variable, or a variable and a field).

```yaml
filter:
  - node: $other
    does not equal:
      same as: $station
  - node:  $dish
    field: prepMinutes
    is at least:
      same as:
        node:  $part
        field: prepMinutes
  - node:  $dish
    field: status
    is one of:
      values: [active, retired]
```

| Word | Holds when |
|---|---|
| `equals`, `does not equal` | the two are equal, or are not equal |
| `is above`, `is below` | the first is greater, or smaller |
| `is at least`, `is at most` | the first is greater or equal, or smaller or equal |
| `is one of` | the first equals one of `values` |
| `is a` | the thing has the named kind (`is a: BaseRecipe`) |
| `is set`, `is not set` | the field has a value, or has no value |

When a condition compares a thing without a field, it compares the id of the
thing. A field with no value fails all comparisons except `is not set`.

A condition cannot see links directly. It reaches links through a nested
block. A nested block has its own `match` and an optional `filter`. It runs
one time for each match of the derivative:

| Word | Holds when |
|---|---|
| `at least one` | the nested block has one or more matches |
| `none` | the nested block has no matches |
| `count` | the number of its matches satisfies the comparison word next to it |
| `any of` | one or more of the listed conditions hold |

There is no general `not`. `does not equal`, `is not set` and `none` cover the
negative cases, and each condition stays easy to read alone.

## Derive

| Adds | Written | Example |
|---|---|---|
| a link | `link`, `from`, `to` | `link: dependsOnStation, from: $station, to: $other` |
| a kind | `node`, `becomes` | `node: $recipe, becomes: BaseRecipe` |
| a field value | `node`, `field`, `set to` with `value:` or `same as:` | `node: $dish, field: status, set to: { value: retired }` |

```yaml
stationDependsOnStation:
  is: Derivative
  says: >
    A station depends on each different station that makes a recipe that one
    of its own recipes uses.
  match:
    - link: makesRecipe
      from: {$station: Station}
      to:   {$recipe: Recipe}
    - link: usesRecipe
      from: {$recipe: Recipe}
      to:   {$part: Recipe}
    - link: madeAtStation
      from: {$part: Recipe}
      to:   {$other: Station}
  filter:
    - node: $other
      does not equal:
        same as: $station
  derive:
    - link: dependsOnStation
      from: $station
      to:   $other
```

```yaml
baseRecipe:
  is: Derivative
  says: >
    A recipe that more than two other recipes use, directly or through other
    recipes, is a base recipe.
  filter:
    - count:
        match:
          - link: usesRecipe
            from: {$user: Recipe}
            to:   {$recipe: Recipe}
      is above:
        value: 2
  derive:
    - node: $recipe
      becomes: BaseRecipe
```

`baseRecipe` has no `match` of its own. Its only variable outside the nested
block is `$recipe`, and `$recipe` represents each recipe in turn
([DV11](#rules)).

A derivative can match a link that a derivative adds. For example,
`allergenFromRecipe` matches `hasAllergen`, and `allergenFromIngredient` adds
`hasAllergen`. When a derived link chains, a derivative adds only the direct
step. `chains: true` on the link adds the longer chains.

```yaml
allergenFromRecipe:
  is: Derivative
  says: A recipe contains each allergen that one of its recipes contains.
  match:
    - link: usesRecipe
      from: {$recipe: Recipe}
      to:   {$part: Recipe}
    - link: hasAllergen
      from: {$part: Recipe}
      to:   {$allergen: Allergen}
  derive:
    - link: hasAllergen
      from: $recipe
      to:   $allergen
```

## Rules

- **DV1.** A derivative MUST have a `says` sentence, and the sentence MUST
  NOT contain a `$` word.
- **DV2.** Every item in `match` MUST have `link`, `from` and `to`, and MAY
  have `repeated: true`. `link` MUST be a link, or an inverse name, that the
  file defines or imports.
- **DV3.** Each `from` and each `to` MUST use the form `{$variable: Kind}`,
  every time and also in nested blocks. Every mention of a variable MUST give
  the same kind.
- **DV4.** The kind MUST be the kind that the link joins at that end, or a
  kind that extends it.
- **DV5.** Conditions and `derive` MUST write a variable without its kind
  (`$recipe`), and every variable that they use MUST appear in a `match` item.
- **DV6.** A variable that appears only in a nested block belongs to that
  block. The derivative MUST NOT use it outside that block.
- **DV7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **DV8.** A comparison MUST compare values of the same type: two fields
  that name the same value type, or two fields that name no value type and
  have the same built-in type. A literal MUST be a value that the field
  allows.
- **DV9.** A field that a condition names, next to `node` or in `same as`,
  MUST be a field of that node's kind, or of a kind that it extends.
- **DV10.** When a condition compares two nodes without fields, the two
  variables MUST have the same kind, or kinds where one extends the other.
- **DV11.** A derivative MAY have no `match` when each variable that it uses
  appears in a nested block. Then each such variable represents, in turn,
  every thing of its kind.
- **DV12.** A derivative MUST have `derive`, with at least one item.
- **DV13.** A derived link MUST join variables whose kinds the link allows
  ([LK3](link.md#rules)). A derived kind MUST be a kind. A derived field MUST
  be a field of the kind of the thing, and its value MUST be a value that the
  field allows.
- **DV14.** A derivative MUST NOT create things. Each variable in `derive`
  MUST be a variable that the derivative binds outside all nested blocks.
- **DV15.** A field MUST have one value at the end. Before anything runs, the
  tool rejects two derivatives that can give different values to the same
  field. The matches of one derivative can still give two values to a field
  of a thing. Then the run stops with an error that shows the two
  explanations.
- **DV16.** A derivative MUST NOT have `unless` or `message`. Derivatives
  reject nothing.

## How derivatives run

1. **Derivatives run over the stated facts,** and each derivative adds what
   it derives for each match. Then they run again over the result, until a
   pass adds nothing. They never create things, merge things or remove facts,
   so this process always ends.
2. **`none`, `count` and `is not set` over derived facts wait.** A derivative
   can use them on facts that a different derivative adds. Then it runs only
   after that other derivative finishes.
3. **Closed world.** A link that is not stated and not derived does not
   exist. A field with no value is not set, not unknown.

- **DV17.** Derivatives MUST NOT depend on each other in a circle through
  `none`, `count` or `is not set`. Such derivatives can never become stable,
  so the tool rejects them.

[Constraints](constraint.md) run after all derivatives finish.

## Explanations

Each derived fact has an explanation. The explanation gives the derivative
that added the fact, the thing that each variable represented, and each fact
that the match used. A fact that is also derived has its own explanation, and
the explanation includes it. This continues down to the facts that a file
states.
