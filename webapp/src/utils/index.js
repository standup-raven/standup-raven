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
    console.log("##########################################################################################");
    console.log(obj);
    console.log(path);
    try {
        let acc = obj;
        for (let x of path.split('.')) {
            acc = acc[x];
            
            console.log(acc);
        }
        console.log("##########################################################################################");
        return acc;
    } catch (e) {
        console.log("##########################################################################################");
        return defaultVal;
    }
}

export default {
    getBaseURL,
    getValueSafely,
};
