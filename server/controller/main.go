package controller

import (
	"net/http"

	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/util"
)

type endpointHandler func(w http.ResponseWriter, r *http.Request) error

type authenticatedEndpointHandler func(userID string, w http.ResponseWriter, r *http.Request) error

type Endpoint struct {
	Path        string
	Method      string
	Execute     endpointHandler
	Middlewares []middleware.Middleware
}

var Endpoints = map[string]*Endpoint{
	getEndpointKey(hook):                     hook,
	getEndpointKey(getStandup):               getStandup,
	getEndpointKey(saveStandup):              saveStandup,
	getEndpointKey(getConfig):                getConfig,
	getEndpointKey(setConfig):                setConfig,
	getEndpointKey(getDefaultTimezone):       getDefaultTimezone,
	getEndpointKey(getActiveStandupChannels): getActiveStandupChannels,
	getEndpointKey(getPluginConfig):          getPluginConfig,
}

func getEndpointKey(endpoint *Endpoint) string {
	return util.GetKeyHash(endpoint.Path + "-" + endpoint.Method)
}

func GetEndpoint(r *http.Request) *Endpoint {
	return Endpoints[util.GetKeyHash(r.URL.Path+"-"+r.Method)]
}

func authenticatedControllerWrapper(handler authenticatedEndpointHandler) endpointHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userID := r.Context().Value(middleware.CtxKeyUserID).(string)
		return handler(userID, w, r)
	}
}
