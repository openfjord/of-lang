---
title: Facts — The OpenFjord modeling language
---

# Facts

> Part of [the reference](README.md).

A fact tells what a thing is: its kind, its fields and its links. Facts that
files contain are stated facts. They occur only in the tests of a package,
where they are the data for examples and questions. Facts that derivatives add
are [derived](#derived-facts). Facts that running systems write are deferred
(see [the research](../../README.md#deferred-facts-at-runtime)).

```yaml
import:
  - from: file:../../defs/table.of.yaml
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
| the fact's id | — | the name of the statement: any text, usually a path such as `db/sales/table/order_lines` |
| `is` | yes | `Fact` |
| `kind` | yes | the kind of the thing |
| the name of a field | as the kind tells | the value of the field |
| the name of a link | as the counts of the kind tell | one id, or a list of ids |

## Rules

- **FC1.** A fact MUST name its `kind`, and the kind MUST be imported.
- **FC2.** Each other key MUST be a field of the kind (inherited fields
  included) or a link that the file imports. Fields do not need an import,
  but links do.
- **FC3.** The facts about a thing, stated and derived together, MUST have:
  - every field that its kind requires,
  - every link that the counts of its kind require,
  - and nothing that its kind does not list.

  An omitted field takes its `default`, if it has one.
- **FC4.** Each field value MUST be a value that its field allows.
- **FC5.** The value of a link MUST be the id of a fact in the same file, or a
  list of such ids. Each of these facts MUST be of the kind where the link
  ends. A fact MAY use an inverse name: `hasTable: db/sales/table/order_lines`
  on a database states the same fact as `inDatabase: db/sales` on the table.
- **FC6.** Facts define nothing: a fact never declares a kind, a link or a
  field.

## Derived facts

A [derivative](derivative.md) adds a derived fact. The `under` or `chains`
key of a link can also add one ([what links imply](link.md#what-links-imply)).
No file declares a derived fact, because then it is not derived. A derived
fact exists only after the derivatives run, and the tool shows it in its
output ([how the tool shows derived facts](../../cli/README.md#derived-facts)).
A test can expect a derived fact, or expect that it is absent, under `expect`
([tests](test.md#examples)). Constraints and questions see stated facts and
derived facts in the same way.

Each derived fact records why it is true. The record contains the derivative
or declaration that added the fact, what each variable stands for, and the
facts that the derived fact comes from. An explanation shows that record.
