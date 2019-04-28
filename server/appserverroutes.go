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
		"GetClientId",
		"GET",
		"/getmyid",
		GetClientId,
	},
	//Route{
	//	"Canonical",
	//	"GET",
	//	"/canonical",
	//	Canonical,
	//},
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
