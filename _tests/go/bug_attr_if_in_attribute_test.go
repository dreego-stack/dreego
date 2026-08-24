package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugAttrIfInAttribute(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<server>cond := true</server>
<body><a class="nav {#if cond}active{/if}">link</a></body>`, "inside attribute value")
	dreegotest.MustFailWith(t, `<server>cond := true</server>
<body class="nav {#if cond}active{/if}"><p>link</p></body>`, "inside attribute value")
}
