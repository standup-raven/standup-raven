<img src="assets/images/banner.png" width="300px">

#

## <img src="https://sentry-brand.storage.googleapis.com/sentry-glyph-black.png" width="40px"> Sentry Configuration

Developers may plugin in their own, organization-owned Sentry accounts with Standup Raven for internal error monitoring. This is very helpful for reporting errors back over here. 

To configure sentry you need to update the `sentry` section of `build_properties.json`. The template already exists in the file so you just need to fill in the details as described below -

* **enabled** - true | false - enables or disables Sentry for the build
* **dsn** - your Sentry project DSN. You need to use the deprecated DSN, the one containing public and private keys. This is the only DSN supported by the Sentry Go client.
* **** 
