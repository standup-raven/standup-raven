package standup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/util"
)

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
	otime.DefaultLocation = location

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

	otime.DefaultLocation = location

	config.SetConfig(mockConfig)

	windowOpenTime, _ := otime.Parse("13:05")
	windowCloseTime, _ := otime.Parse("13:55")

	standupConfig := Config{
		ChannelID:       "channel_id",
		WindowOpenTime:  windowOpenTime,
		WindowCloseTime: windowCloseTime,
		Enabled:         true,
		Members:         []string{"user_id_1", "user_id_2"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Sections:        []string{"section 1", "section 2"},
		Timezone:        "Asia/Kolkata",
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
	}

	rule, err := util.ParseRRuleFromString(standupConfig.RRuleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRule = rule

	assert.Nil(t, standupConfig.IsValid(), "should be valid")

	standupConfig.ChannelID = ""
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as channel ID is empty")
	standupConfig.ChannelID = "channel_id"

	standupConfig.WindowOpenTime = otime.OTime{}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as window open time is empty")

	standupConfig.WindowOpenTime = windowOpenTime
	standupConfig.WindowCloseTime = otime.OTime{}
	assert.NotNil(t, standupConfig.IsValid(), "should be invalid as window close time is empty")

	standupConfig.ChannelID = "channel_id"
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

	standupConfig.Sections = []string{"section_1", "section_2"}
	standupConfig.RRule.Freq = rrule.WEEKLY
	standupConfig.RRule.OrigOptions.Byweekday = []rrule.Weekday{}
	assert.NotNil(t, standupConfig.IsValid(), "should not be valid as no days are specified with weekly standup")
	standupConfig.RRule = rule

	// testing invalid timezone
	standupConfig.Timezone = "Invalid-timezone"
	assert.NotNil(t, standupConfig.IsValid(), "should not be valid as specified timezone is invalid")

	// testing with  empty timezone
	standupConfig.Timezone = ""
	assert.NotNil(t, standupConfig.IsValid(), "should not be valid as specified timezone is empty")
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

	standupConfig := Config{
		ChannelID:       "channel_id",
		WindowOpenTime:  windowOpenTime,
		WindowCloseTime: windowCloseTime,
		Enabled:         true,
		Members:         []string{"user_id_1", "user_id_2"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Sections:        []string{"section 1", "section 2"},
	}

	standupConfig.ToJSON()
}

func TestStandupConfig_PreSave(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	config.Mattermost = mockAPI

	istanbul, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		t.Fatal("istanbul should have loaded successfully", err)
		return
	}

	startDate := time.Date(2020, time.July, 9, 5, 28, 0, 0, istanbul)

	standupConfig := Config{
		Timezone:    "Asia/Kolkata",
		StartDate:   startDate,
		RRuleString: "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH;COUNT=4",
	}

	assert.Nil(t, standupConfig.RRule)

	err = standupConfig.PreSave()

	assert.Nil(t, err)
	assert.Equal(t, "Asia/Kolkata", standupConfig.StartDate.Location().String())
	assert.NotNil(t, standupConfig.RRule)
	assert.Equal(t, rrule.WEEKLY, standupConfig.RRule.Freq)
	assert.Equal(t, 4, len(standupConfig.RRule.Byweekday))
	assert.Equal(t, 1, standupConfig.RRule.Interval)
	assert.Equal(t, 4, standupConfig.RRule.Count)
	assert.Equal(t, 4, len(standupConfig.RRule.All()))

	now := time.Now()
	for _, timeset := range standupConfig.RRule.Timeset {
		assert.Equal(t, now.Year(), timeset.Year())
		assert.Equal(t, now.Month(), timeset.Month())
		assert.Equal(t, now.Day(), timeset.Day())
	}

	// With invalid timezone
	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	standupConfig.Timezone = "Invalid/Timezone"
	assert.NotNil(t, standupConfig.PreSave())
	standupConfig.Timezone = "Asia/Kolkata"

	// with invalid rrule
	standupConfig.RRuleString = "invalid rrule string"
	assert.NotNil(t, standupConfig.PreSave())
	standupConfig.RRuleString = "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH;COUNT=4"
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
	})
	defer monkey.Unpatch(GetStandupChannels)

	monkey.Patch(setStandupChannels, func(channels map[string]string) error {
		return nil
	})
	defer monkey.Unpatch(setStandupChannels)

	assert.Nil(t, AddStandupChannel("channel_2"), "should not produce error")

	monkey.Patch(setStandupChannels, func(channels map[string]string) error {
		return errors.New("")
	})
	assert.NotNil(t, AddStandupChannel("channel_2"), "should produce error as couldn't save standup channels")

	monkey.Patch(GetStandupChannels, func() (map[string]string, error) {
		return nil, errors.New("")
	})
	assert.NotNil(t, AddStandupChannel("channel_2"), "should produce error as couldn't fetch existing standup channels")

	monkey.Patch(GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(setStandupChannels, func(channels map[string]string) error {
		return nil
	})
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

	_, err = GetStandupChannels()
	assert.NotNil(t, err, "error should have been produced as KVGet returned an invalid JSON")
}

func TestSaveUserStandup(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		windowOpenTime := otime.OTime{
			Time: otime.Now("Asia/Kolkata").Add(-55 * time.Minute),
		}
		windowCloseTime := otime.OTime{
			Time: otime.Now("Asia/Kolkata").Add(5 * time.Minute),
		}

		return &Config{
			ChannelID:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:        "Asia/Kolkata",
		}, nil
	})

	userStandup := &UserStandup{
		UserID:    "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string{
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
		UserID:    "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string{
			"section_1": {
				"task_1",
				"task_2",
			},
		},
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(userStandupBytes, nil)

	userStandup, err := GetUserStandup("user_id", "channel_id", otime.Now("Asia/Kolkata"))
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, &UserStandup{
		UserID:    "user_id",
		ChannelID: "channel_id",
		Standup: map[string]*[]string{
			"section_1": {
				"task_1",
				"task_2",
			},
		},
	}, userStandup, "both standups should be same")

	mockAPI = baseMock()
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, util.EmptyAppError())

	userStandup, err = GetUserStandup("user_id", "channel_id", otime.Now("Asia/Kolkata"))
	assert.NotNil(t, err, "error should have been produced as KVGet failed")
	assert.Nil(t, userStandup)

	mockAPI = baseMock()
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte{}, nil)
	userStandup, err = GetUserStandup("user_id", "channel_id", otime.Now("Asia/Kolkata"))
	assert.Nil(t, err, "no error should have been produced")
	assert.Nil(t, userStandup, "no user standup should have been found")
}

func TestSaveStandupConfig(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet", util.GetKeyHash("standup_config_channel_id")).Return([]byte("{\"scheduleEnabled\":false}"), nil)
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockAPI.On("GetChannel", "channel_id").Return(&model.Channel{}, nil)
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)

	standupConfig := &Config{
		ChannelID:       "channel_id",
		Members:         []string{"user_id_1"},
		WindowOpenTime:  otime.Now("Asia/Kolkata"),
		WindowCloseTime: otime.Now("Asia/Kolkata"),
		Sections:        []string{"section 1"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Enabled:         true,
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
	}

	rule, err := util.ParseRRuleFromString(standupConfig.RRuleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRule = rule

	savedStandupConfig, err := SaveStandupConfig(standupConfig)
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, standupConfig, savedStandupConfig, "both standup config should be identical")

	mockAPI = baseMock()
	mockAPI.On("KVGet", util.GetKeyHash("standup_config_channel_id")).Return([]byte("{\"scheduleEnabled\":false}"), nil)
	mockAPI.On("GetChannel", "channel_id").Return(&model.Channel{}, nil)
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(util.EmptyAppError())

	savedStandupConfig, err = SaveStandupConfig(standupConfig)
	assert.Error(t, err, "error should have been produced as KVSet failed")
	assert.Nil(t, savedStandupConfig, "no standup config should have been returned")
}

func TestSaveStandupConfig_DuplicateMembers(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVGet", util.GetKeyHash("standup_config_channel_id")).Return([]byte("{\"scheduleEnabled\":false}"), nil)
	mockAPI.On("GetChannel", "channel_id").Return(&model.Channel{}, nil)
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)

	now := otime.Now("Asia/Kolkata")

	standupConfig := &Config{
		ChannelID:       "channel_id",
		Members:         []string{"user_id_1", "user_id_1"},
		WindowOpenTime:  now,
		WindowCloseTime: now,
		Sections:        []string{"section 1"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Enabled:         true,
	}

	savedStandupConfig, err := SaveStandupConfig(standupConfig)
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, &Config{
		ChannelID:       "channel_id",
		Members:         []string{"user_id_1"},
		WindowOpenTime:  now,
		WindowCloseTime: now,
		Sections:        []string{"section 1"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Enabled:         true,
	}, savedStandupConfig, "both standup config should be identical")
}

func TestGetStandupConfig(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	standupConfigBytes, _ := json.Marshal(&Config{
		ChannelID:       "channel_id",
		Members:         []string{"user_id_1"},
		WindowOpenTime:  otime.Now("Asia/Kolkata"),
		WindowCloseTime: otime.Now("Asia/Kolkata"),
		Sections:        []string{"section 1"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Enabled:         true,
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(standupConfigBytes, nil)

	standupConfig, err := GetStandupConfig("channel_id")
	assert.Nil(t, err, "no error should have been produced")
	assert.Equal(t, "channel_id", standupConfig.ChannelID, "both standup configs should be identical")
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

func TestStandupConfig_GenerateScheduleString(t *testing.T) {
	location, err := time.LoadLocation("Asia/Kolkata")
	assert.Nil(t, err, "location should have loaded successfully")

	otime.DefaultLocation = location

	conf := &config.Configuration{
		TimeZone:                "Asia/Kolkata",
		PermissionSchemaEnabled: false,
		BotUserID:               "",
		Location:                location,
	}

	config.SetConfig(conf)

	windowOpenTime, err := otime.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}

	windowCloseTime, err := otime.Parse("15:00")
	if err != nil {
		t.Fatal(err)
	}

	standupConfig := Config{
		ChannelID:                  "",
		WindowOpenTime:             windowOpenTime,
		WindowCloseTime:            windowCloseTime,
		ReportFormat:               "",
		Members:                    nil,
		Sections:                   nil,
		Enabled:                    false,
		Timezone:                   "",
		WindowOpenReminderEnabled:  false,
		WindowCloseReminderEnabled: false,
		ScheduleEnabled:            false,
	}

	// weekly on all days
	rruleString := "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10"
	rule, err := util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule

	fmt.Println(config.GetConfig())

	standupScheduleString := standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Weekly on MO, TU, WE, TH, FR, SA, SU 10:00 to 15:00", standupScheduleString)

	// weekly, Monday to Friday
	rruleString = "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR;COUNT=10"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Weekly on MO, TU, WE, TH, FR 10:00 to 15:00", standupScheduleString)

	// weekly, Monday to Friday every alternate week
	rruleString = "FREQ=WEEKLY;INTERVAL=2;BYDAY=MO,TU,WE,TH,FR;COUNT=10"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Every 2 weeks on MO, TU, WE, TH, FR 10:00 to 15:00", standupScheduleString)

	// every month on the 1st
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYMONTHDAY=1;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the 1st 10:00 to 15:00", standupScheduleString)

	// testing ordinals - nd
	// every month on the 2nd
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYMONTHDAY=2;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the 2nd 10:00 to 15:00", standupScheduleString)

	// testing ordinals - rd
	// every month on the 3rd
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYMONTHDAY=3;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the 3rd 10:00 to 15:00", standupScheduleString)

	// testing ordinals - th
	// every month on the 4th
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYMONTHDAY=4;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the 4th 10:00 to 15:00", standupScheduleString)

	// every 3 months month on the 2nd
	rruleString = "FREQ=MONTHLY;INTERVAL=3;BYMONTHDAY=2;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Every 3 months on the 2nd 10:00 to 15:00", standupScheduleString)

	// every month on the first Monday
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=1;BYDAY=MO;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the first Monday 10:00 to 15:00", standupScheduleString)

	// every 3 months on the first Monday
	rruleString = "FREQ=MONTHLY;INTERVAL=3;BYSETPOS=1;BYDAY=MO;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Every 3 months on the first Monday 10:00 to 15:00", standupScheduleString)

	// every month on the first day
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the first day 10:00 to 15:00", standupScheduleString)

	// every month on the first day
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=1;BYDAY=MO,TU,WE,TH,FR;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the first weekday 10:00 to 15:00", standupScheduleString)

	// every month on the first weekend
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=1;BYDAY=SA,SU;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the first weekend 10:00 to 15:00", standupScheduleString)

	// every month on the last Monday
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=-1;BYDAY=MO;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the last Monday 10:00 to 15:00", standupScheduleString)

	// every 3 months on the last Monday
	rruleString = "FREQ=MONTHLY;INTERVAL=3;BYSETPOS=-1;BYDAY=MO;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Every 3 months on the last Monday 10:00 to 15:00", standupScheduleString)

	// every month on the last day
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=-1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the last day 10:00 to 15:00", standupScheduleString)

	// every month on the last day
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=-1;BYDAY=MO,TU,WE,TH,FR;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the last weekday 10:00 to 15:00", standupScheduleString)

	// every month on the last weekend
	rruleString = "FREQ=MONTHLY;INTERVAL=1;BYSETPOS=-1;BYDAY=SA,SU;COUNT=5"
	rule, err = util.ParseRRuleFromString(rruleString, time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	standupConfig.RRuleString = rruleString
	standupConfig.RRule = rule
	standupScheduleString = standupConfig.GenerateScheduleString()
	assert.Equal(t, "**Standup Schedule**: Monthly on the last weekend 10:00 to 15:00", standupScheduleString)
}

func TestUpdateChannelHeader(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		Header: "old header",
	}, nil)
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)

	windowOpenTime, err := otime.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}

	windowCloseTime, err := otime.Parse("15:00")
	if err != nil {
		t.Fatal(err)
	}

	parsedRRule, err := util.ParseRRuleFromString("FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10", time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: true,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
			RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
			RRule:           parsedRRule,
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 1)

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: false,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 2)

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: false,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
			RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
			RRule:           parsedRRule,
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: false,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 3)

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 4)

	// no existing channel standup config
	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return nil, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")

	// no existing channel header
	mockAPI = baseMock()
	config.Mattermost = mockAPI
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)
	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		Header: "",
	}, nil)

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")

	// existing standup schedule in header
	mockAPI = baseMock()
	config.Mattermost = mockAPI
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)
	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		Header: "**Standup Schedule**: Weekly on MO 10:00 to 15:00** ** | user-defined header",
	}, nil)
	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: true,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
			RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
			RRule:           parsedRRule,
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")

	// empty header
	mockAPI = baseMock()
	config.Mattermost = mockAPI
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)
	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		Header: "",
	}, nil)
	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: true,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
			RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
			RRule:           parsedRRule,
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
}

func TestUpdateChannelHeader_GetStandupConfig_Error(t *testing.T) {
	defer TearDown()
	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return nil, errors.New("error")
	})

	err := updateChannelHeader(&Config{
		ChannelID: "channel_id_1",
	})

	assert.NotNil(t, err, "error should have been produced as GetStandupConfig failed")
}

func TestUpdateChannelHeader_GetChannel_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetChannel", "channel_id_1").Return(nil, model.NewAppError("", "", nil, "", http.StatusInternalServerError))

	windowOpenTime, err := otime.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}

	windowCloseTime, err := otime.Parse("15:00")
	if err != nil {
		t.Fatal(err)
	}

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: true,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
	})

	assert.NotNil(t, err, "error should have been produced as GetChannel failed")
}

func TestUpdateChannelHeader_UpdateChannel_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		Header: "old header",
	}, nil)
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, model.NewAppError("", "", nil, "", http.StatusInternalServerError))

	windowOpenTime, err := otime.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}

	windowCloseTime, err := otime.Parse("15:00")
	if err != nil {
		t.Fatal(err)
	}

	parsedRRule, err := util.ParseRRuleFromString("FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10", time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: true,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
			RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
			RRule:           parsedRRule,
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.NotNil(t, err, "error should have been produced as UpdateChannel failed")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 1)
}

func TestUpdateChannelHeader_ExistingPipeInHeader(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		Header: "old | header",
	}, nil)
	mockAPI.On("UpdateChannel", mock.Anything).Return(nil, nil)

	windowOpenTime, err := otime.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}

	windowCloseTime, err := otime.Parse("15:00")
	if err != nil {
		t.Fatal(err)
	}

	parsedRRule, err := util.ParseRRuleFromString("FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10", time.Now().Add(-5*24*time.Hour))
	if err != nil {
		t.Fatal("Couldn't parse RRULE", err)
		return
	}

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: true,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
			RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
			RRule:           parsedRRule,
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 1)

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: false,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 2)

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return &Config{
			ScheduleEnabled: false,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Timezone:        "Asia/Kolkata",
		}, nil
	})

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: false,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 3)

	err = updateChannelHeader(&Config{
		ChannelID:       "channel_id_1",
		ScheduleEnabled: true,
		RRuleString:     "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR,SA,SU;COUNT=10",
		RRule:           parsedRRule,
	})

	assert.Nil(t, err, "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 4)
}

func TestUpdateChannelHeader_ArchivedChannel(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetChannel", "channel_id_1").Return(&model.Channel{
		DeleteAt: time.Now().Unix(),
	}, nil)

	monkey.Patch(GetStandupConfig, func(channelID string) (*Config, error) {
		return nil, nil
	})

	err := updateChannelHeader(&Config{
		ChannelID: "channel_id_1",
	})

	assert.Nil(t, err)
	mockAPI.AssertNumberOfCalls(t, "UpdateChannel", 0)
}
