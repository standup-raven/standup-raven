package migration

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/util"
	"github.com/thoas/go-funk"
	"strings"
)

var (
	databaseSchemaVersion = "database_schema_version"
	versionNone           = ""
	version1_4_0          = "1.4.0"
	version1_5_0          = "1.5.0"
	version2_0_0          = "2.0.0"
	version3_0_0          = "3.0.0"
)

var upgradeCompatibility = map[string][]string{
	versionNone:  {},
	version1_5_0: {version1_4_0},
	version3_0_0: {version2_0_0, version1_5_0},
}

type Migratioon func(fromVersion, toVersion string) error

var migrations = []Migratioon{
	upgradeDatabaseToVersion1_5_0,
	upgradeDatabaseToVersion2_0_0,
	upgradeDatabaseToVersion3_0_0,
}

//DatabaseMigration gets the current database schema version and performs
//all the required data migrations.
func DatabaseMigration() error {
	pluginVersion := config.GetConfig().PluginVersion

	schemaVersion, appErr := getCurrentSchemaVersion()
	if appErr != nil {
		return appErr
	}

	if !isUpgradeCompatible(schemaVersion, pluginVersion) {
		msg := fmt.Sprintf(
			"Cannot upgrade Standup Raven from version %s to %s. Please upgrade first to one of versions %s",
			schemaVersion,
			pluginVersion,
			strings.Join(upgradeCompatibility[pluginVersion], ", "),
		)

		logger.Error(msg, nil, nil)
		return errors.New(msg)
	}

	for _, migration := range migrations {
		schemaVersion, err := getCurrentSchemaVersion()
		if err != nil {
			return err
		}

		if err := migration(schemaVersion, pluginVersion); err != nil {
			return err
		}
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

func isUpgradeCompatible(fromVersion, toVersion string) bool {
	return fromVersion == versionNone || funk.Contains(upgradeCompatibility[toVersion], fromVersion)
}
