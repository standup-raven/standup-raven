<img src="assets/images/banner.png" width="300px">

# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method, with the owners of this repository before making a change. 

Please note that we have a code of conduct, please follow it in all your interactions with the project.

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a 
   build.
2. Update the README.md with details of changes to the interface, this includes new environment 
   variables, exposed ports, useful file locations and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this
   Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you 
   do not have permission to do that, you may request the second reviewer to merge it for you.

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as
contributors and maintainers pledge to participate in our project and
our community a harassment-free experience for everyone, regardless of age, body
size, disability, ethnicity, gender identity and expression, level of experience,
nationality, personal appearance, race, religion, or sexual identity and
orientation.

### Our Standards

Examples of behavior that contributes to creating a positive environment
include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other community members

Examples of unacceptable behavior by participants include:

* The use of sexualized language or imagery and unwelcome sexual attention or
advances
* Trolling, insulting/derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or electronic
  address, without explicit permission
* Other conduct that could reasonably be considered inappropriate in a
  professional setting

### Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable
behavior and are expected to take appropriate and fair corrective action in
response to any instances of unacceptable behavior.

Project maintainers have the right and responsibility to remove, edit, or
reject comments, commits, code, wiki edits, issues, and other contributions
that are not aligned to this Code of Conduct, or to ban temporarily or
permanently any contributor for other behaviors that they deem inappropriate,
threatening, offensive, or harmful.

### Scope

This Code of Conduct applies both within project spaces and in public spaces
when an individual is representing the project or its community. Examples of
representing a project or community includes using an official project e-mail
address, posting via an official social media account, or acting as an appointed
representative at an online or offline event. Representation of a project may be
further defined and clarified by project maintainers.

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported by contacting the project team at [INSERT EMAIL ADDRESS]. All
complaints will be reviewed and investigated and will result in a response that
is deemed necessary and appropriate to the circumstances. The project team is
obligated to maintain confidentiality with regard to the reporter of an incident.
Further details of specific enforcement policies may be posted separately.

Project maintainers who do not follow or enforce the Code of Conduct in good
faith may face temporary or permanent repercussions as determined by other
members of the project's leadership.

### Workflow

This is a general workflow for anyone who would like to contribute to the Standup Raven Mattermost plugin.

1. Review the repository structure

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
1. Once a developer has reviewed your pull request, it will either be accepted straight away, or you may receive some review comments.
Review comments help up communicate with you over the code changes you've proposed to make sure we understand each other.
1. Once the review comments are implemented, reply and resolve the comments and a developer will take a second look. Rinse and repeat.
1. On acceptance of your changes, your pull request will be merged and made part of the new release. Your changes will impact 
the thousands (someday) of users in their daily endeavor to fill standups.

### 3rd Party Service URL

* Codecov - https://codecov.io/gh/standup-raven/standup-raven

### Attribution

This Code of Conduct is adapted from the [Contributor Covenant][homepage], version 1.4,
available at [http://contributor-covenant.org/version/1/4][version]

[homepage]: http://contributor-covenant.org
[version]: http://contributor-covenant.org/version/1/4/

