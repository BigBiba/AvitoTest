package model

type PRStatus int

const (
	OPEN PRStatus = iota
	MERGED
)

type PullRequest struct {
	PRId      string
	Name      string
	AuthorID  string
	Status    PRStatus
	Reviewers []string
}
