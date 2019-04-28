package server

import (
	"../p3"
	"github.com/gorilla/mux"
	"net/http"
)

func NewAppServerRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range appServerRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = p3.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
