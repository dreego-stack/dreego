package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSSRContextJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(rec, r)

	c.JSON(http.StatusCreated, map[string]any{"name": "dreego", "n": 1})

	if rec.Code != http.StatusCreated {
		t.Errorf("status: expected %d, got %d", http.StatusCreated, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("content-type: expected application/json; charset=utf-8, got %q", ct)
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("body: invalid json %q: %v", rec.Body.String(), err)
	}
	want := map[string]any{"n": float64(1), "name": "dreego"}
	for k, v := range want {
		if gv, ok := got[k]; !ok || gv != v {
			t.Errorf("body: field %q expected %v, got %v (full %v)", k, v, got[k], got)
		}
	}
}

func TestSSRContextJSONEncodeError(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(rec, r)

	c.JSON(http.StatusOK, struct{ Ch chan int }{Ch: make(chan int)})

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	body := rec.Body.String()
	if body == "" {
		t.Fatal("expected an error body, got empty")
	}
	if !strings.Contains(body, http.StatusText(http.StatusInternalServerError)) {
		t.Errorf("expected generic error body, got %q", body)
	}
	if strings.Contains(body, "unsupported type") {
		t.Errorf("internal cause must not be disclosed, got %q", body)
	}
}

func TestSSRContextXML(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(rec, r)

	c.XML(http.StatusOK, struct {
		XMLName struct{} `xml:"greeting"`
		Text    string   `xml:"text"`
	}{Text: "hello"})

	if rec.Code != http.StatusOK {
		t.Errorf("status: expected %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/xml; charset=utf-8" {
		t.Errorf("content-type: expected application/xml; charset=utf-8, got %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "hello") {
		t.Errorf("expected 'hello' in xml body, got %q", rec.Body.String())
	}
}

func TestSSRContextXMLEncodeError(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(rec, r)

	c.XML(http.StatusOK, struct{ Ch chan int }{Ch: make(chan int)})

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if rec.Body.String() == "" {
		t.Fatal("expected an error body, got empty")
	}
}

func TestSSRContextBind(t *testing.T) {
	body := `{"name":"dreego","n":2}`
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	c := NewSSR(httptest.NewRecorder(), r)

	var target struct {
		Name string `json:"name"`
		N    int    `json:"n"`
	}
	if err := c.Bind(&target); err != nil {
		t.Fatalf("Bind: unexpected error %v", err)
	}
	if target.Name != "dreego" || target.N != 2 {
		t.Errorf("Bind: got %+v", target)
	}
}

func TestSSRContextBindError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"name":`))
	c := NewSSR(httptest.NewRecorder(), r)

	var target struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&target); err == nil {
		t.Fatal("expected an error decoding malformed JSON, got nil")
	}
}

func TestSSRContextBindBodyLimit(t *testing.T) {
	big := strings.Repeat("x", maxBindBodySize+1)
	body := `{"name":"` + big + `"}`
	r := httptest.NewRequest("POST", "/", strings.NewReader(body))
	c := NewSSR(httptest.NewRecorder(), r)

	var target struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&target); err == nil {
		t.Fatal("expected error when binding oversized body, got nil")
	}
}

func TestSSRContextWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(rec, r)

	body := []byte("<h1>hi</h1>")
	c.Write(http.StatusAccepted, "text/html", body)

	if rec.Code != http.StatusAccepted {
		t.Errorf("status: expected %d, got %d", http.StatusAccepted, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("content-type: expected text/html; charset=utf-8, got %q", ct)
	}
	if got := rec.Body.String(); got != "<h1>hi</h1>" {
		t.Errorf("body: expected %q, got %q", "<h1>hi</h1>", got)
	}
}

func TestSSRContextWants(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "text/html, application/json")
	c := NewSSR(httptest.NewRecorder(), r)

	if !c.Wants("text/html") {
		t.Error("expected Wants(text/html) true")
	}
	if !c.Wants("application/json") {
		t.Error("expected Wants(application/json) true")
	}
	if c.Wants("application/xml") {
		t.Error("expected Wants(application/xml) false")
	}
}

func TestWantsEmptyAccept(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(httptest.NewRecorder(), r)
	if c.Wants("text/html") {
		t.Error("expected Wants false on empty Accept header")
	}
}

func TestWantsWithParams(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "text/html;q=0.9, application/json")
	c := NewSSR(httptest.NewRecorder(), r)

	if !c.Wants("text/html") {
		t.Error("expected Wants(text/html) true with q param")
	}
}

func TestStringsContainsMime(t *testing.T) {
	cases := []struct {
		name   string
		accept string
		mime   string
		want   bool
	}{
		{"exact", "text/html", "text/html", true},
		{"with params", "text/html;q=0.9", "text/html", true},
		{"whitespace", " text/html , application/json ", "text/html", true},
		{"no match", "application/json", "text/html", false},
		{"empty accept", "", "text/html", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := stringsContainsMime(tc.accept, tc.mime); got != tc.want {
				t.Errorf("stringsContainsMime(%q, %q): expected %v, got %v", tc.accept, tc.mime, tc.want, got)
			}
		})
	}
}
