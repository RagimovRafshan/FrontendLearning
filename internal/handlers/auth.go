package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"messenger/internal/models"
	"messenger/internal/storage"
	"messenger/pkg/msgpack"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users      sync.Map // Map[login]passwordHash
	userIDs    sync.Map // Map[login]userID
	sessions   sync.Map // Map[userID]sessionInfo
	redis      *storage.RedisClient
	jwtSecret  []byte
}

type sessionInfo struct {
	UserID       string
	Login        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func NewAuthHandler(redis *storage.RedisClient, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		redis:     redis,
		jwtSecret: []byte(jwtSecret),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID != "" {
		msgpack.Respond(w, http.StatusForbidden, models.Response{
			Success: false,
			Error:   "Already authenticated",
		})
		return
	}

	var req models.RegisterRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Login == "" || req.Password == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Login and password are required",
		})
		return
	}

	// Check if user exists
	if _, exists := h.users.Load(req.Login); exists {
		msgpack.Respond(w, http.StatusConflict, models.Response{
			Success: false,
			Error:   "User already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		msgpack.Respond(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to hash password",
		})
		return
	}

	// Generate user ID
	newUserID := uuid.New().String()

	// Store user
	h.users.Store(req.Login, string(hashedPassword))
	h.userIDs.Store(req.Login, newUserID)

	// Generate tokens
	accessToken, refreshToken := h.generateTokens(newUserID, req.Login)

	msgpack.Respond(w, http.StatusOK, models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       newUserID,
		Login:        req.Login,
		ExpiresAt:    time.Now().Add(15 * time.Minute).Unix(),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	if userID != "" {
		msgpack.Respond(w, http.StatusForbidden, models.Response{
			Success: false,
			Error:   "Already authenticated",
		})
		return
	}

	var req models.LoginRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Get stored password hash
	storedHash, exists := h.users.Load(req.Login)
	if !exists {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash.(string)), []byte(req.Password)); err != nil {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// Get user ID
	userID, _ = h.userIDs.Load(req.Login)

	// Generate tokens
	accessToken, refreshToken := h.generateTokens(userID.(string), req.Login)

	msgpack.Respond(w, http.StatusOK, models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID.(string),
		Login:        req.Login,
		ExpiresAt:    time.Now().Add(15 * time.Minute).Unix(),
	})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.Email == "" {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Email is required",
		})
		return
	}

	// In production: send reset link to email
	// For now, just acknowledge
	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data:    "Password reset link sent to email",
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := msgpack.DecodeRequest(r, &req); err != nil {
		msgpack.Respond(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate refresh token and get user info
	// Simplified for demo - in production use JWT validation
	userID := h.validateRefreshToken(req.RefreshToken)
	if userID == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Invalid refresh token",
		})
		return
	}

	// Get login from session
	var login string
	h.sessions.Range(func(key, value interface{}) bool {
		if info, ok := value.(*sessionInfo); ok && info.UserID == userID {
			login = info.Login
			return false
		}
		return true
	})

	if login == "" {
		msgpack.Respond(w, http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Session not found",
		})
		return
	}

	// Generate new tokens
	accessToken, newRefreshToken := h.generateTokens(userID, login)

	msgpack.Respond(w, http.StatusOK, models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		UserID:       userID,
		Login:        login,
		ExpiresAt:    time.Now().Add(15 * time.Minute).Unix(),
	})
}

func (h *AuthHandler) generateTokens(userID, login string) (string, string) {
	accessToken := h.generateToken(32)
	refreshToken := h.generateToken(64)

	expiresAt := time.Now().Add(15 * time.Minute)

	h.sessions.Store(userID, &sessionInfo{
		UserID:       userID,
		Login:        login,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	})

	return accessToken, refreshToken
}

func (h *AuthHandler) generateToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (h *AuthHandler) validateRefreshToken(token string) string {
	var userID string
	h.sessions.Range(func(key, value interface{}) bool {
		if info, ok := value.(*sessionInfo); ok && info.RefreshToken == token {
			userID = info.UserID
			return false
		}
		return true
	})
	return userID
}

func (h *AuthHandler) GetUserByToken(token string) (string, string, bool) {
	var userID, login string
	h.sessions.Range(func(key, value interface{}) bool {
		if info, ok := value.(*sessionInfo); ok && info.AccessToken == token {
			userID = info.UserID
			login = info.Login
			return false
		}
		return true
	})
	return userID, login, userID != ""
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(string)
	
	if userID != "" {
		h.sessions.Delete(userID)
	}
	
	msgpack.Respond(w, http.StatusOK, models.Response{
		Success: true,
		Data:    "Logged out successfully",
	})
}
