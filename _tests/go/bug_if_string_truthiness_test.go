package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// Feedback 3.5: the _docs/forms.md example {#if c.Errors("email")} must
// compile. c.Errors returns a string; the generated condition must route it
// through dreego.Truthy instead of emitting a non-boolean if condition.
func TestBugIfStringConditionCompiles(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/login/+page.dreego": `<server>
    type LoginForm struct {
        Email string ` + "`form:\"email\" validate:\"required,email\"`" + `
    }

    func Login(c dreego.Context, form LoginForm) error {
        return c.Redirect("/login", 303)
    }
</server>

<body>
    <h1>Login</h1>
    {#if c.Errors("email")}<p class="error">{{ c.Errors("email") }}</p>{/if}
    <form g-action="Login" method="post">
        <input name="email" type="email" value="{{ c.Old("email") }}">
        <button type="submit">Login</button>
    </form>
</body>`,
	})
}

// Truthiness matrix at the build level: strings (empty = false), ints
// (zero = false), and slices (empty = false) are valid {#if} conditions.
func TestBugIfTruthinessMatrixCompiles(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<server>
    title := ""
    count := 0
    items := []string{}
    score := 7
</server>
<body>
{#if title}<p>title</p>{/if}
{#if count}<p>count</p>{/if}
{#if items}<p>items</p>{/if}
{#if score}<p>score</p>{/if}
{#if score >= 5}<p>high</p>{/if}
{#if title}<p>a</p>{#else}<p>b</p>{/if}
</body>`,
	})
}

// The runtime semantics must match the compile-time contract: empty string and
// zero render the else branch; non-zero renders the main branch.
func TestBugIfStringTruthinessRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<server>
    title := ""
    count := 0
    items := []string{}
</server>
<body>
{#if title}<p id="title">yes</p>{#else}<p id="title">no</p>{/if}
{#if count}<p id="count">yes</p>{#else}<p id="count">no</p>{/if}
{#if items}<p id="items">yes</p>{#else}<p id="items">no</p>{/if}
</body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `<p id="title">no</p>`)
	dreegotest.MustContainBody(t, body, `<p id="count">no</p>`)
	dreegotest.MustContainBody(t, body, `<p id="items">no</p>`)
}
