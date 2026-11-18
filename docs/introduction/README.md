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

1. [Does the margherita contain nuts?](01-margherita-and-nuts.md) The tool
   gives the answer, and each step that leads to it.
2. [Which dishes contain nuts?](02-dishes-with-nuts.md) The kitchen asks the
   model a question, and the answer follows each change to a recipe.
3. [Stop a menu that hides an allergen](03-menu-that-hides-an-allergen.md).
   The tool rejects the menu, with a message that tells which warning is
   missing.
4. [What stops when the sauce station closes?](04-sauce-station-closes.md)
   The tool finds each station that depends on the closed station.
5. [Which recipes need extra care?](05-recipes-that-need-care.md) The tool
   finds the recipes that many dishes share.
6. [Is the food kept cold enough?](06-food-kept-cold.md) The tool rejects a
   batch in a storage zone that is too warm.
7. [Share the model with purchasing](07-share-with-purchasing.md). Two
   teams have two packages, and the tool checks them together.
8. [Change a rule without breaking the kitchen](08-change-a-rule.md). The
   owner sees the kitchen with the change and without it, and then accepts
   or rejects the change.

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
