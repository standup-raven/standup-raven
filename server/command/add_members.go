package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/thoas/go-funk"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

func commandAddMembers() *Config {
	return &Config{
		AutocompleteData: &model.AutocompleteData{
			Trigger: "addmembers",
			Hint:    "[username 1] [username 2] [username 3]...",
			HelpText: "Adds specified members to the the current channel's standup. " +
				"Members are also automatically added to the current channel if not already part of it.",
			Arguments: []*model.AutocompleteArg{
				{
					Type:     model.AutocompleteArgTypeText,
					Required: true,
					HelpText: "Use @ mentions to quickly refer to a user. For example `@johndoe`",
					Data: &model.AutocompleteTextArg{
						Hint:    "Usernames",
						Pattern: ".+",
					},
				},
			},
		},
		ExtraHelpText: "* usernames can be specified as @ mentions", // TODO what is this helptext needed for?
		Validate:      validateAddMembers,
		Execute:       executeAddMembers,
	}
}

func validateAddMembers(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	// we need at least one member
	if len(args) < 1 {
		return util.SendEphemeralText("Please specify at least one user to add")
	}

	// removing @ from usernames if they were specified using mentions.
	userIds := make(map[string]string, len(args))
	for _, username := range args {
		username = strings.TrimLeft(username, "@")

		// preventing duplicates
		if _, ok := userIds[username]; ok {
			continue
		}

		user, err := config.Mattermost.GetUserByUsername(username)
		if err != nil {
			return util.SendEphemeralText("Couldn't find user with username: " + username)
		}
		userIds[username] = user.Id
	}

	// saving formatted usernames to context for later use
	context.Props["userIds"] = funk.Values(userIds).([]string)
	return nil, nil
}

func executeAddMembers(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	userIds := context.Props["userIds"].([]string)

	// inviting members to standup channel
	addedUsers, notAddedUsers := addChannelMembers(userIds, context.CommandArgs.ChannelId)

	// adding successfully invited members to channel's standup config
	if err := addStandupMembers(addedUsers, context.CommandArgs.ChannelId); err != nil {
		return util.SendEphemeralText("Error occurred while adding standup members: " + err.Error())
	}

	return &model.CommandResponse{
		Type: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: buildSuccessMessage(addedUsers, notAddedUsers),
	}, nil
}

func addChannelMembers(userIds []string, channelID string) ([]string, []string) {
	var addedUsers, notAddedUsers []string

	for _, userID := range userIds {
		if _, appErr := config.Mattermost.AddChannelMember(channelID, userID); appErr != nil {
			logger.Error(fmt.Sprintf("Error adding user [%s] to channel [%s]", userID, channelID), appErr, nil)
			notAddedUsers = append(notAddedUsers, userID)
			continue
		}

		addedUsers = append(addedUsers, userID)
	}

	return addedUsers, notAddedUsers
}

func addStandupMembers(usernames []string, channelID string) error {
	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		return err
	}

	if standupConfig == nil {
		return errors.New("standup is not configured for this channel")
	}

	standupConfig.Members = append(standupConfig.Members, usernames...)
	_, err = standup.SaveStandupConfig(standupConfig)
	if err != nil {
		return err
	}

	return nil
}

func buildSuccessMessage(addedUsers, notAddedUsers []string) string {
	text := fmt.Sprintf("%d users added successfully.", len(addedUsers))
	if len(notAddedUsers) > 0 {
		text += fmt.Sprintf(" Following users couldn't be added: %s", strings.Join(notAddedUsers, ", "))
		text += "\nMake sure these users users exist on the system."
	}

	return text
}
