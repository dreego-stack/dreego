package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentIfEach(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/List.dreego": `Component List (items []string)
<div><ul>
{#each items as item}
    <li>{{ item }}</li>
{/each}
</ul></div>`,
		"www/routes/get.dreego": `<go>items := []string{"a", "b"}</go>
<div><@List items={items}/></div>`,
	})
	dreegotest.MustNotContain(t, gen["www/components/dree.go"], "{#each")
	dreegotest.MustNotContain(t, gen["www/components/dree.go"], "{#if")
}
