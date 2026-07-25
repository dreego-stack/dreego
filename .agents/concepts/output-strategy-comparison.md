
---
type: Concept
title: "Vergleich: Zwei Output-Strategien für dreego generate"
description: "Per-Directory dree.go vs output/-Verzeichnis mit index.json"
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Vergleich: Zwei Output-Strategien für dreego generate

**Datum:** 24.07.2026

## Ansatz A (implementiert): Per-Directory dree.go + gen/routes.go

```
dreego/
├── routes/
│   ├── get.dreego              ← User schreibt
│   ├── dree.go                  ← generiert, gitignored
│   ├── about/
│   │   ├── get.dreego
│   │   └── dree.go              ← generiert, gitignored
│   └── users/_id_/
│       ├── get.dreego
│       └── dree.go              ← generiert, gitignored
├── gen/
│   └── routes.go                ← generiert, committed
│         importiert alle Route-Pakete (init()-basiert)
│
main.go importiert NUR _ "myapp/dreego/gen"
```

**Cache:** Hash-Kommentar in dree.go Kopfzeile: `// hash:{bin:"...", get:"..."}`

**Pro:**
- Go-idiomatisch: jedes Verzeichnis = ein Package
- Inkrementelle Kompilierung: nur geanderte Packages werden neu kompiliert
- Einfache Cache-Prufung (erste Zeile von dree.go lesen)
- Keine extra Metadaten-Datei
- main.go = 1 Import fur alle Routen

**Contra:**
- dree.go-Dateien liegen neben .dreego-Source-Dateien
- routes/-Verzeichnis enthalt generierte Dateien

## Ansatz B (User-Vorschlag): output/ Verzeichnis + index.json

```
dreego/
├── routes/
│   ├── get.dreego              ← User schreibt (NUR Source)
│   ├── about/
│   │   └── get.dreego          ← NUR Source
│   └── users/_id_/
│       └── get.dreego          ← NUR Source
├── output/
│   ├── index.json               ← generiert, committed
│   ├── a1b2c3d4e5f6.go          ← generiert (Hash-Dateiname), gitignored
│   ├── b2c3d4e5f6a7.go          ← generiert
│   └── ...
├── gen/
│   └── routes.go                ← generiert, committed
│         importiert _ "myapp/dreego/output"  (EIN Import)
│
main.go importiert NUR _ "myapp/dreego/gen"
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

**Cache-Logik:**
1. Lade index.json
2. Fur jede .dreego-Datei: Hash berechnen, in index.json nachschlagen
3. Hash matcht → skip
4. Hash neu/geandert → regenerieren, index updaten
5. Binary-Hash geandert → ALLE regenerieren
6. Eintrage in index ohne Source-Datei → Orphan cleanup (löschen)

**Pro:**
- routes/-Verzeichnis bleibt 100% pur (nur .dreego + optionale _middleware.go)
- Saubere Trennung: Source vs generierter Code
- index.json = klare, versionierbare Cache-Metadaten
- Orphan-Detection: geloschte .dreego → generierte Datei wird aufgeraumt
- Keine Hash-Kommentare in Go-Dateien
- Einfacher Reset: rm -rf output/

**Contra:**
- ALLE Dateien in einem Package (output/) → Namespace-Sharing
- Keine inkrementelle Kompilierung: jede Anderung = gesamtes Package neu
- Bei 100+ Routen: go build kompiliert ALLES, nicht nur Geandertes
- Hash-Dateinamen sind nicht aussagekraftig (Debugging)
- Extra index.json zu pflegen
- Output/-Verzeichnis muss aufgeraumt werden

## Entscheidender Unterschied: Go Build Performance

| Aspekt | Ansatz A (per-dir) | Ansatz B (output/) |
|--------|---------------------|---------------------|
| Packages | N Packages (1 pro Route-Verzeichnis) | 1 Package (output/) |
| Inkrem. Build | Nur geandertes Package kompiliert | Gesamtes Package neu |
| 5 Routen | ~gleich schnell | ~gleich schnell |
| 100 Routen | Nur geanderte Routen neu | ALLE 100 neu kompiliert |
| Namespace-Konflikte | Pro Verzeichnis isoliert | Global im output/-Package |
| git diff (index.json) | Keine | index.json andert sich bei jedem generate |

## Empfehlung

Ansatz A fur Projekte mit >10 Routen. Ansatz B ist konzeptionell sauberer (Source/Generated-Trennung), aber der Performance-Nachteil im Go-Build-Pipeline ist signifikant fur wachsende Projekte.

Alternative: Beide kombinieren — output/-Verzeichnis, aber mit Unterverzeichnissen pro Route:
```
output/routes/a1b2.go     (package routes)
output/about/b2c3.go      (package about)
```
Das gibt Source-Trennung + inkrementelle Kompilierung. Die index.json managed die Zuordnung.
