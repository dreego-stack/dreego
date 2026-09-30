---
version: patch
---

- Feat: a website root's `routes/`, `static/`, `layouts/`, and `components/` are now global defaults inherited by every app; an app-local file with the same relative path shadows only that file
- Feat: a root-level `routes/` tree is accepted as the shared route default (no longer a legacy-layout error), so an app without its own `routes/` serves the global routes
- Feat: app-local static files win over the shared `static/` file with the same relative path; other global static files stay available
- Feat: `LAYOUT "path"` selects a layout by path; a `layouts/` directory may hold explicitly named layouts beside the default
- Fix: global route-local layouts under `routes/<sub>/layouts/` resolve as the `root/<sub>` scope and no longer collide
