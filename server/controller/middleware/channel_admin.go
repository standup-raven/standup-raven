package middleware

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/standup-raven/standup-raven/server/config"
	"net/http"
	"strings"
)

func ChannelAdmin(w http.ResponseWriter, r *http.Request) *model.AppError {
	userID := r.Header.Get(config.HeaderMattermostUserId)
	channelID := r.Header.Get(config.HeaderMattermostChannelId)

	if isChannelAdmin, appErr := isChannelAdmin(userID, channelID); appErr != nil {
		return appErr
	} else if isChannelAdmin {
		config.Mattermost.LogInfo("Is channel admin")
		return nil
	}
	
	channel, appErr := config.Mattermost.GetChannel(channelID)
	if appErr != nil {
		return appErr
	}

	if isTeamAdmin, appErr := isTeamAdmin(userID, channel.TeamId); appErr != nil {
		return appErr
	} else if isTeamAdmin {
		config.Mattermost.LogInfo("Is team admin")
		return nil
	}
	
	if isSystemAdmin, appErr := isSystemAdmin(userID); appErr != nil {
		return appErr
	} else if isSystemAdmin {
		config.Mattermost.LogInfo("Is system admin")
		return nil
	}
	
	return model.NewAppError("ChannelAdmin Middleware", "", nil, "", http.StatusUnauthorized)
}

func isChannelAdmin(userID string, channelID string) (bool, *model.AppError) {
	channelMember, appErr := config.Mattermost.GetChannelMember(channelID, userID)
	if appErr != nil {
		return false, appErr
	}
	
	return strings.Contains(channelMember.Roles, model.CHANNEL_ADMIN_ROLE_ID), nil
}

func isTeamAdmin(userID string, teamID string) (bool, *model.AppError) {
	teamMember, appErr := config.Mattermost.GetTeamMember(teamID, userID)
	if appErr != nil {
		return false, appErr
	}

	return strings.Contains(teamMember.Roles, model.TEAM_ADMIN_ROLE_ID), nil
}

func isSystemAdmin(userID string) (bool, *model.AppError) {
	user, appErr := config.Mattermost.GetUser(userID)
	if appErr != nil {
		return false, appErr
	}

	return strings.Contains(user.Roles, model.SYSTEM_ADMIN_ROLE_ID), nil
}
