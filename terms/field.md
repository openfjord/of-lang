---
title: Fields and value types — A modeling language for software people
---

# Fields and value types

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A field's value is described with JSON Schema — the keywords below, with the
meaning JSON Schema gives them. Only field values use it: kinds, links,
constraints and rules are this language's own, because JSON Schema describes
one document and has no links between things.

| Kept | For |
|---|---|
| `type`: `string`, `integer`, `number`, `boolean`, `array`, `object` | the kind of value |
| `enum`, `const`, `default` | fixed values |
| `pattern`, `minLength`, `maxLength`, `format` (`date`, `date-time`, `duration`) | text |
| `minimum`, `maximum` | numbers |
| `items` | the values of an `array` |
| `properties` | the parts of an `object` |

Five differences from JSON Schema, each to keep a decision made elsewhere in
this draft:

| JSON Schema | Here | Why |
|---|---|---|
| `required: [ … ]` lists required parts | every field and every part of an `object` is required unless it says `optional: true` | required is the default everywhere in the language |
| `$ref` points at another schema | `type: Name` names an imported value type | names come from imports |
| `null` is a type | there is none | a field with no value has none; no third state |
| `oneOf`, `anyOf`, `allOf`, `not`, `if`/`then`/`else` combine schemas | not available | hard to read; conditional logic belongs in constraints |
| `additionalProperties` decides whether unlisted parts are allowed | always false | anything not listed is rejected |

`number` is exact: it is exported as `xsd:decimal`, never as a floating-point
type. A structured value is an `object` whose `properties` are its parts:

```yaml
address:
  type: object
  properties:
    street: { type: string, optional: true }
    city:   { type: string }
```

## Value types

**Value types** name a schema once so fields can reuse it. They are entries
like any other, imported like any other:

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

- `type: Name` means the imported value type `Name`. Built-in types need no
  import, and no value type may take a built-in type's name.
- A field may narrow a value type, never widen it: `type: Name` with
  `maxLength: 30` is allowed; with `maxLength: 100` it is an error.
- A value type has a `definition`, so product's pages can say what a field
  holds.
- A value type can be exported from its package and imported elsewhere with
  `pkg:`.
- Conditions compare only values of the same type.
