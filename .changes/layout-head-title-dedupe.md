---
version: none
---

- Bug: dedupe title and meta description in the head merge for body-level layouts, so a route `<title>` wins over the layout title in both layout shapes and the rendered page always contains exactly one `<title>`.
