package utils

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thoas/go-funk"

	"github.com/standup-raven/standup-raven/server/config"
)

func TestUserIcon(t *testing.T) {
	actual := UserIcon("user_id_1")
	expected := "![User Avatar](/api/v4/users/user_id_1/image =20x20)"

	assert.Equal(t, expected, actual)
}

func TestGetUserRoles(t *testing.T) {
	mockAPI := &plugintest.API{}

	mockAPI.On("GetChannelMember", "channel_id_1", "user_id_1").Return(
		&model.ChannelMember{
			Roles: model.CHANNEL_ADMIN_ROLE_ID,
		},
		nil,
	)

	mockAPI.On("GetTeamMember", "team_id_1", "user_id_1").Return(
		&model.TeamMember{Roles: model.TEAM_ADMIN_ROLE_ID},
		nil,
	)

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{Roles: model.SYSTEM_ADMIN_ROLE_ID},
		nil,
	)

	mockAPI.On("GetChannel", "channel_id_1").Return(
		&model.Channel{TeamId: "team_id_1"},
		nil,
	)

	config.Mattermost = mockAPI

	roles, appErr := GetUserRoles("user_id_1", "channel_id_1")
	assert.Nil(t, appErr)
	assert.True(t, funk.Contains(roles, model.TEAM_ADMIN_ROLE_ID))
	assert.True(t, funk.Contains(roles, model.SYSTEM_ADMIN_ROLE_ID))
	assert.True(t, funk.Contains(roles, model.CHANNEL_ADMIN_ROLE_ID))
}

func TestGetUserRoles_GetChannelMembers_Error(t *testing.T) {
	mockAPI := &plugintest.API{}

	mockAPI.On("GetChannelMember", "channel_id_1", "user_id_1").Return(
		nil, model.NewAppError("", "", nil, "", 0),
	)

	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))

	config.Mattermost = mockAPI

	_, appErr := GetUserRoles("user_id_1", "channel_id_1")
	assert.NotNil(t, appErr)
}

func TestGetUserRoles_GetTeamMembers_Error(t *testing.T) {
	mockAPI := &plugintest.API{}

	mockAPI.On("GetChannelMember", "channel_id_1", "user_id_1").Return(
		&model.ChannelMember{
			Roles: model.CHANNEL_ADMIN_ROLE_ID,
		},
		nil,
	)

	mockAPI.On("GetChannel", "channel_id_1").Return(
		&model.Channel{TeamId: "team_id_1"},
		nil,
	)

	mockAPI.On("GetTeamMember", "team_id_1", "user_id_1").Return(
		nil, model.NewAppError("", "", nil, "", 0),
	)

	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))

	config.Mattermost = mockAPI

	_, appErr := GetUserRoles("user_id_1", "channel_id_1")
	assert.NotNil(t, appErr)
}

func TestGetUserRoles_GetUser_Error(t *testing.T) {
	mockAPI := &plugintest.API{}

	mockAPI.On("GetChannelMember", "channel_id_1", "user_id_1").Return(
		&model.ChannelMember{
			Roles: model.CHANNEL_ADMIN_ROLE_ID,
		},
		nil,
	)

	mockAPI.On("GetTeamMember", "team_id_1", "user_id_1").Return(
		&model.TeamMember{Roles: model.TEAM_ADMIN_ROLE_ID},
		nil,
	)

	mockAPI.On("GetChannel", "channel_id_1").Return(
		&model.Channel{TeamId: "team_id_1"},
		nil,
	)

	mockAPI.On("GetUser", "user_id_1").Return(
		nil, model.NewAppError("", "", nil, "", 0),
	)

	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))

	config.Mattermost = mockAPI

	_, appErr := GetUserRoles("user_id_1", "channel_id_1")
	assert.NotNil(t, appErr)
}

func TestGetUserRoles_GetChannel_Error(t *testing.T) {
	mockAPI := &plugintest.API{}

	mockAPI.On("GetChannelMember", "channel_id_1", "user_id_1").Return(
		&model.ChannelMember{
			Roles: model.CHANNEL_ADMIN_ROLE_ID,
		},
		nil,
	)

	mockAPI.On("GetTeamMember", "team_id_1", "user_id_1").Return(
		&model.TeamMember{Roles: model.TEAM_ADMIN_ROLE_ID},
		nil,
	)

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{Roles: model.SYSTEM_ADMIN_ROLE_ID},
		nil,
	)

	mockAPI.On("GetChannel", "channel_id_1").Return(
		nil, model.NewAppError("", "", nil, "", 0),
	)

	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))

	config.Mattermost = mockAPI

	_, appErr := GetUserRoles("user_id_1", "channel_id_1")
	assert.NotNil(t, appErr)
}
