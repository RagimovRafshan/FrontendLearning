package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter middleware для ограничения запросов
type RateLimiter struct {
	limiters sync.Map
	rps      int
	burst    int
}

// NewRateLimiter создает новый rate limiter
func NewRateLimiter(rps, burst int) *RateLimiter {
	return &RateLimiter{
		rps:   rps,
		burst: burst,
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	if v, ok := rl.limiters.Load(key); ok {
		return v.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
	rl.limiters.Store(key, limiter)
	
	// Cleanup old limiters periodically
	go func() {
		time.Sleep(time.Minute)
		rl.limiters.Delete(key)
	}()
	
	return limiter
}

// Middleware возвращает HTTP middleware
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}
		
		limiter := rl.getLimiter(ip)
		
		if !limiter.Allow() {
			http.Error(w, `{"error":{"code":429,"message":"Too many requests"}}`, http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// CORS middleware для обработки CORS запросов
type CORS struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
	maxAge         int
}

// NewCORS создает новый CORS middleware
func NewCORS(origins, methods, headers []string, maxAge int) *CORS {
	return &CORS{
		allowedOrigins: origins,
		allowedMethods: methods,
		allowedHeaders: headers,
		maxAge:         maxAge,
	}
}

// Middleware возвращает HTTP middleware
func (c *CORS) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, o := range c.allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", join(c.allowedMethods))
		w.Header().Set("Access-Control-Allow-Headers", join(c.allowedHeaders))
		w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", c.maxAge))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func join(slice []string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// Recovery middleware для восстановления после паник
type Recovery struct{}

// NewRecovery создает новый recovery middleware
func NewRecovery() *Recovery {
	return &Recovery{}
}

// Middleware возвращает HTTP middleware
func (r *Recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, `{"error":{"code":500,"message":"Internal server error"}}`, http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, req)
	})
}

// Logger middleware для логирования запросов
type Logger struct{}

// NewLogger создает новый logger middleware
func NewLogger() *Logger {
	return &Logger{}
}

// Middleware возвращает HTTP middleware
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		// Log request details here if needed
		_ = duration
	})
}

// contextKey тип для ключей контекста
type contextKey string

const userIDKey contextKey = "user_id"

// WithUserID добавляет user ID в контекст
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext получает user ID из контекста
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// Auth middleware для проверки аутентификации
type Auth struct {
	secretKey []byte
}

// NewAuth создает новый auth middleware
func NewAuth(secretKey string) *Auth {
	return &Auth{
		secretKey: []byte(secretKey),
	}
}

// Middleware возвращает HTTP middleware
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for certain paths
		if r.URL.Path == "/api/v1/auth/login" || 
		   r.URL.Path == "/api/v1/auth/register" ||
		   r.URL.Path == "/api/v1/auth/refresh" {
			next.ServeHTTP(w, r)
			return
		}
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":{"code":401,"message":"Unauthorized"}}`, http.StatusUnauthorized)
			return
		}
		
		// Extract token from "Bearer <token>"
		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
		
		// TODO: Validate JWT token and extract user ID
		// For now, just pass through
		ctx := WithUserID(r.Context(), "user_id_from_token")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
