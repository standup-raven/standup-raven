<div align="center">

<img src="docs/assets/images/banner.png?raw=true" width="70%" max-width="1500px"></img>
#
[![CircleCI](https://circleci.com/gh/standup-raven/standup-raven/tree/master.svg?style=svg)](https://circleci.com/gh/standup-raven/standup-raven/tree/master)
[![codecov](https://codecov.io/gh/standup-raven/standup-raven/branch/master/graph/badge.svg)](https://codecov.io/gh/standup-raven/standup-raven)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/934cb67ed24e42978273489ae17bddef)](https://www.codacy.com/app/harshilsharma/standup-raven?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=standup-raven/standup-raven&amp;utm_campaign=Badge_Grade)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2940/badge)](https://bestpractices.coreinfrastructure.org/projects/2940)

A Mattermost plugin for communicating daily standups across team

<a href="https://www.buymeacoffee.com/harshilsharma63" target="_blank"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>

</div>

<div align="center">
    <img src="docs/assets/images/standup.gif?raw=true"></img>
</div>

## ✨ Features

* Configurable standup window per channel for standup reminders

* Automatic window open reminders

    ![](docs/assets/images/first-reminder.png | width=380)
    
* Automatic window close reminders

    ![](docs/assets/images/second-reminder.png)
    
* Per-channel customizable

    ![](docs/assets/images/config-general.png)
    
    ![](docs/assets/images/config-notifications.png)
    
    ![](docs/assets/images/config-schedule.png)
    
* Automatic standup reports
    
    ![](docs/assets/images/report-user-aggregated.png)

* Multiple standup report formats -

  * User Aggregated - tasks aggregated by individual users

    ![](docs/assets/images/report-user-aggregated.png)
     
  * Type Aggregated - tasks aggregated by type

    ![](docs/assets/images/report-type-aggregated.png)

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

Be responsible and report them to raven@joshtechnologygroup.com .

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

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
