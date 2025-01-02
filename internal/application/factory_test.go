package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewFactory(t *testing.T) {
	// Create test config
	config := Config{
		Database: struct {
			Host     string
			Port     int
			User     string
			Password string
			DBName   string
			SSLMode  string
		}{
			Host:     "localhost",
			Port:     5432,
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
		Redis: struct {
			Host     string
			Port     int
			Password string
			DB       int
		}{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Kafka: struct {
			Brokers []string
			Topic   string
		}{
			Brokers: []string{"localhost:9092"},
			Topic:   "test_topic",
		},
		Auth: struct {
			AccessTokenDuration  int
			RefreshTokenDuration int
			SigningKey          string
			HashingCost         int
		}{
			AccessTokenDuration:  15,
			RefreshTokenDuration: 10080,
			SigningKey:          "test_key",
			HashingCost:         10,
		},
	}

	// Create test logger
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	// Create factory
	factory := NewFactory(config, logger)

	// Verify factory is created with correct config
	assert.NotNil(t, factory)
	assert.Equal(t, config, factory.config)
	assert.Equal(t, logger, factory.logger)
}

func TestCreateUserService(t *testing.T) {
	// Create test config with mock values
	config := Config{
		Database: struct {
			Host     string
			Port     int
			User     string
			Password string
			DBName   string
			SSLMode  string
		}{
			Host:     "localhost",
			Port:     5432,
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
		Redis: struct {
			Host     string
			Port     int
			Password string
			DB       int
		}{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Kafka: struct {
			Brokers []string
			Topic   string
		}{
			Brokers: []string{"localhost:9092"},
			Topic:   "test_topic",
		},
		Auth: struct {
			AccessTokenDuration  int
			RefreshTokenDuration int
			SigningKey          string
			HashingCost         int
		}{
			AccessTokenDuration:  15,
			RefreshTokenDuration: 10080,
			SigningKey:          "test_key",
			HashingCost:         10,
		},
	}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	factory := NewFactory(config, logger)

	// Test service creation
	service, err := factory.CreateUserService()
	
	// We expect an error because we're not actually connecting to the database
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "failed to create database connection")
}

func TestCreateEmailService(t *testing.T) {
	config := Config{}
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	factory := NewFactory(config, logger)

	service, err := factory.CreateEmailService()
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestDefaultCacheConfig(t *testing.T) {
	config := &defaultCacheConfig{}

	// Test default values
	assert.Equal(t, 24*60*60, int(config.GetDefaultTTL().Seconds()))
	assert.Equal(t, 10000, config.GetMaxEntries())
	assert.Equal(t, "identity", config.GetPrefix())
	assert.Equal(t, "users", config.GetNamespace())
}

func TestClose(t *testing.T) {
	config := Config{}
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	factory := NewFactory(config, logger)

	err = factory.Close()
	assert.NoError(t, err)
}
