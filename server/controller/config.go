package controller

import (
	"encoding/json"
	"fmt"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
	"net/http"
)

var getConfig = &Endpoint{
	Path:         "/config",
	Method:       http.MethodGet,
	Execute:      executeGetConfig,
	RequiresAuth: true,
}

var setConfig = &Endpoint{
	Path:         "/config",
	Method:       http.MethodPost,
	Execute:      executeSetConfig,
	RequiresAuth: true,
}

var getTimezone = &Endpoint{
	Path:         "/timezone",
	Method:       http.MethodGet,
	Execute:      executeGetLocation,
	RequiresAuth: true,
}


func executeGetConfig(w http.ResponseWriter, r *http.Request) error {
	channelId := r.URL.Query().Get("channel_id")
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
	fmt.Println("decoder=",decoder)
	conf := &standup.StandupConfig{}
	if err := decoder.Decode(&conf); err != nil {
		logger.Error("Could not decode request body", err, map[string]interface{}{"request": util.DumpRequest(r)})
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return err
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

	return nil
}

func executeGetLocation(w http.ResponseWriter, r *http.Request) error {
	location := config.GetConfig().TimeZone
	
	data, err := json.Marshal(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Couldn't serialize config data", err, map[string]interface{}{"location": location})
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		logger.Error("Error occurred in writing data to HTTP response", err, map[string]interface{}{"data": string(data)})
		return err
	}

	return nil
}
