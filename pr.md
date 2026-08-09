---
version: none
---

- Fix: move pr.md processing from pre-merge (pull_request.yml) to post-merge (merging.yml)
- Fix: prepare job no longer auto-commits to PR branch (local and origin stay in sync)
- Fix: parallel PRs no longer conflict on VERSION (version computed at merge time)
- Feat: merging.yml now processes pr.md after merge: changelog + version + embedded sync + tag + branch delete
- Chore: pull_request.yml is check-only (tests + pr.md validation), no prepare job