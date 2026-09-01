---
version: minor
---

- Feat: target-neutral render foundation with non-HTTP rendering (dreego.Render, dreego.NewContext, render Result)
- Feat: explicit SSR host at core/ssr — generated apps start via ssr.Listen, HTTP-owning APIs moved from core to core/ssr
- Feat: transpiler IR matrix — ir/ model, html input folders (html/html, html/css, html/head) with single output stage, js/js client layer, tokens/lexer/parser subpackages
- Feat: stage model — version minor for stage/* merges, tagged automatically after merge
- Docs: first-party processors decision (md→html, ts→js), v0.1-to-v0.2 migration guide
- Breaking: HTTP host APIs moved to core/ssr (middleware, WithStore, BindForm, ServerConfig) — see _docs/migration-v0.1-to-v0.2.md
