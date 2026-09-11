package context

import (
	stdctx "context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dreego-stack/dreego/internal/session"
)

func TestRenderContextStateContract(t *testing.T) {
	t.Parallel()
	ctx := NewRender(nil)
	if ctx.Err() != nil || ctx.Data("missing") != nil || ctx.Get("missing") != "" {
		t.Fatal("new render context does not expose empty background state")
	}
	ctx.Set("count", 3)
	ctx.Set("label", "ready")
	ctx.Set("error_email", "invalid")
	ctx.Set("old_email", "ada@example.com")
	if ctx.Data("count") != 3 || ctx.Get("count") != "" || ctx.Get("label") != "ready" {
		t.Fatal("render context did not preserve typed and string values")
	}
	if ctx.Errors("email") != "invalid" || ctx.Old("email") != "ada@example.com" {
		t.Fatal("render context validation helpers returned wrong values")
	}
	ctx.Delete("label")
	if ctx.Data("label") != nil {
		t.Fatal("render context delete kept the value")
	}
}

func TestRenderContextPreservesCancellation(t *testing.T) {
	t.Parallel()
	parent, cancel := stdctx.WithCancel(stdctx.Background())
	ctx := NewRender(parent)
	cancel()
	if !errors.Is(ctx.Err(), stdctx.Canceled) {
		t.Fatalf("Err() = %v, want context.Canceled", ctx.Err())
	}
}

func TestSSRContextStateHandlesZeroValueMap(t *testing.T) {
	t.Parallel()
	ctx := &SSRContext{Context: stdctx.Background()}
	if ctx.Data("missing") != nil || ctx.Get("missing") != "" {
		t.Fatal("zero-value data map returned a value")
	}
	ctx.Delete("missing")
	ctx.Set("number", 42)
	ctx.Set("label", "ok")
	if ctx.Data("number") != 42 || ctx.Get("number") != "" || ctx.Get("label") != "ok" {
		t.Fatal("SSR context did not preserve typed and string values")
	}
}

func TestSSRContextValidationStateHelpers(t *testing.T) {
	t.Parallel()
	ctx := NewSSR(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	ctx.Set("error_name", "required")
	ctx.Set("old_name", "Ada")
	if ctx.Errors("name") != "required" || ctx.Old("name") != "Ada" {
		t.Fatal("validation state helpers returned wrong values")
	}
}

func TestSSRContextSessionOperationsWithoutStoreAreNoOps(t *testing.T) {
	t.Parallel()
	ctx := NewSSR(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if ctx.SessionVal("missing") != "" || ctx.CSRFToken() != "" {
		t.Fatal("missing session store returned a value")
	}
	ctx.SetSessionVal("key", "value")
	ctx.DelSessionVal("key")
	ctx.DestroySession()
	if ctx.SessionError() != nil {
		t.Fatalf("SessionError() = %v", ctx.SessionError())
	}
}

func TestSSRContextSessionOperationsUseStore(t *testing.T) {
	t.Parallel()
	store := &recordingStore{values: map[string]string{"user": "Ada"}}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = session.WithStore(req, store)
	ctx := NewSSR(httptest.NewRecorder(), req)
	if ctx.SessionVal("user") != "Ada" {
		t.Fatal("session get did not use the request store")
	}
	ctx.SetSessionVal("theme", "dark")
	ctx.DelSessionVal("user")
	ctx.DestroySession()
	if store.setKey != "theme" || store.setValue != "dark" || store.deleted != "user" || !store.destroyed {
		t.Fatalf("unexpected store calls: %+v", store)
	}
}

func TestSSRContextSessionReadErrorIsObservable(t *testing.T) {
	t.Parallel()
	want := errors.New("store unavailable")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = session.WithStore(req, errorStore{err: want})
	ctx := NewSSR(httptest.NewRecorder(), req)
	if got := ctx.SessionVal("user"); got != "" {
		t.Fatalf("SessionVal() = %q", got)
	}
	if !errors.Is(ctx.SessionError(), want) {
		t.Fatalf("SessionError() = %v", ctx.SessionError())
	}
}

func TestSSRContextSessionWriteErrorsReturnGeneric500(t *testing.T) {
	t.Parallel()
	operations := map[string]func(*SSRContext){
		"set":     func(ctx *SSRContext) { ctx.SetSessionVal("key", "value") },
		"delete":  func(ctx *SSRContext) { ctx.DelSessionVal("key") },
		"destroy": func(ctx *SSRContext) { ctx.DestroySession() },
	}
	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			want := errors.New("write failed")
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = session.WithStore(req, errorStore{err: want})
			ctx := NewSSR(recorder, req)
			operation(ctx)
			if recorder.Code != http.StatusInternalServerError || recorder.Body.String() != "internal server error\n" {
				t.Fatalf("response = %d %q", recorder.Code, recorder.Body.String())
			}
			if !errors.Is(ctx.SessionError(), want) {
				t.Fatalf("SessionError() = %v", ctx.SessionError())
			}
		})
	}
}

type recordingStore struct {
	values           map[string]string
	setKey, setValue string
	deleted          string
	destroyed        bool
}

func (s *recordingStore) Get(_ *http.Request, key string) (string, error) { return s.values[key], nil }
func (s *recordingStore) Set(_ http.ResponseWriter, _ *http.Request, key, value string, _ *session.Options) error {
	s.setKey, s.setValue = key, value
	return nil
}
func (s *recordingStore) Delete(_ http.ResponseWriter, _ *http.Request, key string) error {
	s.deleted = key
	return nil
}
func (s *recordingStore) Destroy(_ http.ResponseWriter, _ *http.Request) error {
	s.destroyed = true
	return nil
}

type errorStore struct{ err error }

func (s errorStore) Get(*http.Request, string) (string, error) { return "", s.err }
func (s errorStore) Set(http.ResponseWriter, *http.Request, string, string, *session.Options) error {
	return s.err
}
func (s errorStore) Delete(http.ResponseWriter, *http.Request, string) error { return s.err }
func (s errorStore) Destroy(http.ResponseWriter, *http.Request) error        { return s.err }
