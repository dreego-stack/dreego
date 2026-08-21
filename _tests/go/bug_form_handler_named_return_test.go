package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugFormHandlerNamedReturn(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<go>
type LoginForm struct {
}
func Save(c *dreego.SSRContext, form LoginForm) (err error) {
    return nil
}
</go>
<div>
  <form g-action="Save" method="post">
    <input name="email">
    <button>Save</button>
  </form>
</div>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "var form LoginForm")
}
