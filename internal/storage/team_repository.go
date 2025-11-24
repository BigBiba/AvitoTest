package storage

import (
	"AvitoTest/internal/model"
	"context"
	"database/sql"
	"errors"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team model.Team) error
	GetTeam(ctx context.Context, teamName string) (team model.Team, err error)
}

type PostgresTeamRepository struct {
	db *sql.DB
}

func NewPostgresTeamRepository(db *sql.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) GetTeam(ctx context.Context, teamName string) (model.Team, error) {
	query := `SELECT u.user_id, u.username, u.is_active 
				FROM teams t JOIN users u ON t.user_id = u.user_id 
				WHERE team_name = $1;`
	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Team{}, model.ErrTeamNotFound
		}
		return model.Team{}, err
	}
	defer rows.Close()
	members := make([]model.TeamMember, 0)
	for rows.Next() {
		var member model.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return model.Team{}, err
		}
		members = append(members, member)
	}
	team := model.Team{
		TeamName: teamName,
		Members:  members,
	}
	return team, nil
}

func (r *PostgresTeamRepository) CreateTeam(ctx context.Context, team model.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, member := range team.Members {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO teams (team_name, user_id) VALUES ($1, $2);`,
			team.TeamName, member.UserID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO users (user_id, username, is_active) VALUES ($1, $2, $3);`,
			member.UserID, member.Username, member.IsActive); err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
