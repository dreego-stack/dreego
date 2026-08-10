---
version: none
---

- Bug: guard `process` job in `pull-request-merging` to only run on merged PRs (or pushes to main), preventing failures when a PR is closed without merging
