package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerMismatchedClose(t *testing.T) {
	dreegotest.MustFail(t, `<go>x:=1</go>
<div>text</go>`)
}
