package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerXSSEscaping(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<go>v := "<script>alert(1)</script>"</go>
<div><p>{{ v }}</p></div>`)
}
