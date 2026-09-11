package wails

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var contextType = reflect.TypeFor[context.Context]()
var errorType = reflect.TypeFor[error]()

func validateServices(services []application.Service) error {
	for _, service := range services {
		instance := service.Instance()
		if instance == nil || reflect.ValueOf(instance).Kind() != reflect.Pointer || reflect.ValueOf(instance).IsNil() {
			return fmt.Errorf("dreego wails: binding service must be a non-nil pointer; Fix: pass application.NewService(&value)")
		}
		if err := validateService(instance); err != nil {
			return err
		}
	}
	return nil
}

func validateService(instance any) error {
	typeOf := reflect.TypeOf(instance)
	skip := lifecycleMethods(instance)
	for index := 0; index < typeOf.NumMethod(); index++ {
		method := typeOf.Method(index)
		if skip[method.Name] {
			continue
		}
		if err := validateMethod(method); err != nil {
			return fmt.Errorf("dreego wails: binding %s.%s: %w; Fix: use exported JSON-safe values", typeOf.Elem().Name(), method.Name, err)
		}
	}
	return nil
}

func lifecycleMethods(instance any) map[string]bool {
	methods := map[string]bool{"ServeHTTP": true}
	if _, ok := instance.(application.ServiceName); ok {
		methods["ServiceName"] = true
	}
	if _, ok := instance.(application.ServiceStartup); ok {
		methods["ServiceStartup"] = true
	}
	if _, ok := instance.(application.ServiceShutdown); ok {
		methods["ServiceShutdown"] = true
	}
	return methods
}

func validateMethod(method reflect.Method) error {
	typeOf := method.Type
	for index := 1; index < typeOf.NumIn(); index++ {
		parameter := typeOf.In(index)
		if parameter == contextType && index == 1 {
			continue
		}
		if err := validateBindingType(parameter, map[reflect.Type]bool{}); err != nil {
			return fmt.Errorf("parameter %d has unsupported type %s: %w", index, parameter, err)
		}
	}
	if typeOf.NumOut() > 2 {
		return fmt.Errorf("has %d results, want at most a value and error", typeOf.NumOut())
	}
	if typeOf.NumOut() == 2 && typeOf.Out(1) != errorType {
		return fmt.Errorf("second result is %s, want error", typeOf.Out(1))
	}
	for index := 0; index < typeOf.NumOut(); index++ {
		result := typeOf.Out(index)
		if result == errorType && index == typeOf.NumOut()-1 {
			continue
		}
		if err := validateBindingType(result, map[reflect.Type]bool{}); err != nil {
			return fmt.Errorf("result %d has unsupported type %s: %w", index+1, result, err)
		}
	}
	return nil
}

func validateBindingType(typeOf reflect.Type, visiting map[reflect.Type]bool) error {
	if visiting[typeOf] {
		return fmt.Errorf("recursive type %s is not supported", typeOf)
	}
	visiting[typeOf] = true
	defer delete(visiting, typeOf)
	switch typeOf.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return nil
	case reflect.Pointer, reflect.Slice, reflect.Array:
		return validateBindingType(typeOf.Elem(), visiting)
	case reflect.Map:
		if typeOf.Key().Kind() != reflect.String {
			return fmt.Errorf("map keys must be strings")
		}
		return validateBindingType(typeOf.Elem(), visiting)
	case reflect.Struct:
		for index := 0; index < typeOf.NumField(); index++ {
			field := typeOf.Field(index)
			if field.PkgPath != "" || strings.Split(field.Tag.Get("json"), ",")[0] == "-" {
				continue
			}
			if err := validateBindingType(field.Type, visiting); err != nil {
				return fmt.Errorf("field %s: %w", field.Name, err)
			}
		}
		return nil
	default:
		return fmt.Errorf("kind %s is not supported", typeOf.Kind())
	}
}
