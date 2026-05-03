package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v9"
)

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}

	// Загружаем базовые значения из env
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Парсим списки хостов вручную, так как env не всегда корректно обрабатывает срезы
	cfg.ScyllaDB.Hosts = parseStringSlice(os.Getenv("SCYLLA_HOSTS"), []string{"localhost:9042"})
	cfg.ElasticSearch.Addresses = parseStringSlice(os.Getenv("ELASTICSEARCH_ADDRESSES"), []string{"http://localhost:9200"})

	// Применяем значения по умолчанию для временных интервалов, если они не были установлены
	setDefaultDurations(cfg)

	return cfg, nil
}

// parseStringSlice парсит строку в срез строк (разделитель - запятая)
func parseStringSlice(value string, defaultVal []string) []string {
	if value == "" {
		return defaultVal
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultVal
	}

	return result
}

// setDefaultDurations устанавливает значения по умолчанию для временных интервалов
func setDefaultDurations(cfg *Config) {
	if cfg.ScyllaDB.Timeout == 0 {
		cfg.ScyllaDB.Timeout = 5 * time.Second
	}
	if cfg.ScyllaDB.ConnectTimeout == 0 {
		cfg.ScyllaDB.ConnectTimeout = 10 * time.Second
	}
	if cfg.Redis.Timeout == 0 {
		cfg.Redis.Timeout = 5 * time.Second
	}
	if cfg.Redis.DialTimeout == 0 {
		cfg.Redis.DialTimeout = 5 * time.Second
	}
	if cfg.Redis.ReadTimeout == 0 {
		cfg.Redis.ReadTimeout = 3 * time.Second
	}
	if cfg.Redis.WriteTimeout == 0 {
		cfg.Redis.WriteTimeout = 3 * time.Second
	}
	if cfg.ElasticSearch.Timeout == 0 {
		cfg.ElasticSearch.Timeout = 10 * time.Second
	}
}

// MustLoad загружает конфигурацию и паникует при ошибке
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

// validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверка ScyllaDB
	if len(c.ScyllaDB.Hosts) == 0 {
		return ErrNoScyllaHosts
	}
	if c.ScyllaDB.Keyspace == "" {
		return ErrNoKeyspace
	}

	// Проверка Redis
	if c.Redis.Addr == "" {
		return ErrNoRedisAddr
	}

	// Проверка MinIO
	if c.MinIO.Endpoint == "" {
		return ErrNoMinIOEndpoint
	}
	if c.MinIO.AccessKeyID == "" || c.MinIO.SecretAccessKey == "" {
		return ErrNoMinIOCredentials
	}

	// Проверка ElasticSearch
	if len(c.ElasticSearch.Addresses) == 0 {
		return ErrNoElasticsearchAddresses
	}

	return nil
}

// Ошибки конфигурации
var (
	ErrNoScyllaHosts         = &configError{"no ScyllaDB hosts specified"}
	ErrNoKeyspace            = &configError{"no keyspace specified for ScyllaDB"}
	ErrNoRedisAddr           = &configError{"no Redis address specified"}
	ErrNoMinIOEndpoint       = &configError{"no MinIO endpoint specified"}
	ErrNoMinIOCredentials    = &configError{"no MinIO credentials specified"}
	ErrNoElasticsearchAddresses = &configError{"no ElasticSearch addresses specified"}
)

type configError struct {
	message string
}

func (e *configError) Error() string {
	return "config error: " + e.message
}
