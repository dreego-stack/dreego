package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestFormActionsGActionBasic(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/get-login.dreego": `<server>
    type LoginForm struct {
        Email string
    }
    func Login(c dreego.Context, form LoginForm) error {
        return nil
    }
</server>
<body>
<form g-action="Login" method="post">
    <input name="email" type="email">
    <button type="submit">Login</button>
</form>
</body>`,
	})
}

func TestFormActionsGActionNoHandler(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-fail.dreego": `<server>
</server>
<body>
<form g-action="Missing" method="post">
    <input name="x">
    <button>OK</button>
</form>
</body>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], `app.Register("POST"`)
}

func TestFormActionsGActionUnexported(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-fail.dreego": `<server>
    type myForm struct {
        X string
    }
    func myAction(c dreego.Context, form myForm) error {
        return nil
    }
</server>
<body>
<form g-action="myAction" method="post">
    <input name="x">
    <button>OK</button>
</form>
</body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "HandleIndexPost")
}

func TestFormActionsGActionWrongArity(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-fail.dreego": `<server>
    type BadForm struct {
        X string
    }
    func Bad(c dreego.Context) error {
        return nil
    }
</server>
<body>
<form g-action="Bad" method="post">
    <input name="x">
    <button>OK</button>
</form>
</body>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], `app.Register("POST"`)
}

func TestFormActionsHandlerSignature(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get-fail.dreego": `<server>
    type BadForm struct {
        X string
    }
    func bad(c dreego.Context, form BadForm) string {
        return "wrong"
    }
</server>
<body>
<form g-action="bad" method="post">
    <input name="x">
    <button>OK</button>
</form>
</body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if dreegotest.BuildInDirOK(t, dir) {
		t.Fatal("should not build with wrong return type")
	}
}

func TestFormActionsMethodPostFile(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/post-login.dreego": `<server>
    type LoginForm struct {
        Email string
    }
    func Login(c dreego.Context, form LoginForm) error {
        return nil
    }
</server>
<body>
<form g-action="Login" method="post">
    <input name="email">
    <button>Login</button>
</form>
</body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "HandleIndexGet")
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "HandleIndexPost")
}

func TestFormActionsNoGAction(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/post-search.dreego": `<server>
    c.Set("title", "Search Page")
</server>
<body>
<form method="post">
    <input name="search">
    <button type="submit">Search</button>
</form>
</body>
`,
	})
}

func TestFormActionsNoValidate(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-form.dreego": `<server>
    type NoValForm struct {
        Email string
    }
    func NoVal(c dreego.Context, form NoValForm) error {
        return nil
    }
</server>
<body>
<form g-action="NoVal" method="post">
    <input name="email">
    <button>OK</button>
</form>
</body>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "dreego.ValidateForm")
}

func TestFormActionsPlainForm(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-plain.dreego": `<server>
    email := c.FormValue("email")
    c.Set("email", email)
</server>
<body>
<form method="post">
    <input name="email">
    <button>Submit</button>
</form>
</body>`,
	})
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], `app.Register("POST"`)
}

func TestFormActionsPlainPostRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/routes/post.dreego": `<server>
    email := c.FormValue("email")
    c.Set("email", email)
</server>
<body>
<form method="post">
    <input name="email">
    <button>Submit</button>
</form>
</body>`,
	}, "app.SetCSRF(false); ")
	code, _, _ := c.Request(t, "POST", "/", "email=hello@test.com", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	dreegotest.MustStatus(t, code, 200)
}

func TestFormActionsStructTags(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-form.dreego": "<server>\n    type MyForm struct {\n        Email string `form:\"email\"`\n    }\n    func DoForm(c dreego.Context, form MyForm) error {\n        return nil\n    }\n</server>\n<body>\n<form g-action=\"DoForm\" method=\"post\">\n    <input name=\"email\">\n    <button>OK</button>\n</form>\n</body>",
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "dreego.BindForm")
}

func TestFormActionsValidateTags(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get-form.dreego": "<server>\n    type ValForm struct {\n        Email string `validate:\"required,email\"`\n    }\n    func DoVal(c dreego.Context, form ValForm) error {\n        return nil\n    }\n</server>\n<body>\n<form g-action=\"DoVal\" method=\"post\">\n    <input name=\"email\">\n    <button>OK</button>\n</form>\n</body>",
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "app.ValidateForm")
}

func TestFormActionsBoolBinding(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/routes/get-news.dreego": "<server>\n    type NewsForm struct {\n        Email    string `validate:\"required\"`\n        Subscribe bool\n    }\n    func SubmitNews(c dreego.Context, form NewsForm) error {\n        if form.Subscribe {\n            return c.Redirect(\"/subscribed\", 303)\n        }\n        return c.Redirect(\"/skipped\", 303)\n    }\n</server>\n<body>\n<form g-action=\"SubmitNews\" method=\"post\">\n    <input name=\"email\" type=\"email\">\n    <input name=\"subscribe\" type=\"checkbox\">\n    <button type=\"submit\">Send</button>\n</form>\n</body>",
	}, "app.SetCSRF(false); ")
	code, _, headers := c.Request(t, "POST", "/", "email=a@b.c&subscribe=on", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("status = %d, want 303 for checked", code)
	}
	if !strings.Contains(headers.Get("Location"), "subscribed") {
		t.Fatalf("expected redirect to /subscribed, got %q", headers.Get("Location"))
	}
	code, _, headers = c.Request(t, "POST", "/", "email=a@b.c", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("status = %d, want 303 for unchecked", code)
	}
	if !strings.Contains(headers.Get("Location"), "skipped") {
		t.Fatalf("expected redirect to /skipped, got %q", headers.Get("Location"))
	}
}

func TestFormActionsIntBinding(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/routes/get-age.dreego": "<server>\n    type AgeForm struct {\n        Age int `validate:\"min=2\"`\n    }\n    func SubmitAge(c dreego.Context, form AgeForm) error {\n        if form.Age == 20 {\n            return c.Redirect(\"/adult\", 303)\n        }\n        return c.Redirect(\"/other\", 303)\n    }\n</server>\n<body>\n<form g-action=\"SubmitAge\" method=\"post\">\n    <input name=\"age\" type=\"number\">\n    <button type=\"submit\">Send</button>\n</form>\n</body>",
	}, "app.SetCSRF(false); ")
	code, _, headers := c.Request(t, "POST", "/", "age=20", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("expected 303 for valid age, got %d", code)
	}
	if !strings.Contains(headers.Get("Location"), "adult") {
		t.Fatalf("expected redirect to /adult, got %q", headers.Get("Location"))
	}
	code, _, _ = c.Request(t, "POST", "/", "age=1", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 200 {
		t.Fatalf("expected 200 re-render when min fails, got %d", code)
	}
}

func TestFormActionsSubmitValid(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/routes/post-login.dreego": "<server>\n    type LoginForm struct {\n        Email string `validate:\"required,email\"`\n    }\n    func Login(c dreego.Context, form LoginForm) error {\n        return c.Redirect(\"/dashboard\", 303)\n    }\n</server>\n<body>\n<form g-action=\"Login\" method=\"post\">\n    <input name=\"email\" type=\"email\">\n    <button type=\"submit\">Login</button>\n</form>\n</body>",
	}, "app.SetCSRF(false); ")
	code, _, _ := c.Request(t, "POST", "/", "email=test@dreego.dev", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("expected 303 redirect, got %d", code)
	}
}

func TestFormActionsSubmitInvalid(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/routes/post-login.dreego": "<server>\n    type LoginForm struct {\n        Email string `validate:\"required,email\"`\n    }\n    func Login(c dreego.Context, form LoginForm) error {\n        return c.Redirect(\"/ok\", 303)\n    }\n</server>\n<body>\n<form g-action=\"Login\" method=\"post\">\n    <input name=\"email\" type=\"email\">\n    <button type=\"submit\">Login</button>\n</form>\n</body>",
	}, "app.SetCSRF(false); ")
	code, _, _ := c.Request(t, "POST", "/", "email=invalid", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	dreegotest.MustStatus(t, code, 200)
}

func TestFormActionsSubmitCSRFPass(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/routes/post-login.dreego": "<server>\n    type LoginForm struct {\n        Email string `validate:\"required\"`\n    }\n    func Login(c dreego.Context, form LoginForm) error {\n        return c.Redirect(\"/ok\", 303)\n    }\n</server>\n<body>\n<form g-action=\"Login\" method=\"post\">\n    <input name=\"email\" type=\"email\">\n    <button type=\"submit\">Login</button>\n</form>\n</body>",
	}, "app.SetSessionStore(dreego.NewCookieStore([]byte(\"01234567890123456789012345678903\"))); ")
	c.Get(t, "/health")
	token := c.Cookie("csrf_token")
	if token == "" {
		t.Fatalf("no csrf_token cookie issued")
	}
	code, _, _ := c.Request(t, "POST", "/", "email=test@dreego.dev&csrf_token="+token, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("expected 303 with valid CSRF, got %d", code)
	}
}
