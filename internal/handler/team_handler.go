package handler

import (
	"AvitoTest/internal/helper"
	"AvitoTest/internal/model"
	"AvitoTest/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TeamHandler struct {
	svc service.TeamService
}

func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamName := r.URL.Query().Get("team_name")
	team, err := h.svc.GetTeam(ctx, teamName)

	if err != nil {
		if errors.Is(err, model.ErrTeamNotFound) {
			helper.WriteJSON(w, http.StatusNotFound,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    "NOT_FOUND",
						Message: "resource not found"}})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: "internal server error"}})
		return
	}
	helper.WriteJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//получить команду, отдекодить, отправить в сервис. Обработать ошибочки
	var team model.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_PARAMS",
				Message: "invalid request",
			},
		})
		return
	}
	err := h.svc.CreateTeam(ctx, team)
	if err != nil {
		if errors.Is(err, model.ErrTeamAlreadyExists) {
			helper.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: ErrorDetail{
					Code:    "TEAM_EXISTS",
					Message: fmt.Sprintf("%s already exists", team.TeamName),
				},
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}
	helper.WriteJSON(w, http.StatusCreated, map[string]any{"team": team})
}
