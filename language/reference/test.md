---
title: Tests — The OpenFjord modeling language
---

# Tests

> Part of [the reference](../README.md#reference).

A package uses tests in its `tests/` folder to prove that it is correct
(reason 16 in [the design notes](../../docs/design.md#why-ontologies-are-hard-for-software-people)). A test is an example or a
question:

- An **example** tells what a set of facts must derive, and which constraints
  the facts must break. An example with no facts tells what the definitions
  derive about kinds and links ([derived terms](#derived-terms)).
- A **question** asks something about a set of facts, and tells what the
  answer must be.

## Examples

An example is a file of facts with one `Test` statement. All the
derivatives and constraints of the package run on the facts of the
example. A question asks about the stated facts and the derived facts
together.

```yaml
# tests/examples/personal-data-through-a-call.of.yaml
import:
  - from: pkg:@acme/org
    use: [Team]
  - from: file:../../defs/team.of.yaml
    use: [dependsOnTeam]
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService, ownedByTeam, callsService, usesTable, writesTable]
  - from: file:../../defs/table.of.yaml
    use: [Table, DataClass, holdsData]
  - from: file:../../defs/privacy.of.yaml
    use: [handlesData, declaresData]

personalDataThroughACall:
  is: Test
  says: >
    The checkout service calls the accounts service, which uses a table with
    personal data. The checkout service also writes the payments table, which
    holds payment data. Thus, the checkout service handles personal data and
    payment data. The privacy notice names the two data classes, so the
    service breaks no constraint. The web team depends on the core team.
  expect:
    derived:
      - link: handlesData
        from: service/checkout
        to:   data/personal
      - link: handlesData
        from: service/checkout
        to:   data/payment
      - link: usesTable
        from: service/checkout
        to:   table/payments
      - link: dependsOnTeam
        from: team/web
        to:   team/core
    notDerived:
      - link: handlesData
        from: service/search
        to:   data/personal
      - link: dependsOnTeam
        from: team/core
        to:   team/web
    rejected: []

team/web:
  is: Fact
  kind: Team
  fields:
    name: web

team/core:
  is: Fact
  kind: Team
  fields:
    name: core

data/personal:
  is: Fact
  kind: DataClass
  fields:
    name: personal

data/payment:
  is: Fact
  kind: DataClass
  fields:
    name: payment

table/customers:
  is: Fact
  kind: Table
  fields:
    name: customers
  links:
    holdsData: data/personal

table/payments:
  is: Fact
  kind: Table
  fields:
    name: payments
  links:
    holdsData: data/payment

table/products:
  is: Fact
  kind: Table
  fields:
    name: products

service/accounts:
  is: Fact
  kind: Service
  fields:
    name: accounts
    repository: https://code.example.com/accounts
  links:
    ownedByTeam: team/core
    usesTable: table/customers

service/catalog:
  is: Fact
  kind: Service
  fields:
    name: catalog
    repository: https://code.example.com/catalog
  links:
    ownedByTeam: team/core
    usesTable: table/products

service/checkout:
  is: Fact
  kind: PublicService
  fields:
    name: checkout
    repository: https://code.example.com/checkout
    availability: 99.9
  links:
    ownedByTeam: team/web
    callsService: [service/catalog, service/accounts]
    writesTable: table/payments
    declaresData: [data/personal, data/payment]

service/search:
  is: Fact
  kind: PublicService
  fields:
    name: search
    repository: https://code.example.com/search
    availability: 99.5
  links:
    ownedByTeam: team/web
    callsService: service/catalog
```

| Key under `expect` | Value |
|---|---|
| `derived` | facts and terms that must hold after the derivatives run |
| `notDerived` | facts and terms that must not hold |
| `rejected` | all constraints that the facts must break |

A fact under `derived` or `notDerived` has the same form as an item under `derive`,
with ids in place of variables:

| Fact | Written |
|---|---|
| a link | `link`, `from`, `to` |
| a kind | `node`, `becomes` |
| a field value | `node`, `field`, `set to` with `value:` |

- **TS1.** An example file MUST hold exactly one `Test` statement, which MUST
  have `says` and `expect`.
- **TS2.** The test passes only if all of these are true:
  - every fact and term under `derived` holds
  - no fact or term under `notDerived` holds
  - the facts break exactly the constraints under `rejected`
  - every fact obeys its kind ([FC3](fact.md#rules)).

  A test with no `rejected` is the same as a test that expects no broken
  constraints.
- **TS3.** The file MUST import each link, kind and constraint that the test
  names.
- **TS15.** A term under `derived` or `notDerived` MUST have one of the
  forms in [derived terms](#derived-terms). A test MAY list only terms. A
  test in a file with no facts MUST list only terms under `expect`, and no
  `rejected`.

## Derived terms

The definitions derive information that is not a fact: what follows for
kinds and links before there are any facts. A test can expect this
information under `derived` and `notDerived`, next to facts or alone:

```yaml
# tests/terms/write.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [Service, PublicService, ownedByTeam, callsService, usesTable, writesTable]

writeIsAUse:
  is: Test
  says: >
    A service that writes a table also uses the table. A public service is
    a service, so a public service can have a team that owns it. A service
    can call a service that calls another service.
  expect:
    derived:
      - kind: PublicService
        is a: Service
      - kind: PublicService
        can have: ownedByTeam
      - link: writesTable
        implies: usesTable
      - link: callsService
        chains: true
    notDerived:
      - link: usesTable
        implies: writesTable
```

| Term | Written | Derived from |
|---|---|---|
| a kind is another kind | `kind`, `is a` | `extends`, in one or more steps, and the kinds that a [union](kind.md#unions) lists |
| two kinds never share a thing | `kind`, `never is` | a union |
| a kind can have a link | `kind`, `can have` | the `from` of the link, or the `to` for an inverse name, with `extends` and unions |
| a link implies another link | `link`, `implies` | `under`, in one or more steps |
| a link chains | `link`, `chains: true` | `chains` |
| a link is another link in the opposite direction | `link`, `inverse of` | the inverse name |

A term names kinds and links only, never ids of facts.

## Questions

A question asks something about the facts of an example, and tells what the
answer must be. It finds things and the links between them, and it uses
conditions to narrow the matches. Then it compares the things that a variable
represented with the expected list.

```yaml
# tests/questions/which-public-services-handle-personal-data.of.yaml
import:
  - from: file:../../defs/service.of.yaml
    use: [PublicService]
  - from: file:../../defs/table.of.yaml
    use: [DataClass]
  - from: file:../../defs/privacy.of.yaml
    use: [handlesData]

whichPublicServicesHandlePersonalData:
  is: Question
  asks: Which public services handle personal data?
  given: file:../examples/personal-data-through-a-call.of.yaml
  match:
    - kind: { name: PublicService, as: $service }
    - kind: { name: DataClass, as: $data }
    - link: handlesData
      from: $service
      to:   $data
  filter:
    - node: $data
      equals:
        value: data/personal
  expect:
    $service: [service/checkout]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Question` |
| `asks` | yes | the question, in plain words |
| `given` | no | `file:` and an example file with the facts to ask about; if not given, the facts in this file |
| `match` | yes | the things and the links to find |
| `filter` | no | conditions that narrow the matches |
| `expect` | yes | for each listed variable, the ids that it must represent |

Each item in `match` is a kind item or a link item. A kind item, such as
`kind: { name: PublicService, as: $service }`, gives a variable to the things of a kind:
the variable represents each thing of that kind in turn. Each variable has
exactly one kind item. A link item is one link that must exist, with
`link`, `from` and `to`. `from` and `to` are each a variable, or a kind
alone for a thing that nothing else refers to. `repeated: true` on a link
item follows the link for one or more steps. `match` finds each combination
of things for which all its items hold. Each combination is one match, which
gives one thing to each variable. An end with a kind alone only has to
exist.

`filter` keeps the matches for which each condition holds. A condition names
a `node`, an optional `field`, a comparison word and what to compare with:

| Word | Takes |
|---|---|
| `equals`, `does not equal`, `is above`, `is below`, `is at least`, `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set`, `is not set` | a `field`, and nothing more |

A condition can also contain a nested block with its own `match` and
`filter`. The words for this are `at least one`, `none`, or `count` with a
comparison word. `any of` holds when one of its listed conditions holds.

A condition can also name a `part` of an `object` field, as in
[constraints](constraint.md#filter-and-unless).

- **TS4.** A question MUST have an `asks` sentence, and the sentence MUST NOT
  contain a `$` word.
- **TS5.** Every item in `match` MUST be a kind item or a link item. A kind
  item MUST have `kind` with `name` and `as`. `name` MUST be a kind that the
  file imports, and `as` MUST be a variable. A link item MUST have `link`,
  `from` and `to`, and MAY have `repeated: true`. `link` MUST be a link, or
  an inverse name, that the file imports.
- **TS6.** Each `from` and each `to` MUST be a variable or a kind. A
  variable there MUST have a kind item, in the same `match` or in a
  `match` that encloses its nested block. The kind at each end MUST be the kind of
  the link at that end, or a kind that extends it. For a variable, this is
  the kind in its kind item.
- **TS7.** Conditions MUST write a variable alone (`$data`), and every
  variable that they use MUST have a kind item. A variable of a nested block
  MUST NOT be used outside that block.
- **TS8.** A comparison MUST compare values of the same type: two fields that
  name the same value type, or two fields that name no value type and have
  the same built-in type. A literal MUST be a value that the field allows. A
  literal compared with a node without a field MUST be the id of a fact. A
  field MUST be a field of the kind of its node, or of a kind that it
  extends. Two nodes that a condition
  compares without fields MUST have the same kind, or kinds where one
  extends the other.
- **TS9.** Each key under `expect` MUST be a variable of the question's
  `match`.
- **TS10.** The question passes only if, for each listed variable, the
  distinct things that it represents in all matches are exactly the listed
  ids.
- **TS11.** `part` MAY appear only next to `field`. Each name in `part`
  MUST be a part that the schema of the field defines. For a `oneOf`, one
  or more of the choices MUST define it. The comparison uses the type of
  that part.
- **TS12.** A question MUST have `match`, with at least one item. Each
  nested block MUST also have `match`, with at least one item.
- **TS13.** Each variable MUST have exactly one kind item. Two kind items
  in one question MUST NOT give the same variable.
- **TS14.** Each variable MUST be used outside its kind item: at two or
  more places, or in a condition or `expect`. A variable at one link end
  and nowhere else refers to nothing, so that end gives the kind alone.
- **TS16.** A condition MUST be one of the forms above. `is set` and `is
  not set` MUST name a field. `is a` MUST name a kind.
- **TS17.** `count` MUST have one comparison word: `equals`, `does not
  equal`, `is above`, `is below`, `is at least` or `is at most`. Its value
  MUST be `value:` with a whole number, 0 or more.