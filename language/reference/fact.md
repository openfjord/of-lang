---
title: Facts — The OpenFjord modeling language
---

# Facts

> Part of [the reference](README.md).

A fact is a thing, stated: its kind, its fields and its links. Facts appear
only in a package's tests, where they are the data examples and questions run
over. Facts written by running systems are set aside for now (see
[the research](../../README.md#deferred-facts-at-runtime)).

```yaml
import:
  - from: file:../../defs/table.yaml
    use: [Table, inDatabase]

db/sales/table/order_lines:
  is: Fact
  kind: Table
  name: order_lines
  visibility: public
  inDatabase: db/sales
```

| Key | Required | Value |
|---|---|---|
| the fact's id | — | the entry's name: any text, usually a path such as `db/sales/table/order_lines` |
| `is` | yes | `Fact` |
| `kind` | yes | the kind of thing it is |
| a field's name | as the kind says | the field's value |
| a link's name | as the kind's counts say | one id, or a list of ids |

## Rules

- **FC1.** A fact MUST name its `kind`, and the kind MUST be imported.
- **FC2.** Each other key MUST be a field of the kind (inherited ones
  included) or a link the file imports. Fields need no import; links do.
- **FC3.** A fact MUST have every field its kind requires and every link its
  kind's counts require, and nothing its kind does not list. A field left out
  takes its `default`, if it has one.
- **FC4.** Every field value MUST be one its field allows.
- **FC5.** A link's value MUST be the id of a fact in the same file, or a
  list of such ids, and each MUST be of the kind the link ends at. An inverse
  name MAY be used: `hasTable: db/sales/table/order_lines` on a database
  states the same fact as `inDatabase: db/sales` on the table.
- **FC6.** Facts define nothing: a fact never declares a kind, a link or a
  field.

## Derived facts

A derived fact is added by a [rule](rule.md), or by a link's `under` or
`chains` ([what links imply](link.md#what-links-imply)); it is never written
in a file. Derived facts live in a graph of their own, apart from stated
facts, and each records why it holds: the rule or declaration that added it,
what each variable stood for, and the facts it came from. That record is what
an explanation shows.
