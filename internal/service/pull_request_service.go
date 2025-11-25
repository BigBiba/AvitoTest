package service

import (
	"AvitoTest/internal/model"
	"AvitoTest/internal/storage"
	"context"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, prID string, prName string, authorID string) (model.PullRequest, error)
	SetMerged(ctx context.Context, prID string) (model.PullRequest, error)
	ReassignReviewers(ctx context.Context, prID string, oldReviewerID string) (model.PullRequest, string, error)
	GetReviewerPullRequests(ctx context.Context, prID string) ([]model.PullRequest, error)
}

type pullRequestService struct {
	prRepo storage.PullRequestRepository
}

func NewPullRequestService(prRepo storage.PullRequestRepository) PullRequestService {
	return &pullRequestService{prRepo: prRepo}
}

func (s *pullRequestService) CreatePullRequest(
	ctx context.Context,
	prID string,
	prName string,
	authorID string,
) (model.PullRequest, error) {
	
}

func (s *pullRequestService) SetMerged(ctx context.Context, prID string) (model.PullRequest, error) {

}

func (s *pullRequestService) ReassignReviewers(
	ctx context.Context,
	prID string,
	oldReviewerID string,
) (model.PullRequest, string, error) {

}

func (s *pullRequestService) GetReviewerPullRequests(ctx context.Context, prID string) ([]model.PullRequest, error) {

}
