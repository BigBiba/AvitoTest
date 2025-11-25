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
	r.HandleFunc("/team/add").Methods(http.MethodPost)
	r.HandleFunc("/team/get").Methods(http.MethodGet)

	r.HandleFunc("/users/setIsActive").Methods(http.MethodPost)
	r.HandleFunc("/users/getReview").Methods(http.MethodGet)

	r.HandleFunc("/pullRequests/create").Methods(http.MethodPost)
	r.HandleFunc("/pullRequests/merge").Methods(http.MethodPost)
	r.HandleFunc("/pullRequests/reassign").Methods(http.MethodPost)

}
