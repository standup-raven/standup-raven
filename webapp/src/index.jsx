import React from 'react';
import ChannelHeaderButtonIcon from './components/channelHeaderButton';
import reducer from './reducer';
import Actions from './actions';
import StandupModal from './components/standupModal';
import ConfigModal from './components/configModal';
import Constants from './constants';
import * as Sentry from '@sentry/browser';
import utils from './utils';
import * as RavenClient from './raven-client';

class StandupRavenPlugin {
    // eslint-disable-next-line class-methods-use-this
    async initialize(registry, store) {
        const siteURL = utils.getValueSafely(store.getState(), 'entities.general.config.SiteURL');
        const pluginConfig = await RavenClient.Config.getPluginConfig(siteURL);

        if (pluginConfig.enableErrorReporting) {
            initSentry(pluginConfig.sentryWebappDSN);
        }

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            (channel) => store.dispatch(Actions.openStandupModal(channel.id)),
            Constants.PLUGIN_DISPLAY_NAME,
            Constants.PLUGIN_DISPLAY_NAME,
        );

        registry.registerRootComponent(StandupModal);
        registry.registerRootComponent(ConfigModal);
        registry.registerWebSocketEventHandler(
            `custom_${Constants.PLUGIN_NAME}_open_config_modal`,
            () => {
                store.dispatch(Actions.openConfigModal());
            },
        );
        registry.registerWebSocketEventHandler(
            `custom_${Constants.PLUGIN_NAME}_open_standup_modal`,
            () => {
                store.dispatch(Actions.openStandupModal());
            },
        );

        registry.registerWebSocketEventHandler(
            `custom_${Constants.PLUGIN_NAME}_add_active_channel`,
            (event) => {
                store.dispatch(Actions.addActiveChannel(event.data.channel_id));
            },
        );

        registry.registerWebSocketEventHandler(
            `custom_${Constants.PLUGIN_NAME}_remove_active_channel`,
            (event) => {
                store.dispatch(Actions.removeActiveChannel(event.data.channel_id));
            },
        );

        registry.registerReducer(reducer);
    }
}

function initSentry(dsn) {
    Sentry.init({
        dsn,
    });

    Sentry.configureScope(((scope) => {
        scope.setTag('pluginComponent', 'webapp');
    }));
}

window.registerPlugin(Constants.PLUGIN_NAME, new StandupRavenPlugin());
