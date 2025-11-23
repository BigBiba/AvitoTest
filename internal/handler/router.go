package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type TeamRequest struct {
}

func AddTeam(w http.ResponseWriter, r *http.Request) {

}

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/team/add")
	mux.Handle("/team/get")
	mux.HandleFunc()

	mux.Handle("/users/setIsActive")
	mux.Handle("/users/getReview")

	mux.Handle("/pullRequests/create")
	mux.Handle("/pullRequests/merge")
	mux.Handle("/pullRequests/reassign")

}
