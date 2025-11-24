package service

import (
	"AvitoTest/internal/model"
	"AvitoTest/internal/storage"
	"context"
)

type userService struct {
	userRepo storage.UserRepository
}

type UserService interface {
	SetIsActivate(ctx context.Context, userId string, isActivate bool) (model.User, error)
}

func NewUserService(userRep storage.UserRepository) UserService {
	return &userService{userRepo: userRep}
}

func (s *userService) SetIsActivate(ctx context.Context, userId string, isActivate bool) (model.User, error) {
	return s.userRepo.SetIsActive(ctx, userId, isActivate)
}
