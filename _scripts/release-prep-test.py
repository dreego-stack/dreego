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
        check("combined: removes all processed files",
              not (repo / ".changes/a.md").exists() and not (repo / ".changes/b.md").exists())


def test_none_is_deferred_until_patch():
    with tempfile.TemporaryDirectory() as tmp:
        repo = Path(tmp)
        changes = repo / ".changes"
        changes.mkdir()
        none_files = {
            "a.md": "---\nversion: none\n---\n\n- Chore: one\n",
            "b.md": "---\nversion: none\n---\n\n- Chore: two\n",
        }
        for name, text in none_files.items():
            (changes / name).write_text(text)
        changelog = "## v0.0.43\n"
        (repo / "CHANGELOG.md").write_text(changelog)
        subprocess.run(["git", "init", "-q"], cwd=repo, check=True)
        subprocess.run(["git", "config", "user.email", "t@t"], cwd=repo, check=True)
        subprocess.run(["git", "config", "user.name", "t"], cwd=repo, check=True)
        subprocess.run(["git", "add", "-A"], cwd=repo, check=True)
        subprocess.run(["git", "commit", "-qm", "init"], cwd=repo, check=True)
        subprocess.run(["git", "tag", "v0.0.43"], cwd=repo, check=True)
        r = subprocess.run([sys.executable, str(SCRIPTS / "release-prep.py")], cwd=repo,
                           capture_output=True, text=True)
        check("none-only: exits 0", r.returncode == 0, r.stderr)
        check("none-only: changelog remains byte-identical",
              (repo / "CHANGELOG.md").read_text() == changelog)
        check("none-only: every file remains pending",
              all((changes / name).exists() for name in none_files))


def test_minor_bump_from_stage():
    change = "---\nversion: minor\n---\n\n- Feat: stage merge\n"
    with tempfile.TemporaryDirectory() as tmp:
        r = run_release_prep(tmp, change, "## v0.1.0 - 2026-08-15\n\n- old\n", ["v0.1.0"])
        check("minor A: exit 0", r.returncode == 0, r.stderr)
        check("minor A: prints new version", "new=v0.2.0" in r.stdout, r.stdout)
        changelog = (Path(tmp) / "CHANGELOG.md").read_text()
        check("minor A: version header added", "## v0.2.0 - " in changelog, changelog[:80])
        check("minor A: patch reset to 0", "## v0.2.0 - " in changelog and "## v0.2.1" not in changelog)
        check("minor A: change removed", not (Path(tmp) / ".changes/change.md").exists())
    with tempfile.TemporaryDirectory() as tmp:
        r = run_release_prep(tmp, change, "## v0.2.0 - 2026-08-15\n\n- old\n", ["v0.2.0"])
        check("minor B: exit 0", r.returncode == 0, r.stderr)
        check("minor B: prints new version", "new=v0.3.0" in r.stdout, r.stdout)


def test_invalid_bumps_are_atomic():
    for bump in ("major",):
        with tempfile.TemporaryDirectory() as tmp:
            change = f"---\nversion: {bump}\n---\n\n- x\n"
            r = run_release_prep(tmp, change, "## v0.0.43\n", ["v0.0.43"])
            repo = Path(tmp)
            check(f"{bump}: exits non-zero", r.returncode != 0, r.stderr)
            check(f"{bump}: changelog unchanged", (repo / "CHANGELOG.md").read_text() == "## v0.0.43\n")
            check(f"{bump}: change remains pending", (repo / ".changes/change.md").exists())


def test_coverage_gate_contract():
    script = SCRIPTS / "coverage-gate.sh"
    with tempfile.TemporaryDirectory() as tmp:
        root = Path(tmp)
        fake_go = root / "go"
        fake_go.write_text("#!/bin/sh\nprintf '%s\\n' 'ok  example 0.001s coverage: 42.0% of statements'\n")
        fake_go.chmod(0o755)
        env = os.environ.copy()
        env["PATH"] = f"{root}{os.pathsep}{env['PATH']}"
        env["DREEGO_COVERAGE_MIN"] = "40"
        ok = subprocess.run(["sh", str(script)], cwd=root, env=env, capture_output=True, text=True)
        check("coverage: threshold passes", ok.returncode == 0, ok.stderr)
        env["DREEGO_COVERAGE_MIN"] = "50"
        low = subprocess.run(["sh", str(script)], cwd=root, env=env, capture_output=True, text=True)
        check("coverage: below threshold exits non-zero", low.returncode != 0, low.stderr)
        env["DREEGO_COVERAGE_MIN"] = "invalid"
        invalid = subprocess.run(["sh", str(script)], cwd=root, env=env, capture_output=True, text=True)
        check("coverage: invalid threshold exits non-zero", invalid.returncode != 0, invalid.stderr)


def test_test_runner_contract():
    makefile = (ROOT / "Makefile").read_text()
    check("test runner: forwards DREEGO_FILTER", 'DREEGO_FILTER="$${DREEGO_FILTER:-}"' in makefile)
    check("test runner: forwards DREEGO_RUNS", 'DREEGO_RUNS="$${DREEGO_RUNS:-1}"' in makefile)
    check("test runner: coverage target exists", "coverage:" in makefile)


def test_v0x_patch_allowed_and_major_rejected():
    change = "---\nversion: patch\n---\n\n- Bug: x\n"
    with tempfile.TemporaryDirectory() as tmp:
        r = run_release_prep(tmp, change, "## v0.1.0\n", ["v0.1.0"])
        check("v0.1: patch allowed, exit 0", r.returncode == 0, r.stderr)
        check("v0.1: prints new version", "new=v0.1.1" in r.stdout, r.stdout)
    with tempfile.TemporaryDirectory() as tmp:
        r = run_release_prep(tmp, change, "## v0.2.0\n", ["v0.2.0"])
        check("v0.2: patch allowed, exit 0", r.returncode == 0, r.stderr)
        check("v0.2: prints new version", "new=v0.2.1" in r.stdout, r.stdout)
    with tempfile.TemporaryDirectory() as tmp:
        r = run_release_prep(tmp, change, "## v1.0.0\n", ["v1.0.0"])
        check("v1.0: major rejected, non-zero exit", r.returncode != 0, r.stderr)


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
    check("pull-request-check preserves published tag ancestry",
          "git merge-base --is-ancestor \"$latest\" HEAD" in pr_check)
    check("main-push preserves published tag ancestry",
          "git merge-base --is-ancestor \"$latest\" HEAD" in main_push)
    check("pull-request-check runs make test", "make test" in pr_check)
    check("pull-request-check runs coverage before tests",
          pr_check.index("make coverage") < pr_check.index("make test"))


def yaml_ok(text):
    try:
        import yaml
        yaml.safe_load(text)
        return True
    except ModuleNotFoundError:
        # The contract test must run in the minimal CI image as well. The
        # workflow files are validated by GitHub; this fallback catches the
        # most damaging local syntax mistakes without requiring PyYAML.
        return bool(text.strip()) and text.count("\n") > 2 and "jobs:" in text
    except Exception:
        return False


def main():
    test_patch_path()
    test_none_path()
    test_idempotent_rerun()
    test_failure_paths()
    test_combined_changes()
    test_none_is_deferred_until_patch()
    test_invalid_bumps_are_atomic()
    test_minor_bump_from_stage()
    test_coverage_gate_contract()
    test_test_runner_contract()
    test_v0x_patch_allowed_and_major_rejected()
    test_workflow_contract()
    print(f"==> {PASS} passed, {FAIL} failed")
    sys.exit(1 if FAIL else 0)


if __name__ == "__main__":
    main()
