package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestSemanticSectionsCompileWithDefaultLanguages(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<head><title>Semantic sections</title></head>
<server>message := "hello"</server>
<body><main>{{ message }}</main></body>
<style>main { color: green; }</style>
<client>console.log("ready")</client>`)
}

func TestSemanticSectionsCompileWithExplicitDefaultLanguages(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<server lang="go">message := "hello"</server>
<body lang="html"><p>{{ message }}</p></body>
<style lang="css">p { color: green; }</style>
<client lang="js">console.log("ready")</client>`)
}

func TestBodyPreservesNestedHTMLScript(t *testing.T) {
	t.Parallel()
	out := dreegotest.Generate(t, `<body><script>window.ready = true</script></body>`)
	dreegotest.MustContain(t, out, "<script>")
	dreegotest.MustContain(t, out, "window.ready = true")
	dreegotest.MustContain(t, out, "</script>")
}

func TestLegacyRootSectionsExplainMigration(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		src  string
		want string
	}{
		{name: "go", src: `<go>message := "hello"</go><body></body>`, want: "replace root <go> with <server>"},
		{name: "div", src: `<div>legacy body</div>`, want: "replace root <div> with <body>"},
		{name: "script", src: `<script>console.log("legacy")</script><body></body>`, want: "replace root <script> with <client>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dreegotest.MustFailWith(t, tc.src, tc.want)
		})
	}
}

func TestUnknownSectionLanguageExplainsProcessorRequirement(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<body lang="markdown"># Hello</body>`, `unsupported language "markdown" for <body>`)
}

func TestUnknownSectionLanguageDiagnosticHasSourceLocation(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body lang="markdown"># Hello</body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate succeeded, want unsupported-language error: %s", out)
	}
	for _, want := range []string{
		"www/routes/get.dreego:1:1",
		`unsupported language "markdown" for <body>`,
		"install a processor",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic missing %q:\n%s", want, out)
		}
	}
}
