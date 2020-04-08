import request from 'superagent';
import Constants from '../constants';

async function getActiveChannels(siteURL) {
    const response = await request
        .get(`${siteURL}/${Constants.URL_ACTIVE_CHANNELS}`)
        .withCredentials();
    return response.body;
}

module.exports = {
    getActiveChannels,
};
