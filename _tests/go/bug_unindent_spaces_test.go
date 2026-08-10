package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugUnindentSpaces(t *testing.T) {
	gen := dreegotest.Build(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
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
	dreegotest.MustNotContain(t, gen["dreego/gen/routes.go"], "    type UserForm struct")
	dreegotest.MustContain(t, gen["dreego/gen/routes.go"], "type UserForm struct")
}
