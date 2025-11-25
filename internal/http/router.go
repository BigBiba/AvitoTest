package http

import (
	"AvitoTest/internal/http/handlers/pull_request"
	"AvitoTest/internal/http/handlers/team"
	"AvitoTest/internal/http/handlers/users"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(server *Server) http.Handler {
	r := mux.NewRouter()
	teamHandler := team.NewTeamHandler(server.teamService)
	r.HandleFunc(fmt.Sprintf("/%s/add", TeamGroup), teamHandler.CreateTeam).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/get", TeamGroup), teamHandler.GetTeam).Methods(http.MethodGet)

	userHandler := users.NewUserHandler(server.userService)
	r.HandleFunc(fmt.Sprintf("/%s/setIsActive", UsersGroup), userHandler.SetIsActive).Methods(http.MethodPost)

	pullReqHandler := pull_request.NewPullRequestHandler(server.pullRequestService)
	r.HandleFunc(fmt.Sprintf("/%s/getReview", UsersGroup), pullReqHandler.GetReviewerPullRequests).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/create", PullRequestGroup), pullReqHandler.CreatePR).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/merge", PullRequestGroup), pullReqHandler.SetMerged).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/reassign", PullRequestGroup), pullReqHandler.Reassign).Methods(http.MethodPost)

	return r
}
