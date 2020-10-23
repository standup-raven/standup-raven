package command

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/standup/notification"
	"github.com/standup-raven/standup-raven/server/util"
	"strings"
	"time"
)

const (
	dateLayout = "02-01-2006"
)

func commandStandup() *Config {
	return &Config{
		AutocompleteData: &model.AutocompleteData{
			Trigger:          "report",
			HelpText: "Generates standup reports for provided dates",
			//AutoCompleteHint: "<dates...> <visibility>",
			Arguments: []*model.AutocompleteArg{
				{
					Name: "Date",
					HelpText: "Date to generate standup report for",
					Type: model.AutocompleteArgTypeText,
					Required: true,
					Data: &model.AutocompleteTextArg{
						Hint:    "Date",
						Pattern: "\\d\\d-\\d\\d-\\d\\d\\d\\d",
					},
				},
			},
		},
		HelpText: "* dates must be in `DD-MM-YYYY` format\n" +
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
	if len(args) == 0 {
		args = []string{otime.Now(standupConfig.Timezone).Format(dateLayout), notification.ReportVisibilityPrivate}
	} else if len(args) == 1 {
		lastArg := strings.ToLower(args[0])

		if lastArg != notification.ReportVisibilityPublic && lastArg != notification.ReportVisibilityPrivate {
			args = append(args, notification.ReportVisibilityPrivate)
		} else {
			args = []string{otime.Now(standupConfig.Timezone).Format(dateLayout), lastArg}
		}
	}

	context.Props["visibility"] = strings.ToLower(args[len(args)-1])

	// processing dates to generate report for
	dates := make([]otime.OTime, len(args)-1)
	count := 0

	for _, arg := range args[0 : len(args)-1] {
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
	channelId := context.CommandArgs.ChannelId
	visibility := context.Props["visibility"].(string)
	userId := context.CommandArgs.UserId

	for _, date := range context.Props["dates"].([]otime.OTime) {
		if err := notification.SendStandupReport([]string{channelId}, date, visibility, userId, false); err != nil {
			// continue
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Standup reports generated",
	}, nil
}
