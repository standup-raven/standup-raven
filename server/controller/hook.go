package controller

import (
	"fmt"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup/notification"
	"net/http"
)

var hook = &Endpoint{
	Path:         "/hook",
	Method:       http.MethodGet,
	Execute:      executeHook,
	RequiresAuth: true,
}

func executeHook(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("hook activated at " + otime.Now().String())

	err := notification.SendNotificationsAndReports()
	if err != nil {
		return err
	}

	return nil
}
