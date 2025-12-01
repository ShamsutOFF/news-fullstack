package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"news-fullstack/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists   = errors.New("user with this email already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UsersRepository struct {
	dbpool *pgxpool.Pool
	logger *slog.Logger
}

func NewUsersRepository(dbpool *pgxpool.Pool, logger *slog.Logger) *UsersRepository {
	return &UsersRepository{
		dbpool: dbpool,
		logger: logger,
	}
}

// CreateUser создает нового пользователя
func (r *UsersRepository) CreateUser(email, name, passwordHash string) error {
	const query = `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int64
	err := r.dbpool.QueryRow(context.Background(), query, email, name, passwordHash).Scan(&id)

	if err != nil {
		// Проверяем на дублирование email
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "already exists") {
			return ErrUserExists
		}

		r.logger.Error("failed to create user", "error", err, "email", email)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("user created", "id", id, "email", email)
	return nil
}

// GetUserByEmail возвращает пользователя по email
func (r *UsersRepository) GetUserByEmail(email string) (*models.User, error) {
	const query = `
		SELECT id, email, name, password_hash, created_at, updated_at,
			   is_active, email_verified
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.dbpool.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
		&user.EmailVerified,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		r.logger.Error("failed to get user", "error", err, "email", email)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Простая проверка существования пользователя
func (r *UsersRepository) UserExists(email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.dbpool.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		r.logger.Error("failed to check user existence", "error", err, "email", email)
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
