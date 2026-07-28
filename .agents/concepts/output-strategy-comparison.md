
---
type: Concept
title: "Comparison: Two Output Strategies for dreego generate"
description: "Per-directory dree.go vs output/ directory with index.json"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Comparison: Two Output Strategies for dreego generate

**Date:** 2026-07-28

## Approach A (implemented): Per-Directory dree.go + gen/routes.go

```
dreego/
├── routes/
│   ├── get.dreego              ← User writes
│   ├── dree.go                  ← generated, gitignored
│   ├── about/
│   │   ├── get.dreego
│   │   └── dree.go              ← generated, gitignored
│   └── users/_id_/
│       ├── get.dreego
│       └── dree.go              ← generated, gitignored
├── gen/
│   └── routes.go                ← generated, committed
│         imports all route packages (init()-based)
│
main.go imports ONLY _ "myapp/dreego/gen"
```

**Cache:** Hash comment in dree.go header line: `// hash:{bin:"...", get:"..."}`

**Pro:**
- Go-idiomatic: each directory = one package
- Incremental compilation: only changed packages are recompiled
- Simple cache check (read first line of dree.go)
- No extra metadata file
- main.go = 1 import for all routes

**Con:**
- dree.go files sit next to .dreego source files
- routes/ directory contains generated files

## Approach B (user suggestion): output/ Directory + index.json

```
dreego/
├── routes/
│   ├── get.dreego              ← User writes (ONLY source)
│   ├── about/
│   │   └── get.dreego          ← ONLY source
│   └── users/_id_/
│       └── get.dreego          ← ONLY source
├── output/
│   ├── index.json               ← generated, committed
│   ├── a1b2c3d4e5f6.go          ← generated (hash filename), gitignored
│   ├── b2c3d4e5f6a7.go          ← generated
│   └── ...
├── gen/
│   └── routes.go                ← generated, committed
│         imports _ "myapp/dreego/output"  (ONE import)
│
main.go imports ONLY _ "myapp/dreego/gen"
```

**index.json:**
```json
{
  "binary": "4354300e1c92",
  "entries": [
    {
      "source": {"path": "dreego/routes/get.dreego", "hash": "961d4658..."},
      "output": {"hash": "dba2aa9bcfbd587af...", "path": "output/a1b2c3d4e5f6.go"}
    }
  ]
}
```

**Cache Logic:**
1. Load index.json
2. For each .dreego file: compute hash, look up in index.json
3. Hash matches → skip
4. Hash new/changed → regenerate, update index
5. Binary hash changed → regenerate ALL
6. Entries in index without source file → orphan cleanup (delete)

**Pro:**
- routes/ directory stays 100% pure (only .dreego + optional _middleware.go)
- Clean separation: source vs generated code
- index.json = clear, versionable cache metadata
- Orphan detection: deleted .dreego → generated file is cleaned up
- No hash comments in Go files
- Easy reset: rm -rf output/

**Con:**
- ALL files in one package (output/) → namespace sharing
- No incremental compilation: every change = entire package recompiled
- With 100+ routes: go build compiles EVERYTHING, not just changed things
- Hash filenames are not descriptive (debugging)
- Extra index.json to maintain
- output/ directory must be cleaned up

## Decisive Difference: Go Build Performance

| Aspect | Approach A (per-dir) | Approach B (output/) |
|--------|---------------------|---------------------|
| Packages | N packages (1 per route directory) | 1 package (output/) |
| Increm. Build | Only changed package compiled | Entire package recompiled |
| 5 Routes | ~equally fast | ~equally fast |
| 100 Routes | Only changed routes recompiled | ALL 100 recompiled |
| Namespace Conflicts | Isolated per directory | Global in output/ package |
| git diff (index.json) | None | index.json changes on every generate |

## Recommendation

Approach A for projects with >10 routes. Approach B is conceptually cleaner (source/generated separation), but the performance penalty in the Go build pipeline is significant for growing projects.

Alternative: Combine both — output/ directory, but with subdirectories per route:
```
output/routes/a1b2.go     (package routes)
output/about/b2c3.go      (package about)
```
This gives source separation + incremental compilation. The index.json manages the mapping.
