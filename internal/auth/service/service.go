package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/chizheg/forum/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo domain.Repository
}

// NewService creates a new auth service
func NewService(repo domain.Repository) domain.Service {
	return &service{repo: repo}
}

func (s *service) Register(username, email, password string) (string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return "", err
	}

	// Create session
	token, err := s.generateToken()
	if err != nil {
		return "", err
	}

	session := &domain.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) Login(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Create session
	token, err := s.generateToken()
	if err != nil {
		return "", err
	}

	session := &domain.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) ValidateToken(token string) (int, error) {
	session, err := s.repo.GetSessionByToken(token)
	if err != nil {
		return 0, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.repo.DeleteSession(token)
		return 0, errors.New("session expired")
	}

	return session.UserID, nil
}

func (s *service) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
