package model

import "errors"

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrUserNotFound      = errors.New("user not found")
)
