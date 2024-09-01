package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"users": []}`))
}

func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/users", usersHandler)
	return r
}
