---
title: Pending changes — The OpenFjord modeling language by example
---

# Pending changes

> Part of [the language by example](README.md).

A pending change is a change to a package that nobody has accepted yet. The
change is in the files of the package, as marks on the statements that it
adds, removes or updates.

<table>
<tr>
<td>

`pending: added` marks a statement that the change adds. The package
without the change does not include it.

</td>
<td>

```yaml
# defs/programs.of.yaml
Job:
  is: Kind
  pending: added
  definition: A program that runs on a schedule.
```

</td>
</tr>
<tr>
<td>

`pending: removed` marks a statement that the change removes. The statement
stays in the file, as written, until somebody accepts the change.

</td>
<td>

```yaml
Script:
  is: Kind
  pending: removed
  definition: A program that a person starts by hand.
```

</td>
</tr>
<tr>
<td>

`pending: updated` marks a statement that the change updates. The statement
shows its version with the change: here, a service has a `tier`.

</td>
<td>

```yaml
Service:
  is: Kind
  pending: updated
  definition: A program that runs in production and answers requests.
  fields:
    name:
      type: string
    tier:
      enum: [critical, standard]
      default: standard
```

</td>
</tr>
<tr>
<td>

`was` keeps the whole statement as it was before the change, without `is`.
Only an updated statement has `was`.

</td>
<td>

```yaml
  was:
    definition: A program that runs in production and answers requests.
    fields:
      name:
        type: string
```

</td>
</tr>
<tr>
<td>

The tests use the package with the change. Thus, a test file can state a
job, and a service with a tier. It cannot state a script.

</td>
<td>

```yaml
# tests/examples/programs-with-the-change.of.yaml
import:
  - from: file:../../defs/programs.of.yaml
    use: [Job, Service]

job/nightly_backup:
  is: Fact
  kind: Job

service/checkout:
  is: Fact
  kind: Service
  fields:
    name: checkout
    tier: critical
```

</td>
</tr>
<tr>
<td>

The test expects that the package with the change rejects nothing.

</td>
<td>

```yaml
programsWithTheChange:
  is: Test
  says: A job and a service with a tier are valid.
  expect:
    rejected: []
```

</td>
</tr>
</table>

All the marks in a package together are its pending change: a package has
one pending change or none. To accept the change is to apply all the marks
in one step, and to reject it is to undo them all. `test` runs the package
with the change:

```console
$ of-lang test pending-changes
pass  test programsWithTheChange  (pending-changes/tests/examples/programs-with-the-change.of.yaml)
1 passed, 0 failed
```

More about this: [pending changes](../../language/reference/pending.md).
