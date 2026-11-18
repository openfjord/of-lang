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

Add each new page to the list in `docs/by-example/README.md`, which also
sets its place in the navigation of the site.

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
