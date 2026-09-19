# Database starting point (plugin-db / docs)

- Area: docs + plugins
- Phase: v0.x
- Goal: Give a user who picks a database a documented starting point. Today
  `plugin-db` is listed as "future" and the docs offer no path.
- Context: an external test reported that `modernc.org/sqlite` failed on arm64
  with build-constraint errors and that CGO `github.com/mattn/go-sqlite3`
  needed `libsqlite3-dev`. Dreego itself has no database opinion, which is a
  defensible choice, but the entry cost is high.
- Proposal: add a short "Using a database" guide with a working driver recipe
  (including the arm64/CGO trade-off and a container note), and decide whether a
  first-party `plugin-db` is worth a concrete proposal.
- Acceptance: docs contain a copy-paste database recipe that builds on arm64;
  plugin-db scope is either specified or explicitly rejected.
- Dependencies: none.
