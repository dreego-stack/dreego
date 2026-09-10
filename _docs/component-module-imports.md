# Component module imports

Dreego resolves component imports during `dreego generate`. The runtime never reads `go.mod`, the module cache, or `.dreego` source files.

Routes can select components from a local website path or a required Go module:

```dreego
from "www/components" import {
    Button,
    Card,
}

from "github.com/dreego-stack/dreego-ui/components/dreegoui" import {
    Navbar,
    PriceCard,
}
```

The generator resolves the module through the application's `go.mod` and Go module resolution, parses the module's `.dreego` sources, and emits ordinary generated Go files in the application's `www/components` package tree. The source files are not copied into the application and are not loaded at runtime.

The Go package of a plugin remains independently importable from `main.go` for runtime registration. This means a plugin can provide both Go runtime behavior and generator-visible `.dreego` component sources without adding a client runtime or a source-file dependency to the final binary.
