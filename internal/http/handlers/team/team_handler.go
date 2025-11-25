package team

import (
	model2 "AvitoTest/internal/domain/model"
	"AvitoTest/internal/helper"
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
		if errors.Is(err, model2.ErrTeamNotFound) {
			helper.WriteJSON(w, http.StatusNotFound,
				model2.ErrorResponse{
					Error: model2.ErrorDetail{
						Code:    "NOT_FOUND",
						Message: "resource not found"}})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, model2.InternalError)
		return
	}
	helper.WriteJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//получить команду, отдекодить, отправить в сервис. Обработать ошибочки
	var team model2.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, model2.ErrorResponse{
			Error: model2.ErrorDetail{
				Code:    "INVALID_PARAMS",
				Message: "invalid request",
			},
		})
		return
	}
	err := h.svc.CreateTeam(ctx, team)
	if err != nil {
		if errors.Is(err, model2.ErrTeamAlreadyExists) {
			helper.WriteJSON(w, http.StatusBadRequest, model2.ErrorResponse{
				Error: model2.ErrorDetail{
					Code:    "TEAM_EXISTS",
					Message: fmt.Sprintf("%s already exists", team.TeamName),
				},
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, model2.InternalError)
		return
	}
	helper.WriteJSON(w, http.StatusCreated, map[string]any{"team": team})
}
