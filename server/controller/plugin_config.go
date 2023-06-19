package controller

import (
	"net/http"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/logger"
)

var getPluginConfig = &Endpoint{
	Path:    "/plugin-config",
	Method:  http.MethodGet,
	Execute: executeGetPluginConfig,
	Middlewares: []middleware.Middleware{
		middleware.Authenticated,
	},
}

func executeGetPluginConfig(w http.ResponseWriter, r *http.Request) error {
	conf := config.GetConfig().Sanitize()

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(conf.ToJSON()); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(conf.ToJSON())})
		return err
	}

	return nil
}
