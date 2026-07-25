
---
type: Decision
title: TypeScript auf V2 verschoben
description: TypeScript-Unterstützung im script-Block wird auf V2 verschoben, V1 nutzt Vanilla JS
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# TypeScript auf V2 verschoben

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

In den `.dreego`-Dateien gibt es einen `<script>`-Block für Client-seitigen Code. Die Frage war, ob dieser TypeScript oder Vanilla JavaScript sein soll.

## Problem

TypeScript in V1 einzubauen hätte enorme Komplexität verursacht:

1. Benötigt JS/TS-Toolchain (esbuild, Node, Bun oder Deno) auf dem Entwickler-PC
2. Zerstört das Go-Versprechen: "Keine externen Dependencies, nur Go installieren"
3. Erhöht die Build-Komplexität massiv
4. Scope Creep — verzögert V1 signifikant

## Lektion aus der Vergangenheit

Elixirs Versuch, nachträglich Types einzuführen, zeigt, wie schwer das ist: Eine dynamische Sprache nach 10 Jahren typisiert zu machen, ist ein massiver, schmerzhafter Prozess. Das ist keine Erfolgsgeschichte — es ist eine Warnung.

**Wichtig für Dreego:** TypeScript im `<script>`-Block ist unkritisch auf V2 zu verschieben, weil:
1. Go selbst ist bereits statisch typisiert (Server-Seite von Tag 1 typisiert)
2. Der `<script>`-Block ist isoliert — er beeinflusst nicht die Core-Architektur
3. Vanilla JS → TypeScript ist ein Upgrade, kein fundamentales Redesign

Aber für ANDERE Entscheidungen gilt: Was wir jetzt in der Architektur festlegen, sollte so designt sein, dass spätere Erweiterungen ohne Breaking Changes möglich sind (Plugin-Interface, Transpiler-Pipeline, Addon-System).

## Entscheidung

**V1: Reines Vanilla JavaScript** im `<script>`-Block. Kein Compiler, kein Bundler, 0 MB Extra-Tooling.

**V2: TypeScript** via esbuild-Integration (esbuild lässt sich als Go-Bibliothek einbinden).

## Begründung

1. Modernes Vanilla JS kann bereits alles Nötige: `import/export`, `async/await`, Shadow DOM
2. Dreego muss JS in V1 nur 1:1 extrahieren und in HTML einbetten
3. 0 Komplexität, maximale Render-Geschwindigkeit
4. Fokus auf den Kern: Transpiler, Routing, Go-Server

## Vorsorge für V2 (Architektur-Design jetzt)

Damit TypeScript später ohne Breaking Changes hinzukommt:

- Der Transpiler hat eine klar definierte Pipeline: Parse → AST → CodeGen
  - Ein `ScriptProcessor`-Interface kann später TS-Compiler statt No-Op sein
- Der `<script>`-Block unterstützt bereits `lang="ts"` als reserviertes Attribut (wird in V1 ignoriert, in V2 aktiviert)
- Types-Sharing (Go-Struct → TS-Interface) ist konzeptionell vorgezeichnet
- `dreego.config.json` hat einen `typescript`-Block (in V1 leer, in V2 befüllt)

## Konsequenzen

- `<script>`-Block erwartet reines JavaScript (kein `lang="ts"` in V1)
- `lang="ts"` wird geparst aber nicht verarbeitet — kein Error, nur Warnung im Dev-Server
- Kein esbuild-Dependency in V1
- V2-Planung: esbuild als Go-Binding für TS-Transpilation und JS-Bundling
- Types-Sharing (Go-Struct → TS-Interface) kommt ebenfalls erst in V2
