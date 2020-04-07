const PLUGIN_NAME = 'standup-raven';

const PLUGIN_DISPLAY_NAME = 'Standup Raven';

const PLUGIN_BASE_URL = `plugins/${PLUGIN_NAME}`;

const PLUGIN_STATIC_DIR_URL = `${PLUGIN_BASE_URL}/static`;

const URL_SUBMIT_USER_STANDUP = `${PLUGIN_BASE_URL}/standup`;

const URL_STANDUP_CONFIG = `${PLUGIN_BASE_URL}/config`;

const URL_SPINNER_ICON = `${PLUGIN_STATIC_DIR_URL}/spinner.svg`;

const URL_GET_TIMEZONE = `${PLUGIN_BASE_URL}/timezone`;

const URL_PLUGIN_CONFIG = `${PLUGIN_BASE_URL}/plugin-config`;

const URL_ACTIVE_CHANNELS = `${PLUGIN_BASE_URL}/active-channels`;

const MATTERMOST_CSRF_COOKIE = 'MMCSRF';

const ACTIONS = {
    OPEN_STANDUP_MODAL: `${PLUGIN_NAME}_open_standup_modal`,
    CLOSE_STANDUP_MODAL: `${PLUGIN_NAME}_close_standup_modal`,
    OPEN_CONFIG_MODAL: `${PLUGIN_NAME}_open_config_modal`,
    CLOSE_CONFIG_MODAL: `${PLUGIN_NAME}_close_config_modal`,
    ADD_ACTIVE_CHANNEL: `${PLUGIN_NAME}_add_active_channel`,
    REMOVE_ACTIVE_CHANNEL: `${PLUGIN_NAME}_remove_active_channel`,
};

export default {
    URL_SUBMIT_USER_STANDUP,
    URL_STANDUP_CONFIG,
    URL_SPINNER_ICON,
    ACTIONS,
    PLUGIN_NAME,
    PLUGIN_DISPLAY_NAME,
    MATTERMOST_CSRF_COOKIE,
    URL_GET_TIMEZONE,
    URL_PLUGIN_CONFIG,
    URL_ACTIVE_CHANNELS,
};
