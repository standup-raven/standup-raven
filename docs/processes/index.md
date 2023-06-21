<img src="../assets/images/banner.png" width="300px">

#

## Standup Raven Release Process

This document contains instructions for releasing new version of Standup Raven.

The goal is to provide clear instructions and procedures for our entire release
process.

### Timeline

We release on **22nd of each month**, or the next working day in case of holiday. 
This is exactly one week after Mattermost's release schedule. The one week gap allows us
sufficient time to test on latest Mattermost release.

### Release Process

The release process is embedded in [release issue template](#release-issue-template) itself 
as a checklist. This allows an easy way to make sure no points are skipped and also serves as a source 
of updates for the community.

### Testing Process

The plugin needs to be tested on all versions of Mattermost starting with the minimum up to the latest supported Mattermost server version version.

The plugin also needs to be tested on all platforms supported by Mattermost Plugins. These include
desktop web browsers, mobile web browsers and the official Mattermost desktop apps.

Exact versions of supported platforms can be found over at [Mattermost client software requirements page](https://docs.mattermost.com/install/requirements.html#client-software).

#### Corner and Edge Cases

Add here any corner or edge cases to keep in mind while testing or during development.

1. Standup window closing at 23:59. 


### Release Issue Template

    ### Release v<x.y.z>
    
    Scheduled Release Date: <scheduled release date>
    
    #### Daily Tasks
    
    ##### Mattermost Release Day
    
    These are to be performed on the day Mattermost release is published.
    
    * [ ] Code freeze to be done on this day. No more changes will be included in this release other than RC bug fixes.
    * [ ] Create a new branch from master named `release-x.y.z` where `x.y.z` is the version being released.
    * [ ] Create release issue on GitHub based on the [release issue template](https://github.com/standup-raven/standup-raven/blob/master/docs/processes/index.md#release-issue-template).
    * [ ] Update release issue with release cut date. This is the date when release branch is first created on remote.
    * [ ] Update release issue with scheduled release date.
    
    ##### Intermediate Days
    
    These are to be performed starting as soon as Mattermost release day tasks are complete, until 
    Standup Raven scheduled release date.
    
    * [ ] Perform end to end testing on minimum up to the latest supported Mattermost version following [recommended testing process](https://github.com/standup-raven/standup-raven/blob/master/docs/processes/index.md#testing-process).
    * [ ] Create GiHub issues for any bugs found.
    * [ ] Once testing is complete, begin with the bug fixes.
    * [ ] Cherrypick any bug fix which is applicable to `master` branch or other release.
    
    ##### Release Day
    
    These are to be done on Standup Raven release day.
    
    * [ ] Generate release notes locally using [what-the-changelog](https://github.com/standup-raven/what-the-changelog). Keep this handy as it will be needed later.
    * [ ] Create tag on release branch as `vx.y.z`. This will trigger the CircleCI release workflow.
    * [ ] Monitor the CircleCI workflow for success or faliure.
    * [ ] Once the CircleCI workflow is complete, verify that a new GitHub release has been created. 
    * [ ] Verify that the GitHub release contains plugin distributions in attachments.
    * [ ] Use the changelog generated in earlier step along with [changelog template](#changelog-template), add the changelog to the GitHub release.
    * [ ] Update plugin version number on `master` for next release.
    * [ ] Close the release issue.


### Changelog Template

    # Changelog  
    
    ###
    
    Please read the upgrade instructions before upgrading the plugin - [upgrade instructions](https://github.com/standup-raven/standup-raven/blob/master/docs/installation.md#upgrade-instructions)
    
    ##
    
    <changelog from what-the-changelog>
