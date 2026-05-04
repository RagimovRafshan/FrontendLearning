# Messenger Backend (Go)

Высоконагруженный backend мессенджера на Go с поддержкой HTTPS и MessagePack.

## Технологии

- **ScyllaDB** - основное хранилище сообщений и чатов
- **Redis** - кэширование сессий и rate limiting
- **MinIO S3** - хранение файлов и медиа
- **ElasticSearch** - поиск пользователей и сообщений
- **WebSocket** - реальное время для новых сообщений и статусов online

## API Endpoints

### Публичные (без авторизации)
| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/auth/register` | Регистрация |
| POST | `/api/v1/auth/login` | Вход |
| POST | `/api/v1/auth/forgot-password` | Восстановление пароля |
| POST | `/api/v1/auth/refresh-token` | Обновление токена |

### Защищенные (требуется Authorization: Bearer <token>)
| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/auth/logout` | Выход |
| POST | `/api/v1/users/search` | Поиск пользователей |
| PUT | `/api/v1/users/profile` | Редактирование профиля |
| GET | `/api/v1/users/profile?user_id=` | Получение профиля |
| POST | `/api/v1/chats/create-by-link` | Создать чат по ссылке |
| POST | `/api/v1/chats/create-with-user` | Создать чат с пользователем |
| POST | `/api/v1/chats/create-group` | Создать групповой чат |
| POST | `/api/v1/chats/create-channel` | Создать канал |
| PUT | `/api/v1/chats/edit-group` | Редактировать группу |
| PUT | `/api/v1/chats/edit-channel` | Редактировать канал |
| GET | `/api/v1/chats/get?chat_id=` | Получить чат |
| POST | `/api/v1/messages/send` | Отправить сообщение |
| POST | `/api/v1/messages/get` | Получить сообщение |
| PUT | `/api/v1/messages/edit` | Редактировать сообщение |
| GET | `/api/v1/messages/chat?chat_id=` | История сообщений |

### WebSocket
| Endpoint | Описание |
|----------|----------|
| GET `/ws?token=<token>` | Постоянное соединение для realtime |

## Формат данных

Все запросы/ответы используют **MessagePack** (`Content-Type: application/msgpack`).
JSON поддерживается как fallback.

## Запуск

```bash
# Копировать конфиг
cp .env.example .env

# Установить зависимости
go mod tidy

# Сгенерировать SSL сертификаты (для разработки)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# Запустить сервер
go run ./cmd/server
```

## WebSocket протокол

### Типы сообщений от клиента:
- `subscribe_chats` - подписка на чаты
- `subscribe_users` - подписка на статусы пользователей
- `unsubscribe_all` - отписаться от всего
- `heartbeat` - проверка соединения

### Типы сообщений от сервера:
- `connection_info` - информация о подключении
- `new_message` - новое сообщение в чате
- `user_status` - изменение статуса пользователя (online/offline)

## Структура проекта

```
messenger/
├── cmd/server/main.go      # Точка входа
├── config/                 # Конфигурация
├── internal/
│   ├── handlers/          # HTTP обработчики
│   ├── middleware/        # Middleware (auth, CORS, rate limit)
│   ├── models/            # Модели данных
│   ├── storage/           # Клиенты БД
│   └── websocket/         # WebSocket менеджер
└── pkg/msgpack/           # MessagePack утилиты
```

## Система авторизации

Пользователи хранятся в `Map<login, passwordHash>`. Пароли хешируются через bcrypt.
При логине выдаются Access и Refresh токены.

Для всех защищенных endpoints требуется заголовок:
```
Authorization: Bearer <access_token>
```
