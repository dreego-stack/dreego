# Style Section

`<style>` contains CSS. `lang="css"` is the default and only supported style
language.

```html
<body><article class="card">Hello</article></body>
<style>
.card { padding: 1rem; border: 1px solid currentColor; }
</style>
```

Route styles are emitted with the route output. Component styles are scoped by
generated attributes so matching selectors do not style unrelated markup on
the page. Dreego rewrites supported component selectors during generation; it
does not ship a browser CSS runtime.

Keep global design tokens and resets in normal static assets. Use component
style sections for rules that belong to one component. See
[Components](components.md) for the scoping contract and [Output Safety](security.md)
for dynamic style attributes inside body templates.

