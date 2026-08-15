---
version: patch
---

- Chore: split core/lexer.go into focused files (lexer_tag.go, lexer_component.go); no behavior change
- Feat: add lexer and parser fuzz targets (FuzzLexer, FuzzParser, FuzzParserPreservesGoSection) checking crashes, determinism, bounded output, and source preservation
