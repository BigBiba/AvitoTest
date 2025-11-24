package service

import (
	"AvitoTest/internal/model"
	"AvitoTest/internal/storage"
	"context"
	"errors"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team model.Team) error
	GetTeam(ctx context.Context, teamName string) (model.Team, error)
}

type teamService struct {
	teamRepo storage.TeamRepository
}

func NewTeamService(teamRep storage.TeamRepository) TeamService {
	return &teamService{teamRepo: teamRep}
}

func (s *teamService) CreateTeam(ctx context.Context, team model.Team) error {
	_, err := s.teamRepo.GetTeam(ctx, team.TeamName)
	if err == nil {
		return model.ErrTeamAlreadyExists
	}
	if !errors.Is(err, model.ErrTeamNotFound) {
		return err
	}
	err = s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return err
	}
	return nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (model.Team, error) {
	team, err := s.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		return model.Team{}, err
	}
	return team, nil
}
