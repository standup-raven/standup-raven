package migration

import (
	"encoding/json"
	"errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

var (
	databaseSchemaVersion = "database_schema_version"
	version1_4_0          = "1.4.0"
	version1_5_0          = "1.5.0"
	version2_0_0          = "2.0.0"
)

//DatabaseMigration gets the current database schema version and performs
//all the required data migrations.
func DatabaseMigration() error {
	if err := ensureSchemaVersion(); err != nil {
		return err
	}

	if upgradeErr := upgrade(); upgradeErr != nil {
		return upgradeErr
	}
	return nil
}

func ensureSchemaVersion() error {
	key, err := config.Mattermost.KVGet(util.GetKeyHash(databaseSchemaVersion))
	if err != nil {
		logger.Error("Couldn't fetch database schema version from KV store", err, nil)
		return err
	}
	if key == nil {
		version, err := json.Marshal(version1_4_0)
		if err != nil {
			logger.Error("Couldn't marshal database schema version", err, nil)
			return err
		}
		if appErr := config.Mattermost.KVSet(util.GetKeyHash(databaseSchemaVersion), version); appErr != nil {
			logger.Error("Couldn't update database version into KV store", appErr, nil)
			return errors.New(appErr.Error())
		}
	}
	return nil
}

func upgrade() error {
	if err := upgradeDatabaseToVersion1_5_0(); err != nil {
		return err
	}

	if err := upgradeDatabaseToVersion2_0_0(); err != nil {
		return err
	}

	return nil
}

func getCurrentSchemaVersion() (string, error) {
	key, err := config.Mattermost.KVGet(util.GetKeyHash(databaseSchemaVersion))
	if err != nil {
		logger.Error("Couldn't fetch database schema version from KV store", err, nil)
		return "", err
	}
	var version string
	appErr := json.Unmarshal(key, &version)
	if appErr != nil {
		logger.Error("Couldn't marshal database schema version", appErr, nil)
		return "", appErr
	}
	return version, nil
}

func updateSchemaVersion(version string) error {
	newVersion, marshalErr := json.Marshal(version)
	if marshalErr != nil {
		logger.Error("Couldn't marshal database schema version", marshalErr, nil)
		return marshalErr
	}
	if appErr := config.Mattermost.KVSet(util.GetKeyHash(databaseSchemaVersion), newVersion); appErr != nil {
		logger.Error("Couldn't update database version into KV store", appErr, nil)
		return errors.New(appErr.Error())
	}
	return nil
}

func upgradeDatabaseToVersion1_5_0() error {
	version, versionErr := getCurrentSchemaVersion()
	if versionErr != nil {
		return versionErr
	}
	if version == version1_4_0 {
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

		if UpdateErr := updateSchemaVersion(version1_5_0); UpdateErr != nil {
			return UpdateErr
		}
	}
	return nil
}

func upgradeDatabaseToVersion2_0_0() error {
	version, versionErr := getCurrentSchemaVersion()
	if versionErr != nil {
		return versionErr
	}
	if version == version1_5_0 {
		// TODO uncomment before release
		//if UpdateErr := updateSchemaVersion(version1_5_0); UpdateErr != nil {
		//	return UpdateErr
		//}
	}
	return nil
}
