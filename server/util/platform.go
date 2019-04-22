package util

import (
	"fmt"
	"github.com/harshilsharma/standup-raven/server/config"
)

func UserIcon(userId string) string {
	return fmt.Sprintf("![]("+config.UserIconURL+" "+config.UserIconSize+")", userId)
}
