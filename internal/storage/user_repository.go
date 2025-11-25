package storage

import (
	model "AvitoTest/internal/domain/model"
	"context"
	"database/sql"
	"errors"
)

type UserRepository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (user model.User, err error)
	GetUser(ctx context.Context, userID string) (user model.User, err error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) (model.User, error) {
	query := `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, is_active
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, isActive, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}
	query = `
		SELECT team_name
		FROM teams
		WHERE user_id = $1
	`
	err = r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Team,
	)
	if err != nil {
		return model.User{}, model.ErrUserNotFound
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUser(ctx context.Context, userID string) (model.User, error) {
	query := `
SELECT user_id, username, is_active
FROM users
WHERE user_id = $1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}
	return user, nil
}
