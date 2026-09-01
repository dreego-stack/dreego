# Change files

Every pull request adds exactly one Markdown file in this directory. Use a
short, unique name such as `fix-session-path.md`:

```markdown
---
version: patch
---

- Bug: keep session cookie paths consistent
```

Use `version: none` only when the release tag must not change; the file stays
pending until a `version: patch` or `version: minor` release. Use
`version: minor` only for pull requests from `stage/*` branches that ship a
planned phase; CI tags the next minor automatically after merge. After merge,
the release workflow combines all pending files into one changelog entry and
removes them.
