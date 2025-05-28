package service

import (
	"testing"
	"time"

	"github.com/chizheg/forum/internal/auth/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock implementation of domain.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(id int) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) CreateSession(session *domain.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockRepository) GetSessionByToken(token string) (*domain.Session, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockRepository) DeleteSession(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestService_Register(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// Test successful registration
	mockRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil)
	mockRepo.On("CreateSession", mock.AnythingOfType("*domain.Session")).Return(nil)

	token, err := svc.Register("testuser", "test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestService_Login(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// Test successful login
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &domain.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)
	mockRepo.On("CreateSession", mock.AnythingOfType("*domain.Session")).Return(nil)

	token, err := svc.Login("testuser", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test invalid password
	token, err = svc.Login("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestService_ValidateToken(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	// Test valid token
	validSession := &domain.Session{
		ID:        1,
		UserID:    1,
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRepo.On("GetSessionByToken", "valid-token").Return(validSession, nil)
	userID, err := svc.ValidateToken("valid-token")
	assert.NoError(t, err)
	assert.Equal(t, 1, userID)

	// Test expired token
	expiredSession := &domain.Session{
		ID:        2,
		UserID:    2,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	mockRepo.On("GetSessionByToken", "expired-token").Return(expiredSession, nil)
	mockRepo.On("DeleteSession", "expired-token").Return(nil)
	userID, err = svc.ValidateToken("expired-token")
	assert.Error(t, err)
	assert.Equal(t, 0, userID)

	mockRepo.AssertExpectations(t)
}
