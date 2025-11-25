package storage

import (
	model "AvitoTest/internal/domain/model"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PullRequestRepository interface {
	GetReviewerPRs(ctx context.Context, user model.User) ([]model.PullRequestShort, error)
	GetPullRequest(ctx context.Context, prID string) (model.PullRequest, error)
	SetMerged(ctx context.Context, prID string) (*time.Time, error)
	ReassignReviewer(ctx context.Context, pr model.PullRequest, oldReviewerID string, newReviewerID string) error
	CreatePullRequest(ctx context.Context, prID string, prName string, authorID string, reviewers []string) (model.PullRequest, error)
	GetTeamMembers(ctx context.Context, userID string) ([]string, error)
	GetCandidates(ctx context.Context, prID string) ([]string, error)
}

type PostgresPullRequestRepository struct {
	db *sql.DB
}

func NewPostgresPullRequestRepository(db *sql.DB) PullRequestRepository {
	return &PostgresPullRequestRepository{db: db}
}

func (r *PostgresPullRequestRepository) GetCandidates(ctx context.Context, prID string) ([]string, error) {
	query := `
		SELECT user_id
		FROM users u 
		WHERE u.team_name = (
		    SELECT team_name
		    FROM teams t
		    WHERE t.user_id = (SELECT author_id FROM pull_requests pr WHERE pr.pr_id = $1)
		)
		AND u.user_id NOT IN (
		    SELECT user_id
		    FROM users_pull_requests
		    WHERE pr_id = $1
		)
`
	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *PostgresPullRequestRepository) GetTeamMembers(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT u.user_id, u.username, t.team_name, u.is_active
		FROM users u JOIN teams t ON u.user_id = t.user_id
		where users.user_id = $1
`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Team,
		&user.IsActive,
	)
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	query = `
		SELECT user_id
		FROM users u JOIN teams t ON u.user_id = t.user_id
		WHERE u.user_id <> $1 and t.team_name = $2
	`
	rows, err := r.db.QueryContext(ctx, query, user.UserID, user.Team)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []string
	for rows.Next() {
		var member string
		if err := rows.Scan(&member); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *PostgresPullRequestRepository) CreatePullRequest(
	ctx context.Context,
	prID string,
	prName string,
	authorID string,
	reviewers []string,
) (model.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.PullRequest{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	var pr model.PullRequest

	query := `
        INSERT INTO pull_requests (pr_id, name, author_id, status)
        VALUES ($1, $2, $3, 'OPEN')
        RETURNING pr_id, pr_name, author_id, status
    `
	err = r.db.QueryRowContext(ctx, query, prID, prName, authorID).Scan(
		&pr.PRId,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
	)

	if err != nil {
		return model.PullRequest{}, err
	}

	for reviewer := range reviewers {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO users_pull_requests (pr_id, user_id) VALUES ($1, $2);`,
			pr.PRId, reviewer); err != nil {
			return model.PullRequest{}, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return model.PullRequest{}, err
	}

	return pr, nil

}

func (r *PostgresPullRequestRepository) GetPullRequest(ctx context.Context, prID string) (model.PullRequest, error) {
	query := `
		SELECT pr_id, name, author_id, status
		FROM pull_requests WHERE pr_id = $1
`
	var pr model.PullRequest
	err := r.db.QueryRowContext(ctx, query, prID).Scan(
		&pr.PRId,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, model.ErrPullRequestNotFound
		}
		return model.PullRequest{}, model.ErrPullRequestNotFound
	}
	query = `SELECT user_id FROM users_pull_requests WHERE pr_id = $1`
	var reviewers []string
	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, model.ErrTeamNotFound
		}
	}
	defer rows.Close()
	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return model.PullRequest{}, err
		}
		reviewers = append(reviewers, reviewer)
	}
	if err := rows.Err(); err != nil {
		return model.PullRequest{}, err
	}
	pr.Reviewers = reviewers
	return pr, nil
}

func (r *PostgresPullRequestRepository) SetMerged(ctx context.Context, prID string) (*time.Time, error) {
	query := `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = NOW()
		WHERE pr_id = $1
		RETURNING merged_at
`
	var mergedAt time.Time
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&mergedAt)
	if err != nil {
		return nil, err
	}
	return &mergedAt, nil
}

func (r *PostgresPullRequestRepository) GetReviewerPRs(
	ctx context.Context,
	user model.User,
) ([]model.PullRequestShort, error) {
	query := `
SELECT pr.pr_id, pr.name, pr.author_id, pr.status
FROM users_pull_requests upr JOIN pull_requests pr ON pr.pr_id == upr.pr_id
WHERE upr.user_id = $1
`
	rows, err := r.db.QueryContext(ctx, query, user.UserID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	defer rows.Close()
	var prs []model.PullRequestShort
	for rows.Next() {
		var pr model.PullRequestShort
		if err := rows.Scan(&pr); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func (r *PostgresPullRequestRepository) ReassignReviewer(
	ctx context.Context,
	pr model.PullRequest,
	oldReviewerID string,
	newReviewerID string,
) error {
	query := `
		UPDATE users_pull_requests
		SET user_id = $1
		WHERE pr_id = $2 AND user_id = $3
`
	_, err := r.db.ExecContext(ctx, query, newReviewerID, pr.PRId, oldReviewerID)
	if err != nil {
		return err
	}

	return nil
}
