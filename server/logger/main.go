package logger

import (
	"github.com/getsentry/raven-go"
	"github.com/harshilsharma/standup-raven/server/config"
)

func Debug(msg string, err error, keyValuePairs ...interface{}) {
	if config.Mattermost != nil {
		errMsg := msg
		if err != nil {
			errMsg = err.Error()
		}

		config.Mattermost.LogDebug(errMsg, msg, keyValuePairs)
	}
}

func Info(msg string, err error, keyValuePairs ...interface{}) {
	if config.Mattermost != nil {
		errMsg := msg
		if err != nil {
			errMsg = err.Error()
		}

		config.Mattermost.LogInfo(errMsg, msg, keyValuePairs)
	}
}

func Error(msg string, err error, extraData map[string]interface{}) {
	raven.CaptureError(raven.WrapWithExtra(err, extraData), map[string]string{"msg": msg})

	if config.Mattermost != nil {
		errMsg := msg
		if err != nil {
			errMsg = err.Error()
		}

		config.Mattermost.LogError(errMsg, msg)
	}
}

func Warn(msg string, err error, keyValuePairs ...interface{}) {
	if config.Mattermost != nil {
		errMsg := msg
		if err != nil {
			errMsg = err.Error()
		}

		config.Mattermost.LogWarn(errMsg, msg, keyValuePairs)
	}
}
