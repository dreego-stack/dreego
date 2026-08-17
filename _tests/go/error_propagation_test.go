package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestErrorPropagationGeneric500NoDisclosure(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
    panic("database: connection to db.internal:5432 failed")
</go>
<div><p>ok</p></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 500 {
		t.Fatalf("status = %d, want 500", code)
	}
	if strings.Contains(body, "database") || strings.Contains(body, "db.internal") {
		t.Fatalf("internal cause disclosed in body: %q", body)
	}
	if !strings.Contains(body, "Internal Server Error") {
		t.Fatalf("expected generic error body, got: %q", body)
	}
}

func TestErrorPropagationComponentRenderFailure500(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Boom.dreego": `Component Boom ()
<go>
    panic("component render failure")
</go>
<div><p>boom</p></div>`,
		"dreego/routes/get.dreego": `<div><@Boom/></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 500 {
		t.Fatalf("status = %d, want 500", code)
	}
	if strings.Contains(body, "component render failure") {
		t.Fatalf("component cause disclosed in body: %q", body)
	}
}

func TestErrorPropagationFormBindGenericError(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
type Form struct {
    Age int
}
func Save(c dreego.Context, form Form) error {
    return nil
}
</go>
<div>
<form g-action="Save" method="post">
    <input name="age">
    <button>Save</button>
</form>
{#if c.Errors("_form") != ""}<p id="formerr">{{ c.Errors("_form") }}</p>{/if}
</div>`,
	}, "app.SetCSRF(false); ")
	code, body, _ := c.Request(t, "POST", "/", "age=notanumber", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 200 {
		t.Fatalf("status = %d, want 200 re-render", code)
	}
	if !strings.Contains(body, `id="formerr"`) {
		t.Fatalf("expected form error field to be set, got: %q", body)
	}
	if strings.Contains(body, "strconv") {
		t.Fatalf("Go type error disclosed in body: %q", body)
	}
}

func TestErrorPropagationFormActionGenericError(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<go>
type Form struct {
    Name string
}
func Save(c dreego.Context, form Form) error {
    return fmt.Errorf("database: insert into users failed")
}
</go>
<div>
<form g-action="Save" method="post">
    <input name="name">
    <button>Save</button>
</form>
</div>`,
	}, "app.SetCSRF(false); ")
	code, body, _ := c.Request(t, "POST", "/", "name=ada", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 500 {
		t.Fatalf("status = %d, want 500", code)
	}
	if strings.Contains(body, "database") || strings.Contains(body, "insert into users") {
		t.Fatalf("action cause disclosed in body: %q", body)
	}
	if !strings.Contains(body, "Internal Server Error") {
		t.Fatalf("expected generic error body, got: %q", body)
	}
}

func TestErrorPropagationSessionFailureGeneric500(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"failstore.go": `package main

import (
	"errors"
	"net/http"

	dreego "github.com/dreego-stack/dreego/core"
)

type failStore struct{}

func (failStore) Get(*http.Request, string) (string, error) {
	return "", errors.New("store read failure")
}
func (failStore) Set(http.ResponseWriter, *http.Request, string, string, *dreego.Options) error {
	return errors.New("store write failure")
}
func (failStore) Delete(http.ResponseWriter, *http.Request, string) error {
	return errors.New("store delete failure")
}
func (failStore) Destroy(http.ResponseWriter, *http.Request) error {
	return errors.New("store destroy failure")
}`,
		"dreego/routes/get.dreego": `<go>
if c.SessionError() != nil {
    panic(c.SessionError())
}
c.SetSessionVal("k", "v")
if c.SessionError() != nil {
    panic(c.SessionError())
}
v := c.SessionVal("k")
</go>
<div><p>{{ v }}</p></div>`,
	}, "app.SetSessionStore(failStore{}); app.SetCSRF(false); ")
	code, body := c.Get(t, "/")
	if code != 500 {
		t.Fatalf("status = %d, want 500", code)
	}
	if strings.Contains(body, "store") {
		t.Fatalf("session cause disclosed in body: %q", body)
	}
}

func TestErrorPropagationCSRFPersistFailure(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"failstore.go": `package main

import (
	"errors"
	"net/http"

	dreego "github.com/dreego-stack/dreego/core"
)

type failStore struct{}

func (failStore) Get(*http.Request, string) (string, error) {
	return "", errors.New("store read failure")
}
func (failStore) Set(http.ResponseWriter, *http.Request, string, string, *dreego.Options) error {
	return errors.New("store write failure")
}
func (failStore) Delete(http.ResponseWriter, *http.Request, string) error {
	return errors.New("store delete failure")
}
func (failStore) Destroy(http.ResponseWriter, *http.Request) error {
	return errors.New("store destroy failure")
}`,
		"dreego/routes/get.dreego": `<div><p>ok</p></div>`,
	}, "app.SetSessionStore(failStore{}); app.SetLogging(false); ")
	code, body := c.Get(t, "/")
	if code != 500 {
		t.Fatalf("status = %d, want 500", code)
	}
	if strings.Contains(body, "store") {
		t.Fatalf("csrf cause disclosed in body: %q", body)
	}
}
