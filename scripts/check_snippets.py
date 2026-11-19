# SPDX-License-Identifier: Apache-2.0
# /// script
# requires-python = ">=3.12"
# dependencies = ["pyyaml==6.0.2"]
# ///
"""Check that every YAML snippet in the reference, the by-example pages and the introduction is valid.

Each snippet starts with a comment that names its file in the platform
example, such as ``# defs/service.of.yaml`` or ``# package.yaml``. A snippet can
hold several files, each starting with such a comment.

A statement that the platform package already has must be exactly the real one, and
so must each manifest key that a snippet shows. For a statement with a pending
change, the version without the change must be the real one. A new statement is merged into
a fresh copy of ``examples/``, with the snippet's imports, and then
``of-lang check`` and ``of-lang test`` must pass on the platform package.

Each folder in ``docs/by-example`` is a package, and the page next to it with
the same name shows it. The YAML blocks of the page, joined in page order,
must be exactly the files of the package, and ``of-lang check`` and ``of-lang
test`` must pass on it. There, only the first block of a file starts with the
comment that names it. A page can leave out the manifest.

``README.md`` shows the packages in ``examples/quickstart``. Each of its YAML
blocks starts with a comment that names a package and a file, such as
``# depends-on/defs/services.of.yaml``, and must be exactly that file. Each of its
console blocks holds one ``of-lang`` command and its output, and the command
must give that output when it runs in ``examples/quickstart``.

Each folder in ``docs/introduction/steps`` is the restaurant package as it is
after the chapter of the introduction with the same number: ``step-01`` after
chapter 01. ``of-lang check``
and ``of-lang test`` must pass on each one, with ``examples/`` as the place of
its dependencies, and the last one must have exactly the files of
``examples/restaurant``. Each snippet of a chapter starts with a comment that
names its file, and each statement, import block and manifest key that it shows
must be exactly the one of the chapter's package. A statement with a pending
change is the exception: the version without the change must be the real one,
and an added statement must be absent. Such a chapter's snippets are merged
into a fresh copy of its package, which must then pass the two commands.

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
PACKAGE = "platform"
FINISHED = "restaurant"
PAGES = (REPO / "language" / "reference",)
BY_EXAMPLE = REPO / "docs" / "by-example"
INTRODUCTION = REPO / "docs" / "introduction"
STEPS = INTRODUCTION / "steps"
STEP_PREFIX = "step-"
QUICKSTART = EXAMPLES / "quickstart"
QUICKSTART_PAGES = (REPO / "README.md", REPO / "docs" / "home.md")
QUICKSTART_HEADER = re.compile(
    r"^# ([a-z0-9-]+/(?:(?:defs|tests)/\S+\.of\.ya?ml|package\.ya?ml))\s*$"
)
CONSOLE_OPEN = "```console"
PROMPT = "$ of-lang "
CHAPTERS = "[0-9][0-9]-*.md"
PACKAGE_FILES = ("package.y*ml", "defs/**/*.of.*", "tests/**/*.of.*")
MANIFESTS = frozenset({"package.yaml", "package.yml"})
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


def has_pending(value: Yaml) -> bool:
    """Tell if a value is a statement that carries a pending change."""
    return isinstance(value, dict) and "pending" in value


def is_added(value: Yaml) -> bool:
    """Tell if a value is a statement that a pending change adds."""
    return (
        isinstance(value, dict)
        and cast("Mapping[str, Yaml]", value).get("pending") == "added"
    )


def differences(
    file: SnippetFile, package: Path = EXAMPLES / PACKAGE, *, whole: bool = False
) -> tuple[str, ...]:
    """Name each statement or manifest key of a snippet file that differs from the one in a package.

    By default, the import block and anything that the package does not have
    are not compared. With ``whole``, the import block must also be the same,
    and the package must have everything that the file shows. Only a statement
    that a pending change adds, and the import block of a file that the
    change adds, can be absent.
    """
    target = package / file.path
    where = package.relative_to(REPO).as_posix()
    if not target.exists() and not whole:
        return ()
    old = load_mapping(target.read_text()) if target.exists() else {}
    new = load_mapping(file.text)
    changed = tuple(
        f"{file.path}: {key} is not the same as in {where}"
        for key, value in new.items()
        if (whole or key != "import")
        and key in old
        and old[key] != as_before_change(value)
    )
    absent = tuple(
        f"{file.path}: {key} is not in {where}"
        for key, value in new.items()
        if whole
        and key not in old
        and not is_added(value)
        and not (key == "import" and not target.exists())
    )
    return changed + absent


def run_tool(
    tool: Path,
    workdir: Path,
    command: str,
    package: str = PACKAGE,
    packages: Path | None = None,
) -> tuple[bool, str]:
    """Run one of-lang command on a package in a workdir.

    ``packages`` is where the tool looks for the dependencies of the package.
    By default, that is the folder that holds the package.
    """
    flags = [] if packages is None else ["--packages", str(packages)]
    done = subprocess.run(  # noqa: S603 — runs the tool this script just built
        [str(tool), command, *flags, package],
        cwd=workdir,
        capture_output=True,
        text=True,
        check=False,
    )
    return done.returncode == 0, (done.stdout + done.stderr).strip()


def check_package(tool: Path, package: Path) -> str:
    """Check and test a package that depends on the examples. Return what the first failure prints."""
    for command in ("check", "test"):
        ok, output = run_tool(tool, package.parent, command, package.name, EXAMPLES)
        if not ok:
            return f"of-lang {command}:\n{output}"
    return ""


def check_merged(tool: Path, package: Path, files: Sequence[SnippetFile]) -> str:
    """Merge snippet files into a fresh copy of a package, then check and test the copy."""
    with tempfile.TemporaryDirectory() as tmp:
        work = Path(tmp) / package.name
        shutil.copytree(package, work)
        for file in files:
            apply(work, file)
        return check_package(tool, work)


def check_snippet(tool: Path, snippet: Snippet) -> Problem | None:
    """Compare a snippet with the real package, then merge it into a fresh copy and check and test it."""
    diffs = [d for file in snippet.files for d in differences(file)]
    if diffs:
        return Problem(snippet, "\n".join(diffs))
    failure = check_merged(tool, EXAMPLES / PACKAGE, snippet.files)
    return Problem(snippet, failure) if failure else None


def shown_files(page: Path, header: re.Pattern[str] = HEADER) -> dict[str, list[str]]:
    """Return the lines that a page shows for each file, in page order.

    A block that starts with a file comment opens that file. A block without
    one continues the file of the block before it.
    """
    files: dict[str, list[str]] = {}
    path = ""
    for _, lines in yaml_blocks(page.read_text()):
        for line in lines:
            match = header.match(line)
            if match:
                path = match.group(1)
                files.setdefault(path, [])
            elif line.strip():
                files.setdefault(path, []).append(line.rstrip())
    return files


def console_blocks(text: str) -> tuple[tuple[int, str, str], ...]:
    """Return each console block of a page: its line, its command and its output."""
    blocks: list[tuple[int, str, str]] = []
    lines = text.splitlines()
    for number, line in enumerate(lines):
        if line.strip() != CONSOLE_OPEN:
            continue
        end = next(
            i for i in range(number + 1, len(lines)) if lines[i].strip() == FENCE_CLOSE
        )
        command, output = lines[number + 1], "\n".join(lines[number + 2 : end])
        if command.startswith(PROMPT):
            blocks.append((number + 1, command.removeprefix(PROMPT), output.strip()))
    return tuple(blocks)


def check_quickstart(tool: Path) -> tuple[str, ...]:
    """Compare the README and the home page with the packages that they show."""
    return tuple(
        problem for page in QUICKSTART_PAGES for problem in check_quickstart_page(tool, page)
    )


def check_quickstart_page(tool: Path, page: Path) -> tuple[str, ...]:
    """Compare one page with the quickstart packages, and run its commands."""
    where = page.relative_to(REPO)
    shown = shown_files(page, QUICKSTART_HEADER)
    real = {
        file.relative_to(QUICKSTART).as_posix(): [
            line.rstrip() for line in file.read_text().splitlines() if line.strip()
        ]
        for package in QUICKSTART.iterdir()
        for pattern in PACKAGE_FILES
        for file in package.glob(pattern)
    }
    problems = [
        f"{where}: {path} is not the same as the file in examples/quickstart"
        if path in real
        else f"{where}: examples/quickstart has no file {path}"
        for path, lines in shown.items()
        if real.get(path) != lines
    ]
    problems.extend(
        f"{where}: the page does not show {path}"
        for path in sorted(real.keys() - shown.keys())
        if Path(path).name not in MANIFESTS
    )
    for line, command, output in console_blocks(page.read_text()):
        done = subprocess.run(  # noqa: S603 — runs the tool this script just built
            [str(tool), *command.split()],
            cwd=QUICKSTART,
            capture_output=True,
            text=True,
            check=False,
        )
        if (done.stdout + done.stderr).strip() != output:
            problems.append(
                f"{where}:{line}: of-lang {command} does not give the output that the page shows"
            )
    return tuple(problems)


def check_example(
    tool: Path, package: Path, page: Path | None = None
) -> tuple[str, ...]:
    """Compare a by-example page with its package, then check and test the package."""
    page = page or package.with_suffix(".md")
    where = page.relative_to(REPO)
    if not page.exists():
        return (f"{package.relative_to(REPO)}: the package has no page {page.name}",)
    shown = shown_files(page)
    if "" in shown:
        return (
            f"{where}: the first snippet does not start with a comment that names its file",
        )
    real = {
        str(file.relative_to(package)): [
            line.rstrip() for line in file.read_text().splitlines() if line.strip()
        ]
        for pattern in PACKAGE_FILES
        for file in package.glob(pattern)
    }
    problems = [
        f"{where}: {path} is not the same as the file in the package"
        if path in real
        else f"{where}: the package has no file {path}"
        for path, lines in shown.items()
        if real.get(path) != lines
    ]
    problems.extend(
        f"{where}: the page does not show {path}"
        for path in sorted(real.keys() - shown.keys() - MANIFESTS)
    )
    for command in ("check", "test"):
        ok, output = run_tool(tool, package.parent, command, package.name)
        if not ok:
            problems.append(
                f"{package.relative_to(REPO)}: of-lang {command}:\n{output}"
            )
    return tuple(problems)


def all_files(package: Path) -> dict[str, bytes]:
    """Return the content of every file of a package, by its path in the package."""
    return {
        file.relative_to(package).as_posix(): file.read_bytes()
        for file in package.rglob("*")
        if file.is_file()
    }


def check_finished(step: Path) -> tuple[str, ...]:
    """Compare the last package of the introduction with the restaurant example."""
    where = step.relative_to(REPO).as_posix()
    shown, real = all_files(step), all_files(EXAMPLES / FINISHED)
    return (
        *(
            f"{where}: {path} is not the same as in examples/{FINISHED}"
            for path in sorted(shown.keys() & real.keys())
            if shown[path] != real[path]
        ),
        *(
            f"{where}: {path} is not in examples/{FINISHED}"
            for path in sorted(shown.keys() - real.keys())
        ),
        *(
            f"{where}: the package does not have {path} of examples/{FINISHED}"
            for path in sorted(real.keys() - shown.keys())
        ),
    )


def check_step(tool: Path, step: Path) -> tuple[str, ...]:
    """Check and test one package of the introduction, and find its chapter."""
    where = step.relative_to(REPO).as_posix()
    failure = check_package(tool, step)
    chapters = tuple(INTRODUCTION.glob(f"{step.name.removeprefix(STEP_PREFIX)}-*.md"))
    return (
        *((f"{where}: {failure}",) if failure else ()),
        *(() if chapters else (f"{where}: the package has no chapter",)),
    )


def check_chapter(tool: Path, page: Path) -> tuple[int, tuple[str, ...]]:
    """Compare the snippets of a chapter of the introduction with its package.

    Return the number of snippets and the problems. When a snippet shows a
    pending change, the snippets of the chapter together are that change: they
    are merged into a fresh copy of the package, which is then checked and
    tested.
    """
    where = page.relative_to(REPO).as_posix()
    step = STEPS / f"{STEP_PREFIX}{page.name[:2]}"
    if not (step / "package.yaml").exists():
        return 0, (f"{where}: the chapter has no package steps/{step.name}",)
    snippets, unnamed = find_snippets(page)
    problems = [
        *unnamed,
        *(
            f"{where}:{snippet.line}: {difference}"
            for snippet in snippets
            for file in snippet.files
            for difference in differences(file, step, whole=True)
        ),
    ]
    files = tuple(file for snippet in snippets for file in snippet.files)
    changes = any(
        has_pending(value)
        for file in files
        for value in load_mapping(file.text).values()
    )
    if changes and not problems:
        failure = check_merged(tool, step, files)
        problems.extend(
            (f"{where}: with its pending change: {failure}",) if failure else ()
        )
    return len(snippets), tuple(problems)


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
    examples = sorted(p.parent for p in BY_EXAMPLE.glob("*/package.y*ml"))
    steps = sorted(p.parent for p in STEPS.glob("*/package.y*ml"))
    chapters = sorted(INTRODUCTION.glob(CHAPTERS))
    shown = len(snippets)
    with tempfile.TemporaryDirectory() as tmp:
        tool = build_tool(Path(tmp))
        for example in examples:
            problems.extend(check_example(tool, example))
        problems.extend(check_quickstart(tool))
        for snippet in snippets:
            problem = check_snippet(tool, snippet)
            if problem is not None:
                where = (
                    f"{problem.snippet.page.relative_to(REPO)}:{problem.snippet.line}"
                )
                problems.append(f"{where}: {problem.text}")
        for step in steps:
            problems.extend(check_step(tool, step))
        problems.extend(check_finished(steps[-1]) if steps else ())
        for chapter in chapters:
            count, wrong = check_chapter(tool, chapter)
            shown += count
            problems.extend(wrong)
    for p in problems:
        print(p, file=sys.stderr)
    print(
        f"{shown} snippets, {len(examples)} by-example pages and {len(chapters)} chapters of the introduction checked, {len(problems)} problems"
    )
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(main())
