---
title: Pending changes — The OpenFjord modeling language
---

# Pending changes

> Part of [the reference](../README.md#reference).

A pending change is a change to a package that nobody has accepted yet. The
change is in the files of the package, as marks on the statements that it
adds, removes or updates. Thus, the tool and a user interface can show what the
change does without git.

```yaml
# defs/service.of.yaml
Engineer:                    # the change adds this statement
  is: Kind
  pending: added
  definition: A person who builds and runs services in one or more teams.

serviceHasOneOwner:          # the change removes this statement
  is: Constraint
  pending: removed
  says: Exactly one team owns each service.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - count:
        match:
          - link: ownedBy
            from: $service
            to:   Team
      does not equal:
        value: 1
  message: "{$service} must have exactly one team that owns it."

SharedService:               # the change updates this statement
  is: Kind
  pending: updated
  extends: Service
  definition: A service that three or more other services depend on, so a change to it needs extra care.
  was:
    extends: Service
    definition: A service that many other services depend on, so a change to it needs extra care.
```

`pending` and `was` are annotations: keys that the language adds to any
statement, next to the keys of the statement itself. They never clash with
the names of fields or links, because those names are always one level lower
([statements](files.md#statements)).

| Key | Required | Value |
|---|---|---|
| `pending` | no | `added`, `removed` or `updated` |
| `was` | only with `pending: updated` | the whole statement as it was before the change, without `is`, `pending` and `was` |

An updated statement shows its new version. `was` keeps the old version, so
that the tool can show the difference and can reject the change. In this
example, the update changes the definition of `SharedService`.

## One change at a time

All the marks in a package together are its pending change. A package has
one pending change or none. The change has no name.

The package has two versions:

| Version | Added statements | Removed statements | Updated statements |
|---|---|---|---|
| without the change | not included | included, as written | the `was` version |
| with the change | included | not included | the new version |

The tool checks the two versions. The [change report](../../cli/README.md#the-change-report)
compares them: what the change states, what newly follows, what no longer
follows, and what is rejected.

Other packages see only the version without the change. A pending change has
no effect outside its package until somebody accepts it.

## Accepting and rejecting

To accept the change, the tool does these steps as one operation:

1. It deletes each removed statement.
2. It deletes `pending` and `was` from each added and updated statement.

To reject the change, the tool does these steps as one operation:

1. It deletes each added statement.
2. It replaces each updated statement with its `was` version.
3. It deletes `pending` from each removed statement.

## Rules

- **PD1.** Any statement in `defs/` or `tests/` MAY have `pending`, with the
  value `added`, `removed` or `updated`.
- **PD2.** A statement with `pending: updated` MUST have `was`. `was` MUST
  hold the whole statement as it was before the change, without `is`,
  `pending` and `was`. A statement with a different `pending` value MUST NOT
  have `was`.
- **PD3.** An update MUST NOT change `is`. To change what a statement is,
  the change removes the statement and adds a different one.
- **PD4.** The version without the change and the version with the change
  MUST each pass every rule in this specification.

The manifest and the Markdown files in `docs/` have no
pending marks. A pending change to them is deferred.
