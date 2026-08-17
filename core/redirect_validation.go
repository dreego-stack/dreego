package core

import (
	"fmt"
	"net/http"
	"strings"
)

var validRedirectStatus = map[int]struct{}{
	http.StatusMovedPermanently:  {},
	http.StatusFound:             {},
	http.StatusSeeOther:          {},
	http.StatusTemporaryRedirect: {},
	http.StatusPermanentRedirect: {},
}

func validateRedirect(from, to string, status int) error {
	if err := validateRedirectPattern(from); err != nil {
		return err
	}
	if err := validateRedirectTarget(to); err != nil {
		return err
	}
	if _, ok := validRedirectStatus[status]; !ok {
		return fmt.Errorf("dreego: invalid redirect status %d (allowed: 301,302,303,307,308)", status)
	}
	if err := validateRedirectPair(from, to); err != nil {
		return err
	}
	if redirectLoops(from, to) {
		return fmt.Errorf("dreego: redirect %q -> %q loops back to itself", from, to)
	}
	return nil
}

func validateRewrite(from, to string) error {
	if err := validateRedirectPattern(from); err != nil {
		return err
	}
	if err := validateRedirectTarget(to); err != nil {
		return err
	}
	if err := validateRedirectPair(from, to); err != nil {
		return err
	}
	if redirectLoops(from, to) {
		return fmt.Errorf("dreego: rewrite %q -> %q loops back to itself", from, to)
	}
	return nil
}

func validateRedirectPair(from, to string) error {
	if strings.HasSuffix(from, "/*") && to == "/" {
		return fmt.Errorf("dreego: redirect/rewrite target %q with wildcard pattern %q emits //-prefixed targets", to, from)
	}
	if !strings.HasSuffix(from, "/*") && strings.HasSuffix(to, "/*") {
		return fmt.Errorf("dreego: redirect/rewrite target %q has wildcard but pattern %q is exact", to, from)
	}
	return nil
}

func validateRedirectPattern(p string) error {
	if p == "" {
		return fmt.Errorf("dreego: redirect/rewrite pattern is empty")
	}
	if !strings.HasPrefix(p, "/") {
		return fmt.Errorf("dreego: redirect/rewrite pattern %q must start with /", p)
	}
	if strings.Contains(p, "//") {
		return fmt.Errorf("dreego: redirect/rewrite pattern %q contains //", p)
	}
	if strings.HasSuffix(p, "/") && p != "/" {
		return fmt.Errorf("dreego: redirect/rewrite pattern %q must not end with /", p)
	}
	wildcard := strings.HasSuffix(p, "/*")
	base := p
	if wildcard {
		base = strings.TrimSuffix(p, "/*")
	}
	if wildcard && base == "" {
		return fmt.Errorf("dreego: redirect/rewrite pattern %q has empty prefix", p)
	}
	if !wildcard && strings.Contains(p, "*") {
		return fmt.Errorf("dreego: redirect/rewrite pattern %q invalid wildcard (use /*)", p)
	}
	if wildcard && strings.Count(p, "/*") > 1 {
		return fmt.Errorf("dreego: redirect/rewrite pattern %q has multiple wildcards", p)
	}
	return nil
}

func validateRedirectTarget(p string) error {
	if p == "" {
		return fmt.Errorf("dreego: redirect/rewrite target is empty")
	}
	if !strings.HasPrefix(p, "/") {
		return fmt.Errorf("dreego: redirect/rewrite target %q must start with /", p)
	}
	if strings.Contains(p, "//") {
		return fmt.Errorf("dreego: redirect/rewrite target %q contains //", p)
	}
	if strings.HasSuffix(p, "/") && p != "/" {
		return fmt.Errorf("dreego: redirect/rewrite target %q must not end with /", p)
	}
	if strings.Contains(p, "*") && !strings.HasSuffix(p, "/*") {
		return fmt.Errorf("dreego: redirect/rewrite target %q invalid wildcard (use /*)", p)
	}
	if strings.Count(p, "/*") > 1 {
		return fmt.Errorf("dreego: redirect/rewrite target %q has multiple wildcards", p)
	}
	return nil
}

func redirectLoops(from, to string) bool {
	if from == to {
		return true
	}
	if !strings.HasSuffix(from, "/*") {
		return false
	}
	fromPrefix := strings.TrimSuffix(from, "/*")
	if !strings.HasSuffix(to, "/*") {
		return to == fromPrefix || strings.HasPrefix(to, fromPrefix+"/")
	}
	toPrefix := strings.TrimSuffix(to, "/*")
	return toPrefix == fromPrefix || strings.HasPrefix(toPrefix, fromPrefix+"/")
}
