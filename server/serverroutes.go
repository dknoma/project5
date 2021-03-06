package server

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var appServerRoutes = Routes{
	Route{
		"CreateAccount",
		"GET",
		"/account/create",
		CreateAccount,
	},
	Route{
		"CreateRequest", // Allow sellers to create trade requests
		"POST",          // POST seller id, equipment json, demand json
		"/trade/request/create",
		CreateRequest,
	},
	Route{
		"ViewRequest", // Allow players to view info on a trade request
		"GET",         // GET to view info on a trade request
		"/trade/request/{id}",
		ViewRequest,
	},
	Route{
		"FulfillRequest",      // Allow buyers to potentially fulfill trade requests
		"POST",                // POST buyer id to check their currency, if has enough then can fulfill request
		"/trade/fulfill/{id}", // id is request id
		FulfillRequest,        // NOTE: monolithic=same service, microservice=fulfillment service
	},
	Route{
		"GetPendingTransactions",
		"GET",
		"/trade/fulfillments/peek",
		GetPendingTransactions,
	},
	Route{
		"UpdatePendingTransactions",
		"POST",
		"/trade/fulfillments/update",
		UpdatePendingTransactions,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
}
