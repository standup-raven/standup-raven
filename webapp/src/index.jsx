import React from 'react';
import ChannelHeaderButtonIcon from './components/channelHeaderButton';
import reducer from './reducer';
import Actions from './actions';
import StandupModal from './components/standupModal';
import ConfigModal from './components/configModal';
import Constants from './constants';
import * as Sentry from '@sentry/browser';
const buildProperties = require('../../build_properties.json');

class StandupRavenPlugin {
    // eslint-disable-next-line class-methods-use-this
    initialize(registry, store) {
        if (buildProperties.sentryEnabled) {
            initSentry();
        }

        registry.registerRootComponent(StandupModal);
        registry.registerRootComponent(ConfigModal);
        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            (channel) => store.dispatch(Actions.openStandupModal(channel.id)),
            Constants.PLUGIN_DISPLAY_NAME,
            Constants.PLUGIN_DISPLAY_NAME,
        );
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

        registry.registerReducer(reducer);
    }
}

function initSentry() {
    Sentry.init({
        dsn: buildProperties.sentry.publicDsn,
    });

    Sentry.configureScope(((scope) => {
        scope.setTag('pluginComponent', 'webapp');
    }));
}

window.registerPlugin(Constants.PLUGIN_NAME, new StandupRavenPlugin());
