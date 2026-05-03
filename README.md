# Backend мессенджера на Go

## Структура проекта

```
messenger/
├── config/
│   ├── config.go      # Структуры конфигурации
│   └── loader.go      # Загрузка и валидация конфигурации
├── internal/
│   └── storage/
│       ├── scylladb.go      # Клиент ScyllaDB
│       ├── redis.go         # Клиент Redis
│       ├── minio.go         # Клиент MinIO S3
│       └── elasticsearch.go # Клиент ElasticSearch
├── go.mod
└── README.md
```

## Конфигурация

Все сервисы настроены на подключение к localhost по умолчанию.

### Переменные окружения

#### ScyllaDB
- SCYLLA_HOSTS=localhost:9042
- SCYLLA_KEYSPACE=messenger
- SCYLLA_TIMEOUT=5s
- SCYLLA_CONNECT_TIMEOUT=10s

#### Redis
- REDIS_ADDR=localhost:6379
- REDIS_DB=0
- REDIS_POOL_SIZE=100

#### MinIO S3
- MINIO_ENDPOINT=localhost:9000
- MINIO_ACCESS_KEY_ID=minioadmin
- MINIO_SECRET_ACCESS_KEY=minioadmin
- MINIO_BUCKET_NAME=messenger-files

#### ElasticSearch
- ELASTICSEARCH_ADDRESSES=http://localhost:9200
- ELASTICSEARCH_INDEX_PREFIX=messenger

## Требования

- Go 1.19+
- ScyllaDB (localhost:9042)
- Redis (localhost:6379)
- MinIO (localhost:9000)
- ElasticSearch (localhost:9200)
