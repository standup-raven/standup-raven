<img src="assets/images/banner.png" width="300px">

#

## 🐞 Sentry Configuration

Developers may plugin in their own, organization-owned Sentry accounts with Standup Raven for internal error monitoring. This is very helpful for reporting errors back here. 

To configure Sentry, you need to update the `sentry` section of `build_properties.json`. The template already exists in the file, so you just need to fill in the details as described below -

* **enabled** - true | false - Enables or disables Sentry for the build
* **dsn** - Your Sentry project DSN with public and private keys. This is used with the server component of the plugin. This will soon be replaced with the public DSN once the updated Sentry Go client is released.
* **publicDsn** - Your Sentry public DSN. This is used with the webapp component of the plugin.
* **server_url** - Sentry server URL. Use `https://sentry.io` if using the Sentry-hosted cloud instance or your own Sentry instance URL.
* **org** - Sentry organization name to use.
* **project** - Sentry project name to use.
* **auth_token** - Sentry auth token. This is used by the Sentry CLI to upload webapp sources to Sentry.
