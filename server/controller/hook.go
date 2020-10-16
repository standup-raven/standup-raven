package controller

import (
	"fmt"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller/middleware"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup/notification"
	"net/http"
)

var hook = &Endpoint{
	Path:    "/hook",
	Method:  http.MethodGet,
	Execute: executeHook,
	Middlewares: []middleware.Middleware{
		middleware.Authenticate,
	},
}

func executeHook(w http.ResponseWriter, r *http.Request) error {
	logger.Debug(fmt.Sprintf("Fetching notification status for channel: %s", otime.Now(config.GetConfig().TimeZone).String()), nil)
	notification.SendNotificationsAndReports()
	return nil
}
