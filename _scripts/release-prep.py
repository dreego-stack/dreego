#!/usr/bin/env python3
"""Apply pending change files to CHANGELOG.md and print the next version.

Reads .changes/*.md from the current directory, validates each version field,
computes the next version from the latest git tag (vX.Y.Z), prepends a
CHANGELOG entry, and removes the processed files. The VERSION file is no longer used —
the git tag is the single source of truth. When only version=none files are
pending, nothing is applied: the files stay pending and are processed together
with the next version=patch file.

Changelog format:
- version=none: deferred — the change file stays pending until a patch
  file arrives; then all pending files (none + patch) are applied together
- version=patch: prepend a version block (blank line,
  '## vX.Y.Z - YYYY-MM-DD', blank line) followed by the changelog lines.
- version=minor: only valid for stage/* branch merges; produces the next
  minor tag (v0.1.0 -> v0.2.0) and resets the patch to 0.

Prints 'new=vX.Y.Z' or 'new=none' on stdout for the workflow to consume.

Usage: python3 _scripts/release-prep.py [--verify-tags]
Exit 0 on success, non-zero on validation error.
"""

import re
import subprocess
import sys
from datetime import date
from pathlib import Path

ROOT = Path.cwd()
CHANGES = ROOT / ".changes"
CHANGELOG = ROOT / "CHANGELOG.md"

VALID_VERSIONS = ("none", "patch", "minor")
MODULE_PATH = "github.com/dreego-stack/dreego"
TAG_PREFIXES = ("", "core/", "adapter/ssr/", "dreegotest/", "cmd/dreego/")


def fail(msg):
    print(f"error: {msg}", file=sys.stderr)
    sys.exit(1)


def latest_tag():
    out = subprocess.run(
        ["git", "tag", "-l", "v[0-9]*"],
        capture_output=True, text=True, cwd=ROOT,
    )
    tags = [t for t in out.stdout.splitlines() if re.match(r"^v\d+\.\d+\.\d+$", t)]
    if not tags:
        return "v0.0.0"
    return sorted(tags, key=lambda t: [int(x) for x in t[1:].split(".")])[-1]


def release_tags(version):
    return [f"{prefix}{version}" for prefix in TAG_PREFIXES]


def git_tags():
    result = subprocess.run(
        ["git", "tag", "--list"], capture_output=True, text=True, cwd=ROOT, check=True
    )
    return set(result.stdout.splitlines())


def tag_commit(tag):
    result = subprocess.run(
        ["git", "rev-list", "-n", "1", tag],
        capture_output=True,
        text=True,
        cwd=ROOT,
        check=True,
    )
    return result.stdout.strip()


def verify_tag_groups():
    tags = git_tags()
    allowed = re.compile(
        r"^(?:core/|adapter/ssr/|dreegotest/|cmd/dreego/)?v\d+\.\d+\.\d+$"
    )
    for tag in tags:
        if "/v" in tag and not allowed.match(tag):
            fail(f"unexpected module release tag: {tag}")
    roots = sorted(
        tag for tag in tags
        if re.match(r"^v\d+\.\d+\.\d+$", tag)
        and tuple(int(part) for part in tag[1:].split(".")) >= (0, 8, 0)
    )
    for root_tag in roots:
        expected = release_tags(root_tag)
        missing = [tag for tag in expected if tag not in tags]
        if missing:
            fail(f"incomplete coordinated release {root_tag}; missing: {' '.join(missing)}")
        commits = {tag_commit(tag) for tag in expected}
        if len(commits) != 1:
            fail(f"coordinated release tags do not share one commit: {' '.join(expected)}")
    print(f"verified={len(roots)}")


def update_internal_requirements(version):
    pattern = re.compile(
        rf"^(\s*(?:require\s+)?)({re.escape(MODULE_PATH)}(?:/[^\s]+)?)(\s+)v\d+\.\d+\.\d+(\s*(?://.*)?)$",
        re.MULTILINE,
    )
    changed = []
    for path in sorted(ROOT.rglob("go.mod")):
        if ".git" in path.parts or ".worktrees" in path.parts:
            continue
        old = path.read_text()
        new = pattern.sub(rf"\1\2\3{version}\4", old)
        if new != old:
            path.write_text(new)
            changed.append(path.relative_to(ROOT))
    workspace = ROOT / "go.work"
    if workspace.exists():
        old = workspace.read_text()
        workspace_pattern = re.compile(
            rf"^(\s*replace\s+)({re.escape(MODULE_PATH)}(?:/[^\s]+)?)(\s+)v\d+\.\d+\.\d+(\s+=>\s+.+)$",
            re.MULTILINE,
        )
        new = workspace_pattern.sub(rf"\1\2\3{version}\4", old)
        if new != old:
            workspace.write_text(new)
            changed.append(workspace.relative_to(ROOT))
    return changed


def parse_change(text, source):
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n?(.*)$", text, re.DOTALL)
    if not m:
        fail(f"{source} must start with YAML frontmatter (--- ... ---)")
    fm, body = m.group(1), m.group(2)
    version = None
    for line in fm.strip().splitlines():
        line = line.strip()
        if line.startswith("version:"):
            version = line.split(":", 1)[1].strip()
    if version is None:
        fail(f"{source} frontmatter must contain 'version: none|patch'")
    if version not in VALID_VERSIONS:
        fail(f"{source} version must be one of {VALID_VERSIONS}, got '{version}'")
    lines = [l.strip() for l in body.splitlines() if l.strip()]
    if not lines:
        fail(f"{source} has no changelog lines")
    lines = [l[2:].strip() if l.startswith("- ") else l for l in lines]
    return version, lines


def next_version(current, bump):
    m = re.match(r"^v?(\d+)\.(\d+)\.(\d+)$", current.strip())
    if not m:
        fail(f"latest tag has invalid format: '{current}'")
    major, minor, patch = (int(g) for g in m.groups())
    if major != 0:
        fail(f"latest tag must be v0.x.y, got '{current}'")
    if bump == "patch":
        patch += 1
    elif bump == "minor":
        minor += 1
        patch = 0
    return f"v{major}.{minor}.{patch}"


def main():
    if sys.argv[1:] == ["--verify-tags"]:
        verify_tag_groups()
        return
    if sys.argv[1:]:
        fail("usage: release-prep.py [--verify-tags]")
    files = sorted(p for p in CHANGES.glob("*.md") if p.name != "README.md")
    if not files:
        fail("no pending .changes/*.md files found")
    parsed = [parse_change(p.read_text(), p.relative_to(ROOT)) for p in files]
    if any(v == "minor" for v, _ in parsed):
        version = "minor"
    elif any(v == "patch" for v, _ in parsed):
        version = "patch"
    else:
        version = "none"
    lines = [line for _, change_lines in parsed for line in change_lines]

    if version == "none":
        for path in files:
            print(f"deferred: {path.relative_to(ROOT)} (no version bump)", file=sys.stderr)
        print("new=none")
        return

    current = latest_tag()
    new_version = next_version(current, version)

    updated_modules = update_internal_requirements(new_version)

    today = date.today().isoformat()
    old = CHANGELOG.read_text() if CHANGELOG.exists() else ""
    if old and not old.startswith("\n"):
        old = "\n" + old

    lines_text = "\n".join(f"- {l}" for l in lines) + "\n"

    entry = f"\n## {new_version} - {today}\n\n{lines_text}"
    if f"## {new_version} -" in old:
        print(f"new={new_version}")
        print(f"tags={' '.join(release_tags(new_version))}")
        print(f"skipped: version {new_version} already in CHANGELOG", file=sys.stderr)
        for path in files:
            path.unlink()
        return

    CHANGELOG.write_text(entry + old)

    for path in files:
        path.unlink()
    print(f"new={new_version}")
    print(f"tags={' '.join(release_tags(new_version))}")
    print(f"applied: version={version} new={new_version} lines={len(lines)}", file=sys.stderr)
    for path in updated_modules:
        print(f"updated: {path}", file=sys.stderr)


if __name__ == "__main__":
    main()
