# Component module imports

Dreego resolves component imports during `dreego generate`. The runtime never reads `go.mod`, the module cache, or `.dreego` source files.

Routes can select components from a local website path or a required Go module:

```dreego
COMPONENT "www/components" IMPORT {
    Button,
    Card,
}

COMPONENT "github.com/dreego-stack/dreego-ui/components/dreegoui" IMPORT {
    Navbar,
    PriceCard,
}
```

Each `COMPONENT` directive names a path and lists the component names to import.
An entry may alias a name with `as`:

```dreego
COMPONENT "www/components" IMPORT { Card, Card as ProductCard }
```

The component name is the name of the source `.dreego` file; the alias becomes
the callable `<@ProductCard>` name.

Alias resolution is generator-global by design, consistent with the global
component registry: components found under the website root are callable from
every route, layout, and component, independent of the file that imported them.
An alias declared in one file therefore resolves in every other file. Two files
must not declare the same alias for different targets, and an alias must not
shadow a real component; both collisions fail at `dreego generate`.

The legacy `from "<path>" import { ... }` form is rejected at `dreego generate`
with a `file:line:col` diagnostic pointing to
`COMPONENT "<path>" IMPORT { ... }`.

The generator resolves the module through the application's `go.mod` and Go module resolution, parses the module's `.dreego` sources, and emits ordinary generated Go files in the application's `www/components` package tree. The source files are not copied into the application and are not loaded at runtime.

The Go package of a plugin remains independently importable from `main.go` for runtime registration. This means a plugin can provide both Go runtime behavior and generator-visible `.dreego` component sources without adding a client runtime or a source-file dependency to the final binary.
