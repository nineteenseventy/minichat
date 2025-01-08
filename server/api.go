package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/http/middleware"
	"github.com/nineteenseventy/minichat/server/api"
	serverMiddleware "github.com/nineteenseventy/minichat/server/http/middleware"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func getRoutes() []func(r chi.Router) {
	return []func(r chi.Router){
		api.UsersRouter,
		api.ChannelsRouter,
		api.MessagesRouter,
	}
}

func getMiddleware() []func(http.Handler) http.Handler {
	args := serverutil.GetArgs()
	return []func(http.Handler) http.Handler{
		middleware.LoggerMiddlewareFactory(),
		middleware.AuthenticationMiddlewareFactory(middleware.AuthenticationMiddlewareOptions{
			Domain:   args.Auth0Domain,
			Audience: args.Auth0Audience,
		}),
		serverMiddleware.UserMiddlewareFactory(),
	}
}

func ApiRouter() chi.Router {
	router := chi.NewRouter()
	for _, middleware := range getMiddleware() {
		router.Use(middleware)
	}
	for _, route := range getRoutes() {
		router.Group(route)
	}
	return router
}
