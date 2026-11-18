---
title: Enums — The OpenFjord modeling language by example
---

# Enums

> Part of [the language by example](README.md).

A field with `enum` has a fixed list of values. A fact gives one value from
the list.

<table>
<tr>
<td>

Each table has a `name`, which can be any text.

</td>
<td>

```yaml
# defs/table.of.yaml
Table:
  is: Kind
  definition: A set of rows with the same columns, in a database.
  fields:
    name:
      type: string
```

</td>
</tr>
<tr>
<td>

`enum` lists the values that `sensitivity` can have. The data in a table
is public, internal or confidential, and nothing else.

</td>
<td>

```yaml
    sensitivity:
      enum: [public, internal, confidential]
```

</td>
</tr>
<tr>
<td>

A test file states two tables. Each one gives a sensitivity from the list.

</td>
<td>

```yaml
# tests/examples/two-tables.of.yaml
import:
  - from: file:../../defs/table.of.yaml
    use: [Table]

table/products:
  is: Fact
  kind: Table
  fields:
    name: products
    sensitivity: public

table/customers:
  is: Fact
  kind: Table
  fields:
    name: customers
    sensitivity: confidential
```

</td>
</tr>
<tr>
<td>

A test passes only if every fact obeys its kind. This one passes, so the
field allows the two values.

</td>
<td>

```yaml
twoTables:
  is: Test
  says: A public table and a confidential table are both valid.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A question asks which tables are confidential, given those facts. A
condition compares the field with a value from the list.

</td>
<td>

```yaml
# tests/questions/which-tables-are-confidential.of.yaml
import:
  - from: file:../../defs/table.of.yaml
    use: [Table]

whichTablesAreConfidential:
  is: Question
  asks: Which tables are confidential?
  given: file:../examples/two-tables.of.yaml
  match:
    - kind: { name: Table, as: $table }
  filter:
    - node: $table
      field: sensitivity
      equals:
        value: confidential
```

</td>
</tr>
<tr>
<td>

The answer is the customers table.

</td>
<td>

```yaml
  expect:
    $table: [table/customers]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test enums
pass  test twoTables  (enums/tests/examples/two-tables.of.yaml)
pass  question whichTablesAreConfidential  (enums/tests/questions/which-tables-are-confidential.of.yaml)
2 passed, 0 failed
```

More about this: [the keywords of a field](../../language/reference/kind.md#keywords)
and [the field rules](../../language/reference/kind.md#field-rules).
