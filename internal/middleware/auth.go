package middleware

import (
"context"
"net/http"
"strconv"
"strings"
"sync"
"time"
)

type AuthMiddleware struct {
authHandler interface {
GetUserByToken(token string) (string, string, bool)
}
}

func NewAuthMiddleware(authHandler interface{ GetUserByToken(string) (string, string, bool) }) *AuthMiddleware {
return &AuthMiddleware{
authHandler: authHandler,
}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
http.Error(w, "Missing authorization header", http.StatusUnauthorized)
return
}

parts := strings.Split(authHeader, " ")
if len(parts) != 2 || parts[0] != "Bearer" {
http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
return
}

token := parts[1]
userID, login, valid := m.authHandler.GetUserByToken(token)
if !valid {
http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
return
}

ctx := context.WithValue(r.Context(), "user_id", userID)
ctx = context.WithValue(ctx, "login", login)
ctx = context.WithValue(ctx, "token", token)

next.ServeHTTP(w, r.WithContext(ctx))
})
}

type CORSMiddleware struct {
allowedOrigins []string
allowedMethods []string
allowedHeaders []string
}

func NewCORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders []string) *CORSMiddleware {
if len(allowedOrigins) == 0 {
allowedOrigins = []string{"*"}
}
if len(allowedMethods) == 0 {
allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
}
if len(allowedHeaders) == 0 {
allowedHeaders = []string{"Authorization", "Content-Type", "X-Request-ID"}
}

return &CORSMiddleware{
allowedOrigins: allowedOrigins,
allowedMethods: allowedMethods,
allowedHeaders: allowedHeaders,
}
}

func (m *CORSMiddleware) Middleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
origin := r.Header.Get("Origin")
allowed := false
for _, o := range m.allowedOrigins {
if o == "*" || o == origin {
allowed = true
break
}
}

if allowed {
w.Header().Set("Access-Control-Allow-Origin", origin)
if origin == "" {
w.Header().Set("Access-Control-Allow-Origin", "*")
}
}

w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.allowedMethods, ", "))
w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.allowedHeaders, ", "))
w.Header().Set("Access-Control-Allow-Credentials", "true")
w.Header().Set("Access-Control-Max-Age", "86400")

if r.Method == "OPTIONS" {
w.WriteHeader(http.StatusOK)
return
}

next.ServeHTTP(w, r)
})
}

type RateLimitMiddleware struct {
requests map[string][]int64
limit    int
window   int64
mu       sync.Mutex
}

func NewRateLimitMiddleware(limit int, windowSeconds int64) *RateLimitMiddleware {
return &RateLimitMiddleware{
requests: make(map[string][]int64),
limit:    limit,
window:   windowSeconds,
}
}

func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
clientIP := r.RemoteAddr
if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
clientIP = strings.Split(forwarded, ",")[0]
}

now := time.Now().Unix()

m.mu.Lock()
var validRequests []int64
for _, t := range m.requests[clientIP] {
if now-t < m.window {
validRequests = append(validRequests, t)
}
}

if len(validRequests) >= m.limit {
m.mu.Unlock()
w.Header().Set("Retry-After", strconv.FormatInt(m.window, 10))
http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
return
}

m.requests[clientIP] = append(validRequests, now)
m.mu.Unlock()

next.ServeHTTP(w, r)
})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
defer func() {
if err := recover(); err != nil {
http.Error(w, "Internal server error", http.StatusInternalServerError)
}
}()
next.ServeHTTP(w, r)
})
}

func LoggingMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
start := time.Now()
next.ServeHTTP(w, r)
_ = start
})
}
