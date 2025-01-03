package services

import (
	"context"
	"time"
)

// CacheService defines the interface for caching operations
type CacheService interface {
	// Set stores a value in the cache with the given key and expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get retrieves a value from the cache by key
	Get(ctx context.Context, key string, dest interface{}) error

	// Delete removes a value from the cache by key
	Delete(ctx context.Context, key string) error

	// Clear removes all values from the cache
	Clear(ctx context.Context) error

	// SetNX sets a value in the cache only if the key doesn't exist
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}

// CacheSettings represents the configuration settings for cache operations
type CacheSettings struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	MaxEntries        int
}
