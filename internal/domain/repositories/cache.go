package repositories

import "time"

// CacheService defines the interface for caching operations
type CacheService interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, value interface{}) error
	Delete(key string) error
}
