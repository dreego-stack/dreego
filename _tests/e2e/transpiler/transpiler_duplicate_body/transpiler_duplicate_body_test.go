package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerDuplicateBody(t *testing.T) {
	t.Parallel()
	dreegotest.MustFail(t, `<body><p>first</p></body>
<body><p>second</p></body>`)
}
