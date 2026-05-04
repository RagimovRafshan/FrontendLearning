package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"messenger/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
)

// AuthService управляет аутентификацией
type AuthService struct {
	// Map<login, passwordHash> - как просил пользователь
	users     map[string]string
	userData  map[string]*model.User
	mu        sync.RWMutex
	
	// Токены: token -> userID
	tokens    map[string]tokenInfo
	tokenMu   sync.RWMutex
	
	// Сессии для refresh токенов
	refreshTokens map[string]string // token -> userID
	refreshMu     sync.RWMutex
}

type tokenInfo struct {
	UserID    string
	ExpiresAt time.Time
}

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

// NewAuthService создает сервис аутентификации
func NewAuthService() *AuthService {
	return &AuthService{
		users:         make(map[string]string),
		userData:      make(map[string]*model.User),
		tokens:        make(map[string]tokenInfo),
		refreshTokens: make(map[string]string),
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Проверяем существует ли пользователь
	if _, exists := s.users[req.Username]; exists {
		return nil, ErrUserExists
	}
	
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// Создаем пользователя
	userID := generateID()
	user := &model.User{
		ID:          userID,
		Username:    req.Username,
		Email:       req.Email,
		CreatedAt:   time.Now(),
		IsOnline:    false,
	}
	
	if req.Email != "" {
		user.DisplayName = req.Username
	}
	
	// Сохраняем в Map<login, passwordHash>
	s.users[req.Username] = string(hashedPassword)
	s.userData[userID] = user
	
	// Генерируем токены
	accessToken, refreshToken := s.generateTokens(userID)
	
	return &model.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Ищем пользователя по логину в Map
	passwordHash, exists := s.users[req.Username]
	if !exists {
		return nil, ErrUserNotFound
	}
	
	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}
	
	// Находим данные пользователя
	var userData *model.User
	for _, user := range s.userData {
		if user.Username == req.Username {
			userData = user
			break
		}
	}
	
	if userData == nil {
		return nil, ErrUserNotFound
	}
	
	// Генерируем токены
	accessToken, refreshToken := s.generateTokens(userData.ID)
	
	return &model.AuthResponse{
		User:         *userData,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ForgotPassword инициирует восстановление пароля
func (s *AuthService) ForgotPassword(ctx context.Context, req *model.ForgotPasswordRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// В реальном приложении здесь была бы отправка email
	// Для демо просто проверяем существование пользователя
	_, exists := s.users[req.EmailOrUsername]
	if !exists {
		// Не сообщаем是否存在 пользователя из соображений безопасности
		return nil
	}
	
	// TODO: Отправить email с токеном восстановления
	return nil
}

// RefreshToken обновляет access токен
func (s *AuthService) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.AuthResponse, error) {
	s.refreshMu.RLock()
	userID, exists := s.refreshTokens[req.RefreshToken]
	s.refreshMu.RUnlock()
	
	if !exists {
		return nil, ErrInvalidToken
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, exists := s.userData[userID]
	if !exists {
		return nil, ErrUserNotFound
	}
	
	// Генерируем новые токены
	accessToken, newRefreshToken := s.generateTokens(userID)
	
	// Удаляем старый refresh токен
	s.refreshMu.Lock()
	delete(s.refreshTokens, req.RefreshToken)
	s.refreshTokens[newRefreshToken] = userID
	s.refreshMu.Unlock()
	
	return &model.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// ValidateToken проверяет validity токена
func (s *AuthService) ValidateToken(token string) (string, error) {
	s.tokenMu.RLock()
	defer s.tokenMu.RUnlock()
	
	info, exists := s.tokens[token]
	if !exists {
		return "", ErrInvalidToken
	}
	
	if time.Now().After(info.ExpiresAt) {
		return "", ErrTokenExpired
	}
	
	return info.UserID, nil
}

// GetUserByID получает пользователя по ID
func (s *AuthService) GetUserByID(userID string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, exists := s.userData[userID]
	if !exists {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

// SearchUsers ищет пользователей по запросу
func (s *AuthService) SearchUsers(query string, limit int) []*model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if limit <= 0 {
		limit = 10
	}
	
	var results []*model.User
	count := 0
	
	for _, user := range s.userData {
		if count >= limit {
			break
		}
		
		// Поиск по username или displayName
		if containsIgnoreCase(user.Username, query) || 
		   containsIgnoreCase(user.DisplayName, query) {
			results = append(results, user)
			count++
		}
	}
	
	return results
}

// UpdateProfile обновляет профиль пользователя
func (s *AuthService) UpdateProfile(userID string, req *model.UpdateProfileRequest) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	user, exists := s.userData[userID]
	if !exists {
		return nil, ErrUserNotFound
	}
	
	if req.Username != "" {
		// Проверяем уникальность username
		if _, exists := s.users[req.Username]; exists && req.Username != user.Username {
			return nil, ErrUserExists
		}
		// Обновляем ключ в мапе
		delete(s.users, user.Username)
		user.Username = req.Username
		s.users[req.Username] = s.users[user.Username] // сохраняем тот же хеш
	}
	
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	
	return user, nil
}

// GenerateTokens генерирует пару токенов
func (s *AuthService) generateTokens(userID string) (string, string) {
	accessToken := generateSecureToken()
	refreshToken := generateSecureToken()
	
	// Сохраняем access токен
	s.tokenMu.Lock()
	s.tokens[accessToken] = tokenInfo{
		UserID:    userID,
		ExpiresAt: time.Now().Add(AccessTokenExpiry),
	}
	s.tokenMu.Unlock()
	
	// Сохраняем refresh токен
	s.refreshMu.Lock()
	s.refreshTokens[refreshToken] = userID
	s.refreshMu.Unlock()
	
	// Очищаем старые токены
	go s.cleanupTokens()
	
	return accessToken, refreshToken
}

// cleanupTokens удаляет просроченные токены
func (s *AuthService) cleanupTokens() {
	s.tokenMu.Lock()
	now := time.Now()
	for token, info := range s.tokens {
		if now.After(info.ExpiresAt) {
			delete(s.tokens, token)
		}
	}
	s.tokenMu.Unlock()
}

// generateSecureToken генерирует криптографически безопасный токен
func generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// generateID генерирует уникальный ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// containsIgnoreCase проверяет содержит ли строка подстроку (без учета регистра)
func containsIgnoreCase(s, substr string) bool {
	if s == "" || substr == "" {
		return substr == ""
	}
	return len(s) >= len(substr) && 
		   (s[:len(substr)] == substr || 
		    findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
