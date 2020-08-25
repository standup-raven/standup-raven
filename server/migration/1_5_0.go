package migration

import (
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
)

func upgradeDatabaseToVersion1_5_0(fromVersion string) error {
	if fromVersion == version1_4_0 {
		channelIDs, err := standup.GetStandupChannels()
		if err != nil {
			return err
		}

		defaultTimezone := config.GetConfig().TimeZone
		for channelID := range channelIDs {
			standupConfig, err := standup.GetStandupConfig(channelID)
			if err != nil {
				return err
			}
			if standupConfig == nil {
				logger.Error("Unable to find standup config for channel", nil, map[string]interface{}{"channelID": channelID})
				continue
			}

			standupConfig.Timezone = defaultTimezone
			standupConfig.WindowOpenReminderEnabled = true
			standupConfig.WindowCloseReminderEnabled = true
			if _, configErr := standup.SaveStandupConfig(standupConfig); configErr != nil {
				return configErr
			}
		}
	}

	if UpdateErr := updateSchemaVersion(version1_5_0); UpdateErr != nil {
		return UpdateErr
	}

	return nil
}
