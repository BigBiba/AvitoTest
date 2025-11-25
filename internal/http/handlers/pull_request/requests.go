package pull_request

type PRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type ReassignRequest struct {
	prID          string `json:"pull_request_id"`
	oldReviewerID string `json:"old_reviewer_id"`
}
