package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTemplateComponentNestedIfElse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Grade.dreego": `Component Grade (score int)
<body class="grade">
{#if score >= 90}
A
{#else}
{#if score >= 80}
B
{#else}
C
{/if}
D
{/if}
</body>`,
		"www/routes/get.dreego": `<server>score := 85</server>
<body>
<@Grade score={score}/>
</body>`,
	})
}

func TestTemplateEachElse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>items := []string{}</server>
<body>{#each items as item}<p>{{ item }}</p>{#each else}<p>empty</p>{/each}</body>`,
	})
}

func TestTemplateEachEmpty(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>items:=[]string{}</server>
<body>{#each items as item}<span>{{ item }}</span>{/each}<p>done</p></body>`,
	})
}

func TestTemplateEachLoopVar(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>items := []string{"a", "b", "c"}</server>
<body>{#each items as item}<p>{{ $loop.Index }}: {{ item }}</p>{/each}</body>`,
	})
}

func TestTemplateEachLoop(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>items := []string{"a", "b"}</server>
<body><ul>{#each items as item}<li>{{ item }}</li>{/each}</ul></body>`,
	})
}

func TestTemplateEachWithIf(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>items:=[]string{"a","","c"}</server>
<body>{#each items as item}{#if item != ""}<span>{{ item }}</span>{/if}{/each}</body>`,
	})
}

func TestTemplateElseOutsideIf(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body>{#else}</body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err == nil {
		t.Fatalf("expected generate failure but succeeded: %s", out)
	}
}

func TestTemplateExpression(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<head><title>T</title></head>
<server>x := "world"</server>
<body><h1>Hello {{ x }}</h1></body>`,
	})
}

func TestTemplateFilters(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>rawHtml := "<b>bold</b>"</server>
<body><p>{{ rawHtml|raw }}</p><p>{{ rawHtml }}</p></body>`,
	})
}

func TestTemplateIfElse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>show := false</server>
<body>{#if show}<p>yes</p>{#else}<p>no</p>{/if}</body>`,
	})
}

func TestTemplateIfFalse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>x := false</server>
<body>{#if x}<strong>yes</strong>{/if}<p>no</p></body>`,
	})
}

func TestTemplateIfTrue(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>x := true</server>
<body>{#if x}<strong>yes</strong>{/if}</body>`,
	})
}

func TestTemplateMissingVar(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body><p>{{ undefined }}</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if dreegotest.BuildInDirOK(t, dir) {
		t.Fatal("expected build failure but succeeded")
	}
}

func TestTemplateNestedIf(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<server>
    x := true
    y := true
</server>
<body>{#if x}{#if y}<strong>both</strong>{/if}{/if}</body>`,
	})
}

func TestTemplateVerbatim(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get.dreego": `<body><p>before</p>{#verbatim}<script>var x = {a: 1};</script>{/verbatim}<p>after</p></body>`,
	})
}
