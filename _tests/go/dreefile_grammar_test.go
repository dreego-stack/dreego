package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDreefileGrammarComponentFilenameAndImport(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Greeting.dreego": `DREEFILE component (name string)
<body><p>Hello {{ name }}</p></body>`,
		"www/routes/+page.dreego": `COMPONENT "www/components" IMPORT { Greeting }
<body><@Greeting name="grammar"/></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<p>Hello grammar</p>") {
		t.Fatalf("expected filename-derived component rendered, got: %s", body)
	}
}

func TestDreefileGrammarComponentAlias(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Card.dreego": `DREEFILE component (title string)
<body><article><h2>{{ title }}</h2></article></body>`,
		"www/routes/+page.dreego": `COMPONENT "www/components" IMPORT { Card as ProductCard }
<body><@ProductCard title="Aliased"/></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<h2>Aliased</h2>") {
		t.Fatalf("expected aliased component rendered, got: %s", body)
	}
}

func TestDreefileGrammarLayout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `DREEFILE layout
<body><html><body><nav>Nav</nav>{#slot}<footer>Foot</footer></body></html></body>`,
		"www/routes/+page.dreego": `<body><p>page body</p></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, want := range []string{"<nav>Nav</nav>", "page body", "<footer>Foot</footer>"} {
		if !strings.Contains(body, want) {
			t.Fatalf("layout response missing %q, got: %s", want, body)
		}
	}
}

func TestDreefileGrammarGoImportGenerates(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `GOIMPORT { strings }
<body><p>goimport route</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate with GOIMPORT failed: %v\n%s", err, out)
	}
}

func TestDreefileGrammarComponentAliasShadowsRealComponent(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Card.dreego": `DREEFILE component (title string)
<body><article><h2>{{ title }}</h2></article></body>`,
		"www/components/ProductCard.dreego": `DREEFILE component (title string)
<body><p>{{ title }}</p></body>`,
		"www/routes/+page.dreego": `COMPONENT "www/components" IMPORT { Card as ProductCard }
<body><@ProductCard title="clash"/></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate silently accepted an alias shadowing a real component:\n%s", out)
	}
	if !strings.Contains(out, "ProductCard") {
		t.Fatalf("clash diagnostic must name ProductCard, got:\n%s", out)
	}
}

func TestDreefileGrammarComponentAliasConflictFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Card.dreego": `DREEFILE component (title string)
<body><article><h2>{{ title }}</h2></article></body>`,
		"www/components/Other.dreego": `DREEFILE component (title string)
<body><p>{{ title }}</p></body>`,
		"www/routes/a.dreego": `COMPONENT "www/components" IMPORT { Card as Widget }
<body><@Widget title="a"/></body>`,
		"www/routes/b.dreego": `COMPONENT "www/components" IMPORT { Other as Widget }
<body><@Widget title="b"/></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate silently picked one of two conflicting aliases:\n%s", out)
	}
	if !strings.Contains(out, "Widget") {
		t.Fatalf("conflict diagnostic must name Widget, got:\n%s", out)
	}
}

func TestDreefileGrammarRejectsUnknownDreefileValue(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "DREEFILE Component\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an unknown DREEFILE value:\n%s", out)
	}
	for _, want := range []string{"www/routes/+page.dreego", "1:", "page, component, layout"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarImportPathContainingImport(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Card.dreego": `DREEFILE component (title string)
<body><article><h2>{{ title }}</h2></article></body>`,
		"www/routes/+page.dreego": `COMPONENT "www/IMPORTant/components" IMPORT { Card }
<body><@Card title="import path"/></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<h2>import path</h2>") {
		t.Fatalf("expected component rendered despite IMPORT in path, got: %s", body)
	}
}

func TestDreefileGrammarRejectsInvalidComponentFilename(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/product-card.dreego": "DREEFILE component ()\n<body><p>x</p></body>",
		"www/routes/+page.dreego":            `<body><p>page</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an invalid component filename:\n%s", out)
	}
	for _, want := range []string{"product-card.dreego", "Card.dreego"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarFmtUnbalancedDirectiveKeepsBody(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "DREEFILE component (title string)\n\nCOMPONENT \"www/components\" IMPORT { Card\n\n<body><p>{{ title }}</p></body>\n",
	})
	out, err := dreegotest.RunCLI(t, dir, "fmt", "--stdout", "www/routes/+page.dreego")
	if err != nil {
		t.Fatalf("fmt: %v\n%s", err, out)
	}
	if strings.Count(out, "<body>") != 1 {
		t.Fatalf("fmt moved or duplicated the body, got:\n%s", out)
	}
	bodyIdx := strings.Index(out, "<body>")
	if idx := strings.Index(out, "COMPONENT"); idx < 0 || idx > bodyIdx {
		t.Fatalf("COMPONENT directive must stay in the header, got:\n%s", out)
	}
}

func TestDreefileGrammarFmtCheckStaysLenientOnLegacyHeader(t *testing.T) {
	t.Parallel()
	src := "Component Legacy (name string)\n\nimport dreego github.com/dreego-stack/dreego\n\n<body><p>x</p></body>\n"
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": src,
	})
	formatted, err := dreegotest.RunCLI(t, dir, "fmt", "--stdout", "www/routes/+page.dreego")
	if err != nil {
		t.Fatalf("fmt: %v\n%s", err, formatted)
	}
	if err := os.WriteFile(filepath.Join(dir, "www/routes/+page.dreego"), []byte(formatted), 0644); err != nil {
		t.Fatalf("write formatted file: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "fmt", "--check", "www/routes/+page.dreego"); err != nil {
		t.Fatalf("fmt --check must stay lenient on a legacy header: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err == nil {
		t.Fatalf("generate must keep rejecting the legacy header:\n%s", out)
	}
}

func TestDreefileGrammarFmtBareKeywordMatchesGenerate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "DREEFILE\n<body><p>x</p></body>\n",
	})
	out, err := dreegotest.RunCLI(t, dir, "fmt", "--stdout", "www/routes/+page.dreego")
	if err != nil {
		t.Fatalf("fmt: %v\n%s", err, out)
	}
	bodyIdx := strings.Index(out, "<body>")
	if idx := strings.Index(out, "DREEFILE"); idx < 0 || idx > bodyIdx {
		t.Fatalf("fmt must keep a bare DREEFILE in the header, got:\n%s", out)
	}
}

func TestDreefileGrammarLegacyFormsRejected(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "component line",
			src:  "Component Legacy (name string)\n<body><p>x</p></body>",
			want: "DREEFILE component",
		},
		{
			name: "bare import",
			src:  "import dreego github.com/dreego-stack/dreego\n<body><p>x</p></body>",
			want: "GOIMPORT",
		},
		{
			name: "from import",
			src:  "from \"www/components\" import {\n\tLegacy,\n}\n<body><p>x</p></body>",
			want: "COMPONENT",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := dreegotest.ProjectDir(t, map[string]string{
				"www/routes/+page.dreego": tc.src,
			})
			out, err := dreegotest.RunCLI(t, dir, "generate")
			if err == nil {
				t.Fatalf("generate accepted legacy header form %q:\n%s", tc.src, out)
			}
			for _, want := range []string{"+page.dreego", "1:1:", tc.want} {
				if !strings.Contains(out, want) {
					t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
				}
			}
		})
	}
}
