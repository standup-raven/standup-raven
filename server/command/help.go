package command

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

func commandHelp() *Config {
	return &Config{
		AutocompleteData: &model.AutocompleteData{
			Trigger:  "help",
			HelpText: "Shows help on various standup commands",
		},
		Validate: validateCommandHelp,
		Execute:  executeCommandHelp,
	}
}

func validateCommandHelp(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	return nil, nil
}

func executeCommandHelp(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	helpText := generateHelpText([]*Config{
		commandConfig(),
		commandAddMembers(),
		commandViewConfig(),
		commandRemoveMembers(),
		commandStandup(),
		commandHelp(),
	})

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         helpText,
	}, nil
}

func generateHelpText(commands []*Config) string {
	text := "TODO"

	// TODO
	//for _, command := range commands {
	//	text += fmt.Sprintf("* `%s %s` - %s \n\t%s\n", command.Command.Trigger, command.Command.AutoCompleteHint, command.Command.AutoCompleteDesc, strings.Replace(command.HelpText, "\n", "\n\t", -1))
	//}

	return text
}
