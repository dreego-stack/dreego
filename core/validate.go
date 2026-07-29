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
		if val != "" {
			v.Field(i).SetString(val)
		}
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
		if len(val) < atoi(min) {
			return "must be at least " + min + " characters"
		}
	case strings.HasPrefix(rule, "max="):
		max := strings.TrimPrefix(rule, "max=")
		if len(val) > atoi(max) {
			return "must be at most " + max + " characters"
		}
	}
	return ""
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
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
