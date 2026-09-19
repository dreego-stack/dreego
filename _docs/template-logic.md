# Template Logic

Dreego templates contain a deliberately small control-flow language. Complex
work remains ordinary Go in the matching `<server>` section; template constructs
compile to generated Go and add no browser runtime.

## Expressions

Use double braces in text and attributes:

```html
<h1>{{ user.Name }}</h1>
<a href="{{ profileURL }}">Profile</a>
```

Dreego escapes expressions for their output context. Text, quoted attributes,
URLs, event attributes, and inline style values do not share one generic
escaping rule. See [Output Safety](security.md).

Filters form a pipe-separated chain:

```html
<p>{{ name|upper }}</p>
<div>{{ trustedHTML|raw }}</div>
```

`raw` bypasses normal escaping and must receive application-trusted HTML.

## Conditions

```html
{#if user.IsAdmin}
    <a href="/admin">Admin</a>
{#else if user.IsMember}
    <a href="/account">Account</a>
{#else}
    <a href="/login">Sign in</a>
{/if}
```

`{#elseif condition}` is accepted as an alternative spelling. Conditions are
Go expressions using values available to the template.

## Iteration

```html
{#each users as user}
    <p>{{ $loop.Index }}: {{ user.Name }}</p>
{#each else}
    <p>No users</p>
{/each}
```

`$loop.Index` is zero-based. `$loop.First`, `$loop.Last`, `$loop.Even`, and
`$loop.Odd` describe the current iteration. The optional `{#each else}` branch
renders when the collection is empty.

## Components and slots

```html
<@Card title="Welcome">
    {#slot header}<strong>News</strong>{/slot}
    <p>Default slot content</p>
</@Card>
```

Inside the component, `{#slot}` renders default content and
`{#slot header}{/slot}` marks a named slot. Slots are passed as render inputs;
they are not stored in request-global string keys. See [Components](components.md).

## Verbatim regions

```html
{#verbatim}
<template>{{ handled_by_another_tool }}</template>
{/verbatim}
```

Verbatim content is emitted literally and receives no Dreego interpolation or
escaping. Use it only for developer-authored content.

## Boundaries

Control-flow blocks cannot begin inside an HTML attribute. Wrap the complete
element in the condition or loop instead. Dreego has no catch, await, switch,
or template-local variable blocks; implement that logic in Go.
