package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mibrahim2344/identity-service/internal/application"
)

// LoadConfig loads configuration from environment variables and/or config file
func LoadConfig(configPath string) (application.Config, error) {
	var config application.Config

	// First try to load from config file if provided
	if configPath != "" {
		if err := loadFromFile(configPath, &config); err != nil {
			return config, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables if present
	loadFromEnv(&config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return config, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(path string, config *application.Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	type configAlias application.Config
	var rawConfig struct {
		configAlias
		Cache struct {
			DefaultTTL int    `json:"defaultTTL"`
			MaxEntries int    `json:"maxEntries"`
			Prefix     string `json:"prefix"`
			Namespace  string `json:"namespace"`
		} `json:"cache"`
	}

	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	var newConfig application.Config
	newConfig = application.Config(rawConfig.configAlias)
	newConfig.Cache.DefaultTTL = time.Duration(rawConfig.Cache.DefaultTTL) * time.Second
	newConfig.Cache.MaxEntries = rawConfig.Cache.MaxEntries
	newConfig.Cache.Prefix = rawConfig.Cache.Prefix
	newConfig.Cache.Namespace = rawConfig.Cache.Namespace
	*config = newConfig

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *application.Config) {
	// Database configuration
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}
	if maxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConns != "" {
		if mic, err := strconv.Atoi(maxIdleConns); err == nil {
			config.Database.MaxIdleConns = mic
		}
	}
	if maxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConns != "" {
		if moc, err := strconv.Atoi(maxOpenConns); err == nil {
			config.Database.MaxOpenConns = moc
		}
	}
	if connMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME_MINUTES"); connMaxLifetime != "" {
		if cml, err := strconv.Atoi(connMaxLifetime); err == nil {
			config.Database.ConnMaxLifetimeMinutes = cml
		}
	}

	// Redis configuration
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			config.Redis.DB = d
		}
	}

	// Kafka configuration
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		config.Kafka.Brokers = strings.Split(brokers, ",")
	}
	if topic := os.Getenv("KAFKA_TOPIC"); topic != "" {
		config.Kafka.Topic = topic
	}

	// Auth configuration
	if duration := os.Getenv("AUTH_ACCESS_TOKEN_DURATION"); duration != "" {
		if d, err := strconv.Atoi(duration); err == nil {
			config.Auth.AccessTokenDuration = d
		}
	}
	if duration := os.Getenv("AUTH_REFRESH_TOKEN_DURATION"); duration != "" {
		if d, err := strconv.Atoi(duration); err == nil {
			config.Auth.RefreshTokenDuration = d
		}
	}
	if key := os.Getenv("AUTH_SIGNING_KEY"); key != "" {
		config.Auth.SigningKey = key
	}
	if cost := os.Getenv("AUTH_HASHING_COST"); cost != "" {
		if c, err := strconv.Atoi(cost); err == nil {
			config.Auth.HashingCost = c
		}
	}
}

// validateConfig validates the configuration
func validateConfig(config application.Config) error {
	// Database validation
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port == 0 {
		return fmt.Errorf("database port is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Redis validation
	if config.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}
	if config.Redis.Port == 0 {
		return fmt.Errorf("redis port is required")
	}

	// Kafka validation
	if len(config.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka brokers are required")
	}
	if config.Kafka.Topic == "" {
		return fmt.Errorf("kafka topic is required")
	}

	// Auth validation
	if config.Auth.AccessTokenDuration == 0 {
		return fmt.Errorf("access token duration is required")
	}
	if config.Auth.RefreshTokenDuration == 0 {
		return fmt.Errorf("refresh token duration is required")
	}
	if config.Auth.SigningKey == "" {
		return fmt.Errorf("auth signing key is required")
	}
	if config.Auth.HashingCost == 0 {
		config.Auth.HashingCost = 10 // Set default bcrypt cost
	}

	return nil
}
