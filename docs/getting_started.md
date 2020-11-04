<img src="assets/images/banner.png" width="300px">

#

## 🚦 Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

See [deployment notes](#%EF%B8%8F%EF%B8%8Fdeployment) on how to deploy the project on a live system.

### 🔑 Prerequisites

Set up your development environment for building, running, and testing the Standup Raven.

#### 👨‍💻 Obtaining Source

    $ git clone git@github.com:standup-raven/standup-raven.git

#### Go

Requires go version 1.15

    https://golang.org/doc/install
    
#### NodeJS

Recommended NodeJS version 10.17

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

Once you have fetched the repo, simply run `$ make dist` from the repo.

This will produce three artifacts in `/dist` directory -

| Flavor  | Distribution |
|-------- | ------------ |
| Linux   | `mattermost-plugin-standup-raven-vx.y.z-linux-amd64.tar.gz`  |
| MacOS   | `mattermost-plugin-standup-raven-vx.y.z-darwin-amd64.tar.gz` |
| Windows | `mattermost-plugin-standup-raven-vx.y.z-windows-amd64.tar.gz`|

This will also install, Glide - the Go package manager.

## 💯 Running Tests

Following command will run all server and webapp tests -

    $ make test
    
## 👞 Running Style Check

This will run server and webapp style checks -

    $ make check-style
    
You can also run style check for server and webapp individually

    $ make check-style-server # server style check
    $ make check-style-webapp # webapp style check
