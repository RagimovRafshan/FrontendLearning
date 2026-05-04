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

type UserHandler struct {
	elastic   *storage.ElasticsearchClient
	onlineUsers sync.Map // Map[userID]lastSeen
}

func NewUserHandler(elastic *storage.ElasticsearchClient) *UserHandler {
	return &UserHandler{
		elastic: elastic,
	}
}

func (h *UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.FindUserRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Query == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Query is required",
		})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Search in Elasticsearch
	users, err := h.elastic.SearchUsers(r.Context(), req.Query, req.Limit)
	if err != nil {
		msgpack.Respond(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Search failed: " + err.Error(),
		})
		return
	}

	// Add online status
	var results []models.UserResponse
	for _, user := range users {
		lastSeen, isOnline := h.onlineUsers.Load(user.UserID)
		results = append(results, models.UserResponse{
			UserID:      user.UserID,
			Login:       user.Login,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			Bio:         user.Bio,
			IsOnline:    isOnline && time.Since(lastSeen.(time.Time)) < 5*time.Minute,
			LastSeen:    lastSeen.(time.Time).UnixMilli(),
		})
	}

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data:    results,
	})
}

func (h *UserHandler) EditProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	login, _ := ctx.Value("login").(string)
	
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req models.EditProfileRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update in Elasticsearch
	err := h.elastic.UpdateUserProfile(r.Context(), userID, login, req.DisplayName, req.Avatar, req.Bio)
	if err != nil {
		msgpack.Respond(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to update profile: " + err.Error(),
		})
		return
	}

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data:    "Profile updated successfully",
	})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestedUserID := r.URL.Query().Get("user_id")
	
	if requestedUserID == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "user_id query parameter is required",
		})
		return
	}

	user, err := h.elastic.GetUserByID(r.Context(), requestedUserID)
	if err != nil {
		msgpack.Respond(w, http.StatusNotFound, models.Response{
			Success: false,
			Error:   "User not found: " + err.Error(),
		})
		return
	}

	lastSeen, _ := h.onlineUsers.Load(requestedUserID)
	isOnline := lastSeen != nil && time.Since(lastSeen.(time.Time)) < 5*time.Minute

	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data: models.UserResponse{
			UserID:      user.UserID,
			Login:       user.Login,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			Bio:         user.Bio,
			IsOnline:    isOnline,
			LastSeen:    lastSeen.(time.Time).UnixMilli(),
		},
	})
}

func (h *UserHandler) SetOnline(userID string, online bool) {
	if online {
		h.onlineUsers.Store(userID, time.Now())
	} else {
		h.onlineUsers.Delete(userID)
	}
}

func (h *UserHandler) IsUserOnline(userID string) bool {
	lastSeen, exists := h.onlineUsers.Load(userID)
	if !exists {
		return false
	}
	return time.Since(lastSeen.(time.Time)) < 5*time.Minute
}

// User data structure for Elasticsearch
type UserData struct {
	UserID      string `json:"user_id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Bio         string `json:"bio,omitempty"`
	CreatedAt   int64  `json:"created_at"`
}
