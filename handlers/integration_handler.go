package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/igralkin/go-highload/services"
	"github.com/igralkin/go-highload/utils"
)

type IntegrationHandler struct {
	userService        *services.UserService
	integrationService *services.IntegrationService
}

func NewIntegrationHandler(
	userService *services.UserService,
	integrationService *services.IntegrationService,
) *IntegrationHandler {
	return &IntegrationHandler{
		userService:        userService,
		integrationService: integrationService,
	}
}

func (h *IntegrationHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/integration/save-users", h.SaveUsersToMinio).Methods(http.MethodPost)
	r.HandleFunc("/api/integration/list-objects", h.ListObjects).Methods(http.MethodGet)
}

func (h *IntegrationHandler) SaveUsersToMinio(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users := h.userService.GetAll()
	objectName, err := h.integrationService.SaveUsers(ctx, users)
	if err != nil {
		http.Error(w, "failed to save users to MinIO: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"status":      "ok",
		"object_name": objectName,
		"user_count":  len(users),
	})
}

func (h *IntegrationHandler) ListObjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	objects, err := h.integrationService.ListObjects(ctx)
	if err != nil {
		http.Error(w, "failed to list objects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"objects": objects,
	})
}
