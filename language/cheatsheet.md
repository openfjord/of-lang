---
title: Cheat sheet — The OpenFjord modeling language
---

# Cheat sheet

> Part of [the OpenFjord modeling language](README.md). All of the syntax is
> on one page, in the default YAML form. [The reference](README.md#reference) gives the exact rules.

## A package

The name of every file in `defs/` and `tests/` ends in `.of.yaml` or `.of.yml`
(in the JSON form, `.of.json`; a package uses one form). The manifest is always YAML:
`package.yaml` or `package.yml`.

```
my-package/
  package.yaml       # required: the manifest
  defs/              # required: Kind, Link, ValueType, Constraint, Derivative, Doc
  docs/              # optional: Markdown that links to or embeds definitions; each file has a Doc, or other docs refer to it
  tests/             # optional: Fact, Test, Question
```

```yaml
# package.yaml
package: "@acme/platform"          # required; quoted because of the leading @
version: 0.1.0                     # required
owner: platform                    # required; owns everything in the package
rdfNamespace: platform             # required; used only in exported RDF
dependsOn:
  "@acme/org": ^1.2
exports:
  kinds: [Service]
  links: [ownedByTeam]
  valueTypes: [Name]
openFor:
  - extend: Service
  - under: usesTable
tags:                                # each label that definitions can carry
  - name: privacy
    description: What a privacy review needs to tell customers which services handle which sorts of data.
```

## A file

```yaml
import:
  - from: file:names.of.yaml                     # relative to this file
    use: [Name]
  - from: pkg:@acme/org                          # must be in dependsOn and exported there
    use: [Team, { name: Region, as: SalesRegion }]

SomeStatement:
  is: Kind        # or Link, ValueType, Constraint, Derivative, Doc — Fact, Test, Question in tests/
```

The import schemes are `file:` and `pkg:`. A later version can add more.
A file imports each name that it uses but does not define, and no other
name. A file does not import field names.

All free-form text (`definition`, `reads`, `examples`, `says`, `asks`,
`message`, docs) follows ASD-STE100 Simplified Technical English.

## Definitions

```yaml
Service:                                 # UpperCamelCase
  is: Kind
  definition: A program that runs in production and answers requests.   # required
  tags: [privacy]                        # optional; names from package.yaml; on any definition but a Doc
  examples:    [the checkout service]
  notExamples: [a script that a developer runs one time and does not keep]
  fields:
    name:           { type: Name }                       # required by default
    repository:     { type: URL }                        # where the code is
    timeoutSeconds: { type: Seconds, optional: true }
  docs: [PlatformGuide]                  # optional; Docs of this package

PlatformGuide:                           # UpperCamelCase
  is: Doc
  path: file:../docs/platform.md         # a Markdown file in docs/, relative to this file

PublicService:
  is: Kind
  extends: Service                       # optional
  definition: A service that customers call.
  fields:
    availability: { type: Percent }
    status:       { enum: [active, retired], default: active }

Resource:
  is: Kind
  definition: Something that a service needs to do its work.
  oneOf: [Service, Table]                # a union: each resource is exactly one; closed; no fields or links

callsService:                            # lowerCamelCase, two words or more
  is: Link
  from: Service
  to: Service
  reads: "{from} calls {to}"             # required, with both {from} and {to}
  definition: The first service needs an answer from the second service, directly or through other services.
  chains: true                           # optional: a→b and b→c means a→c
  inverse: calledByService               # optional: the name read the other way

writesTable:
  is: Link
  from: Service
  to: Table
  reads: "{from} writes {to}"
  definition: The service adds, changes or deletes rows of the table.
  under: usesTable                       # optional: every writesTable is a usesTable

Name:
  is: ValueType
  definition: A lowercase identifier for services, tables, data classes and regions.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"

Seconds:
  is: ValueType
  definition: A number of seconds that a service waits for an answer.
  type: integer                          # a type and its keywords
  minimum: 1
  maximum: 3600

URL:
  is: ValueType
  definition: The address of a page on the web.
  type: string
  format: uri                            # an absolute URI, such as https://…
```

By default, the RDF export holds everything. An optional list of tags limits it. Tags do not pass through
`extends` or `under`. The export does not include a reference to an untagged
definition. It also does not include a link whose `from` or `to` kind is
untagged.

Field keywords: `type` (`string` `integer` `number` `boolean` `array` `object`
or a value type), `enum` `const` `default` `pattern` `minLength` `maxLength`
`format` (`date` `date-time` `duration` `uri`) `minimum` `maximum` `items`
`properties` `optional` `oneOf` (variants: exactly one choice matches; object
choices carry a `const` part such as `variant`). The language does not have
`null`, `$ref`, `required`, `anyOf`, `allOf`, `not` or `if`.

## Match and filter

Derivatives, constraints and questions each define their own `match` and
`filter`. All three use the same syntax, but each can change apart from the
others.

```yaml
match:                                             # things and links that must exist
  - kind: { name: PublicService, as: $service }    # each variable has exactly one kind item
  - kind: { name: Service, as: $called }
  - link: ownsService                              # a link or an inverse name
    from: Team                                     # a kind alone: nothing else refers to this end
    to:   $service
  - link: callsService
    from: $service
    to:   $called                                  # callsService chains, so this also finds longer chains
filter:                                            # conditions; all must hold
  - node: $called
    does not equal: { same as: $service }
  - node: $service
    field: status
    equals: { value: active }
```

`repeated: true` on a link item follows a link that does not chain for one or
more steps.

A match can hold only a kind item. The variable then represents each thing
of the kind:

```yaml
match:
  - kind: { name: PublicService, as: $service }    # each public service in turn
filter:
  - node: $service
    field: availability
    is below: { value: 99.9 }
```

| Comparison | Takes |
|---|---|
| `equals` `does not equal` `is above` `is below` `is at least` `is at most` | `value:` or `same as:` |
| `is one of` | `values: [...]` |
| `is a` | a kind |
| `is set` `is not set` | (only a field) |

A condition on an object field can add `part: variant`, or `part: [a, b]`
for a part in a part.

| Nested | Holds when |
|---|---|
| `at least one: { match, filter }` | it matches at least once |
| `none: { match, filter }` | it never matches |
| `count: { match, filter }` + a comparison | the number of matches satisfies the comparison |
| `any of: [conditions]` | one of them holds |

## Constraints and derivatives

```yaml
publicServiceDeclaresItsData:
  is: Constraint
  says: The privacy notice must name each data class that an active public service handles.   # no $ words
  match: [ ... ]
  filter: [ ... ]
  unless: [ ... ]                        # optional; a match is a violation unless these hold
  message: "The privacy notice must tell that {$service} handles {$data}."

dataFromTable:
  is: Derivative
  says: A service handles each data class that one of its tables holds.
  match: [ ... ]
  derive:
    - { link: handlesData, from: $service, to: $data }                    # a link
    - { node: $service, becomes: SharedService }                          # a kind
    - { node: $service, field: timeoutSeconds, set to: { value: 30 } }    # a field value
```

## Tests

```yaml
# tests/examples/some-example.of.yaml — facts plus exactly one Test
service/checkout:                        # a fact's id
  is: Fact
  kind: PublicService
  fields:
    name: checkout
    repository: https://code.example.com/checkout
    availability: 99.9
  links:
    ownedByTeam: team/web
    callsService: [service/catalog, service/accounts]   # a link: an id, or a list of ids

someExample:
  is: Test
  says: What this example shows.
  expect:
    derived:    [ { link: handlesData, from: service/checkout, to: data/personal } ]
    notDerived: [ { node: service/accounts, becomes: SharedService } ]
    rejected:   [publicServiceDeclaresItsData]

# tests/terms/write.of.yaml — no facts: what the definitions derive about kinds and links
writeIsAUse:
  is: Test
  says: A service that writes a table also uses the table.
  expect:
    derived:
      - { kind: PublicService, is a: Service }        # also: never is, can have
      - { link: writesTable, implies: usesTable }     # also: chains, inverse of
    notDerived:
      - { link: usesTable, implies: writesTable }

# tests/questions/which-public-services-handle-personal-data.of.yaml
whichPublicServicesHandlePersonalData:
  is: Question
  asks: Which public services handle personal data?
  given: file:../examples/personal-data-through-a-call.of.yaml
  match: [ ... ]
  filter: [ ... ]
  expect:
    $service: [service/checkout]
```

## Pending changes

```yaml
Cache:
  is: Kind
  pending: added                         # or removed, or updated
  definition: A copy of rows that a service keeps in memory.

Service:
  is: Kind
  pending: updated
  tags: [privacy]
  definition: A program that runs in production and gives answers to other programs.   # the new version
  examples:
    - the checkout service
    - the catalog service
  notExamples:
    - a table, which holds data and answers no requests
    - a script that a developer runs one time and does not keep
  fields:
    name:
      type: Name
    repository:
      type: URL
    timeoutSeconds:
      type: Seconds
      optional: true
  was:                                   # the whole old version, without is
    tags: [privacy]
    definition: A program that runs in production and answers requests.
    examples:
      - the checkout service
      - the catalog service
    notExamples:
      - a table, which holds data and answers no requests
      - a script that a developer runs one time and does not keep
    fields:
      name:
        type: Name
      repository:
        type: URL
      timeoutSeconds:
        type: Seconds
        optional: true
```

A package has one pending change or none. Accept applies all marks in one
step, and reject undoes them. Other packages see the package without the
change.

## How it runs

The derivatives run repeatedly until no new fact follows. `under` and `chains`
also add their facts. Then the constraints check everything. Anything that is
not stated or derived is false. Every derived fact and every rejection has
its explanation.
