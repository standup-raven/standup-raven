<div align="center">

<img src="docs/assets/images/banner.png?raw=true" width="70%" max-width="1500px"></img>

#
A Mattermost plugin for communicating daily standups across teams


[![CircleCI](https://circleci.com/gh/standup-raven/standup-raven/tree/master.svg?style=svg)](https://circleci.com/gh/standup-raven/standup-raven/tree/master)
[![codecov](https://codecov.io/gh/standup-raven/standup-raven/branch/master/graph/badge.svg)](https://codecov.io/gh/standup-raven/standup-raven)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2940/badge)](https://bestpractices.coreinfrastructure.org/projects/2940)

</div>

<div align="center">
    <img src="docs/assets/images/standup.gif?raw=true"></img>
</div>

## ✨ Features

* Configurable standup window per channel for standup reminders

* Automatic window open reminders

    <img src="docs/assets/images/first-reminder.png?raw=true" width="380px"></img>
    
* Automatic window close reminders

    <img src="docs/assets/images/second-reminder.png?raw=true" width="600px"></img>
    
* Per-channel customizable

    <img src="docs/assets/images/config-general.png?raw=true" width="650px"></img>
    
    <img src="docs/assets/images/config-notifications.png?raw=true" width="650px"></img>
    
    <img src="docs/assets/images/config-schedule.png?raw=true" width="650px"></img>
    
* Automatic standup reports
    
    <img src="docs/assets/images/report-user-aggregated.png?raw=true" width="500px"></img>

* Multiple standup report formats -

  * User Aggregated - Tasks aggregated by individual users

    <img src="docs/assets/images/report-user-aggregated.png?raw=true" width="500px"></img>
     
  * Type Aggregated - Tasks aggregated by type

    <img src="docs/assets/images/report-type-aggregated.png?raw=true" width="500px"></img>

* Ability to preview a standup report without publishing it in the channel
* Ability to manually generate standup reports for any arbitrary date

## 🧰 Functionality

* Customize standup sections on per-channel basis, so team members can make it suite their style.

* Multiple report formats to choose from.

* Receive a window open notification at the configured window open time to remind you to fill your standup.

* Receive a reminder at the completion of 80% of the configured window duration to remind you to fill your standup. 
This message tags those members who haven't yet filled their standup.

* Receive an auto-generated standup report at the end of the configured window close time. 
The generated standup contains the names of members who have yet to fill their standup.

* Allow or restrict standup configuration modification to channel admins (Requires Mattermost EE).

## Guides

### User Guide

* 👩‍💼 [User Guide](docs/user_guide.md)

### Developer Guide

* 🚦 [Getting Started](docs/getting_started.md)

* 🐞 [Integrating Sentry](docs/sentry.md)

### Ops Guide

* ⬇ [Installing](docs/installation.md)

* 🏌️‍♀️[️Deployment](docs/deployment.md)

* ⚙ [Plugin Configurations](docs/configuration.md)

* ⁉ [Troubleshooting](docs/troubleshooting.md)

### TODO

* [x] Permissions
* [ ] Vacation
* [ ] Periodic reports

## Reporting Security Vulnerabilities

Due to the sensitive nature of such vulnerabilities, please refrain from posting them publically
over GitHub issues or any other medium.

Be responsible and report them to hello@standupraven.com .

## 🌟 Attribution

<div>Project logo (the Raven) is made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/jatinjtg"><img src="https://avatars2.githubusercontent.com/u/50952137?v=4" width="100px;" alt="jatinjtg"/><br /><sub><b>jatinjtg</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=jatinjtg" title="Code">💻</a> <a href="https://github.com/Harshil Sharma/Standup Raven/issues?q=author%3Ajatinjtg" title="Bug reports">🐛</a> <a href="#ideas-jatinjtg" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=jatinjtg" title="Documentation">📖</a> <a href="#infra-jatinjtg" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=jatinjtg" title="Tests">⚠️</a></td>
    <td align="center"><a href="https://github.com/goku321"><img src="https://avatars2.githubusercontent.com/u/12848015?v=4" width="100px;" alt="Deepak Sah"/><br /><sub><b>Deepak Sah</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=goku321" title="Code">💻</a></td>
    <td align="center"><a href="http://sandipagarwal.in"><img src="https://avatars0.githubusercontent.com/u/988003?v=4" width="100px;" alt="Sandip Agarwal"/><br /><sub><b>Sandip Agarwal</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=sandipagarwal" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/chetanyakan"><img src="https://avatars3.githubusercontent.com/u/35728906?v=4" width="100px;" alt="Chetanya Kandhari"/><br /><sub><b>Chetanya Kandhari</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=chetanyakan" title="Code">💻</a> <a href="https://github.com/Harshil Sharma/Standup Raven/issues?q=author%3Achetanyakan" title="Bug reports">🐛</a> <a href="#ideas-chetanyakan" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=chetanyakan" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/ayadav"><img src="https://avatars2.githubusercontent.com/u/154998?v=4" width="100px;" alt="Amit Yadav"/><br /><sub><b>Amit Yadav</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=ayadav" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/SezalAgrawal"><img src="https://avatars1.githubusercontent.com/u/13785694?v=4" width="100px;" alt="SezalAgrawal"/><br /><sub><b>SezalAgrawal</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=SezalAgrawal" title="Code">💻</a></td>
    <td align="center"><a href="http://TheodoreLindsey.io"><img src="https://avatars3.githubusercontent.com/u/6985440?v=4" width="100px;" alt="Theodore S Lindsey"/><br /><sub><b>Theodore S Lindsey</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=RagingRoosevelt" title="Code">💻</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/Amalkh5"><img src="https://avatars2.githubusercontent.com/u/20528562?v=4" width="100px;" alt="Amal Alkhamees"/><br /><sub><b>Amal Alkhamees</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/commits?author=Amalkh5" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/henzai"><img src="https://avatars2.githubusercontent.com/u/25699758?v=4" width="100px;" alt="henzai"/><br /><sub><b>henzai</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/issues?q=author%3Ahenzai" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://www.hardwario.com/"><img src="https://avatars0.githubusercontent.com/u/19538414?v=4" width="100px;" alt="Pavel Hübner"/><br /><sub><b>Pavel Hübner</b></sub></a><br /><a href="#ideas-hubpav" title="Ideas, Planning, & Feedback">🤔</a> <a href="#userTesting-hubpav" title="User Testing">📓</a> <a href="#talk-hubpav" title="Talks">📢</a></td>
    <td align="center"><a href="https://github.com/tgly307"><img src="https://avatars3.githubusercontent.com/u/25153311?v=4" width="100px;" alt="tgly307"/><br /><sub><b>tgly307</b></sub></a><br /><a href="#ideas-tgly307" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/Harshil Sharma/Standup Raven/issues?q=author%3Atgly307" title="Bug reports">🐛</a></td>
    <td align="center"><a href="http://tzonkovs.wixsite.com/alex"><img src="https://avatars1.githubusercontent.com/u/4975715?v=4" width="100px;" alt="Alex Tzonkov"/><br /><sub><b>Alex Tzonkov</b></sub></a><br /><a href="#ideas-attzonko" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/Harshil Sharma/Standup Raven/issues?q=author%3Aattzonko" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/sonam-singh"><img src="https://avatars1.githubusercontent.com/u/10594597?v=4" width="100px;" alt="Sonam Singh"/><br /><sub><b>Sonam Singh</b></sub></a><br /><a href="https://github.com/Harshil Sharma/Standup Raven/issues?q=author%3Asonam-singh" title="Bug reports">🐛</a> <a href="#ideas-sonam-singh" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/skyscooby"><img src="https://avatars1.githubusercontent.com/u/4450935?v=4" width="100px;" alt="Andrew Greenwood"/><br /><sub><b>Andrew Greenwood</b></sub></a><br /><a href="#ideas-skyscooby" title="Ideas, Planning, & Feedback">🤔</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/mihai-satmarean"><img src="https://avatars0.githubusercontent.com/u/4729542?v=4" width="100px;" alt="mihai-satmarean"/><br /><sub><b>mihai-satmarean</b></sub></a><br /><a href="#ideas-mihai-satmarean" title="Ideas, Planning, & Feedback">🤔</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind are welcome!
