package core

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"
)

func TestDecryptPayloadInvalidBase64(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	_, ok := decryptPayload(key, []byte("!!!not-base64!!!"))
	if ok {
		t.Error("expected decrypt to fail on invalid base64")
	}
}

func TestDecryptPayloadTooShort(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	// Valid base64 but shorter than the GCM nonce size (12 bytes).
	short := base64.RawURLEncoding.EncodeToString([]byte("short"))
	_, ok := decryptPayload(key, []byte(short))
	if ok {
		t.Error("expected decrypt to fail on payload shorter than nonce")
	}
}

func TestDecryptPayloadWrongKey(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	enc, err := encryptPayload(key, []byte("secret"), rand.Reader)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	wrong := make([]byte, 32)
	if _, err := rand.Read(wrong); err != nil {
		t.Fatal(err)
	}
	_, ok := decryptPayload(wrong, enc)
	if ok {
		t.Error("expected decrypt to fail with wrong key")
	}
}

func TestDecryptPayloadRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	plaintext := []byte("round trip payload")
	enc, err := encryptPayload(key, plaintext, rand.Reader)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	pt, ok := decryptPayload(key, enc)
	if !ok {
		t.Fatal("expected decrypt to succeed")
	}
	if !bytes.Equal(pt, plaintext) {
		t.Errorf("expected %q, got %q", plaintext, pt)
	}
}

func TestEncryptPayloadNonceReadError(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	_, err := encryptPayload(key, []byte("data"), &shortReader{})
	if err == nil {
		t.Error("expected encrypt to fail when nonce read fails")
	}
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}
