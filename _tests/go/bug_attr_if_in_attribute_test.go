package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugAttrIfInAttribute(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<go>cond := true</go>
<div><a class="nav {#if cond}active{/if}">link</a></div>`, "inside attribute value")
	dreegotest.MustFailWith(t, `<go>cond := true</go>
<div class="nav {#if cond}active{/if}"><p>link</p></div>`, "inside attribute value")
}
