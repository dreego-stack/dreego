---
version: minor
---

- Feat: Markdown body processor — `<body lang="md">` renders Markdown to HTML at generation time (stdlib-only, hand-written parser)
- Feat: protected Dreego constructs inside Markdown bodies (expressions, conditions, components stay intact)
- Feat: blog demo app with Markdown posts (demo/blog) demonstrating the new processor
- Feat: blog layout with document shell (doctype/html/head/body) + prose styling for markdown content
- Fix: layout style sections now emitted into generated output (previously dropped)
- Fix: head-merge helpers no longer duplicated with multiple route subdirectories
- Feat: demo path routing (/blog, /saas) with trailing-slash redirects
- Docs: Markdown guide (_docs/markdown.md), roadmap + ADR status updated
- Feat: markdown pipe tables, nested lists, images, and footnotes (GFM-style)
