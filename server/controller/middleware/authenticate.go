package middleware

import (
	"context"
	"net/http"

	util "github.com/standup-raven/standup-raven/server/utils"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/standup-raven/standup-raven/server/config"
)

// ContextKey is the custom type for storing
// data in http request context
type ContextKey string

const (
	// CtxKeyUserID is the key corresponding to user ID, stored as a string.
	CtxKeyUserID = ContextKey("user_id")
	// CtxKeyUserRoles stores the high-level user role types - guest and admin.
	CtxKeyUserRoles = ContextKey("user_roles")

	// RoleTypeEffectiveChannelAdmin tells whether the user is effectively a channel admin or not.
	RoleTypeEffectiveChannelAdmin = "isEffectiveChannelAdmin"
	// RoleTypeGuest tells whether the user is a guest user or not
	RoleTypeGuest = "isGuest"
)

// Authenticated middleware verifies the request was made by a logged in Mattermost user.
// this is checked by the presence of Mattermost-User-Id HTTP header.
func Authenticated(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError) {
	userID := r.Header.Get(config.HeaderMattermostUserID)

	if userID == "" {
		return nil, model.NewAppError("MiddlewareAuthenticate", "", nil, "Unauthorized", http.StatusUnauthorized)
	}

	ctxWithUser := context.WithValue(r.Context(), CtxKeyUserID, userID)
	rWithUser := r.WithContext(ctxWithUser)
	return rWithUser, nil
}

// SetUserRoles middleware fetches the user roles and sets them in the request context.
func SetUserRoles(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError) {
	rawUserID := r.Context().Value(CtxKeyUserID)
	if rawUserID == nil {
		return nil, model.NewAppError("SetUserRoles", "couldn't find user ID in context", nil, "Couldn't authenticate user.", http.StatusInternalServerError)
	}

	userID := rawUserID.(string)
	channelID := r.URL.Query().Get("channel_id")
	userRoles, err := util.GetUserRoles(userID, channelID)
	if err != nil {
		return nil, model.NewAppError("MiddlewareSetUserRoles", err.Error(), map[string]interface{}{"userID": userID, "channelID": channelID}, "Couldn't verify user roles.", http.StatusInternalServerError)
	}

	userRoleTypes := map[string]bool{}

	userRolesMap := make(map[string]bool, len(userRoles))
	for _, role := range userRoles {
		userRolesMap[role] = true
	}

	userRoleTypes[RoleTypeEffectiveChannelAdmin] = userRolesMap[model.SYSTEM_ADMIN_ROLE_ID] || userRolesMap[model.TEAM_ADMIN_ROLE_ID] || userRolesMap[model.CHANNEL_ADMIN_ROLE_ID]
	userRoleTypes[RoleTypeGuest] = userRolesMap[model.SYSTEM_GUEST_ROLE_ID]

	ctxWithUserRoles := context.WithValue(r.Context(), CtxKeyUserRoles, userRoleTypes)
	rWithUserRoles := r.WithContext(ctxWithUserRoles)
	return rWithUserRoles, nil
}

// DisallowGuests middleware prevents guest users from accessing the endpoint.
func DisallowGuests(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError) {
	rawUserRoleTypes := r.Context().Value(CtxKeyUserRoles)
	if rawUserRoleTypes == nil {
		return nil, model.NewAppError("DisallowGuests", "couldn't find user roles in context", nil, "Couldn't verify user roles.", http.StatusInternalServerError)
	}

	userRoleTypes := rawUserRoleTypes.(map[string]bool)

	if userRoleTypes[RoleTypeGuest] {
		return nil, model.NewAppError("DisallowGuests", "", nil, "Guest users are not allowed to perform this operation.", http.StatusForbidden)
	}

	return r, nil
}

// HandlePermissionSchema middleware checks for effective-channel-admin permission
// when permission schema is enabled.
func HandlePermissionSchema(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError) {
	if !config.GetConfig().PermissionSchemaEnabled {
		return r, nil
	}

	rawUserRoleTypes := r.Context().Value(CtxKeyUserRoles)
	if rawUserRoleTypes == nil {
		return nil, model.NewAppError("DisallowGuests", "couldn't find user roles in context", nil, "Couldn't verify user roles.", http.StatusInternalServerError)
	}

	userRoleType := rawUserRoleTypes.(map[string]bool)

	if !userRoleType[RoleTypeEffectiveChannelAdmin] {
		return r, model.NewAppError("HandlePermissionSchema", "", map[string]interface{}{"userID": r.Context().Value(CtxKeyUserID)}, "You do not have permission to perform this operation.", http.StatusForbidden)
	}

	return r, nil
}
