package api

import (
	"github.com/go-chi/chi/v5"
)

func ApiRouter() chi.Router {
	r := chi.NewRouter()
	r.Mount("/users", UserRouter())
	return r
}
