package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugEachLoopInCond(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/app/routes/+page.dreego": `<server>items := []string{"a", "b", "c"}</server>
<body>
{#each items as item}<div>{#if !$loop.Last}, {/if}{{ item }}</div>{/each}
</body>`,
	})
	dreegotest.MustContain(t, gen["www/app/routes/dree.go"], "loop.Last")
	dreegotest.MustNotContain(t, gen["www/app/routes/dree.go"], "$loop")
}
