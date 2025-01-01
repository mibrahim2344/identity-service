package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/mibrahim2344/identity-service/internal/domain/services"
)

// CacheService implements the domain.CacheService interface using Redis
type CacheService struct {
	client *redis.Client
	config services.CacheConfig
}

// NewCacheService creates a new Redis cache service
func NewCacheService(client *redis.Client, config services.CacheConfig) *CacheService {
	return &CacheService{
		client: client,
		config: config,
	}
}

// Set stores a value in the cache with the given key and expiration
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	if err := s.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	return nil
}

// Get retrieves a value from the cache by key
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return services.ErrCacheKeyNotFound
		}
		return fmt.Errorf("failed to get cache value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// Delete removes a value from the cache by key
func (s *CacheService) Delete(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache key: %w", err)
	}
	return nil
}

// Clear removes all values from the cache
func (s *CacheService) Clear(ctx context.Context) error {
	if err := s.client.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

// SetNX sets a value in the cache only if the key doesn't exist
func (s *CacheService) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal cache value: %w", err)
	}

	success, err := s.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set cache value with NX: %w", err)
	}

	return success, nil
}

// GetWithTTL retrieves a value and its remaining TTL from the cache
func (s *CacheService) GetWithTTL(ctx context.Context, key string, dest interface{}) (time.Duration, error) {
	pipe := s.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			return 0, services.ErrCacheKeyNotFound
		}
		return 0, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	data, err := getCmd.Bytes()
	if err != nil {
		return 0, fmt.Errorf("failed to get cache value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return 0, fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return ttlCmd.Val(), nil
}

// Increment atomically increments a counter by the given value
func (s *CacheService) Increment(ctx context.Context, key string, value int64) (int64, error) {
	result, err := s.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}
	return result, nil
}

// Decrement atomically decrements a counter by the given value
func (s *CacheService) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	result, err := s.client.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement counter: %w", err)
	}
	return result, nil
}
