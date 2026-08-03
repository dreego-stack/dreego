package core

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"testing"
)

func TestKeyDerivationUsesHMAC(t *testing.T) {
	secret := []byte("my-secret")
	keys := deriveKeys(secret)

	if bytes.Equal(keys.sig, keys.enc) {
		t.Error("signing and encryption keys must differ")
	}

	macSig := hmac.New(sha256.New, secret)
	macSig.Write([]byte("dreego-session-sig"))
	expectedSig := macSig.Sum(nil)

	macEnc := hmac.New(sha256.New, secret)
	macEnc.Write([]byte("dreego-session-enc"))
	expectedEnc := macEnc.Sum(nil)

	if !bytes.Equal(keys.sig, expectedSig) {
		t.Error("signing key not derived via HMAC-SHA256 of secret with label")
	}
	if !bytes.Equal(keys.enc, expectedEnc) {
		t.Error("encryption key not derived via HMAC-SHA256 of secret with label")
	}
}
