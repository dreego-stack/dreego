package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerUnclosedDiv(t *testing.T) {
	t.Parallel()
	dreegotest.MustFail(t, `<div>no end`)
}
