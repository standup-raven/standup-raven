package utils

import "github.com/mattermost/mattermost-server/v5/model"

func EmptyAppError() *model.AppError {
	return model.NewAppError("", "", nil, "", 0)
}
