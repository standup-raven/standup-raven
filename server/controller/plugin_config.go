package controller

import (
	"encoding/json"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/logger"
	"net/http"
)

var getPluginConfig = &Endpoint{
	Path:        "/plugin-config",
	Method:      http.MethodGet,
	Execute:     executeGetPluginConfig,
	Middlewares: []middleware.Middleware{
		//middleware.Authenticate,
	},
}

func executeGetPluginConfig(w http.ResponseWriter, r *http.Request) error {
	conf := config.GetConfig()

	pluginConfig := map[string]interface{}{
		"disableChannelHeaderButton": conf.DisableChannelHeaderButton,
	}

	data, err := json.Marshal(pluginConfig)
	if err != nil {
		logger.Error("Couldn't serialize plugin config data", err, map[string]interface{}{"pluginConfig": pluginConfig})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}

	return nil
}
