<img src="assets/images/banner.png" width="300px">

#

## ⁉ Troubleshooting

* ##### I submitted my standup but it's not showing up in the report.

    Make sure you submit the standup, in the same channel as the report gets generated in.

* ##### I filled out my standup report but had to make some changes to it. However, I'm seeing a blank standup form on opening it.

    Make sure you are in the same channel as when you originally filled out your standup.
    
* ##### I'm seeing the "Standup not configured for this channel" message on opening the standup modal.

    Make sure you are filling the standup in the right channel or that the standup has been configured in the channel.
    
* ##### I'm seeing the "You are not a part of this channel's standup" message on opening the standup modal. 

    You are not part of the current channel's standup. Make sure you are filling standup in the right channel or that you were correctly added to the channel's standup.
    
* ##### I'm seeing the "No members configured for this channel's standup" message on opening standup modal.

    Make sure you've added some members to the channel's standup.
    
* ##### I'm seeing the "Standup is disabled for this channel" message on opening the standup modal.

    The channel standup is disabled. Run `/standup config` and enable standup from the modal that opens.
    
* ##### I run `/standup`, but no dialog popups for submitting standup.

    Verify that you have set the value for [`Site URL`](https://docs.mattermost.com/administration/config-settings.html#site-url) setting in your Mattermost server configuration. Also verify that the value is correct (should contain protocol, host and port).

    ![Site URL Verification Demo](/docs/assets/images/test-live-url.gif)
    
* ##### I run `/standup config` but no configuration doalog opens up.

    Verify that you have set the value for [`Site URL`](https://docs.mattermost.com/administration/config-settings.html#site-url) setting in your Mattermost server configuration. Also verify that the value is correct (should contain protocol, host and port).

    ![Site URL Verification Demo](/docs/assets/images/test-live-url.gif)

* ##### I think the plugin is awesome and super cool.

    Hey, that's not a problem! It was designed that way 😎. 
