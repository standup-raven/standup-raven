package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateCommandHelp(t *testing.T) {
	response, appErr := validateCommandHelp([]string{}, Context{})
	assert.Nil(t, response)
	assert.Nil(t, appErr)
}

func Test_executeCommandHelp(t *testing.T) {
	response, appErr := executeCommandHelp([]string{}, Context{})
	assert.NotNil(t, response)
	assert.Nil(t, appErr)
}
