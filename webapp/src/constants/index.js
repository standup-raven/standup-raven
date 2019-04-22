import Utils from '../utils';

const PLUGIN_NAME = 'standup-raven';

const PLUGIN_DISPLAY_NAME = 'Standup Raven';

const PLUGIN_BASE_URL = `${Utils.getBaseURL()}/plugins/${PLUGIN_NAME}`;

const PLUGIN_STATIC_DIR_URL = `${PLUGIN_BASE_URL}/static`;

const URL_SUBMIT_USER_STANDUP = `${PLUGIN_BASE_URL}/standup`;

const URL_STANDUP_CONFIG = `${PLUGIN_BASE_URL}/config`;

const URL_PLUGIN_ICON = `${PLUGIN_STATIC_DIR_URL}/logo.png`;

const URL_SPINNER_ICON = `${PLUGIN_STATIC_DIR_URL}/spinner.svg`;

const ACTIONS = {
    OPEN_STANDUP_MODAL: `${PLUGIN_NAME}_open_standup_modal`,
    CLOSE_STANDUP_MODAL: `${PLUGIN_NAME}_close_standup_modal`,
    OPEN_CONFIG_MODAL: `${PLUGIN_NAME}_open_config_modal`,
    CLOSE_CONFIG_MODAL: `${PLUGIN_NAME}_close_config_modal`,
};

export default {
    URL_SUBMIT_USER_STANDUP,
    URL_STANDUP_CONFIG,
    URL_PLUGIN_ICON,
    URL_SPINNER_ICON,
    ACTIONS,
    PLUGIN_NAME,
    PLUGIN_DISPLAY_NAME,
};
