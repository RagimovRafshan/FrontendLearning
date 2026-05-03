package config

import (
	"time"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	ScyllaDB   ScyllaDBConfig
	Redis      RedisConfig
	MinIO      MinIOConfig
	ElasticSearch ElasticSearchConfig
}

// ScyllaDBConfig конфигурация подключения к ScyllaDB
type ScyllaDBConfig struct {
	Hosts        []string      `env:"SCYLLA_HOSTS" env-default:"localhost:9042"`
	Keyspace     string        `env:"SCYLLA_KEYSPACE" env-default:"messenger"`
	Timeout      time.Duration `env:"SCYLLA_TIMEOUT" env-default:"5s"`
	ConnectTimeout time.Duration `env:"SCYLLA_CONNECT_TIMEOUT" env-default:"10s"`
	Username     string        `env:"SCYLLA_USERNAME"`
	Password     string        `env:"SCYLLA_PASSWORD"`
	Consistency  string        `env:"SCYLLA_CONSISTENCY" env-default:"LOCAL_QUORUM"`
}

// RedisConfig конфигурация подключения к Redis
type RedisConfig struct {
	Addr         string        `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password     string        `env:"REDIS_PASSWORD"`
	DB           int           `env:"REDIS_DB" env-default:"0"`
	PoolSize     int           `env:"REDIS_POOL_SIZE" env-default:"100"`
	MinIdleConns int           `env:"REDIS_MIN_IDLE_CONNS" env-default:"10"`
	Timeout      time.Duration `env:"REDIS_TIMEOUT" env-default:"5s"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" env-default:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" env-default:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" env-default:"3s"`
}

// MinIOConfig конфигурация подключения к MinIO S3
type MinIOConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID" env-default:"minioadmin"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY" env-default:"minioadmin"`
	BucketName      string `env:"MINIO_BUCKET_NAME" env-default:"messenger-files"`
	UseSSL          bool   `env:"MINIO_USE_SSL" env-default:"false"`
	Region          string `env:"MINIO_REGION" env-default:"us-east-1"`
}

// ElasticSearchConfig конфигурация подключения к ElasticSearch
type ElasticSearchConfig struct {
	Addresses    []string      `env:"ELASTICSEARCH_ADDRESSES" env-default:"http://localhost:9200"`
	Username     string        `env:"ELASTICSEARCH_USERNAME"`
	Password     string        `env:"ELASTICSEARCH_PASSWORD"`
	IndexPrefix  string        `env:"ELASTICSEARCH_INDEX_PREFIX" env-default:"messenger"`
	Sniff        bool          `env:"ELASTICSEARCH_SNIFF" env-default:"false"`
	HealthCheck  bool          `env:"ELASTICSEARCH_HEALTH_CHECK" env-default:"true"`
	Timeout      time.Duration `env:"ELASTICSEARCH_TIMEOUT" env-default:"10s"`
}
