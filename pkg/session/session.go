package session

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type ctxKey struct{}

type storeCtxKey struct{}

func WithStore(r *http.Request, s Store) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), storeCtxKey{}, s))
}

func StoreFromCtx(ctx context.Context) Store {
	s, _ := ctx.Value(storeCtxKey{}).(Store)
	return s
}

type Store interface {
	Get(r *http.Request, key string) (string, error)
	Set(w http.ResponseWriter, r *http.Request, key, value string, opts *Options) error
	Delete(w http.ResponseWriter, r *http.Request, key string) error
	Destroy(w http.ResponseWriter, r *http.Request) error
}

type Options struct {
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Path     string
}

type CookieStore struct {
	secret []byte
	name   string
}

func NewCookieStore(secret []byte) *CookieStore {
	return &CookieStore{secret: secret, name: "dreego_session"}
}

func (s *CookieStore) Get(r *http.Request, key string) (string, error) {
	m := s.load(r)
	return m[key], nil
}

func (s *CookieStore) Set(w http.ResponseWriter, r *http.Request, key, value string, opts *Options) error {
	m := s.load(r)
	if value == "" {
		delete(m, key)
	} else {
		m[key] = value
	}
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKey{}, m))
	encoded := s.sign(m)
	http.SetCookie(w, &http.Cookie{
		Name:     s.name,
		Value:    encoded,
		MaxAge:   opt(opts, func(o *Options) int { return o.MaxAge }),
		Secure:   opt(opts, func(o *Options) bool { return o.Secure }),
		HttpOnly: opt(opts, func(o *Options) bool { return o.HttpOnly }),
		Path:     optStr(opts, func(o *Options) string { return o.Path }, "/"),
	})
	return nil
}

func (s *CookieStore) load(r *http.Request) map[string]string {
	if m, ok := r.Context().Value(ctxKey{}).(map[string]string); ok {
		return m
	}
	ck, err := r.Cookie(s.name)
	if err != nil {
		return map[string]string{}
	}
	data, ok := s.verify(ck.Value)
	if !ok {
		return map[string]string{}
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]string{}
	}
	return m
}

func (s *CookieStore) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	return s.Set(w, r, key, "", nil)
}

func (s *CookieStore) Destroy(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:   s.name,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	return nil
}

func (s *CookieStore) sign(m map[string]string) string {
	data, _ := json.Marshal(m)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(data)
	sig := mac.Sum(nil)
	combined := append(sig, byte('.'))
	combined = append(combined, data...)
	return base64.RawURLEncoding.EncodeToString(combined)
}

func (s *CookieStore) verify(value string) ([]byte, bool) {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, false
	}
	idx := strings.IndexByte(string(decoded), '.')
	if idx < 0 || idx != sha256.Size {
		return nil, false
	}
	sig := decoded[:idx]
	data := decoded[idx+1:]
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(data)
	if !hmac.Equal(mac.Sum(nil), sig) {
		return nil, false
	}
	return data, true
}

func opt[T any](opts *Options, fn func(*Options) T) T {
	if opts == nil {
		var zero T
		return zero
	}
	return fn(opts)
}

func optStr(opts *Options, fn func(*Options) string, def string) string {
	if opts == nil {
		return def
	}
	v := fn(opts)
	if v == "" {
		return def
	}
	return v
}
