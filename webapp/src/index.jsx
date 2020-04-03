import React from 'react';
import ChannelHeaderButtonIcon from './components/channelHeaderButton';
import reducer from './reducer';
import Actions from './actions';
import StandupModal from './components/standupModal';
import ConfigModal from './components/configModal';
import Constants from './constants';
import * as Sentry from '@sentry/browser';
import utils from './utils';
import request from 'superagent';

const buildProperties = require('../../build_properties.json');

class StandupRavenPlugin {
    // eslint-disable-next-line class-methods-use-this
    async initialize(registry, store) {
        const pluginConfig = await getPluginConfig(store);

        if (!pluginConfig.disableChannelHeaderButton) {
            registry.registerChannelHeaderButtonAction(
                <ChannelHeaderButtonIcon/>,
                (channel) => store.dispatch(Actions.openStandupModal(channel.id)),
                Constants.PLUGIN_DISPLAY_NAME,
                Constants.PLUGIN_DISPLAY_NAME,
            );
        }

        if (buildProperties.sentryEnabled) {
            initSentry();
        }

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

async function getPluginConfig(store) {
    const siteURL = utils.getValueSafely(store.getState(), 'entities.general.config.SiteURL');
    const response = await request
        .get(`${siteURL}/${Constants.URL_PLUGIN_CONFIG}`)
        .withCredentials();
    return response.body;
}

window.registerPlugin(Constants.PLUGIN_NAME, new StandupRavenPlugin());
