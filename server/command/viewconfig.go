package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

func commandViewConfig() *Config {
	return &Config{
		Command: &model.Command{
			Trigger:          "viewconfig",
			AutoCompleteDesc: "View standup settings for this channel.",
			AutoComplete:     true,
		},
		HelpText: "",
		Validate: validateViewConfig,
		Execute:  executeViewConfig,
	}
}

func validateViewConfig(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	return nil, nil
}

func executeViewConfig(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	message := ""

	standupConfig, err := standup.GetStandupConfig(context.CommandArgs.ChannelId)
	if err != nil {
		return util.SendEphemeralText("Error occurred while extracting standup config")
	} else if standupConfig == nil {
		message = "Standup not configured for this channel"
	} else {
		membersString := "no members present"
		if len(standupConfig.Members) > 0 {
			members := make([]string, len(standupConfig.Members))

			for i, userId := range standupConfig.Members {
				user, err := config.Mattermost.GetUser(userId)
				if err != nil {
					logger.Error("Couldn't fetch details for user with ID: "+userId, err, nil)
					return util.SendEphemeralText("Couldn't fetch details for user with ID: " + userId)
				}

				members[i] = user.Username
			}

			membersString = strings.Join(members, ", ")
		}
		windowOpenReminder := "disabled"
		windowCloseReminder := "disabled"
		if standupConfig.WindowOpenReminderEnabled {
			windowOpenReminder = "enabled"
		}
		if standupConfig.WindowCloseReminderEnabled {
			windowCloseReminder = "enabled"
		}
		message = fmt.Sprintf(
			"Window open time: %s \nWindow close time: %s \nTimezone: %s \nWindow open reminder: %s \nWindow close reminder: %s \nReport format: %s \nMembers: %s",
			standupConfig.WindowOpenTime.GetTimeString(),
			standupConfig.WindowCloseTime.GetTimeString(),
			standupConfig.Timezone,
			windowOpenReminder,
			windowCloseReminder,
			standupConfig.ReportFormat,
			membersString,
		)
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         message,
	}, nil
}
