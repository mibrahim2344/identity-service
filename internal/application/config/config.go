package config

import "time"

// Config holds all the configuration needed for the application services
type Config struct {
	Database struct {
		Host                   string
		Port                   int
		User                   string
		Password               string
		DBName                 string
		SSLMode                string
		MaxIdleConns           int
		MaxOpenConns           int
		ConnMaxLifetimeMinutes int
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	Cache struct {
		DefaultTTL time.Duration
		MaxEntries int
		Prefix     string
		Namespace  string
	}
	Kafka struct {
		Brokers []string
		Topic   string
	}
	Auth struct {
		AccessTokenDuration  time.Duration
		RefreshTokenDuration time.Duration
		SigningKey           string
		HashingCost          int
	}
	Server struct {
		Host string
		Port int
	}
	WebApp struct {
		URL string
	}
}
