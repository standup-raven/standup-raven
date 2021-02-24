package standup

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/teambition/rrule-go"

	"github.com/thoas/go-funk"

	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/util"
)

const (
	standupSectionsMinLength       = 1
	channelHeaderScheduleSeparator = "|"
	standupScheduleEndMarker       = "** **"
)

var (
	standupScheduleRegex = regexp.MustCompile(`^\*\*Standup Schedule\*\*: .+\*\* \*\*$`)
	weekRanks            = map[int]string{
		-1: "last",
		1:  "first",
		2:  "second",
		3:  "third",
		4:  "fourth",
	}
)

type UserStandup struct {
	UserID    string               `json:"userId"`
	ChannelID string               `json:"channelId"`
	Standup   map[string]*[]string `json:"standup"`
}

func (us *UserStandup) IsValid() error {
	if us.UserID == "" {
		return errors.New("no user ID specified in standup")
	}

	if us.ChannelID == "" {
		return errors.New("no channels ID specified in standup")
	}

	if _, err := config.Mattermost.GetChannel(us.ChannelID); err != nil {
		return errors.New("No channel found with channel ID: " + us.ChannelID)
	}

	maxLen := 0
	for _, sectionTasks := range us.Standup {
		maxLen = util.Max(maxLen, len(*sectionTasks))
	}

	if maxLen == 0 {
		return errors.New("no tasks found. Please specify tasks for at least one section")
	}

	return nil
}

type Config struct {
	ChannelID                  string       `json:"channelId"`
	WindowOpenTime             otime.OTime  `json:"windowOpenTime"`
	WindowCloseTime            otime.OTime  `json:"windowCloseTime"`
	ReportFormat               string       `json:"reportFormat"`
	Members                    []string     `json:"members"`
	Sections                   []string     `json:"sections"`
	Enabled                    bool         `json:"enabled"`
	Timezone                   string       `json:"timezone"`
	WindowOpenReminderEnabled  bool         `json:"windowOpenReminderEnabled"`
	WindowCloseReminderEnabled bool         `json:"windowCloseReminderEnabled"`
	ScheduleEnabled            bool         `json:"scheduleEnabled"`
	RRule                      *rrule.RRule `json:"rrule"`
	RRuleString                string       `json:"rruleString"`
	StartDate                  time.Time    `json:"startDate"`
}

func (sc *Config) IsValid() error {
	if sc.ChannelID == "" {
		return errors.New("channel ID cannot be empty")
	}

	emptyTime := otime.OTime{}

	if sc.WindowOpenTime == emptyTime {
		return errors.New("window open time cannot be empty")
	}

	if sc.WindowCloseTime == emptyTime {
		return errors.New("window close time cannot be empty")
	}

	if sc.WindowOpenTime.Time.After(sc.WindowCloseTime.Time) {
		return errors.New("window open time cannot be after window close time")
	}

	if sc.Timezone == "" {
		return errors.New("timezone cannot be empty")
	}

	reportFormat := sc.ReportFormat
	if !funk.Contains(config.ReportFormats, reportFormat) {
		return fmt.Errorf("invalid report format specified. Report format should be one of: \"%s\"", strings.Join(config.ReportFormats, "\", \""))
	}

	if _, err := time.LoadLocation(sc.Timezone); err != nil {
		return fmt.Errorf("invalid timezone specified : \"%s\"", sc.Timezone)
	}

	if len(sc.Sections) < standupSectionsMinLength {
		return fmt.Errorf("too few sections in standup. Required at least %d section%s", standupSectionsMinLength, util.SingularPlural(standupSectionsMinLength))
	}

	if duplicateSection, hasDuplicate := util.ContainsDuplicates(&sc.Sections); hasDuplicate {
		return errors.New("Duplicate sections are not allowed. Contains duplicate section '" + duplicateSection + "'")
	}

	if duplicateMember, hasDuplicate := util.ContainsDuplicates(&sc.Members); hasDuplicate {
		return errors.New("Duplicate members are not allowed. Contains duplicate member '" + duplicateMember + "'")
	}

	if sc.RRule.Freq == rrule.WEEKLY && (sc.RRule.OrigOptions.Byweekday == nil || len(sc.RRule.OrigOptions.Byweekday) == 0) {
		return errors.New("at least one day must be selected for weekly standup")
	}

	return nil
}

func (sc *Config) ToJSON() string {
	b, _ := json.Marshal(sc)
	return string(b)
}

func (sc *Config) PreSave() error {
	if err := sc.setStartDateLocation(); err != nil {
		return err
	}

	if err := sc.initializeRRule(); err != nil {
		return err
	}

	sc.fixRRuleTimezone()
	return nil
}

// setStartDateLocation sets timezone of start date to
// be same as standup timezone.
func (sc *Config) setStartDateLocation() error {
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

	return nil
}

// initializeRRule initialized RRULE by parsing the RRULE string.
func (sc *Config) initializeRRule() error {
	rule, err := util.ParseRRuleFromString(sc.RRuleString, sc.StartDate)
	if err != nil {
		logger.Error("unable to parse rrule string in standup config pre-save", err, map[string]interface{}{
			"rrule":     sc.RRuleString,
			"channelID": sc.ChannelID,
		})
		return err
	}

	sc.RRule = rule
	return nil
}

// fixRRuleTimezone fix issue in RRULE caused by countries having
// different timezones in different points in history, specifically
// in the year 0001.
// Timezones are date dependent, for example India had the timezone (at least
// in Go timezone database) +553 in the year 0001 as opposed +530 being followed now.
// RRULE Timesets are internally used by RRULE library for extracting only the time component from date,
// but since date and time are tied up, having the incorrect date causes incorrect time due to timezone
// dependency on date.
//
// So here we bring all timesets to current date (no alterting of time) to
// get the current timezone picked up.
func (sc *Config) fixRRuleTimezone() {
	today := time.Now()
	for i := range sc.RRule.Timeset {
		sc.RRule.Timeset[i] = sc.RRule.Timeset[i].AddDate(today.Year()-1, int(today.Month())-1, today.Day()-1)
	}
}

// GenerateScheduleString generates a user-friendly, string representation of standup schedule.
func (sc *Config) GenerateScheduleString() string {
	_ = config.GetConfig()

	windowOpenTime := sc.WindowOpenTime.Format("15:04")
	windowCloseTime := sc.WindowCloseTime.Format("15:04")

	var frequencyString string

	switch sc.RRule.Freq {
	case rrule.WEEKLY:
		frequencyString = sc.generateWeeklySchedule()
	case rrule.MONTHLY:
		frequencyString = sc.generateMonthlySchedule()
	}

	return fmt.Sprintf("**Standup Schedule**: %s %s to %s", frequencyString, windowOpenTime, windowCloseTime)
}

func (sc *Config) generateWeeklySchedule() string {
	prefix := ""

	if sc.RRule.Interval == 1 {
		prefix = "Weekly"
	} else {
		prefix = fmt.Sprintf("Every %d weeks", sc.RRule.Interval)
	}

	daysOfWeek := make([]string, len(sc.RRule.Byweekday))
	for i, day := range sc.RRule.Byweekday {
		daysOfWeek[i] = strings.ToUpper(time.Weekday((day + 1) % 7).String()[:2])
	}

	return fmt.Sprintf("%s on %s", prefix, strings.Join(daysOfWeek, ", "))
}

func (sc *Config) generateMonthlySchedule() string {
	var prefix, suffix string

	if sc.RRule.Interval == 1 {
		prefix = "Monthly"
	} else {
		prefix = fmt.Sprintf("Every %d months", sc.RRule.Interval)
	}

	// this indicates "on date" mode,
	// i.e. event occurs on specific day of month
	if len(sc.RRule.Bymonthday) > 0 {
		suffix = humanize.Ordinal(sc.RRule.Bymonthday[0])
	} else if len(sc.RRule.Bysetpos) > 0 {
		weekOrdinal := weekRanks[sc.RRule.Bysetpos[0]]

		var dayOfWeek string
		switch len(sc.RRule.Byweekday) {
		case 1:
			// single day
			dayOfWeek = time.Weekday(sc.RRule.Byweekday[0] + 1).String()
		case 2:
			// weekend
			dayOfWeek = "weekend"
		case 5:
			// weekday
			dayOfWeek = "weekday"
		case 7:
			// any day
			dayOfWeek = "day"
		}

		suffix = weekOrdinal + " " + dayOfWeek
	}

	return fmt.Sprintf("%s on the %s", prefix, suffix)
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

// RemoveStandupChannels removes all specified channels from list of standup channels.
// This is later user for iterating over all standup channels.
func RemoveStandupChannels(channelIDs []string) error {
	logger.Debug(fmt.Sprintf("Removing standup channels: %v", channelIDs), nil)

	channels, err := GetStandupChannels()
	if err != nil {
		return err
	}

	for _, channelID := range channelIDs {
		delete(channels, channelID)
	}

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
func SaveStandupConfig(standupConfig *Config) (*Config, error) {
	logger.Debug(fmt.Sprintf("Saving standup config for channel: %s", standupConfig.ChannelID), nil)

	standupConfig.Members = funk.UniqString(standupConfig.Members)
	serializedStandupConfig, err := json.Marshal(standupConfig)
	if err != nil {
		logger.Error("Couldn't marshal standup config", err, nil)
		return nil, err
	}

	if err := updateChannelHeader(standupConfig); err != nil {
		return nil, err
	}

	key := config.CacheKeyPrefixTeamStandupConfig + standupConfig.ChannelID
	if err := config.Mattermost.KVSet(util.GetKeyHash(key), serializedStandupConfig); err != nil {
		logger.Error("Couldn't save channel standup config in KV store", err, map[string]interface{}{"channelID": standupConfig.ChannelID})
		return nil, err
	}

	return standupConfig, nil
}

func updateChannelHeader(newConfig *Config) error {
	oldConfig, err := GetStandupConfig(newConfig.ChannelID)
	if err != nil {
		return err
	}

	// no old config is equivalent to having standup schedule disabled in old config
	if oldConfig == nil {
		oldConfig = &Config{
			ScheduleEnabled: false,
		}
	}

	channel, appErr := config.Mattermost.GetChannel(newConfig.ChannelID)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	// Updating an archived channel causes error.
	// Skip if channel is archived.
	if channel.DeleteAt != 0 {
		return nil
	}

	switch {
	case oldConfig.ScheduleEnabled && !newConfig.ScheduleEnabled:
		channel.Header = removeChannelHeaderSchedule(channel.Header)
	case !oldConfig.ScheduleEnabled && newConfig.ScheduleEnabled:
		channel.Header = addChannelHeaderSchedule(channel.Header, newConfig.GenerateScheduleString())
	case oldConfig.ScheduleEnabled && newConfig.ScheduleEnabled:
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
	var userDefinedHeader string

	components := strings.Split(channelHeader, channelHeaderScheduleSeparator)
	if standupScheduleRegex.MatchString(strings.TrimSpace(components[0])) {
		userDefinedHeader = strings.Join(components[1:], channelHeaderScheduleSeparator)
	} else {
		userDefinedHeader = channelHeader
	}

	return userDefinedHeader
}

func addChannelHeaderSchedule(channelHeader string, schedule string) string {
	if channelHeader == "" {
		return schedule + standupScheduleEndMarker
	}

	return schedule + standupScheduleEndMarker + " | " + channelHeader
}

// GetStandupConfig fetches standup config for the specified channel
func GetStandupConfig(channelID string) (*Config, error) {
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

	var standupConfig *Config
	if len(data) > 0 {
		standupConfig = &Config{}
		if err := json.Unmarshal(data, standupConfig); err != nil {
			logger.Error("Couldn't unmarshal data into standup config", err, nil)
			return nil, err
		}
	}

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
