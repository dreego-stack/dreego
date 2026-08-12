package core

import (
	"context"
	"io"
)

// Storage is the frozen v1 contract for file storage backends. Like
// database/sql, core defines the contract and plugins implement it
// (S3, R2, Local). Core code never imports a plugin.
type Storage interface {
	// Put streams r under key. Implementations must not read from r after
	// the call returns; the caller must not reuse r afterwards. All data
	// is stored before Put returns.
	Put(ctx context.Context, key string, r io.Reader) error
	// Get returns a stream for key. The returned stream must be closed by
	// the CALLER. Implementations must return an error if key does not
	// exist.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete removes key. It is idempotent: deleting a missing key is not
	// an error.
	Delete(ctx context.Context, key string) error
	// List returns all keys with the given prefix. The order is
	// implementation-defined. v1 has no pagination: implementations must
	// return the complete result in one call.
	List(ctx context.Context, prefix string) ([]string, error)
	// URL returns a usable URL for key. Implementations may return a
	// signed URL or a public one, as defined by the backend.
	URL(ctx context.Context, key string) (string, error)
}
