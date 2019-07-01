package command

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
	"strings"
)

func commandRemoveMembers() *Config {
	return &Config{
		Command: &model.Command{
			Trigger:          "removemembers",
			AutoComplete:     true,
			AutoCompleteDesc: "Removes specified members from this channel's standup.",
			AutoCompleteHint: "usernames...",
		},
		HelpText: "* doesn't remove the users from the channel\n" +
			"	* usernames can be specified as @ mentions",
		Validate: validateRemoveMembers,
		Execute:  executeRemoveMembers,
	}
}

func validateRemoveMembers(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	// we need at least one member
	if len(args) < 1 {
		return util.SendEphemeralText("Please specify at least one user to remove")
	}

	// removing @ from usernames if they were specified using mentions.
	// We do allow it. Pretty cool!
	usernamesNotFound := []string{}
	userIDs := []string{}

	for _, arg := range args {
		argToUse := arg
		if arg[0] == '@' {
			argToUse = arg[1:]
		}

		user, err := config.Mattermost.GetUserByUsername(argToUse)
		if err != nil {
			usernamesNotFound = append(usernamesNotFound, argToUse)
		} else {
			userIDs = append(userIDs, user.Id)
		}
	}

	// saving formatted usernames to context for later use
	context.Props["usernamesNotFound"] = usernamesNotFound
	context.Props["userIDs"] = userIDs
	return nil, nil
}

func executeRemoveMembers(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	userIDs := context.Props["userIDs"].([]string)
	usersNotInStandup, usersRemoved, err := removeMembersFromStandup(userIDs, context.CommandArgs.ChannelId)
	if err != nil {
		return util.SendEphemeralText("An error occurred while removing members from standup")
	}

	for i := range usersRemoved {
		user, err := config.Mattermost.GetUser(usersRemoved[i])
		if err != nil {
			logger.Error("User not found", err, nil)
		} else {
			usersRemoved[i] = user.Username
		}
	}

	for i := range usersNotInStandup {
		user, err := config.Mattermost.GetUser(usersNotInStandup[i])
		if err != nil {
			logger.Error("User not found", err, nil)
		} else {
			usersNotInStandup[i] = user.Username
		}
	}

	text := ""

	if len(usersRemoved) > 0 {
		text += "Removed users from standup: " + strings.Join(usersRemoved, ", ")
	}

	if len(usersNotInStandup) > 0 {
		text += "\nUsers not in standup: " + strings.Join(usersNotInStandup, ", ")
	}

	if len(context.Props["usernamesNotFound"].([]string)) > 0 {
		text += "\nUsers not found: " + strings.Join(context.Props["usernamesNotFound"].([]string), ", ")
	}

	return &model.CommandResponse{
		Type: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: text,
	}, nil
}

func removeMembersFromStandup(userIDs []string, channelID string) ([]string, []string, error) {
	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		return nil, nil, err
	}

	if standupConfig == nil {
		return nil, nil, errors.New("Unable to find standup config for channel: " + channelID)
	}

	membersNotInStandup := util.Difference(userIDs, standupConfig.Members)
	membersRemovedFromStandup := util.Difference(standupConfig.Members, util.Difference(standupConfig.Members, userIDs))

	standupConfig.Members = util.Difference(standupConfig.Members, userIDs)

	_, err = standup.SaveStandupConfig(standupConfig)
	if err != nil {
		return nil, nil, err
	}

	return membersNotInStandup, membersRemovedFromStandup, nil
}
