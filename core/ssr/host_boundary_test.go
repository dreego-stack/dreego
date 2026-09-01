package ssr

import (
	"reflect"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestAppDoesNotOwnHTTPServerLifecycle(t *testing.T) {
	typeOfApp := reflect.TypeOf(*dreego.New())
	for _, name := range []string{"server", "serverConfig", "shutdownDone"} {
		if _, ok := typeOfApp.FieldByName(name); ok {
			t.Fatalf("core App still owns SSR lifecycle field %q", name)
		}
	}
	for _, name := range []string{"Listen", "Shutdown", "SetServerConfig"} {
		if _, ok := reflect.TypeOf(dreego.New()).MethodByName(name); ok {
			t.Fatalf("core App still exposes SSR lifecycle method %q", name)
		}
	}
}
