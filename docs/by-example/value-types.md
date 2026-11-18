---
title: Value types — The OpenFjord modeling language by example
---

# Value types

> Part of [the language by example](README.md).

A value type gives a name to a schema one time, so that many fields can use
the schema again.

<table>
<tr>
<td>

A value type is a statement with `is: ValueType`. Like a kind, it needs a
`definition`.

</td>
<td>

```yaml
# defs/types.of.yaml
Minutes:
  is: ValueType
  definition: A number of minutes.
```

</td>
</tr>
<tr>
<td>

The other keys are the schema: a type and the keywords that go with it. A
number of minutes is a whole number from 1 to 1440.

</td>
<td>

```yaml
  type: integer
  minimum: 1
  maximum: 1440
```

</td>
</tr>
<tr>
<td>

A file imports each name that it uses but does not define. The file of the
job imports `Minutes`.

</td>
<td>

```yaml
# defs/job.of.yaml
import:
  - from: file:types.of.yaml
    use: [Minutes]
```

</td>
</tr>
<tr>
<td>

A field names the value type as its `type`. Two fields of the job use
`Minutes`: the time between two starts, and the longest time that one run
can take.

</td>
<td>

```yaml
Job:
  is: Kind
  definition: A program that runs on a schedule.
  fields:
    name:
      type: string
    everyMinutes:
      type: Minutes
    timeoutMinutes:
      type: Minutes
```

</td>
</tr>
<tr>
<td>

A test file states two jobs. Each value is a valid number of minutes, so
the test passes.

</td>
<td>

```yaml
# tests/examples/two-jobs.of.yaml
import:
  - from: file:../../defs/job.of.yaml
    use: [Job]

job/backup:
  is: Fact
  kind: Job
  fields:
    name: backup
    everyMinutes: 1440
    timeoutMinutes: 120

job/report:
  is: Fact
  kind: Job
  fields:
    name: report
    everyMinutes: 15
    timeoutMinutes: 30

twoJobs:
  is: Test
  says: The times of the backup job and of the report job are valid numbers of minutes.
  expect:
    rejected: []
```

</td>
</tr>
<tr>
<td>

A value type is a type of its own. A comparison needs two values of the
same type, and the two fields both name `Minutes`. Thus, a question can
compare them.

</td>
<td>

```yaml
# tests/questions/which-jobs-can-run-into-their-next-start.of.yaml
import:
  - from: file:../../defs/job.of.yaml
    use: [Job]

whichJobsCanRunIntoTheirNextStart:
  is: Question
  asks: Which jobs can run for longer than the time between two starts?
  given: file:../examples/two-jobs.of.yaml
  match:
    - kind: { name: Job, as: $job }
  filter:
    - node: $job
      field: timeoutMinutes
      is above:
        same as:
          node: $job
          field: everyMinutes
  expect:
    $job: [job/report]
```

</td>
</tr>
</table>

`test` runs the test and the question:

```console
$ of-lang test value-types
pass  test twoJobs  (value-types/tests/examples/two-jobs.of.yaml)
pass  question whichJobsCanRunIntoTheirNextStart  (value-types/tests/questions/which-jobs-can-run-into-their-next-start.of.yaml)
2 passed, 0 failed
```

More about this: [value types](../../language/reference/kind.md#value-types).
