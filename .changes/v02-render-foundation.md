---
version: minor
---

- Feat: target-neutral typed component rendering with RenderContext, dreego.Render, and structured Result
- Feat: explicit SSR lifecycle host at core/ssr — generated apps start via ssr.Listen and Host owns server limits and shutdown
- Feat: transpiler IR matrix — ir/ model, html input folders (html/html, html/css, html/head) with single output stage, js/js client layer, tokens/lexer/parser subpackages
- Feat: stage model — version minor for stage/* merges, tagged automatically after merge
- Docs: first-party processors decision (md→html, ts→js), v0.1-to-v0.2 migration guide
- Breaking: App server lifecycle moved to core/ssr Host; generated components now return Result — see _docs/migration-v0.1-to-v0.2.md
