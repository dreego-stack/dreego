# YAML Frontmatter

A `.dreego` file may start with an optional YAML-like frontmatter block delimited by `---` on its own line, both opening and closing. Dreego parses it and exposes the `key: value` pairs as typed metadata.

## Syntax

A frontmatter block sits at the very top of the file:

```dreego
---
title: My Page
author: Lukas
tags: [go, web]
---

<div>
  <h1>{title}</h1>
</div>
```

Rules:

- The opening `---\n` must be the very first characters of the file.
- The block is closed by a `---` on its own line.
- Each line is `key: value`. The **first** `:` splits key from value, so a `:` inside the value is preserved (e.g. `url: https://example.com`).
- Quoted values are unwrapped (`title: "My Page"` → `My Page`).
- List values `tags: [go, web]` are normalized to a comma-joined string (`go, web`).
- Blank lines and lines starting with `#` are ignored.
- A file without a leading `---` block has no frontmatter.

## API

`dreego.ParseFrontmatter(src)` splits the source and returns the metadata plus the remaining body:

```go
fm, body := dreego.ParseFrontmatter(src)
// fm   → map[string]string{"title": "My Page", ...}
// body → the .dreego source without the frontmatter block
```

If there is no leading frontmatter, `fm` is `nil` and `body` is the whole `src`. A block without a closing `---` yields `nil` too.

## See Also

- [Docs Index](https://codeberg.org/dreego/dreego/src/branch/main/_docs/index.md)
- [Getting Started](https://codeberg.org/dreego/dreego/src/branch/main/_docs/getting-started.md)
