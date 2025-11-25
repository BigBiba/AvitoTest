package model

import "time"

type PRStatus int

const (
	OPEN PRStatus = iota
	MERGED
)

type PullRequest struct {
	PRId      string     `json:"pull_request_id"`
	Name      string     `json:"pull_request_name"`
	AuthorID  string     `json:"author_id"`
	Status    PRStatus   `json:"status"`
	Reviewers []string   `json:"assigned_reviewers"`
	MergedAt  *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PRId     string   `json:"pull_request_id"`
	Name     string   `json:"pull_request_name"`
	AuthorId string   `json:"author_id"`
	Status   PRStatus `json:"status"`
}
