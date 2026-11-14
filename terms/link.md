---
title: Links — A modeling language for software people
---

# Links

> Part of [a modeling language for software people](../README.md), a first
> draft. Reason numbers refer to the table there.

A link is declared once, as its own entry, and the kind it comes out of says
how many it carries. The link usually sits in the same file as that kind, so
an object's fields, the links that come out of it and the limits on them read
together.

```yaml
inDatabase:
  is: Link
  from: Table
  to: Database
  reads: "{from} is in {to}"
  definition: The database that stores the table.
  inverse: hasTable

jobOf:
  is: Link
  from: Job
  to: Application
  reads: "{from} is a job of {to}"
  definition: The application that runs the job.
  inverse: hasJob

readsTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} reads {to}"
  definition: The job reads rows from the table when it runs.
  under: usesTable                   # every readsTable is also a usesTable
  inverse: readByJob

writesTable:
  is: Link
  from: Job
  to: Table
  reads: "{from} writes {to}"
  definition: The job writes rows to the table when it runs.
  under: usesTable
  inverse: writtenByJob

feedsTable:
  is: Link
  from: Table
  to: Table
  reads: "{from} feeds {to}"
  definition: >
    Rows of the first table flow, directly or through other tables, into the
    second.
  chains: true                       # a feeds b and b feeds c means a feeds c
```

Rules on links:

- **The name says what it links and how** — two words at least, never a bare
  verb: `readsTable`, never `reads` (reason 2).
- **A definition is required**, and so is `reads`, the sentence template that
  product sees (reason 7).
- **`from` and `to` are checks.** A team that `readsTable` is rejected. No
  link ever makes something into a kind.
- **`under`** places a link beneath a broader one, so specific links can be
  written and general questions still answered.
- **`inverse`** and **`chains`** are the only built-in inference on links.
