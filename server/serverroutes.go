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
		"GiveClientId",
		"GET",
		"/getmyid",
		GiveClientId,
	}, Route{
		"CreateRequest",
		"POST",
		"/trade/request/create",
		CreateRequest,
	},
	Route{
		"ViewRequest",
		"POST",
		"/trade/request/{id}",
		ViewRequest,
	},
	Route{
		"FulfillRequest",
		"POST",
		"/trade/fulfill/{id}",
		FulfillRequest,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	//Route{
	//	"Show",
	//	"GET",
	//	"/show",
	//	Show,
	//},
	//Route{
	//	"Upload",
	//	"POST",
	//	"/upload",
	//	Upload,
	//},
	//Route{
	//	"UploadBlock",
	//	"GET",
	//	"/block/{height}/{hash}",
	//	UploadBlock,
	//},
	//Route{
	//	"HeartBeatReceive",
	//	"POST",
	//	"/heartbeat/receive",
	//	HeartBeatReceive,
	//},
	//Route{
	//	"Start",
	//	"GET",
	//	"/start",
	//	Start,
	//},
}
