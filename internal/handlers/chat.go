package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"messenger/internal/models"
	"messenger/internal/storage"
	"messenger/pkg/msgpack"
)

type ChatHandler struct {
	scylla      *storage.ScyllaDBClient
	elastic     *storage.ElasticsearchClient
	chats       sync.Map // Map[chatID]chatData
	chatMembers sync.Map // Map[chatID][]memberIDs
}

type chatData struct {
	ChatID      string
	Type        string // private, group, channel
	Name        string
	Description string
	OwnerID     string
	MemberIDs   []string
	CreatedAt   int64
	InviteLink  string
}

func NewChatHandler(scylla *storage.ScyllaDBClient, elastic *storage.ElasticsearchClient) *ChatHandler {
	return &ChatHandler{
		scylla:  scylla,
		elastic: elastic,
	}
}

func (h *ChatHandler) CreateChatByLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.CreateChatByLinkRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	chatID := uuid.New().String()
	inviteLink := generateInviteLink()

	chat := &chatData{
		ChatID:      chatID,
		Type:        "private",
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		MemberIDs:   []string{userID},
		CreatedAt:   time.Now().UnixMilli(),
		InviteLink:  inviteLink,
	}

	h.chats.Store(chatID, chat)
	h.chatMembers.Store(chatID, chat.MemberIDs)

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:      chat.ChatID,
			Type:        chat.Type,
			Name:        chat.Name,
			Description: chat.Description,
			MemberIDs:   chat.MemberIDs,
			OwnerID:     chat.OwnerID,
			CreatedAt:   chat.CreatedAt,
			InviteLink:  chat.InviteLink,
		},
	})
}

func (h *ChatHandler) CreateChatWithUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.CreateChatWithUserRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.UserID == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "user_id is required",
		})
		return
	}

	if req.UserID == userID {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Cannot create chat with yourself",
		})
		return
	}

	chatID := uuid.New().String()
	chatName := req.Name
	if chatName == "" {
		chatName = "Private Chat"
	}

	chat := &chatData{
		ChatID:      chatID,
		Type:        "private",
		Name:        chatName,
		OwnerID:     userID,
		MemberIDs:   []string{userID, req.UserID},
		CreatedAt:   time.Now().UnixMilli(),
	}

	h.chats.Store(chatID, chat)
	h.chatMembers.Store(chatID, chat.MemberIDs)

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:    chat.ChatID,
			Type:      chat.Type,
			Name:      chat.Name,
			MemberIDs: chat.MemberIDs,
			OwnerID:   chat.OwnerID,
			CreatedAt: chat.CreatedAt,
		},
	})
}

func (h *ChatHandler) CreateGroupChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.CreateGroupChatRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Name == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "name is required",
		})
		return
	}

	// Add owner to members if not present
	memberIDs := make([]string, 0)
	ownerFound := false
	for _, memberID := range req.MemberIDs {
		if memberID == userID {
			ownerFound = true
		}
		memberIDs = append(memberIDs, memberID)
	}
	if !ownerFound {
		memberIDs = append(memberIDs, userID)
	}

	chatID := uuid.New().String()

	chat := &chatData{
		ChatID:      chatID,
		Type:        "group",
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		MemberIDs:   memberIDs,
		CreatedAt:   time.Now().UnixMilli(),
	}

	h.chats.Store(chatID, chat)
	h.chatMembers.Store(chatID, memberIDs)

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:      chat.ChatID,
			Type:        chat.Type,
			Name:        chat.Name,
			Description: chat.Description,
			MemberIDs:   chat.MemberIDs,
			OwnerID:     chat.OwnerID,
			CreatedAt:   chat.CreatedAt,
		},
	})
}

func (h *ChatHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.CreateChannelRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Name == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "name is required",
		})
		return
	}

	chatID := uuid.New().String()
	chatType := "channel"
	if req.IsPublic {
		chatType = "public_channel"
	}

	chat := &chatData{
		ChatID:      chatID,
		Type:        chatType,
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		MemberIDs:   []string{userID},
		CreatedAt:   time.Now().UnixMilli(),
	}

	h.chats.Store(chatID, chat)
	h.chatMembers.Store(chatID, chat.MemberIDs)

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:      chat.ChatID,
			Type:        chat.Type,
			Name:        chat.Name,
			Description: chat.Description,
			MemberIDs:   chat.MemberIDs,
			OwnerID:     chat.OwnerID,
			CreatedAt:   chat.CreatedAt,
		},
	})
}

func (h *ChatHandler) EditGroupChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.EditGroupChatRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.ChatID == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "chat_id is required",
		})
		return
	}

	chatInterface, exists := h.chats.Load(req.ChatID)
	if !exists {
		msgpack.Respond(w, http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Chat not found",
		})
		return
	}

	chat := chatInterface.(*chatData)

	// Check if user is owner
	if chat.OwnerID != userID {
		msgpack.Respond(w, http.StatusForbidden, models.Response{
			Success: false,
			Error:   "Only owner can edit group chat",
		})
		return
	}

	if req.Name != "" {
		chat.Name = req.Name
	}
	if req.Description != "" {
		chat.Description = req.Description
	}

	h.chats.Store(req.ChatID, chat)

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:      chat.ChatID,
			Type:        chat.Type,
			Name:        chat.Name,
			Description: chat.Description,
			MemberIDs:   chat.MemberIDs,
			OwnerID:     chat.OwnerID,
			CreatedAt:   chat.CreatedAt,
		},
	})
}

func (h *ChatHandler) EditChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.EditChannelRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.ChannelID == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "channel_id is required",
		})
		return
	}

	chatInterface, exists := h.chats.Load(req.ChannelID)
	if !exists {
		msgpack.Respond(w, http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Channel not found",
		})
		return
	}

	chat := chatInterface.(*chatData)

	// Check if user is owner
	if chat.OwnerID != userID {
		msgpack.Respond(w, http.StatusForbidden, models.Response{
			Success: false,
			Error:   "Only owner can edit channel",
		})
		return
	}

	if req.Name != "" {
		chat.Name = req.Name
	}
	if req.Description != "" {
		chat.Description = req.Description
	}

	h.chats.Store(req.ChannelID, chat)

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:      chat.ChatID,
			Type:        chat.Type,
			Name:        chat.Name,
			Description: chat.Description,
			MemberIDs:   chat.MemberIDs,
			OwnerID:     chat.OwnerID,
			CreatedAt:   chat.CreatedAt,
		},
	})
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "chat_id query parameter is required",
		})
		return
	}

	chatInterface, exists := h.chats.Load(chatID)
	if !exists {
		msgpack.Respond(w, http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Chat not found",
		})
		return
	}

	chat := chatInterface.(*chatData)

	// Check if user is member
	isMember := false
	for _, memberID := range chat.MemberIDs {
		if memberID == userID {
			isMember = true
			break
		}
	}

	if !isMember && chat.Type != "public_channel" {
		msgpack.Respond(w, http.StatusForbidden, models.Response{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.ChatResponse{
			ChatID:      chat.ChatID,
			Type:        chat.Type,
			Name:        chat.Name,
			Description: chat.Description,
			MemberIDs:   chat.MemberIDs,
			OwnerID:     chat.OwnerID,
			CreatedAt:   chat.CreatedAt,
			InviteLink:  chat.InviteLink,
		},
	})
}

func generateInviteLink() string {
	return "https://messenger.app/join/" + uuid.New().String()[:8]
}
