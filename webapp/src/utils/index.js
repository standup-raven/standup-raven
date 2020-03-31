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

export default {
    getBaseURL,
    getValueSafely,
};
