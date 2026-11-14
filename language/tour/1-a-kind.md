---
title: 1. A kind — Tour
---

# 1. A kind

> Part of [the tour](README.md). Next: [2. Facts and a test](2-facts-and-a-test.md).

## A package

Everything lives in a package: a folder with a manifest, `package.yaml`, and a
`defs/` folder for definitions.

```
data-model/
  package.yaml
  defs/
```

```yaml
# package.yaml
package: "@openfjord/data-model"
version: 0.1.0
owner: product
rdfNamespace: data
```

- `package` is the name. `@openfjord` is its namespace, as in npm. YAML
  needs the quotes because of the leading `@`.
- `owner` owns everything in the package. Any change needs their approval.
- `rdfNamespace` is used only when exporting to RDF. You never see it in the
  files.

## A kind

A kind is a kind of thing. Write the first one in `defs/team.yaml`:

```yaml
Team:
  is: Kind
  definition: A group of people who answer for databases and applications.
  fields:
    name:
      type: string
```

Every statement in a file has a name (`Team`) and says what it is with `is:`.
Every kind needs a `definition`: one sentence that product reads.

`fields` lists the values a team has. A field is described with JSON Schema,
the same keywords you know from API specs. Every field is required unless it
says `optional: true` or has a `default`.

## A value type

Teams, databases and tables all have names that follow the same rule. Give the
rule a name once, in `defs/names.yaml`:

```yaml
Name:
  is: ValueType
  definition: A lowercase identifier used for teams, databases, tables, applications and jobs.
  type: string
  pattern: "^[a-z][a-z0-9_]*$"
  maxLength: 63
```

and use it from `team.yaml`, which must now import it:

```yaml
import:
  - from: file:names.yaml
    use: [Name]

Team:
  is: Kind
  definition: A group of people who answer for databases and applications.
  fields:
    name:
      type: Name
```

A file may use a name only if it defines it or imports it, even from a file
next to it. The top of every file therefore lists everything it depends on.

## Check it

```
$ of-lang check data-model
ok
```

Misspell a key and the check fails. Nothing unknown is ignored:

```yaml
Team:
  is: Kind
  defintion: A group of people who answer for databases and applications.
```

```
$ of-lang check data-model
data-model/defs/team.yaml:1:1: Team must have a definition (KD1)
data-model/defs/team.yaml:3:3: Team has no key defintion (FL2)
```

Each error names the rule it enforces, so you can look it up in
[the reference](../reference/README.md).
