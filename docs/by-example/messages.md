---
title: Messages — The OpenFjord modeling language by example
---

# Messages

> Part of [the language by example](README.md).

The message of a constraint is what a violation says. It quotes variables in
braces, and each match gives their values.

<table>
<tr>
<td>

A batch states the highest temperature that is safe for it, and the
temperature of the place where the cooks keep it. The shelf is optional.

</td>
<td>

```yaml
# defs/batch.of.yaml
Batch:
  is: Kind
  definition: An amount of one ingredient that the kitchen received on one day.
  fields:
    keepAtMostCelsius:
      type: integer
    keptAtCelsius:
      type: integer
    shelf:
      type: string
      optional: true
```

</td>
</tr>
<tr>
<td>

The constraint matches each batch that is too warm.

</td>
<td>

```yaml
batchKeptColdEnough:
  is: Constraint
  says: The cooks keep each batch cold enough.
  match:
    - kind: { name: Batch, as: $batch }
  filter:
    - node:  $batch
      field: keptAtCelsius
      is above:
        same as:
          node:  $batch
          field: keepAtMostCelsius
```

</td>
</tr>
<tr>
<td>

The message states the rule that the match breaks, with the word "must".
`{$batch}` gives the id of the thing. `{$batch: shelf}` gives the value of
a field.

</td>
<td>

```yaml
  message: >-
    {$batch} on shelf {$batch: shelf} must be kept at
    {$batch: keepAtMostCelsius} °C or colder, but it is at
    {$batch: keptAtCelsius} °C.
```

</td>
</tr>
<tr>
<td>

A test file states two batches that are too warm. The first one states its
shelf.

</td>
<td>

```yaml
# tests/examples/two-warm-batches.of.yaml
import:
  - from: file:../../defs/batch.of.yaml
    use: [Batch, batchKeptColdEnough]

batch/parmesan_monday:
  is: Fact
  kind: Batch
  fields:
    keepAtMostCelsius: 8
    keptAtCelsius: 12
    shelf: B2
```

</td>
</tr>
<tr>
<td>

The second one states no shelf.

</td>
<td>

```yaml
batch/cream_tuesday:
  is: Fact
  kind: Batch
  fields:
    keepAtMostCelsius: 4
    keptAtCelsius: 6
```

</td>
</tr>
<tr>
<td>

The test expects that the facts break the constraint.

</td>
<td>

```yaml
twoWarmBatches:
  is: Test
  says: The two batches are warmer than is safe for them, so the constraint rejects them.
  expect:
    rejected: [batchKeptColdEnough]
```

</td>
</tr>
</table>

`run` prints one line for each violation: the name of the constraint, then
the message with the values of that match. A field with no value gives the
words "not set":

```console
$ of-lang run messages messages/tests/examples/two-warm-batches.of.yaml
derived:
  nothing
rejected:
  batchKeptColdEnough: batch/parmesan_monday on shelf B2 must be kept at 8 °C or colder, but it is at 12 °C.
  batchKeptColdEnough: batch/cream_tuesday on shelf not set must be kept at 4 °C or colder, but it is at 6 °C.
```

More about this: [messages](../../language/reference/constraint.md#messages).
