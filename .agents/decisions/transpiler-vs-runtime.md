---
type: Decision
title: Compile-Time Transpiler statt Runtime-Parsing
description: Build-Zeit Code-Generation, kein Laufzeit-Parsing
tags: [transpiler, v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---

# Compile-Time Transpiler statt Runtime-Parsing

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

`.dreego`-Dateien müssen in ausführbaren Code umgewandelt werden. Zwei Ansätze stehen zur Wahl:

1. **Runtime-Parsing:** Server liest `.dreego`-Dateien zur Laufzeit (ähnlich `html/template`)
2. **Compile-Time Transpiler:** `dreego generate` wandelt `.dreego` → `.go` vor dem Build

## Entscheidung

**Compile-Time Transpiler** (Weg A aus dem Gemini-Chat).

`dreego generate` liest `.dreego`-Dateien und generiert daraus Go-Code.

## Begründung

| Kriterium              | Runtime-Parsing      | Compile-Time (gewählt)     |
|------------------------|---------------------|----------------------------|
| Performance            | Langsamer (Parsing) | Maximal (kein Laufzeit-Overhead) |
| Single Binary          | via `//go:embed`    | Alles im Binary, kein Parsing   |
| Fehlererkennung        | Zur Laufzeit (Crash)| Zur Build-Zeit (`go build` bricht ab) |
| Type-Safety            | Keine               | Volle Go-Typensicherheit   |
| DevX                   | Kein Build-Step     | `dreego generate` im Watcher |
| Debugging              | Schwer              | Normal (generierter Go-Code) |

## Konsequenzen

- Build-Step: `dreego generate` muss vor `go build` laufen
- Generierte `*_dreego.go`-Dateien werden nicht committed
- 100% Compile-Time Safety: Kein Template-Fehler erreicht Production
- Dev-Server führt `dreego generate` automatisch bei Dateiänderungen aus
- Go hat keine Makros — Code Generation ist der Go-idiomatische Weg
