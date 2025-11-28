package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/igralkin/go-highload/handlers"
	"github.com/igralkin/go-highload/metrics"
	"github.com/igralkin/go-highload/services"
	"github.com/igralkin/go-highload/utils"
)

func main() {
	userService := services.NewUserService()

	auditLogger := utils.NewAuditLogger(1000)
	notificationService := services.NewNotificationService(1000)

	userHandler := handlers.NewUserHandler(userService, auditLogger, notificationService)

	r := mux.NewRouter()

	// Здесь важен порядок: сначала метрики, потом rate limiting.
	// Так метрики будут видеть и успешные запросы, и 429.
	r.Use(metrics.MetricsMiddleware)
	r.Use(utils.RateLimitMiddleware)

	// CRUD-роуты
	userHandler.RegisterRoutes(r)

	// Endpoint для Prometheus
	r.Handle("/metrics", metrics.Handler()).Methods(http.MethodGet)

	addr := ":8080"
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
