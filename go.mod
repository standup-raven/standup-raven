module github.com/standup-raven/standup-raven

go 1.14

require (
	bou.ke/monkey v1.0.2
	github.com/dustin/go-humanize v1.0.0
	github.com/getsentry/sentry-go v0.7.0
	github.com/mattermost/mattermost-server/v5 v5.26.1
	github.com/mitchellh/gox v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/teambition/rrule-go v1.6.0
	github.com/thoas/go-funk v0.7.0
	go.uber.org/atomic v1.6.0
)

replace github.com/teambition/rrule-go => github.com/standup-raven/rrule-go v1.5.1-0.20200606021409-a2ced8306e77
