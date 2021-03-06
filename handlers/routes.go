package handlers

import "net/http"

//Route is the structure for an http route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//getRoutes maps url patterns to handlers
func getRoutes(h *Handler) []Route {
	return []Route{
		Route{
			"GetMessages",
			"GET",
			"/messages",
			h.GetMessages,
		},
		Route{
			"CreateMessages",
			"POST",
			"/messages",
			h.CreateMessage,
		},
		Route{
			"DeleteMessage",
			"DELETE",
			"/messages/{id}",
			h.DeleteMessage,
		},
		Route{
			"GetMessage",
			"GET",
			"/messages/{id}",
			h.GetMessage,
		},
	}
}
