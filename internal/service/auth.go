package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/you/messenger/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

// AuthService сервис аутентификации
type AuthService struct {
	users      sync.Map // Map[login]passwordHash
	tokens     sync.Map // Map[token]tokenData
	mu         sync.RWMutex
	bcryptCost int
	tokenExpiry time.Duration
}

type tokenData struct {
	userID    string
	login     string
	expiresAt time.Time
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(bcryptCost int, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		bcryptCost: bcryptCost,
		tokenExpiry: tokenExpiry,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if user already exists
	if _, exists := s.users.Load(req.Login); exists {
		return nil, ErrUserAlreadyExists
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}
	
	user := &model.User{
		ID:        generateID(),
		Login:     req.Login,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Store password hash (not the plain password)
	s.users.Store(req.Login, string(hashedPassword))
	
	return user, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Get password hash
	passwordHash, exists := s.users.Load(req.Login)
	if !exists {
		return nil, ErrInvalidCredentials
	}
	
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash.(string)), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	
	// Generate tokens
	accessToken, refreshToken := generateTokens()
	expiresAt := time.Now().Add(s.tokenExpiry)
	
	// Store token data
	s.tokens.Store(accessToken, &tokenData{
		userID:    "user_id_from_login_" + req.Login,
		login:     req.Login,
		expiresAt: expiresAt,
	})
	
	s.tokens.Store(refreshToken, &tokenData{
		userID:    "user_id_from_login_" + req.Login,
		login:     req.Login,
		expiresAt: expiresAt.Add(24 * time.Hour * 7), // Refresh token valid for 7 days
	})
	
	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: &model.User{
			ID:    "user_id_from_login_" + req.Login,
			Login: req.Login,
		},
	}, nil
}

// RefreshToken обновляет access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data, exists := s.tokens.Load(refreshToken)
	if !exists {
		return nil, ErrInvalidToken
	}
	
	tokenData := data.(*tokenData)
	
	// Check if refresh token is expired
	if time.Now().After(tokenData.expiresAt) {
		s.tokens.Delete(refreshToken)
		return nil, ErrInvalidToken
	}
	
	// Generate new tokens
	newAccessToken, newRefreshToken := generateTokens()
	expiresAt := time.Now().Add(s.tokenExpiry)
	
	// Store new tokens
	s.tokens.Store(newAccessToken, &tokenData{
		userID:    tokenData.userID,
		login:     tokenData.login,
		expiresAt: expiresAt,
	})
	
	s.tokens.Store(newRefreshToken, &tokenData{
		userID:    tokenData.userID,
		login:     tokenData.login,
		expiresAt: time.Now().Add(24 * time.Hour * 7),
	})
	
	// Delete old tokens
	s.tokens.Delete(refreshToken)
	
	return &model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User: &model.User{
			ID:    tokenData.userID,
			Login: tokenData.login,
		},
	}, nil
}

// Logout удаляет токены
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	s.tokens.Delete(accessToken)
	s.tokens.Delete(refreshToken)
	return nil
}

// ValidateToken проверяет токен и возвращает данные пользователя
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*tokenData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data, exists := s.tokens.Load(token)
	if !exists {
		return nil, ErrInvalidToken
	}
	
	tokenData := data.(*tokenData)
	
	// Check if token is expired
	if time.Now().After(tokenData.expiresAt) {
		return nil, ErrInvalidToken
	}
	
	return tokenData, nil
}

// generateID генерирует уникальный ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// generateTokens генерирует пару токенов
func generateTokens() (string, string) {
	access := make([]byte, 32)
	refresh := make([]byte, 32)
	rand.Read(access)
	rand.Read(refresh)
	
	return base64.URLEncoding.EncodeToString(access), 
		   base64.URLEncoding.EncodeToString(refresh)
}
