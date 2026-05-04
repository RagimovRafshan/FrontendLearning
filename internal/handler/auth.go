package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/you/messenger/internal/model"
	"github.com/you/messenger/internal/service"
	"github.com/you/messenger/pkg/msgpack"
)

// AuthHandler обработчик запросов аутентификации
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler создает новый обработчик аутентификации
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register обрабатывает запрос на регистрацию
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.RegisterRequest
	if err := h.decodeRequest(r, &req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrUserAlreadyExists:
			h.writeError(w, "User already exists", http.StatusConflict)
		default:
			h.writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.writeSuccess(w, user, http.StatusCreated)
}

// Login обрабатывает запрос на вход
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.LoginRequest
	if err := h.decodeRequest(r, &req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			h.writeError(w, "Invalid login or password", http.StatusUnauthorized)
		default:
			h.writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.writeSuccess(w, resp, http.StatusOK)
}

// RefreshToken обрабатывает запрос на обновление токена
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.RefreshTokenRequest
	if err := h.decodeRequest(r, &req); err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidToken:
			h.writeError(w, "Invalid token", http.StatusUnauthorized)
		default:
			h.writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.writeSuccess(w, resp, http.StatusOK)
}

// Logout обрабатывает запрос на выход
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get tokens from headers or body
	accessToken := r.Header.Get("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	var req struct {
		RefreshToken string `msgpack:"refresh_token"`
	}
	
	body, _ := io.ReadAll(r.Body)
	msgpack.Unmarshal(body, &req)

	if err := h.authService.Logout(r.Context(), accessToken, req.RefreshToken); err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccess(w, map[string]string{"status": "ok"}, http.StatusOK)
}

// decodeRequest декодирует запрос из MessagePack или JSON
func (h *AuthHandler) decodeRequest(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return nil
	}

	// Try MessagePack first (default for our API)
	if contentType == "application/msgpack" || contentType == "" {
		return msgpack.Unmarshal(body, v)
	}

	// Fallback to JSON
	if contentType == "application/json" {
		return json.Unmarshal(body, v)
	}

	// Default to MessagePack
	return msgpack.Unmarshal(body, v)
}

// writeSuccess записывает успешный ответ
func (h *AuthHandler) writeSuccess(w http.ResponseWriter, data interface{}, status int) {
	response := model.SuccessResponse{Data: data}
	
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		contentType = "application/msgpack"
	}
	
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)

	if contentType == "application/msgpack" {
		bytes, err := msgpack.Marshal(response)
		if err != nil {
			h.writeError(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
	} else {
		json.NewEncoder(w).Encode(response)
	}
}

// writeError записывает ответ с ошибкой
func (h *AuthHandler) writeError(w http.ResponseWriter, message string, code int) {
	response := model.ErrorResponse{
		Error: model.APIError{
			Code:    code,
			Message: message,
		},
	}
	
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		contentType = "application/msgpack"
	}
	
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)

	if contentType == "application/msgpack" {
		bytes, err := msgpack.Marshal(response)
		if err != nil {
			w.Write([]byte(`{"error":{"code":500,"message":"Failed to encode error"}}`))
			return
		}
		w.Write(bytes)
	} else {
		json.NewEncoder(w).Encode(response)
	}
}
