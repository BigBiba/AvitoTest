package http

import (
	"AvitoTest/internal/service"
	"AvitoTest/internal/storage"
	"database/sql"
)

type Server struct {
	teamService        service.TeamService
	userService        service.UserService
	pullRequestService service.PullRequestService
}

func NewServer(db *sql.DB) *Server {
	teamRepo := storage.NewPostgresTeamRepository(db)
	userRepo := storage.NewPostgresUserRepository(db)
	pullRequestRepo := storage.NewPostgresPullRequestRepository(db)

	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	pullRequestService := service.NewPullRequestService(pullRequestRepo, userRepo)

	return &Server{
		teamService:        teamService,
		userService:        userService,
		pullRequestService: pullRequestService,
	}
}
