package wails

import (
	"context"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type bindingSnapshot struct {
	Ready bool   `json:"ready"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type validBinding struct{}

func (*validBinding) Snapshot(context.Context) (bindingSnapshot, error) {
	return bindingSnapshot{}, nil
}

func (*validBinding) Reset() error {
	return nil
}

type channelBinding struct{}

func (*channelBinding) Stream() chan string {
	return nil
}

type functionBinding struct{}

func (*functionBinding) Apply(func()) {}

type interfaceBinding struct{}

func (*interfaceBinding) Value() any {
	return nil
}

type multipleValueBinding struct{}

func (*multipleValueBinding) Values() (string, bool) {
	return "", false
}

type recursiveBindingValue struct {
	Next *recursiveBindingValue `json:"next"`
}

type recursiveBinding struct{}

func (*recursiveBinding) Value() recursiveBindingValue {
	return recursiveBindingValue{}
}

func TestValidateServicesAcceptsTypedContract(t *testing.T) {
	service := application.NewService(&validBinding{})
	if err := validateServices([]application.Service{service}); err != nil {
		t.Fatalf("validateServices: %v", err)
	}
}

func TestValidateServicesRejectsUnsupportedValues(t *testing.T) {
	for _, test := range []struct {
		name    string
		service application.Service
		want    string
	}{
		{name: "channel", service: application.NewService(&channelBinding{}), want: "chan string"},
		{name: "function", service: application.NewService(&functionBinding{}), want: "func()"},
		{name: "interface", service: application.NewService(&interfaceBinding{}), want: "interface {}"},
		{name: "multiple values", service: application.NewService(&multipleValueBinding{}), want: "want error"},
		{name: "nil pointer", service: application.NewService((*validBinding)(nil)), want: "non-nil pointer"},
		{name: "recursive", service: application.NewService(&recursiveBinding{}), want: "recursive type"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validateServices([]application.Service{test.service})
			if err == nil || !strings.Contains(err.Error(), test.want) || !strings.Contains(err.Error(), "Fix:") {
				t.Fatalf("validateServices error = %v", err)
			}
		})
	}
}

func TestApplicationOptionsRegistersOnlyExplicitServices(t *testing.T) {
	service := application.NewService(&validBinding{})
	configured := applicationOptions(Options{Services: []application.Service{service}})
	if len(configured.Services) != 1 || configured.Services[0].Instance() != service.Instance() {
		t.Fatalf("configured services = %#v", configured.Services)
	}
	defaults := applicationOptions(Options{})
	if len(defaults.Services) != 0 {
		t.Fatalf("default services = %#v", defaults.Services)
	}
}
