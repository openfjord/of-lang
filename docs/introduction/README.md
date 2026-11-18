---
title: Introduction — The OpenFjord modeling language
---

# Introduction to the OpenFjord modeling language

> Part of [the OpenFjord modeling language](../../language/README.md).

The introduction teaches the language with one package: the kitchen of a
small restaurant. Each chapter is one question that the kitchen wants an
answer to. The chapter adds to the package only what the answer needs, and
then shows the answer that the tool gives. After the last chapter, the
package is the same as [the example package](../../examples/restaurant/).

Read the chapters in order. Each chapter uses only what the chapters before
it explain.

1. [Get an answer that nobody wrote down](01-margherita-and-nuts.md). Does
   the margherita contain nuts? The tool gives the answer, and each step
   that leads to it.
2. [Ask the model a question](02-dishes-with-nuts.md). Which dishes contain
   nuts? The answer follows each change to a recipe.
3. [Reject what breaks a rule](03-menu-that-hides-an-allergen.md). A menu
   hides an allergen. The tool rejects it, with a message that tells which
   warning is missing.
4. [Find what depends on what](04-sauce-station-closes.md). What stops when
   the sauce station closes? The tool finds each station that depends on
   it.
5. [Find what many things share](05-recipes-that-need-care.md). Which
   recipes need extra care? The tool finds the recipes that many dishes
   share.
6. [Check one value against another](06-food-kept-cold.md). Is the food kept
   cold enough? The tool rejects a batch in a storage zone that is too
   warm.
7. [Share a model between teams](07-share-with-purchasing.md). The kitchen
   and purchasing have two packages, and the tool checks them together.
8. [Change a rule safely](08-change-a-rule.md). The owner sees the kitchen
   with the change and without it, and then accepts or rejects the change.

## The package after each chapter

The folder [`steps/`](steps/) holds the package as it is after each chapter:
`steps/step-01` after chapter 1, `steps/step-02` after chapter 2, and so on. Each one
is a full package that the tool accepts. `steps/step-08` holds the same files as
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
