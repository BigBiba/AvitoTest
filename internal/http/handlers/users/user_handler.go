package users

import (
	"AvitoTest/internal/domain/model"
	"AvitoTest/internal/helper"
	http2 "AvitoTest/internal/http"
	"AvitoTest/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type UserSetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req UserSetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, http2.ErrorResponse{
			Error: http2.ErrorDetail{
				Code:    "INVALID_PARAMS",
				Message: "invalid request",
			},
		})
		return
	}
	user, err := h.svc.SetIsActivate(ctx, req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			helper.WriteJSON(w, http.StatusNotFound, http2.ErrorResponse{
				Error: http2.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, http2.InternalError)
		return
	}
	helper.WriteJSON(w, http.StatusOK, user)
}
