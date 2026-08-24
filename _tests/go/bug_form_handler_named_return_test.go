package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugFormHandlerNamedReturn(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<server>
type LoginForm struct {
}
func Save(c *dreego.SSRContext, form LoginForm) (err error) {
    return nil
}
</server>
<body>
  <form g-action="Save" method="post">
    <input name="email">
    <button>Save</button>
  </form>
</body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "var form LoginForm")
}
