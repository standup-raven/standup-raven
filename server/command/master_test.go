package command

import (
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/assert"
)

func TearDown() {
	monkey.UnpatchAll()
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
			"subCommand":     dummyCommand,
			"subCommandArgs": []string{"some-command"},
		},
	}

	response, err := command.Execute([]string{}, context)
	assert.Nil(t, err)
	assert.Nil(t, response)
}
