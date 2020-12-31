package middleware

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

// Middleware type implements any logic required to be performed
// before the endpoint implementation is executed.
type Middleware func(w http.ResponseWriter, r *http.Request) (*http.Request, *model.AppError)
