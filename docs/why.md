---
title: Why a model — The OpenFjord modeling language
---

# Why a model

## The problem

Every software product has things, and rules about those things. Its things
are services, databases, tables, jobs and the teams that own them. Some of
its rules are:

- A service that depends on a second service also depends on everything
  that the second service depends on.
- Each service must have a team that owns it.
- A table that many services read needs extra care when it changes.

Usually, nobody writes these things and rules down in one place. They are in
the heads of the engineers, in wiki pages, in diagrams and in the code. This
causes problems:

- **The same word has different meanings.** One team says "customer" for a
  person who has an account. A different team uses the same word for a
  company that pays. Nobody sees the difference until a report gives a wrong
  number.
- **Nothing checks a rule that is only text.** A rule in a wiki page does
  not stop a change that breaks it. An engineer adds a service without an
  owner, and nothing objects. People find the problem later, during an
  incident, when nobody answers.
- **Nobody sees the effect of a change.** The orders database goes down for
  maintenance. Which services stop? To get the answer, somebody must read
  the code and the diagrams by hand.
- **The people who own the meaning cannot read where it is written.** A
  product manager decides what an "active customer" is. But the decision
  ends up in a query in the code, which the product manager does not read.
- **Each tool keeps its own copy.** The service catalog, the data catalog
  and the monitoring system each describe the same services in their own
  way. The copies do not agree.
- **AI agents guess.** A developer asks an AI agent: "Is it safe to remove
  this column?" The agent reads the code, finds no exact rules, and guesses.
  A wrong guess looks the same as a correct answer.

## What a model gives

A model writes these things down one time, in files next to the code. It
contains four types of information:

| Information | Example |
|---|---|
| the kinds of things, and their fields | a service has a name and a tier |
| the links between things | a service depends on a service; a team owns a service |
| what follows from other facts | if checkout depends on orders, and orders depends on the database, then checkout depends on the database |
| what is not allowed | a service must not be without an owner |

Then a tool can do things that text in a wiki page cannot do:

- **Check.** The tool compares data with the rules, and rejects data that
  breaks a rule. It tells which rule and why.
- **Derive.** The tool finds facts that nobody wrote. Nobody wrote
  "checkout depends on the database". The tool finds it from the two
  dependencies that somebody did write, and shows the steps.
- **Answer questions.** "Which services stop when the database is down?" The tool answers from the
  model, not from a search by hand.
- **Show the effect of a change.** Before a change is merged, the tool shows
  the new facts that follow, the facts that no longer follow, and the rules
  that the change breaks.
- **Give the owners of the meaning a page to read.** The tool makes a page for
  each term, in sentences and examples. The owners review the page, not the
  files.
- **Give other tools one source.** The tool exports the model in standard
  formats. A service catalog or a monitoring system reads the export, so there is one
  source, not many copies.
- **Give AI agents exact knowledge.** An AI agent reads the model to learn
  what a term means and which rules apply. It asks the tool questions, and
  gets answers with their steps. The tool checks each change that the agent
  makes, and rejects a change that breaks a rule.

In `of-lang` 0.2.0, `check`, `run` and `explain` check and derive, `test`
answers the questions that a package writes as tests, and `diff` shows what
a change to the facts does. The report for a change to the model, the pages
and the export are planned; [the tool's page](../cli/README.md) describes them.

A database schema or a JSON Schema does not do this. A database schema tells
how one system stores its data. A JSON Schema checks one document. Neither
describes things across many systems, the links between those things, or
what follows from the links.

## Who it is for

- **Developers** who write the model in files next to their code, review it
  like code, and check it in CI (continuous integration).
- **Product people** who own what the terms mean. Each definition is a
  plain sentence with examples, and each link reads as a sentence, so they
  need no background in logic.
- **The tools and AI agents** that need one exact source for what a term
  means and which rules apply.

No one needs to know RDF, OWL or any other ontology standard to use it.

## Why a new language

A model of this type has a name: an **ontology**. Standards for ontologies
already exist. The W3C (World Wide Web Consortium) published them, and many
tools read them:

- **RDF** (Resource Description Framework) stores each fact as three parts:
  a thing, a link and a different thing. "checkout depends on orders" is one fact.
- **RDFS** (RDF Schema) and **OWL** (Web Ontology Language) describe the
  kinds of things and the links, and what follows from them.
- **SHACL** (Shapes Constraint Language) describes what is not allowed.

These standards can describe a product exactly. But few product people or
developers read or write them, for eighteen reasons that
[the design of the language](design.md#why-ontologies-are-hard-for-software-people)
lists. This language is one that software people can read and write. It
keeps what the standards do well, inference above all, and its tool will
make the standard formats from it.
