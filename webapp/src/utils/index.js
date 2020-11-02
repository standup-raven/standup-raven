const channelAdminRole = 'channel_admin';
const teamAdminRole = 'team_admin';
const systemAdminRole = 'system_admin';

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

    return userRoleMap[channelAdminRole] || userRoleMap[teamAdminRole] || userRoleMap[systemAdminRole];
}

export default {
    getBaseURL,
    getValueSafely,
    trimTrailingSlash,
    isEffectiveChannelAdmin,
};
