package main

import (
"context"
"crypto/tls"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/you/messenger/config"
"github.com/you/messenger/internal/handler"
"github.com/you/messenger/internal/middleware"
"github.com/you/messenger/internal/service"
"github.com/you/messenger/internal/storage"
)

func main() {
cfg, err := config.LoadFromFile(".env")
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

scyllaDB, err := storage.NewScyllaDB(cfg.ScyllaDB)
if err != nil {
log.Printf("Warning: Failed to connect to ScyllaDB: %v", err)
}
defer func() {
if scyllaDB != nil {
scyllaDB.Close()
}
}()

redisClient, err := storage.NewRedis(cfg.Redis)
if err != nil {
log.Printf("Warning: Failed to connect to Redis: %v", err)
}
defer func() {
if redisClient != nil {
redisClient.Close()
}
}()

minioClient, err := storage.NewMinIO(cfg.MinIO)
if err != nil {
log.Printf("Warning: Failed to connect to MinIO: %v", err)
}
_ = minioClient

elasticClient, err := storage.NewElasticSearch(cfg.ElasticSearch)
if err != nil {
log.Printf("Warning: Failed to connect to ElasticSearch: %v", err)
}
_ = elasticClient

authService := service.NewAuthService(cfg.Auth.BcryptCost, cfg.Auth.TokenExpiry)
authHandler := handler.NewAuthHandler(authService)

rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.RequestsPerSec, cfg.RateLimit.BurstSize)
cors := middleware.NewCORS(cfg.CORS.AllowedOrigins, cfg.CORS.AllowedMethods, cfg.CORS.AllowedHeaders, cfg.CORS.MaxAge)
recovery := middleware.NewRecovery()
logger := middleware.NewLogger()
auth := middleware.NewAuth("your-secret-key")

mux := http.NewServeMux()
mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
mux.HandleFunc("/api/v1/auth/refresh", authHandler.RefreshToken)
mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status":"ok"}`))
})

var handler http.Handler = mux
handler = logger.Middleware(handler)
handler = recovery.Middleware(handler)
handler = cors.Middleware(handler)
handler = rateLimiter.Middleware(handler)
handler = auth.Middleware(handler)

server := &http.Server{
Addr:           cfg.Server.Host + ":" + cfg.Server.Port,
Handler:        handler,
ReadTimeout:    cfg.Server.ReadTimeout,
WriteTimeout:   cfg.Server.WriteTimeout,
IdleTimeout:    cfg.Server.IdleTimeout,
MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
TLSConfig: &tls.Config{
MinVersion: tls.VersionTLS12,
CipherSuites: []uint16{
tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
},
},
}

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

log.Printf("Starting HTTPS server on %s:%s", cfg.Server.Host, cfg.Server.Port)
log.Println("API endpoints:")
log.Println("  POST /api/v1/auth/register - Register new user")
log.Println("  POST /api/v1/auth/login    - Login user")
log.Println("  POST /api/v1/auth/refresh  - Refresh token")
log.Println("  POST /api/v1/auth/logout   - Logout user")
log.Println("  GET  /health               - Health check")
log.Println("\nContent-Type: application/msgpack (default) or application/json")

if _, err := os.Stat(cfg.Server.CertFile); os.IsNotExist(err) {
log.Printf("Certificate file not found: %s", cfg.Server.CertFile)
log.Println("Please generate SSL certificates:")
log.Println("  mkdir -p certs")
log.Println("  openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes")
log.Fatal("Exiting...")
}

if err := server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil && err != http.ErrServerClosed {
log.Fatalf("Server error: %v", err)
}
log.Println("Server stopped")
}
