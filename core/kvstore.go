package core

import (
	"context"
	"time"
)

// KVStore is a key-value store contract in the style of database/sql: core
// defines the interface, plugins provide the implementation (Redis, Ristretto,
// in-memory, ...). It is deliberately distinct from Storage (blobs): KV holds
// small values with an optional TTL, Storage holds larger opaque blobs.
// All methods respect ctx cancellation.
type KVStore interface {
	// Get returns the value stored under key. It returns an error if key
	// does not exist or its ttl has expired.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set stores val under key with the given ttl. A ttl <= 0 means no
	// expiry (keep forever).
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	// Delete removes key. It is idempotent: deleting a missing key is not
	// an error.
	Delete(ctx context.Context, key string) error
	// Expire sets or adjusts the ttl on an existing key. It returns an
	// error if key does not exist.
	Expire(ctx context.Context, key string, ttl time.Duration) error
}
