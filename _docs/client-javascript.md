# JavaScript Client Section

`<client>` contains browser code. `lang="js"` is the default.

```html
<body><button id="save" type="button">Save</button></body>
<client>
const button = document.querySelector("#save")
button?.addEventListener("click", () => {
    button.textContent = "Saved"
})
</client>
```

JavaScript is passed through as the configured client output language. Dreego
escapes script-end sequences when it embeds client code, but it does not type
check JavaScript or add a component runtime.

A root client section is collected as page client code. A plain `<script>`
nested in an HTML body remains in place as HTML. For checked client code use
[TypeScript](client-typescript.md); for Lua authoring use [Browser Lua](lua.md).

The current JavaScript compatibility target is listed in
[Client Languages](client-languages.md).
