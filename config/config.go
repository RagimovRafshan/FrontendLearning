package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	Server       ServerConfig
	ScyllaDB     ScyllaDBConfig
	Redis        RedisConfig
	MinIO        MinIOConfig
	ElasticSearch ElasticSearchConfig
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Port     int
	CertFile string
	KeyFile  string
}

// ScyllaDBConfig конфигурация ScyllaDB
type ScyllaDBConfig struct {
	Hosts    []string
	Port     int
	Keyspace string
	Username string
	Password string
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// MinIOConfig конфигурация MinIO S3
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// ElasticSearchConfig конфигурация ElasticSearch
type ElasticSearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:     getEnvInt("SERVER_PORT", 8443),
			CertFile: getEnv("SERVER_CERT_FILE", "cert.pem"),
			KeyFile:  getEnv("SERVER_KEY_FILE", "key.pem"),
		},
		ScyllaDB: ScyllaDBConfig{
			Hosts:    getEnvSlice("SCYLLADB_HOSTS", "localhost"),
			Port:     getEnvInt("SCYLLADB_PORT", 9042),
			Keyspace: getEnv("SCYLLADB_KEYSPACE", "messenger"),
			Username: getEnv("SCYLLADB_USERNAME", ""),
			Password: getEnv("SCYLLADB_PASSWORD", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		MinIO: MinIOConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretAccessKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:          getEnvBool("MINIO_USE_SSL", false),
			BucketName:      getEnv("MINIO_BUCKET", "messenger"),
		},
		ElasticSearch: ElasticSearchConfig{
			Addresses: getEnvSlice("ELASTICSEARCH_ADDRESSES", "http://localhost:9200"),
			Username:  getEnv("ELASTICSEARCH_USERNAME", ""),
			Password:  getEnv("ELASTICSEARCH_PASSWORD", ""),
		},
	}

	return cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if len(c.ScyllaDB.Hosts) == 0 {
		return fmt.Errorf("scylladb hosts cannot be empty")
	}

	if c.Redis.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}

	if c.MinIO.Endpoint == "" {
		return fmt.Errorf("minio endpoint cannot be empty")
	}

	if len(c.ElasticSearch.Addresses) == 0 {
		return fmt.Errorf("elasticsearch addresses cannot be empty")
	}

	return nil
}

// GetRedisAddr возвращает адрес Redis
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetScyllaDBAddrs возвращает адреса ScyllaDB
func (c *ScyllaDBConfig) GetAddrs() []string {
	addrs := make([]string, len(c.Hosts))
	for i, host := range c.Hosts {
		addrs[i] = fmt.Sprintf("%s:%d", host, c.Port)
	}
	return addrs
}

// Вспомогательные функции

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvSlice(key, defaultValue string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return splitString(value, ",")
	}
	return splitString(defaultValue, ",")
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i = start - 1
		}
	}
	result = append(result, s[start:])
	return result
}
