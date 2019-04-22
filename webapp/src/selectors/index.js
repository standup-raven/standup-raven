import Constants from '../constants';

const getPluginState = (state) => state[`plugins-${Constants.PLUGIN_NAME}`] || {};

export const isStandupModalVisible = (state) => getPluginState(state).standupModalVisible;

export const isConfigModalVisible = (state) => getPluginState(state).configModalVisible;

export default {
    isStandupModalVisible,
    isConfigModalVisible,
};
