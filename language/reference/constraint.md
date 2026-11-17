---
title: Constraints — The OpenFjord modeling language
---

# Constraints

> Part of [the reference](../README.md#reference).

A constraint rejects facts. It finds things and the links between them, and
it uses conditions to narrow the matches. It reports each match as a
violation. A match that satisfies the conditions under `unless` is not a
violation.

## With `unless`

```yaml
dishListsItsAllergens:
  is: Constraint
  says: >
    The menu must give a warning about each allergen in an active dish.
  match:
    - kind: { name: Dish, as: $dish }
    - kind: { name: Allergen, as: $allergen }
    - link: hasAllergen
      from: $dish
      to:   $allergen
  filter:
    - node:  $dish
      field: status
      equals:
        value: active
  unless:
    - at least one:
        match:
          - link: listsAllergen
            from: $dish
            to:   $allergen
  message: "The menu must warn that {$dish} contains {$allergen}."
```

Read it as follows. A dish contains an allergen. When the dish is active, this
is a violation, unless the menu gives a warning about the allergen next to
the dish.

## Without `unless`

```yaml
noRecipeLoop:
  is: Constraint
  says: A recipe must not use itself, directly or through other recipes.
  match:
    - kind: { name: Recipe, as: $recipe }
    - link: usesRecipe
      from: $recipe
      to:   $recipe
  message: "{$recipe} must not use itself."
```

Each match is a violation. The link `usesRecipe` chains, so this constraint
also finds a recipe that uses itself through other recipes.

Four principles apply to constraints:

1. **Readability comes before brevity.** Keys and words state each part
   in full. There are no expression strings to parse.
2. **No implicit variables.** Each thing that a constraint refers to has a
   name and a kind. A constraint does not reach a thing through a dotted path
   that skips an unnamed thing.
3. **Variables start with `$`,** so a reader cannot confuse a variable with a
   kind, a link, a field or a value.
4. **The sentence is plain words.** Product reads `says`, so `says` names no
   variables. Variables are part of the logic.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Constraint` |
| `tags` | no | the [tags](tag.md) of the definition |
| `says` | yes | the rule, in plain words |
| `match` | yes | the things and the links to find |
| `filter` | no | conditions that narrow the matches |
| `unless` | no | conditions that excuse a match |
| `message` | yes | what a violation says |

## Match

`match` lists the things and the links to find. It has two sorts of item.

A kind item names a kind, and gives a variable to the things of that kind:

```yaml
match:
  - kind: { name: Dish, as: $dish }
```

The variable represents each thing of that kind in turn. A variable is `$`,
then a lowercase letter, then letters and digits: `$recipe`, `$part`. Each
variable has exactly one kind item, so its kind is written in one place.
Every other mention of the variable is the variable alone.

A link item is one link that must exist between two things:

```yaml
match:
  - kind: { name: Recipe, as: $recipe }
  - kind: { name: Ingredient, as: $ingredient }
  - link: usesIngredient
    from: $recipe
    to:   $ingredient
```

`from` and `to` are each a variable or a kind. A variable is for a thing
that the constraint refers to again. A kind alone is for a thing that nothing
else refers to, as in `from: Recipe`. A constraint has no unused variables.

`match` finds each combination of things for which all its items hold. Each
combination is one match: a row that gives one thing to each end. Two
matches that give the same things to the same ends are one match. Two
variables can represent the same thing, unless a condition prevents it. One
variable that occurs two times, as in `noRecipeLoop`, is the same thing at
each occurrence.

The inverse name of a link reads the link in the opposite direction:
`link: makesRecipe, from: $station, to: $recipe` matches exactly the same
things as `link: madeAtStation, from: $recipe, to: $station`.

`repeated: true` matches the link one or more times in a sequence. The things
between the two ends have no names. Thus a longer chain between the same two
things is still one match. Explanations list the steps.

```yaml
- link: usesRecipe
  from: $recipe
  to:   $part
  repeated: true
```

A narrower kind limits the matches. For example, with
`kind: { name: BaseRecipe, as: $part }`, the link item above matches only
recipes that use base recipes.

## Filter and unless

`filter` and `unless` both contain conditions, and both receive the matches
one at a time. `filter` selects the matches that the constraint applies to.
`unless` tells which of these matches to excuse. Each condition in a list must
hold. Thus `unless` excuses a match only when all its conditions hold.

A condition applies to a thing (`node` alone) or to one of its fields (`node`
and `field`). Then it gives a comparison word and what to compare with:

- `value:` (a literal)
- `values:` (a list, for `is one of`)
- `same as:` (a variable, or a variable and a field).

```yaml
unless:
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

A condition on an `object` field can name one of its parts with `part`,
for example `field: plating, part: garnish`. For a part in a part, `part`
is a list of the names, from the outer part to the inner part. When the
value has no such part, the condition holds only for `is not set`. A
comparison holds only when the value has the type of the other side.

A condition cannot see links directly. It reaches links through a nested
block. A nested block has its own `match` and an optional `filter`. It runs
one time for each match of the constraint:

| Word | Holds when |
|---|---|
| `at least one` | the nested block has one or more matches |
| `none` | the nested block has no matches |
| `count` | the number of its matches satisfies the comparison word next to it |
| `any of` | one or more of the listed conditions hold |

There is no general `not`. `does not equal`, `is not set` and `none` cover the
negative cases, and each condition stays easy to read alone.

## Messages

A message states the rule that the match breaks. It uses the words "must" or
"must not". It quotes variables in braces, and the match gives their values.
`{$dish}` gives the id of the thing. `{$dish: status}` gives the value of a
field.

## Rules

- **CN1.** A constraint MUST have a `says` sentence, and the sentence MUST
  NOT contain a `$` word.
- **CN2.** Every item in `match` MUST be a kind item or a link item. A kind
  item MUST have `kind` with `name` and `as`. `name` MUST be a kind that the
  file defines or imports, and `as` MUST be a variable. A link item MUST have
  `link`, `from` and `to`, and MAY have `repeated: true`. `link` MUST be a
  link, or an inverse name, that the file defines or imports.
- **CN3.** Each `from` and each `to` MUST be a variable or a kind. A
  variable there MUST have a kind item, in the same `match` or in the
  `match` around its nested block.
- **CN4.** The kind at each end of a link item MUST be the kind that the
  link joins at that end, or a kind that extends it. For a variable, this is
  the kind in its kind item.
- **CN5.** Conditions and messages MUST write a variable alone (`$recipe`), and
  every variable that they use MUST have a kind item.
- **CN6.** A variable that appears only in a nested block belongs to that
  block. The constraint MUST NOT use it outside that block, and this includes
  the message.
- **CN7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **CN8.** A comparison MUST compare values of the same type: two fields
  that name the same value type, or two fields that name no value type and
  have the same built-in type. A literal MUST be a value that the field
  allows.
- **CN9.** A field can occur in a condition, next to `node` or in `same as`.
  A field can also occur in a message. Such a field MUST be a field of that
  variable's kind, or of a kind that it extends.
- **CN10.** When a condition compares two nodes without fields, the two
  variables MUST have the same kind, or kinds where one extends the other.

- **CN12.** A constraint MUST have a `message`. The message SHOULD use the
  words "must" or "must not".
- **CN13.** A constraint MUST NOT have `derive`. Constraints add no facts.
- **CN15.** `part` MAY appear only next to `field`. Each name in `part`
  MUST be a part that the schema of the field defines. For a `oneOf`, one
  or more of the choices MUST define it. The comparison uses the type of
  that part.
- **CN16.** A constraint MUST have `match`, with at least one item. Each
  nested block MUST also have `match`, with at least one item.
- **CN17.** Each variable MUST have exactly one kind item. A kind item
  MUST NOT give a variable that another kind item gives, in the same `match`
  or in the `match` around its nested block.
- **CN18.** Each variable MUST occur two or more times in the constraint:
  in ends, kind items, conditions and the message. A variable that occurs one
  time refers to nothing, so a link end gives only the kind instead.

## How constraints run

Constraints run after each [derivative](derivative.md) finishes. They run over
stated facts and derived facts, and they add no facts. A link that is not
stated and not derived does not exist. A field with no value is not set, not unknown.

- **CN14.** The tool reports each match that is a violation one time, with
  its message and its explanation. The explanation gives the conditions under
  `unless` that failed, if any, and the facts that the match used.
