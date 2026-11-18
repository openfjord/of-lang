# Contributing

## Changelog

Every change that a user of the language or the tool would notice gets an
entry in the changelog for its major version, `changelog/vN.md`, under the
next version's heading:

```markdown
## [0.2.0] - 2026-12-01

### Added
- **of-lang**: ...

### Changed
- **language**: ...
```

Group entries under `Added`, `Changed`, `Fixed` and `Removed`, and start each
with the part it touches in bold: `language`, `of-lang`, `examples` or
`docs`. The
entry becomes the release notes, so links in it must be full URLs. A new major
version starts a new file, and gets a line in
[`changelog/README.md`](changelog/README.md).

## Snippets in the reference

Every YAML snippet in `language/reference/` is a file of the restaurant
example, or a part of one. It starts with a comment that names the file, such
as `# defs/recipe.of.yaml`. A statement that the restaurant has must be the
same in the snippet, and a new statement must keep the package valid. Check
them with:

```bash
uv run scripts/check_snippets.py
```

The snippets workflow (`.github/workflows/snippets.yml`) runs this check on
each push to `main` and on each pull request.

## By-example pages

Each page in `docs/by-example/` shows one thing in the language. The folder
with the name of the page is a package: `defaults.md` shows the package in
`defaults/`. The page is a table with one row for each step, with the text
in the left cell and the code in the right cell:

````markdown
<table>
<tr>
<td>

`status` has a default. A dish that gives no status is `active`.

</td>
<td>

```yaml
    status:
      enum: [active, retired]
      default: active
```

</td>
</tr>
</table>
````

The code cells of a page, joined in order, must be exactly the files of the
package. The first block of each file starts with a comment that names the
file, such as `# defs/dish.of.yaml`, and the blocks after it continue that
file. A page can leave out the manifest. The snippets check also runs
`of-lang check` and `of-lang test` on each package.

The top `README.md` shows the three packages in `examples/quickstart`. Each
of its YAML blocks starts with a comment that names a package and a file,
such as `# with-nuts/defs/menu.of.yaml`, and must be exactly that file. Each
console block holds one `of-lang` command and its output. The snippets check
runs each command in `examples/quickstart` and compares the output.

Add each new page to the list in `docs/by-example/README.md`, under the
heading of its group. The list also sets the place of the page in the
navigation of the site, and each heading is a group there.

## The introduction

The eight chapters in `docs/introduction/` grow the restaurant package one
result at a time. Each chapter is one question that the kitchen wants an
answer to, such as "Which dishes contain nuts?". Its title gives the result
in general words, such as "Ask a question": it names neither the
restaurant nor a part of the language. The chapter opens with the question
of the kitchen. The chapter adds to the package
only what the answer needs, explains each part of the language where the
answer first needs it, and shows the answer of the tool.

`steps/step-NN` holds the package as it is after chapter `NN`, and the last
one has exactly the files of `examples/restaurant`. An earlier package is
that package with less in it: a chapter adds files, statements, and keys to
statements, and a reader never meets a statement before the chapter that
explains it. A file of facts must hold a test, so the first test starts
small and gets its full form in chapter 4.

A chapter is prose with snippets between the paragraphs. Each YAML snippet
starts with a comment that names its file, such as `# defs/recipe.of.yaml`,
and shows whole statements. Each statement, import block and manifest key
that a snippet shows must be the same as in the package of its chapter. A
statement with a pending change follows the rule for the reference: the
version without the change must be the same. Tool output goes in a `console`
block, with the command and what it prints.

The snippets check compares each chapter with its package, compares the last
package with `examples/restaurant`, and runs `of-lang check` and `of-lang
test` on each package. The packages find `@bistro/suppliers` in `examples/`.

Add each new chapter to the list in `docs/introduction/README.md`, which
also sets its place in the navigation of the site.

## The site

The documentation site holds each Markdown file that starts with a block
that gives its `title`. Build it into `_site/` with:

```bash
uv run scripts/build_site.py
```

The build fails when a link points at a file or a heading that does not
exist, or when a page is missing from the navigation. The template, the
stylesheet and the theme script are in `scripts/site/`. The groups of the
navigation are at the top of the script.

Text that makes sense on GitHub but not on the site, such as the link to
the site, goes between a line `<!-- not on the site -->` and a line
`<!-- end -->`. GitHub shows it, and the site leaves it out.

A diagram is an SVG file whose top element has the class `d-root`, such as
`docs/images/of-lang-flow.svg`. Its colors are variables, and one block
gives the light colors for a reader whose system is set to light. GitHub
shows the file as an image. The site draws it into the page, where the same
block follows the light and dark switch of the site.

The site workflow (`.github/workflows/site.yml`) builds the site on each
pull request, and publishes it to GitHub Pages on each push to `main`.

## Icons of the vocabulary

On the site, each word of the [vocabulary](language/README.md#vocabulary)
has an icon of its own. The icon is in the vocabulary table, and next to the
title and the navigation entry of each page of the reference that is about
that one word. The
icons come from [Tabler Icons](https://tabler.io/icons), under the MIT
license (`scripts/site/tabler-icons.LICENSE`).

`ICONS` in `scripts/build_site.py` gives each word the name of its icon and
its code point in the font, and `PAGE_WORDS` gives each such page its word.
The build fails when a word of the vocabulary has no icon, or when two words
have the same icon.

`scripts/site/tabler-icons.woff2` holds only the icons in use. To add or
change an icon, make the file again from the full font of Tabler Icons, with
the code points in `ICONS`:

```bash
uvx --from fonttools --with brotli pyftsubset tabler-icons.woff2 \
  --unicodes=U+eaff,U+ebcc --flavor=woff2 \
  --output-file=scripts/site/tabler-icons.woff2
```

## Releasing

One version covers the whole repository. To release version `X.Y.Z`:

1. Make sure `changelog/vN.md` has a `## [X.Y.Z] - <date>` entry, with the
   release date, on `main`.
2. Tag that commit and push the tag:

   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

The release workflow (`.github/workflows/release.yml`) then:

- fails if the changelog has no entry for the version;
- checks the code (`gofmt`, `go vet`, `go test`);
- builds `of-lang` for macOS and Linux, on amd64 and arm64;
- adds the tag `cli/vX.Y.Z`, which Go needs because the module is in
  the `cli` folder;
- publishes a GitHub release with the binaries, their checksums, and the
  changelog entry as its notes.

Then update `Formula/of-lang.rb` in
[openfjord/homebrew-tap](https://github.com/openfjord/homebrew-tap): set `url`
to the new tag's source archive and `sha256` to its checksum.
