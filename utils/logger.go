package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/igralkin/go-highload/models"
)

type AuditLogger struct {
	ch chan string
}

// NewAuditLogger создаёт логгер и запускает воркер в отдельной goroutine.
func NewAuditLogger(bufferSize int) *AuditLogger {
	l := &AuditLogger{
		ch: make(chan string, bufferSize),
	}
	go l.run()
	return l
}

// Log формирует сообщение и отправляет его в канал (не блокируясь, если канал забит).
func (l *AuditLogger) Log(action string, user models.User) {
	msg := fmt.Sprintf("action=%s user_id=%d name=%s email=%s time=%s",
		action, user.ID, user.Name, user.Email, time.Now().Format(time.RFC3339))

	select {
	case l.ch <- msg:
		// успешно отправили в канал
	default:
		// канал переполнен — можно дропнуть сообщение, чтобы не блокироваться
		log.Println("audit logger channel is full, dropping log message")
	}
}

// run — воркер, который крутится в отдельной goroutine и пишет логи.
func (l *AuditLogger) run() {
	for msg := range l.ch {
		log.Println("[AUDIT]", msg)
	}
}
