package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTemplateComponentNestedIfElse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/components/Grade.dreego": `Component Grade (score int)
<div class="grade">
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
</div>`,
		"dreego/routes/get.dreego": `<go>score := 85</go>
<div>
<@Grade score={score}/>
</div>`,
	})
}

func TestTemplateEachElse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>items := []string{}</go>
<div>{#each items as item}<p>{{ item }}</p>{#each else}<p>empty</p>{/each}</div>`,
	})
}

func TestTemplateEachEmpty(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>items:=[]string{}</go>
<div>{#each items as item}<span>{{ item }}</span>{/each}<p>done</p></div>`,
	})
}

func TestTemplateEachLoopVar(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>items := []string{"a", "b", "c"}</go>
<div>{#each items as item}<p>{{ $loop.Index }}: {{ item }}</p>{/each}</div>`,
	})
}

func TestTemplateEachLoop(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>items := []string{"a", "b"}</go>
<div><ul>{#each items as item}<li>{{ item }}</li>{/each}</ul></div>`,
	})
}

func TestTemplateEachWithIf(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>items:=[]string{"a","","c"}</go>
<div>{#each items as item}{#if item != ""}<span>{{ item }}</span>{/if}{/each}</div>`,
	})
}

func TestTemplateElseOutsideIf(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div>{#else}</div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err == nil {
		t.Fatalf("expected generate failure but succeeded: %s", out)
	}
}

func TestTemplateExpression(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<head><title>T</title></head>
<go>x := "world"</go>
<div><h1>Hello {{ x }}</h1></div>`,
	})
}

func TestTemplateFilters(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>rawHtml := "<b>bold</b>"</go>
<div><p>{{ rawHtml|raw }}</p><p>{{ rawHtml }}</p></div>`,
	})
}

func TestTemplateIfElse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>show := false</go>
<div>{#if show}<p>yes</p>{#else}<p>no</p>{/if}</div>`,
	})
}

func TestTemplateIfFalse(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>x := false</go>
<div>{#if x}<strong>yes</strong>{/if}<p>no</p></div>`,
	})
}

func TestTemplateIfTrue(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<go>x := true</go>
<div>{#if x}<strong>yes</strong>{/if}</div>`,
	})
}

func TestTemplateMissingVar(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>{{ undefined }}</p></div>`,
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
		"dreego/routes/get.dreego": `<go>
    x := true
    y := true
</go>
<div>{#if x}{#if y}<strong>both</strong>{/if}{/if}</div>`,
	})
}

func TestTemplateVerbatim(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>before</p>{#verbatim}<script>var x = {a: 1};</script>{/verbatim}<p>after</p></div>`,
	})
}
