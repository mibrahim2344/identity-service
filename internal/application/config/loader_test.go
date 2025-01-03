package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mibrahim2344/identity-service/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"database": {
			"host": "db.example.com",
			"port": 5432,
			"user": "test_user",
			"password": "test_password",
			"dbname": "test_db",
			"sslmode": "disable",
			"maxIdleConns": 10,
			"maxOpenConns": 100,
			"connMaxLifetimeMinutes": 60
		},
		"redis": {
			"host": "redis.example.com",
			"port": 6379,
			"password": "redis_password",
			"db": 1
		},
		"kafka": {
			"brokers": ["kafka1:9092", "kafka2:9092"],
			"topic": "test_topic"
		},
		"auth": {
			"accessTokenDuration": 30,
			"refreshTokenDuration": 20160,
			"signingKey": "test_signing_key",
			"hashingCost": 12
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Run("Load from file", func(t *testing.T) {
		config, err := LoadConfig(configPath)
		require.NoError(t, err)

		// Verify config values
		assert.Equal(t, "db.example.com", config.Database.Host)
		assert.Equal(t, 5432, config.Database.Port)
		assert.Equal(t, "test_user", config.Database.User)
		assert.Equal(t, "test_password", config.Database.Password)
		assert.Equal(t, "test_db", config.Database.DBName)
		assert.Equal(t, "disable", config.Database.SSLMode)
		assert.Equal(t, 10, config.Database.MaxIdleConns)
		assert.Equal(t, 100, config.Database.MaxOpenConns)
		assert.Equal(t, 60, config.Database.ConnMaxLifetimeMinutes)

		assert.Equal(t, "redis.example.com", config.Redis.Host)
		assert.Equal(t, 6379, config.Redis.Port)
		assert.Equal(t, "redis_password", config.Redis.Password)
		assert.Equal(t, 1, config.Redis.DB)

		assert.Equal(t, []string{"kafka1:9092", "kafka2:9092"}, config.Kafka.Brokers)
		assert.Equal(t, "test_topic", config.Kafka.Topic)

		assert.Equal(t, 30, config.Auth.AccessTokenDuration)
		assert.Equal(t, 20160, config.Auth.RefreshTokenDuration)
		assert.Equal(t, "test_signing_key", config.Auth.SigningKey)
		assert.Equal(t, 12, config.Auth.HashingCost)
	})

	t.Run("Override with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("DB_HOST", "db2.example.com")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_MAX_IDLE_CONNS", "20")
		os.Setenv("DB_MAX_OPEN_CONNS", "200")
		os.Setenv("DB_CONN_MAX_LIFETIME_MINUTES", "120")
		os.Setenv("REDIS_PASSWORD", "new_password")
		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_MAX_IDLE_CONNS")
			os.Unsetenv("DB_MAX_OPEN_CONNS")
			os.Unsetenv("DB_CONN_MAX_LIFETIME_MINUTES")
			os.Unsetenv("REDIS_PASSWORD")
		}()

		config, err := LoadConfig(configPath)
		require.NoError(t, err)

		// Verify environment variables override file config
		assert.Equal(t, "db2.example.com", config.Database.Host)
		assert.Equal(t, 5433, config.Database.Port)
		assert.Equal(t, 20, config.Database.MaxIdleConns)
		assert.Equal(t, 200, config.Database.MaxOpenConns)
		assert.Equal(t, 120, config.Database.ConnMaxLifetimeMinutes)
		assert.Equal(t, "new_password", config.Redis.Password)
	})

	t.Run("Invalid config file path", func(t *testing.T) {
		_, err := LoadConfig("nonexistent.json")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config file")
	})

	t.Run("Invalid config file content", func(t *testing.T) {
		invalidConfigPath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(invalidConfigPath, []byte("invalid json"), 0644)
		require.NoError(t, err)

		_, err = LoadConfig(invalidConfigPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config file")
	})

	t.Run("Missing required fields", func(t *testing.T) {
		emptyConfigPath := filepath.Join(tmpDir, "empty.json")
		err := os.WriteFile(emptyConfigPath, []byte("{}"), 0644)
		require.NoError(t, err)

		_, err = LoadConfig(emptyConfigPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid configuration")
	})
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      func() application.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: func() application.Config {
				return application.Config{
					Database: struct {
						Host     string
						Port     int
						User     string
						Password string
						DBName   string
						SSLMode  string
						MaxIdleConns       int
						MaxOpenConns       int
						ConnMaxLifetimeMinutes int
					}{
						Host:   "localhost",
						Port:   5432,
						User:   "user",
						DBName: "dbname",
						MaxIdleConns:       10,
						MaxOpenConns:       100,
						ConnMaxLifetimeMinutes: 60,
					},
					Redis: struct {
						Host     string
						Port     int
						Password string
						DB       int
					}{
						Host: "localhost",
						Port: 6379,
					},
					Kafka: struct {
						Brokers []string
						Topic   string
					}{
						Brokers: []string{"localhost:9092"},
						Topic:   "topic",
					},
					Auth: struct {
						AccessTokenDuration  int
						RefreshTokenDuration int
						SigningKey           string
						HashingCost          int
					}{
						AccessTokenDuration:  15,
						RefreshTokenDuration: 10080,
						SigningKey:           "key",
					},
				}
			},
			expectError: false,
		},
		{
			name: "Missing database host",
			config: func() application.Config {
				c := application.Config{}
				c.Database.Port = 5432
				c.Database.User = "user"
				c.Database.DBName = "dbname"
				return c
			},
			expectError: true,
			errorMsg:    "database host is required",
		},
		{
			name: "Missing database port",
			config: func() application.Config {
				c := application.Config{}
				c.Database.Host = "localhost"
				c.Database.User = "user"
				c.Database.DBName = "dbname"
				return c
			},
			expectError: true,
			errorMsg:    "database port is required",
		},
		{
			name: "Missing redis host",
			config: func() application.Config {
				c := application.Config{}
				c.Redis.Port = 6379
				return c
			},
			expectError: true,
			errorMsg:    "redis host is required",
		},
		{
			name: "Missing kafka brokers",
			config: func() application.Config {
				c := application.Config{}
				c.Kafka.Topic = "topic"
				return c
			},
			expectError: true,
			errorMsg:    "kafka brokers are required",
		},
		{
			name: "Default hashing cost",
			config: func() application.Config {
				c := application.Config{}
				c.Database.Host = "localhost"
				c.Database.Port = 5432
				c.Database.User = "user"
				c.Database.DBName = "dbname"
				c.Redis.Host = "localhost"
				c.Redis.Port = 6379
				c.Kafka.Brokers = []string{"localhost:9092"}
				c.Kafka.Topic = "topic"
				c.Auth.AccessTokenDuration = 15
				c.Auth.RefreshTokenDuration = 10080
				c.Auth.SigningKey = "key"
				return c
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config())
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
