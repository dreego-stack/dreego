---
version: patch
---

- Feat: remove VERSION file — the latest git tag is now the single source of truth for the CLI version
- Feat: CLI version derives from git tag at build time (`-ldflags -X main.version=$(git describe --tags --abbrev=0)`) or from build info (`go install pkg@tag`)
- Fix: merging workflow creates the tag from the version computed by release-prep.py (no more VERSION file read)
- Fix: test environment injects the version via build arg (Dockerfile + Makefile test + test.sh), version-drift test compares against DREEGO_VERSION
- Chore: pull-request-check.yml fetches tags (fetch-depth: 0) so the version-drift test is meaningful