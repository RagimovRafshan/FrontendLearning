package model

import "time"

// User модель пользователя
type User struct {
	ID        string    `json:"id" msgpack:"id"`
	Login     string    `json:"login" msgpack:"login"`
	Password  string    `json:"-" msgpack:"-"`
	CreatedAt time.Time `json:"created_at" msgpack:"created_at"`
	UpdatedAt time.Time `json:"updated_at" msgpack:"updated_at"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Login    string `msgpack:"login"`
	Password string `msgpack:"password"`
}

// LoginResponse ответ при входе
type LoginResponse struct {
	AccessToken  string    `msgpack:"access_token"`
	RefreshToken string    `msgpack:"refresh_token"`
	ExpiresAt    time.Time `msgpack:"expires_at"`
	User         *User     `msgpack:"user,omitempty"`
}

// RefreshTokenRequest запрос на обновление токена
type RefreshTokenRequest struct {
	RefreshToken string `msgpack:"refresh_token"`
}

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Login    string `msgpack:"login"`
	Password string `msgpack:"password"`
}

// Message модель сообщения
type Message struct {
	ID        string    `json:"id" msgpack:"id"`
	FromID    string    `json:"from_id" msgpack:"from_id"`
	ToID      string    `json:"to_id" msgpack:"to_id"`
	Content   string    `json:"content" msgpack:"content"`
	CreatedAt time.Time `json:"created_at" msgpack:"created_at"`
	Read      bool      `json:"read" msgpack:"read"`
}

// Chat модель чата
type Chat struct {
	ID        string    `json:"id" msgpack:"id"`
	Participants []string `json:"participants" msgpack:"participants"`
	CreatedAt time.Time `json:"created_at" msgpack:"created_at"`
	LastMessageID string `json:"last_message_id" msgpack:"last_message_id"`
}

// APIError ошибка API
type APIError struct {
	Code    int    `msgpack:"code"`
	Message string `msgpack:"message"`
	Field   string `msgpack:"field,omitempty"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error APIError `msgpack:"error"`
}

// SuccessResponse успешный ответ
type SuccessResponse struct {
	Data interface{} `msgpack:"data"`
}
