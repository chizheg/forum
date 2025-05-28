package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chizheg/forum/internal/auth/domain"
)

type repository struct {
	db *sql.DB
}

// NewRepository creates a new PostgreSQL repository
func NewRepository(db *sql.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *repository) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (r *repository) GetUserByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (r *repository) CreateSession(session *domain.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		session.UserID,
		session.Token,
		session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	return nil
}

func (r *repository) GetSessionByToken(token string) (*domain.Session, error) {
	session := &domain.Session{}
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = $1`

	err := r.db.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return session, nil
}

func (r *repository) DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	result, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}
