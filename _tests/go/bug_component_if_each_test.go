package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentIfEach(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/components/List.dreego": `Component List (items []string)
<div><ul>
{#each items as item}
    <li>{{ item }}</li>
{/each}
</ul></div>`,
		"dreego/routes/get.dreego": `<go>items := []string{"a", "b"}</go>
<div><@List items={items}/></div>`,
	})
	dreegotest.MustNotContain(t, gen["dreego/gen/components.go"], "{#each")
	dreegotest.MustNotContain(t, gen["dreego/gen/components.go"], "{#if")
}
