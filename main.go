package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/igralkin/go-highload/handlers"
	"github.com/igralkin/go-highload/services"
)

func main() {
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	addr := ":8080"
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
