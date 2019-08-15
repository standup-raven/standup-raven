<img src="assets/images/banner.png" width="300px">

#

## 👩‍💼 User Guide

Once the plugin is installed in your Mattermost instance, enabling teams to use it is super easy.
Just follow these steps and you'll be ready in no time.

1. **Creating channels for standup** - create a new channel, or use existing one, for each team who wants to use Standup Raven for their standup.

1. **Configuring channel standup** - for each channel any member can enter configurations for the channel standup. If you are on Mattermost Enterprise Edition and have *Permission Schema* enabled, only a channel admin, team admin or a system admin can perform this operation.
    
    Running the following slash command allows specifying team specific settings -
    
        /standup config
        
    In the dialog box presented, following settings are required -
    
    * **Status** - `Enabled` to enable standup for your channel or `Disable` to disable it.
    
    * **Window open time** - the time at which standup reminder will be sent in the channel.
    
    * **Window close time** - the time at which an automated standup report will be sent in the channel. The report
    will include standups for all members who have filled their standups until this time.
    An additional reminder notification is sent in the channel at 80% completion of the window duration.
    This message tags members who have not yet filled their standups.
    
    * **Timezone** - channel specific timezone to follow for standup notifications.
    
    * **Window Open Reminder** - enable of disable window open reminder.
    
    * **Window Close Reminder** - enable or disable window close reminder.
     
    * **Sections** - sections defile the type of tasks users will fill in their standup.
    For example, if your team fills their standup at beginning of their work day, suggested sections would be
    `yesterday`, `today` and maybe `blockers`.
        
        At least one section is required to be specified.
        
1. **Saving standup config** - save the standup config that you filled.

1. **Adding standup members** - the following slash command allows you to add members to the channel standup -

        /standup addmembers
        
    You can specify multiple members together, separated by space. Members who are not present in the channel will
    be automatically added to the channel as well.
    
1. **Filling your standup** - once all configuration is complete, click on the Standup Raven button in
    channel header to bring up a modal for filling your standup.
    
    Once saved, you can click on the Standup Raven button again to bring back your filled standup, allowing you
    to make updates to it.
     
