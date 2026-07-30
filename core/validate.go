package core

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type FieldError struct {
	Field   string
	Message string
}

func BindForm(r *http.Request, target any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	v := reflect.ValueOf(target).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		val := r.FormValue(tag)
		if val == "" {
			continue
		}
		if field.Type.Kind() != reflect.String {
			return fmt.Errorf("unsupported field type %s for field %s", field.Type.Kind(), field.Name)
		}
		v.Field(i).SetString(val)
	}
	return nil
}

func ValidateForm(form any) map[string]string {
	t := reflect.TypeOf(form)
	v := reflect.ValueOf(form)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	errs := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		formTag := field.Tag.Get("form")
		if formTag == "" {
			formTag = strings.ToLower(field.Name)
		}
		val := v.Field(i).String()
		for _, rule := range strings.Split(tag, ",") {
			rule = strings.TrimSpace(rule)
			if msg := applyRule(rule, val); msg != "" {
				errs[formTag] = msg
				break
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func applyRule(rule string, val string) string {
	switch {
	case rule == "required":
		if strings.TrimSpace(val) == "" {
			return "is required"
		}
	case rule == "email":
		if !strings.Contains(val, "@") || !strings.Contains(val, ".") {
			return "must be a valid email"
		}
	case strings.HasPrefix(rule, "min="):
		min := strings.TrimPrefix(rule, "min=")
		n, err := atoi(min)
		if err != nil {
			return "min must be a valid number"
		}
		if len(val) < n {
			return "must be at least " + min + " characters"
		}
	case strings.HasPrefix(rule, "max="):
		max := strings.TrimPrefix(rule, "max=")
		n, err := atoi(max)
		if err != nil {
			return "max must be a valid number"
		}
		if len(val) > n {
			return "must be at most " + max + " characters"
		}
	}
	return ""
}

func atoi(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty number")
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid digit %q", c)
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func SaveOld(c *SSRContext, form any) {
	v := reflect.ValueOf(form)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("form")
		if tag == "" {
			tag = strings.ToLower(t.Field(i).Name)
		}
		c.Set("old_"+tag, fmt.Sprint(v.Field(i).Interface()))
	}
}

func SaveErrors(c *SSRContext, errs map[string]string) {
	for k, v := range errs {
		c.Set("error_"+k, v)
	}
}
