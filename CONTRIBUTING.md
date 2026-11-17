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
with the part it touches in bold: `language`, `of-lang` or `examples`. The
entry becomes the release notes, so links in it must be full URLs. A new major
version starts a new file, and gets a line in
[`changelog/README.md`](changelog/README.md).

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
