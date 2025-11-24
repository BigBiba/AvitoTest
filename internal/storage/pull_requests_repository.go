package storage

import (
	"AvitoTest/internal/model"
	"context"
)

type PullRequestRepository interface {
	CreatePullRequest(ctx context.Context, prID string, prName string, authorID string) (pr PullRequest, err error)
	SetMerged(ctx context.Context, prID string) (pr model.PullRequest, err error)
	ReassignReviewer(ctx context.Context, prID string) (pr model.PullRequest, err error)
	GetReview(ctx context.Context, userID string) (prs []model.PullRequest, err error)
}
