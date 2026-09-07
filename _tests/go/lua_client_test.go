package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestLuaRootClientProducesJavaScriptAndRuntime(t *testing.T) {
	out := dreegotest.Generate(t, `<body><main id="status">Ready</main></body>
<client lang="lua">
local status = document.querySelector("#status")
status.textContent = "Loaded"
print(status.textContent)
</client>`)
	for _, want := range []string{
		`<script src="/_dreego/lua.js"></script>`,
		`let status = document.querySelector("#status");`,
		`status.textContent = "Loaded";`,
		"dreegoLua.print(status.textContent);",
	} {
		dreegotest.MustContain(t, out, want)
	}
}

func TestLuaGenerateBuildsOneFeatureLinkedRuntimeAsset(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego":   `<body></body><client lang="lua">print("ready")</client>`,
		"www/routes/about.dreego": `<body></body><client lang="lua">local title = "About"</client>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	root, err := os.ReadFile(filepath.Join(dir, "www", "dree.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := string(root)
	if strings.Count(generated, `/_dreego/lua.js`) != 1 || !strings.Contains(generated, "console.log") {
		t.Fatalf("generated runtime asset:\n%s", generated)
	}
	if strings.Contains(generated, "concatValue") || strings.Contains(generated, "Lua arithmetic") {
		t.Fatalf("unused Lua helpers were linked:\n%s", generated)
	}
}

func TestLuaInlineBodyScriptProducesJavaScript(t *testing.T) {
	out := dreegotest.Generate(t, `<body>
<script lang="lua">local enabled = true
if enabled then print("ready") end</script>
</body>`)
	dreegotest.MustContain(t, out, "let enabled = true;")
	dreegotest.MustContain(t, out, "if (dreegoLua.truthy(enabled))")
}

func TestLuaFunctionsCompileBrowserCallbacks(t *testing.T) {
	out := dreegotest.Generate(t, `<body><button id="count">Count</button></body>
<client lang="lua">
local button = document:querySelector("#count")
local count = 0
local function render()
    button.textContent = "Count " .. count
end
button:addEventListener("click", function()
    count = count + 1
    render()
end)
</client>`)
	for _, want := range []string{
		`document.querySelector("#count")`,
		"let render = () => {",
		`button.addEventListener("click", () => {`,
		"count = dreegoLua.add(count, 1);",
		`button.textContent = dreegoLua.concat("Count ", count);`,
	} {
		dreegotest.MustContain(t, out, want)
	}
}

func TestLuaFunctionReturnCompilesThroughGenerate(t *testing.T) {
	out := dreegotest.Generate(t, `<body><output id="result"></output></body>
<client lang="lua">
local function describe(value)
    if value then
        return "truthy"
    end
    return "falsey"
end
local result = document:querySelector("#result")
result.textContent = describe(0)
</client>`)
	dreegotest.MustContain(t, out, `return "truthy";`)
	dreegotest.MustContain(t, out, `return "falsey";`)
	dreegotest.MustContain(t, out, "result.textContent = describe(0);")
}

func TestLuaApplicationBuilds(t *testing.T) {
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<body><main>Lua</main></body>
<client lang="lua">local ready = true
if ready then print("ready") end</client>`,
	})
}

func TestLuaUnsupportedFeatureFailsGeneration(t *testing.T) {
	_, err := dreegotest.MustGenerate(t, `<body></body>
<client lang="lua">local values = {1, 2}</client>`)
	if err == nil || !strings.Contains(err.Error(), "tables are not supported") {
		t.Fatalf("error = %v", err)
	}
}

func TestLuaDiagnosticUsesDreegoSourceLine(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body><main>Ready</main></body>
<client lang="lua">
local values = {1, 2}
</client>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate succeeded:\n%s", out)
	}
	if !regexp.MustCompile(`www/routes/get\.dreego\(3,\d+\): Lua: tables are not supported`).MatchString(out) {
		t.Fatalf("diagnostic does not identify the Lua source line:\n%s", out)
	}
}
