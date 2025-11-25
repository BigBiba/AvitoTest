package model

import "errors"

var (
	ErrTeamNotFound             = errors.New("team not found")
	ErrTeamAlreadyExists        = errors.New("team already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrPullRequestAlreadyExists = errors.New("pull request already exists")
	ErrPullRequestNotFound      = errors.New("pull request not found")
	ErrReassignAfterMerged      = errors.New("cannot reassign on merged PR")
	ErrUserIsNotReviewer        = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidates             = errors.New("no active replacement candidate in team")
)
