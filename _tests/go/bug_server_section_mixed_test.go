package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// A <server> section that mixes declarations and statements must compile:
// declarations go to package level, statements stay inside the render function.
func TestBugServerSectionMixedDeclarationsAndStatements(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `<server>
type T struct{ V int }
func (t T) Get() int { return t.V }
x := T{V: 5}
</server>
<body>{{ x.Get() }}</body>`,
	})
	out := gen["www/routes/dree.go"]
	dreegotest.MustContain(t, out, "type T struct")
	dreegotest.MustContain(t, out, "func (t T) Get() int")
}

// func declarations plus statements in one section (feedback 3.2, second case).
func TestBugServerSectionFuncPlusStatements(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
func helperA() int { return 41 }
n := helperA() + 1
</server>
<body>{{ n }}</body>`,
	})
}

// Mixing a form struct, its handler, and a derived local value (feedback 3.3).
func TestBugServerSectionFormActionMixed(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
type F struct { A string ` + "`form:\"a\" validate:\"required\"`" + ` }
func Act(c dreego.Context, f F) error { return c.Redirect("/", 303) }
derived := "computed"
</server>
<body>
<form g-action="Act" method="post"><input name="a"><button>Go</button></form>
<p>{{ derived }}</p>
</body>`,
	})
}

// Declarations-only sections must also land at package level, so route files in
// the same directory can share the type.
func TestBugServerSectionDeclarationsSharedAcrossRoutes(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `<server>
type Item struct{ Name string }
</server>
<body>{{ (Item{Name: "a"}).Name }}</body>`,
		"www/routes/other.dreego": `<server>
item := Item{Name: "b"}
</server>
<body>{{ item.Name }}</body>`,
	})
	out := gen["www/routes/dree.go"]
	if strings.Count(out, "type Item struct") != 1 {
		t.Fatalf("shared declaration must be emitted exactly once, got:\n%s", out)
	}
}

// Two route files that hoist the same package-level name must fail generation
// with a dreego diagnostic, not a raw compiler redeclaration error.
func TestBugServerSectionDuplicateDeclarationDiagnostic(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `<server>
type Product struct{ Name string }
</server>
<body>{{ (Product{Name: "a"}).Name }}</body>`,
		"www/routes/other.dreego": `<server>
type Product struct{ Name string }
</server>
<body>{{ (Product{Name: "b"}).Name }}</body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate must reject duplicate package-level declarations, got:\n%s", out)
	}
	for _, want := range []string{"duplicate package-level declaration", "Product", "Fix:"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

// A pure statements section must still compile inside the render function.
func TestBugServerSectionStatementsOnlyStillCompiles(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
msg := "hello"
</server>
<body>{{ msg }}</body>`,
	})
}
