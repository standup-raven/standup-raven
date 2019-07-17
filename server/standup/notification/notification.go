package notification

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/otime"
	"github.com/standup-raven/standup-raven/server/standup"
	"github.com/standup-raven/standup-raven/server/util"
)

type ChannelNotificationStatus struct {
	WindowOpenNotificationSent  bool `json:"windowOpenNotificationSent"`
	WindowCloseNotificationSent bool `json:"windowCloseNotificationSent"`
	StandupReportSent           bool `json:"standupReportSent"`
}

const (
	// statuses for standup notification in a channel
	ChannelNotificationStatusSent   = "sent"
	ChannelNotificationStatusNotYet = "not_yet"
	ChannelNotificationStatusSend   = "send"

	// standup report visibilities
	ReportVisibilityPublic  = "public"
	ReportVisibilityPrivate = "private"
)

// SendNotificationsAndReports checks for all standup channels and sends
// notifications and standup reports as needed.
// This is the entry point of the whole standup cycle.
func SendNotificationsAndReports() error {
	channelIDs, err := standup.GetStandupChannels()
	if err != nil {
		return err
	}

	channels, err := channelsWorkDay(channelIDs)
	if err != nil {
		return err
	}

	a, b, c, err := filterChannelNotification(channels)
	if err != nil {
		return err
	}

	sendWindowOpenNotification(a)
	if err := sendWindowCloseNotification(b); err != nil {
		return err
	}
	if err := sendAllStandupReport(c); err != nil {
		return err
	}

	return nil
}

func sendAllStandupReport(channelIDs []string) error {
	for _, channelID := range channelIDs {
		standupConfig, err := standup.GetStandupConfig(channelID)
		if err != nil {
			return err
		}
		if standupConfig == nil {
			return errors.New("standup not configured for channel: " + channelID)
		}
		standupReportError := SendStandupReport([]string{channelID}, otime.Now(standupConfig.Timezone), ReportVisibilityPublic, "", true)
		if standupReportError != nil {
			return standupReportError
		}
	}
	return nil
}

//channelsWorkDay return channels that have working day today
func channelsWorkDay(channels map[string]string) (map[string]string, error) {
	channelIDs := map[string]string{}
	for channelID := range channels {
		standupConfig, err := standup.GetStandupConfig(channelID)
		if err != nil {
			return nil, err
		}
		if standupConfig == nil {
			continue
		}

		// don't send notifications if it's not a work week.
		if isWorkDay(standupConfig.Timezone) {
			channelIDs[channelID] = channelID
		}
	}
	return channelIDs, nil
}

// GetNotificationStatus gets the notification status for specified channel
func GetNotificationStatus(channelID string) (*ChannelNotificationStatus, error) {
	logger.Debug(fmt.Sprintf("Fetching notification status for channel: %s", channelID), nil)
	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		return nil, err
	}
	if standupConfig == nil {
		return nil, errors.New("standup not configured for channel: " + channelID)
	}
	key := fmt.Sprintf("%s_%s_%s", config.CacheKeyPrefixNotificationStatus, channelID, util.GetCurrentDateString(standupConfig.Timezone))
	data, appErr := config.Mattermost.KVGet(util.GetKeyHash(key))
	if appErr != nil {
		logger.Error("Couldn't get notification status from KV store", appErr, nil)
		return nil, errors.New(appErr.Error())
	} else if len(data) == 0 {
		return &ChannelNotificationStatus{}, nil
	}

	status := &ChannelNotificationStatus{}
	if err := json.Unmarshal(data, status); err != nil {
		logger.Error("Couldn't unmarshal notification status data into struct", err, map[string]interface{}{"channelID": channelID, "data": string(data)})
		return nil, err
	}

	logger.Debug(fmt.Sprintf("notification status for channel: %s, %v", channelID, status), nil)
	return status, nil
}

// SendStandupReport sends standup report for all channel IDs specified
func SendStandupReport(channelIDs []string, date otime.OTime, visibility string, userId string, updateStatus bool) error {
	for _, channelID := range channelIDs {
		logger.Info("Sending standup report for channel: "+channelID+" time: "+date.GetDateString(), nil)

		standupConfig, err := standup.GetStandupConfig(channelID)
		if err != nil {
			return err
		}

		if standupConfig == nil {
			return errors.New("standup not configured for channel: " + channelID)
		}

		// standup of all channel standup members
		var members []*standup.UserStandup

		// names of channel standup members who haven't yet submitted their standup
		membersNoStandup := []string{}
		for _, userID := range standupConfig.Members {
			userStandup, err := standup.GetUserStandup(userID, channelID, date)
			if err != nil {
				return err
			} else if userStandup == nil {
				// if user has not submitted standup
				logger.Info("Could not fetch standup for user: "+userID, nil)

				user, appErr := config.Mattermost.GetUser(userID)
				if appErr != nil {
					logger.Error("Couldn't fetch user", appErr, map[string]interface{}{"userID": userId})
					return errors.New(appErr.Error())
				}

				membersNoStandup = append(membersNoStandup, user.Username)

				continue
			}

			members = append(members, userStandup)
		}

		var post *model.Post

		if standupConfig.ReportFormat == config.ReportFormatTypeAggregated {
			post, err = generateTypeAggregatedStandupReport(standupConfig, members, membersNoStandup, channelID, date)
		} else if standupConfig.ReportFormat == config.ReportFormatUserAggregated {
			post, err = generateUserAggregatedStandupReport(standupConfig, members, membersNoStandup, channelID, date)
		} else {
			err = errors.New("Unknown report format encountered for channel: " + channelID + ", report format: " + standupConfig.ReportFormat)
			logger.Error("Unknown report format encountered for channel", err, nil)
		}

		if err != nil {
			return err
		}

		if visibility == ReportVisibilityPrivate {
			config.Mattermost.SendEphemeralPost(userId, post)
		} else {
			_, appErr := config.Mattermost.CreatePost(post)
			if appErr != nil {
				logger.Error("Couldn't create standup report post", appErr, nil)
				return errors.New(appErr.Error())
			}
		}

		if updateStatus {
			notificationStatus, err := GetNotificationStatus(channelID)
			if err != nil {
				continue
			}

			notificationStatus.StandupReportSent = true
			if err := SetNotificationStatus(channelID, notificationStatus); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetNotificationStatus sets provided notification status for the specified channel ID.
func SetNotificationStatus(channelID string, status *ChannelNotificationStatus) error {
	standupConfig, err := standup.GetStandupConfig(channelID)
	if err != nil {
		return err
	}
	if standupConfig == nil {
		return errors.New("standup not configured for channel: " + channelID)
	}
	key := fmt.Sprintf("%s_%s_%s", config.CacheKeyPrefixNotificationStatus, channelID, util.GetCurrentDateString(standupConfig.Timezone))
	serializedStatus, err := json.Marshal(status)
	if err != nil {
		logger.Error("Couldn't marshal standup status data", err, nil)
		return err
	}

	if appErr := config.Mattermost.KVSet(util.GetKeyHash(key), serializedStatus); appErr != nil {
		logger.Error("Couldn't save standup status data into KV store", appErr, nil)
		return errors.New(appErr.Error())
	}

	return nil
}

// filterChannelNotification filters all provided standup channels into three categories -
// 		1. channels requiring window open notification
// 		2. channels requiring window close notification
//		3. channels requiring standup report
func filterChannelNotification(channelIDs map[string]string) ([]string, []string, []string, error) {
	logger.Debug("Filtering channels for sending notifications", nil)

	var windowOpenNotificationChannels, windowCloseNotificationChannels, standupReportChannels []string

	for channelID := range channelIDs {
		logger.Debug(fmt.Sprintf("Processing channel: %s", channelID), nil)

		notificationStatus, err := GetNotificationStatus(channelID)
		if err != nil {
			return nil, nil, nil, err
		}

		standupConfig, err := standup.GetStandupConfig(channelID)
		if err != nil {
			return nil, nil, nil, err
		}

		if standupConfig == nil {
			logger.Error("Unable to find standup config for channel", nil, map[string]interface{}{"channelID": channelID})
			continue
		}

		if !standupConfig.Enabled {
			continue
		}

		// we check in opposite order of time and check for just one notification to send.
		// This prevents expired notifications from being sent in case some of
		// the notifications were missed in the past
		if status := shouldSendStandupReport(notificationStatus, standupConfig); status == ChannelNotificationStatusSend {
			logger.Debug(fmt.Sprintf("Channel [%s] needs standup report", channelID), nil)
			standupReportChannels = append(standupReportChannels, channelID)
		} else if status == ChannelNotificationStatusSent {
			// pass
		} else if shouldSendWindowCloseNotification(notificationStatus, standupConfig) == ChannelNotificationStatusSend {
			if standupConfig.WindowCloseReminderEnabled {
				logger.Debug(fmt.Sprintf("Channel [%s] needs window close notification", channelID), nil)
				windowCloseNotificationChannels = append(windowCloseNotificationChannels, channelID)
			}
		} else if status == ChannelNotificationStatusSent {
			// pass
		} else if shouldSendWindowOpenNotification(notificationStatus, standupConfig) == ChannelNotificationStatusSend {
			if standupConfig.WindowOpenReminderEnabled {
				logger.Debug(fmt.Sprintf("Channel [%s] needs window open notification", channelID), nil)
				windowOpenNotificationChannels = append(windowOpenNotificationChannels, channelID)
			}
		}
	}

	logger.Debug(fmt.Sprintf(
		"Notifications filtered: open: %d, close: %d, reports: %d",
		len(windowOpenNotificationChannels),
		len(windowCloseNotificationChannels),
		len(standupReportChannels),
	), nil)
	return windowOpenNotificationChannels, windowCloseNotificationChannels, standupReportChannels, nil
}

// shouldSendWindowOpenNotification checks if window open notification should
// be sent to the channel with specified notification status
func shouldSendWindowOpenNotification(notificationStatus *ChannelNotificationStatus, standupConfig *standup.StandupConfig) string {
	if notificationStatus.WindowOpenNotificationSent {
		return ChannelNotificationStatusSent
	} else if otime.Now(standupConfig.Timezone).GetTimeWithSeconds(standupConfig.Timezone).After(standupConfig.WindowOpenTime.GetTimeWithSeconds(standupConfig.Timezone).Time) {
		return ChannelNotificationStatusSend
	} else {
		return ChannelNotificationStatusNotYet
	}
}

// shouldSendWindowCloseNotification checks if window close notification should
// be sent to the channel with specified notification status
func shouldSendWindowCloseNotification(notificationStatus *ChannelNotificationStatus, standupConfig *standup.StandupConfig) string {
	if notificationStatus.WindowCloseNotificationSent {
		return ChannelNotificationStatusSent
	}

	windowDuration := standupConfig.WindowCloseTime.GetTime(standupConfig.Timezone).Time.Sub(standupConfig.WindowOpenTime.GetTime(standupConfig.Timezone).Time)
	targetDurationSeconds := windowDuration.Seconds() * config.WindowCloseNotificationDurationPercentage
	targetDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", targetDurationSeconds))

	// now we just need to check if current time is targetDuration seconds after window open time
	if otime.Now(standupConfig.Timezone).GetTimeWithSeconds(standupConfig.Timezone).After(standupConfig.WindowOpenTime.GetTimeWithSeconds(standupConfig.Timezone).Add(targetDuration)) {
		return ChannelNotificationStatusSend
	} else {
		return ChannelNotificationStatusNotYet
	}
}

// shouldSendStandupReport checks if standup report should
// be sent to the channel with specified notification status
func shouldSendStandupReport(notificationStatus *ChannelNotificationStatus, standupConfig *standup.StandupConfig) string {
	if notificationStatus.StandupReportSent {
		return ChannelNotificationStatusSent
	} else if otime.Now(standupConfig.Timezone).GetTimeWithSeconds(standupConfig.Timezone).After(standupConfig.WindowCloseTime.GetTimeWithSeconds(standupConfig.Timezone).Time) {
		return ChannelNotificationStatusSend
	} else {
		return ChannelNotificationStatusNotYet
	}
}

// sendWindowOpenNotification sends window open notification to the specified channels
func sendWindowOpenNotification(channelIDs []string) {
	for _, channelID := range channelIDs {
		post := &model.Post{
			ChannelId: channelID,
			UserId:    config.GetConfig().BotUserID,
			Type:      model.POST_DEFAULT,
			Message:   "Please start filling your standup!",
		}

		if _, appErr := config.Mattermost.CreatePost(post); appErr != nil {
			logger.Error("Error sending window open notification for channel", appErr, map[string]interface{}{"channelID": channelID})
			continue
		}

		notificationStatus, err := GetNotificationStatus(channelID)
		if err != nil {
			continue
		}

		notificationStatus.WindowOpenNotificationSent = true
		if err := SetNotificationStatus(channelID, notificationStatus); err != nil {
			continue
		}
	}
}

// sendWindowCloseNotification sends window close notification to the specified channels
func sendWindowCloseNotification(channelIDs []string) error {
	for _, channelID := range channelIDs {
		standupConfig, err := standup.GetStandupConfig(channelID)
		if err != nil {
			return err
		}

		if standupConfig == nil {
			logger.Error("Unable to find standup config for channel", nil, map[string]interface{}{"channelID": channelID})
			continue
		}

		logger.Debug("Fetching members with pending standup reports", nil)

		var usersPendingStandup []string
		for _, userId := range standupConfig.Members {
			userStandup, err := standup.GetUserStandup(userId, channelID, otime.Now(standupConfig.Timezone))
			if err != nil {
				return err
			}

			if userStandup == nil {
				usersPendingStandup = append(usersPendingStandup, userId)
			}
		}

		logger.Debug("Fetching usernames for users with pending standup", nil)

		for i := range usersPendingStandup {
			user, err := config.Mattermost.GetUser(usersPendingStandup[i])
			if err != nil {
				logger.Error("Couldn't find user with user ID", err, map[string]interface{}{"userID": usersPendingStandup[i]})
				return err
			}

			usersPendingStandup[i] = user.Username
		}

		// no need to send reminder if everyone has filled their standup
		if len(usersPendingStandup) == 0 {
			logger.Debug("Not sending window close notification. No pending standups found.", nil, nil)
			return nil
		}

		var message string
		if len(usersPendingStandup) > 0 {
			message = fmt.Sprintf("@%s - a gentle reminder to fill your standup.", strings.Join(usersPendingStandup, ", @"))
		} else {
			message = "A gentle reminder to fill your standup."
		}

		post := &model.Post{
			ChannelId: channelID,
			UserId:    config.GetConfig().BotUserID,
			Type:      model.POST_DEFAULT,
			Message:   message,
		}

		if _, appErr := config.Mattermost.CreatePost(post); appErr != nil {
			logger.Error("Error sending window open notification for channel", appErr, map[string]interface{}{"channelID": channelID})
			continue
		}

		notificationStatus, err := GetNotificationStatus(channelID)
		if err != nil {
			continue
		}

		notificationStatus.WindowCloseNotificationSent = true
		if err := SetNotificationStatus(channelID, notificationStatus); err != nil {
			return err
		}
	}

	return nil
}

// generateTypeAggregatedStandupReport generates a Type Aggregated standup report
func generateTypeAggregatedStandupReport(
	standupConfig *standup.StandupConfig,
	userStandups []*standup.UserStandup,
	membersNoStandup []string,
	channelID string,
	date otime.OTime,
) (*model.Post, error) {
	logger.Debug("Generating type aggregated standup report for channel: "+channelID, nil)

	userTasks := map[string]string{}
	userNoTasks := map[string][]string{}

	for _, userStandup := range userStandups {
		for _, sectionTitle := range standupConfig.Sections {
			userDisplayName, err := getUserDisplayName(userStandup.UserID)
			if err != nil {
				logger.Debug("Couldn't fetch display name for user", err, map[string]string{"userID": userStandup.UserID})
				return nil, err
			}

			header := fmt.Sprintf("##### %s %s", util.UserIcon(userStandup.UserID), userDisplayName)

			if userStandup.Standup[sectionTitle] != nil && len(*userStandup.Standup[sectionTitle]) > 0 {
				userTasks[sectionTitle] += fmt.Sprintf("%s\n1. %s\n", header, strings.Join(*userStandup.Standup[sectionTitle], "\n1. "))
			} else {
				userNoTasks[sectionTitle] = append(userNoTasks[sectionTitle], userDisplayName)
			}
		}
	}

	text := fmt.Sprintf("#### Standup Report for *%s*\n\n", date.Format("2 Jan 2006"))

	if len(userStandups) > 0 {
		if len(membersNoStandup) > 0 {
			text += fmt.Sprintf("%s %s not submitted their standup.\n", strings.Join(membersNoStandup, ", "), util.HasHave(len(membersNoStandup)))
		}

		for _, sectionTitle := range standupConfig.Sections {
			text += "##### ** " + sectionTitle + "**\n\n" + userTasks[sectionTitle] + "\n"
			if len(userNoTasks[sectionTitle]) > 0 {
				text += fmt.Sprintf(
					"%s %s no open items for %s\n",
					strings.Join(userNoTasks[sectionTitle], ", "),
					util.HasHave(len(userNoTasks[sectionTitle])),
					sectionTitle,
				)
			}
		}
	} else {
		text += ":warning: **No user has submitted their standup.**"
	}

	return &model.Post{
		ChannelId: channelID,
		UserId:    config.GetConfig().BotUserID,
		Message:   text,
	}, nil
}

// generateUserAggregatedStandupReport generates a User Aggregated standup report
func generateUserAggregatedStandupReport(
	standupConfig *standup.StandupConfig,
	userStandups []*standup.UserStandup,
	membersNoStandup []string,
	channelID string,
	date otime.OTime,
) (*model.Post, error) {
	logger.Debug("Generating user aggregated standup report for channel: "+channelID, nil)

	userTasks := ""

	for _, userStandup := range userStandups {
		userDisplayName, err := getUserDisplayName(userStandup.UserID)
		if err != nil {
			logger.Debug("Couldn't fetch display name for user", err, map[string]string{"userID": userStandup.UserID})
			return nil, err
		}

		header := fmt.Sprintf("#### %s %s", util.UserIcon(userStandup.UserID), userDisplayName)

		userTask := header + "\n\n"

		for _, sectionTitle := range standupConfig.Sections {
			if userStandup.Standup[sectionTitle] == nil || len(*userStandup.Standup[sectionTitle]) == 0 {
				continue
			}

			userTask += fmt.Sprintf("##### %s\n", sectionTitle)
			userTask += "1. " + strings.Join(*userStandup.Standup[sectionTitle], "\n1. ") + "\n\n"
		}

		userTasks += userTask
	}

	text := fmt.Sprintf("#### Standup Report for *%s*\n", date.Format("2 Jan 2006"))

	if len(userStandups) > 0 {
		if len(membersNoStandup) > 0 {
			text += fmt.Sprintf("\n@%s %s not submitted their standup\n\n", strings.Join(membersNoStandup, ", @"), util.HasHave(len(membersNoStandup)))
		}

		text += userTasks
	} else {
		text += ":warning: **No user has submitted their standup.**"
	}

	conf := config.GetConfig()

	return &model.Post{
		ChannelId: channelID,
		UserId:    conf.BotUserID,
		Message:   text,
	}, nil
}

func isWorkDay(timezone string) bool {
	conf := config.GetConfig()
	dayOfWeek := int(otime.Now(timezone).Time.Weekday())

	workWeekStart, _ := strconv.Atoi(conf.WorkWeekStart)
	workWeekEnd, _ := strconv.Atoi(conf.WorkWeekEnd)

	if workWeekStart < workWeekEnd {
		return workWeekStart <= dayOfWeek && dayOfWeek <= workWeekEnd
	} else {
		return dayOfWeek >= workWeekStart || dayOfWeek <= workWeekEnd
	}
}

func getUserDisplayName(userID string) (string, error) {
	user, appErr := config.Mattermost.GetUser(userID)
	if appErr != nil {
		return "", errors.New(appErr.Error())
	}
	return user.GetDisplayName(model.SHOW_FULLNAME), nil
}
