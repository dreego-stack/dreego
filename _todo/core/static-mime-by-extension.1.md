# Static files: standard extensions served as application/octet-stream

- Area: compiler
- Phase: asset system
- Goal: serve every static file whose extension the Go standard library knows with the correct Content-Type.
- Gap: `MimeByExt` (`internal/dreefile/static.go`) hard-codes a short extension list and returns `application/octet-stream` for everything else, so `.txt`, `.xml`, `.csv`, `.webmanifest` and other standard formats are mislabelled (e.g. a generated `robots.txt` or `sitemap.xml`).
- Fix: keep the pinned extensions (fonts, icons, `css`/`js`/`svg` overrides) and fall back to `mime.TypeByExtension`, which consults the built-in table first and only then the host database.
- Acceptance: `robots.txt` → `text/plain`, `sitemap.xml` → `text/xml`, unknown extension → `application/octet-stream`; covered by a compiler unit test and an integration test.
