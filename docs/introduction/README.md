---
title: Introduction — The OpenFjord modeling language
---

# Introduction to the OpenFjord modeling language

> Part of [the OpenFjord modeling language](../../language/README.md).

The introduction teaches the language with one package: the kitchen of a
small restaurant. The package starts with a manifest and five kinds. Each
chapter tells what the kitchen needs next, and adds the statements that give
it. After the last chapter, the package is the same as
[the example package](../../examples/restaurant/).

Read the chapters in order. Each chapter uses only what the chapters before
it explain.

1. [Kinds](01-kinds.md): the things of the kitchen, their fields, and the
   files that hold them.
2. [Links](02-links.md): how the things relate to each other.
3. [Facts and tests](03-facts-and-tests.md): one real dish, and a test that
   proves what the package tells about it.
4. [Derivatives](04-derivatives.md): the facts that follow from other facts,
   and questions about them.
5. [Constraints](05-constraints.md): what the kitchen does not allow.
6. [Packages](06-packages.md): a dependency, exports, tags and docs.
7. [Pending changes](07-pending-changes.md): a change that nobody has
   accepted yet.

## The package after each chapter

The folder [`steps/`](steps/) holds the package as it is after each chapter:
`steps/step-01` after chapter 1, `steps/step-02` after chapter 2, and so on. Each one
is a full package that the tool accepts. `steps/step-07` holds the same files as
the example package.

The chapters show the tool, `of-lang`, at work. To run the same commands,
[install the tool](../../cli/README.md#install), get a clone of the
repository, and go to the `steps` folder:

```console
$ cd docs/introduction/steps
$ of-lang check step-01
ok
```

## Other ways to learn

- [The language by example](../by-example/README.md) has one short page for
  each thing in the language.
- [The cheat sheet](../../language/cheatsheet.md) shows all of the syntax on
  one page.
- [The reference](../../language/README.md#reference) gives each rule.
