package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"messenger/internal/model"
)

// ChatService управляет чатами и сообщениями
type ChatService struct {
	chats      map[string]*model.ChatInfo
	messages   map[string][]*model.MessageData // chat_id -> messages
	userChats  map[string][]string             // user_id -> chat_ids
	mu         sync.RWMutex
	
	msgIDCounter int
}

// NewChatService создает сервис чатов
func NewChatService() *ChatService {
	return &ChatService{
		chats:      make(map[string]*model.ChatInfo),
		messages:   make(map[string][]*model.MessageData),
		userChats:  make(map[string][]string),
		msgIDCounter: 0,
	}
}

// CreateChatByLink создает чат по ссылке
func (s *ChatService) CreateChatByLink(ctx context.Context, userID string, req *model.CreateChatByLinkRequest) (*model.ChatInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// В реальном приложении здесь была бы проверка ссылки
	chatID := generateChatID()
	chat := &model.ChatInfo{
		ID:        chatID,
		Type:      "private",
		Name:      "Chat via link",
		CreatedAt: time.Now(),
		Link:      req.Link,
		Members:   []string{userID},
	}
	
	s.chats[chatID] = chat
	s.userChats[userID] = append(s.userChats[userID], chatID)
	
	return chat, nil
}

// CreateChatWithUser создает чат с другим пользователем
func (s *ChatService) CreateChatWithUser(ctx context.Context, userID string, req *model.CreateChatWithUserRequest) (*model.ChatInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Проверяем существует ли уже чат с этим пользователем
	for _, chatID := range s.userChats[userID] {
		chat := s.chats[chatID]
		if chat.Type == "private" && len(chat.Members) == 2 {
			for _, member := range chat.Members {
				if member == req.UserID {
					return chat, nil // Чат уже существует
				}
			}
		}
	}
	
	chatID := generateChatID()
	chat := &model.ChatInfo{
		ID:        chatID,
		Type:      "private",
		Name:      req.Name,
		CreatedAt: time.Now(),
		Members:   []string{userID, req.UserID},
	}
	
	s.chats[chatID] = chat
	s.userChats[userID] = append(s.userChats[userID], chatID)
	s.userChats[req.UserID] = append(s.userChats[req.UserID], chatID)
	
	return chat, nil
}

// CreateGroupChat создает групповой чат
func (s *ChatService) CreateGroupChat(ctx context.Context, userID string, req *model.CreateGroupChatRequest) (*model.ChatInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	chatID := generateChatID()
	members := append([]string{userID}, req.MemberIDs...)
	
	chat := &model.ChatInfo{
		ID:          chatID,
		Type:        "group",
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		CreatedAt:   time.Now(),
		Members:     members,
		AdminID:     userID,
	}
	
	s.chats[chatID] = chat
	for _, member := range members {
		s.userChats[member] = append(s.userChats[member], chatID)
	}
	
	return chat, nil
}

// CreateChannel создает канал
func (s *ChatService) CreateChannel(ctx context.Context, userID string, req *model.CreateChannelRequest) (*model.ChatInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	chatID := generateChatID()
	chat := &model.ChatInfo{
		ID:          chatID,
		Type:        "channel",
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		CreatedAt:   time.Now(),
		Members:     []string{userID},
		AdminID:     userID,
		Link:        generateInviteLink(),
	}
	
	s.chats[chatID] = chat
	s.userChats[userID] = append(s.userChats[userID], chatID)
	
	return chat, nil
}

// EditGroupChat редактирует групповой чат
func (s *ChatService) EditGroupChat(ctx context.Context, userID string, req *model.EditGroupChatRequest) (*model.ChatInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	chat, exists := s.chats[req.ChatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	
	if chat.AdminID != userID {
		return nil, ErrNotAdmin
	}
	
	if req.Name != "" {
		chat.Name = req.Name
	}
	if req.Description != "" {
		chat.Description = req.Description
	}
	if req.Avatar != "" {
		chat.Avatar = req.Avatar
	}
	
	return chat, nil
}

// EditChannel редактирует канал
func (s *ChatService) EditChannel(ctx context.Context, userID string, req *model.EditChannelRequest) (*model.ChatInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	chat, exists := s.chats[req.ChatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	
	if chat.AdminID != userID {
		return nil, ErrNotAdmin
	}
	
	if req.Name != "" {
		chat.Name = req.Name
	}
	if req.Description != "" {
		chat.Description = req.Description
	}
	if req.Avatar != "" {
		chat.Avatar = req.Avatar
	}
	
	return chat, nil
}

// GetChats получает список чатов пользователя
func (s *ChatService) GetChats(ctx context.Context, userID string, req *model.GetChatsRequest) ([]*model.ChatInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	chatIDs := s.userChats[userID]
	var chats []*model.ChatInfo
	
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	
	offset := req.Offset
	end := offset + limit
	if end > len(chatIDs) {
		end = len(chatIDs)
	}
	
	for i := offset; i < end; i++ {
		if chat, ok := s.chats[chatIDs[i]]; ok {
			chats = append(chats, chat)
		}
	}
	
	return chats, nil
}

// GetChatInfo получает информацию о чате
func (s *ChatService) GetChatInfo(ctx context.Context, userID string, req *model.GetChatInfoRequest) (*model.ChatInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	chat, exists := s.chats[req.ChatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	
	// Проверяем доступ
	isMember := false
	for _, member := range chat.Members {
		if member == userID {
			isMember = true
			break
		}
	}
	
	if !isMember && chat.Type != "channel" {
		return nil, ErrAccessDenied
	}
	
	return chat, nil
}

// SendMessage отправляет сообщение
func (s *ChatService) SendMessage(ctx context.Context, userID string, req *model.SendMessageRequest) (*model.MessageData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	chat, exists := s.chats[req.ChatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	
	// Проверяем доступ
	isMember := false
	for _, member := range chat.Members {
		if member == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, ErrAccessDenied
	}
	
	s.msgIDCounter++
	msgID := generateMessageID(s.msgIDCounter)
	
	msg := &model.MessageData{
		ID:        msgID,
		ChatID:    req.ChatID,
		SenderID:  userID,
		Content:   req.Content,
		MessageType: req.MessageType,
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
	}
	
	s.messages[req.ChatID] = append(s.messages[req.ChatID], msg)
	
	return msg, nil
}

// EditMessage редактирует сообщение
func (s *ChatService) EditMessage(ctx context.Context, userID string, req *model.EditMessageRequest) (*model.MessageData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Ищем сообщение
	var foundMsg *model.MessageData
	
	for _, msgs := range s.messages {
		for _, msg := range msgs {
			if msg.ID == req.MessageID {
				foundMsg = msg
				break
			}
		}
		if foundMsg != nil {
			break
		}
	}
	
	if foundMsg == nil {
		return nil, ErrMessageNotFound
	}
	
	// Проверяем что пользователь автор сообщения
	if foundMsg.SenderID != userID {
		return nil, ErrNotMessageOwner
	}
	
	foundMsg.Content = req.Content
	now := time.Now()
	foundMsg.EditedAt = &now
	
	return foundMsg, nil
}

// GetMessages получает сообщения чата
func (s *ChatService) GetMessages(ctx context.Context, userID string, req *model.GetMessagesRequest) ([]*model.MessageData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	chat, exists := s.chats[req.ChatID]
	if !exists {
		return nil, ErrChatNotFound
	}
	
	// Проверяем доступ
	isMember := false
	for _, member := range chat.Members {
		if member == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, ErrAccessDenied
	}
	
	messages := s.messages[req.ChatID]
	if len(messages) == 0 {
		return []*model.MessageData{}, nil
	}
	
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	
	// Применяем пагинацию
	start := 0
	end := len(messages)
	
	if req.AfterID != "" {
		for i, msg := range messages {
			if msg.ID == req.AfterID {
				start = i + 1
				break
			}
		}
	}
	
	if req.BeforeID != "" {
		for i, msg := range messages {
			if msg.ID == req.BeforeID {
				end = i
				break
			}
		}
	}
	
	if start > len(messages) {
		start = len(messages)
	}
	if end > len(messages) {
		end = len(messages)
	}
	
	// Ограничиваем количество
	if end-start > limit {
		end = start + limit
	}
	
	result := make([]*model.MessageData, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, messages[i])
	}
	
	return result, nil
}

// DeleteMessage удаляет сообщение
func (s *ChatService) DeleteMessage(ctx context.Context, userID string, req *model.DeleteMessageRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Ищем и удаляем сообщение
	for cid, msgs := range s.messages {
		for i, msg := range msgs {
			if msg.ID == req.MessageID {
				if msg.SenderID != userID {
					return ErrNotMessageOwner
				}
				s.messages[cid] = append(msgs[:i], msgs[i+1:]...)
				return nil
			}
		}
	}
	
	return ErrMessageNotFound
}

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrMessageNotFound   = errors.New("message not found")
	ErrAccessDenied      = errors.New("access denied")
	ErrNotAdmin          = errors.New("not admin")
	ErrNotMessageOwner   = errors.New("not message owner")
)

func generateChatID() string {
	return generateID()
}

func generateMessageID(counter int) string {
	return generateID()
}

func generateInviteLink() string {
	return "https://messenger.app/join/" + generateID()
}
