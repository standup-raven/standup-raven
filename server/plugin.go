package main

import (
	"github.com/getsentry/raven-go"
	"github.com/standup-raven/standup-raven/server/logger"
	"github.com/standup-raven/standup-raven/server/standup/notification"
	"net/http"
	"strconv"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/standup-raven/standup-raven/server/command"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/controller"
	"github.com/standup-raven/standup-raven/server/util"
	"os"
	"path/filepath"
	"github.com/pkg/profile"
)

var SentryEnabled string

var SentryDSN string

type Plugin struct {
	plugin.MattermostPlugin

	handler http.Handler
	running bool
}

func (p *Plugin) OnActivate() error {
	config.Mattermost = p.API

	if err := p.initSentry(); err != nil {
		config.Mattermost.LogError(err.Error())
	}

	if err := p.setupStaticFileServer(); err != nil {
		return err
	}

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	if err := p.RegisterCommands(); err != nil {
		return err
	}

	p.Run()

	return nil
}

func (p *Plugin) setupStaticFileServer() error {
	exe, err := os.Executable()
	if err != nil {
		logger.Error("Couldn't find plugin executable path", err, nil)
		return err
	}
	p.handler = http.FileServer(http.Dir(filepath.Dir(exe) + config.ServerExeToWebappRootPath))
	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	if config.Mattermost != nil {
		var configuration config.Configuration

		if err := config.Mattermost.LoadPluginConfiguration(&configuration); err != nil {
			logger.Error("Error occurred during loading plugin configuraton from Mattermost", err, nil)
			return err
		}

		if err := configuration.IsValid(); err != nil {
			return err
		}

		if err := configuration.ProcessConfiguration(); err != nil {
			config.Mattermost.LogError(err.Error())
			return err
		}

		config.SetConfig(&configuration)
	}
	return nil
}

func (p *Plugin) RegisterCommands() error {
	for _, c := range command.Commands {
		if err := config.Mattermost.RegisterCommand(c.Command); err != nil {
			logger.Error("Cound't register command", err, map[string]interface{}{"command": c.Command.Trigger})
			return err
		}
	}

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	// cant use strings.split as it includes empty string if deliminator
	// is the last character in input string
	split, argErr := util.SplitArgs(args.Command)
	if argErr != nil {
		return &model.CommandResponse{
			Type: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: argErr.Error(),
		}, nil
	}

	function := split[0]
	var params []string

	if len(split) > 1 {
		params = split[1:]
	}

	commandConfig := command.Commands[function]
	if commandConfig == nil {
		return nil, &model.AppError{Message: "Unknown command: [" + function + "] encountered"}
	}

	context := p.prepareContext(args)
	if response, err := commandConfig.Validate(params, context); response != nil {
		return response, err
	}

	// todo add error logs here
	return commandConfig.Execute(params, context)
}

func (p *Plugin) prepareContext(args *model.CommandArgs) command.Context {
	return command.Context{
		CommandArgs: args,
		Props:       make(map[string]interface{}),
	}
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	conf := config.GetConfig()
	if err := conf.IsValid(); err != nil {
		logger.Error("Plugin is not configured", err, nil)
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	d := util.DumpRequest(r)
	endpoint := controller.GetEndpoint(r)

	if endpoint == nil {
		p.handler.ServeHTTP(w, r)
	} else if !endpoint.RequiresAuth || controller.Authenticated(w, r) {
		if err := endpoint.Execute(w, r); err != nil {
			logger.Error("Error occurred processing "+r.URL.String(), err, map[string]interface{}{"request": string(d)})
			raven.CaptureError(err, nil)
		}
	}
}

func (p *Plugin) Run() {
	if !p.running {
		p.running = true
		p.runner()
	}
}

func (p *Plugin) runner() {
	go func() {
		<-time.NewTimer(config.RunnerInterval).C
		if err := notification.SendNotificationsAndReports(); err != nil {
			config.Mattermost.LogError(err.Error())
		}
		if !p.running {
			return
		}
		p.runner()
	}()
}

func (p *Plugin) initSentry() error {
	var err error

	if enabled, _ := strconv.ParseBool(SentryEnabled); enabled {
		err = raven.SetDSN(SentryDSN)
	}

	raven.SetTagsContext(map[string]string{"pluginComponent": "server"})

	return err
}

func main() {
	defer profile.Start().Stop()
	plugin.ClientMain(&Plugin{})
}
