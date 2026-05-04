package config

import "time"

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Host            string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port            string        `env:"SERVER_PORT" env-default:"8443"`
	CertFile        string        `env:"CERT_FILE" env-default:"./certs/server.crt"`
	KeyFile         string        `env:"KEY_FILE" env-default:"./certs/server.key"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" env-default:"120s"`
	MaxHeaderBytes  int           `env:"MAX_HEADER_BYTES" env-default:"1048576"`
	MaxConnsPerHost int           `env:"MAX_CONNS_PER_HOST" env-default:"10000"`
}

// AuthConfig конфигурация аутентификации
type AuthConfig struct {
	TokenExpiry    time.Duration `env:"TOKEN_EXPIRY" env-default:"24h"`
	RefreshExpiry  time.Duration `env:"REFRESH_EXPIRY" env-default:"7d"`
	PasswordMinLen int           `env:"PASSWORD_MIN_LEN" env-default:"8"`
	BcryptCost     int           `env:"BCRYPT_COST" env-default:"12"`
}

// RateLimitConfig конфигурация ограничения запросов
type RateLimitConfig struct {
	Enabled       bool          `env:"RATE_LIMIT_ENABLED" env-default:"true"`
	RequestsPerSec int          `env:"RATE_LIMIT_RPS" env-default:"1000"`
	BurstSize     int           `env:"RATE_LIMIT_BURST" env-default:"2000"`
}

// CORSConfig конфигурация CORS
type CORSConfig struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" env-default:"*"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS" env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" env-default:"Content-Type,Authorization,X-Request-ID"`
	MaxAge         int      `env:"CORS_MAX_AGE" env-default:"86400"`
}

// PoolConfig конфигурация пулов соединений
type PoolConfig struct {
	MinWorkers int `env:"POOL_MIN_WORKERS" env-default:"10"`
	MaxWorkers int `env:"POOL_MAX_WORKERS" env-default:"1000"`
	QueueSize  int `env:"POOL_QUEUE_SIZE" env-default:"10000"`
}
