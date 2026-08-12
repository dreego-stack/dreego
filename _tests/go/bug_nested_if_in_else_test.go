package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugNestedIfInElse(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>score := 85</go>
<div>
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
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "if score >= 80")
}
