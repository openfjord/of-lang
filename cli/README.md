---
title: of-lang — The OpenFjord modeling language
---

# `of-lang`

A command-line tool for [the OpenFjord modeling language](../language/README.md):
it checks packages, runs their derivatives over facts, explains what it derived,
compares two sets of facts, and runs a package's tests.

> **It lags the draft.** The tool implements an earlier version of the
> language: `links` and `where` instead of `match` and `filter`, variables
> named in `says`, kinds given only on a variable's first mention, `Rule`
> instead of `Derivative`, `require` instead of `unless`, `definitions/` instead of `defs/`, and `.yaml`
> instead of `.of.yaml`. It does not run the
> [example package](../examples/restaurant/) as it stands. It will be brought
> up to date once the draft settles.

## Build

```bash
cd cli && go build -o of-lang ./cmd/of-lang
```

## Commands

```
of-lang check <package-dir>...
of-lang run <package-dir> <facts.yaml>
of-lang explain <package-dir> <facts.yaml> <subject> <link-or-kind> [object]
of-lang diff <package-dir> <before.yaml> <after.yaml>
of-lang test <package-dir>
```

| Command | Does |
|---|---|
| `check` | checks definitions and docs of one or more packages |
| `run` | derives what follows from a facts file, then checks constraints |
| `explain` | shows why one fact holds: a link (`subject link object`) or a kind (`subject Kind`) |
| `diff` | shows what changes between two facts files: stated, derived and rejected |
| `test` | runs every example and question in the package's `tests/` |

| Flag | Means |
|---|---|
| `--packages <dir>` | where to look for dependencies; repeatable; by default the folder that holds the package |
| `--format text\|json` | the output format; `text` by default |

| Exit code | Means |
|---|---|
| 0 | everything passes |
| 1 | problems found: a check failed, a constraint rejected something, or a test failed |
| 2 | the tool could not run: a bad command line, or a file it cannot read |

## Output

All examples below run in the folder that holds the package.

### `check`

```
$ of-lang check restaurant
ok
```

Each problem is one line, with the file, line and column:

```
$ of-lang check restaurant
restaurant/defs/recipe.yaml:59:3: reads: of madeAtStation must use both {from} and {to}
```

The draft has each error name the rule it enforces, such as `(LK2)`; the tool
does not yet.

### `run`

Every derived fact as its `reads` sentence, with what derived it: a derivative, or
`under` or `chains` on a link. Then every rejection, with its message.

```
$ of-lang run restaurant restaurant/tests/examples/allergen-through-a-sauce.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
  station/pizza depends on station/sauce  (rule stationDependsOnStation)
  recipe/pesto contains allergen/nuts  (rule allergenFromIngredient)
  dish/margherita contains allergen/milk  (rule allergenFromIngredient)
  dish/margherita contains allergen/nuts  (rule allergenFromRecipe)
rejected:
  nothing
```

```
$ of-lang run restaurant restaurant/tests/examples/dish-without-allergen-warning.yaml
derived:
  ...
rejected:
  dishListsItsAllergens: The menu must warn that dish/margherita contains allergen/nuts.
```

A rejection makes `run` exit with 1.

### `explain`

The fact, the derivative's `says`, what each variable stood for, and the facts
the match used. A derived fact among them unfolds into its own explanation,
indented below it.

```
$ of-lang explain restaurant restaurant/tests/examples/allergen-through-a-sauce.yaml dish/margherita hasAllergen allergen/nuts
dish/margherita contains allergen/nuts  (rule allergenFromRecipe)
  A recipe contains each allergen that one of its recipes contains.
  with $allergen = allergen/nuts, $part = recipe/pesto, $recipe = dish/margherita
  because:
    dish/margherita uses recipe/pesto  (written at restaurant/tests/examples/allergen-through-a-sauce.yaml:99:3)
    recipe/pesto contains allergen/nuts  (rule allergenFromIngredient)
      A recipe contains each allergen that one of its ingredients contains.
      with $allergen = allergen/nuts, $ingredient = ingredient/pine_nuts, $recipe = recipe/pesto
      because:
        recipe/pesto uses ingredient/pine_nuts  (written at restaurant/tests/examples/allergen-through-a-sauce.yaml:84:3)
        ingredient/pine_nuts contains allergen/nuts  (written at restaurant/tests/examples/allergen-through-a-sauce.yaml:66:3)
```

```
$ of-lang explain restaurant restaurant/tests/examples/base-recipe.yaml recipe/tomato_sauce BaseRecipe
recipe/tomato_sauce is a BaseRecipe  (rule baseRecipe)
  A recipe that more than two other recipes use, directly or through other recipes, is a base recipe.
  with $recipe = recipe/tomato_sauce
```

### `diff`

What changes from one facts file to the other, in three parts: facts stated,
facts derived, and rejections. `+` is added, `-` is gone. A field value is a
fact of its own, so a changed value shows as one `-` and one `+`.

In this example, a copy of the test file adds the pesto to the marinara and
increases its price. The file is `marinara-with-pesto.yaml`.

```
$ of-lang diff restaurant restaurant/tests/examples/allergen-through-a-sauce.yaml marinara-with-pesto.yaml
asserted:
  + dish/marinara has price 11
  + dish/marinara uses recipe/pesto
  - dish/marinara has price 10
derived:
  + dish/marinara contains allergen/nuts
rejected:
  + dishListsItsAllergens: The menu must warn that dish/marinara contains allergen/nuts.
```

`diff` always exits with 0: a difference is not a problem.

### `test`

One line per test or question, then a total. A failure says what was
expected and did not happen.

```
$ of-lang test restaurant
FAIL  test allergenThroughASauce  (restaurant/tests/examples/allergen-through-a-sauce.yaml)
      expected derived, but not: dish/margherita contains allergen/milk
pass  test tomatoSauceInThreeDishes  (restaurant/tests/examples/base-recipe.yaml)
pass  test dishWithoutAllergenWarning  (restaurant/tests/examples/dish-without-allergen-warning.yaml)
pass  question whichDishesContainNuts  (restaurant/tests/questions/which-dishes-contain-nuts.yaml)
3 passed, 1 failed
```

### JSON

`--format json` prints the same results as one JSON object, for other tools
to read. A derived link, for example:

```json
{
  "type": "link",
  "sentence": "dish/margherita contains allergen/nuts",
  "from": "dish/margherita",
  "link": "hasAllergen",
  "to": "allergen/nuts"
}
```

and a rejection carries the constraint, its message, and what each variable
stood for:

```json
{
  "constraint": "dishListsItsAllergens",
  "message": "The menu must warn that dish/margherita contains allergen/nuts.",
  "bindings": [
    { "var": "$dish", "id": "dish/margherita" },
    { "var": "$allergen", "id": "allergen/nuts" }
  ]
}
```

### Derived facts

When the tool writes facts out, a derived fact is a fact with
`derived: true`; `derived` is `false` by default. The flag exists only in
the tool's output: the language has no such key, so a file that writes it is
rejected like any unknown key.

A thing that is partly stated and partly derived has two facts under one id:
the stated one, as written, and one derived fact beside it holding only what
was derived, never repeating the stated one. Its `kind` is the kind a
derivative made the thing, or the stated kind if no derivative did.

```yaml
recipe/tomato_sauce:
  is: Fact
  kind: BaseRecipe                 # a derivative made it a base recipe
  derived: true

recipe/pesto:
  is: Fact
  kind: Recipe                     # the stated kind; no derivative changed it
  derived: true
  hasAllergen: [allergen/nuts]     # from allergenFromIngredient
```

The tool does not write facts in this form yet.

## Planned outputs

Not built yet; the draft describes them so the language can be judged by what
it produces.

### The change report

Every change to a package returns one report, for the facts in its tests and
examples, before it lands. It joins `diff` and `explain`: what was stated,
what newly follows and why, what no longer follows, and what is rejected.

```
Stated (2)
  + dish/margherita usesRecipe recipe/pesto
  + dish/marinara usesRecipe recipe/pesto

Derived (3)
  + station/pizza dependsOnStation station/sauce
      derivative  stationDependsOnStation
      match   $station = station/pizza, $recipe = dish/margherita,
              $part = recipe/pesto, $other = station/sauce
      links   station/pizza makesRecipe dish/margherita
              dish/margherita usesRecipe recipe/pesto
              recipe/pesto madeAtStation station/sauce
      filter  station/sauce does not equal station/pizza
  + dish/margherita hasAllergen allergen/nuts
      derivative  allergenFromRecipe
  + (1 more, grouped: hasAllergen)

No longer derived (0)

Rejected (2)
  ✗ dishListsItsAllergens
    The menu must warn that dish/margherita contains allergen/nuts.
  ✗ dishListsItsAllergens
    The menu must warn that dish/marinara contains allergen/nuts.
```

Large results are grouped by link, with new kinds of conclusion and
rejections listed first.

### Product's pages

A page per definition for product: the definition, examples and
non-examples, every count as a sentence, and every constraint and derivative as its
`says`, shown as written.

> **Recipe** — The steps that make one food from ingredients and other recipes.
> - Every recipe is made at exactly one station.
> - A recipe contains each allergen that one of its ingredients contains.
>   *(derivative)*
> - The menu must give a warning about each allergen in an active dish.
>   *(constraint)*
>
> Example: the tomato sauce. Not a recipe: a new dish that a cook tries one
> time and does not keep.

A change to a product-owned package shows the sentences that changed, for
approval on that page.

### The RDF export

RDF (Resource Description Framework), OWL (Web Ontology Language) and SHACL
(Shapes Constraint Language) are generated for exchange and for existing
tools. Nobody writes them.

| In the language | Exported as |
|---|---|
| package | an ontology with a version and `owl:imports`; its statements in their own named graph |
| kind | `owl:Class`, plus a closed `sh:NodeShape` for its fields and links |
| `extends` | `rdfs:subClassOf` |
| field | `owl:DatatypeProperty`, checked by a `sh:property` built from its JSON Schema keywords |
| value type | `rdfs:Datatype` with restrictions (`xsd:pattern`, `xsd:maxLength`, `xsd:minInclusive` and the like) |
| `string`, `integer`, `number`, `boolean`, `date`, `date-time`, `duration` | `xsd:string`, `xsd:integer`, `xsd:decimal`, `xsd:boolean`, `xsd:date`, `xsd:dateTime`, `xsd:duration` |
| link | `owl:ObjectProperty`; `from` and `to` as SHACL checks, never `rdfs:domain` or `rdfs:range` |
| `under` | `rdfs:subPropertyOf` |
| `inverse` | `owl:inverseOf` |
| `chains` | `owl:TransitiveProperty` |
| counts | `sh:minCount`, `sh:maxCount`, never an OWL cardinality restriction |
| constraint | a SHACL shape with a SPARQL (SPARQL Protocol and RDF Query Language) check |
| derivative | a SHACL rule |
| definition, `reads` | `rdfs:comment`, and the generated sentence as a label for readers |
| facts | RDF triples in the package's named graph |
| derived facts | RDF triples in a derived graph, with PROV-O (the W3C provenance vocabulary) linking each to its derivative |

Each IRI (Internationalized Resource Identifier) is built from the package's
`rdfNamespace` and the definition's name.

`--tags` exports only the definitions carrying at least one of the named
[tags](../language/reference/tag.md); without it, the whole package is
exported. A reference to a definition that is not exported is left out,
and the export lists each one:

```
$ of-lang export --tags menu examples/restaurant
left out: link listsAllergen (Dish to Allergen): Allergen is not exported
left out: Dish extends Recipe: Recipe is not exported
```

## Planned: a language server

`of-lang lsp` will be a language server, speaking LSP (the Language Server
Protocol) over standard input and output, so any editor that speaks LSP gets
the language with one command to launch. It is one server, depending on no
other language server, including the YAML one.

### Two layers in one server

- **Structure, from a schema.** One JSON Schema per statement kind, chosen by
  `is:`, plus one for `package.yaml`, hand-written and embedded in the
  binary. The schemas are the source of truth for the shape of statements:
  every key, its description, and the fixed words it allows (`exactly 1`,
  `is at most`, `becomes`). The field keywords are part of them, so a field
  definition is checked and completed like any other part of a statement.
  The checker's structural errors (unknown keys, wrong shapes, words not in
  the list) are this validation, mapped back to a line and column.
- **Meaning, from the checker.** Everything that needs names (imports,
  kinds, links, fields, variables, fact ids, docs links) comes from the same
  loader and checker as `of-lang check`.

```
cli/internal/schema/
  package.schema.json
  kind.schema.json
  link.schema.json
  valuetype.schema.json
  constraint.schema.json
  derivative.schema.json
  fact.schema.json
  test.schema.json
  question.schema.json
```

The schemas can also be published for other tools to use; the server never
needs them to.

### What it offers

| Feature | What it does |
|---|---|
| Diagnostics | the check's errors as you type, each with its rule code, linked to that rule in the reference |
| Go to definition | from any name to where it is defined: names in `use:`, `kind:`, `link:`, `extends`, `under`, `type:`, `{$x: Kind}`, fact ids, docs links such as `recipe.of.yaml#Recipe`, and a `$variable` to its first mention |
| Find references | where a definition is used, including the docs that link to it or embed it |
| Hover | a key's description from the schema; a definition's `definition`; a link's `reads`; a variable's kind; on a fact, what derivatives add to it |
| Completion | the keys and words the schema allows at that place; after `link:`, the links the file imports; in `{$x: …}`, only kinds that fit that end of the link; after `field:`, only fields of that node's kind; a name not yet imported, with its import |
| Code actions | add a missing import; rename a definition and the imports that use it, warning when the name is exported |
| Inlay hints | in test files, the derived facts beside the facts they are about |
| Code lens | run a `Test` or `Question` from above it |

### What it needs from the tool

1. **Positions for every name,** where it is defined and wherever it is used.
2. **A loader that keeps going after errors,** giving a partial model of a
   broken file, since a file being typed is broken most of the time.
3. **Unsaved files,** read from the editor instead of the disk.
4. **Re-checking the package** after each pause in typing; packages are
   small enough that incremental checking can wait.

### Order

1. The schemas, with structural diagnostics, key completion and hover.
2. Diagnostics from the checker, go to definition, hover on names.
3. Completion of names, adding imports, rename.
4. Inlay hints for derived facts, and running tests from the editor.

## Planned: a user interface

`of-lang ui` will start a local server that serves a web app. The web app
edits packages. Each change goes to the files of the package, in whichever
form those files are. The rest of the design is not yet discussed.
