package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugNestedIfInElse(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": `<server>score := 85</server>
<body>
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
	})
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], "if dreego.Truthy(score >= 80)")
}
