package redis

import "time"

// CacheConfig implements the services.CacheConfig interface
type CacheConfig struct {
	DefaultTTL time.Duration
	MaxEntries int
	Prefix     string
	Namespace  string
}

// GetDefaultTTL returns the default time-to-live for cache entries
func (c *CacheConfig) GetDefaultTTL() time.Duration {
	return c.DefaultTTL
}

// GetMaxEntries returns the maximum number of entries allowed in the cache
func (c *CacheConfig) GetMaxEntries() int {
	return c.MaxEntries
}

// GetPrefix returns the prefix to use for all cache keys
func (c *CacheConfig) GetPrefix() string {
	return c.Prefix
}

// GetNamespace returns the namespace for the cache
func (c *CacheConfig) GetNamespace() string {
	return c.Namespace
}

// NewCacheConfig creates a new cache configuration
func NewCacheConfig(defaultTTL time.Duration, maxEntries int, prefix, namespace string) *CacheConfig {
	return &CacheConfig{
		DefaultTTL: defaultTTL,
		MaxEntries: maxEntries,
		Prefix:     prefix,
		Namespace:  namespace,
	}
}
