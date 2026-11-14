---
title: 7. Many packages — Tour
---

# 7. Many packages

> Part of [the tour](README.md). Previous: [6. Questions and docs](6-questions-and-docs.md).

One owner owns everything in a package. When teams belong to the people side
of the organization rather than to product, `Team` moves to a package of its
own, with its own owner.

## A second package

```
org-model/
  package.yaml
  defs/
    team.yaml
```

```yaml
# org-model/package.yaml
package: "@openfjord/org-model"
version: 1.2.0
owner: people-ops
rdfNamespace: org
exports:
  kinds: [Team]
```

A package's definitions are internal unless it lists them under `exports`.
Other packages may import only what is exported.

## Depending on it

`data-model` declares the dependency, with a version range:

```yaml
# data-model/package.yaml
package: "@openfjord/data-model"
version: 0.2.0
owner: product
rdfNamespace: data
dependsOn:
  "@openfjord/org-model": ^1.2
```

and imports `Team` with `pkg:` instead of `file:`:

```yaml
# data-model/defs/database.yaml
import:
  - from: pkg:@openfjord/org-model
    use: [Team]
```

`^1.2` means 1.2 or later, below 2.0. Two packages that each define `Team`
can both be imported; `as` renames one:

```yaml
    use: [{ name: Team, as: OrgTeam }]
```

## Extension points

A package says things only about its own definitions. `data-model` cannot add
a field to `Team`, or a kind of team. If `org-model` wants others to extend
it, it opens an extension point:

```yaml
openFor:
  - extend: Team        # others may declare a kind of team
```

`under: <link>` opens a link the same way, so another package may declare a
narrower link under it.

## Versions

A new version must say what changed. Removing or renaming a definition, or
tightening a constraint, breaks the packages that depend on it and needs a
new major version. Adding a definition or loosening a constraint does not.

```
$ of-lang check data-model --packages .
ok
```

`--packages` says where to find dependencies. By default, the tool looks in
the folder that holds the package.

## Where to go next

You have seen every part of the language. [The cheat sheet](../cheatsheet.md)
puts it on one page; [the reference](../reference/README.md) gives the exact
rules for each part.
