# Semantic Sections Migration

Dreego renamed three root sections before v0.1 so section names describe their
purpose instead of their implementation language.

## Required changes

| Legacy root section | Current root section | Default language |
| --- | --- | --- |
| `<go>...</go>` | `<server>...</server>` | Go |
| `<div>...</div>` | `<body>...</body>` | HTML |
| `<script>...</script>` | `<client>...</client>` | JavaScript |

`<head>` and `<style>` keep their names. There are no legacy aliases: update
the source and regenerate it with the matching Dreego CLI version.

Before:

```html
<go>message := "Hello"</go>
<div><h1>{{ message }}</h1></div>
<script>console.log("ready")</script>
```

After:

```html
<server>message := "Hello"</server>
<body><h1>{{ message }}</h1></body>
<client>console.log("ready")</client>
```

## Language attributes

Built-in languages can be stated explicitly:

```html
<server lang="go"></server>
<head lang="html"></head>
<body lang="html"></body>
<style lang="css"></style>
<client lang="js"></client>
```

The attributes may be omitted for these defaults. Other section/language pairs
fail during generation until a compatible language processor is installed.

## HTML scripts inside the body

An HTML `<script>` inside `<body>` remains normal HTML and is not renamed:

```html
<body>
    <script type="application/ld+json">{"name":"Dreego"}</script>
</body>
```

Only a root `<client>` section represents JavaScript source collected by
Dreego.

## Repository migration

Update application routes, layouts, components, tests, editor grammar, and
snippets in the same change. Run:

```text
dreego fmt --check
dreego generate --check
```

Then build and test the application. A legacy root section produces a source
location and the exact replacement. An unsupported language reports its
section and the missing processor requirement.
