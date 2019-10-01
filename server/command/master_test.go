package command

import (
	"github.com/bouk/monkey"
	"github.com/mattermost/mattermost-server/model"
	"github.com/standup-raven/standup-raven/server/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TearDown() {
	monkey.UnpatchAll()
}

func TestCommandMaster_Validation(t *testing.T) {
	defer TearDown()
	
	command := Master()
	context := Context{
		Props: map[string]interface{}{},
	}
	
	response, err := command.Validate([]string{commandViewConfig().Command.Trigger}, context)
	assert.Nil(t, err)
	assert.Nil(t, response)

	response, err = command.Validate([]string{"non-existing-command"}, context)
	assert.Nil(t, err)
	assert.Equal(t, "Invalid command: non-existing-command", response.Text)

	response, err = command.Validate([]string{}, context)
	assert.Nil(t, err)
	assert.Nil(t, response)
	
	monkey.Patch(commandViewConfig().Validate, func([]string, Context) (*model.CommandResponse, *model.AppError) {
		return util.SendEphemeralText("error")
	})

	response, err = command.Validate([]string{commandViewConfig().Command.Trigger}, context)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "error", response.Text)
	
}

func TestCommandMaster_Execution(t *testing.T) {
	defer TearDown()
	
	command := Master()
	dummyCommand := &Config{
		Execute: func([]string, Context) (*model.CommandResponse, *model.AppError) {
			return nil, nil
		},
	}
	context := Context{
		Props: map[string]interface{}{
			"subCommand": dummyCommand,
			"subCommandArgs": []string{"some-command"},
		},
	}
	
	response, err := command.Execute([]string{}, context)
	assert.Nil(t, err)
	assert.Nil(t, response)

	//dummyCommand = &Config{
	//	Validate: func([]string, Context) (*model.CommandResponse, *model.AppError) {
	//		return util.SendEphemeralText("error")	
	//	},
	//}
	//context = Context{
	//	Props: map[string]interface{}{
	//		"subCommand":     dummyCommand,
	//		"subCommandArgs": []string{"some-command"},
	//	},
	//}
	//
	//response, err = command.Execute([]string{}, context)
	//assert.Nil(t, err)
	//assert.NotNil(t, response)
	//assert.Equal(t, "error", response.Text)
}
