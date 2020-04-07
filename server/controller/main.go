package controller

import (
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/util"
	"net/http"
)

type Endpoint struct {
	Path        string
	Method      string
	Execute     func(w http.ResponseWriter, r *http.Request) error
	Middlewares []middleware.Middleware
}

var Endpoints = map[string]*Endpoint{
	getEndpointKey(hook):                     hook,
	getEndpointKey(getStandup):               getStandup,
	getEndpointKey(saveStandup):              saveStandup,
	getEndpointKey(getConfig):                getConfig,
	getEndpointKey(setConfig):                setConfig,
	getEndpointKey(getDefaultTimezone):       getDefaultTimezone,
	getEndpointKey(getPluginConfig):          getPluginConfig,
	getEndpointKey(getActiveStandupChannels): getActiveStandupChannels,
}

func getEndpointKey(endpoint *Endpoint) string {
	return util.GetKeyHash(endpoint.Path + "-" + endpoint.Method)
}

func GetEndpoint(r *http.Request) *Endpoint {
	return Endpoints[util.GetKeyHash(r.URL.Path+"-"+r.Method)]
}
