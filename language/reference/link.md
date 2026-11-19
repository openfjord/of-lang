---
title: Links — The OpenFjord modeling language
---

# Links

> Part of [the reference](../README.md#reference).

A link is a named relationship from a kind to a kind, in one direction. A link
is declared one time, as its own statement. Usually, it is in the same file as
the kind where it starts. A thing can have each link that starts at its kind,
at a kind that its kind extends, or at a union
([links of a thing](kind.md#links-of-a-thing)).

```yaml
# defs/service.of.yaml
ownedBy:
  is: Link
  from: Service
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that builds and runs the service.
  inverse: ownsService

writesTable:
  is: Link
  tags: [privacy]
  from: Service
  to: Table
  reads: "{from} writes {to}"
  definition: The service adds, changes or deletes rows of the table.
  under: readsTable                   # every writesTable is also a readsTable

dependsOn:
  is: Link
  tags: [privacy]
  from: Service
  to: Service
  reads: "{from} depends on {to}"
  definition: The first service needs an answer from the second service, directly or through other services.
  chains: true                       # a depends on b, b on c: a depends on c
  inverse: neededBy
```

| Key | Required | Value |
|---|---|---|
| `is` | yes | `Link` |
| `from` | yes | the kind where the link starts |
| `to` | yes | the kind where the link ends |
| `reads` | yes | the sentence that product sees, with `{from}` and `{to}` |
| `definition` | yes | one sentence that tells what the link means |
| `under` | no | the broader link; this link is a narrower sort of it |
| `inverse` | no | a name for the link in the opposite direction |
| `chains` | no | `true` when a chain of the link implies the link |
| `tags` | no | the [tags](tag.md) that the link belongs to |

A link has only its two ends. It has no fields. Any data about a link becomes
a thing in the middle, for example a `Grant` between a service and a table
that it uses.

## Rules

- **LK1.** A link's name, and its inverse name, MUST be lowerCamelCase, with
  two or more words that tell what the link connects and how: `dependsOn`,
  never `depends` (reason 2 in [the design notes](../../docs/design.md#why-ontologies-are-hard-for-software-people)).
- **LK2.** A link MUST have a `definition` and a `reads` sentence. `reads`
  MUST contain both `{from}` and `{to}` (reason 7).
- **LK3.** `from` and `to` MUST each name a kind that the file defines or
  imports. A name that is not a kind, or that names nothing, is an error.
  `from` and `to` are checks: the language
  rejects a fact that links a team with `dependsOn`. A link never gives a
  thing a kind.
- **LK4.** `under` MUST name a declared link, not an inverse name. The `from`
  of the narrower link MUST be the `from` of the broader link, or a kind that
  extends it. The same rule applies to `to`. A link MUST NOT be under itself,
  directly or through other links.
- **LK5.** In a link with `chains: true`, `to` MUST be the same kind as
  `from`, or a kind that extends it.

## What links imply

These are the only inferences that the language makes without a derivative.
Each fact that they add is a [derived fact](fact.md#derived-facts), and
records its reason.

| Declared | Adds |
|---|---|
| `inverse: ownsService` on `ownedBy` | nothing new: `ownsService` is every `ownedBy` in the opposite direction, in all derivatives, constraints, questions and facts |
| `under: readsTable` on `writesTable` | a `readsTable` for every `writesTable` |
| `chains: true` on `dependsOn` | `a dependsOn c` each time that `a dependsOn b` and `b dependsOn c` |
