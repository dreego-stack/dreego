package core

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"strings"
	"testing"
)

// Compile-time assertion: fakeStorage must satisfy the full Storage contract.
var _ Storage = (*fakeStorage)(nil)

type fakeStorage struct {
	data map[string][]byte
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{data: make(map[string][]byte)}
}

func (f *fakeStorage) Put(ctx context.Context, key string, r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	f.data[key] = b
	return nil
}

func (f *fakeStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	b, ok := f.data[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return io.NopCloser(bytes.NewReader(b)), nil
}

func (f *fakeStorage) Delete(ctx context.Context, key string) error {
	delete(f.data, key)
	return nil
}

func (f *fakeStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	for k := range f.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys, nil
}

func (f *fakeStorage) URL(ctx context.Context, key string) (string, error) {
	return "https://example.com/" + key, nil
}

func TestStorageDeleteIsIdempotent(t *testing.T) {
	s := newFakeStorage()
	ctx := context.Background()

	if err := s.Delete(ctx, "missing"); err != nil {
		t.Fatalf("first delete of missing key: %v", err)
	}
	if err := s.Delete(ctx, "missing"); err != nil {
		t.Fatalf("second delete of missing key: %v", err)
	}
}

func TestStorageListFiltersByPrefix(t *testing.T) {
	s := newFakeStorage()
	ctx := context.Background()

	for _, key := range []string{"a/1", "a/2", "b/1"} {
		if err := s.Put(ctx, key, strings.NewReader(key)); err != nil {
			t.Fatalf("put %q: %v", key, err)
		}
	}

	keys, err := s.List(ctx, "a/")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	want := []string{"a/1", "a/2"}
	if len(keys) != len(want) {
		t.Fatalf("list returned %v, want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("list returned %v, want %v", keys, want)
		}
	}
}

func TestStorageGetMissingKeyReturnsError(t *testing.T) {
	s := newFakeStorage()

	rc, err := s.Get(context.Background(), "missing")
	if err == nil {
		rc.Close()
		t.Fatal("get of missing key returned nil error, want error")
	}
}

func TestStoragePutGetRoundtrip(t *testing.T) {
	s := newFakeStorage()
	ctx := context.Background()
	key := "dir/file.txt"
	content := "hello storage"

	if err := s.Put(ctx, key, strings.NewReader(content)); err != nil {
		t.Fatalf("put: %v", err)
	}

	rc, err := s.Get(ctx, key)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != content {
		t.Fatalf("roundtrip content = %q, want %q", got, content)
	}
}

func TestStorageURLReturnsNonEmpty(t *testing.T) {
	s := newFakeStorage()

	u, err := s.URL(context.Background(), "dir/file.txt")
	if err != nil {
		t.Fatalf("url: %v", err)
	}
	if u == "" {
		t.Fatal("url returned empty string, want non-empty")
	}
}
