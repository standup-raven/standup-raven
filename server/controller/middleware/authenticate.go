package middleware

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/standup-raven/standup-raven/server/config"
	"net/http"
)

// Authenticate middleware verifies the request was made by a logged in Mattermost user.
// this is checked by the presence of Mattermost-User-Id HTTP header. 
func Authenticate(w http.ResponseWriter, r *http.Request) *model.AppError {
	userID := r.Header.Get(config.HeaderMattermostUserId)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return model.NewAppError("MiddlewareAuthenticate", "", nil, "Unauthorized", http.StatusUnauthorized)
	}
	return nil
}
