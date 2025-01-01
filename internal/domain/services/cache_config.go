package services

import "time"

// CacheConfig defines the configuration interface for cache services
type CacheConfig interface {
	// GetDefaultTTL returns the default time-to-live for cache entries
	GetDefaultTTL() time.Duration

	// GetMaxEntries returns the maximum number of entries allowed in the cache
	GetMaxEntries() int

	// GetPrefix returns the prefix to use for all cache keys
	GetPrefix() string

	// GetNamespace returns the namespace for the cache (e.g., "users", "sessions")
	GetNamespace() string
}
