package command

import (
	"fmt"
	"github.com/mattermost/mattermost-server/model"
	"github.com/harshilsharma/standup-raven/server/config"
	"github.com/harshilsharma/standup-raven/server/otime"
	"github.com/harshilsharma/standup-raven/server/standup/notification"
	"github.com/harshilsharma/standup-raven/server/util"
	"strings"
	"time"
)

const (
	dateLayout = "02-01-2006"
)

func commandStandup() *Config {
	return &Config{
		Command: &model.Command{
			Trigger:          config.CommandPrefix + "report",
			AutoComplete:     true,
			AutoCompleteDesc: "Generates standup reports for provided dates",
			AutoCompleteHint: "<dates...> <visibility>",
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
	if len(args) == 0 {
		args = []string{time.Now().Format(dateLayout), notification.ReportVisibilityPrivate}
	} else if len(args) == 1 {
		lastArg := strings.ToLower(args[len(args)-1])

		if lastArg != notification.ReportVisibilityPublic && lastArg != notification.ReportVisibilityPrivate {
			args = append(args, notification.ReportVisibilityPrivate)
		} else {
			args = []string{time.Now().Format(dateLayout), lastArg}
		}
	}

	context.Props["visibility"] = strings.ToLower(args[len(args)-1])

	// processing dates to generate report for
	context.Props["dates"] = make([]otime.OTime, len(args))
	count := 0

	for _, arg := range args[0 : len(args)-1] {
		t, err := time.Parse(dateLayout, arg)
		if err != nil {
			return util.SendEphemeralText(fmt.Sprintf("Error parsing this date: %s. Please specify date in format: DD-MM-YYYY", arg))
		}

		context.Props["dates"].([]otime.OTime)[count] = otime.OTime{t}
		count++
	}

	context.Props["dates"] = context.Props["dates"].([]otime.OTime)[0:count]

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
