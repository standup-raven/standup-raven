<img src="assets/images/banner.png" width="300px">

#

## Contribution Guidelines

### Workflow

This is a general workflow for anyone who would like to contribute to Standup Raven Mattermost plugin.

1. Review the respository structure

        .
        ├── docs - contains all documentations, including this file
        ├── server  - contains server conponent the plugin
        │   ├── command - contains logic for all the slash commands
        │   ├── config - contains all plugin configurations
        │   ├── controller - contains all HTTP endpoints, used by webapp componnet to communicate
        │   ├── logger - contains logger binded with Sentry API
        │   ├── otime - custom time class with functions for printing time in specific formats
        │   ├── standup - contains the core logic for all standup management
        │   ├── util - general, abstract utilities
        └── webapp - contains webapp component of the plugin
            └── src
                ├── actions - contains the Redux actions
                ├── assets - contains miscellaneous assets such as images
                ├── components - contains the custom components used in webapp
                ├── constants - contains general constants such as HTTP endpoint paths
                ├── reducer - contains Redux reducers
                ├── selectors - contains Redux selectors
                └── utils - general, abstract utilities
                
1. On your fork, create a branch `GH-###` where `###` is the GitHub issue ID.
1. For any code changes, make sure you write or modify test cases where ever appropriate.
1. Run unit tests by running `make test` from the project root. Make sure all tests are passing before submitting a pull request.
1. Run style check by running `make check-style` from the project root. Make sure your modifications comply with the set up style guidelines.
1. Create a pull request from `your_fork/GH-###` to `standup-raven/master`. Make sure to update the checklist
in the pull request template.
1. Once a developer has reviewed your pull request, it will either be accepted straight away, of you may receive some review comments.
Review comments help up communicate with you over the code changes you've proposed to make sure we understand each other.
1. Once the review comments are implemented, reply and resolve on the comments and a developer will take a second look. Rinse and repeat.
1. On acceptance of your changes, your pull request will be merged and made part of new release. Your changes will impact 
the thousands (,some day) of users in their daily endavour to fill standups.

### 3rd Party Service URL

* CircleCI - https://circleci.com/gh/standup-raven/standup-raven
* Codecov - https://codecov.io/gh/standup-raven/standup-raven
* Codacy - https://app.codacy.com/project/harshilsharma/standup-raven/dashboard

