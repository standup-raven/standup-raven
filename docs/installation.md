<img src="assets/images/banner.png" width="300px">

#

## ⬇ Installation

Upload the plugin binary for your platform in Mattermost `System Console` > `Plugins (BETA)` > `Management`.

## Upgrade Instructions

1. Make sure to backup Standup Raven data from the query specified below
using your organization's standard procedures for backing up MySQL or PostgreSQL

    PostgreSQL -
    
        select * from pluginkeyvaluestore where pluginid='standup-raven';
        
    MySQL -
    
        select * from PluginKeyValueStore where pluginid="standup-raven";
