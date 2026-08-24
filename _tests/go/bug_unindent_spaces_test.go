package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugUnindentSpaces(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<server>
    type UserForm struct {
    }
    func Save(c *dreego.SSRContext, form UserForm) error {
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
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "    type UserForm struct")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "type UserForm struct")
}
