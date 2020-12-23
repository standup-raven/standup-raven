package middleware

import (
	"context"
	"github.com/standup-raven/standup-raven/server/util"
	"github.com/thoas/go-funk"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/standup-raven/standup-raven/server/config"
)

type ContextKey string

const (
	CtxKeyUserID = ContextKey("user_id")
	
	RoleTypeEffectiveChannelAdmin = "isEffectiveChannelAdmin"
	RoleTypeGuest = "isGuest"
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

func SetUserRoles(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError) {
	userID := r.Context().Value(CtxKeyUserID).(string)
	
	// LOL
	// TODO pass channel_id query param in set config API
	channelID := r.URL.Query().Get("channel_id")
	
	userRoles, err := util.GetUserRoles(userID, channelID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return nil, model.NewAppError("MiddlewareSetUserRoles", "", map[string]interface{}{"userID": userID, "channelID": channelID}, "", http.StatusInternalServerError)
	}

	userRoleTypes := map[string]bool{
		"isEffectiveChannelAdmin": false,
		"isGuest": false,
	}
	
	userRoleTypes[RoleTypeEffectiveChannelAdmin] = funk.Contains(userRoles, )
}
