package service

import (
	"AvitoTest/internal/model"
	"AvitoTest/internal/storage"
	"context"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, prID string, prName string, authorID string) (model.PullRequest, error)
	SetMerged(ctx context.Context, prID string) (model.PullRequest, error)
	ReassignReviewers(ctx context.Context, prID string, oldReviewerID string) (model.PullRequest, error)
	GetReviewerPullRequests(ctx context.Context, prID string) ([]model.PullRequest, error)
}

type pullRequestService struct {
	prRepo storage.PullRequestRepository
}
