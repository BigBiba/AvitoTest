package service

import (
	"AvitoTest/internal/domain/model"
	"AvitoTest/internal/helper"
	"AvitoTest/internal/storage"
	"context"
	"database/sql"
	"errors"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, prID string, prName string, authorID string) (model.PullRequest, error)
	SetMerged(ctx context.Context, prID string) (model.PullRequest, error)
	ReassignReviewers(ctx context.Context, prID string, oldReviewerID string) (model.PullRequest, string, error)
	GetReviewerPullRequests(ctx context.Context, userID string) ([]model.PullRequestShort, error)
}

type pullRequestService struct {
	prRepo   storage.PullRequestRepository
	userRepo storage.UserRepository
}

func NewPullRequestService(prRepo storage.PullRequestRepository, userRepo storage.UserRepository) PullRequestService {
	return &pullRequestService{prRepo: prRepo, userRepo: userRepo}
}

func (s *pullRequestService) CreatePullRequest(
	ctx context.Context,
	prID string,
	prName string,
	authorID string,
) (model.PullRequest, error) {
	reviewers, err := s.prRepo.GetTeamMembers(ctx, authorID)
	if err != nil {
		return model.PullRequest{}, err
	}
	return s.prRepo.CreatePullRequest(ctx, prID, prName, authorID, reviewers)
}

func (s *pullRequestService) SetMerged(ctx context.Context, prID string) (model.PullRequest, error) {
	pr, err := s.prRepo.GetPullRequest(ctx, prID)
	if err != nil {
		return model.PullRequest{}, err
	}
	timeMerged, err := s.prRepo.SetMerged(ctx, prID)
	if err != nil {
		return model.PullRequest{}, err
	}
	pr.MergedAt = timeMerged
	pr.Status = model.MERGED
	return pr, nil
}

func (s *pullRequestService) ReassignReviewers(
	ctx context.Context,
	prID string,
	oldReviewerID string,
) (model.PullRequest, string, error) {
	oldReviewer, err := s.userRepo.GetUser(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, "", model.ErrUserNotFound
		}
		return model.PullRequest{}, "", err
	}
	pr, err := s.prRepo.GetPullRequest(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, "", model.ErrPullRequestNotFound
		}
		return model.PullRequest{}, "", err
	}
	if pr.Status == model.MERGED {
		return model.PullRequest{}, "", model.ErrReassignAfterMerged
	}
	idx := -1
	for i, reviewerId := range pr.Reviewers {
		if reviewerId == oldReviewer.UserID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return model.PullRequest{}, "", model.ErrUserIsNotReviewer
	}
	candidates, err := s.prRepo.GetCandidates(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, "", model.ErrNoCandidates
		}
		return model.PullRequest{}, "", err
	}
	candidate := helper.GetRandomElement(candidates)
	err = s.prRepo.ReassignReviewer(ctx, pr, oldReviewer.UserID, candidate)
	if err != nil {
		return model.PullRequest{}, "", err
	}
	pr.Reviewers[idx] = candidate
	return pr, candidate, nil
}

func (s *pullRequestService) GetReviewerPullRequests(ctx context.Context, userID string) ([]model.PullRequestShort, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.prRepo.GetReviewerPRs(ctx, user)
}
