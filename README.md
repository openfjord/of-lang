---
title: A modeling language for software people
---

# A modeling language for software people

The OpenFjord modeling language describes the things of a business, the
links between them, and the rules about them, in YAML files next to the code.
Its tool, `of-lang`, checks those files, derives the facts that nobody wrote
down, explains each one, and runs the model's tests.

<p align="center">
  <img src="docs/images/of-lang-flow.svg" alt="Packages go into the OpenFjord modeling language, which checks, derives, explains and tests. An agent, an MCP server, a coding factory, a knowledge base and a semantic layer use what it knows." width="100%">
</p>

Packages go into the language. An agent, an MCP (Model Context Protocol)
server, a coding factory, a knowledge base and a semantic layer use what it
knows.

> **Version 0.1.0.** Before version 1.0, the language can change in ways
> that break a package.

<!-- not on the site -->
The documentation is also a site: <https://lang.openfjord.io/>.
<!-- end -->

## Install

With [Homebrew](https://brew.sh):

```bash
brew install openfjord/tap/of-lang
```

Or download a binary for macOS or Linux from the
[releases](https://github.com/openfjord/of-lang/releases).

## Three small examples

Each example is complete, and needs nothing from the other two. Each one is
a **package**: a folder with a manifest, definitions and facts.

```
wrong-value/
  package.yaml           # the manifest: the name and the owner of the package
  defs/
    services.of.yaml     # the definitions: what can exist, and the rules
  tests/
    facts.of.yaml        # the facts: what does exist, in one example
```

The manifest is four lines. Only the name is different in each example:

```yaml
# wrong-value/package.yaml
package: "@acme/wrong-value"
version: 0.1.0
owner: platform
rdfNamespace: wrong-value
```

The packages are also in [`examples/quickstart`](examples/quickstart/). Run
the commands in the folder that holds the package.

### 1. Catch a wrong value

A **kind** is a category of things. A **field** is a value that each thing
of the kind has. Each service has a tier, which is one of two words.

```yaml
# wrong-value/defs/services.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.
  fields:
    tier:
      enum: [critical, standard]
```

A **fact** states one real thing. This service has a tier that is not one of
the two words. The `import` block tells where the name `Service` comes from.

```yaml
# wrong-value/tests/facts.of.yaml
import:
  - from: file:../defs/services.of.yaml
    use: [Service]

service/checkout:
  is: Fact
  kind: Service
  fields:
    tier: important
```

The tool finds the wrong value, and gives the file, the line and the rule:

```console
$ of-lang run wrong-value wrong-value/tests/facts.of.yaml
wrong-value/tests/facts.of.yaml:9:5: tier: important is not one of critical, standard (FC4)
```

### 2. Find what breaks when the database is down

A **link** is a named relationship between two things. `chains: true` tells
that the link passes on: if A depends on B, and B depends on C, then A
depends on C.

```yaml
# depends-on/defs/services.of.yaml
Service:
  is: Kind
  definition: A program that runs in production and answers requests.

dependsOn:
  is: Link
  from: Service
  to: Service
  reads: "{from} depends on {to}"
  definition: The first service cannot work when the second service is down.
  chains: true
```

The facts state that checkout depends on orders, and that orders depends on
the database. No fact states that checkout depends on the database.

```yaml
# depends-on/tests/facts.of.yaml
import:
  - from: file:../defs/services.of.yaml
    use: [Service, dependsOn]

service/database:
  is: Fact
  kind: Service

service/orders:
  is: Fact
  kind: Service
  links:
    dependsOn: service/database

service/checkout:
  is: Fact
  kind: Service
  links:
    dependsOn: service/orders
```

The tool concludes it:

```console
$ of-lang run depends-on depends-on/tests/facts.of.yaml
derived:
  service/checkout depends on service/database  (chains)
rejected:
  nothing
```

### 3. Make sure that each service has an owner

A **constraint** is a rule that rejects facts. `match` finds each service.
`filter` keeps the services that have no `ownedBy` link to a team. Each
service that is left breaks the rule.

```yaml
# no-owner/defs/services.of.yaml
Team:
  is: Kind
  definition: A group of people who build and run services.

Service:
  is: Kind
  definition: A program that runs in production and answers requests.

ownedBy:
  is: Link
  from: Service
  to: Team
  reads: "{from} is owned by {to}"
  definition: The team that answers for the service.

serviceHasOwner:
  is: Constraint
  says: Each service must have an owner.
  match:
    - kind: { name: Service, as: $service }
  filter:
    - none:
        match:
          - link: ownedBy
            from: $service
            to:   Team
  message: "{$service} has no owner."
```

The facts state two services. Only one has an owner.

```yaml
# no-owner/tests/facts.of.yaml
import:
  - from: file:../defs/services.of.yaml
    use: [Team, Service, ownedBy]

team/payments:
  is: Fact
  kind: Team

service/checkout:
  is: Fact
  kind: Service
  links:
    ownedBy: team/payments

service/search:
  is: Fact
  kind: Service
```

The tool rejects the other one, with the message of the constraint:

```console
$ of-lang run no-owner no-owner/tests/facts.of.yaml
derived:
  nothing
rejected:
  serviceHasOwner: service/search has no owner.
```

## Where to go next

- [The introduction](docs/introduction/README.md) builds one larger model
  in eight chapters, with the kitchen of a restaurant as its example. Each
  chapter is a question, and ends with the answer of the tool.
- [The language by example](docs/by-example/README.md) has one short page
  for each thing in the language, with the text next to the code.
- [The language specification](language/README.md) tells what the language
  is, with its principles and its vocabulary.
  - [The cheat sheet](language/cheatsheet.md) has the whole language on one
    page.
  - [The reference](language/README.md#reference) has every rule.
- [The `of-lang` tool](cli/README.md) tells about each command and its
  output, and what is planned.
- [The full example package](examples/restaurant/) is the model that the
  introduction builds. It depends on
  [a second package](examples/suppliers/), the restaurant's suppliers.
- [Why a model](docs/why.md) tells which problems a model solves, who it is
  for, and why this is a language of its own.
- [The design of the language](docs/design.md) tells what the language
  leaves out, the roads not taken, and the open questions.
- [The changelog](CHANGELOG.md) and [how to contribute](CONTRIBUTING.md).

## License

Licensed under the [Apache License 2.0](LICENSE).
