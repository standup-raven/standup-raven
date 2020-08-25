package migration

import (
	"bou.ke/monkey"
	"errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpgradeDatabaseToVersion1_5_0(t *testing.T) {
	defer TearDown()
	baseMock()
	
	// incompatible upgrade path
	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})
	
	err := upgradeDatabaseToVersion1_5_0(version2_0_0)
	assert.Nil(t, err)
	assert.Equal(t, 1, updateSchemaVersionCount)


	updateSchemaVersionCount = 0
	monkey.Patch(standup.GetStandupChannels, func () (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	
	monkey.Patch(config.GetConfig, func () *config.Configuration {
		return &config.Configuration{
			TimeZone: "Asia/Kolkata",
		}
	})
	
	monkey.Patch(standup.GetStandupConfig, func (channelID string) (*standup.StandupConfig, error) {
		return &standup.StandupConfig{}, nil
	})
	
	monkey.Patch(standup.SaveStandupConfig, func (standupConfig *standup.StandupConfig) (*standup.StandupConfig, error) {
		return &standup.StandupConfig{}, nil
	})
	
	err = upgradeDatabaseToVersion1_5_0(version1_4_0)
	assert.Nil(t, err)
}

func TestUpgradeDatabaseToVersion1_5_0_GetStandupChannels_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	// incompatible upgrade path
	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})
	
	monkey.Patch(standup.GetStandupChannels, func () (map[string]string, error) {
		return nil, errors.New("")
	})
	
	err := upgradeDatabaseToVersion1_5_0(version1_4_0)
	assert.NotNil(t, err)

}

func TestUpgradeDatabaseToVersion1_5_0_GetStandupConfig_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	// incompatible upgrade path
	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})

	monkey.Patch(standup.GetStandupChannels, func () (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(config.GetConfig, func () *config.Configuration {
		return &config.Configuration{
			TimeZone: "Asia/Kolkata",
		}
	})

	monkey.Patch(standup.GetStandupConfig, func (channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})

	err := upgradeDatabaseToVersion1_5_0(version1_4_0)
	assert.NotNil(t, err)
}

func TestUpgradeDatabaseToVersion1_5_0_SaveStandupConfig_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	// incompatible upgrade path
	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return nil
	})
	
	monkey.Patch(standup.GetStandupChannels, func () (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(config.GetConfig, func () *config.Configuration {
		return &config.Configuration{
			TimeZone: "Asia/Kolkata",
		}
	})

	monkey.Patch(standup.GetStandupConfig, func (channelID string) (*standup.StandupConfig, error) {
		return &standup.StandupConfig{}, nil
	})

	monkey.Patch(standup.SaveStandupConfig, func (standupConfig *standup.StandupConfig) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})

	err := upgradeDatabaseToVersion1_5_0(version1_4_0)
	assert.NotNil(t, err)
}

func TestUpgradeDatabaseToVersion1_5_0_updateSchemaVersion_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	// incompatible upgrade path
	updateSchemaVersionCount := 0
	monkey.Patch(updateSchemaVersion, func(version string) error {
		updateSchemaVersionCount++
		return errors.New("")
	})
	
	monkey.Patch(standup.GetStandupChannels, func () (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(config.GetConfig, func () *config.Configuration {
		return &config.Configuration{
			TimeZone: "Asia/Kolkata",
		}
	})

	monkey.Patch(standup.GetStandupConfig, func (channelID string) (*standup.StandupConfig, error) {
		return &standup.StandupConfig{}, nil
	})

	monkey.Patch(standup.SaveStandupConfig, func (standupConfig *standup.StandupConfig) (*standup.StandupConfig, error) {
		return &standup.StandupConfig{}, nil
	})

	err := upgradeDatabaseToVersion1_5_0(version1_4_0)
	assert.NotNil(t, err)
}
