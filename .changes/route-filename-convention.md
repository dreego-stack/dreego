---
version: patch
---

- Breaking: map only `+page.dreego` and `index.dreego` to their directory URL; all other route filenames now become literal URL segments and default to GET.
- Breaking: remove the former `get.dreego`, `post.dreego`, `put.dreego`, and `delete.dreego` method-file convention; HTTP methods must now be declared with method-specific sections.
- Fix: reject directories containing both `+page.dreego` and `index.dreego` with a duplicate-route error.
