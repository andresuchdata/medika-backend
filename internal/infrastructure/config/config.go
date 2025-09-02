package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Auth          AuthConfig          `mapstructure:"auth"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Log           LogConfig           `mapstructure:"log"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Prefork      bool          `mapstructure:"prefork"`
}

type DatabaseConfig struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	URL           string        `mapstructure:"url"`
	MaxRetries    int           `mapstructure:"max_retries"`
	MinIdleConns  int           `mapstructure:"min_idle_conns"`
	PoolSize      int           `mapstructure:"pool_size"`
	ReadTimeout   time.Duration `mapstructure:"read_timeout"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"`
}

type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	TokenDuration time.Duration `mapstructure:"token_duration"`
}

type ObservabilityConfig struct {
	Tracing TracingConfig `mapstructure:"tracing"`
	Metrics MetricsConfig `mapstructure:"metrics"`
}

type TracingConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	JaegerEndpoint  string `mapstructure:"jaeger_endpoint"`
	ServiceName     string `mapstructure:"service_name"`
	SamplingRate    float64 `mapstructure:"sampling_rate"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    string `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set defaults
	setDefaults()

	// Read from environment variables
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, use defaults and env vars
		fmt.Printf("Config file not found, using defaults and environment variables: %v\n", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", "5s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.prefork", false)

	// Database defaults
	viper.SetDefault("database.url", "postgres://localhost:5432/medika?sslmode=disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "30m")

	// Redis defaults
	viper.SetDefault("redis.url", "redis://localhost:6379")
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.min_idle_conns", 2)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")

	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "change-this-in-production")
	viper.SetDefault("auth.token_duration", "24h")

	// Observability defaults
	viper.SetDefault("observability.tracing.enabled", true)
	viper.SetDefault("observability.tracing.jaeger_endpoint", "http://localhost:14268/api/traces")
	viper.SetDefault("observability.tracing.service_name", "medika-api")
	viper.SetDefault("observability.tracing.sampling_rate", 1.0)
	viper.SetDefault("observability.metrics.enabled", true)
	viper.SetDefault("observability.metrics.port", "9090")
	viper.SetDefault("observability.metrics.path", "/metrics")

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
}
