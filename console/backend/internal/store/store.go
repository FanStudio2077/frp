package store

import (
	"context"
	"errors"
)

// ErrNotImplemented marks scaffold methods.
var ErrNotImplemented = errors.New("store: not implemented yet")

// Store is the persistence boundary for the console.
// SQLite is the default implementation; keep this interface free of driver
// details so Postgres can be added later without touching handlers.
type Store interface {
	Close() error
	Ping(ctx context.Context) error
}
