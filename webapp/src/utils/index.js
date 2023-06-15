const channelAdminRole = 'channel_admin';
const teamAdminRole = 'team_admin';
const systemAdminRole = 'system_admin';
const systemGuestRole = 'system_guest';

/**
 * Returns the base url of the plugin
 * installation.
 *
 * @return {string} instance base URL
 */
function getBaseURL() {
    const url = new URL(window.location.href);
    return `${url.protocol}//${url.host}`;
}

function getValueSafely(obj, path, defaultVal) {
    try {
        let acc = obj;
        for (const x of path.split('.')) {
            acc = acc[x];
        }
        return acc;
    } catch (e) {
        return defaultVal;
    }
}

function trimTrailingSlash(url) {
    return url.replace(/\/+$/, '');
}

function isEffectiveChannelAdmin(userRoles) {
    const userRoleMap = {};
    userRoles.forEach((role) => {
        userRoleMap[role] = true;
    });

    return Boolean(userRoleMap[channelAdminRole] || userRoleMap[teamAdminRole] || userRoleMap[systemAdminRole]);
}

/**
 * Get current user's system and channel level roles.
 *
 * @param state the global application state.
 */
function getCurrentUserRoles(state) {
    const userRoles = [];

    // check for system admin
    if (state.entities.users.profiles[state.entities.users.currentUserId].roles.indexOf(systemAdminRole) > -1) {
        userRoles.push(systemAdminRole);
    }

    // check channel-level roles
    const currentChannelID = state.entities.channels.currentChannelId;
    const currentChannel = state.entities.channels.myMembers[currentChannelID];
    if (currentChannel !== undefined) {
        // this is undefined at Mattermost initialization
        userRoles.push(...(currentChannel.roles.split(' ')));
    }

    return userRoles;
}

function isGuestUser(state, userID) {
    return state.entities.users.profiles[userID] ? state.entities.users.profiles[userID].roles.indexOf(systemGuestRole) > -1 : false;
}

export default {
    getBaseURL,
    getValueSafely,
    trimTrailingSlash,
    isEffectiveChannelAdmin,
    channelAdminRole,
    teamAdminRole,
    systemAdminRole,
    getCurrentUserRoles,
    isGuestUser,
};
