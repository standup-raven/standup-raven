package migration

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
)

func upgradeDatabaseToVersion3_0_0(fromVersion string) error {
	if fromVersion != version2_0_0 {
		return nil
	}

	standupChannels, err := standup.GetStandupChannels()
	if err != nil {
		return err
	}

	rruleString, err := generateRRuleStringByWorkWeek()
	if err != nil {
		return err
	}

	config.Mattermost.LogInfo("Upgrading Standup Raven to v3.0.0. Standup Raven will automatically be disabled for channels which fail to be upgraded.")

	failedChannelIDs := []string{}

	for channelID := range standupChannels {
		err := upgradeChannel(channelID, rruleString)

		if err == nil {
			config.Mattermost.LogInfo(fmt.Sprintf("Successfully upgraded channel '%s'.", channelID))
		} else {
			config.Mattermost.LogError(fmt.Sprintf("Failed to upgrade channel '%s'. Error: %s.", channelID, err))
			failedChannelIDs = append(failedChannelIDs, channelID)
		}
	}

	if len(failedChannelIDs) > 0 {
		if err := standup.RemoveStandupChannels(failedChannelIDs); err != nil {
			return err
		}

		for _, channelID := range failedChannelIDs {
			_ = standup.ArchiveStandupChannels(channelID)
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
		logger.Error("Couldn't parse integer in old config work week start, defaulting to using 1 (Monday)", err, nil)
		workWeekStart = 1
	}

	workWeekEnd, err := strconv.Atoi(oldConf.WorkWeekEnd)
	if err != nil {
		logger.Error("Couldn't parse integer in old config work week end, defaulting to using 5 (Friday)", err, nil)
		workWeekEnd = 5
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

func upgradeChannel(channelID, rruleString string) error {
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

	return nil
}
