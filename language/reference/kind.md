---
title: Kinds, fields and value types — The OpenFjord modeling language
---

# Kinds, fields and value types

> Part of [the reference](../README.md#reference).

A kind is a category of things, for example a service, a table or a
data class. A kind tells what the thing is and gives examples of it. It also
lists the fields of the thing.

```yaml
# defs/service.of.yaml
Service:
  is: Kind
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

PublicService:
  is: Kind
  tags: [privacy]
  extends: Service
  definition: A service that customers call.
  examples:
    - the checkout service
  notExamples:
    - the catalog service, which only other services call
  fields:
    availability:
      type: Percent
    status:
      enum: [active, retired]
      default: active
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Kind` |
| `definition` | yes | one sentence that tells what the kind is |
| `extends` | no | the broader kind; this kind is a narrower sort of it |
| `oneOf` | no | the kinds that together make up this kind, with no overlap ([unions](#unions)) |
| `examples` | no | a list of things of this kind, in words |
| `notExamples` | no | a list of things that look like things of this kind, but are not |
| `fields` | no | the [fields](#fields) of the kind, by name |
| `tags` | no | the [tags](tag.md) that the kind belongs to |
| `docs` | no | the [docs](package.md#docs) that tell more about the kind |

## Kind rules

- **KD1.** A kind MUST have a `definition`.
- **KD2.** `extends` MUST name a kind. A kind MUST NOT extend itself,
  directly or through other kinds.
- **KD3.** When kind B extends kind A:
  - B has every field of A.
  - B MAY add fields of its own.
  - Each thing of kind B is also a thing of kind A.

  For example, `PublicService` extends `Service`. A public service has the
  fields and links of a service, and also an availability. Each public
  service is also a service.
- **KD13.** A kind MUST NOT declare a field with the name of a field of a
  kind that it extends, directly or through other kinds. A rule that is
  stricter for the narrower kind goes in a [constraint](constraint.md), for
  example "each public service must have a timeout".
- **KD14.** Each name in `docs` MUST be a `Doc` of the same package, which
  the file defines or imports.

## Links of a thing

Each [link](link.md) names the kind where it starts and the kind where it
ends. A thing can have each link that starts at its kind, at a kind that its
kind extends, or at a union that lists its kind or a kind that it extends.

A [constraint](constraint.md) with `count` sets how many links a thing has:

```yaml
# defs/service.of.yaml
serviceHasOneOwner:
  is: Constraint
  says: Exactly one team owns each service.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - count:
        match:
          - link: ownedByTeam
            from: $service
            to:   Team
      does not equal:
        value: 1
  message: "{$service} must have exactly one team that owns it."
```

A service without a team is an error, and so is a service with two
teams. A link with no such constraint can occur any number of times,
or not at all.

Every field is required, unless it has `optional: true` or a `default`
([FD3](#field-rules)).

## Unions

Sometimes a kind is exactly one of a fixed list of kinds, and nothing else.
For example, the thing that a service needs is a service or a table. Each
resource is one of the two, and never both. `oneOf` states this:

```yaml
# defs/service.of.yaml
Resource:
  is: Kind
  definition: Something that a service needs to do its work.
  oneOf: [Service, Table]

needsResource:
  is: Link
  from: Service
  to: Resource
  reads: "{from} needs {to}"
  definition: The service cannot do its work without the resource.
```

`Service` and `Table` do not extend `Resource`, and they do not get anything
from it. `Resource` gives a name to the two kinds together. Thus, a link, a
match or a condition can use `Resource` where a service or a table can
occur.

Each thing of a listed kind is also a thing of the union kind. Each rule
that accepts a narrower kind in place of a kind also accepts a listed kind
in place of its union kind. For example, `needsResource` can link a
service to a table.

An alert shows a different use. It is one message that tells people about a
problem in production. An alert goes to one target: a pager or a chat
channel.

```yaml
# defs/alerts.of.yaml
Pager:
  is: Kind
  definition: A device that wakes the person who is on call.

ChatChannel:
  is: Kind
  definition: A place in the chat tool where a team reads messages.

AlertTarget:
  is: Kind
  definition: The place that gets an alert.
  oneOf: [Pager, ChatChannel]

Alert:
  is: Kind
  definition: One message that tells people about a problem in production.

sentToTarget:
  is: Link
  from: Alert
  to: AlertTarget
  reads: "{from} is sent to {to}"
  definition: The target that gets the alert.
```

An alert can go to a pager or to a chat channel. A target is
never both, because each target is exactly one listed kind. An alert
cannot go to a service, because a service is not a target. A [constraint](constraint.md) with `count` can make sure that each
alert goes to exactly one target
([links of a thing](#links-of-a-thing)).

A union is closed and final. Only the statement itself lists its kinds. No
other kind, and no other package, can add a kind to it. No kind extends a
union, and a union extends no kind.

- **KD7.** `oneOf` MUST list two or more kinds. A listed kind MUST NOT extend a
  different listed kind, directly or through other kinds, because then each
  thing of the narrower kind is two listed kinds at once.
- **KD8.** A kind with `oneOf` MUST NOT have `fields`. Its listed kinds
  keep their own fields.
- **KD9.** A kind MUST NOT extend a kind with `oneOf`.
- **KD10.** A fact MUST NOT name a kind with `oneOf` as its `kind`, and a
  derivative MUST NOT name it in `becomes`. They name one of the listed
  kinds.
- **KD11.** After the derivatives run, each thing of a union kind MUST be a
  thing of exactly one listed kind. The tool rejects a thing that a
  derivative makes a thing of two listed kinds.
- **KD12.** A kind with `oneOf` MUST NOT have `extends`.

## Fields

A field is a value of the thing itself, for example a name, an availability or a
number of seconds. A subset of JSON Schema describes each field, and each keyword
has the meaning that JSON Schema gives it. Only field values use JSON Schema.
Kinds, links, constraints and derivatives use the syntax of the language
itself, because JSON Schema describes one document and has no links between
things.

```yaml
# defs/jobs.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name]

Job:
  is: Kind
  definition: A program that runs at set times and then stops, and that answers no requests.
  fields:
    name:
      type: Name                 # a value type, imported
    status:
      enum: [active, retired]
      default: active
    attempts:
      type: integer
      minimum: 1
      optional: true
    schedule:
      type: object
      properties:
        timezone: { type: string, optional: true }
        times:    { type: string }
```

### Keywords

| Keyword | For |
|---|---|
| `type`: `string`, `integer`, `number`, `boolean`, `array`, `object`, or a value type's name | the type of the value |
| `enum`, `const`, `default` | fixed values |
| `pattern`, `minLength`, `maxLength`, `format` (`date`, `date-time`, `duration`, `uri`) | text |
| `minimum`, `maximum` | numbers |
| `items` | the values of an `array` |
| `properties` | the parts of an `object` |
| `oneOf` | [variants](#variants): the value matches exactly one of the listed schemas |
| `optional: true` | a fact can omit the field |

There are five differences from JSON Schema. Each difference agrees with a
decision in a different part of the language:

| JSON Schema | Here | Why |
|---|---|---|
| `required: [ … ]` lists required parts | every field and every part of an `object` is required, unless it has `optional: true` | required is the default in all places |
| `$ref` refers to a different schema | `type: Name` names an imported value type | names come from imports |
| `null` is a type | there is no null type | a field with no value has no value; there is no third state |
| `oneOf`, `anyOf`, `allOf`, `not`, `if`/`then`/`else` combine schemas | only `oneOf` | `oneOf` gives variants; value types give reuse; conditional logic goes in constraints |
| `additionalProperties` controls if parts that are not listed are allowed | always false | the language rejects all items that are not listed |

`number` is exact, so `0.1 + 0.2` is `0.3`. The export keeps it as an exact
decimal number, never as a floating-point number.

### Variants

Some values have one of a small number of shapes. `oneOf` lists the shapes,
and each value matches exactly one of them. When the shapes are objects,
each shape SHOULD have a part with a fixed value, so that a reader, a
constraint and a derivative can see which shape a value has
([FD8](#field-rules)). Constraints and derivatives read a part with the
[`part` key](constraint.md#filter-and-unless).

A choice between sorts of things, such as a pager or a chat channel, is not a
field. It is a [union of kinds](#unions).

### Field rules

- **FD1.** A field MUST use only the keywords in the table above.
- **FD2.** `required`, `$ref`, `null`, `anyOf`, `allOf`, `not`,
  `if`, `then`, `else` and `additionalProperties` MUST NOT appear. The error
  tells what to write instead.
- **FD3.** A field is required, unless it has `optional: true` or a
  `default`. A field SHOULD NOT have both, because a default already makes
  the field optional.
- **FD4.** `pattern` MUST be a valid regular expression. `format` MUST be
  `date`, `date-time` or `duration` (ISO 8601), or `uri`: an absolute URI
  (Uniform Resource Identifier, RFC 3986), with a scheme such as `https:`.
- **FD5.** The rest of the field's schema MUST allow each value in `enum`,
  the `default` and the `const`.
- **FD6.** There is no null. A field with no value is omitted, and never has
  the value `null`.
- **FD7.** Each choice in `oneOf` MUST be a schema that uses the keywords
  above, but not `optional` or `default`. A value MUST match exactly one
  choice.
- **FD8.** When the choices of a `oneOf` are objects, each choice SHOULD have
  a part with the same name and a different `const` value.

## Value types

A value type gives a name to a schema one time, so that many fields can use
the schema again. A value type can hold a `oneOf`. Then all the fields that
name it share the same variants. A value type is a statement like all other statements, and
is imported in the same way:

```yaml
# defs/names.of.yaml
Name:
  is: ValueType
  tags: [privacy, backups]
  definition: A lowercase identifier for services, tables, data classes and regions.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"
  maxLength: 63

Seconds:
  is: ValueType
  tags: [privacy]
  definition: A number of seconds that a service waits for an answer.
  type: integer
  minimum: 1
  maximum: 3600

URL:
  is: ValueType
  tags: [privacy]
  definition: The address of a page on the web.
  type: string
  format: uri

Percent:
  is: ValueType
  tags: [privacy]
  definition: A share of a whole, in parts of one hundred.
  type: number
  minimum: 0
```

```yaml
# defs/service.of.yaml
import:
  - from: file:names.of.yaml
    use: [Name, Seconds, Percent, URL]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `ValueType` |
| `definition` | yes | one sentence that tells what the value holds; product's pages show it |
| `tags` | no | the [tags](tag.md) that the value type belongs to |
| any keyword above, but not `optional` | — | the schema |

- **VT1.** A value type MUST have a `definition`.
- **VT2.** A value type MUST NOT use the name of a built-in type (`string`,
  `integer`, `number`, `boolean`, `array`, `object`). Built-in types do not
  need an import.
- **VT3.** A field that names a value type MAY narrow it, but MUST NOT widen
  it. For example, `type: Name` with `maxLength: 30` is allowed, but with
  `maxLength: 100` it is an error. This rule also applies to `minLength`,
  `minimum`, `maximum` and `format`.
- **VT4.** The language checks a value against the value type and against
  the field's own keywords. The value must satisfy both.
- **VT5.** A value type is a type of its own, not only a narrower schema. Two
  fields have the same type only if one of these conditions is true:
  - Both fields name the same value type.
  - Neither field names a value type, and both fields have the same built-in
    type.

  `Seconds` and a plain `integer` are different types. A field that
  narrows a value type ([VT3](#value-types)) keeps that value type.
- **VT6.** A field that names a value type with `oneOf` MUST NOT add
  keywords to it, other than `optional` and `default`.
