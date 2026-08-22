#!/usr/bin/env python3
"""Contract tests for the PR-driven release workflow.

Covers:
- release-prep.py behavior: combined patch, none, and failure paths
- workflow contract: expected workflow files exist, are named per AGENTS.md,
  serialize via concurrency groups, and tag only after make test

Usage: python3 _scripts/release-prep-test.py
Exit 0 on success, non-zero on any failed check.
"""

import os
import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SCRIPTS = ROOT / "_scripts"
WORKFLOWS = ROOT / ".github" / "workflows"

PASS = 0
FAIL = 0


def check(name, cond, detail=""):
    global PASS, FAIL
    if cond:
        PASS += 1
        print(f"PASS: {name}")
    else:
        FAIL += 1
        print(f"FAIL: {name} {detail}")


def run_release_prep(workdir, change, changelog, tags):
    repo = Path(workdir)
    changes = repo / ".changes"
    changes.mkdir()
    if change is not None:
        (changes / "change.md").write_text(change)
    (repo / "CHANGELOG.md").write_text(changelog)
    subprocess.run(["git", "init", "-q"], cwd=repo, check=True)
    subprocess.run(["git", "config", "user.email", "t@t"], cwd=repo, check=True)
    subprocess.run(["git", "config", "user.name", "t"], cwd=repo, check=True)
    subprocess.run(["git", "add", "-A"], cwd=repo, check=True)
    subprocess.run(["git", "commit", "-qm", "init"], cwd=repo, check=True)
    for tag in tags:
        subprocess.run(["git", "tag", tag], cwd=repo, check=True)
    return subprocess.run(
        [sys.executable, str(SCRIPTS / "release-prep.py")],
        cwd=repo, capture_output=True, text=True,
    )


def test_patch_path():
    with tempfile.TemporaryDirectory() as tmp:
        pr = "---\nversion: patch\n---\n\n- Feat: add widget\n"
        r = run_release_prep(tmp, pr, "## v0.0.43 - 2026-08-15\n\n- old\n", ["v0.0.43"])
        check("patch: exit 0", r.returncode == 0, r.stderr)
        check("patch: prints new version", "new=v0.0.44" in r.stdout, r.stdout)
        changelog = (Path(tmp) / "CHANGELOG.md").read_text()
        check("patch: version header added", "## v0.0.44 - " in changelog, changelog[:80])
        check("patch: changelog line added", "- Feat: add widget" in changelog)
        check("patch: old entry preserved", "## v0.0.43 - 2026-08-15" in changelog)
        check("patch: change removed", not (Path(tmp) / ".changes/change.md").exists())


def test_none_path():
    with tempfile.TemporaryDirectory() as tmp:
        pr = "---\nversion: none\n---\n\n- Chore: bump dep\n"
        r = run_release_prep(tmp, pr, "## v0.0.43 - 2026-08-15\n\n- old\n", ["v0.0.43"])
        check("none: exit 0", r.returncode == 0, r.stderr)
        check("none: prints none", "new=none" in r.stdout, r.stdout)
        changelog = (Path(tmp) / "CHANGELOG.md").read_text()
        check("none: no version header", "## v0.0.44" not in changelog)
        check("none: changelog unchanged", changelog == "## v0.0.43 - 2026-08-15\n\n- old\n", changelog[:60])
        check("none: change remains pending", (Path(tmp) / ".changes/change.md").exists())


def test_idempotent_rerun():
    with tempfile.TemporaryDirectory() as tmp:
        pr = "---\nversion: patch\n---\n\n- Feat: add widget\n"
        r1 = run_release_prep(tmp, pr, "## v0.0.43 - 2026-08-15\n\n- old\n", ["v0.0.43"])
        check("rerun: first run ok", r1.returncode == 0, r1.stderr)
        repo = Path(tmp)
        (repo / ".changes/change.md").write_text(pr)
        r2 = subprocess.run(
            [sys.executable, str(SCRIPTS / "release-prep.py")],
            cwd=tmp, capture_output=True, text=True,
        )
        check("rerun: exit 0", r2.returncode == 0, r2.stderr)
        changelog = (repo / "CHANGELOG.md").read_text()
        check("rerun: no duplicate header", changelog.count("## v0.0.44 - ") == 1, changelog[:200])


def test_failure_paths():
    cases = [
        ("missing change file", None, "## v0.0.43\n", ["v0.0.43"]),
        ("invalid version minor", "---\nversion: minor\n---\n\n- x\n", "## v0.0.43\n", ["v0.0.43"]),
        ("invalid version major", "---\nversion: major\n---\n\n- x\n", "## v0.0.43\n", ["v0.0.43"]),
        ("no changelog lines", "---\nversion: patch\n---\n", "## v0.0.43\n", ["v0.0.43"]),
        ("malformed frontmatter", "version: patch\n\n- x\n", "## v0.0.43\n", ["v0.0.43"]),
    ]
    for name, pr, changelog, tags in cases:
        with tempfile.TemporaryDirectory() as tmp:
            r = run_release_prep(tmp, pr, changelog, tags)
            check(f"failure {name}: non-zero exit", r.returncode != 0, f"exit={r.returncode}")


def test_combined_changes():
    with tempfile.TemporaryDirectory() as tmp:
        repo = Path(tmp)
        (repo / ".changes").mkdir()
        (repo / ".changes/a.md").write_text("---\nversion: none\n---\n\n- Chore: one\n")
        (repo / ".changes/b.md").write_text("---\nversion: patch\n---\n\n- Bug: two\n")
        (repo / "CHANGELOG.md").write_text("## v0.0.43\n")
        subprocess.run(["git", "init", "-q"], cwd=repo, check=True)
        subprocess.run(["git", "config", "user.email", "t@t"], cwd=repo, check=True)
        subprocess.run(["git", "config", "user.name", "t"], cwd=repo, check=True)
        subprocess.run(["git", "add", "-A"], cwd=repo, check=True)
        subprocess.run(["git", "commit", "-qm", "init"], cwd=repo, check=True)
        subprocess.run(["git", "tag", "v0.0.43"], cwd=repo, check=True)
        r = subprocess.run([sys.executable, str(SCRIPTS / "release-prep.py")], cwd=repo, capture_output=True, text=True)
        text = (repo / "CHANGELOG.md").read_text()
        check("combined: exits 0", r.returncode == 0, r.stderr)
        check("combined: one patch bump", "new=v0.0.44" in r.stdout)
        check("combined: includes every line", "Chore: one" in text and "Bug: two" in text)


def test_rejects_non_v0_0_tag():
    with tempfile.TemporaryDirectory() as tmp:
        change = "---\nversion: patch\n---\n\n- Bug: x\n"
        r = run_release_prep(tmp, change, "## v0.1.0\n", ["v0.1.0"])
        check("pre-v0.1: rejects v0.1 tag", r.returncode != 0, r.stderr)


def test_workflow_contract():
    expected = ["main-push.yml", "pull-request-check.yml"]
    for name in expected:
        check(f"workflow {name} exists", (WORKFLOWS / name).exists())
    check("workflow pull_request.yml removed", not (WORKFLOWS / "pull_request.yml").exists())
    check("workflow release-prep.yml removed", not (WORKFLOWS / "release-prep.yml").exists())
    check("workflow release.yml removed", not (WORKFLOWS / "release.yml").exists())

    for name in expected:
        text = (WORKFLOWS / name).read_text()
        check(f"workflow {name} valid yaml", yaml_ok(text), name)

    main_push = (WORKFLOWS / "main-push.yml").read_text()
    check("main-push runs make test before change processing",
          main_push.index("make test") < main_push.index("release-prep.py"))
    check("main-push processes change files on main", "release-prep.py" in main_push)
    check("main-push retries stale pushes", "git fetch origin main --tags" in main_push and "seq 1 5" in main_push)
    check("main-push serialized globally", "group: main-push" in main_push)

    pr_check = (WORKFLOWS / "pull-request-check.yml").read_text()
    check("pull-request-check validates change file", ".changes" in pr_check)
    check("pull-request-check isolates PR change file", "git merge-base HEAD origin/main" in pr_check)
    check("pull-request-check runs make test", "make test" in pr_check)


def yaml_ok(text):
    try:
        import yaml
        yaml.safe_load(text)
        return True
    except Exception:
        return False


def main():
    test_patch_path()
    test_none_path()
    test_idempotent_rerun()
    test_failure_paths()
    test_combined_changes()
    test_rejects_non_v0_0_tag()
    test_workflow_contract()
    print(f"==> {PASS} passed, {FAIL} failed")
    sys.exit(1 if FAIL else 0)


if __name__ == "__main__":
    main()
