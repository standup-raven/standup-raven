package middleware

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/standup-raven/standup-raven/server/config"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request) *model.AppError {
	userId := r.Header.Get(config.HeaderMattermostUserId)

	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return model.NewAppError("MiddlewareAuthenticate", "", nil, "Unauthorized", http.StatusUnauthorized)
	}
	return nil
}
