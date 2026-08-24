package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerMismatchedClose(t *testing.T) {
	t.Parallel()
	dreegotest.MustFail(t, `<server>x:=1</server>
<body>text</server>`)
}
