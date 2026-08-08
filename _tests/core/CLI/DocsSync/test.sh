#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Verify the embedded docs in cmd/dreego/embedded/ match _docs/, README.md, and CHANGELOG.md
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

# Source dirs: repo-root docs vs. the mirror embedded into the CLI.
src_docs="$realrepo/_docs"
src_readme="$realrepo/README.md"
src_changelog="$realrepo/CHANGELOG.md"

embedded_dir="$realrepo/cmd/dreego/embedded"
emb_docs="$embedded_dir/_docs"
emb_readme="$embedded_dir/README.md"
emb_changelog="$embedded_dir/CHANGELOG.md"

fail() {
    echo "FAIL: $1"
    echo "docs out of sync, run _scripts/sync-embedded-docs.sh"
    exit 1
}

# Compare two directories: identical relative file lists AND identical content.
compare_dirs() {
    a="$1"
    b="$2"
    label="$3"

    [ -d "$a" ] || fail "$label source dir missing: $a"
    [ -d "$b" ] || fail "$label embedded dir missing: $b"

    (cd "$a" && find . -type f | sort) > "$workdir/list_a"
    (cd "$b" && find . -type f | sort) > "$workdir/list_b"

    if ! diff -q "$workdir/list_a" "$workdir/list_b" >/dev/null; then
        fail "$label file lists differ (missing or extra files)"
    fi

    while IFS= read -r rel; do
        [ -n "$rel" ] || continue
        if ! diff -q "$a/$rel" "$b/$rel" >/dev/null; then
            fail "$label content differs: $rel"
        fi
    done < "$workdir/list_a"
}

# Compare two single files: identical content.
compare_files() {
    a="$1"
    b="$2"
    label="$3"

    [ -f "$a" ] || fail "$label source file missing: $a"
    [ -f "$b" ] || fail "$label embedded file missing: $b"
    if ! diff -q "$a" "$b" >/dev/null; then
        fail "$label content differs: $b"
    fi
}

compare_dirs "$src_docs" "$emb_docs" "docs"
compare_files "$src_readme" "$emb_readme" "README"
compare_files "$src_changelog" "$emb_changelog" "CHANGELOG"

echo "PASS: embedded docs in sync with _docs/, README.md, CHANGELOG.md"
