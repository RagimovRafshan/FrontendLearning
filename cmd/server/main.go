package main

import (
"context"
"crypto/tls"
"fmt"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"messenger/config"
"messenger/internal/handlers"
"messenger/internal/middleware"
"messenger/internal/storage"
"messenger/internal/websocket"

"github.com/gorilla/mux"
)

func main() {
// Загрузка конфигурации
cfg, err := config.Load()
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

// Инициализация хранилищ
scylla, err := storage.NewScyllaDBClient(cfg.ScyllaDB.Hosts, cfg.ScyllaDB.Keyspace)
if err != nil {
log.Printf("Warning: ScyllaDB connection failed: %v", err)
}

redis, err := storage.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Password)
if err != nil {
log.Printf("Warning: Redis connection failed: %v", err)
}

minio, err := storage.NewMinIOClient(cfg.MinIO.Endpoint, cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, cfg.MinIO.UseSSL)
if err != nil {
log.Printf("Warning: MinIO connection failed: %v", err)
}

elastic, err := storage.NewElasticsearchClient(cfg.Elasticsearch.URL)
if err != nil {
log.Printf("Warning: Elasticsearch connection failed: %v", err)
}

// Инициализация обработчиков
authHandler := handlers.NewAuthHandler(redis, cfg.JWT.Secret)
userHandler := handlers.NewUserHandler(elastic)
chatHandler := handlers.NewChatHandler(scylla, elastic)
messageHandler := handlers.NewMessageHandler(scylla)
wsManager := websocket.NewWSManager()

// Middleware
authMiddleware := middleware.NewAuthMiddleware(authHandler)
corsMiddleware := middleware.NewCORSMiddleware(nil, nil, nil)
rateLimitMiddleware := middleware.NewRateLimitMiddleware(1000, 60)
recoveryMiddleware := middleware.RecoveryMiddleware

// Создание роутера
r := mux.NewRouter()

// Применяем middleware
r.Use(recoveryMiddleware)
r.Use(corsMiddleware.Middleware)
r.Use(rateLimitMiddleware.Middleware)

// Public routes (без авторизации)
publicRouter := r.PathPrefix("/api/v1").Subrouter()
publicRouter.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
publicRouter.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
publicRouter.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword).Methods("POST")
publicRouter.HandleFunc("/auth/refresh-token", authHandler.RefreshToken).Methods("POST")

// Protected routes (с авторизацией)
protectedRouter := r.PathPrefix("/api/v1").Subrouter()
protectedRouter.Use(authMiddleware.Middleware)

// Auth
protectedRouter.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")

// Users
protectedRouter.HandleFunc("/users/search", userHandler.FindUsers).Methods("POST")
protectedRouter.HandleFunc("/users/profile", userHandler.EditProfile).Methods("PUT")
protectedRouter.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")

// Chats
protectedRouter.HandleFunc("/chats/create-by-link", chatHandler.CreateChatByLink).Methods("POST")
protectedRouter.HandleFunc("/chats/create-with-user", chatHandler.CreateChatWithUser).Methods("POST")
protectedRouter.HandleFunc("/chats/create-group", chatHandler.CreateGroupChat).Methods("POST")
protectedRouter.HandleFunc("/chats/create-channel", chatHandler.CreateChannel).Methods("POST")
protectedRouter.HandleFunc("/chats/edit-group", chatHandler.EditGroupChat).Methods("PUT")
protectedRouter.HandleFunc("/chats/edit-channel", chatHandler.EditChannel).Methods("PUT")
protectedRouter.HandleFunc("/chats/get", chatHandler.GetChat).Methods("GET")

// Messages
protectedRouter.HandleFunc("/messages/send", messageHandler.SendMessage).Methods("POST")
protectedRouter.HandleFunc("/messages/get", messageHandler.GetMessage).Methods("POST")
protectedRouter.HandleFunc("/messages/edit", messageHandler.EditMessage).Methods("PUT")
protectedRouter.HandleFunc("/messages/chat", messageHandler.GetChatMessages).Methods("GET")

// WebSocket для реального времени
r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// Проверяем авторизацию через токен в query param или заголовке
token := r.URL.Query().Get("token")
if token == "" {
token = r.Header.Get("Authorization")
if len(token) > 7 && token[:7] == "Bearer " {
token = token[7:]
}
}

if token == "" {
http.Error(w, "Unauthorized", http.StatusUnauthorized)
return
}

userID, login, valid := authHandler.GetUserByToken(token)
if !valid {
http.Error(w, "Invalid token", http.StatusUnauthorized)
return
}

// Добавляем в контекст и передаем в WS handler
ctx := context.WithValue(r.Context(), "user_id", userID)
ctx = context.WithValue(ctx, "login", login)
wsManager.HandleWS(w, r.WithContext(ctx))
}).Methods("GET")

// Health check
r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status":"ok"}`))
}).Methods("GET")

// HTTPS сервер
tlsConfig := &tls.Config{
MinVersion: tls.VersionTLS12,
CipherSuites: []uint16{
tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
},
PreferServerCipherSuites: true,
}

server := &http.Server{
Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
Handler:      r,
TLSConfig:    tlsConfig,
ReadTimeout:  15 * time.Second,
WriteTimeout: 15 * time.Second,
IdleTimeout:  60 * time.Second,
}

// Graceful shutdown
go func() {
sigint := make(chan os.Signal, 1)
signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
<-sigint

log.Println("Shutting down server...")
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
log.Printf("Server shutdown error: %v", err)
}
}()

// Запуск сервера
log.Printf("Starting HTTPS server on port %d...", cfg.Server.Port)

// Для продакшена нужны реальные сертификаты
// Для разработки можно использовать self-signed
certFile := cfg.Server.CertFile
keyFile := cfg.Server.KeyFile

if certFile == "" || keyFile == "" {
log.Println("Warning: No TLS certificates configured. Using HTTP for development.")
log.Printf("Server started at http://localhost:%d", cfg.Server.Port)
if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatalf("Server error: %v", err)
}
} else {
log.Printf("Server started at https://localhost:%d", cfg.Server.Port)
if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
log.Fatalf("Server error: %v", err)
}
}
}
