# DreeJS typed client props

- Area: compiler
- Phase: DreeJS
- Goal: serialize only explicitly declared component data across the server-to-client boundary.
- Acceptance: supported Go values retain generated TypeScript types, HTML and script contexts are escaped safely, unsupported or secret-bearing values fail with source diagnostics, representations for time, numbers, optionals, and errors are documented, and payload minimization has regression tests.
- Depends on: dreejs-lifecycle.1
