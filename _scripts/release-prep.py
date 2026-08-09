#!/usr/bin/env python3
"""Apply a PR's pr.md to CHANGELOG.md and VERSION.

Reads pr.md from the current directory, validates the version field,
computes the next version from the current VERSION file, prepends a
CHANGELOG entry, updates VERSION, and removes pr.md.

Changelog format:
- version=none: prepend changelog lines at the very top of the file
- version=patch|minor|major: prepend a version block (blank line,
  '## vX.Y.Z - YYYY-MM-DD', blank line) followed by the changelog
  lines. Update VERSION file.

Usage: python3 _scripts/release-prep.py
Exit 0 on success, non-zero on validation error.
"""

import re
import sys
from datetime import date
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
PR_MD = ROOT / "pr.md"
CHANGELOG = ROOT / "CHANGELOG.md"
VERSION_FILE = ROOT / "VERSION"

VALID_VERSIONS = ("none", "patch", "minor", "major")


def fail(msg):
    print(f"error: {msg}", file=sys.stderr)
    sys.exit(1)


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
        fail("pr.md frontmatter must contain 'version: none|patch|minor|major'")
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
        fail(f"VERSION file has invalid format: '{current}'")
    major, minor, patch = (int(g) for g in m.groups())
    if bump == "patch":
        patch += 1
    elif bump == "minor":
        minor += 1
        patch = 0
    elif bump == "major":
        major += 1
        minor = 0
        patch = 0
    return f"v{major}.{minor}.{patch}"


def main():
    if not PR_MD.exists():
        fail("pr.md not found in PR branch")
    version, lines = parse_pr_md(PR_MD.read_text())

    current = VERSION_FILE.read_text().strip() if VERSION_FILE.exists() else "v0.0.0"
    new_version = next_version(current, version) if version != "none" else None

    today = date.today().isoformat()
    old = CHANGELOG.read_text() if CHANGELOG.exists() else ""
    if old and not old.startswith("\n"):
        old = "\n" + old

    lines_text = "\n".join(f"- {l}" for l in lines) + "\n"

    if new_version:
        entry = f"\n## {new_version} - {today}\n\n{lines_text}"
    else:
        entry = lines_text

    CHANGELOG.write_text(entry + old)

    if new_version:
        VERSION_FILE.write_text(new_version + "\n")

    PR_MD.unlink()
    print(f"applied: version={version} new={new_version or '(none)'} lines={len(lines)}")


if __name__ == "__main__":
    main()