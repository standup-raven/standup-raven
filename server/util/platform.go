package util

import (
	"fmt"
	"github.com/standup-raven/standup-raven/server/config"
)

func UserIcon(userId string) string {
	return fmt.Sprintf("![User Avatar]("+config.UserIconURL+" "+config.UserIconSize+")", userId)
}
