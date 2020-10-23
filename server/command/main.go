package command

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Context struct {
	CommandArgs *model.CommandArgs
	Props       map[string]interface{}
}

type Config struct {
	//Command  *model.Command
	AutocompleteData *model.AutocompleteData
	HelpText         string
	Execute          func([]string, Context) (*model.CommandResponse, *model.AppError)
	Validate         func([]string, Context) (*model.CommandResponse, *model.AppError)
}

func (c *Config) Syntax() string {
	return fmt.Sprintf("/%s %s", c.AutocompleteData.Trigger, c.AutocompleteData.HelpText)
}

var commands = map[string]*Config{
	commandViewConfig().AutocompleteData.Trigger:    commandViewConfig(),
	commandConfig().AutocompleteData.Trigger:        commandConfig(),
	commandAddMembers().AutocompleteData.Trigger:    commandAddMembers(),
	commandRemoveMembers().AutocompleteData.Trigger: commandRemoveMembers(),
	commandStandup().AutocompleteData.Trigger:       commandStandup(),
	commandHelp().AutocompleteData.Trigger:          commandHelp(),
}
