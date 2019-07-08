package notification

import (
	"encoding/json"
	"fmt"
	"github.com/bouk/monkey"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	"github.com/pkg/errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
	"time"
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

	return mockAPI
}

func TearDown() {
	monkey.UnpatchAll()
}

func TestSendNotificationsAndReports(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		} else if channelID == "channel_2" {
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
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})

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

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_NoStandupChannels(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) { return map[string]string{}, nil })

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_GetStandupChannels_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return nil, errors.New("")
	})

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
}

func TestSendNotificationsAndReports_SendStandupReport_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return errors.New("")
	})

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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		} else if channelID == "channel_2" {
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
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})

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

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_GetNotificationStatus_NoData(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		} else if channelID == "channel_2" {
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
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})

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

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_GetUser_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(nil, &model.AppError{})

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-55 * time.Minute)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

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

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return nil, nil
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

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_GetStandupConfig_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(nil, &model.AppError{})

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})

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

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return nil, nil
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

	monkey.Patch(config.Mattermost.GetUser, func(string) (*model.User, *model.AppError) {
		return nil, model.NewAppError("", "", nil, "", http.StatusInternalServerError)
	})

	assert.NotNil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_GetStandupConfig_Nil(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})

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

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1" || userID == "user_id_2" {
				return nil, nil
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

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced as no standup config found is handled")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_WindowOpenNotificationSent_Sent(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		if channelID == "channel_1" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		} else if channelID == "channel_2" {
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
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})

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

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	//mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_NotWorkDay(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
	}
	
	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})
	
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	config.SetConfig(mockConfig)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_Integration(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)
	mockAPI.On("KVGet", "uScyewRiWEwQavauYw9iOK76jISl+5Qq0mV+Cn/jFPs=").Return(
		[]byte("{\"channel_1\": \":channel_1\"}"), nil,
	)
	mockAPI.On("KVGet", util.GetKeyHash(fmt.Sprintf("%s_%s_%s", config.CacheKeyPrefixNotificationStatus, "channel_1", util.GetCurrentDateString("Asia/Kolkata")))).Return(nil, nil)

	windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
	windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}
	standupConfig, _ := json.Marshal(&standup.StandupConfig{
		ChannelId:       "channel_1",
		WindowOpenTime:  windowOpenTime,
		WindowCloseTime: windowCloseTime,
		Enabled:         true,
		Members:         []string{"user_id_1", "user_id_2"},
		ReportFormat:    config.ReportFormatUserAggregated,
		Sections:        []string{"section 1", "section 2"},
		Timezone:		 "Asia/Kolkata",
	})
	mockAPI.On("KVGet", "UzFgbepiypG8qfVARBfHu154LDNiZOw7Mr6Ue4kNZrk=").Return(standupConfig, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_sendWindowCloseNotification_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(nil, util.EmptyAppError())
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
			StandupReportSent:           false,
			WindowOpenNotificationSent:  false,
			WindowCloseNotificationSent: false,
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

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{}, nil
	})

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_FilterChannelNotifications_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(nil, util.EmptyAppError())
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
			StandupReportSent:           false,
			WindowOpenNotificationSent:  false,
			WindowCloseNotificationSent: false,
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

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return nil, errors.New("")
	})

	assert.NotNil(t, SendNotificationsAndReports())
}

func TestSendNotificationsAndReports_Standup_Disabled(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
			StandupReportSent:           false,
			WindowOpenNotificationSent:  false,
			WindowCloseNotificationSent: false,
		}, nil
	})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

		return &standup.StandupConfig{
			ChannelId:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         false,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
}

func TestSendNotificationsAndReports_StandupReport_Sent(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
			StandupReportSent:           true,
			WindowOpenNotificationSent:  true,
			WindowCloseNotificationSent: true,
		}, nil
	})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}
	
		return &standup.StandupConfig{
			ChannelId:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
}

func TestSendNotificationsAndReports_SendWindowOpenNotification_CreatePost_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, model.NewAppError("", "", nil, "", 0))
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  false,
				WindowCloseNotificationSent: false,
			}, nil
	})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{}, nil
	})

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestSendNotificationsAndReports_ShouldSendWindowOpenNotification_NotYet(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  false,
				WindowCloseNotificationSent: false,
			}, nil
	})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(2 * time.Hour)}

		return &standup.StandupConfig{
			ChannelId:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 0)
}

func TestSendNotificationsAndReports_WindowCloseNotification(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
		}, nil
	})

	monkey.Patch(SendStandupReport, func(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
		return nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{
			StandupReportSent:           false,
			WindowOpenNotificationSent:  false,
			WindowCloseNotificationSent: false,
		}, nil
	})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-55 * time.Minute)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

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

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return nil, nil
	})

	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestGetNotificationStatus(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	notificationStatusJSON, _ := json.Marshal(ChannelNotificationStatus{
		WindowOpenNotificationSent:  true,
		WindowCloseNotificationSent: false,
		StandupReportSent:           true,
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(notificationStatusJSON, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	actualNotificationStatus, err := GetNotificationStatus("channel_1")
	assert.Nil(t, err, "no error should have been produced")

	expectedNotificationStatus := &ChannelNotificationStatus{
		WindowOpenNotificationSent:  true,
		WindowCloseNotificationSent: false,
		StandupReportSent:           true,
	}
	assert.Equal(t, expectedNotificationStatus, actualNotificationStatus)
}

func TestGetNotificationStatus_KVGet_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, model.NewAppError("", "", nil, "", 0))

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	actualNotificationStatus, err := GetNotificationStatus("channel_1")
	assert.NotNil(t, err, "error should have been produced as KVGet failed")
	assert.Nil(t, actualNotificationStatus)
}

func TestGetNotificationStatus_Json_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	notificationStatusJSON, _ := json.Marshal(ChannelNotificationStatus{
		WindowOpenNotificationSent:  true,
		WindowCloseNotificationSent: false,
		StandupReportSent:           true,
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(notificationStatusJSON[0:len(notificationStatusJSON)-10], nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	actualNotificationStatus, err := GetNotificationStatus("channel_1")
	assert.NotNil(t, err, "error should have been produced as inbalid JSOn was returned by KVGet")
	assert.Nil(t, actualNotificationStatus)
}

func TestGetNotificationStatus_KVSet_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	notificationStatusJSON, _ := json.Marshal(ChannelNotificationStatus{
		WindowOpenNotificationSent:  true,
		WindowCloseNotificationSent: false,
		StandupReportSent:           true,
	})
	mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(notificationStatusJSON, nil)
	mockAPI.On("KVSet")

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	actualNotificationStatus, err := GetNotificationStatus("channel_1")
	assert.Nil(t, err, "no error should have been produced")

	expectedNotificationStatus := &ChannelNotificationStatus{
		WindowOpenNotificationSent:  true,
		WindowCloseNotificationSent: false,
		StandupReportSent:           true,
	}
	assert.Equal(t, expectedNotificationStatus, actualNotificationStatus)
}

func TestSendStandupReport(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "shouldn't produce error as standup with no members is a valid case")
}

func TestSendStandupReport_GetUserStandup_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return nil, errors.New("")
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetUserStandup failed")
}

func TestSendStandupReport_GetUserStandup_Nil(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("LogInfo", mock.Anything, mock.Anything, mock.Anything).Return()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return nil, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "shouldn't produce error as standup with no members is a valid case")
}

func TestSendStandupReport_GetUserStandup_Nil_GetUser_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("LogInfo", mock.Anything, mock.Anything, mock.Anything).Return()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		nil, model.NewAppError("", "", nil, "", 0),
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return nil, nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetUser failed")
}

func TestSendStandupReport_ReportFormatUserAggregated(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "shouldn't produce error as standup with no members is a valid case")
}

func TestSendStandupReport_UnknownReportFormat(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    "some_unknown_report_format",
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.NotNil(t, err, "should produce error as report format was unknown")
}

func TestSendStandupReport_ReportVisibility_Public(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("CreatePost", mock.AnythingOfType("*model.Post"), mock.Anything).Return(&model.Post{}, nil)

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.Nil(t, err, "shouldn't produce error as standup with no members is a valid case")
}

func TestSendStandupReport_ReportVisibility_Public_CreatePost_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("CreatePost", mock.AnythingOfType("*model.Post"), mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 0))

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.NotNil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", false)
	assert.NotNil(t, err, "shouldn't produce error as standup with no members is a valid case")
}

func TestSendStandupReport_UpdateStatus_True(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return &ChannelNotificationStatus{}, nil
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", true)
	assert.Nil(t, err, "should not produce any error")

	// no standup channels specified
	err = SendStandupReport([]string{}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", false)
	assert.Nil(t, err, "should not produce any error")

	// error in GetStandupConfig
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, errors.New("")
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", true)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig failed")

	// no standup config
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		return nil, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", true)
	assert.NotNil(t, err, "should produce any error as GetStandupConfig didn't return any standup config")

	// standup with no members
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", true)
	assert.Nil(t, err, "shouldn't produce error as standup with no members is a valid case")
}

func TestSendStandupReport_UpdateStatus_True_GetNotificationStatus_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()

	mockAPI.On("GetUser", "user_id_1").Return(
		&model.User{
			FirstName: "Foo",
			LastName:  "Bar",
		}, nil,
	)

	mockAPI.On("GetUser", "user_id_2").Return(
		&model.User{
			FirstName: "John",
			LastName:  "Doe",
		}, nil,
	)

	mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(&model.Post{})

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-5 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       channelID,
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			ReportFormat:    config.ReportFormatTypeAggregated,
			Sections:        []string{"section_1", "section_2"},
			Members:         []string{"user_id_1", "user_id_2"},
			Enabled:         true,
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return &standup.UserStandup{
			UserID:    userID,
			ChannelID: channelID,
			Standup: map[string]*[]string{
				"section_1": {"task_1", "task_2"},
				"section_2": {"task_3", "task_4"},
			},
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		return nil, errors.New("")
	})

	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return nil
	})

	err := SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", true)
	assert.Nil(t, err, "should not produce any error")

	monkey.Unpatch(GetNotificationStatus)
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
	monkey.Patch(SetNotificationStatus, func(channelID string, status *ChannelNotificationStatus) error {
		return errors.New("")
	})

	err = SendStandupReport([]string{"channel_1", "channel_2"}, otime.Now("Asia/Kolkata"), ReportVisibilityPrivate, "user_1", true)
	assert.NotNil(t, err, "should not produce any error")
}

func TestSetNotificationStatus(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_id_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	assert.Nil(t, SetNotificationStatus("channel_id_1", &ChannelNotificationStatus{}))
}

func TestSetNotificationStatus_JsonMarshal_Error(t *testing.T) {
	defer TearDown()
	baseMock()

	monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, errors.New("")
	})
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_id_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})

	assert.NotNil(t, SetNotificationStatus("channel_id_1", &ChannelNotificationStatus{}))
}

func TestSetNotificationStatus_KVSet_Error(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(model.NewAppError("", "", nil, "", 0))
	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
		windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

		return &standup.StandupConfig{
			ChannelId:       "channel_id_1",
			WindowOpenTime:  windowOpenTime,
			WindowCloseTime: windowCloseTime,
			Enabled:         true,
			Members:         []string{"user_id_1", "user_id_2"},
			ReportFormat:    config.ReportFormatUserAggregated,
			Sections:        []string{"section 1", "section 2"},
			Timezone:		 "Asia/Kolkata",
		}, nil
	})
	assert.NotNil(t, SetNotificationStatus("channel_id_1", &ChannelNotificationStatus{}))
}

func TestSendNotificationsAndReports_GetUserStandup_Nodata(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

	monkey.Patch(GetNotificationStatus, func(channelID string) (*ChannelNotificationStatus, error) {
		if channelID == "channel_1" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: true,
			}, nil
		} else if channelID == "channel_2" {
			return &ChannelNotificationStatus{
				StandupReportSent:           false,
				WindowOpenNotificationSent:  true,
				WindowCloseNotificationSent: true,
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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		} else if channelID == "channel_2" {
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
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})

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
	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		return nil, nil
	})
	err :=SendStandupReport([]string{"channel_1", "channel_2", "channel_3"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", true)
	assert.Nil(t, err, "should not produce any error")
	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	mockAPI.AssertNumberOfCalls(t, "CreatePost", 3)
}

func TestSendNotificationsAndReports_MemberNoStandup(t *testing.T) {
	defer TearDown()
	mockAPI := baseMock()
	mockAPI.On("CreatePost", mock.AnythingOfType(model.Post{}.Type)).Return(&model.Post{}, nil)
	mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Username: "username"}, nil)

	location, _ := time.LoadLocation("Asia/Kolkata")
	mockConfig := &config.Configuration{
		Location:      location,
		WorkWeekStart: strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) - 1),
		WorkWeekEnd:   strconv.Itoa(int(otime.Now("Asia/Kolkata").Time.Weekday()) + 1),
	}

	config.SetConfig(mockConfig)

	monkey.Patch(standup.GetStandupChannels, func() (map[string]string, error) {
		return map[string]string{
			"channel_1": "channel_1",
			"channel_2": "channel_2",
			"channel_3": "channel_3",
		}, nil
	})

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

	monkey.Patch(standup.GetStandupConfig, func(channelID string) (*standup.StandupConfig, error) {
		if channelID == "channel_1" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_1",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatTypeAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		} else if channelID == "channel_2" {
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
		} else if channelID == "channel_3" {
			windowOpenTime := otime.OTime{otime.Now("Asia/Kolkata").Add(-1 * time.Hour)}
			windowCloseTime := otime.OTime{otime.Now("Asia/Kolkata").Add(1 * time.Minute)}

			return &standup.StandupConfig{
				ChannelId:       "channel_3",
				WindowOpenTime:  windowOpenTime,
				WindowCloseTime: windowCloseTime,
				Enabled:         true,
				Members:         []string{"user_id_1", "user_id_2"},
				ReportFormat:    config.ReportFormatUserAggregated,
				Sections:        []string{"section 1", "section 2"},
				Timezone:		 "Asia/Kolkata",
			}, nil
		}

		t.Fatal("unknown argument encountered: " + channelID)
		return nil, nil
	})

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

	monkey.Patch(standup.GetUserStandup, func(userID, channelID string, date otime.OTime) (*standup.UserStandup, error) {
		if channelID == "channel_1" {
			if userID == "user_id_1"  {
				return nil, nil
			} else if userID == "user_id_2" {
				return &standup.UserStandup{}, nil
			}
		} else if channelID == "channel_2" {
			if userID == "user_id_1"  {
				return nil, nil
			} else if userID == "user_id_2" {
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
	err :=SendStandupReport([]string{"channel_1", "channel_2", "channel_3"}, otime.Now("Asia/Kolkata"), ReportVisibilityPublic, "user_1", true)
	assert.Nil(t, err, "should not produce any error")
	assert.Nil(t, SendNotificationsAndReports(), "no error should have been produced")
	
}
