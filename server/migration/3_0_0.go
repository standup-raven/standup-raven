package migration

import (
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"strconv"
	"strings"
	"time"
)

func upgradeDatabaseToVersion3_0_0(fromVersion, toVersion string) error {
	if fromVersion == version2_0_0 || toVersion == version3_0_0 {

		standupChannels, err := standup.GetStandupChannels()
		if err != nil {
			return err
		}

		rruleString, err := generateRRuleStringByWorkWeek()
		if err != nil {
			return err
		}

		for channelID := range standupChannels {
			channelConfig, err := standup.GetStandupConfig(channelID)
			if err != nil {
				return err
			}

			channelConfig.RRuleString = rruleString
			if err := channelConfig.PreSave(); err != nil {
				return err
			}

			if err := channelConfig.IsValid(); err != nil {
				return err
			}

			if _, err := standup.SaveStandupConfig(channelConfig); err != nil {
				return err
			}
		}
	}

	if UpdateErr := updateSchemaVersion(version3_0_0); UpdateErr != nil {
		return UpdateErr
	}
	return nil
}

func generateRRuleStringByWorkWeek() (string, error) {
	type oldConfigType struct {
		TimeZone                string `json:"timeZone"`
		WorkWeekStart           string `json:"workWeekStart"`
		WorkWeekEnd             string `json:"workWeekEnd"`
		PermissionSchemaEnabled bool   `json:"permissionSchemaEnabled"`
		EnableErrorReporting    bool   `json:"enableErrorReporting"`
		SentryServerDSN         string `json:"sentryServerDSN"`
		SentryWebappDSN         string `json:"sentryWebappDSN"`
	}

	var oldConf *oldConfigType
	err := config.Mattermost.LoadPluginConfiguration(&oldConf)
	if err != nil {
		logger.Error("couldn't fetch old configuration", err, nil)
		return "", err
	}

	workWeekStart, err := strconv.Atoi(oldConf.WorkWeekStart)
	if err != nil {
		logger.Error("Couldn't parse integer in old config work week start", err, nil)
		return "", err
	}

	workWeekEnd, err := strconv.Atoi(oldConf.WorkWeekEnd)
	if err != nil {
		logger.Error("Couldn't parse integer in old config work week end", err, nil)
		return "", err
	}

	weekdays := getStandupWeekDays(workWeekStart, workWeekEnd)
	rruleString := "FREQ=WEEKLY;INTERVAL=1;BYDAY=" + strings.Join(weekdays, ",")
	return rruleString, nil
}

func getStandupWeekDays(workWeekStart, workWeekEnd int) []string {
	weekdays := []string{}
	if workWeekStart > workWeekEnd {
		for i := workWeekStart; i <= 6; i++ {
			weekdays = append(weekdays, strings.ToUpper(time.Weekday(i).String())[:2])
		}

		for i := 0; i <= workWeekEnd; i++ {
			weekdays = append(weekdays, strings.ToUpper(time.Weekday(i).String())[:2])
		}
	} else {
		for i := workWeekStart; i <= workWeekEnd; i++ {
			weekdays = append(weekdays, strings.ToUpper(time.Weekday(i).String())[:2])
		}
	}

	return weekdays
}
