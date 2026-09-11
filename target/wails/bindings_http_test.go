package wails

import (
	"net/http"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type handlerBinding struct{}

func (*handlerBinding) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (*handlerBinding) Ready() bool {
	return true
}

func TestValidateServicesSkipsWailsInternalHTTPMethod(t *testing.T) {
	service := application.NewService(&handlerBinding{})
	if err := validateServices([]application.Service{service}); err != nil {
		t.Fatalf("validateServices: %v", err)
	}
}
