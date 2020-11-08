package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Context struct {
	CommandArgs *model.CommandArgs
	Props       map[string]interface{}
}

type Config struct {
	AutocompleteData *model.AutocompleteData
	ExtraHelpText    string
	Execute          func([]string, Context) (*model.CommandResponse, *model.AppError)
	Validate         func([]string, Context) (*model.CommandResponse, *model.AppError)
}

func (c *Config) Syntax() string {
	return fmt.Sprintf("/%s %s", c.AutocompleteData.Trigger, c.AutocompleteData.HelpText)
}

func (c *Config) GetHelpText() string {
	helpText := fmt.Sprintf(
		"* `%s %s` - %s",
		c.AutocompleteData.Trigger,
		c.AutocompleteData.Hint,
		c.AutocompleteData.HelpText,
	)

	if c.ExtraHelpText != "" {
		helpText += "\n\t" + strings.ReplaceAll(c.ExtraHelpText, "\n", "\n\t")
	}

	helpText += "\n\n"
	return helpText
}

// Remember to add any new command to `executeCommandHelp` as well for
// generating help text.
// executeCommandHelp doesn't use this map to prevent circular imports.
var commands = map[string]*Config{
	commandConfig().AutocompleteData.Trigger:        commandConfig(),
	commandAddMembers().AutocompleteData.Trigger:    commandAddMembers(),
	commandRemoveMembers().AutocompleteData.Trigger: commandRemoveMembers(),
	commandStandup().AutocompleteData.Trigger:       commandStandup(),
	commandHelp().AutocompleteData.Trigger:          commandHelp(),
}
