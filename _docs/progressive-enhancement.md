# Progressive Enhancement

Dreego is SSR-first. Every page renders complete HTML on the server, and the
client receives finished markup. Interactivity is added on top of that HTML
with HTMX, Alpine.js, or plain JavaScript — there is no internal client
runtime before v0.1.

This guide explains the model and walks through one worked example that works
with JavaScript disabled, with plain JavaScript, and with HTMX or Alpine.js
enhancements.

## The Model

- The server renders HTML. State lives on the server (sessions, databases,
  `<go>` block variables).
- The client enhances that HTML. Nothing the client does is required for the
  page to be usable.
- HTMX swaps HTML fragments returned by the server (partial page updates
  without a full reload).
- Alpine.js handles local UI state that never needs the server (dropdowns,
  tabs, modals, toggles).
- Plain JavaScript is always available for anything the other two do not cover.
- There is no framework runtime, no hydration, and no client-side router in
  the current v0.1 core. The planned DreeJS phases later add optional modular
  local, fetch, poll, stream, and live behavior; see the roadmap.

## The Worked Example: A Comment Form

A comment form on a blog post. The requirements:

1. Works without JavaScript (full page reload).
2. Works with HTMX (AJAX submit, no full reload).
3. Validates on the server in both cases.
4. Never trusts the client.

### Step 1: The server-rendered form (no JavaScript)

```dreego
<!-- www/routes/posts/[id]/get.dreego -->
<go>
    post, err := loadPost(c.Param("id"))
    if err != nil {
        return "", err
    }
    comments, err := loadComments(post.ID)
    if err != nil {
        return "", err
    }
</go>

<div>
    <h1>{{ post.Title }}</h1>
    <ul id="comments">
        {#each comments as comment}
            <li>{{ comment.Author }}: {{ comment.Body }}</li>
        {/each}
    </ul>

    <form g-action="CreateComment" method="post">
        <input type="hidden" name="csrf_token" value="{{ c.CSRFToken() }}">
        <label for="author">Name</label>
        <input id="author" name="author" required>
        <label for="body">Comment</label>
        <textarea id="body" name="body" required></textarea>
        <button type="submit">Post comment</button>
    </form>
</div>
```

This is a plain HTML form. With JavaScript disabled it still works: the POST
goes to the server, the server validates, and the response is a full page
reload. The `required` attributes are a convenience, not a security boundary —
the server validates again.

The `g-action="CreateComment"` attribute is what makes the form work: Dreego
generates a POST handler for this route that parses the form, maps it to the
`CommentForm` struct, validates it, and calls `CreateComment`. Without
`g-action` the form is a plain HTML form and no POST handler is generated.

### Step 2: The server handler

The `g-action` handler definition lives in the POST route file (`post.dreego`) for the same URL — Dreego's method-filename routing maps `post.dreego` to the POST method on that route.

```dreego
<!-- www/routes/posts/[id]/post.dreego -->
<go>
    type CommentForm struct {
        Author string `form:"author" validate:"required,max=80"`
        Body   string `form:"body" validate:"required,max=2000"`
    }

    func CreateComment(c dreego.Context, form CommentForm) error {
        if err := storeComment(c.Param("id"), form.Author, form.Body); err != nil {
            return err
        }
        return c.Redirect("/posts/"+c.Param("id"), 303)
    }
</go>
```

The `g-action` pipeline parses the form, maps it to the struct, and validates
with the built-in validators. On failure it re-renders the page with errors;
on success it redirects (Post-Redirect-Get). CSRF is checked by middleware
before the handler runs.

### Step 3: HTMX enhancement (AJAX submit)

HTMX is loaded once, for example in the layout `<head>`:

```html
<script src="https://unpkg.com/htmx.org@2" defer></script>
```

The form gains one attribute:

```dreego
<form g-action="CreateComment" method="post"
      hx-boost="true">
```

`hx-boost` converts the form submit into an AJAX request. The server still
runs the same `g-action` pipeline: it validates, stores the comment, and
redirects with 303. HTMX follows the redirect and swaps the resulting full
page into the body — the page updates without a browser reload. The no-JS
path is unchanged: without JavaScript the form submits normally.

HTMX does **not** send the CSRF token on its own. The middleware accepts the
token either as the `X-CSRF-Token` header or as the `csrf_token` form field
(see [Forms](forms.md)). The hidden field in Step 1 covers both paths, so no
extra configuration is needed. If you prefer the header, configure it once on
the `<body>`:

```html
<body hx-headers='{"X-CSRF-Token": "{{ c.CSRFToken() }}"}'>
```

The important rule: **the server response is always valid HTML**, and the
no-JS path never depends on the enhancement.

### Step 4: Alpine.js enhancement (local UI state)

Alpine.js handles things that never need the server. For example, a
"preview" toggle next to the textarea:

```html
<script src="https://unpkg.com/alpinejs@3" defer></script>
```

```dreego
<div x-data="{ preview: false }">
    <button type="button" @click="preview = !preview">
        <span x-show="!preview">Preview</span>
        <span x-show="preview">Edit</span>
    </button>
    <textarea id="body" name="body" required x-show="!preview"></textarea>
    <div x-show="preview" x-text="document.getElementById('body').value"></div>
</div>
```

With JavaScript disabled the button does nothing and the textarea is always
visible — the form still works. Alpine.js only enhances local behavior; it
never replaces a server round-trip.

### Step 5: Plain JavaScript

Plain JavaScript is the fallback for anything HTMX and Alpine.js do not
cover. It follows the same rule: enhance, never require.

```html
<script>
    document.querySelectorAll("form[data-confirm]").forEach(function (form) {
        form.addEventListener("submit", function (event) {
            if (!window.confirm(form.dataset.confirm)) {
                event.preventDefault();
            }
        });
    });
</script>
```

## Security

- **CSRF:** Every state-changing request is checked by the CSRF middleware.
  Plain HTML forms send the token as a hidden field; HTMX requests send it
  either as the same hidden field or as the `X-CSRF-Token` header via
  `hx-headers`. Never disable CSRF for a form that changes state.
- **XSS:** Template output is HTML-escaped by default (`{{ expression }}`).
  Only `{{ expression|raw }}` bypasses escaping — use it rarely and only for
  trusted content. Never render user input with `|raw`.
- **URLs:** Escaped URL attributes still require scheme validation. Do not
  put user input into `href` or `action` without checking the scheme
  (`http`, `https`, `mailto`).
- **Validation:** Client-side validation (`required`, Alpine.js checks) is a
  UX convenience. The server validates every request with the built-in
  validators. Never trust the client.
- **CSP:** The default Content-Security-Policy allows `unsafe-inline` for
  scripts and styles so HTMX, Alpine.js, and scoped CSS work out of the box.
  If you tighten the CSP with `app.SetCSP`, you must allow the scripts you
  actually load — see [Middleware](middleware.md).
- **CDN origins:** If scripts are loaded from a CDN (e.g. `unpkg.com`), the
  CDN origin must be included in `script-src` (e.g.
  `script-src 'self' 'unsafe-inline' https://unpkg.com`).

## Failure Behavior

- **JavaScript fails to load:** The page is still fully rendered HTML. Forms
  submit with a full page reload, links navigate normally, and Alpine.js
  widgets degrade to their static content.
- **HTMX request fails:** The server returns an error page or an error
  fragment. The page is not left in a broken state — the last successful
  render stays visible, and a full reload always recovers.
- **Server error:** The App error path renders the error page. Errors are
  never disclosed to clients (no stack traces, no internal paths).
- **Network failure:** A plain form submit or HTMX request that never reaches
  the server leaves the page unchanged. The user can retry.

## No-JavaScript Behavior

The no-JavaScript path is the baseline, not an afterthought:

- Every form works with a plain `method="post"` submit and a server redirect.
- Every link is a real `<a href>`.
- Every interactive element has a server-rendered default state.
- Accessibility is a release quality gate: do not claim that Dreego can make
  arbitrary user applications automatically accessible. Verify keyboard
  operability, focus behavior, and screen-reader output for the components
  you build.

## Rules of Thumb

1. Build the no-JS version first. It is the contract.
2. Add HTMX for server round-trips without reloads.
3. Add Alpine.js for local UI state only.
4. Use plain JavaScript for everything else.
5. Keep the server as the single source of truth for state and validation.

## See Also

- [Forms](forms.md) — `g-action`, validation, CSRF
- [Middleware](middleware.md) — CSP and security headers
- [Components](components.md) — component system
- [Roadmap](roadmap.md) — progressive enhancement in the v0.0.x phase
