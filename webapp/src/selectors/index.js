import Constants from '../constants';

const getPluginState = (state) => state[`plugins-${Constants.PLUGIN_NAME}`] || {};

export const isStandupModalVisible = (state) => getPluginState(state).standupModalVisible || false;

export const isConfigModalVisible = (state) => getPluginState(state).configModalVisible || false;

export default {
    isStandupModalVisible,
    isConfigModalVisible,
};
