package core

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

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
	encoded := s.sign(m, opt(opts, func(o *Options) bool { return o.Encrypt }))
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

func (s *CookieStore) sign(m map[string]string, encrypt bool) string {
	data, _ := json.Marshal(m)
	key := deriveKeys(s.secret)
	if encrypt {
		data = append([]byte{encMarker}, encryptPayload(key.enc, data)...)
	}
	mac := hmac.New(sha256.New, key.sig)
	mac.Write(data)
	sig := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(append(sig, data...))
}

func (s *CookieStore) verify(value string) ([]byte, bool) {
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, false
	}
	if len(decoded) < sha256.Size {
		return nil, false
	}
	sig := decoded[:sha256.Size]
	data := decoded[sha256.Size:]
	key := deriveKeys(s.secret)
	mac := hmac.New(sha256.New, key.sig)
	mac.Write(data)
	if !hmac.Equal(mac.Sum(nil), sig) {
		return nil, false
	}
	if len(data) > 0 && data[0] == encMarker {
		return decryptPayload(key.enc, data[1:])
	}
	return data, true
}
