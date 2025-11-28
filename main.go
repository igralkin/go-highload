package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/igralkin/go-highload/handlers"
	"github.com/igralkin/go-highload/metrics"
	"github.com/igralkin/go-highload/services"
	"github.com/igralkin/go-highload/utils"
)

func main() {
	userService := services.NewUserService()

	// Конфиг для MinIO из переменных окружения
	minioEndpoint := getenv("MINIO_ENDPOINT", "minio:9000")
	minioAccessKey := getenv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getenv("MINIO_SECRET_KEY", "minioadmin")
	minioUseSSL := false
	minioBucket := getenv("MINIO_BUCKET", "go-highload")

	integrationService, err := services.NewIntegrationService(
		minioEndpoint,
		minioAccessKey,
		minioSecretKey,
		minioBucket,
		minioUseSSL,
	)
	if err != nil {
		log.Fatalf("failed to init integration service: %v", err)
	}

	auditLogger := utils.NewAuditLogger(1000)
	notificationService := services.NewNotificationService(1000)

	userHandler := handlers.NewUserHandler(userService, auditLogger, notificationService)
	integrationHandler := handlers.NewIntegrationHandler(userService, integrationService)

	r := mux.NewRouter()
	r.Use(metrics.MetricsMiddleware)
	r.Use(utils.RateLimitMiddleware)

	userHandler.RegisterRoutes(r)
	integrationHandler.RegisterRoutes(r)

	r.Handle("/metrics", metrics.Handler()).Methods(http.MethodGet)

	addr := ":8080"
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
