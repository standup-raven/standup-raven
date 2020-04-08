import Constants from '../constants';

const getPluginState = (state) => state[`plugins-${Constants.PLUGIN_NAME}`] || {};

export const isStandupModalVisible = (state) => getPluginState(state).standupModalVisible || false;

export const isConfigModalVisible = (state) => getPluginState(state).configModalVisible || false;

export const addedActiveChannel = (state) => getPluginState(state).addedActiveChannel || '';

export const removedActiveChannel = (state) => getPluginState(state).removedActiveChannel || '';

export default {
    isStandupModalVisible,
    isConfigModalVisible,
    addedActiveChannel,
    removedActiveChannel,
};
