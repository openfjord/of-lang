---
title: Derivatives — The OpenFjord modeling language
---

# Derivatives

> Part of [the reference](../README.md#reference).

A derivative adds facts that follow from other facts. It finds things and the
links between them, and it uses conditions to narrow the matches. For each
match, it adds a link, a kind or a field value.

```yaml
# defs/privacy.of.yaml
dataFromTable:
  is: Derivative
  tags: [privacy]
  says: A service handles each data class that one of its tables holds.
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Table, as: $table }
    - kind: { name: DataClass, as: $data }
    - link: readsTable
      from: $service
      to:   $table
    - link: holdsData
      from: $table
      to:   $data
  derive:
    - link: handlesData
      from: $service
      to:   $data
```

Read it as follows: for each service `$service` that reads a table
`$table` that holds a data class `$data`, add "`$service` handles
`$data`".

Four principles apply to derivatives:

1. **Readability comes before brevity.** Keys and words state each part
   in full. There are no expression strings to parse.
2. **No implicit variables.** Each variable has a kind item. A derivative
   never reaches a thing through a dotted path. A thing that the derivative
   refers to again has a variable.
3. **Variables start with `$`,** so a reader cannot confuse a variable with a
   kind, a link, a field or a value.
4. **The sentence is plain words.** Product people read `says`, so `says` names no
   variables. Variables are part of the logic.

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Derivative` |
| `tags` | no | the [tags](tag.md) of the definition |
| `says` | yes | what the derivative concludes, in plain words |
| `match` | yes | the things and the links to find |
| `filter` | no | conditions that narrow the matches |
| `derive` | yes | what to add for each match |

## Match

`match` lists the things and the links to find. It has two sorts of item.

A kind item names a kind, and gives a variable to the things of that kind,
as in `kind: { name: Service, as: $service }` in `dataFromTable`. The
variable represents each thing of that kind in turn. A variable is `$`,
then a lowercase letter, then letters and digits: `$service`, `$called`. Each
variable has exactly one kind item, so its kind is written in one place.
Every other mention of the variable is the variable alone.

A link item is one link that must exist between two things, such as the
`readsTable` link from `$service` to `$table` in
`dataFromTable`.

`from` and `to` are each a variable or a kind. A variable is for a thing
that the derivative refers to again. A kind alone is for a thing that nothing
else refers to, as in `from: Service`. A derivative has no unused variables.

`match` finds each combination of things for which all its items hold. Each
combination is one match: a row that gives one thing to each variable. Two
matches that give the same things to the same variables are one match. An
end with a kind alone only has to exist. In a nested block, a row also gives
one thing to each end with a kind alone. Thus `count` counts the different
things there, for example the services that depend on a service in `sharedServiceFromCalls`. Two
variables can represent the same thing, unless a condition prevents it.

The inverse name of a link reads the link in the opposite direction:
`link: ownsService, from: $team, to: $service` matches exactly the same
things as `link: ownedBy, from: $service, to: $team`.

`repeated: true` matches the link one or more times in a sequence. The things
between the two ends have no names. Thus a longer chain between the same two
things is still one match. Explanations list the steps. The constraint
`noTeamLoop` in [Constraints](constraint.md#match) matches
`dependsOnTeam` with `repeated: true`.

`repeated` is for a link that does not chain. A link with `chains: true`
already holds for each longer chain.

A narrower kind limits the matches. For example, with
`kind: { name: SharedService, as: $called }`, a `dependsOn` link item from
`$service` to `$called` matches only services that call shared services.

## Filter

`filter` receives the matches one at a time. It keeps the matches for which
each condition holds.

A condition applies to a thing (`node` alone) or to one of its fields (`node`
and `field`). Then it gives a comparison word and what to compare with:

- `value:` (a literal)
- `values:` (a list, for `is one of`)
- `same as:` (a variable, or a variable and a field).

```yaml
# defs/team.of.yaml
teamDependsOnTeam:
  is: Derivative
  says: >
    A team depends on each different team that owns a service that one of
    its own services depends on.
  match:
    - kind: { name: Team, as: $team }
    - kind: { name: Service, as: $service }
    - kind: { name: Service, as: $called }
    - kind: { name: Team, as: $other }
    - link: ownsService
      from: $team
      to:   $service
    - link: dependsOn
      from: $service
      to:   $called
    - link: ownedBy
      from: $called
      to:   $other
  filter:
    - node: $other
      does not equal:
        same as: $team
  derive:
    - link: dependsOnTeam
      from: $team
      to:   $other
```

Read it as follows. A team owns a service that depends on a second service, and
another team owns the second service. When the two teams are different, the
first team depends on the second.

| Word | Holds when |
|---|---|
| `equals`, `does not equal` | the two are equal, or are not equal |
| `is above`, `is below` | the first is greater, or smaller |
| `is at least`, `is at most` | the first is greater or equal, or smaller or equal |
| `is one of` | the first equals one of `values` |
| `is a` | the thing has the named kind (`is a: SharedService`) |
| `is set`, `is not set` | the field has a value, or has no value |

When a condition compares a thing without a field, it compares the id of the
thing. A field with no value fails all comparisons except `is not set`.

A condition on an `object` field can name one of its parts with `part`,
for example `field: schedule, part: timezone`. For a part in a part, `part`
is a list of the names, from the outer part to the inner part. When the
value has no such part, the condition holds only for `is not set`. A
comparison holds only when the value has the type of the other side.

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
| a link | `link`, `from`, `to` | `link: dependsOnTeam, from: $team, to: $other` |
| a kind | `node`, `becomes` | `node: $service, becomes: SharedService` |
| a field value | `node`, `field`, `set to` with `value:` or `same as:` | `node: $service, field: status, set to: { value: retired }` |

```yaml
# defs/service.of.yaml
sharedServiceFromCalls:
  is: Derivative
  says: >
    A service that more than two services depend on, directly or through other
    services, is a shared service.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - count:
        match:
          - link: dependsOn
            from: Service
            to:   $service
      is above:
        value: 2
  derive:
    - node: $service
      becomes: SharedService
```

`$service` represents each service in turn, and the nested block counts the
services that depend on it. The end `from: Service` has no variable, because nothing
else refers to the service that depends on `$service`.

A derivative can match a link that a derivative adds. For example,
`dataFromService` matches `handlesData`, and `dataFromTable` adds
`handlesData`. When a derived link chains, a derivative adds only the direct
step. `chains: true` on the link adds the longer chains.

```yaml
# defs/privacy.of.yaml
dataFromService:
  is: Derivative
  tags: [privacy]
  says: A service that depends on a second service handles each data class that the second service handles.
  match:
    - kind: { name: Service, as: $service }
    - kind: { name: Service, as: $called }
    - kind: { name: DataClass, as: $data }
    - link: dependsOn
      from: $service
      to:   $called
    - link: handlesData
      from: $called
      to:   $data
  derive:
    - link: handlesData
      from: $service
      to:   $data
```

## Rules

- **DV1.** A derivative MUST have a `says` sentence, and the sentence MUST
  NOT contain a `$` word.
- **DV2.** Every item in `match` MUST be a kind item or a link item. A kind
  item MUST have `kind` with `name` and `as`. `name` MUST be a kind that the
  file defines or imports, and `as` MUST be a variable. A link item MUST have
  `link`, `from` and `to`, and MAY have `repeated: true`. `link` MUST be a
  link, or an inverse name, that the file defines or imports.
- **DV3.** Each `from` and each `to` MUST be a variable or a kind. A
  variable there MUST have a kind item, in the same `match` or in a
  `match` that encloses its nested block.
- **DV4.** The kind at each end of a link item MUST be the kind that the
  link joins at that end, or a kind that extends it. For a variable, this is
  the kind in its kind item.
- **DV5.** Conditions and `derive` MUST write a variable alone (`$service`), and
  every variable that they use MUST have a kind item.
- **DV6.** A variable belongs to the `match` that has its kind item: the
  derivative's own `match`, or the `match` of a nested block. A variable of a
  nested block MUST NOT be used outside that block.
- **DV7.** A condition MUST be one of the forms above. `is set` and `is not
  set` MUST name a field. `is a` MUST name a kind.
- **DV8.** A comparison MUST compare values of the same type: two fields
  that name the same value type, or two fields that name no value type and
  have the same built-in type. A literal MUST be a value that the field
  allows. A literal that a condition compares with a node without a field
  MUST be the id of a fact.
- **DV9.** A field that a condition names, next to `node` or in `same as`,
  MUST be a field of that node's kind, or of a kind that it extends.
- **DV10.** When a condition compares two nodes without fields, the two
  variables MUST have the same kind, or kinds where one extends the other.

- **DV12.** A derivative MUST have `derive`, with at least one item.
- **DV13.** A derived link MUST join variables whose kinds the link allows
  ([LK3](link.md#rules)). A derived kind MUST be a kind that extends the kind of the variable,
  directly or through other kinds. A derived field MUST
  be a field of the kind of the thing, and its value MUST be a value that the
  field allows.
- **DV14.** A derivative MUST NOT create things. Each variable in `derive`
  MUST be a variable whose kind item is in the derivative's own `match`, not
  in a nested block.
- **DV15.** A field MUST have one value at the end. Before anything runs, the
  tool rejects two derivatives that can give different values to the same
  field. The matches of one derivative can still give two values to a field
  of a thing. Then the run stops with an error that shows the two
  explanations.
- **DV16.** A derivative MUST NOT have `unless` or `message`. Derivatives
  reject nothing.
- **DV18.** `part` MAY appear only next to `field`. Each name in `part`
  MUST be a part that the schema of the field defines. For a `oneOf`, one
  or more of the choices MUST define it. The comparison uses the type of
  that part.
- **DV19.** A derivative MUST have `match`, with at least one item. Each
  nested block MUST also have `match`, with at least one item.
- **DV20.** Each variable MUST have exactly one kind item. Two kind items
  in one derivative MUST NOT give the same variable.
- **DV21.** Each variable MUST be used outside its kind item: at two or
  more places, or in a condition or `derive`. A variable at one link end and nowhere else
  refers to nothing, so that end gives the kind alone.
- **DV22.** `count` MUST have one comparison word: `equals`, `does not
  equal`, `is above`, `is below`, `is at least` or `is at most`. Its value
  MUST be `value:` with a whole number, 0 or more.

## How derivatives run

1. **Derivatives run over the stated facts,** and each derivative adds what
   it derives for each match. Then they run again over the result, until a
   pass adds nothing. They never create things, merge things or remove facts,
   so this process always ends.
2. **`none`, `count` and `is not set` over derived facts wait.** A derivative
   can use them on facts that a different derivative adds. Then it runs only
   after that other derivative finishes.
3. **What is not stated is false.** A link that is not stated and not
   derived does not exist. A field with no value is not set, not unknown.

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
