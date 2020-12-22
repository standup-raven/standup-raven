package middleware

import (
	"context"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/standup-raven/standup-raven/server/config"
)

type ContextKey string

const (
	CtxKeyUserID = ContextKey("user_id")
)

// Authenticated middleware verifies the request was made by a logged in Mattermost user.
// this is checked by the presence of Mattermost-User-Id HTTP header.
func Authenticated(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError) {
	userID := r.Header.Get(config.HeaderMattermostUserID)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, model.NewAppError("MiddlewareAuthenticate", "", nil, "Unauthorized", http.StatusUnauthorized)
	}

	ctxWithUser := context.WithValue(r.Context(), CtxKeyUserID, userID)
	rWithUser := r.WithContext(ctxWithUser)
	return rWithUser, nil
}
