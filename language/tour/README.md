---
title: Tour — The OpenFjord modeling language
---

# Tour

> Part of [the OpenFjord modeling language](../README.md).

Seven short steps, each adding one idea to one package. By the end you will
have written the [example package](../../examples/data-model/): teams own
databases, databases hold tables, applications run jobs, and jobs read and
write tables.

| Step | You learn | You add |
|---|---|---|
| [1. A kind](1-a-kind.md) | packages, files, kinds, fields, value types | `Team` |
| [2. Facts and a test](2-facts-and-a-test.md) | facts, examples, `of-lang test` | a test with one team |
| [3. Links](3-links.md) | links, counts, inverses, imports between files | `Database`, `Table` |
| [4. Rules](4-rules.md) | patterns, `derive`, `under`, `chains`, explanations | who a table belongs to; jobs; critical tables |
| [5. Constraints](5-constraints.md) | `require`, messages, rejections | applications, contracts, reading across applications |
| [6. Questions and docs](6-questions-and-docs.md) | questions, docs that embed definitions | who breaks if an application goes |
| [7. Many packages](7-many-packages.md) | namespaces, dependencies, exports, extension points, versions | a second package |

Each step ends with what the tool prints. The output follows the formats of
[the `of-lang` tool](../../cli/README.md), written for the language as this
draft specifies it; the tool itself still implements an earlier version.

When a step says "see KD3", that is a rule in [the reference](../reference/README.md).
You do not need the reference to follow the tour.
