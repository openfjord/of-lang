# The OpenFjord modeling language

A way to write down the knowledge of your business: what its words mean, how
its things connect, and what follows from what. Its tool, `of-lang`, checks
that knowledge, reasons with it, and explains each conclusion. The result is a
model that people, code and AI agents can all use. The technical name for such
a model is an ontology.

<p align="center">
  <img src="docs/images/of-lang-flow.svg" alt="Knowledge and facts go into the OpenFjord modeling language, which produces a model." width="100%">
</p>

**Documentation: <https://lang.openfjord.io/>**

> **Version 0.2.0.** Before version 1.0, the language can change in ways
> that break a package.

## Install

```bash
brew install openfjord/tap/of-lang
```

Or download a binary for macOS or Linux from the
[releases](https://github.com/openfjord/of-lang/releases).

## An example

The definitions say what a service is, and that "depends on" passes on:

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

The facts give three services:

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

No fact says that checkout depends on the database. The tool concludes it:

```console
$ of-lang run depends-on depends-on/tests/facts.of.yaml
derived:
  service/checkout depends on service/database  (chains)
rejected:
  nothing
```

The package is in [`examples/quickstart/depends-on`](examples/quickstart/depends-on/).

## Read more

- [What it is and what it is good for](https://lang.openfjord.io/)
- [The introduction](https://lang.openfjord.io/docs/introduction/), in eight chapters
- [The language by example](https://lang.openfjord.io/docs/by-example/)
- [The reference](https://lang.openfjord.io/language/)
- [The `of-lang` tool](cli/README.md)
- [The changelog](CHANGELOG.md) and [how to contribute](CONTRIBUTING.md)

## License

Licensed under the [Apache License 2.0](LICENSE).
