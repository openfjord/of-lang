---
title: of-lang — The OpenFjord modeling language
---

# `of-lang`

A command-line tool for [the OpenFjord modeling language](../../language/README.md):
it checks packages, runs their derivatives over facts, explains what it derived,
compares two sets of facts, and runs a package's tests.

> **It lags the draft.** The tool implements an earlier version of the
> language: `links` and `where` instead of `match` and `filter`, variables
> named in `says`, kinds given on a variable's first mention instead of in kind items, `Rule`
> instead of `Derivative`, `require` instead of `unless`, `definitions/` instead of `defs/`, and `.yaml`
> instead of `.of.yaml`. It does not run the
> [example package](../../examples/restaurant/) as it stands. It will be brought
> up to date once the draft settles.

## Build

```bash
go build -o of-lang ./cmd/of-lang
```

## Commands

```
of-lang check <package-dir>...
of-lang run <package-dir> <facts-file>
of-lang explain <package-dir> <facts-file> <subject> <link-or-kind> [object]
of-lang diff <package-dir> <before-file> <after-file>
of-lang test <package-dir>
```

| Command | Does |
|---|---|
| `check` | checks definitions and docs of one or more packages |
| `run` | derives what follows from a facts file, then checks constraints |
| `explain` | shows why one fact holds: a link (`subject link object`) or a kind (`subject Kind`) |
| `diff` | shows what changes between two facts files: stated, derived and rejected |
| `test` | runs every test and question in the package's `tests/`, including tests of derived terms |

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
restaurant/defs/recipe.of.yaml:49:3: the reads of madeAtStation must use both {from} and {to} (LK2)
```

The draft has each error name the rule it enforces, such as `(LK2)`; the tool
does not yet.

### `run`

Every derived fact as its `reads` sentence, with what derived it: a derivative, or
`under` or `chains` on a link. Then every rejection, with its message.

```
$ of-lang run restaurant restaurant/tests/examples/allergen-through-a-sauce.of.yaml
derived:
  dish/margherita uses ingredient/parmesan  (under)
  station/pizza depends on station/sauce  (derivative stationDependsOnStation)
  recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
  dish/margherita contains allergen/milk  (derivative allergenFromIngredient)
  dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
rejected:
  nothing
```

```
$ of-lang run restaurant restaurant/tests/examples/dish-without-allergen-warning.of.yaml
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
$ of-lang explain restaurant restaurant/tests/examples/allergen-through-a-sauce.of.yaml dish/margherita hasAllergen allergen/nuts
dish/margherita contains allergen/nuts  (derivative allergenFromRecipe)
  A recipe contains each allergen that one of its recipes contains.
  with $allergen = allergen/nuts, $part = recipe/pesto, $recipe = dish/margherita
  because:
    dish/margherita uses recipe/pesto  (written at restaurant/tests/examples/allergen-through-a-sauce.of.yaml:114:5)
    recipe/pesto contains allergen/nuts  (derivative allergenFromIngredient)
      A recipe contains each allergen that one of its ingredients contains.
      with $allergen = allergen/nuts, $ingredient = ingredient/pine_nuts, $recipe = recipe/pesto
      because:
        recipe/pesto uses ingredient/pine_nuts  (written at restaurant/tests/examples/allergen-through-a-sauce.of.yaml:95:5)
        ingredient/pine_nuts contains allergen/nuts  (written at restaurant/tests/examples/allergen-through-a-sauce.of.yaml:72:5)
```

```
$ of-lang explain restaurant restaurant/tests/examples/base-recipe.of.yaml recipe/tomato_sauce BaseRecipe
recipe/tomato_sauce is a BaseRecipe  (derivative baseRecipe)
  A recipe that more than two recipes use, directly or through other recipes, is a base recipe.
  with $recipe = recipe/tomato_sauce
```

### `diff`

What changes from one facts file to the other, in three parts: facts stated,
facts derived, and rejections. `+` is added, `-` is gone. A field value is a
fact of its own, so a changed value shows as one `-` and one `+`.

In this example, a copy of the test file adds the pesto to the marinara and
increases its price. The file is `marinara-with-pesto.of.yaml`.

```
$ of-lang diff restaurant restaurant/tests/examples/allergen-through-a-sauce.of.yaml marinara-with-pesto.of.yaml
stated:
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
expected and did not happen. In this example, somebody removed the parmesan
garnish from the margherita in `allergen-through-a-sauce.of.yaml`:

```
$ of-lang test restaurant
FAIL  test allergenThroughASauce  (restaurant/tests/examples/allergen-through-a-sauce.of.yaml)
      expected derived, but not: dish/margherita contains allergen/milk
pass  test tomatoSauceInThreeDishes  (restaurant/tests/examples/base-recipe.of.yaml)
pass  test dishWithoutAllergenWarning  (restaurant/tests/examples/dish-without-allergen-warning.of.yaml)
pass  test garnishIsAnIngredient  (restaurant/tests/terms/garnish.of.yaml)
pass  question whichDishesContainNuts  (restaurant/tests/questions/which-dishes-contain-nuts.of.yaml)
4 passed, 1 failed
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
  links:
    hasAllergen: [allergen/nuts]   # from allergenFromIngredient
```

The tool does not write facts in this form yet.

## Planned outputs

Not built yet; the draft describes them so the language can be judged by what
it produces.

### The change report

The [pending change](../../language/reference/pending.md) of a package returns
one report, for the facts in its tests and examples. The report compares the
package without the change and with the change, so git is not necessary. It joins `diff` and `explain`: what was stated,
what newly follows and why, what no longer follows, and what is rejected.

```
Stated (2)
  + dish/margherita uses recipe/pesto
  + dish/marinara uses recipe/pesto

Derived (3)
  + station/pizza depends on station/sauce
      derivative  stationDependsOnStation
      match   $station = station/pizza, $recipe = dish/margherita,
              $part = recipe/pesto, $other = station/sauce
      links   dish/margherita is made at station/pizza
              dish/margherita uses recipe/pesto
              recipe/pesto is made at station/sauce
      filter  station/sauce does not equal station/pizza
  + dish/margherita contains allergen/nuts
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
non-examples, and every constraint and derivative as its
`says`, shown as written.

> **Recipe** — The steps that make one food from ingredients and other recipes.
> - Cooks make each recipe at exactly one station.
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

The export does not translate each statement into one RDF item. One
statement can give many RDF items, or none, and several statements can give
one. The export holds the same knowledge as the package and its
dependencies. The check for this: each question and each test of the
package gets the same result from the exported RDF, with the same example
facts, as from the package.

The export keeps these properties:

- The `from` and `to` of a link are checks. They never give a thing a kind.
- Constraints reject. They never add facts.
- The statements of each package are in their own named graph.
- Each derived fact keeps the reason that it is true.

Each IRI (Internationalized Resource Identifier) is built from the package's
`rdfNamespace` and the definition's name.

`of-lang export <package-dir> [--tags a,b]` writes the export. `--tags`
exports only the definitions carrying at least one of the named
[tags](../../language/reference/tag.md); without it, the whole package is
exported. A tagged subset is always complete, because each definition that
a tagged definition refers to has the same tag
([TG5](../../language/reference/tag.md#rules)). `of-lang check` rejects a
package that breaks this rule. For example, if `Dish` had the tag
`allergens` and `Recipe` did not:

```
$ of-lang check restaurant
restaurant/defs/recipe.of.yaml:25:3: Dish has the tag allergens and extends Recipe, which does not have it (TG5)
```

## Planned: a language server

`of-lang lsp` will be a language server, speaking LSP (the Language Server
Protocol) over standard input and output, so any editor that speaks LSP gets
the language with one command to launch. It is one server, depending on no
other language server, including the YAML one.

### Two layers in one server

- **Structure, from a schema.** One JSON Schema per statement type, chosen by
  `is:`, plus one for the manifest (`package.yaml` or `package.yml`), hand-written and embedded in the
  binary. The schemas are the source of truth for the shape of statements:
  every key, its description, and the fixed words it allows (`is at most`,
  `becomes`). The field keywords are part of them, so a field
  definition is checked and completed like any other part of a statement.
  The checker's structural errors (unknown keys, wrong shapes, words not in
  the list) are this validation, mapped back to a line and column.
- **Meaning, from the checker.** Everything that needs names (imports,
  kinds, links, fields, variables, fact ids, docs links) comes from the same
  loader and checker as `of-lang check`.

```
internal/schema/
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
| Go to definition | from any name to where it is defined: names in `use:`, `kind:`, `link:`, `extends`, `under`, `type:`, `kind: { name: … }`, fact ids, docs links such as `recipe.of.yaml#Recipe`, and a `$variable` to its kind item |
| Find references | where a definition is used, including the docs that link to it or embed it |
| Hover | a key's description from the schema; a definition's `definition`; a link's `reads`; a variable's kind; on a fact, what derivatives add to it |
| Completion | the keys and words the schema allows at that place; after `link:`, the links the file imports; in a link item's `from` and `to`, only variables and kinds that fit that end of the link; in `kind: { name: … }`, only kinds; after `field:`, only fields of that node's kind; a name not yet imported, with its import |
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
