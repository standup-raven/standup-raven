package config

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/plugin"
	"go.uber.org/atomic"

	"github.com/standup-raven/standup-raven/server/otime"
)

const (
	PluginName                   = "standup-raven"
	CommandPrefix                = "standup"
	ServerExeToStaticDirRootPath = "/../webapp/static"

	URLPluginBase = "/plugins/" + PluginName
	URLStaticBase = URLPluginBase

	HeaderMattermostUserID = "Mattermost-User-Id"

	ReportFormatUserAggregated = "user_aggregated"
	ReportFormatTypeAggregated = "type_aggregated"

	CacheKeyPrefixNotificationStatus = "notif_status"
	CacheKeyPrefixTeamStandupConfig  = "standup_config_"

	CacheKeyAllStandupChannels = "all_standup_channels"

	WindowCloseNotificationDurationPercentage = 0.8 // 80%

	UserIconURL  = "/api/v4/users/%s/image"
	UserIconSize = "=20x20"

	// Ensure two full cycles can run in a under a minute
	// to handle the special case of 23:59 window close time.
	// If first cycle starts at 23:58:59, second at 23:59:xx1,
	// third will probably run at 00:00:xx2 causing no automated standup reports as
	// the date changed between 23:59 and 00:00:xx2.
	RunnerInterval = 25 * time.Second

	BotUsername     = "raven"
	BotDisplayName  = "Raven"
	OverrideIconURL = URLStaticBase + "/logo.png"
)

var (
	config        atomic.Value
	Mattermost    plugin.API
	ReportFormats = []string{ReportFormatUserAggregated, ReportFormatTypeAggregated}
)

type Configuration struct {
	Location                *time.Location `json:"location"`        // derived attributes
	BotUserID               string         `json:"botUserId"`       // derived attributes
	SentryServerDSN         string         `json:"sentryServerDSN"` // derived attributes
	SentryWebappDSN         string         `json:"sentryWebappDSN"` // derived attributes
	TimeZone                string         `json:"timeZone"`
	PluginVersion           string         `json:"plugin_version"`
	PermissionSchemaEnabled bool           `json:"permissionSchemaEnabled"`
	EnableErrorReporting    bool           `json:"enableErrorReporting"`
}

func GetConfig() *Configuration {
	return config.Load().(*Configuration)
}

func SetConfig(c *Configuration) {
	config.Store(c)
}

func (c *Configuration) ProcessConfiguration() error {
	location, err := time.LoadLocation(c.TimeZone)
	if err != nil {
		Mattermost.LogError("Couldn't load location in time " + err.Error())
		return err
	}

	c.SentryServerDSN = strings.TrimSpace(c.SentryServerDSN)

	if c.EnableErrorReporting && len(c.SentryServerDSN) == 0 {
		Mattermost.LogError("Sentry Server DSN cannot be empty if error reporting is enabled")
		return errors.New("sentry server DSN cannot be empty if error reporting is enabled")
	}

	c.SentryWebappDSN = strings.TrimSpace(c.SentryWebappDSN)

	if c.EnableErrorReporting && len(c.SentryWebappDSN) == 0 {
		Mattermost.LogError("Sentry Webapp DSN cannot be empty if error reporting is enabled")
		return errors.New("sentry webapp DSN cannot be empty if error reporting is enabled")
	}

	c.Location = location
	otime.DefaultLocation = location
	return nil
}

func (c *Configuration) ToJSON() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *Configuration) Clone() *Configuration {
	var clone Configuration
	_ = json.Unmarshal(c.ToJSON(), &clone)
	return &clone
}

func (c *Configuration) Sanitize() *Configuration {
	clone := c.Clone()
	clone.BotUserID = ""
	clone.Location = nil
	clone.SentryServerDSN = ""
	return clone
}
