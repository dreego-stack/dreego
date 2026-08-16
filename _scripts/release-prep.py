#!/usr/bin/env python3
"""Apply a PR's pr.md to CHANGELOG.md and print the next version.

Reads pr.md from the current directory, validates the version field,
computes the next version from the latest git tag (vX.Y.Z), prepends a
CHANGELOG entry, and removes pr.md. The VERSION file is no longer used —
the git tag is the single source of truth.

Changelog format:
- version=none: prepend changelog lines at the very top of the file
- version=patch: prepend a version block (blank line,
  '## vX.Y.Z - YYYY-MM-DD', blank line) followed by the changelog lines.

Prints 'new=vX.Y.Z' or 'new=none' on stdout for the workflow to consume.

Usage: python3 _scripts/release-prep.py
Exit 0 on success, non-zero on validation error.
"""

import re
import subprocess
import sys
from datetime import date
from pathlib import Path

ROOT = Path.cwd()
PR_MD = ROOT / "pr.md"
CHANGELOG = ROOT / "CHANGELOG.md"

VALID_VERSIONS = ("none", "patch")


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


def parse_pr_md(text):
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n?(.*)$", text, re.DOTALL)
    if not m:
        fail("pr.md must start with YAML frontmatter (--- ... ---)")
    fm, body = m.group(1), m.group(2)
    version = None
    for line in fm.strip().splitlines():
        line = line.strip()
        if line.startswith("version:"):
            version = line.split(":", 1)[1].strip()
    if version is None:
        fail("pr.md frontmatter must contain 'version: none|patch' (minor/major are blocked in the v0.0.x phase)")
    if version not in VALID_VERSIONS:
        fail(f"pr.md version must be one of {VALID_VERSIONS}, got '{version}'")
    lines = [l.strip() for l in body.splitlines() if l.strip()]
    if not lines:
        fail("pr.md has no changelog lines")
    lines = [l[2:].strip() if l.startswith("- ") else l for l in lines]
    return version, lines


def next_version(current, bump):
    m = re.match(r"^v?(\d+)\.(\d+)\.(\d+)$", current.strip())
    if not m:
        fail(f"latest tag has invalid format: '{current}'")
    major, minor, patch = (int(g) for g in m.groups())
    if bump == "patch":
        patch += 1
    return f"v{major}.{minor}.{patch}"


def main():
    if not PR_MD.exists():
        fail("pr.md not found in PR branch")
    version, lines = parse_pr_md(PR_MD.read_text())

    current = latest_tag()
    new_version = next_version(current, version) if version != "none" else None

    today = date.today().isoformat()
    old = CHANGELOG.read_text() if CHANGELOG.exists() else ""
    if old and not old.startswith("\n"):
        old = "\n" + old

    lines_text = "\n".join(f"- {l}" for l in lines) + "\n"

    if new_version:
        entry = f"\n## {new_version} - {today}\n\n{lines_text}"
        if f"## {new_version} -" in old:
            print(f"new={new_version}")
            print(f"skipped: version {new_version} already in CHANGELOG", file=sys.stderr)
            PR_MD.unlink()
            return
    else:
        entry = lines_text

    CHANGELOG.write_text(entry + old)

    PR_MD.unlink()
    print(f"new={new_version or 'none'}")
    print(f"applied: version={version} new={new_version or '(none)'} lines={len(lines)}", file=sys.stderr)


if __name__ == "__main__":
    main()