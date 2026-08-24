package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugNestedIfInElse(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<server>score := 85</server>
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
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "if score >= 80")
}
