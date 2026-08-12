package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerDuplicateHead(t *testing.T) {
	t.Parallel()
	dreegotest.MustFail(t, `<head><title>A</title></head>
<head><title>B</title></head>
<div><p>hi</p></div>`)
}
