package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugEachLoopInCond(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>items := []string{"a", "b", "c"}</go>
<div>
{#each items as item}<div>{#if !$loop.Last}, {/if}{item}</div>{/each}
</div>`,
	})
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "loop.Last")
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "$loop")
}
