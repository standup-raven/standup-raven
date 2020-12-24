module github.com/standup-raven/standup-raven

go 1.15

require (
	bou.ke/monkey v1.0.2
	github.com/dustin/go-humanize v1.0.0
	github.com/getsentry/sentry-go v0.7.0
	github.com/go-ldap/ldap v3.0.3+incompatible // indirect
	github.com/mattermost/mattermost-plugin-api v0.0.12
	github.com/mattermost/mattermost-server v5.11.1+incompatible
	github.com/mattermost/mattermost-server/v5 v5.27.0
	github.com/nicksnyder/go-i18n v1.10.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/teambition/rrule-go v1.6.0
	github.com/thoas/go-funk v0.7.0
	go.uber.org/atomic v1.6.0
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
)

replace github.com/teambition/rrule-go => github.com/standup-raven/rrule-go v1.5.1-0.20200606021409-a2ced8306e77
