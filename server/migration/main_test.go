package migration

import (
	"encoding/json"
	"errors"
	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func baseMock() *plugintest.API {
	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI
	
	monkey.Patch(logger.Debug, func(msg string, err error, keyValuePairs ...interface{}) {})
	monkey.Patch(logger.Error, func(msg string, err error, extraData map[string]interface{}) {})
	
	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location: location,
		PluginVersion: version3_0_0,
	}
	config.SetConfig(mockConfig)
	return mockAPI
}

func TearDown() {
	monkey.UnpatchAll()
}

func TestDatebaseMigration_getCurrentSchemaVersion_Error(t *testing.T) {
	defer TearDown()
	baseMock()
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "", errors.New("")
	})
	
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_KVGet_error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 0))
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_EnsureSchemaVersion_error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 0))
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_KVSet_error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return(nil,nil)
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return( model.NewAppError("", "", nil, "", 0))
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_JsonMarshal_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return(nil, nil)
	monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, errors.New("")
	})
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_JsonUnmarshal_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return([]byte("1.4.0"), nil)
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	monkey.Patch(json.Unmarshal, func(data []byte, v interface{}) error{
		return errors.New("")
	})
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_GetStandupChannels_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return([]byte("1.4.0"), nil)
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return nil, errors.New("")
	})
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_GetStandupConfig_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return([]byte("1.4.0"), nil)
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil,errors.New("")
	})
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_GetStandupConfig_Nil(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
	
	conf := config.GetConfig()
	conf.PluginVersion = version1_5_0
	config.SetConfig(conf)
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil,nil
	})
	err := DatabaseMigration()
	assert.Nil(t, err)
}

func TestDatabaseMigration_SaveStandupConfig_Error(t *testing.T) {
	defer TearDown()
	baseMock()
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_2",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	
	monkey.Patch(standup.SaveStandupConfig, func(standupConfig *standup.StandupConfig) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("KVGet","mS93mHcYcKvlwjnt1DvUXsRwcIuoOO+mKsCRNZl/Ht4=", mock.Anything).Return([]byte("1.4.0"), nil)
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)

	conf := config.GetConfig()
	conf.PluginVersion = version1_5_0
	config.SetConfig(conf)
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil,nil
	})
	err := DatabaseMigration()
	assert.Nil(t, err)
}

func TestDatabaseMigration_updateSchemaVersion_Error(t *testing.T) {
	defer TearDown()
	baseMock()
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil,nil
	})
	monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, errors.New("")
	})
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_updateSchemaVersion_KVSet_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.Anything, mock.Anything).Return( model.NewAppError("", "", nil, "", 0))
	
	monkey.Patch(getCurrentSchemaVersion, func() (string, error) {
		return "1.4.0", nil
	})
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil,nil
	})
	err := DatabaseMigration()
	assert.NotNil(t, err)
}

func TestDatabaseMigration_getSChemaVersion_Error(t *testing.T) {
	defer TearDown()
	baseMock()
	
	monkey.Patch(getCurrentSchemaVersion, func () (string, error) {
		return "", errors.New("")
	})

	conf := config.GetConfig()
	conf.PluginVersion = version1_5_0
	config.SetConfig(conf)
	
	err := DatabaseMigration()
	assert.NotNil(t, err)
}
