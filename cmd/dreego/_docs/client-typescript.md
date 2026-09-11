# TypeScript Client Code

TypeScript is supported in a root client section and in inline body scripts.

```html
<client lang="ts">
const message: string = "ready"
document.title = message
</client>
```

```html
<body>
<button id="save" type="button">Save</button>
<script lang="ts">
const button = document.querySelector<HTMLButtonElement>("#save")
button?.addEventListener("click", () => { button.disabled = true })
</script>
</body>
```

`dreego generate` type checks and compiles each section to JavaScript. Errors
refer back to the `.dreego` source location. Generation never downloads tools
implicitly; install the pinned compiler explicitly:

```text
dreego tools install typescript
```

Generated Go model declarations available to the client are checked together
with the TypeScript source where Dreego owns that schema. Browser APIs retain
their normal TypeScript DOM types.

The exact compiler, strictness, and JavaScript output targets are versioned in
[Client Languages](client-languages.md).
