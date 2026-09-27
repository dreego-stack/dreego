package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTranspilerXSSEscaping(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompile(t, `<server>v := "<client>alert(1)</client>"</server>
<body><p>{{ v }}</p></body>`)
}
