package controller

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
	"net/http"
	"strings"
)

var getConfig = &Endpoint{
	Path:    "/config",
	Method:  http.MethodGet,
	Execute: executeGetConfig,
	Middlewares: []middleware.Middleware{
		middleware.Authenticate,
	},
}

var setConfig = &Endpoint{
	Path:    "/config",
	Method:  http.MethodPost,
	Execute: executeSetConfig,
	Middlewares: []middleware.Middleware{
		middleware.Authenticate,
	},
}

var getDefaultTimezone = &Endpoint{
	Path:    "/timezone",
	Method:  http.MethodGet,
	Execute: executeGetDefaultTimezone,
}

var getActiveStandupChannels = &Endpoint{
	Path: "/active-channels",
	Method: http.MethodGet,
	Execute: executeGetActiveStandupChannels,
	Middlewares: []middleware.Middleware{
		middleware.Authenticate,
	}, 
}

func executeGetConfig(w http.ResponseWriter, r *http.Request) error {
	channelId := r.URL.Query().Get("channel_id")
	userID := r.Header.Get(config.HeaderMattermostUserId)

	// verifying if user is an effective channel admin
	source := r.URL.Query().Get("source")
	if config.GetConfig().PermissionSchemaEnabled && source != "standup-modal" {
		isAdmin, appErr := isEffectiveAdmin(userID, channelId)

		if appErr != nil {
			http.Error(w, "An error occurred while verifying user permissions", appErr.StatusCode)
			logger.Error("An error occurred while verifying user permissions", errors.New(appErr.Error()), nil)
			return appErr
		}

		if !isAdmin {
			http.Error(w, "You do not have permission to perform this operation", http.StatusUnauthorized)
			return nil
		}
	}

	c, err := standup.GetStandupConfig(channelId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if c == nil {
		http.Error(w, "Standup not configured for this channel", http.StatusNotFound)
		return nil
	}

	// TODO: make use of ToJson function for sending conf in response
	data, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Couldn't serialize config data", err, map[string]interface{}{"config": c.ToJson()})
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}

	return nil
}

func executeSetConfig(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	conf := &standup.StandupConfig{}
	if err := decoder.Decode(&conf); err != nil {
		logger.Error("Could not decode request body", err, map[string]interface{}{"request": util.DumpRequest(r)})
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return err
	}

	userID := r.Header.Get(config.HeaderMattermostUserId)
	channelID := conf.ChannelId

	// verifying if user is an effective channel admin
	if config.GetConfig().PermissionSchemaEnabled {
		isAdmin, appErr := isEffectiveAdmin(userID, channelID)

		if appErr != nil {
			http.Error(w, "An error occurred while verifying user permissions", appErr.StatusCode)
			logger.Error("An error occurred while verifying user permissions", errors.New(appErr.Error()), nil)
			return appErr
		}

		if !isAdmin {
			http.Error(w, "You do not have permission to perform this operation", http.StatusUnauthorized)
			return nil
		}
	}

	if err := conf.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	conf, err := standup.SaveStandupConfig(conf)
	if err != nil {
		http.Error(w, "Error occurred while saving standup conf", http.StatusInternalServerError)
		return err
	}

	if err := standup.AddStandupChannel(conf.ChannelId); err != nil {
		http.Error(w, "Error occurred while saving standup conf", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(conf.ToJson())); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"config": conf.ToJson()})
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
			"channel_id": conf.ChannelId,
		},
		&model.WebsocketBroadcast{
			UserId: r.Header.Get(config.HeaderMattermostUserId),
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

func isEffectiveAdmin(userID string, channelID string) (bool, *model.AppError) {
	if isChannelAdmin, appErr := isChannelAdmin(userID, channelID); appErr != nil {
		return false, appErr
	} else if isChannelAdmin {
		config.Mattermost.LogDebug("User is channel admin", "userID", userID)
		return true, nil
	}

	channel, appErr := config.Mattermost.GetChannel(channelID)
	if appErr != nil {
		return false, appErr
	}

	if isTeamAdmin, appErr := isTeamAdmin(userID, channel.TeamId); appErr != nil {
		return false, appErr
	} else if isTeamAdmin {
		config.Mattermost.LogDebug("User is team admin", "userID", userID)
		return true, nil
	}

	if isSystemAdmin, appErr := isSystemAdmin(userID); appErr != nil {
		return false, appErr
	} else if isSystemAdmin {
		config.Mattermost.LogDebug("User is system admin", "userID", userID)
		return true, nil
	}

	return false, nil
}

func isChannelAdmin(userID string, channelID string) (bool, *model.AppError) {
	channelMember, appErr := config.Mattermost.GetChannelMember(channelID, userID)
	if appErr != nil {
		return false, appErr
	}

	return strings.Contains(channelMember.Roles, model.CHANNEL_ADMIN_ROLE_ID), nil
}

func isTeamAdmin(userID string, teamID string) (bool, *model.AppError) {
	teamMember, appErr := config.Mattermost.GetTeamMember(teamID, userID)
	if appErr != nil {
		return false, appErr
	}

	return strings.Contains(teamMember.Roles, model.TEAM_ADMIN_ROLE_ID), nil
}

func isSystemAdmin(userID string) (bool, *model.AppError) {
	user, appErr := config.Mattermost.GetUser(userID)
	if appErr != nil {
		return false, appErr
	}

	return strings.Contains(user.Roles, model.SYSTEM_ADMIN_ROLE_ID), nil
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
		http.Error(w,"An error occurred serializing active standup channel list", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}
	
	return nil
}
