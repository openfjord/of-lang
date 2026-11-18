---
title: Imports — The OpenFjord modeling language by example
---

# Imports

> Part of [the language by example](README.md).

A file can use a name only if the file defines that name or imports it. The
`import` block at the top of a file shows the source of each name.

<table>
<tr>
<td>

This file defines the kind `Team`. It uses no other name, so it has no
`import` block.

</td>
<td>

```yaml
# defs/team.of.yaml
Team:
  is: Kind
  definition: A group of people who build and run services.
```

</td>
</tr>
<tr>
<td>

A second file uses `Team`, so it imports the name. `from` gives the
source: `file:`, then a path relative to this file. `use` lists the names.

</td>
<td>

```yaml
# defs/service.of.yaml
import:
  - from: file:team.of.yaml
    use: [Team]
```

</td>
</tr>
<tr>
<td>

The file defines `Service` itself, so it does not import that name.

</td>
<td>

```yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
```

</td>
</tr>
<tr>
<td>

The link uses the two names: one that the file defines and one that it
imports.

</td>
<td>

```yaml
ownedBy:
  is: Link
  from: Service
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that builds and runs the service.
```

</td>
</tr>
<tr>
<td>

A test file imports each name that it uses, with one item for each source
file. `../` goes to the parent folder.

</td>
<td>

```yaml
# tests/examples/checkout-has-an-owner.of.yaml
import:
  - from: file:../../defs/team.of.yaml
    use: [Team]
  - from: file:../../defs/service.of.yaml
    use: [Service, ownedBy]
```

</td>
</tr>
<tr>
<td>

The facts use the two kinds and the link, so the file imports all three
names.

</td>
<td>

```yaml
team/payments:
  is: Fact
  kind: Team

service/checkout:
  is: Fact
  kind: Service
  links:
    ownedBy: team/payments

checkoutHasAnOwner:
  is: Test
  says: A service that a team owns is valid.
  expect:
    rejected: []
```

</td>
</tr>
</table>

`test` reads the files with their imports and runs the test:

```console
$ of-lang test imports
pass  test checkoutHasAnOwner  (imports/tests/examples/checkout-has-an-owner.of.yaml)
1 passed, 0 failed
```

More about this: [imports](../../language/reference/files.md#imports).
