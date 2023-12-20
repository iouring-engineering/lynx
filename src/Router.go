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
	baseRouter.HandleFunc("/create",
		BaseMW(CreateShortLink)).Methods(http.MethodPost)
	baseRouter.HandleFunc("/{shortcode}",
		BaseMW(GetSourceLink)).Methods(http.MethodGet)
	baseRouter.HandleFunc("/data/{shortcode}",
		BaseMW(GetData)).Methods(http.MethodGet)
}
