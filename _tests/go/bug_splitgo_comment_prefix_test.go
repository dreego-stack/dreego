package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugSplitGoCommentPrefix(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
// UserForm holds the login data
type UserForm struct {
}
func Save(c *dreego.SSRContext, form UserForm) error {
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
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "type UserForm struct")
}
