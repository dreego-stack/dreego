---
version: patch
---

- Chore: align CI with the PR-driven release process — pull_request.yml validates pr.md, release-prep.yml applies changelog+version to the PR branch, release.yml creates the tag after merge
- Chore: serialize release-prep and release runs with concurrency groups so concurrent merges cannot race version, changelog, commit, or tag creation
- Chore: tag only a fully checked commit — release.yml runs make test on the merged commit before creating the tag
- Chore: add release-prep contract tests (patch, none, idempotent rerun, failure paths, workflow contract) wired into make test
- Chore: make release-prep.py idempotent and resolve paths from the working directory
- Chore: fix pr.md.example frontmatter so it passes pr.md validation