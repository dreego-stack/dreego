package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugComponentIfEach(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/List.dreego": `Component List (items []string)
<body><ul>
{#each items as item}
    <li>{{ item }}</li>
{/each}
</ul></body>`,
		"www/routes/get.dreego": `<server>items := []string{"a", "b"}</server>
<body><@List items={items}/></body>`,
	})
	dreegotest.MustNotContain(t, gen["www/components/dree.go"], "{#each")
	dreegotest.MustNotContain(t, gen["www/components/dree.go"], "{#if")
}
