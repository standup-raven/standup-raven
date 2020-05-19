package standup

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/teambition/rrule-go"
	"strconv"
	"strings"
	"time"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/util"
	"github.com/thoas/go-funk"
)

const (
	standupSectionsMinLength       = 1
	channelHeaderScheduleSeparator = "|"
)

type UserStandup struct {
	UserID    string               `json:"userId"`
	ChannelID string               `json:"channelId"`
	Standup   map[string]*[]string `json:"standup"`
}

func (us *UserStandup) IsValid() error {
	if us.UserID == "" {
		return errors.New("No user ID specified in standup")
	}

	if us.ChannelID == "" {
		return errors.New("No channels ID specified in standup")
	}

	if _, err := config.Mattermost.GetChannel(us.ChannelID); err != nil {
		return errors.New("No channel found with channel ID: " + us.ChannelID)
	}

	maxLen := 0
	for _, sectionTasks := range us.Standup {
		maxLen = util.Max(maxLen, len(*sectionTasks))
	}

	if maxLen == 0 {
		return errors.New("No tasks found. Please specify tasks for at least one section")
	}

	return nil
}

type StandupConfig struct {
	ChannelId                  string      `json:"channelId"`
	WindowOpenTime             otime.OTime `json:"windowOpenTime"`
	WindowCloseTime            otime.OTime `json:"windowCloseTime"`
	ReportFormat               string      `json:"reportFormat"`
	Members                    []string    `json:"members"`
	Sections                   []string    `json:"sections"`
	Enabled                    bool        `json:"enabled"`
	Timezone                   string      `json:"timezone"`
	WindowOpenReminderEnabled  bool        `json:"windowOpenReminderEnabled"`
	WindowCloseReminderEnabled bool        `json:"windowCloseReminderEnabled"`
	ScheduleEnabled            bool        `json:"scheduleEnabled"`
	RRule                      rrule.RRule `json:"rrule"`
	RRuleString                string      `json:"rruleString"`
	StartDate                  time.Time   `json:"startDate"`
}

func (sc *StandupConfig) IsValid() error {
	if sc.ChannelId == "" {
		return errors.New("Channel ID cannot be empty")
	}

	emptyTime := otime.OTime{}

	if sc.WindowOpenTime == emptyTime {
		return errors.New("window open time cannot be empty")
	}

	if sc.WindowCloseTime == emptyTime {
		return errors.New("window close time cannot be empty")
	}

	if sc.WindowOpenTime.Time.After(sc.WindowCloseTime.Time) {
		return errors.New("Window open time cannot be after window close time")
	}

	if sc.Timezone == "" {
		return errors.New("Timezone cannot be empty")
	}

	reportFormat := sc.ReportFormat
	if !funk.Contains(config.ReportFormats, reportFormat) {
		return errors.New(fmt.Sprintf(
			"Invalid report format specified. Report format should be one of: \"%s\"",
			strings.Join(config.ReportFormats, "\", \"")),
		)
	}

	if _, err := time.LoadLocation(sc.Timezone); err != nil {
		return errors.New(fmt.Sprintf(
			"Invalid timezone specified : \"%s\"", sc.Timezone),
		)
	}

	if len(sc.Sections) < standupSectionsMinLength {
		return errors.New(fmt.Sprintf("Too few sections in standup. Required at least %d section%s.", standupSectionsMinLength, util.SingularPlural(standupSectionsMinLength)))
	}

	if duplicateSection, hasDuplicate := util.ContainsDuplicates(&sc.Sections); hasDuplicate {
		return errors.New("duplicate sections are not allowed. Contains duplicate section '" + duplicateSection + "'")
	}

	if duplicateMember, hasDuplicate := util.ContainsDuplicates(&sc.Members); hasDuplicate {
		return errors.New("duplicate members are not allowed. Contains duplicate member '" + duplicateMember + "'")
	}

	return nil
}

func (sc *StandupConfig) ToJson() string {
	b, _ := json.Marshal(sc)
	return string(b)
}

func (sc *StandupConfig) PreSave() error {
	// parse location
	location, err := time.LoadLocation(sc.Timezone)
	if err != nil {
		logger.Error("Unable to parse standup location", err, map[string]interface{}{"location": sc.Timezone})
		return err
	}

	// remove time from start date
	sc.StartDate = time.Date(
		sc.StartDate.Year(),
		sc.StartDate.Month(),
		sc.StartDate.Day(),
		0,
		0,
		0,
		0,
		location,
	)

	// parse rrule
	rruleOptions, err := rrule.StrToROption(sc.RRuleString)
	if err != nil {
		logger.Error("unable to parse rrule string ini standup config pre-save", err, map[string]interface{}{
			"rrule":     sc.RRuleString,
			"channelID": sc.ChannelId,
		})
		return err
	}

	rrule, err := rrule.NewRRule(*rruleOptions)
	if err != nil {
		logger.Error("unable to create new rrule from options", err, map[string]interface{}{
			"rrule":     sc.RRuleString,
			"channelID": sc.ChannelId,
		})
		return err
	}

	rrule.DTStart(sc.StartDate)

	return nil
}

// GenerateScheduleString generates a user-friendly, string representation of standup schedule.
func (sc *StandupConfig) GenerateScheduleString() string {
	pluginConfig := config.GetConfig()

	windowOpenTime := sc.WindowOpenTime.Format("15:04")
	windowCloseTime := sc.WindowCloseTime.Format("15:04")

	workWeekStartNumber, _ := strconv.Atoi(pluginConfig.WorkWeekStart)
	workWeekEndNumber, _ := strconv.Atoi(pluginConfig.WorkWeekEnd)
	workWeekStart := time.Weekday(workWeekStartNumber).String()
	workWeekEnd := time.Weekday(workWeekEndNumber).String()

	return fmt.Sprintf("**Standup Schedule**: %s to %s, %s to %s", workWeekStart, workWeekEnd, windowOpenTime, windowCloseTime)
}

// AddStandupChannel adds the specified channel to the list of standup channels.
// This is later user for iterating over all standup channels.
func AddStandupChannel(channelID string) error {
	logger.Debug(fmt.Sprintf("Adding standup channel: %s", channelID), nil)

	channels, err := GetStandupChannels()
	if err != nil {
		return err
	}

	channels[channelID] = channelID
	return setStandupChannels(channels)
}

// GetStandupChannels fetches all channels where standup is configured.
// Returns a map of channel ID to channel ID for maintaining uniqueness.
func GetStandupChannels() (map[string]string, error) {
	logger.Debug("Fetching all standup channels", nil)

	data, appErr := config.Mattermost.KVGet(util.GetKeyHash(config.CacheKeyAllStandupChannels))
	if appErr != nil {
		logger.Error("Couldn't fetch standup channel list from KV store", appErr, nil)
		return nil, errors.New(appErr.Error())
	}

	channels := map[string]string{}

	if len(data) > 0 {
		err := json.Unmarshal(data, &channels)
		if err != nil {
			logger.Error("Couldn't unmarshal standup channel list into map", err, map[string]interface{}{"data": string(data)})
			return nil, err
		}
	}

	logger.Debug(fmt.Sprintf("Found %d standup channels", len(channels)), nil)
	return channels, nil
}

// SaveUserStandup saves a user's standup for a channel
func SaveUserStandup(userStandup *UserStandup) error {
	// span across time zones.
	standupConfig, err := GetStandupConfig(userStandup.ChannelID)
	if err != nil {
		return err
	}
	if standupConfig == nil {
		return errors.New("standup not configured for channel: " + userStandup.ChannelID)
	}
	key := otime.Now(standupConfig.Timezone).GetDateString() + "_" + userStandup.ChannelID + userStandup.UserID
	bytes, err := json.Marshal(userStandup)
	if err != nil {
		logger.Error("Error occurred in serializing user standup", err, nil)
		return err
	}

	if appErr := config.Mattermost.KVSet(util.GetKeyHash(key), bytes); appErr != nil {
		logger.Error("Error occurred in saving user standup in KV store", errors.New(appErr.Error()), nil)
		return appErr
	}

	return nil
}

// GetUserStandup fetches a user's standup for the specified channel and date.
func GetUserStandup(userID, channelID string, date otime.OTime) (*UserStandup, error) {
	key := date.GetDateString() + "_" + channelID + userID
	data, appErr := config.Mattermost.KVGet(util.GetKeyHash(key))
	if appErr != nil {
		logger.Error("Couldn't fetch user standup from KV store", appErr, map[string]interface{}{"userID": userID, "channelID": channelID})
		return nil, errors.New(appErr.Error())
	}

	if len(data) == 0 {
		return nil, nil
	}

	userStandup := &UserStandup{}
	if err := json.Unmarshal(data, userStandup); err != nil {
		logger.Error("Couldn't unmarshal user standup data", err, nil)
		return nil, err
	}

	return userStandup, nil
}

// TODO this should return the set config
// SaveStandupConfig saves standup config for the specified channel
func SaveStandupConfig(standupConfig *StandupConfig) (*StandupConfig, error) {
	logger.Debug(fmt.Sprintf("Saving standup config for channel: %s", standupConfig.ChannelId), nil)

	standupConfig.Members = funk.UniqString(standupConfig.Members)
	serializedStandupConfig, err := json.Marshal(standupConfig)
	if err != nil {
		logger.Error("Couldn't marshal standup config", err, nil)
		return nil, err
	}

	if err := updateChannelHeader(standupConfig); err != nil {
		return nil, err
	}

	key := config.CacheKeyPrefixTeamStandupConfig + standupConfig.ChannelId
	if err := config.Mattermost.KVSet(util.GetKeyHash(key), serializedStandupConfig); err != nil {
		logger.Error("Couldn't save channel standup config in KV store", err, map[string]interface{}{"channelID": standupConfig.ChannelId})
		return nil, err
	}

	return standupConfig, nil
}

func updateChannelHeader(newConfig *StandupConfig) error {
	oldConfig, err := GetStandupConfig(newConfig.ChannelId)
	if err != nil {
		return err
	}

	// no old config is equivalent to having standup schedule disabled in old config
	if oldConfig == nil {
		oldConfig = &StandupConfig{
			ScheduleEnabled: false,
		}
	}

	channel, appErr := config.Mattermost.GetChannel(newConfig.ChannelId)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	if oldConfig.ScheduleEnabled && !newConfig.ScheduleEnabled {
		channel.Header = removeChannelHeaderSchedule(channel.Header)
	} else if !oldConfig.ScheduleEnabled && newConfig.ScheduleEnabled {
		channel.Header = addChannelHeaderSchedule(channel.Header, newConfig.GenerateScheduleString())
	} else if oldConfig.ScheduleEnabled && newConfig.ScheduleEnabled {
		channel.Header = removeChannelHeaderSchedule(channel.Header)
		channel.Header = addChannelHeaderSchedule(channel.Header, newConfig.GenerateScheduleString())
	}

	_, appErr = config.Mattermost.UpdateChannel(channel)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

func removeChannelHeaderSchedule(channelHeader string) string {
	components := strings.Split(channelHeader, channelHeaderScheduleSeparator)
	if len(components) < 2 {
		return channelHeader
	}

	updatedHeader := strings.TrimLeft(components[1], " ")
	if len(components) > 2 {
		updatedHeader = updatedHeader + "|" + strings.Join(components[2:], "|")
	}

	return updatedHeader
}

func addChannelHeaderSchedule(channelHeader string, schedule string) string {
	return schedule + " | " + channelHeader
}

// GetStandupConfig fetches standup config for the specified channel
func GetStandupConfig(channelID string) (*StandupConfig, error) {
	logger.Debug(fmt.Sprintf("Fetching standup config for channel: %s", channelID), nil)

	key := config.CacheKeyPrefixTeamStandupConfig + channelID
	data, appErr := config.Mattermost.KVGet(util.GetKeyHash(key))
	if appErr != nil {
		logger.Error("Couldn't fetch standup config for channel from KV store", appErr, map[string]interface{}{"channelID": channelID})
		return nil, errors.New(appErr.Error())
	}

	if len(data) == 0 {
		logger.Debug(fmt.Sprintf("Counldn't find standup config for channel: %s", channelID), nil)
		return nil, nil
	}

	var standupConfig *StandupConfig
	if len(data) > 0 {
		standupConfig = &StandupConfig{}
		if err := json.Unmarshal(data, standupConfig); err != nil {
			logger.Error("Couldn't unmarshal data into standup config", err, nil)
			return nil, err
		}
	}

	logger.Debug(fmt.Sprintf("Standup config for channel: %s, %v", channelID, standupConfig), nil)
	return standupConfig, nil
}

// setStandupChannels saves the provided list of standup channels in the KV store
func setStandupChannels(channels map[string]string) error {
	logger.Debug("Saving standup channels", nil)

	data, err := json.Marshal(channels)
	if err != nil {
		return err
	}

	appErr := config.Mattermost.KVSet(util.GetKeyHash(config.CacheKeyAllStandupChannels), data)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}
