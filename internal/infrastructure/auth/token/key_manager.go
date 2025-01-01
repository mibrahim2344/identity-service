package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/mibrahim2344/identity-service/internal/domain/services"
)

// KeyManager defines the interface for managing signing keys
type KeyManager interface {
	// GetSigningKey returns the signing key for the given token type
	GetSigningKey(ctx context.Context, tokenType services.TokenType) ([]byte, error)
	
	// RotateKey rotates the signing key for the given token type
	RotateKey(ctx context.Context, tokenType services.TokenType) error
}

// LocalKeyManager implements KeyManager using local storage
type LocalKeyManager struct {
	keys  map[services.TokenType][]byte
	mutex sync.RWMutex
}

// NewLocalKeyManager creates a new LocalKeyManager
func NewLocalKeyManager() *LocalKeyManager {
	return &LocalKeyManager{
		keys: make(map[services.TokenType][]byte),
	}
}

// GetSigningKey returns the signing key for the given token type
func (m *LocalKeyManager) GetSigningKey(ctx context.Context, tokenType services.TokenType) ([]byte, error) {
	m.mutex.RLock()
	key, exists := m.keys[tokenType]
	m.mutex.RUnlock()

	if !exists {
		// Generate a new key if one doesn't exist
		if err := m.RotateKey(ctx, tokenType); err != nil {
			return nil, err
		}
		m.mutex.RLock()
		key = m.keys[tokenType]
		m.mutex.RUnlock()
	}

	return key, nil
}

// RotateKey rotates the signing key for the given token type
func (m *LocalKeyManager) RotateKey(ctx context.Context, tokenType services.TokenType) error {
	key := make([]byte, 32) // 256 bits
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	m.mutex.Lock()
	m.keys[tokenType] = key
	m.mutex.Unlock()

	return nil
}

// RedisKeyManager implements KeyManager using Redis for distributed key management
type RedisKeyManager struct {
	cache services.CacheService
	local *LocalKeyManager
}

// NewRedisKeyManager creates a new RedisKeyManager
func NewRedisKeyManager(cache services.CacheService) *RedisKeyManager {
	return &RedisKeyManager{
		cache: cache,
		local: NewLocalKeyManager(),
	}
}

// GetSigningKey returns the signing key for the given token type
func (m *RedisKeyManager) GetSigningKey(ctx context.Context, tokenType services.TokenType) ([]byte, error) {
	var encodedKey string
	err := m.cache.Get(ctx, fmt.Sprintf("signing_key:%s", tokenType), &encodedKey)
	if err != nil {
		// Fallback to local key if Redis is unavailable
		return m.local.GetSigningKey(ctx, tokenType)
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	return key, nil
}

// RotateKey rotates the signing key for the given token type
func (m *RedisKeyManager) RotateKey(ctx context.Context, tokenType services.TokenType) error {
	key := make([]byte, 32) // 256 bits
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)
	err := m.cache.Set(ctx, fmt.Sprintf("signing_key:%s", tokenType), encodedKey, 0)
	if err != nil {
		// Fallback to local key management if Redis is unavailable
		m.local.mutex.Lock()
		m.local.keys[tokenType] = key
		m.local.mutex.Unlock()
	}

	return nil
}
