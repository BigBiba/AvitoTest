package service

import (
	"AvitoTest/internal/model"
	"AvitoTest/internal/storage"
	"context"
)

type userService struct {
	userRepo storage.UserRepository
	teamRepo storage.TeamRepository
}

type UserService interface {
	SetIsActivate(ctx context.Context, userId string, isActivate bool) (model.User, error)
}

func (s *userService) SetIsActivate(ctx context.Context, userId string, isActivate bool) (model.User, error) {
	
}
