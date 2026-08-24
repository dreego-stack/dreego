package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestComponentCallNamedPropsOrder(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"www/components/Greet.dreego": `Component Greet (first string, second string)
<body><p>{{ first }} {{ second }}</p></body>`,
		"www/routes/get.dreego": `<body><@Greet second="World" first="Hello"/></body>`,
	}
	gen := dreegotest.Build(t, files)
	if !strings.Contains(gen["www/components/dree.go"], "Greet(") {
		t.Fatal("generated component function not found")
	}
}

func TestComponentCallMissingProp(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (first string, second string)
<body><p>{{ first }} {{ second }}</p></body>`,
		"www/routes/get.dreego": `<body><@Greet first="Hello"/></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for missing prop, got success: %s", out)
	}
	if !strings.Contains(out, "Greet") {
		t.Fatalf("error must name component Greet, got: %s", out)
	}
	if !strings.Contains(out, "second") {
		t.Fatalf("error must name missing prop second, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") {
		t.Fatalf("error must reference the calling source path, got: %s", out)
	}
}

func TestComponentCallUnknownProp(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (first string)
<body><p>{{ first }}</p></body>`,
		"www/routes/get.dreego": `<body><@Greet first="Hello" second="World"/></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for unknown prop, got success: %s", out)
	}
	if !strings.Contains(out, "Greet") {
		t.Fatalf("error must name component Greet, got: %s", out)
	}
	if !strings.Contains(out, "second") {
		t.Fatalf("error must name unknown prop second, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") {
		t.Fatalf("error must reference the calling source path, got: %s", out)
	}
}

func TestComponentCallDuplicateProp(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (first string)
<body><p>{{ first }}</p></body>`,
		"www/routes/get.dreego": `<body><@Greet first="Hello" first="Again"/></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for duplicate prop, got success: %s", out)
	}
	if !strings.Contains(out, "Greet") {
		t.Fatalf("error must name component Greet, got: %s", out)
	}
	if !strings.Contains(out, "first") {
		t.Fatalf("error must name duplicated prop first, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") {
		t.Fatalf("error must reference the calling source path, got: %s", out)
	}
}

func TestComponentCallNamedPropsHTTP(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Greet.dreego": `Component Greet (first string, second string)
<body><p>{{ first }} {{ second }}</p></body>`,
		"www/routes/get.dreego": `<body><@Greet second="World" first="Hello"/></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello") {
		t.Fatalf("first prop not rendered, got: %s", body)
	}
	if !strings.Contains(body, "World") {
		t.Fatalf("second prop not rendered, got: %s", body)
	}
	if !strings.Contains(body, "<p>Hello World</p>") {
		t.Fatalf("expected ordered props rendered as <p>Hello World</p>, got: %s", body)
	}
}
