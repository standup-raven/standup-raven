package main

import (
	"errors"
	"io/ioutil"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/stretchr/testify/assert"
)

func TearDown() {
	monkey.UnpatchAll()
}

func TestSetUpBot(t *testing.T) {
	defer TearDown()
	bot := &model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: "Bot for Standup Raven.",
	}
	p := &Plugin{}
	helpers := &plugintest.Helpers{}
	helpers.On("EnsureBot", bot).Return("botID", nil)
	api := &plugintest.API{}
	api.On("GetBundlePath").Return("tmp/", nil)
	monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return []byte{}, nil
	})
	api.On("SetProfileImage", "botID", []byte{}).Return(nil)
	p.SetAPI(api)
	p.SetHelpers(helpers)
	_, err := p.setUpBot()
	assert.Nil(t, err, "no error should have been produced")
}

func TestSetUpBot_EnsureBot_Error(t *testing.T) {
	defer TearDown()
	bot := &model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: "Bot for Standup Raven.",
	}
	p := &Plugin{}
	helpers := &plugintest.Helpers{}
	helpers.On("EnsureBot", bot).Return("", errors.New(""))
	p.SetAPI(&plugintest.API{})
	p.SetHelpers(helpers)

	_, err := p.setUpBot()
	assert.NotNil(t, err)
}

func TestSetUpBot_GetBundlePath_Error(t *testing.T) {
	defer TearDown()
	bot := &model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: "Bot for Standup Raven.",
	}
	p := &Plugin{}
	helpers := &plugintest.Helpers{}
	helpers.On("EnsureBot", bot).Return("botID", nil)
	api := &plugintest.API{}
	api.On("GetBundlePath").Return("", errors.New(""))
	p.SetAPI(api)
	p.SetHelpers(helpers)
	_, err := p.setUpBot()
	assert.NotNil(t, err)
}

func TestSetUpBot_Readfile_Error(t *testing.T) {
	defer TearDown()
	bot := &model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: "Bot for Standup Raven.",
	}
	p := &Plugin{}
	helpers := &plugintest.Helpers{}
	helpers.On("EnsureBot", bot).Return("botID", nil)
	api := &plugintest.API{}
	api.On("GetBundlePath").Return("tmp/", nil)
	p.SetAPI(api)
	p.SetHelpers(helpers)
	monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return nil, errors.New("")
	})
	_, err := p.setUpBot()
	assert.NotNil(t, err)
}

func TestSetUpBot_SetProfileImage_Error(t *testing.T) {
	defer TearDown()
	bot := &model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: "Bot for Standup Raven.",
	}
	p := &Plugin{}
	helpers := &plugintest.Helpers{}
	helpers.On("EnsureBot", bot).Return("botID", nil)
	api := &plugintest.API{}
	api.On("GetBundlePath").Return("tmp/", nil)
	monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
		return []byte{}, nil
	})
	api.On("SetProfileImage", "botID", []byte{}).Return(&model.AppError{})
	p.SetAPI(api)
	p.SetHelpers(helpers)
	_, err := p.setUpBot()
	assert.NotNil(t, err)
}
