# HTML Body Section

`<body>` is the rendered template of a route, layout, or component.
`lang="html"` is the default.

```html
<server>items := []string{"one", "two"}</server>
<body>
<main>
    {#if len(items) > 0}
        <ul>{#each items as item}<li>{{ item }}</li>{/each}</ul>
    {/if}
</main>
</body>
```

The section contains the application's body markup; Dreego preserves the body
wrapper in generated HTML. Expressions use `{{ expression }}`. Control flow,
component calls, slots, filters, and raw output are described in
[Template Logic](template-logic.md).

## Inline processors

An HTML body can contain selected nested language regions:

```html
<body>
<md># Markdown inside HTML</md>
<script lang="ts">const ready: boolean = true</script>
<script lang="lua">print("ready")</script>
</body>
```

`<md>` becomes HTML during generation. Typed and Lua scripts are compiled to
JavaScript and stay at their original body position. An ordinary nested
`<script>` remains ordinary HTML and is not treated as a root client section.

For a complete Markdown body use `<body lang="md">`; see
[Markdown Bodies and Regions](../../cmd/dreego/_docs/markdown.md).
