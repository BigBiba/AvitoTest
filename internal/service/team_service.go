package service

import (
	"AvitoTest/internal/model"
	"AvitoTest/internal/storage"
	"context"
)

type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, members []model.TeamMember) (model.Team, error)
	GetTeam(ctx context.Context, teamName string) (model.Team, error)
}

type teamService struct {
	teamRepo storage.TeamRepository
	userRepo storage.UserRepository
}

func (s *teamService) CreateTeam(ctx context.Context, teamName string, members []model.TeamMember) (model.Team, error) {

}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (model.Team, error) {
	//пойти в репу поискать команду с таким названием, если такой команды нет, то ошибка, команда не найдена
	team, err := s.teamRepo.GetTeam(ctx, teamName)
}
