package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/igralkin/go-highload/services"
	"github.com/igralkin/go-highload/utils"
)

type UserHandler struct {
	service             *services.UserService
	auditLogger         *utils.AuditLogger
	notificationService *services.NotificationService
}

func NewUserHandler(
	service *services.UserService,
	auditLogger *utils.AuditLogger,
	notificationService *services.NotificationService,
) *UserHandler {
	return &UserHandler{
		service:             service,
		auditLogger:         auditLogger,
		notificationService: notificationService,
	}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/users", h.GetAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/users/{id}", h.GetUserByID).Methods(http.MethodGet)
	r.HandleFunc("/api/users", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/api/users/{id}", h.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/api/users/{id}", h.DeleteUser).Methods(http.MethodDelete)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAll()
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user := h.service.Create(req.Name, req.Email)
	utils.WriteJSON(w, http.StatusCreated, user)

	h.auditLogger.Log("CREATE", user)
	h.notificationService.Notify(services.Notification{
		Type: services.NotificationUserCreated,
		User: user,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := h.service.Update(id, req.Name, req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)

	h.auditLogger.Log("UPDATE", user)
	h.notificationService.Notify(services.Notification{
		Type: services.NotificationUserUpdated,
		User: user,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	h.auditLogger.Log("DELETE", user)
	h.notificationService.Notify(services.Notification{
		Type: services.NotificationUserDeleted,
		User: user,
	})
}

func parseID(raw string) (int, error) {
	return strconv.Atoi(raw)
}
