package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/http/util"
)

type HealthMessage struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	msg := HealthMessage{Status: "ok"}
	util.JSONResponse(w, msg)
}

func HealthRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", healthHandler)
	return r
}
