# Entscheidung: Kein `<catch>`-Tag — Fehler via Go-Idiome

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

Ursprünglich wurde ein `<catch>`-Block vorgeschlagen, um Fehler im Template abzufangen (inspiriert von JavaScripts try/catch).

## Problem

`<catch>` ist ein Konzept aus Sprachen mit Exceptions (JavaScript, Java) und fühlt sich in der Go-Welt fremd an.

Go-Entwickler behandeln Fehler explizit (`if err != nil`). Ein separates Error-Handling-Tag würde gegen Go-Idiome verstoßen und die Lernkurve für Go-Entwickler unnötig erhöhen.

## Entscheidung

**Kein `<catch>`-Tag.** Fehler werden im `<go>`-Block behandelt und als Variablen ans Template übergeben:

```html
<go>
    user, err := db.GetUser(id)
    hasError := err != nil
</go>

<div class="profile">
    {#if hasError}
        <p class="error">Benutzer konnte nicht geladen werden.</p>
    {#else}
        <h1>Hallo, {user.Name}!</h1>
    {/if}
</div>
```

## Begründung

1. Hält die Template-Sprache extrem schlank
2. Go-Entwickler müssen nichts Neues lernen
3. Keine Framework-Magie — nur Variablen und `{#if}`
4. Explizite Fehlerbehandlung ist idiomatisches Go

## Konsequenzen

- Template-Logik (`{#if}`) übernimmt auch Error-Rendering
- Entwickler müssen Fehler selbst behandeln (wie in Go üblich)
- Kein implizites Error-Handling im Framework
