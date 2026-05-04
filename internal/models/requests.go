package models

import "time"

// Auth requests
type RegisterRequest struct {
	Login    string `msgpack:"login"`
	Password string `msgpack:"password"`
	Email    string `msgpack:"email,omitempty"`
}

type LoginRequest struct {
	Login    string `msgpack:"login"`
	Password string `msgpack:"password"`
}

type ForgotPasswordRequest struct {
	Email string `msgpack:"email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `msgpack:"refresh_token"`
}

// User requests
type FindUserRequest struct {
	Query string `msgpack:"query"`
	Limit int    `msgpack:"limit,omitempty"`
}

type EditProfileRequest struct {
	DisplayName string `msgpack:"display_name,omitempty"`
	Avatar      string `msgpack:"avatar,omitempty"`
	Bio         string `msgpack:"bio,omitempty"`
}

// Chat requests
type CreateChatByLinkRequest struct {
	Name        string `msgpack:"name,omitempty"`
	Description string `msgpack:"description,omitempty"`
}

type CreateChatWithUserRequest struct {
	UserID string `msgpack:"user_id"`
	Name   string `msgpack:"name,omitempty"`
}

type CreateGroupChatRequest struct {
	Name        string   `msgpack:"name"`
	Description string   `msgpack:"description,omitempty"`
	MemberIDs   []string `msgpack:"member_ids"`
}

type CreateChannelRequest struct {
	Name        string `msgpack:"name"`
	Description string `msgpack:"description,omitempty"`
	IsPublic    bool   `msgpack:"is_public"`
}

type EditGroupChatRequest struct {
	ChatID      string `msgpack:"chat_id"`
	Name        string `msgpack:"name,omitempty"`
	Description string `msgpack:"description,omitempty"`
}

type EditChannelRequest struct {
	ChannelID   string `msgpack:"channel_id"`
	Name        string `msgpack:"name,omitempty"`
	Description string `msgpack:"description,omitempty"`
}

// Message requests
type SendMessageRequest struct {
	ChatID      string `msgpack:"chat_id"`
	Content     string `msgpack:"content"`
	MessageType string `msgpack:"message_type,omitempty"` // text, image, file, etc.
	MediaURL    string `msgpack:"media_url,omitempty"`
}

type GetMessageRequest struct {
	MessageID string `msgpack:"message_id"`
}

type EditMessageRequest struct {
	MessageID string `msgpack:"message_id"`
	Content   string `msgpack:"content"`
}

// Common response
type Response struct {
	Success bool        `msgpack:"success"`
	Data    interface{} `msgpack:"data,omitempty"`
	Error   string      `msgpack:"error,omitempty"`
}

// Auth response
type AuthResponse struct {
	AccessToken  string `msgpack:"access_token"`
	RefreshToken string `msgpack:"refresh_token"`
	UserID       string `msgpack:"user_id"`
	Login        string `msgpack:"login"`
	ExpiresAt    int64  `msgpack:"expires_at"`
}

// User response
type UserResponse struct {
	UserID      string `msgpack:"user_id"`
	Login       string `msgpack:"login"`
	DisplayName string `msgpack:"display_name,omitempty"`
	Avatar      string `msgpack:"avatar,omitempty"`
	Bio         string `msgpack:"bio,omitempty"`
	IsOnline    bool   `msgpack:"is_online"`
	LastSeen    int64  `msgpack:"last_seen,omitempty"`
}

// Chat response
type ChatResponse struct {
	ChatID      string   `msgpack:"chat_id"`
	Type        string   `msgpack:"type"` // private, group, channel
	Name        string   `msgpack:"name"`
	Description string   `msgpack:"description,omitempty"`
	MemberIDs   []string `msgpack:"member_ids,omitempty"`
	OwnerID     string   `msgpack:"owner_id,omitempty"`
	CreatedAt   int64    `msgpack:"created_at"`
	InviteLink  string   `msgpack:"invite_link,omitempty"`
}

// Message response
type MessageResponse struct {
	MessageID   string `msgpack:"message_id"`
	ChatID      string `msgpack:"chat_id"`
	SenderID    string `msgpack:"sender_id"`
	SenderName  string `msgpack:"sender_name,omitempty"`
	Content     string `msgpack:"content"`
	MessageType string `msgpack:"message_type"`
	MediaURL    string `msgpack:"media_url,omitempty"`
	CreatedAt   int64  `msgpack:"created_at"`
	EditedAt    int64  `msgpack:"edited_at,omitempty"`
}

// WebSocket messages
type WSMessage struct {
	Type      string      `msgpack:"type"`
	Payload   interface{} `msgpack:"payload,omitempty"`
	RequestID string      `msgpack:"request_id,omitempty"`
}

type WSNewMessage struct {
	MessageID   string `msgpack:"message_id"`
	ChatID      string `msgpack:"chat_id"`
	SenderID    string `msgpack:"sender_id"`
	Content     string `msgpack:"content"`
	MessageType string `msgpack:"message_type"`
	CreatedAt   int64  `msgpack:"created_at"`
}

type WSUserStatus struct {
	UserID   string `msgpack:"user_id"`
	IsOnline bool   `msgpack:"is_online"`
	LastSeen int64  `msgpack:"last_seen,omitempty"`
}

type WSSubscribeRequest struct {
	ChatIDs []string `msgpack:"chat_ids,omitempty"`
	UserIDs []string `msgpack:"user_ids,omitempty"`
}

type WSConnectionInfo struct {
	UserID          string   `msgpack:"user_id"`
	SubscribedChats []string `msgpack:"subscribed_chats"`
	SubscribedUsers []string `msgpack:"subscribed_users"`
	ConnectedAt     int64    `msgpack:"connected_at"`
}

// Pagination
type PaginatedRequest struct {
	Limit  int `msgpack:"limit,omitempty"`
	Offset int `msgpack:"offset,omitempty"`
}

type PaginatedResponse struct {
	Items   interface{} `msgpack:"items"`
	Total   int         `msgpack:"total"`
	Limit   int         `msgpack:"limit"`
	Offset  int         `msgpack:"offset"`
	HasMore bool        `msgpack:"has_more"`
}

// Timestamps
func Now() int64 {
	return time.Now().UnixNano() / 1e6
}
