package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugElseIfControlFlow(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>score := 85</go>
<div>
{#if score >= 90}
A
{#else if score >= 80}
B
{#else}
C
{/if}
</div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "{#else if")
}
