package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/server/api"
	"github.com/nineteenseventy/minichat/server/middleware"
)

func getRoutes() []func() chi.Router {
	return [](func() chi.Router){
		api.UserRouter,
	}
}

func getMiddleware() []func(http.Handler) http.Handler {
	args := GetArgs()
	return []func(http.Handler) http.Handler{
		middleware.LoggerMiddleware(),
		middleware.AuthenticationMiddleware(middleware.AuthenticationMiddlewareOptions{
			Domain:   args.Auth0Domain,
			Audience: args.Auth0Audience,
		}),
	}
}

func ApiRouter() chi.Router {
	r := chi.NewRouter()
	for _, middleware := range getMiddleware() {
		r.Use(middleware)
	}
	for _, route := range getRoutes() {
		r.Mount("/", route())
	}
	return r
}
