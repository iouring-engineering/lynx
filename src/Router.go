package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct{}

func (router *Router) initializeRouter(baseRouter *mux.Router) {
	var subRouter = baseRouter.PathPrefix("/lynx").Subrouter()
	subRouter.HandleFunc("/create",
		BaseMW(CreateShortLink)).Methods(http.MethodPost)
}
