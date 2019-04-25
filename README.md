<div align="center">

<img src="docs/assets/images/banner.png?raw=true" width="70%" max-width="1500px"></img>
#
[![CircleCI](https://circleci.com/gh/standup-raven/standup-raven/tree/master.svg?style=svg)](https://circleci.com/gh/standup-raven/standup-raven/tree/master)
[![codecov](https://codecov.io/gh/standup-raven/standup-raven/branch/master/graph/badge.svg)](https://codecov.io/gh/standup-raven/standup-raven)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/934cb67ed24e42978273489ae17bddef)](https://www.codacy.com/app/harshilsharma/standup-raven?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=standup-raven/standup-raven&amp;utm_campaign=Badge_Grade)

A Mattermost plugin for communicating daily standups across team

</div>

<div align="center">
    <img src="docs/assets/images/standup.gif?raw=true"></img>
</div>

## ✨ Features

* Configurable standup window per channel for standup reminders

* Automatic window open reminders

    ![](docs/assets/images/window_open_notification.png)
    
* Automatic window close reminders

    ![](docs/assets/images/window_close_notification.png)
    
* Per-channel customizable

    ![](docs/assets/images/standup_config.png)
    
* Automatic standup reports
    
    ![](docs/assets/images/user_aggregated_report.png)

* Multiple standup report formats -

  * User Aggregated - tasks aggregated by individual users

    ![](docs/assets/images/user_aggregated_report.png)
     
  * Type Aggregated - tasks aggregated by type

    ![](docs/assets/images/type_aggregated_report.png)

* Ability to preview standup report without publishing it in channel
* Ability to manually generate standup reports for any arbitrary date

## 🧰 Functionality

* Customize standup sections on per-channel basis, so team members can make it suite their style.

* Multiple report formats to choose from.

* Receive a window open notification at the configured window open time to remind filling your standup.

* Receive a reminder at completion of 80% of configured window duration to remind filling your standup. 
This message tags members who haven't yet filled their standup.

* Receive auto-generated standup report at the end of configured window close time. 
The generated standup contains names of members who have yet not filled their standup.

### TODO

* [ ] Permissions
* [ ] Vacation
* [ ] Periodic reports

## 🚦 Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

See [deployment notes](#%EF%B8%8F%EF%B8%8Fdeployment) on how to deploy the project on a live system.

### 🔑 Prerequisites

Set up your development environment for building, running, and testing the Standup Raven.

#### 👨‍💻 Obtaining Source

    # TODO - test this, and update to use new repo path
    $ go get -u github.com/standup-raven/mattermost-standup-plugin/..

#### Go

Requires go version 1.12

    https://golang.org/doc/install
    
#### NodeJS

Recommended NodeJS version 10.11

    https://nodejs.org/download/release/v10.11.0/

#### Make

On Ubuntu -

    $ sudo apt-get install build-essential
    
On MacOS, install XCode command line tools. 

#### HTTPie

You need this only if you want to use `$ make deploy` for deployments to Mattermost instance.

On MacOS

    $ brew install httpie
    
On Ubuntu

    $ apt-get install httpie
    
For other platform refer the [official installation guide](https://github.com/jakubroztocil/httpie#id3).

### 👨‍💻 Building

Once you have cloned the repo in the correct path - `$GOPATH/JoshLabs`, simply run `$ make dist` from the cloned repo.

This will produce three artifacts in `/dist` directory -

| Flavor  | Distribution |
|-------- | ------------ |
| Linux   | `mattermost-plugin-standup-raven-vx.y.z-linux-amd64.tar.gz`  |
| MacOS   | `mattermost-plugin-standup-raven-vx.y.z-darwin-amd64.tar.gz` |
| Windows | `mattermost-plugin-standup-raven-vx.y.z-windows-amd64.tar.gz`|

This will also install, Glide - the Go package manager.

### ⬇ Installing

Upload the plugin binary for your platform in Mattermost `System Console` > `Plugins (BETA)` > `Management`. 

## 💯 Running Tests

Following command will run all server and webapp tests -

    $ make test
    
## 👞 Running Style Check

This will run server and webapp style checks -

    $ make check-style
    
You can also run style check for server and webapp individually

    $ make check-style-server # server style check
    $ make check-style-webapp # webapp style check
      

## 🏌️‍♀️Deployment

The plugin can be deployed to Mattermost directly via the `deploy` make command. You need to expose the following
environment variable for it to work -

    $ export MM_SERVICESETTINGS_SITEURL="<mattermost-site0url>"; \
    export MM_ADMIN_USERNAME="<username-to-upload-via>"; \
    read -s MM_ADMIN_PASSWORD; export MM_ADMIN_PASSWORD; \
    export PLATFORM="<target-mattermost-platform>";

## 🖊 Usage Instructions

1. Create a channel for your team standup or use an existing one.

1. Add configurations for your standup -

        /standupconfig
        
    this opens a modal where you can enter your channel's configurations.

1. Add members to standup -

        /standupaddmembers <usernames...>
        
    Usernames can be specified as @ mentions.
    
1. You may verify saved config if you want by executing -

        /standupviewconfig
        
1. Fill your standup by clicking on the Standup Raven icon in the channel header bar. The icon may be hidden in ellipsis icon.

    ![](docs/assets/images/channel_header_button.png)
    
1. Execute help command anytime to access plugin commands help -

        /standuphelp 

## ⚙ Plugin Configurations

* `Bot Username`: User account to be used for sending all automated posts from. It's recommended to use a separate, bot account for the purpose.
* `Time Zone`: The time zone your team is working in. This is to make sure all datetimes you enter are interpreted in your timezone and not in server's.
* `Work Week Start`: Day on which your work week starts.
* `Work Week End`: Day on which your work week ends.

## ⁉ Troubleshooting

* ##### I submitted my standup but it's not showing up in the report.

    Make sure you submitted the standup in the same channel as the report gets generated in.

* ##### I filled my standup report but had to make some changes to it. However, I'm seeing a blank standup form on opening it.

    Make sure you are in the same channel as you originally filled your standup it.
    
* ##### I'm seeing "Standup not configured for this channel" message on opening the standup modal.

    Make sure you are filling the standup in the right channel or that standup has been configured in the channel.
    
* ##### I'm seeing "You are not a part of this channel's standup" message on opening the standup modal. 

    You are not the part of current channel's standup. Make sure you are filling standup in the right channel or that you were correctly added to the channel's standup.
    
* ##### I'm seeing "No members configured for this channel's standup" message on opening standup modal.

    Make sure you've added some members to the channel's standup.
    
* ##### I'm seeing "Standup is disabled for this channel" message on opening standup modal.

    The channel standup is disabled. Run `/standupconfig` and enable standup from the modal that opens. 

* ##### I think the plugin is awesome and super cool.

    Hey, that's not a problem! It was designed that way 😎. 

## 🌟 Attribution

<div>Project logo (the Raven) is made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
