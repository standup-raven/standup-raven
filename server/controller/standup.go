package controller

import (
	"encoding/json"
	"errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"net/http"
)

var getStandup = &Endpoint{
	Path:    "/standup",
	Method:  http.MethodGet,
	Execute: executeGetStandup,
	Middlewares: []middleware.Middleware{
		middleware.Authenticate,
	},
}

var saveStandup = &Endpoint{
	Path:    "/standup",
	Method:  http.MethodPost,
	Execute: executeSaveStandup,
	Middlewares: []middleware.Middleware{
		middleware.Authenticate,
	},
}

func executeSaveStandup(w http.ResponseWriter, r *http.Request) error {
	userStandup := &standup.UserStandup{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(userStandup); err != nil {
		logger.Error("Couldn't decode request body into user standup object", err, nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	userStandup.UserID = r.Header.Get(config.HeaderMattermostUserId)

	if err := userStandup.IsValid(); err != nil {
		logger.Info("user standup validation failed", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := standup.SaveUserStandup(userStandup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if _, err := w.Write([]byte("ok")); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, nil)
		return err
	}

	return nil
}

func executeGetStandup(w http.ResponseWriter, r *http.Request) error {
	userID := r.Header.Get(config.HeaderMattermostUserId)
	channelID := r.URL.Query().Get("channel_id")
	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		http.Error(w, "Error occurred while fetching standup config", http.StatusInternalServerError)
		return err
	}
	if standupConfig == nil {
		http.Error(w, "Standup not configured for channel", http.StatusNotFound)
		return errors.New("standup not configured for channel: " + channelID)
	}

	userStandup, err := standup.GetUserStandup(userID, channelID, otime.Now(standupConfig.Timezone))
	if err != nil {
		http.Error(w, "Error occurred while fetching user standup", http.StatusInternalServerError)
		return err
	} else if userStandup == nil {
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	data, err := json.Marshal(userStandup)
	if err != nil {
		logger.Error("Error occurred while marshaling user standup", err, nil)
		http.Error(w, "Error occurred while marshaling user standup", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, nil)
		return err
	}

	return nil
}
