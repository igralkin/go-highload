package services

import (
	"log"

	"github.com/igralkin/go-highload/models"
)

type NotificationType string

const (
	NotificationUserCreated NotificationType = "USER_CREATED"
	NotificationUserUpdated NotificationType = "USER_UPDATED"
	NotificationUserDeleted NotificationType = "USER_DELETED"
)

type Notification struct {
	Type NotificationType
	User models.User
}

type NotificationService struct {
	ch chan Notification
}

// NewNotificationService создаёт сервис и запускает воркер.
func NewNotificationService(bufferSize int) *NotificationService {
	ns := &NotificationService{
		ch: make(chan Notification, bufferSize),
	}
	go ns.run()
	return ns
}

// Notify отправляет уведомление в канал (также неблокирующе).
func (ns *NotificationService) Notify(n Notification) {
	select {
	case ns.ch <- n:
	default:
		log.Println("notification channel is full, dropping notification")
	}
}

// run — воркер, который имитирует отправку уведомлений.
func (ns *NotificationService) run() {
	for n := range ns.ch {
		// Здесь могла бы быть интеграция с email / SMS / очередью.
		log.Printf("[NOTIFY] type=%s user_id=%d name=%s email=%s\n",
			n.Type, n.User.ID, n.User.Name, n.User.Email)
	}
}
