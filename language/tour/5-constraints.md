---
title: 5. Constraints — Tour
---

# 5. Constraints

> Part of [the tour](README.md). Previous: [4. Rules](4-rules.md).
> Next: [6. Questions and docs](6-questions-and-docs.md).

Rules derive. Constraints reject. Nothing does both.

## Applications and contracts

`defs/application.yaml` adds applications, which run jobs, and contracts: an
application's promise to keep some tables stable for others.

```yaml
Application:
  is: Kind
  definition: A deployable program that runs jobs.
  fields:
    name:
      type: Name
  links:
    dependsOnApp: any number
    consumesContract: any number

Contract:
  is: Kind
  definition: An application's promise to keep some of its tables stable for others to read.
  fields:
    name:
      type: Name
  links:
    coversTable: at least 1
```

with the links `consumesContract` (application to contract) and
`coversTable` (contract to table), and `jobOf` (job to application, inverse
`hasJob`) in `job.yaml`.

## A constraint

Teams kept breaking each other by reading internal tables. The rule: reading
another application's table needs a public table and a contract. It concerns
several kinds, so it sits a level up, in
`defs/reading-across-applications.yaml`:

```yaml
crossAppReadNeedsContract:
  is: Constraint
  says: >
    A job of one application may read a table that a job of another
    application writes, only if the table is public and the reading
    application consumes a contract that covers the table.
  match:
    - link: jobOf
      from: {$job: Job}
      to:   {$jobApp: Application}
    - link: readsTable
      from: {$job: Job}
      to:   {$table: Table}
    - link: writesTable
      from: {$writer: Job}
      to:   {$table: Table}
    - link: jobOf
      from: {$writer: Job}
      to:   {$writerApp: Application}
  filter:
    - node: $writerApp
      does not equal:
        same as: $jobApp
  require:
    - node:  $table
      field: visibility
      equals:
        value: public
    - at least one:
        match:
          - link: consumesContract
            from: {$jobApp: Application}
            to:   {$contract: Contract}
          - link: coversTable
            from: {$contract: Contract}
            to:   {$table: Table}
  message: "{$job} must not read {$table} from {$writerApp} without a contract."
```

Read it as: for every job reading a table that a job of another application
writes, the table must be public, and a contract the reading application
consumes must cover it.

- `match` and `filter` find the cases the constraint is about.
- `require` lists what every case must satisfy. `at least one` holds when its
  nested pattern matches at least once.
- `message` is what a violation says, worded with "must" or "must not", with
  variables filled in from the case.

## A rejection

Give revenue reporting a job that reads order intake's internal table, with
no contract:

```
$ of-lang run data-model data-model/tests/examples/cross-app-read-without-contract.yaml
derived:
  ...
rejected:
  crossAppReadNeedsContract: app/revenue_report/job/compute_revenue must not read db/sales/table/order_lines from app/order_intake without a contract.
```

`run` exits with 1 when anything is rejected, so CI stops the change. A test
can expect the rejection, which proves the constraint catches what it should:

```yaml
crossAppReadWithoutContract:
  is: Test
  says: >
    Revenue reporting reads an internal table of order intake with no
    contract, so the read is rejected.
  expect:
    rejected: [crossAppReadNeedsContract]
```

## Without `require`

Some rules read better as what must not happen. Leave out `require`, and every
match is a violation:

```yaml
noFeedbackLoop:
  is: Constraint
  says: A table must not feed itself, directly or through other tables.
  match:
    - link: feedsTable
      from: {$table: Table}
      to:   {$table: Table}
  message: "{$table} must not feed itself."
```

Using `$table` at both ends means the same table. Because `feedsTable` chains,
this also catches a loop through other tables.

Constraints run last, after every rule, over stated and derived facts alike.
Whatever is neither stated nor derived is false: a missing contract is
missing, not unknown.
