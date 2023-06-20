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
* [ ] Create tag on release branch as `vx.y.z`. This will trigger the GitHub actions job.
* [ ] Monitor the GitHub actions for success or faliure.
* [ ] Once the GitHub actions are complete, verify that a new GitHub release has been created. 
* [ ] Verify that the GitHub release contains plugin distributions in attachments.
* [ ] Use the changelog generated in earlier step along with [changelog template](#changelog-template), add the changelog to the GitHub release.
* [ ] Update plugin version number on `master` for next release.
* [ ] Close the release issue.
