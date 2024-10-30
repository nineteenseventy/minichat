package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/http/util"
)

type HealthMessage struct {
	Status string `json:"status"`
}

func healthHandler(writer http.ResponseWriter, request *http.Request) {
	msg := HealthMessage{Status: "ok"}
	util.JSONResponse(writer, msg)
}

func HealthRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/health", healthHandler)
	return router
}
