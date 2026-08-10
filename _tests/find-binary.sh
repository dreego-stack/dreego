#!/bin/sh
# Detect binary files in the repo. Every repo file should be readable as text;
# only allowlisted binary types (images, fonts, archives) are permitted.
# Pure filesystem scan — no git involved. A file is binary when it contains a
# NUL byte (portable text/binary heuristic, independent of grep -I windowing).
set -e

# cd from _tests/ to the repo root.
cd "$(dirname "$0")/.."

# Allowlisted binary/image extensions (lowercased). Add types here when a
# binary file is intentionally part of the repo.
ALLOWED_BIN_EXT="svg png jpg jpeg gif ico webp woff woff2 ttf eot pdf zip gz tar"

fail=0

while IFS= read -r file; do
    case "$file" in
        */.DS_Store)
            continue
            ;;
    esac

    ext=$(echo "$file" | sed 's/.*\.//' | tr '[:upper:]' '[:lower:]')
    case " $ALLOWED_BIN_EXT " in
        *" $ext "*)
            continue
            ;;
    esac

    # Empty files (e.g. .gitkeep) are text.
    if [ ! -s "$file" ]; then
        continue
    fi

    # Binary = contains at least one NUL byte. Compare the raw size against
    # the size after stripping NULs; a difference means a NUL was present.
    total=$(wc -c < "$file" | tr -d ' ')
    without_nul=$(tr -d '\000' < "$file" | wc -c | tr -d ' ')
    if [ "$total" != "$without_nul" ]; then
        echo "-> FAIL -> Found binary $file"
        fail=1
    fi
done < <(find . \
    \( -path ./.git -o -path ./.kilo -o -path ./.tmp \) -prune -o \
    -type f -print 2>/dev/null)

if [ "$fail" -ne 0 ]; then
    echo "FAIL: binary files found in repo"
    exit 1
fi

echo "PASS: no binary files in repo"
