package middleware

import (
	"github.com/mattermost/mattermost-server/model"
	"net/http"
)

type Middleware func(w http.ResponseWriter, r *http.Request) *model.AppError
