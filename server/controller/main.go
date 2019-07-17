package controller

import (
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/util"
	"net/http"
)

type Endpoint struct {
	Path         string
	Method       string
	Execute      func(w http.ResponseWriter, r *http.Request) error
	RequiresAuth bool
}

var Endpoints = map[string]*Endpoint{
	getEndpointKey(hook):               hook,
	getEndpointKey(getStandup):         getStandup,
	getEndpointKey(saveStandup):        saveStandup,
	getEndpointKey(getConfig):          getConfig,
	getEndpointKey(setConfig):          setConfig,
	getEndpointKey(getDefaultTimezone): getDefaultTimezone,
}

func getEndpointKey(endpoint *Endpoint) string {
	return util.GetKeyHash(endpoint.Path + "-" + endpoint.Method)
}

func GetEndpoint(r *http.Request) *Endpoint {
	return Endpoints[util.GetKeyHash(r.URL.Path+"-"+r.Method)]
}

// verifies if provided request is performed by a logged in Mattermost user.
func Authenticated(w http.ResponseWriter, r *http.Request) bool {
	userId := r.Header.Get(config.HeaderMattermostUserId)

	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}
