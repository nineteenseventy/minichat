package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/http/middleware"
	"github.com/nineteenseventy/minichat/server/api"
	serverMiddleware "github.com/nineteenseventy/minichat/server/http/middleware"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func getRoutes() []func() (string, chi.Router) {
	return [](func() (string, chi.Router)){
		api.UserRouter,
		api.ChannelRouter,
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
	// for _, route := range getRoutes() {
	// 	prefix, router := route()
	// 	fmt.Println(prefix, router)
	// 	router.Mount(prefix, router)
	// }
	e, lmao := api.UserRouter()
	fmt.Println(e)
	router.Mount(e, lmao)
	e, lmao = api.ChannelRouter()
	fmt.Println(e)
	router.Mount(e, lmao)
	return router
}
