package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct{}

func (router *Router) initializeRouter(baseRouter *mux.Router) {
	baseRouter.HandleFunc("/.well-known/apple-app-site-association",
		BaseMW(IosVerify)).Methods(http.MethodGet)
	baseRouter.HandleFunc("/.well-known/assetlinks.json",
		BaseMW(AndroidVerify)).Methods(http.MethodGet)

	var subRouter = baseRouter.PathPrefix("/lynx").Subrouter()
	subRouter.HandleFunc("/create",
		BaseMW(CreateShortLink)).Methods(http.MethodPost)
	subRouter.HandleFunc("/{shorturl}",
		BaseMW(GetSourceLink)).Methods(http.MethodGet)
}
