package migration

import (
	"encoding/json"
	"errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

var(
	databaseSchemaVersion = "database_schema_version"
	VERSION_1_4_0 = "1.4.0"
	VERSION_1_5_0 =	"1.5.0"
	OLD_VERSION = "1.4.0"
) 

func DatabaseMigration() error {
	key,err := config.Mattermost.KVGet(util.GetKeyHash(databaseSchemaVersion))
	if err!=nil {
		logger.Error("Couldn't fetch database schema version from KV store", err, nil)
		return err
	}
	if key == nil {
		version, err := json.Marshal(OLD_VERSION)
		if err != nil {
			logger.Error("Couldn't marshal database schema version", err, nil)
			return err
		}
		if appErr := config.Mattermost.KVSet(util.GetKeyHash(databaseSchemaVersion), version); appErr != nil {
			logger.Error("Couldn't update database version into KV store", appErr, nil)
			return errors.New(appErr.Error())
		}
	}
	upgradeDatabaseToversion15()
	return nil
}

func upgradeDatabaseToversion15() error{
	key,err := config.Mattermost.KVGet(util.GetKeyHash(databaseSchemaVersion))
	if err!=nil {
		logger.Error("Couldn't fetch database schema version from KV store", err, nil)
		return err
	}
	var version string
	appErr := json.Unmarshal(key,&version)
	if appErr != nil {
		logger.Error("Couldn't marshal database schema version", appErr, nil)
			return appErr
	}
	if version == VERSION_1_4_0 {
		channelIDs, err := standup.GetStandupChannels()
		if err != nil {
			return err
		}
		for channelID := range channelIDs {
			standupConfig, err := standup.GetStandupConfig(channelID)
			if err != nil {
				return err
			}
			if standupConfig == nil {
				logger.Error("Unable to find standup config for channel", nil, map[string]interface{}{"channelID": channelID})
				continue
			}
	
			if !standupConfig.Enabled {
				continue
			}
			standupConfig.Timezone = config.GetConfig().TimeZone
			standup.SaveStandupConfig(standupConfig)
		}
	}
	newVersion, marshalErr := json.Marshal(VERSION_1_5_0)
	if marshalErr != nil {
		logger.Error("Couldn't marshal database schema version", marshalErr, nil)
		return err
	}
	if appErr := config.Mattermost.KVSet(util.GetKeyHash(databaseSchemaVersion), newVersion); appErr != nil {
		logger.Error("Couldn't update database version into KV store", appErr, nil)
		return errors.New(appErr.Error())
	}
	return nil
}
