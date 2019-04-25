package notification

import (
	"encoding/json"
	"fmt"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
	"github.com/bouk/monkey"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func baseMock() *plugintest.API {
	mockAPI := &plugintest.API{}
	config.Mattermost = mockAPI
	mockAPI.On("LogDebug", mock.AnythingOfType("string"))
	mockAPI.On("LogDebug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything)
	mockAPI.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	mockAPI.On("LogInfo", mock.AnythingOfType("string"))

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location: location,
	}

	config.SetConfig(mockConfig)

	return mockAPI
}

func TestSendNotificationsAndReports(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupChannels)

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})
	defer monkey.Unpatch(SendStandupReport)

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		if channelID == "channel_1" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  false,
				WindowCloseNotificationSent: false,
			}, nil
		} else if channelID == "channel_2" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: false,
			}, nil
		} else if channelID == "channel_3" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: true,
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})
	defer monkey.Unpatch(GetNotificationStatus)

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		} else if channelID == "channel_2" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_2",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetStandupConfig)

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		if channelID == "channel_1" {
			return nil
		} else if channelID == "channel_2" {
			return nil
		} else if channelID == "channel_3" {
			return nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil
	})
	defer monkey.Unpatch(SetNotificationStatus)

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_2" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_3" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		}

		panic(t)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetUserStandup)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_NoStandupChannels(t *testing.T) {
	mockAPI := baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) { return map[string]string{}, nil })
	defer monkey.Unpatch(standup.GetStandupChannels)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_GetStandupChannels_Error(t *testing.T) {
	baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return nil, errors.New("")
	})
	defer monkey.Unpatch(standup.GetStandupChannels)

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
}

func TestSendNotificationsAndReports_SendStandupReport_Error(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupChannels)

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return errors.New("")
	})
	defer monkey.Unpatch(SendStandupReport)

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		if channelID == "channel_1" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  false,
				WindowCloseNotificationSent: false,
			}, nil
		} else if channelID == "channel_2" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: false,
			}, nil
		} else if channelID == "channel_3" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: true,
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})
	defer monkey.Unpatch(GetNotificationStatus)

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		} else if channelID == "channel_2" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_2",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetStandupConfig)

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		if channelID == "channel_1" {
			return nil
		} else if channelID == "channel_2" {
			return nil
		} else if channelID == "channel_3" {
			return nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil
	})
	defer monkey.Unpatch(SetNotificationStatus)

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_2" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_3" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		}

		panic(t)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetUserStandup)

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_GetNotificationStatus_NoData(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupChannels)

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})
	defer monkey.Unpatch(SendStandupReport)

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})
	defer monkey.Unpatch(GetNotificationStatus)

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		} else if channelID == "channel_2" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_2",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now().Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetStandupConfig)

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		if channelID == "channel_1" {
			return nil
		} else if channelID == "channel_2" {
			return nil
		} else if channelID == "channel_3" {
			return nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil
	})
	defer monkey.Unpatch(SetNotificationStatus)

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_2" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_3" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		}

		panic(t)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetUserStandup)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_GetStandupConfig_Error(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupChannels)

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})
	defer monkey.Unpatch(SendStandupReport)

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		if channelID == "channel_1" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  false,
				WindowCloseNotificationSent: false,
			}, nil
		} else if channelID == "channel_2" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: false,
			}, nil
		} else if channelID == "channel_3" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: true,
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})
	defer monkey.Unpatch(GetNotificationStatus)

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	defer monkey.Unpatch(standup.GetStandupConfig)

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		if channelID == "channel_1" {
			return nil
		} else if channelID == "channel_2" {
			return nil
		} else if channelID == "channel_3" {
			return nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil
	})
	defer monkey.Unpatch(SetNotificationStatus)

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_2" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_3" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		}

		panic(t)
		return nil, nil
	})
	defer monkey.Unpatch(standup.GetUserStandup)

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_NotWorkDay(t *testing.T) {
	mockAPI := baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
	}

	config.SetConfig(mockConfig)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_Integration(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return(
		[]byte("{\"channel_1\": \":channel_1\"}"), nil,
	)
	mockAPI.On("KVGet", util.GetKeyHash(fmt.Sprintf("%s_%s_%s", config.CacheKeyPrefixNotificationStatus, "channel_1", util.GetCurrentDateString()))).Return(nil, nil)

	windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
	windowCloseTime := otime.OTime{otime.Now().Add(2 * time.Hour)}
	standupConfig, _ := json.Marshal(&standup.StandupConfig{
		ChannelId:       "channel_1",
		WindowOpenTime:  windowOpenTime,
		WindowCloseTime: windowCloseTime,
		Enabled:         true,
		Members:         []string{"user_id_1", "user_id_2"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Sections:        []string{"section 1", "section 2"},
	})
	mockAPI.On("KVGet", "UzFgbepiypG8qfVARBfHu154LDNiZOw7Mr6Ue4kNZrk=").Return(standupConfig, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_sendWindowCloseNotification_Error(t *testing.T) {
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(nil, util.EmptyAppError())
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now().Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now().Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupChannels)

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})
	defer monkey.Unpatch(SendStandupReport)

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
			StandupReportSent:           false,
			WindowOpenNotificationSent:  false,
			WindowCloseNotificationSent: false,
		}, nil
	})
	defer monkey.Unpatch(GetNotificationStatus)

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now().Add(1 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_2",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupConfig)

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})
	defer monkey.Unpatch(SetNotificationStatus)

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{}, nil
	})
	defer monkey.Unpatch(standup.GetUserStandup)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendStandupReport(t *testing.T) {
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName: "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName: "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now().Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId: channelID,
			WindowOpenTime: windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat: config.ReportFormatTypeAggregated,
			Sections: []string{"section_1", "section_2"},
			Members: []string{"user_id_1", "user_id_2"},
			Enabled: true,
		}, nil
	})
	defer monkey.Unpatch(standup.GetStandupConfig)

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID: userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})
	defer monkey.Unpatch(standup.GetUserStandup)

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})
	defer monkey.Unpatch(GetNotificationStatus)

	monkey.Patch(SetNotificationStatus, func (channelID string, status *ChannelNotificationStatus) error {
		return nil
	})
	defer monkey.Unpatch(SetNotificationStatus)

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now(), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now(), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now(), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now(), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now().Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now().Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId: channelID,
			WindowOpenTime: windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat: config.ReportFormatTypeAggregated,
			Sections: []string{"section_1", "section_2"},
			Members: []string{},
			Enabled: true,
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now(), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "shouldn't produce error as standup with no members is a valid case")
}
