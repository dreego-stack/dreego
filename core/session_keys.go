package core

import (
	"crypto/sha256"
)

type sessionKeys struct {
	sig []byte
	enc []byte
}

func deriveKeys(secret []byte) sessionKeys {
	h := sha256.New()
	h.Write(secret)
	h.Write([]byte("dreego-session-sig"))
	sig := h.Sum(nil)

	h2 := sha256.New()
	h2.Write(secret)
	h2.Write([]byte("dreego-session-enc"))
	enc := h2.Sum(nil)

	return sessionKeys{sig: sig, enc: enc}
}
