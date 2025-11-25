package pull_request

import (
	"AvitoTest/internal/domain/model"
	"AvitoTest/internal/helper"
	"AvitoTest/internal/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type PullRequestHandler struct {
	svc service.PullRequestService
}

func NewPullRequestHandler(svc service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{svc: svc}
}

func (h *PullRequestHandler) GetReviewerPullRequests(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := r.URL.Query().Get("user_id")
	prs, err := h.svc.GetReviewerPullRequests(ctx, userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, model.InternalError)
	}
	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}

func (h *PullRequestHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req PRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "INVALID_PARAMS",
				Message: "invalid request body",
			},
		})
		return
	}
	pr, err := h.svc.CreatePullRequest(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) || errors.Is(err, model.ErrPullRequestNotFound) {
			helper.WriteJSON(w, http.StatusNotFound, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
			return
		}
		if errors.Is(err, model.ErrPullRequestAlreadyExists) {
			helper.WriteJSON(w, http.StatusConflict, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
				},
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, model.InternalError)
		return
	}
	helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

func (h *PullRequestHandler) SetMerged(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req struct {
		prID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "INVALID_PARAMS",
				Message: "invalid request body",
			},
		})
	}
	pr, err := h.svc.SetMerged(ctx, req.prID)
	if err != nil {
		if errors.Is(err, model.ErrPullRequestNotFound) {
			helper.WriteJSON(w, http.StatusNotFound, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, model.InternalError)
		return
	}
	helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

func (h *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "INVALID_PARAMS",
				Message: "invalid request body",
			},
		})
		return
	}
	pr, newReviewerId, err := h.svc.ReassignReviewers(ctx, req.prID, req.oldReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound) || errors.Is(err, model.ErrPullRequestNotFound):
			helper.WriteJSON(w, http.StatusNotFound, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
			return
		case errors.Is(err, model.ErrReassignAfterMerged):
			helper.WriteJSON(w, http.StatusConflict, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR",
				},
			})
			return
		case errors.Is(err, model.ErrUserIsNotReviewer):
			helper.WriteJSON(w, http.StatusConflict, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "NOT_ASSIGNED",
					Message: "reviewer is not assigned to this PR",
				},
			})
			return
		case errors.Is(err, model.ErrNoCandidates):
			helper.WriteJSON(w, http.StatusConflict, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "NO_CANDIDATES",
					Message: "no active replacement candidate in team",
				},
			})
			return
		default:
			helper.WriteJSON(w, http.StatusInternalServerError, model.InternalError)
			return
		}
	}
	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": newReviewerId,
	})
}
