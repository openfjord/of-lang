# SPDX-License-Identifier: Apache-2.0
# /// script
# requires-python = ">=3.12"
# dependencies = ["markdown-it-py==3.0.0", "pygments==2.18.0", "pyyaml==6.0.2"]
# ///
"""Build the documentation site from the Markdown files of this repository.

A Markdown file is a page of the site when it starts with a block that gives
its ``title``. ``README.md`` becomes the ``index.html`` of its folder, and any
other file keeps its name, with ``.html``. Links between pages are rewritten,
and a link to any other file of the repository goes to GitHub. The site
leaves out the text between a line ``<!-- not on the site -->`` and a line
``<!-- end -->``, which GitHub shows. The build fails
when a link points at a file or a heading that does not exist, or when a page
is missing from the navigation.

Usage: uv run scripts/build_site.py [output folder, default _site]
"""

import html
import posixpath
import re
import shutil
import sys
from collections.abc import Mapping, Sequence
from dataclasses import dataclass
from pathlib import Path
from typing import cast

import yaml
from markdown_it import MarkdownIt
from markdown_it.token import Token
from pygments import highlight
from pygments.formatters import HtmlFormatter
from pygments.lexers import get_lexer_by_name
from pygments.util import ClassNotFound

REPO = Path(__file__).resolve().parents[1]
ASSETS = Path(__file__).resolve().parent / "site"
GITHUB = "https://github.com/openfjord/of-lang"
SITE_NAME = "The OpenFjord modeling language"
HOME = "README.md"
SKIPPED_FOLDERS = frozenset({"_site", "node_modules"})
IMAGE_ENDINGS = frozenset({".png", ".svg", ".jpg", ".jpeg", ".gif", ".webp"})
HAS_SCHEME = re.compile(r"^[a-zA-Z][a-zA-Z0-9+.-]*:|^//")
HTML_LINK = re.compile(r'\b(src|href)="([^"]*)"')
FRONT_MATTER = re.compile(r"\A---\n(.*?)\n---\n", re.DOTALL)
NOT_ON_SITE = re.compile(
    r"^<!-- not on the site -->\n.*?^<!-- end -->\n\n?", re.DOTALL | re.MULTILINE
)
DARK_STYLE = "github-dark"
LIGHT_STYLE = "default"


@dataclass(frozen=True)
class Section:
    """One group of the navigation.

    The group holds the listed pages, and then each page in ``folder`` that
    the ``index`` page links to, in the order of the links. With ``grouped``,
    the linked pages go in smaller groups, one for each heading of the index
    page that has links under it.
    """

    title: str
    pages: tuple[str, ...] = ()
    index: str = ""
    folder: str = ""
    grouped: bool = False


SECTIONS = (
    Section("", pages=(HOME,)),
    Section(
        "Introduction",
        pages=("docs/introduction/README.md",),
        index="docs/introduction/README.md",
        folder="docs/introduction/",
    ),
    Section("Language", pages=("language/README.md", "language/cheatsheet.md")),
    Section("Reference", index="language/README.md", folder="language/reference/"),
    Section(
        "By example",
        pages=("docs/by-example/README.md",),
        index="docs/by-example/README.md",
        folder="docs/by-example/",
        grouped=True,
    ),
    Section("More", pages=("docs/design.md", "cli/README.md")),
)


@dataclass(frozen=True)
class Page:
    """One Markdown file of the site, parsed."""

    path: str
    title: str
    heading: str
    tokens: tuple[Token, ...]
    anchors: frozenset[str]
    links: tuple[str, ...]
    link_headings: tuple[str, ...]


@dataclass(frozen=True)
class Target:
    """Where a link goes: the address to write, or the problem with it."""

    url: str = ""
    asset: str = ""
    problem: str = ""


def output_path(path: str) -> str:
    """Return where a page lands in the site."""
    folder, name = posixpath.split(path)
    stem = "index" if name == "README.md" else name.removesuffix(".md")
    return posixpath.join(folder, stem + ".html")


def relative_url(from_page: str, to_file: str) -> str:
    """Return the address of a site file as seen from a page."""
    return posixpath.relpath(to_file, posixpath.dirname(output_path(from_page)) or ".")


def slug(text: str, seen: int) -> str:
    """Return the id that GitHub gives a heading: its text, and a number when the text repeats."""
    base = re.sub(r"[^\w\- ]", "", text.strip().lower()).replace(" ", "-")
    return base if seen == 0 else f"{base}-{seen}"


def plain_text(inline: Token) -> str:
    """Return the text of an inline token without its markup."""
    return "".join(
        child.content
        for child in inline.children or ()
        if child.type in {"text", "code_inline"}
    )


def highlighted(code: str, language: str, _attrs: str) -> str:
    """Return a code block as HTML with a class on each kind of word."""
    try:
        lexer = get_lexer_by_name(language) if language else None
    except ClassNotFound:
        lexer = None
    if lexer is None:
        return html.escape(code)
    return cast("str", highlight(code, lexer, HtmlFormatter(nowrap=True)))


def markdown() -> MarkdownIt:
    """Return a reader that follows the rules GitHub uses for Markdown files."""
    reader = MarkdownIt("commonmark", {"html": True, "highlight": highlighted})
    return reader.enable(["table", "strikethrough"])


def parse(path: str, text: str) -> Page | None:
    """Parse a Markdown file. None when it does not give a title."""
    front = FRONT_MATTER.match(text)
    meta = cast("object", yaml.safe_load(front.group(1))) if front else None
    title = (
        cast("Mapping[str, object]", meta).get("title")
        if isinstance(meta, dict)
        else None
    )
    if front is None or not isinstance(title, str):
        return None
    shown = NOT_ON_SITE.sub("", text[front.end() :])
    body = re.sub(r"(?m)^<table>$", '<table class="steps">', shown)
    tokens = markdown().parse(body)
    counts: dict[str, int] = {}
    heading = ""
    links: list[str] = []
    link_headings: list[str] = []
    section = ""
    for index, token in enumerate(tokens):
        if token.type == "heading_open":
            text_of = plain_text(tokens[index + 1])
            section = text_of if token.tag == "h2" else section
            base = slug(text_of, 0)
            token.attrSet("id", slug(text_of, counts.get(base, 0)))
            counts[base] = counts.get(base, 0) + 1
            heading = heading or text_of
        found = [
            str(child.attrGet("href"))
            for child in token.children or ()
            if child.type == "link_open"
        ]
        links.extend(found)
        link_headings.extend(section for _ in found)
    anchors = frozenset(
        str(t.attrGet("id")) for t in tokens if t.type == "heading_open"
    )
    return Page(
        path, title, heading, tuple(tokens), anchors, tuple(links), tuple(link_headings)
    )


def repo_path(page: str, link: str) -> str:
    """Return the file of the repository that a relative link names."""
    start = "" if link.startswith("/") else posixpath.dirname(page)
    joined = posixpath.normpath(posixpath.join(start, link.lstrip("/")))
    return "" if joined == "." else joined


def resolve(
    page: Page, link: str, pages: Mapping[str, Page], files: frozenset[str]
) -> Target:
    """Decide where a link of a page goes."""
    if HAS_SCHEME.match(link):
        return Target(url=link)
    path_part, _, fragment = link.partition("#")
    if not path_part:
        ok = not fragment or fragment in page.anchors
        return Target(url=link) if ok else Target(problem=f"no heading #{fragment}")
    return resolve_file(page, repo_path(page.path, path_part), fragment, pages, files)


def resolve_file(
    page: Page,
    path: str,
    fragment: str,
    pages: Mapping[str, Page],
    files: frozenset[str],
) -> Target:
    """Decide where a link to a file or a folder of the repository goes."""
    is_folder = any(f.startswith(path + "/") for f in files) or not path
    page_path = posixpath.join(path, "README.md") if is_folder else path
    if page_path in pages:
        if fragment and fragment not in pages[page_path].anchors:
            return Target(problem=f"{page_path} has no heading #{fragment}")
        url = relative_url(page.path, output_path(page_path))
        return Target(url=url + (f"#{fragment}" if fragment else ""))
    if is_folder:
        return Target(url=f"{GITHUB}/tree/main/{path}")
    if path not in files:
        return Target(problem=f"no file {path}")
    if posixpath.splitext(path)[1].lower() in IMAGE_ENDINGS:
        return Target(url=relative_url(page.path, path), asset=path)
    return Target(url=f"{GITHUB}/blob/main/{path}")


def rewritten(
    page: Page, pages: Mapping[str, Page], files: frozenset[str]
) -> tuple[frozenset[str], tuple[str, ...]]:
    """Point every link of a page at its place. Return the assets used and the problems."""
    assets: set[str] = set()
    problems: list[str] = []

    def address(link: str) -> str:
        target = resolve(page, link, pages, files)
        if target.problem:
            problems.append(f"{page.path}: link {link}: {target.problem}")
            return link
        if target.asset:
            assets.add(target.asset)
        return target.url

    def in_html(found: re.Match[str]) -> str:
        return (
            f'{found.group(1)}="{html.escape(address(html.unescape(found.group(2))))}"'
        )

    for token in page.tokens:
        for part in (token, *(token.children or ())):
            if part.type == "link_open":
                part.attrSet("href", address(str(part.attrGet("href"))))
            elif part.type == "image":
                part.attrSet("src", address(str(part.attrGet("src"))))
            elif part.type in {"html_block", "html_inline"}:
                part.content = HTML_LINK.sub(in_html, part.content)
    return frozenset(assets), tuple(problems)


@dataclass(frozen=True)
class Group:
    """One group of the navigation: its pages, and then its smaller groups."""

    title: str
    pages: tuple[str, ...]
    parts: tuple[tuple[str, tuple[str, ...]], ...] = ()


def linked_files(index: Page | None) -> tuple[tuple[str, str], ...]:
    """Return each file that a page links to, in order, with the heading above the link."""
    if index is None:
        return ()
    return tuple(
        (heading, repo_path(index.path, link.partition("#")[0]))
        for heading, link in zip(index.link_headings, index.links, strict=True)
        if not HAS_SCHEME.match(link) and link.partition("#")[0]
    )


def navigation(pages: Mapping[str, Page]) -> tuple[Group, ...]:
    """Return each group of the navigation with the pages it holds, in order."""
    groups: list[Group] = []
    for section in SECTIONS:
        listed = tuple(p for p in section.pages if p in pages)
        linked: dict[str, str] = {}
        for heading, path in linked_files(pages.get(section.index)):
            if path.startswith(section.folder) and path in pages and path not in listed:
                linked.setdefault(path, heading)
        if not section.grouped:
            groups.append(Group(section.title, (*listed, *linked)))
            continue
        headings = tuple(dict.fromkeys(linked.values()))
        parts = tuple(
            (heading, tuple(p for p, h in linked.items() if h == heading))
            for heading in headings
        )
        groups.append(Group(section.title, listed, parts))
    return tuple(groups)


def group_pages(group: Group) -> tuple[str, ...]:
    """Return every page of a group, also the pages of its smaller groups."""
    return (*group.pages, *(path for _, paths in group.parts for path in paths))


def label(page: Page, group: str) -> str:
    """Return the words that the navigation shows for a page."""
    if page.path == HOME:
        return "Home"
    if page.path.endswith("README.md") and group not in {"Language", "More"}:
        return "Overview"
    return page.heading or page.title


def links_html(
    paths: Sequence[str], group: str, pages: Mapping[str, Page], current: Page
) -> str:
    """Return a list of links to pages, as seen from one page."""
    items: list[str] = []
    for path in paths:
        here = ' aria-current="page"' if path == current.path else ""
        url = html.escape(relative_url(current.path, output_path(path)))
        name = html.escape(label(pages[path], group))
        items.append(f'<li><a href="{url}"{here}>{name}</a></li>')
    return "<ul>\n" + "\n".join(items) + "\n</ul>"


def navigation_html(
    groups: Sequence[Group], pages: Mapping[str, Page], current: Page
) -> str:
    """Return the navigation as seen from one page.

    A smaller group is closed, unless it holds the page.
    """
    parts: list[str] = []
    for group in groups:
        if group.title:
            parts.append(f'<p class="group">{html.escape(group.title)}</p>')
        parts.append(links_html(group.pages, group.title, pages, current))
        for title, paths in group.parts:
            is_open = " open" if current.path in paths else ""
            parts.append(f"<details{is_open}><summary>{html.escape(title)}</summary>")
            parts.append(links_html(paths, group.title, pages, current))
            parts.append("</details>")
    return "\n".join(parts)


def code_colors() -> str:
    """Return the style rules that color code, for the dark and the light theme."""

    def rules(style: str, prefix: str) -> str:
        defs = cast("str", HtmlFormatter(style=style).get_style_defs(prefix))
        return "\n".join(
            line for line in defs.splitlines() if re.match(r"^\S.* \.\w+ \{", line)
        )

    light = ':root[data-theme="light"] pre'
    reset = (
        f"{light} span {{ color: inherit; font-weight: inherit; font-style: inherit; }}"
    )
    return "\n".join((rules(DARK_STYLE, "pre"), reset, rules(LIGHT_STYLE, light)))


def page_html(page: Page, nav: str, template: str) -> str:
    """Return the whole HTML file of a page."""
    body = markdown().renderer.render(list(page.tokens), markdown().options, {})
    values = {
        "title": html.escape(page.title),
        "site": html.escape(SITE_NAME),
        "home": html.escape(relative_url(page.path, output_path(HOME))),
        "assets": html.escape(relative_url(page.path, "assets")),
        "source": html.escape(f"{GITHUB}/blob/main/{page.path}"),
        "nav": nav,
        "body": body,
    }
    return re.sub(r"\{\{(\w+)\}\}", lambda found: values[found.group(1)], template)


def repo_files() -> frozenset[str]:
    """Return every file of the repository, as a path from its top."""
    return frozenset(
        path.relative_to(REPO).as_posix()
        for path in REPO.rglob("*")
        if path.is_file()
        and not any(
            part.startswith(".") or part in SKIPPED_FOLDERS
            for part in path.relative_to(REPO).parts
        )
    )


def main(arguments: Sequence[str]) -> int:
    """Build the site and print each problem."""
    out = Path(arguments[0]).resolve() if arguments else REPO / "_site"
    files = repo_files()
    parsed = (
        parse(path, (REPO / path).read_text())
        for path in sorted(files)
        if path.endswith(".md")
    )
    pages = {page.path: page for page in parsed if page is not None}
    groups = navigation(pages)
    in_navigation = {path for group in groups for path in group_pages(group)}
    problems = [
        f"{path}: the page is not in the navigation"
        for path in sorted(pages.keys() - in_navigation)
    ]
    template = (ASSETS / "page.html").read_text()

    shutil.rmtree(out, ignore_errors=True)
    (out / "assets").mkdir(parents=True)
    for name in ("style.css", "theme.js"):
        shutil.copy(ASSETS / name, out / "assets" / name)
    (out / "assets" / "code.css").write_text(code_colors())
    (out / ".nojekyll").write_text("")
    for page in pages.values():
        assets, found = rewritten(page, pages, files)
        problems.extend(found)
        target = out / output_path(page.path)
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(
            page_html(page, navigation_html(groups, pages, page), template)
        )
        for asset in assets:
            (out / asset).parent.mkdir(parents=True, exist_ok=True)
            shutil.copy(REPO / asset, out / asset)

    for problem in problems:
        print(problem, file=sys.stderr)
    print(f"{len(pages)} pages built in {out}, {len(problems)} problems")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
