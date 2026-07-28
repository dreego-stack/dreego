package core

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

func (c *SSRContext) JSON(status int, data any) {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(status)
	if err := json.NewEncoder(c.W).Encode(data); err != nil {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

func (c *SSRContext) XML(status int, data any) {
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(status)
	if err := xml.NewEncoder(c.W).Encode(data); err != nil {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

func (c *SSRContext) Bind(target any) error {
	return json.NewDecoder(c.R.Body).Decode(target)
}

func (c *SSRContext) Write(status int, contentType string, body []byte) {
	c.W.Header().Set("Content-Type", contentType+"; charset=utf-8")
	c.W.WriteHeader(status)
	c.W.Write(body)
}

func (c *SSRContext) Wants(mime string) bool {
	accept := c.R.Header.Get("Accept")
	if accept == "*/*" || accept == "" {
		return true
	}
	return stringsContainsMime(accept, mime)
}

func stringsContainsMime(accept, mime string) bool {
	for _, part := range strings.Split(accept, ",") {
		t := strings.TrimSpace(part)
		if idx := strings.IndexByte(t, ';'); idx >= 0 {
			t = t[:idx]
		}
		if t == mime || t == "*/*" {
			return true
		}
	}
	return false
}
