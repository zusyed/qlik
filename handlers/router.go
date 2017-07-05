package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

//NewRouter returns a new router object for the specified routes.
func NewRouter(h *Handler) *mux.Router {
	routes := getRoutes(h)
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
