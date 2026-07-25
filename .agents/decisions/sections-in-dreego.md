
---
type: Decision
title: 5 Sektionen in .dreego-Dateien
description: Dreego-Dateien haben 5 klar getrennte Sektionen für Server- und Client-Code
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# 5 Sektionen in .dreego-Dateien

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

Eine `.dreego`-Datei muss klar trennen, was auf dem Server und was im Browser läuft. Viele Frameworks scheitern an dieser Trennung.

## Entscheidung

Eine `.dreego`-Datei wird in **5 klar getrennte Sektionen** unterteilt:

1. **`<head>`** — Komponenten-spezifische Meta-Tags, Scripts, CSS-Links
2. **`<go>`** — Server-seitiger Go-Code (Daten-Fetching, Logik)
3. **Template (HTML)** — Das Markup mit Dreego-Template-Syntax
4. **`<script>`** — Client-seitiges JavaScript (V1: Vanilla JS)
5. **`<style>`** — Scoped CSS (automatisch mit Hashes)

## Begründung

1. **Keine Verwirrung:** Es ist immer klar, welcher Code wo läuft
2. **Volle Go-Power auf dem Server:** `<go>` hat DB-Zugriff, Request-Kontext, etc.
3. **Echtes JavaScript für den Browser:** `<script>` wird 1:1 an den Client geschickt
4. **Komponenten-basierte Assets:** `<head>` lädt Skripte nur, wenn die Komponente gerendert wird
5. **Scoped CSS:** `<style>` verschmutzt nicht den globalen Namespace

## Die `<head>`-Innovation

Das `<head>`-Tag ist eine Kern-Innovation für Addons und Performance:

- `dreego-map` deklariert Mapbox-Skripte nur in seinem `<head>`
- Der Dreego-Transpiler injiziert diese nur, wenn die Komponente tatsächlich gerendert wird
- Kein globales Laden schwerer Libraries für alle Seiten

## Konsequenzen

- Der Transpiler muss die 5 Sektionen parsen und trennen können
- Jede Sektion wird anders verarbeitet:
  - `<go>` → Go-Code (Server)
  - Template → Go-Code (HTML-Generierung)
  - `<style>` → Gesammelt, gescoped, in CSS-Datei
  - `<script>` → Extrahiert, in HTML eingebettet
  - `<head>` → Dynamisch in finalen HTML-Head injiziert
