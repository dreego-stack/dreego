# Server Section

`<server>` contains request-time Go. `lang="go"` is the default and only
supported server language.

```html
<server>
name := c.Param("name")
if name == "" {
    name = "World"
}
</server>

<body><h1>Hello {{ name }}</h1></body>
```

Route server code receives `c`, the generated route context. It exposes request
data, route parameters, response helpers, sessions, validation data, and other
documented SSR operations. Values declared in a server section can be used by
the matching head and body templates.

Component server code receives the component's typed props and render context;
it is not a second HTTP handler. See [Components](components.md) for the smaller
component boundary.

## Methods

Sections without `method` belong to GET. A route file can define other methods
explicitly:

```html
<server method="post">
result := saveForm(c)
</server>
<body method="post"><p>{{ result }}</p></body>
```

The server and body sections for one method form one route response. See
[Routing](routing.md) and [Forms](forms.md).

## Typed responses

`type="json"`, `type="xml"`, and custom response helpers allow a route to
respond without an HTML body. Content negotiation and response methods are
documented in [Routing](routing.md) and [Runtime API](runtime.md).
