## 🏌️‍♀️Deployment

The plugin can be deployed to Mattermost directly via the `deploy` make command. You need to expose the following
environment variable for it to work -

    $ export MM_SERVICESETTINGS_SITEURL="<mattermost-site0url>"; \
    export MM_ADMIN_USERNAME="<username-to-upload-via>"; \
    read -s MM_ADMIN_PASSWORD; export MM_ADMIN_PASSWORD; \
    export PLATFORM="<target-mattermost-platform>";

## 🖊 Usage Instructions

1. Create a channel for your team standup or use an existing one.

1. Add configurations for your standup -

        /standupconfig
        
    this opens a modal where you can enter your channel's configurations.

1. Add members to standup -

        /standupaddmembers <usernames...>
        
    Usernames can be specified as @ mentions.
    
1. You may verify saved config if you want by executing -

        /standupviewconfig
        
1. Fill your standup by clicking on the Standup Raven icon in the channel header bar. The icon may be hidden in ellipsis icon.

    ![](assets/images/channel_header_button.png)
    
1. Execute help command anytime to access plugin commands help -

        /standuphelp 
