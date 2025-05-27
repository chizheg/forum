package domain

import (
	"time"
)

// ForumService определяет расширенный интерфейс форума,
// включающий все функции чата
type ForumService interface {
	// Методы чата
	SendMessage(userID int, content string) error
	GetMessages(limit int) ([]*Message, error)
	DeleteOldMessages(maxAge time.Duration) error
	JoinChat(userID int) error
	LeaveChat(userID int) error
	IsParticipant(userID int) (bool, error)

	// Здесь могут быть добавлены дополнительные методы форума
}
