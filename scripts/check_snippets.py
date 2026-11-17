# SPDX-License-Identifier: Apache-2.0
# /// script
# requires-python = ">=3.12"
# dependencies = ["pyyaml==6.0.2"]
# ///
"""Check that every YAML snippet in the language reference is valid.

Each snippet starts with a comment that names its file in the restaurant
example, such as ``# defs/recipe.of.yaml`` or ``# package.yaml``. A snippet can
hold several files, each starting with such a comment.

A statement that the restaurant already has must be exactly the real one, and
so must each manifest key that a snippet shows. For a statement with a pending
change, the version without the change must be the real one. A new statement is merged into
a fresh copy of ``examples/``, with the snippet's imports, and then
``of-lang check`` and ``of-lang test`` must pass on the restaurant package.

Usage: uv run scripts/check_snippets.py
"""

import re
import shutil
import subprocess
import sys
import tempfile
from collections.abc import Mapping, Sequence
from dataclasses import dataclass
from pathlib import Path
from typing import cast

import yaml

REPO = Path(__file__).resolve().parents[1]
EXAMPLES = REPO / "examples"
PACKAGE = "restaurant"
PAGES = (REPO / "language" / "reference",)
HEADER = re.compile(r"^# ((?:defs|tests)/\S+\.of\.ya?ml|package\.ya?ml)\s*$")
FENCE_OPEN = "```yaml"
FENCE_CLOSE = "```"

type Yaml = object


@dataclass(frozen=True)
class SnippetFile:
    """One file of a snippet: its path in the package and its text."""

    path: str
    text: str


@dataclass(frozen=True)
class Snippet:
    """One YAML block of a page, where it is, and the files it holds."""

    page: Path
    line: int
    files: tuple[SnippetFile, ...]


@dataclass(frozen=True)
class Problem:
    """Why one snippet is not valid."""

    snippet: Snippet
    text: str


def yaml_blocks(text: str) -> tuple[tuple[int, tuple[str, ...]], ...]:
    """Return each YAML block of a Markdown page with the line it starts on."""
    blocks: list[tuple[int, tuple[str, ...]]] = []
    current: list[str] | None = None
    start = 0
    for number, line in enumerate(text.splitlines(), start=1):
        if current is None:
            if line.strip() == FENCE_OPEN:
                current, start = [], number
            continue
        if line.strip() == FENCE_CLOSE:
            blocks.append((start, tuple(current)))
            current = None
            continue
        current.append(line)
    return tuple(blocks)


def split_files(lines: Sequence[str]) -> tuple[SnippetFile, ...] | None:
    """Split a block into its files. None when it does not start with a file comment."""
    first = next((line for line in lines if line.strip()), "")
    if not HEADER.match(first):
        return None
    files: list[SnippetFile] = []
    path = ""
    body: list[str] = []
    for line in lines:
        match = HEADER.match(line)
        if match:
            if path:
                files.append(SnippetFile(path, "\n".join(body)))
            path, body = match.group(1), []
            continue
        body.append(line)
    files.append(SnippetFile(path, "\n".join(body)))
    return tuple(files)


def find_snippets(page: Path) -> tuple[tuple[Snippet, ...], tuple[str, ...]]:
    """Return the snippets of a page, and the blocks that name no file."""
    snippets: list[Snippet] = []
    unnamed: list[str] = []
    for line, lines in yaml_blocks(page.read_text()):
        files = split_files(lines)
        if files is None:
            unnamed.append(
                f"{page.relative_to(REPO)}:{line}: the snippet does not start with a comment that names its file, such as # defs/recipe.of.yaml"
            )
            continue
        snippets.append(Snippet(page, line, files))
    return tuple(snippets), tuple(unnamed)


def merged_imports(old: Sequence[Yaml], new: Sequence[Yaml]) -> list[dict[str, Yaml]]:
    """Merge import blocks: entries with the same from join their use lists."""
    out = [dict(cast("Mapping[str, Yaml]", e)) for e in old]
    for entry in new:
        e = cast("Mapping[str, Yaml]", entry)
        same = next((o for o in out if o.get("from") == e.get("from")), None)
        if same is None:
            out.append(dict(e))
            continue
        uses = cast("list[Yaml]", same.get("use", []))
        same["use"] = uses + [
            u for u in cast("list[Yaml]", e.get("use", [])) if u not in uses
        ]
    return out


def merged_file(old: Mapping[str, Yaml], new: Mapping[str, Yaml]) -> dict[str, Yaml]:
    """Merge a file of statements: a statement replaces the one of the same name."""
    out = dict(old)
    if "import" in new:
        out["import"] = merged_imports(
            cast("list[Yaml]", old.get("import", [])), cast("list[Yaml]", new["import"])
        )
    out.update({k: v for k, v in new.items() if k != "import"})
    if "import" in out:
        out = {"import": out.pop("import"), **out}
    return out


def load_mapping(text: str) -> dict[str, Yaml]:
    """Parse a YAML file that holds a mapping, or nothing."""
    data = cast("Yaml", yaml.safe_load(text))
    return dict(cast("Mapping[str, Yaml]", data)) if isinstance(data, dict) else {}


def apply(package: Path, file: SnippetFile) -> None:
    """Write one snippet file into a package, merged with what is there."""
    target = package / file.path
    old = load_mapping(target.read_text()) if target.exists() else {}
    new = load_mapping(file.text)
    result: Yaml = (
        {**old, **new} if file.path.startswith("package.") else merged_file(old, new)
    )
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(yaml.safe_dump(result, sort_keys=False, allow_unicode=True))


def as_before_change(value: Yaml) -> Yaml:
    """Return a statement as it is without its pending change.

    A removed statement is the statement without its mark, and an updated
    statement is its ``was`` version. Any other value stays as it is.
    """
    if not isinstance(value, dict):
        return value
    statement = cast("Mapping[str, Yaml]", value)
    match statement.get("pending"):
        case "removed":
            return {k: v for k, v in statement.items() if k != "pending"}
        case "updated":
            return {
                "is": statement.get("is"),
                **cast("Mapping[str, Yaml]", statement.get("was", {})),
            }
        case _:
            return value


def differences(file: SnippetFile) -> tuple[str, ...]:
    """Name each statement or manifest key of a snippet file that differs from the real one."""
    target = EXAMPLES / PACKAGE / file.path
    if not target.exists():
        return ()
    old = load_mapping(target.read_text())
    new = load_mapping(file.text)
    return tuple(
        f"{file.path}: {key} is not the same as in examples/{PACKAGE}"
        for key, value in new.items()
        if key != "import" and key in old and old[key] != as_before_change(value)
    )


def run_tool(tool: Path, workdir: Path, command: str) -> tuple[bool, str]:
    """Run one of-lang command on the restaurant package in a workdir."""
    done = subprocess.run(  # noqa: S603 — runs the tool this script just built
        [str(tool), command, PACKAGE],
        cwd=workdir,
        capture_output=True,
        text=True,
        check=False,
    )
    return done.returncode == 0, (done.stdout + done.stderr).strip()


def check_snippet(tool: Path, snippet: Snippet) -> Problem | None:
    """Compare a snippet with the real package, then merge it into a fresh copy and check and test it."""
    diffs = [d for file in snippet.files for d in differences(file)]
    if diffs:
        return Problem(snippet, "\n".join(diffs))
    with tempfile.TemporaryDirectory() as tmp:
        work = Path(tmp) / "examples"
        shutil.copytree(EXAMPLES, work)
        for file in snippet.files:
            apply(work / PACKAGE, file)
        for command in ("check", "test"):
            ok, output = run_tool(tool, work, command)
            if not ok:
                return Problem(snippet, f"of-lang {command}:\n{output}")
    return None


def build_tool(into: Path) -> Path:
    """Build the of-lang tool from this repository."""
    tool = into / "of-lang"
    command = ["go", "build", "-o", str(tool), "./cmd/of-lang"]
    subprocess.run(command, cwd=REPO / "cli", check=True)  # noqa: S603 — runs go from PATH
    return tool


def main() -> int:
    """Check every snippet and print each problem."""
    snippets: list[Snippet] = []
    problems: list[str] = []
    for folder in PAGES:
        for page in sorted(folder.glob("*.md")):
            found, unnamed = find_snippets(page)
            snippets.extend(found)
            problems.extend(unnamed)
    with tempfile.TemporaryDirectory() as tmp:
        tool = build_tool(Path(tmp))
        for snippet in snippets:
            problem = check_snippet(tool, snippet)
            if problem is not None:
                where = (
                    f"{problem.snippet.page.relative_to(REPO)}:{problem.snippet.line}"
                )
                problems.append(f"{where}: {problem.text}")
    for p in problems:
        print(p, file=sys.stderr)
    print(f"{len(snippets)} snippets checked, {len(problems)} problems")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(main())
