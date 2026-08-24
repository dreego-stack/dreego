package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerUnclosedBody(t *testing.T) {
	t.Parallel()
	dreegotest.MustFail(t, `<body>no end`)
}
