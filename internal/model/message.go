package model

import (
	"time"

	"github.com/google/uuid"
)

// MessageType определяет тип сообщения в протоколе
type MessageType string

const (
	// Авторизация
	MsgTypeRegister       MessageType = "register"
	MsgTypeLogin          MessageType = "login"
	MsgTypeForgotPassword MessageType = "forgot_password"
	MsgTypeRefreshToken   MessageType = "refresh_token"

	// Пользователи
	MsgTypeSearchUser    MessageType = "search_user"
	MsgTypeUpdateProfile MessageType = "update_profile"
	MsgTypeGetProfile    MessageType = "get_profile"

	// Чаты
	MsgTypeCreateChatByLink   MessageType = "create_chat_by_link"
	MsgTypeCreateChatWithUser MessageType = "create_chat_with_user"
	MsgTypeCreateGroupChat    MessageType = "create_group_chat"
	MsgTypeCreateChannel      MessageType = "create_channel"
	MsgTypeEditGroupChat      MessageType = "edit_group_chat"
	MsgTypeEditChannel        MessageType = "edit_channel"
	MsgTypeGetChats           MessageType = "get_chats"
	MsgTypeGetChatInfo        MessageType = "get_chat_info"

	// Сообщения
	MsgTypeSendMessage  MessageType = "send_message"
	MsgTypeEditMessage  MessageType = "edit_message"
	MsgTypeGetMessages  MessageType = "get_messages"
	MsgTypeDeleteMessage MessageType = "delete_message"

	// Статусы и уведомления
	MsgTypeGetOnlineStatus MessageType = "get_online_status"
	MsgTypeSubscribeEvents MessageType = "subscribe_events"

	// Ответы сервера
	MsgTypeResponse      MessageType = "response"
	MsgTypeError         MessageType = "error"
	MsgTypeNewMessage    MessageType = "new_message"
	MsgTypeUserOnline    MessageType = "user_online"
	MsgTypeUserOffline   MessageType = "user_offline"
	MsgTypeTypingStatus  MessageType = "typing_status"
)

// Request - базовая структура запроса от клиента
type Request struct {
	Type      MessageType     `msgpack:"type"`
	RequestID string          `msgpack:"request_id"`
	Payload   interface{}     `msgpack:"payload,omitempty"`
}

// Response - базовая структура ответа сервера
type Response struct {
	Type      MessageType     `msgpack:"type"`
	RequestID string          `msgpack:"request_id"`
	Success   bool            `msgpack:"success"`
	Data      interface{}     `msgpack:"data,omitempty"`
	Error     *ErrorData      `msgpack:"error,omitempty"`
}

// ErrorData - данные об ошибке
type ErrorData struct {
	Code    int    `msgpack:"code"`
	Message string `msgpack:"message"`
}

// ========== Авторизация ==========

// RegisterRequest - регистрация
type RegisterRequest struct {
	Username string `msgpack:"username"`
	Password string `msgpack:"password"`
	Email    string `msgpack:"email,omitempty"`
}

// LoginRequest - вход
type LoginRequest struct {
	Username string `msgpack:"username"`
	Password string `msgpack:"password"`
}

// ForgotPasswordRequest - восстановление пароля
type ForgotPasswordRequest struct {
	EmailOrUsername string `msgpack:"email_or_username"`
}

// RefreshTokenRequest - обновление токена
type RefreshTokenRequest struct {
	RefreshToken string `msgpack:"refresh_token"`
}

// AuthResponse - ответ авторизации
type AuthResponse struct {
	User         User   `msgpack:"user"`
	AccessToken  string `msgpack:"access_token"`
	RefreshToken string `msgpack:"refresh_token"`
}

// ========== Пользователи ==========

// SearchUserRequest - поиск пользователя
type SearchUserRequest struct {
	Query string `msgpack:"query"`
	Limit int    `msgpack:"limit,omitempty"`
}

// UpdateProfileRequest - обновление профиля
type UpdateProfileRequest struct {
	Username    string `msgpack:"username,omitempty"`
	DisplayName string `msgpack:"display_name,omitempty"`
	Avatar      string `msgpack:"avatar,omitempty"`
	Bio         string `msgpack:"bio,omitempty"`
}

// GetProfileRequest - получение профиля
type GetProfileRequest struct {
	UserID string `msgpack:"user_id,omitempty"`
}

// ========== Чаты ==========

// CreateChatByLinkRequest - создание чата по ссылке
type CreateChatByLinkRequest struct {
	Link string `msgpack:"link"`
}

// CreateChatWithUserRequest - создание чата с пользователем
type CreateChatWithUserRequest struct {
	UserID string `msgpack:"user_id"`
	Name   string `msgpack:"name,omitempty"`
}

// CreateGroupChatRequest - создание группового чата
type CreateGroupChatRequest struct {
	Name        string   `msgpack:"name"`
	Description string   `msgpack:"description,omitempty"`
	MemberIDs   []string `msgpack:"member_ids"`
	Avatar      string   `msgpack:"avatar,omitempty"`
}

// CreateChannelRequest - создание канала
type CreateChannelRequest struct {
	Name        string `msgpack:"name"`
	Description string `msgpack:"description,omitempty"`
	IsPublic    bool   `msgpack:"is_public"`
	Avatar      string `msgpack:"avatar,omitempty"`
}

// EditGroupChatRequest - редактирование группы
type EditGroupChatRequest struct {
	ChatID      string `msgpack:"chat_id"`
	Name        string `msgpack:"name,omitempty"`
	Description string `msgpack:"description,omitempty"`
	Avatar      string `msgpack:"avatar,omitempty"`
}

// EditChannelRequest - редактирование канала
type EditChannelRequest struct {
	ChatID      string `msgpack:"chat_id"`
	Name        string `msgpack:"name,omitempty"`
	Description string `msgpack:"description,omitempty"`
	Avatar      string `msgpack:"avatar,omitempty"`
}

// GetChatsRequest - получение списка чатов
type GetChatsRequest struct {
	Limit  int    `msgpack:"limit,omitempty"`
	Offset int    `msgpack:"offset,omitempty"`
}

// GetChatInfoRequest - получение информации о чате
type GetChatInfoRequest struct {
	ChatID string `msgpack:"chat_id"`
}

// ChatInfo - информация о чате
type ChatInfo struct {
	ID          string    `msgpack:"id"`
	Type        string    `msgpack:"type"` // "private", "group", "channel"
	Name        string    `msgpack:"name"`
	Description string    `msgpack:"description,omitempty"`
	Avatar      string    `msgpack:"avatar,omitempty"`
	CreatedAt   time.Time `msgpack:"created_at"`
	Members     []string  `msgpack:"members,omitempty"`
	AdminID     string    `msgpack:"admin_id,omitempty"`
	Link        string    `msgpack:"link,omitempty"`
}

// ========== Сообщения ==========

// SendMessageRequest - отправка сообщения
type SendMessageRequest struct {
	ChatID    string                 `msgpack:"chat_id"`
	Content   string                 `msgpack:"content"`
	MessageType string               `msgpack:"message_type,omitempty"` // "text", "image", "file", etc.
	Metadata  map[string]interface{} `msgpack:"metadata,omitempty"`
}

// EditMessageRequest - редактирование сообщения
type EditMessageRequest struct {
	MessageID string `msgpack:"message_id"`
	Content   string `msgpack:"content"`
}

// GetMessagesRequest - получение сообщений
type GetMessagesRequest struct {
	ChatID   string `msgpack:"chat_id"`
	Limit    int    `msgpack:"limit,omitempty"`
	Offset   int    `msgpack:"offset,omitempty"`
	BeforeID string `msgpack:"before_id,omitempty"`
	AfterID  string `msgpack:"after_id,omitempty"`
}

// DeleteMessageRequest - удаление сообщения
type DeleteMessageRequest struct {
	MessageID string `msgpack:"message_id"`
}

// MessageData - данные сообщения
type MessageData struct {
	ID        string                 `msgpack:"id"`
	ChatID    string                 `msgpack:"chat_id"`
	SenderID  string                 `msgpack:"sender_id"`
	Sender    *User                  `msgpack:"sender,omitempty"`
	Content   string                 `msgpack:"content"`
	MessageType string               `msgpack:"message_type"`
	Metadata  map[string]interface{} `msgpack:"metadata,omitempty"`
	CreatedAt time.Time              `msgpack:"created_at"`
	EditedAt  *time.Time             `msgpack:"edited_at,omitempty"`
}

// ========== Статусы ==========

// GetOnlineStatusRequest - проверка статуса online
type GetOnlineStatusRequest struct {
	UserIDs []string `msgpack:"user_ids"`
}

// OnlineStatusResponse - статус online пользователей
type OnlineStatusResponse struct {
	OnlineUsers map[string]bool `msgpack:"online_users"` // user_id -> is_online
}

// SubscribeEventsRequest - подписка на события
type SubscribeEventsRequest struct {
	Events []string `msgpack:"events"` // "new_message", "user_online", "typing", etc.
}

// TypingStatusRequest - статус набора текста
type TypingStatusRequest struct {
	ChatID   string `msgpack:"chat_id"`
	IsTyping bool   `msgpack:"is_typing"`
}

// ========== Пользователь ==========

// User - данные пользователя
type User struct {
	ID          string    `msgpack:"id"`
	Username    string    `msgpack:"username"`
	DisplayName string    `msgpack:"display_name,omitempty"`
	Email       string    `msgpack:"email,omitempty"`
	Avatar      string    `msgpack:"avatar,omitempty"`
	Bio         string    `msgpack:"bio,omitempty"`
	CreatedAt   time.Time `msgpack:"created_at"`
	IsOnline    bool      `msgpack:"is_online,omitempty"`
}

// NewRequest создает новый запрос с уникальным ID
func NewRequest(msgType MessageType, payload interface{}) *Request {
	return &Request{
		Type:      msgType,
		RequestID: uuid.New().String(),
		Payload:   payload,
	}
}

// NewResponse создает новый ответ
func NewResponse(reqID string, msgType MessageType, success bool, data interface{}, err *ErrorData) *Response {
	return &Response{
		Type:      msgType,
		RequestID: reqID,
		Success:   success,
		Data:      data,
		Error:     err,
	}
}

// NewErrorResponse создает ответ с ошибкой
func NewErrorResponse(reqID string, code int, message string) *Response {
	return &Response{
		Type:      MsgTypeError,
		RequestID: reqID,
		Success:   false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
		},
	}
}
