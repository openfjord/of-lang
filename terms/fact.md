---
title: Facts — A modeling language for software people
---

# Facts

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

Facts appear only in a package's tests and examples. A fact is an entry
keyed by its id, with `is: Fact` and its kind:

```yaml
import:
  - from: file:../../defs/table.yaml
    use: [Table, inDatabase, hasColumn]

db/sales/table/order_lines:
  is: Fact
  kind: Table
  name: order_lines
  format: columnar
  visibility: public
  inDatabase: db/sales
  hasColumn: [order_id, sku, quantity]
```

- `kind` names an imported kind, and no field may be called `kind`.
- Fields belong to the kind, so they need no import. Links are entries of
  their own, so every link used as a key is imported.
- A link's value is one id or a list of ids, each a fact in the same set of
  facts, whose kind the link allows.
- The fact is checked against its kind: every field and link it must have,
  every value its field type allows, nothing that is not listed.

## Derived facts

A derived fact is added by a [rule](rule.md), never written in a file.
Derived facts live in a graph of their own, each recording the rule and the
match that produced it; see [how patterns run](pattern.md#how-they-run).
