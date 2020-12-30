package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

var getConfig = &Endpoint{
	Path:    "/config",
	Method:  http.MethodGet,
	Execute: authenticatedControllerWrapper(executeGetConfig),
	Middlewares: []middleware.Middleware{
		middleware.Authenticated,
	},
}

var setConfig = &Endpoint{
	Path:    "/config",
	Method:  http.MethodPost,
	Execute: authenticatedControllerWrapper(executeSetConfig),
	Middlewares: []middleware.Middleware{
		middleware.Authenticated,
		middleware.SetUserRoles,
		middleware.DisallowGuests,
		middleware.HandlePermissionSchema,
	},
}

var getDefaultTimezone = &Endpoint{
	Path:    "/timezone",
	Method:  http.MethodGet,
	Execute: executeGetDefaultTimezone,
}

var getActiveStandupChannels = &Endpoint{
	Path:    "/active-channels",
	Method:  http.MethodGet,
	Execute: executeGetActiveStandupChannels,
	Middlewares: []middleware.Middleware{
		middleware.Authenticated,
	},
}

func executeGetConfig(userID string, w http.ResponseWriter, r *http.Request) error {
	channelID := r.URL.Query().Get("channel_id")

	c, err := standup.GetStandupConfig(channelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if c == nil {
		http.Error(w, "Standup not configured for this channel", http.StatusNotFound)
		return nil
	}

	// TODO: make use of ToJSON function for sending conf in response
	data, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Couldn't serialize config data", err, map[string]interface{}{"config": c.ToJSON()})
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}

	return nil
}

func executeSetConfig(userID string, w http.ResponseWriter, r *http.Request) error {
	// get config data from body
	decoder := json.NewDecoder(r.Body)
	conf := &standup.Config{}
	if err := decoder.Decode(&conf); err != nil {
		logger.Error("Could not decode request body", err, map[string]interface{}{"request": util.DumpRequest(r)})
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return err
	}

	channelID := conf.ChannelID
	channelIDParam := r.URL.Query().Get("channel_id")

	if channelID != channelIDParam {
		http.Error(w, "Mismatched channel ID", http.StatusBadRequest)
		return errors.New("channel ID provided in config body does not match with the value in query params")
	}

	if err := conf.PreSave(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := conf.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	json, err := json.Marshal(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(json))
	}

	conf, err = standup.SaveStandupConfig(conf)
	if err != nil {
		http.Error(w, "Error occurred while saving standup conf", http.StatusInternalServerError)
		return err
	}

	// if `SaveStandupConfig` succeed and this failed-
	// 	1. Standup Raven channel header button won't show up ever in
	//			this channel, even if standup is configured in the channel.
	// 	2. Scheduled standup reports won't work for this channel.
	if err := standup.AddStandupChannel(conf.ChannelID); err != nil {
		http.Error(w, "Error occurred while saving standup conf", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(conf.ToJSON())); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"config": conf.ToJSON()})
		return err
	}

	event := ""
	if conf.Enabled {
		event = "add_active_channel"
	} else {
		event = "remove_active_channel"
	}

	config.Mattermost.PublishWebSocketEvent(
		event,
		map[string]interface{}{
			"channel_id": conf.ChannelID,
		},
		&model.WebsocketBroadcast{
			UserId: userID,
		},
	)

	return nil
}

func executeGetDefaultTimezone(w http.ResponseWriter, r *http.Request) error {
	timezone := config.GetConfig().TimeZone

	data, err := json.Marshal(timezone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Couldn't serialize config data", err, map[string]interface{}{"timezone": timezone})
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}

	return nil
}

func executeGetActiveStandupChannels(w http.ResponseWriter, r *http.Request) error {
	standupChannels, err := standup.GetStandupChannels()
	if err != nil {
		logger.Error("An error occurred while fetching standup channels", err, nil)
		http.Error(w, "An error occurred while fetching standup channels", http.StatusInternalServerError)
		return err
	}

	activeStandupChannels := []string{}

	for _, channelID := range standupChannels {
		standupConfig, err := standup.GetStandupConfig(channelID)
		if err != nil {
			logger.Error("An error occurred while fetching standup config for channel", err, map[string]interface{}{"channelID": channelID})
			http.Error(w, "An error occurred while fetching standup config", http.StatusInternalServerError)
			return err
		}

		if standupConfig.Enabled {
			activeStandupChannels = append(activeStandupChannels, channelID)
		}
	}

	data, err := json.Marshal(activeStandupChannels)
	if err != nil {
		logger.Error("An error occurred serializing active standup channel list", err, nil)
		http.Error(w, "An error occurred serializing active standup channel list", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}

	return nil
}
