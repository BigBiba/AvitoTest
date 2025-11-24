package model

import "errors"

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamAlreadyExists = errors.New("team already exists")
)
