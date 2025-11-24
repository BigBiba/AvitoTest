package storage

import (
	"AvitoTest/internal/model"
	"context"
	"database/sql"
	"errors"
)

type UserRepository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (user model.User, err error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
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
