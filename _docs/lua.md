# Browser Lua

Dreego compiles a deliberately small Lua language for browser behavior. The
compiler is built into Dreego and has no external dependencies, VM, WebAssembly
runtime, Node installation, or package manager.

Lua may be used in the root client section or in an inline body script:

```html
<client lang="lua">
local status = document.querySelector("#status")

if status then
    status.textContent = "Ready"
    print(status.textContent)
end
</client>
```

```html
<body>
    <button id="save">Save</button>
    <script lang="lua">
    local button = document.querySelector("#save")
    button.disabled = false
    </script>
</body>
```

`dreego generate` translates each block to JavaScript. When Lua behavior differs
from JavaScript, the compiler records a semantic feature. Dreego combines the
features from every route, layout, and component and generates one minimal
`/_dreego/lua.js` asset. Direct-only programs do not generate that asset.

Each generated Lua block executes in its own lexical scope. Reusing a local name
in a route, layout, or component therefore cannot cause a JavaScript global
redeclaration error. Script-end sequences in generated client code are escaped
before HTML emission so a source string cannot terminate its script element.

## MVP syntax

The initial compiler supports:

- `local` declarations and assignments;
- strings using single or double quotes;
- numbers, `true`, `false`, and `nil`;
- `+`, `-`, `*`, `/`, `%`, and `..`;
- `==`, `~=`, `<`, `<=`, `>`, and `>=`;
- Lua `and`, `or`, and `not` semantics;
- `if`, `elseif`, `else`, and `end`;
- local and anonymous functions with lexical closures;
- function calls and dotted browser object access;
- browser method calls using `object:method(arguments)`;
- line comments beginning with `--`;
- `print`, linked to the generated browser runtime.

Zero and empty strings are true. Logical operators return operands and evaluate
their right side lazily. Arithmetic rejects non-number operands rather than
using JavaScript coercion. Modulo by zero raises an error instead of leaking a
JavaScript `NaN` value. Concatenation accepts only strings and numbers.

## Browser boundary

The compiler rejects tables for now rather than implementing an incomplete
table model. Loops, iterators, varargs, multiple returns, metatables, and
coroutines are planned only through tested patch releases. Browser method calls
use Lua's colon spelling but preserve the native JavaScript receiver instead of
injecting an additional Lua `self` argument.

JavaScript `undefined` values received from browser APIs are treated like Lua
`nil` for truthiness. Generated helper calls use the global runtime explicitly,
so a Lua local named `dreegoLua` cannot shadow it. JavaScript reserved words are
mapped to identifiers that cannot collide with valid Lua names.

The following server or dynamic-loading functions are not available:

- `require`, `load`, `loadfile`, and `dofile`;
- filesystem, process, socket, native-module, bytecode, and debug APIs;
- `collectgarbage`.

Dreego does not claim complete Lua 5.x compatibility. Unsupported behavior is
reported during generation where it can be detected; the supported contract is
expanded without silently changing existing semantics.
