package core

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

var _ KVStore = (*fakeKVStore)(nil)

type kvEntry struct {
	value  []byte
	expiry time.Time
}

type fakeKVStore struct {
	mu    sync.Mutex
	now   func() time.Time
	store map[string]kvEntry
}

func newFakeKVStore(now func() time.Time) *fakeKVStore {
	return &fakeKVStore{
		now:   now,
		store: make(map[string]kvEntry),
	}
}

func (f *fakeKVStore) Get(ctx context.Context, key string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.store[key]
	if !ok {
		return nil, errors.New("kvstore: key not found")
	}
	if !e.expiry.IsZero() && !f.now().Before(e.expiry) {
		delete(f.store, key)
		return nil, errors.New("kvstore: key not found")
	}
	return e.value, nil
}

func (f *fakeKVStore) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	e := kvEntry{value: val}
	if ttl > 0 {
		e.expiry = f.now().Add(ttl)
	}
	f.store[key] = e
	return nil
}

func (f *fakeKVStore) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.store, key)
	return nil
}

func (f *fakeKVStore) Expire(ctx context.Context, key string, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.store[key]
	if !ok {
		return errors.New("kvstore: key not found")
	}
	if ttl > 0 {
		e.expiry = f.now().Add(ttl)
	} else {
		e.expiry = time.Time{}
	}
	f.store[key] = e
	return nil
}

func TestKVStoreGetSetRoundtrip(t *testing.T) {
	now := time.Now()
	f := newFakeKVStore(func() time.Time { return now })
	ctx := context.Background()

	if err := f.Set(ctx, "name", []byte("dreego"), time.Minute); err != nil {
		t.Fatalf("set: %v", err)
	}
	got, err := f.Get(ctx, "name")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if string(got) != "dreego" {
		t.Fatalf("got %q, want %q", got, "dreego")
	}
}

func TestKVStoreGetMissingKeyReturnsError(t *testing.T) {
	f := newFakeKVStore(time.Now)
	_, err := f.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("get on missing key: want error, got nil")
	}
}

func TestKVStoreGetExpiredReturnsError(t *testing.T) {
	now := time.Now()
	f := newFakeKVStore(func() time.Time { return now })
	ctx := context.Background()

	if err := f.Set(ctx, "tmp", []byte("data"), time.Second); err != nil {
		t.Fatalf("set: %v", err)
	}
	f.now = func() time.Time { return now.Add(2 * time.Second) }

	_, err := f.Get(ctx, "tmp")
	if err == nil {
		t.Fatal("get on expired key: want error, got nil")
	}
}

func TestKVStoreSetTTLZeroKeepsForever(t *testing.T) {
	now := time.Now()
	f := newFakeKVStore(func() time.Time { return now })
	ctx := context.Background()

	if err := f.Set(ctx, "forever", []byte("data"), 0); err != nil {
		t.Fatalf("set: %v", err)
	}
	f.now = func() time.Time { return now.Add(100 * time.Hour) }

	got, err := f.Get(ctx, "forever")
	if err != nil {
		t.Fatalf("get after long time: %v", err)
	}
	if string(got) != "data" {
		t.Fatalf("got %q, want %q", got, "data")
	}
}

func TestKVStoreDeleteIdempotent(t *testing.T) {
	f := newFakeKVStore(time.Now)
	ctx := context.Background()

	if err := f.Delete(ctx, "missing"); err != nil {
		t.Fatalf("first delete: %v", err)
	}
	if err := f.Delete(ctx, "missing"); err != nil {
		t.Fatalf("second delete: %v", err)
	}
}

func TestKVStoreExpireUpdatesTTL(t *testing.T) {
	now := time.Now()
	f := newFakeKVStore(func() time.Time { return now })
	ctx := context.Background()

	if err := f.Set(ctx, "key", []byte("data"), time.Second); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := f.Expire(ctx, "key", time.Hour); err != nil {
		t.Fatalf("expire: %v", err)
	}

	f.now = func() time.Time { return now.Add(30 * time.Second) }
	got, err := f.Get(ctx, "key")
	if err != nil {
		t.Fatalf("get after re-expire: %v", err)
	}
	if string(got) != "data" {
		t.Fatalf("got %q, want %q", got, "data")
	}
}

func TestKVStoreExpireMissingKeyReturnsError(t *testing.T) {
	f := newFakeKVStore(time.Now)
	err := f.Expire(context.Background(), "missing", time.Minute)
	if err == nil {
		t.Fatal("expire on missing key: want error, got nil")
	}
}
