## Standup Raven Release Process

This document contains instructions for releasing new version of Standup Raven.

The goal is to provide clear instructions and procedures for our entire release
process.

### Timeline

We release on **22nd of each month**, or the next working day in case of holiday. 
This is exactly one week after Mattermost's release schedule. The one week gap allows us
sufficient time to test on latest Mattermost release.

#### Daily Tasks

##### Mattermost Release Day

These are to be performed on the day Mattermost release is published.

1. Code freeze to be done on this day. No more changes will be included in this release other than RC bug fixes.
1. Create a new branch from master named `release-x.y.z` where `x.y.z` is the version being released.
1. Create release issue on GitHub based on the [release issue template](#release-issue-template).
1. Update release issue with release cut date. This is the date when release branch is first created on remote.
1. Update release issue with scheduled release date.

##### Intermediate Days

These are to be performed starting as soon as Mattermost release day tasks are complete, until 
Standup Raven scheduled release date.

1. Perform end to end testing on minimum supported Mattermost version following [recommended testing process](#testing-process).
1. Create GiHub issues for any bugs found.
1. Once testing is complete, begin with the bug fixes.
1. Cherrypick any bug fix which is applicable to `master` branch or other release.

##### Release Day

These are to be done on Standup Raven release day.

1. Generate release notes locally using [what-the-changelog](https://github.com/standup-raven/what-the-changelog). Keep this handy as it will be needed later.
1. Create tag on release branch. This will trigger the CircleCI release job.
1. Monitor the CircleCI job for success or faliure.
1. Once the CircleCI job is complete, verify that a new GitHub release has been created. 
1. Verify that the GitHub release contains plugin distributions in attachments.
1. Use the changelog generated in earlier step along with [changelog template](#changelog-template), add the changelog to the GitHub release.
1. Update plugin version number on `master` for next release.



### Testing Process

The plugin needs to be tested on all versions of Mattermost starting with the minimum up to the latest supported Mattermost server version version.

The plugin also needs to be tested on all platforms supported by Mattermost Plugins. These include
desktop web browsers, mobile web browsers and the official Mattermost dekstop apps.

Exact versions of supported platforms can be found over at [Mattermost client software requirements page](https://docs.mattermost.com/install/requirements.html#client-software).

#### Corner and Edge Cases

Add here any corner or edge cases to keep in mind while testing or during development.

1. Standup window closing at 23:59. 


### Release Issue Template 



### Changelog Template

    # Changelog  
    
    ###
    
    Please read the upgrade instructions before upgrading the plugin - [upgrade instructions](https://github.com/standup-raven/standup-raven/blob/master/docs/installation.md#upgrade-instructions)
    
    ##
    ##
