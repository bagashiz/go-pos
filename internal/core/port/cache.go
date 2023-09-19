package port

import (
	"context"
	"time"
)

// CacheService is an interface for interacting with cache-related business logic
type CacheService interface {
	// Set stores the value in the cache
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Get retrieves the value from the cache
	Get(ctx context.Context, key string) ([]byte, error)
	// Delete removes the value from the cache
	Delete(ctx context.Context, key string) error
	// DeleteByPrefix removes the value from the cache with the given prefix
	DeleteByPrefix(ctx context.Context, prefix string) error
	// Close closes the connection to the cache server
	Close() error
}
