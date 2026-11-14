---
title: Fields and value types — The OpenFjord modeling language
---

# Fields and value types

> Part of [the reference](README.md).

A field is a value of the thing itself: a name, a format, a number of days.
Each field is described with a subset of JSON Schema, with the meaning JSON
Schema gives each keyword. Only field values use JSON Schema; kinds, links,
constraints and rules are the language's own, because JSON Schema describes
one document and has no links between things.

```yaml
fields:
  name:
    type: Name                 # a value type, imported
  visibility:
    enum: [public, internal]
    default: internal
  retentionDays:
    type: integer
    minimum: 1
    optional: true
  address:
    type: object
    properties:
      street: { type: string, optional: true }
      city:   { type: string }
```

## Keywords

| Keyword | For |
|---|---|
| `type`: `string`, `integer`, `number`, `boolean`, `array`, `object`, or a value type's name | the kind of value |
| `enum`, `const`, `default` | fixed values |
| `pattern`, `minLength`, `maxLength`, `format` (`date`, `date-time`, `duration`) | text |
| `minimum`, `maximum` | numbers |
| `items` | the values of an `array` |
| `properties` | the parts of an `object` |
| `optional: true` | the field may be left out |

Five differences from JSON Schema, each to keep a decision made elsewhere in
the language:

| JSON Schema | Here | Why |
|---|---|---|
| `required: [ … ]` lists required parts | every field and every part of an `object` is required unless it says `optional: true` | required is the default everywhere |
| `$ref` points at another schema | `type: Name` names an imported value type | names come from imports |
| `null` is a type | there is none | a field with no value has none; no third state |
| `oneOf`, `anyOf`, `allOf`, `not`, `if`/`then`/`else` combine schemas | not available | hard to read; conditional logic belongs in constraints |
| `additionalProperties` decides whether unlisted parts are allowed | always false | anything not listed is rejected |

`number` is exact: `0.1 + 0.2` is `0.3`, and it is exported as
`xsd:decimal`, never as a floating-point type.

## Rules

- **FD1.** A field MUST use only the keywords in the table above.
- **FD2.** `required`, `$ref`, `null`, `oneOf`, `anyOf`, `allOf`, `not`,
  `if`, `then`, `else` and `additionalProperties` MUST NOT appear; the error
  says what to write instead.
- **FD3.** A field is required unless it says `optional: true` or has a
  `default`. A field SHOULD NOT have both: a default already makes it
  optional.
- **FD4.** `pattern` MUST be a valid regular expression. `format` MUST be
  `date`, `date-time` or `duration` (ISO 8601).
- **FD5.** Every value in `enum`, and the `default` and `const`, MUST be
  allowed by the rest of the field's schema.
- **FD6.** There is no null: a field with no value is left out, never set to
  `null`.

## Value types

A value type names a schema once so fields can reuse it. It is an entry like
any other, imported like any other:

```yaml
# defs/names.yaml

Name:
  is: ValueType
  definition: A lowercase identifier used for databases, tables and jobs.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"
  maxLength: 63

RetentionDays:
  is: ValueType
  definition: How long rows are kept, in days.
  type: integer
  minimum: 1
  maximum: 3650
```

```yaml
# defs/table.yaml

import:
  - from: file:names.yaml
    use: [Name, RetentionDays]
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `ValueType` |
| `definition` | yes | one sentence saying what the value holds, shown on product's pages |
| any keyword above except `optional` | — | the schema |

- **VT1.** A value type MUST have a `definition`.
- **VT2.** A value type MUST NOT take a built-in type's name (`string`,
  `integer`, `number`, `boolean`, `array`, `object`). Built-in types need no
  import.
- **VT3.** A field that names a value type MAY narrow it but MUST NOT widen
  it: `type: Name` with `maxLength: 30` is allowed; with `maxLength: 100` it
  is an error. The same holds for `minLength`, `minimum`, `maximum` and
  `format`.
- **VT4.** A value is checked against both the value type and the field's
  own keywords, and must satisfy both.
- **VT5.** Conditions compare only values of the same type.
