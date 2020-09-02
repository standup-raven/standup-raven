package middleware

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"net/http"
)

// Middleware type implements any logic required to be performed
// before the endpoint implementation is executed.
type Middleware func(w http.ResponseWriter, r *http.Request) *model.AppError
