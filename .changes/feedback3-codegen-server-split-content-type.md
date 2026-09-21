---
version: patch
---

- Bug: `<server>` sections that mix Go declarations and statements now compile; the leading declaration block (type/func/var/const) and any top-level func are emitted at package level and the remaining statements stay inside the render function, instead of emitting the whole section at package level
- Bug: declarations at the top of a `<server>` section are hoisted to package level, so route files in one directory can share types, consts, funcs, and stores
- Bug: request-local `var` declarations that follow a statement stay inside the render function, so they are not turned into shared package state
- Bug: the generated GET handler no longer overwrites a `Content-Type` already written by a `<server type="custom">` route
