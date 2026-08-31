package render

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/context"
)

func TestComponentRendersHTML(t *testing.T) {
	res, err := Component(func(c *context.SSRContext) (string, error) {
		return "<p>hi</p>", nil
	})
	if err != nil {
		t.Fatalf("Component returned error: %v", err)
	}
	if !bytes.Equal(res.HTML, []byte("<p>hi</p>")) {
		t.Errorf("HTML = %q, want %q", res.HTML, "<p>hi</p>")
	}
	if res.Head != nil {
		t.Errorf("Head = %q, want nil", res.Head)
	}
	if res.Assets != nil {
		t.Errorf("Assets = %v, want nil", res.Assets)
	}
}

func TestComponentPropagatesError(t *testing.T) {
	res, err := Component(func(c *context.SSRContext) (string, error) {
		return "", errors.New("boom")
	})
	if err == nil {
		t.Fatal("expected component error to propagate")
	}
	if res.HTML != nil || res.Head != nil || res.Assets != nil {
		t.Errorf("expected zero Result on error, got %+v", res)
	}
}
