package standup

import (
	"encoding/json"
	"github.com/bouk/monkey"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	"github.com/pkg/errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/util"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"
import "github.com/mattermost/mattermost-server/plugin/plugintest"

func baseMock() *plugintest.API {
	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI

	monkey.Patch(logger.Debug, func(msg string, err error, keyValuePairs ...interface{}) {})
	monkey.Patch(logger.Error, func(msg string, err error, extraData map[string]interface{}) {})
	monkey.Patch(logger.Info, func(msg string, err error, keyValuePairs ...interface{}) {})
	monkey.Patch(logger.Warn, func(msg string, err error, keyValuePairs ...interface{}) {})

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location: location,
	}

	config.SetConfig(mockConfig)

	return mockAPI
}

func TearDown() {
	monkey.UnpatchAll()
}

func TestUserStandup_IsValid(t *testing.T) {
	defer TearDown()
	userStandup := UserStandup{
		UserID:    "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string{
			"section_1": {"task_1", "task_2"},
		},
	}

	mockAPI := baseMock()
	config.Mattermost = mockAPI

	mockAPI.On("GetChannel", "channel_id").Return(&model.Channel{}, nil)
	mockAPI.On("GetChannel", "non_existing_channel_id").Return(nil, model.NewAppError("", "", nil, "", 0))

	assert.Nil(t, userStandup.IsValid(), "should be valid")

	userStandup.UserID = ""
	assert.NotNil(t, userStandup.IsValid(), "should be invalid as user ID is empty")
	userStandup.UserID = "user_id"

	userStandup.ChannelID = ""
	assert.NotNil(t, userStandup.IsValid(), "should be invalid as channel ID is empty")
	userStandup.ChannelID = "channel_id"

	userStandup.ChannelID = "non_existing_channel_id"
	assert.NotNil(t, userStandup.IsValid(), "should be invalid as channel doesn exist")
	userStandup.ChannelID = "channel_id"

	userStandup.Standup = nil
	assert.NotNil(t, userStandup.IsValid(), "should be invalid as standup are nil")

	userStandup.Standup = map[string]*[]string{}
	assert.NotNil(t, userStandup.IsValid(), "should be invalid as standup contain no entries")

	userStandup.Standup = map[string]*[]string{
		"section_1": {},
	}
	assert.NotNil(t, userStandup.IsValid(), "should be invalid as standup contain entries with no tasks")
}

func TestStandupConfig_IsValid(t *testing.T) {
	defer TearDown()
	// mocking Mattermost API
	mockAPI := baseMock()
	config.Mattermost = mockAPI

	// mocking plugin config
	location, err := time.LoadLocation("Asia/Kolkata")
	assert.NoError(t, err, "Time zone should parse successfully")

	mockConfig := &config.Configuration{
		Location: location,
	}

	config.SetConfig(mockConfig)

	windowOpenTime, _ := otime.Parse("13:05")
	windowCloseTime, _ := otime.Parse("13:55")

	standupConfig := StandupConfig{
		ChannelId:       "channel_id",
		WindowOpenTime:  windowOpenTime,
		WindowCloseTime: windowCloseTime,
		Enabled:         true,
		Members:         []string{"user_id_1", "user_id_2"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Sections:        []string{"section 1", "section 2"},
	}

	assert.Nil(t, standupConfig.IsValid(), "should be valid")

	standupConfig.ChannelId = ""
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as channel ID is empty")
	standupConfig.ChannelId = "channel_id"

	standupConfig.WindowOpenTime = otime.OTime{}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as window open time is empty")

	standupConfig.WindowOpenTime = windowOpenTime
	standupConfig.WindowCloseTime = otime.OTime{}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as window close time is empty")

	standupConfig.ChannelId = "channel_id"
	standupConfig.WindowOpenTime, _ = otime.Parse("10:00")
	standupConfig.WindowCloseTime, _ = otime.Parse("09:00")
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as window open time is after window close time")

	standupConfig.WindowOpenTime = windowOpenTime
	standupConfig.WindowCloseTime = windowCloseTime
	standupConfig.Members = nil
	assert.Nil(t, standupConfig.IsValid(), "should be valid as lack of members is a valid state")

	standupConfig.Members = []string{}
	assert.Nil(t, standupConfig.IsValid(), "should be valid as lack of members is a valid state")

	standupConfig.Members = []string{"member_1", "member_1"}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as duplicate members are added")
	standupConfig.Members = []string{"member_1"}

	standupConfig.ReportFormat = "invalid_report_format"
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid is report format is not one of the allowed values")

	standupConfig.ReportFormat = config.ReportFormatTypeAggregated
	standupConfig.Sections = nil
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as sections are nil")

	standupConfig.Sections = []string{}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as sections are empty")

	standupConfig.Sections = []string{"section_1", "section_1"}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as duplicate sections are added")
}

func TestStandupConfig_ToJson(t *testing.T) {
	defer TearDown()
	// mocking Mattermost API
	mockAPI := baseMock()
	config.Mattermost = mockAPI

	// mocking plugin config
	location, err := time.LoadLocation("Asia/Kolkata")
	assert.NoError(t, err, "Time zone should parse successfully")

	mockConfig := &config.Configuration{
		Location: location,
	}

	config.SetConfig(mockConfig)

	windowOpenTime, _ := otime.Parse("13:05")
	windowCloseTime, _ := otime.Parse("13:55")

	standupConfig := StandupConfig{
		ChannelId:       "channel_id",
		WindowOpenTime:  windowOpenTime,
		WindowCloseTime: windowCloseTime,
		Enabled:         true,
		Members:         []string{"user_id_1", "user_id_2"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Sections:        []string{"section 1", "section 2"},
	}

	standupConfig.ToJson()
}

func TestAddStandupChannel(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	config.Mattermost = mockAPI

	mockAPI.On("LogDebug", mock.AnythingOfType("string"))

	monkey.Patch(GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	}, )
	defer monkey.Unpatch(GetStandupChannels)

	monkey.Patch(setStandupChannels, func(channels map[string]string) error {
		return nil
	}, )
	defer monkey.Unpatch(setStandupChannels)

	assert.Nil(t, AddStandupChannel("channel_2"), "should not produce error")

	monkey.Patch(setStandupChannels, func(channels map[string]string) error {
		return errors.New("")
	}, )
	assert.NotNil(t, AddStandupChannel("channel_2"), "should produce error as couldn't save standup channels")

	monkey.Patch(GetStandupChannels, func() (map[string]string, error) {
		return nil, errors.New("")
	}, )
	assert.NotNil(t, AddStandupChannel("channel_2"), "should produce error as couldn't fetch existing standup channels")

	monkey.Patch(GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	}, )

	monkey.Patch(setStandupChannels, func(channels map[string]string) error {
		return nil
	}, )
	assert.Nil(t, AddStandupChannel("channel_1"), "shouldn't fail if adding a channel which was already a standup channel")
}

func TestAddStandupChannel_IntegrationTest(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	config.Mattermost = mockAPI

	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return([]byte("{\"channel_1\":\"channel_1\"}"), nil)
	mockAPI.On("KVSet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=", mock.Anything).Return(nil)

	assert.Nil(t, AddStandupChannel("channel_2"), "should not produce error")

	mockAPI = &plugintest.API{}
	config.Mattermost = mockAPI

	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return([]byte("{\"channel_1\":\"channel_1\"}"), nil)
	mockAPI.On("KVSet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=", mock.Anything).Return(model.NewAppError("", "", nil, "", 0))
	assert.NotNil(t, AddStandupChannel("channel_2"), "should produce error as couldn't save standup channels")

	mockAPI = &plugintest.API{}
	config.Mattermost = mockAPI

	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return(nil, model.NewAppError("", "", nil, "", 0))
	assert.NotNil(t, AddStandupChannel("channel_2"), "should produce error as couldn't fetch existing standup channels")

	mockAPI = &plugintest.API{}
	config.Mattermost = mockAPI

	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return([]byte("{\"channel_1\":\"channel_1\"}"), nil)
	mockAPI.On("KVSet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=", mock.Anything).Return(nil)
	assert.Nil(t, AddStandupChannel("channel_1"), "should fail if adding a channel which was already a standup channel")
}

func TestGetStandupChannels(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return([]byte("{\"channel_1\":\"channel_1\"}"), nil)

	standupChannels, err := GetStandupChannels()
	assert.Nil(t, err, "no error should be produced")
	assert.Equal(t, standupChannels, map[string]string{"channel_1": "channel_1"}, "one standup channel, 'channel_1' should be returned")

	mockAPI = &plugintest.API{}
	config.Mattermost = mockAPI

	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return(nil, model.NewAppError("", "", nil, "", 0))

	standupChannels, err = GetStandupChannels()
	assert.NotNil(t, err, "error should be produced as KVGet failed")
	assert.Nil(t, standupChannels, "no standup channels should be returned")

	mockAPI = &plugintest.API{}
	config.Mattermost = mockAPI
	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return(nil, nil)

	standupChannels, err = GetStandupChannels()
	assert.Nil(t, err, "no error should be produced")
	assert.Equal(t, standupChannels, map[string]string{}, "empty standup channel map should have been returned")

	mockAPI = &plugintest.API{}
	config.Mattermost = mockAPI
	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return([]byte("{\"channel_1\":\"channel_1\""), nil)

	standupChannels, err = GetStandupChannels()
	assert.NotNil(t, err, "error should have been produced as KVGet returned an invalid JSON")
}

func TestSaveUserStandup(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)

	userStandup := &UserStandup{
		UserID: "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string {
			"section_1": {
				"task_1",
				"task_2",
			},
		},
	}

	assert.Nil(t, SaveUserStandup(userStandup), "should not return any error")

	mockAPI = baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(util.EmptyAppError())

	assert.NotNil(t, SaveUserStandup(userStandup), "should return error as KVSet failed")
}

func TestGetUserStandup(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	userStandupBytes, _ := json.Marshal(&UserStandup{
		UserID: "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string {
			"section_1": {
				"task_1",
				"task_2",
			},
		},
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(userStandupBytes, nil)

	userStandup, err := GetUserStandup("user_id", "channel_id", otime.Now())
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, &UserStandup{
		UserID: "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string {
			"section_1": {
				"task_1",
				"task_2",
			},
		},
	}, userStandup, "both standups should be same")

	mockAPI = baseMock()
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, util.EmptyAppError())

	userStandup, err = GetUserStandup("user_id", "channel_id", otime.Now())
	assert.NotNil(t, err, "error should have been produced as KVGet failed")
	assert.Nil(t, userStandup)

	mockAPI = baseMock()
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte{}, nil)
	userStandup, err = GetUserStandup("user_id", "channel_id", otime.Now())
	assert.Nil(t, err, "no error should have been produced")
	assert.Nil(t, userStandup, "no user standup should have been found")

}

func TestSaveStandupConfig(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)

	standupConfig := &StandupConfig{
		ChannelId: "channel_id",
		 Members: []string{"user_id_1"},
		 WindowOpenTime: otime.Now(),
		 WindowCloseTime: otime.Now(),
		 Sections: []string{"section 1"},
		 ReportFormat: config.ReportFormatUserAggregated,
		 Enabled: true,
	}

	savedStandupConfig, err := SaveStandupConfig(standupConfig)
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, standupConfig, savedStandupConfig, "both standup config should be identical")

	mockAPI = baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(util.EmptyAppError())

	savedStandupConfig, err = SaveStandupConfig(standupConfig)
	assert.Error(t, err, "error should have been produced as KVSet failed")
	assert.Nil(t, savedStandupConfig, "no standup config should have been returned")


}

func TestGetStandupConfig(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	standupConfigBytes, _ := json.Marshal(&StandupConfig{
		ChannelId: "channel_id",
		Members: []string{"user_id_1"},
		WindowOpenTime: otime.Now(),
		WindowCloseTime: otime.Now(),
		Sections: []string{"section 1"},
		ReportFormat: config.ReportFormatUserAggregated,
		Enabled: true,
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(standupConfigBytes, nil)

	standupConfig, err := GetStandupConfig("channel_id")
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, "channel_id", standupConfig.ChannelId, "both standup configs should be identical")
	assert.Equal(t, []string{"user_id_1"}, standupConfig.Members, "both standup configs should be identical")
	assert.Equal(t, []string{"section 1"}, standupConfig.Sections, "both standup configs should be identical")
	assert.Equal(t, config.ReportFormatUserAggregated, standupConfig.ReportFormat, "both standup configs should be identical")
	assert.Equal(t, true, standupConfig.Enabled, "both standup configs should be identical")

	mockAPI = baseMock()
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, util.EmptyAppError())
	standupConfig, err = GetStandupConfig("channel_id")
	assert.NotNil(t, err, "error should have been produced as KVGet failed")
	assert.Nil(t, standupConfig, "no standup config should have been returned")

	mockAPI = baseMock()
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte{}, nil)
	standupConfig, err = GetStandupConfig("channel_id")
	assert.Nil(t, err, "no error should have been produced as lack of standup config is valid scenerio")
	assert.Nil(t, standupConfig, "no standup config should have been returned")

}
