package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/standup/notification"
	util "github.com/standup-raven/standup-raven/server/utils"

)

const (
	dateLayout = "02-01-2006"
)

func commandStandup() *Config {
	return &Config{
		AutocompleteData: &model.AutocompleteData{
			Trigger:  "report",
			HelpText: "Generate standup reports for provided dates",
			RoleID:   model.SYSTEM_USER_ROLE_ID,
			Arguments: []*model.AutocompleteArg{
				{
					HelpText: "Report visibility",
					Type:     model.AutocompleteArgTypeStaticList,
					Required: true,
					Data: model.AutocompleteStaticListArg{
						PossibleArguments: []model.AutocompleteListItem{
							{
								Item:     "Public",
								HelpText: "Generated report(s) will be visible to everyone in this channel.",
							},
							{
								Item:     "Private",
								HelpText: "Generated report(s) will only be visible to you.",
							},
						},
					},
				},
				{
					HelpText: "Date to generate standup report for. Dates must be in `DD-MM-YYYY` format.",
					Type:     model.AutocompleteArgTypeText,
					Required: true,
					Data: &model.AutocompleteTextArg{
						Hint:    "[date 1] [date 2] [date 3]...",
						Pattern: "\\d\\d-\\d\\d-\\d\\d\\d\\d",
					},
				},
			},
		},
		ExtraHelpText: "* dates must be in `DD-MM-YYYY` format\n" +
			"* visibility can be one of the following -\n" +
			"	* `public` - generated report is visible to everyone in the channel\n" +
			"	* `private` - generated report is visible only to you",
		Validate: validateCommandStandup,
		Execute:  executeCommandStandup,
	}
}

func validateCommandStandup(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	standupConfig, err := standup.GetStandupConfig(context.CommandArgs.ChannelId)
	if err != nil {
		return util.SendEphemeralText("Error getting standup config of the channel")
	}

	if standupConfig == nil {
		return util.SendEphemeralText("Standup not configured for the channel")
	}
	if len(args) < 2 {
		return util.SendEphemeralText("Please specify report format and dates to generate report for.")
	}

	context.Props["visibility"] = strings.ToLower(args[0])

	// processing dates to generate report for
	dates := make([]otime.OTime, len(args)-1)
	count := 0

	for _, arg := range args[1:] {
		t, err := time.Parse(dateLayout, arg)
		if err != nil {
			return util.SendEphemeralText(fmt.Sprintf("Error parsing this date: %s. Please specify date in format: DD-MM-YYYY", arg))
		}

		dates[count] = otime.OTime{Time: t}
		count++
	}

	context.Props["dates"] = dates[0:count]
	return nil, nil
}

func executeCommandStandup(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	channelID := context.CommandArgs.ChannelId
	visibility := context.Props["visibility"].(string)
	userID := context.CommandArgs.UserId

	for _, date := range context.Props["dates"].([]otime.OTime) {
		_ = notification.SendStandupReport([]string{channelID}, date, visibility, userID, false)
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Standup reports generated",
	}, nil
}
