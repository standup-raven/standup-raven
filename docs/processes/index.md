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

* Code freeze to be done on this day. No more changes will be included in this release other than RC bug fixes.
* Create a new branch from master named `release-x.y.z` where `x.y.z` is the version being released.
* Perform end to end testing on minimum supported Mattermost version following [recommended testing process](#testing-process).




### Testing Process
