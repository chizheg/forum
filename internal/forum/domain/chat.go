package domain

import (
	"time"
)

// Message represents a chat message
type Message struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Participant represents a chat participant
type Participant struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

// Repository defines the interface for chat data access
type Repository interface {
	SaveMessage(msg *Message) error
	GetMessages(limit int, before time.Time) ([]*Message, error)
	DeleteOldMessages(before time.Time) error
	AddParticipant(userID int) error
	RemoveParticipant(userID int) error
	IsParticipant(userID int) (bool, error)
}

// Service defines the interface for chat business logic
type Service interface {
	SendMessage(userID int, content string) error
	GetMessages(limit int) ([]*Message, error)
	DeleteOldMessages(maxAge time.Duration) error
	JoinChat(userID int) error
	LeaveChat(userID int) error
	IsParticipant(userID int) (bool, error)
}

// WebsocketMessage represents a message sent over websocket
type WebsocketMessage struct {
	Type    string          `json:"type"`
	Payload map[string]any `json:"payload"`
} 