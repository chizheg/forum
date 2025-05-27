package domain

import "time"

// User represents the user entity
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session represents a user session
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository defines the interface for user data access
type Repository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateSession(session *Session) error
	GetSessionByToken(token string) (*Session, error)
	DeleteSession(token string) error
}

// Service defines the interface for user business logic
type Service interface {
	Register(username, email, password string) (string, error)
	Login(username, password string) (string, error)
	ValidateToken(token string) (int, error)
} 