# Messenger Backend

Высоконагруженный backend мессенджера на Go с использованием ScyllaDB, Redis, MinIO S3 и ElasticSearch.

## Особенности

- **HTTPS** сервер с TLS 1.2+
- **MessagePack** формат для передачи данных (быстрее JSON)
- **Аутентификация** на основе токенов
- **Rate limiting** для защиты от DDoS
- **CORS** поддержка
- **Graceful shutdown**

## Структура проекта

```
├── cmd/server          # Точка входа
├── config              # Конфигурация
├── internal/
│   ├── handler         # HTTP обработчики
│   ├── middleware      # Middleware (auth, rate limit, cors)
│   ├── model           # Модели данных
│   ├── service         # Бизнес логика
│   └── storage         # Клиенты БД и хранилищ
├── pkg/msgpack         # MessagePack сериализация
└── .env.example        # Пример конфигурации
```

## Быстрый старт

### 1. Запуск зависимостей (Docker Compose)

```bash
docker-compose up -d
```

### 2. Генерация SSL сертификатов

```bash
mkdir -p certs
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes
```

### 3. Настройка конфигурации

```bash
cp .env.example .env
```

### 4. Установка зависимостей

```bash
go mod tidy
```

### 5. Запуск сервера

```bash
go run ./cmd/server
```

## API Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| POST | /api/v1/auth/register | Регистрация пользователя |
| POST | /api/v1/auth/login | Вход |
| POST | /api/v1/auth/refresh | Обновление токена |
| POST | /api/v1/auth/logout | Выход |
| GET | /health | Проверка здоровья |

## Формат запросов

По умолчанию используется **MessagePack**:

```bash
# Регистрация
curl -k -X POST https://localhost:8443/api/v1/auth/register \
  -H "Content-Type: application/msgpack" \
  --data-binary $(python3 -c "import msgpack; print(msgpack.packb({'login': 'user', 'password': 'pass'}).hex())" | xxd -r -p)

# Вход
curl -k -X POST https://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/msgpack" \
  --data-binary $(python3 -c "import msgpack; print(msgpack.packb({'login': 'user', 'password': 'pass'}).hex())" | xxd -r -p)
```

Также поддерживается JSON:

```bash
curl -k -X POST https://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user","password":"pass"}'
```

## Переменные окружения

### ScyllaDB
- `SCYLLA_HOSTS` - хосты (default: localhost:9042)
- `SCYLLA_KEYSPACE` - keyspace (default: messenger)

### Redis
- `REDIS_ADDR` - адрес (default: localhost:6379)
- `REDIS_PASSWORD` - пароль

### MinIO
- `MINIO_ENDPOINT` - endpoint (default: localhost:9000)
- `MINIO_ACCESS_KEY` - access key (default: minioadmin)
- `MINIO_SECRET_KEY` - secret key (default: minioadmin)

### ElasticSearch
- `ES_ADDRESSES` - адреса (default: http://localhost:9200)

### Сервер
- `SERVER_PORT` - порт HTTPS (default: 8443)
- `CERT_FILE` - путь к сертификату
- `KEY_FILE` - путь к ключу

## Архитектура аутентификации

Система использует двухтокенную схему:
- **Access Token** - короткоживущий токен для доступа к API
- **Refresh Token** - долгоживущий токен для обновления access token

Пароли хранятся в хешированном виде (bcrypt).

## Производительность

- Rate limiting: 1000 RPS по умолчанию
- Connection pooling для всех клиентов
- Async graceful shutdown
