package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerDuplicateDiv(t *testing.T) {
	t.Parallel()
	dreegotest.MustFail(t, `<div><p>first</p></div>
<div><p>second</p></div>`)
}
