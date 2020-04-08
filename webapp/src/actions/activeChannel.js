import Constants from '../constants';

export const addActiveChannel = (channelID) => (dispatch) => {
    console.log('addActiveChannel: ' + channelID);
    dispatch({
        type: Constants.ACTIONS.ADD_ACTIVE_CHANNEL,
        channelID,
    });
};

export const removeActiveChannel = (channelID) => (dispatch) => {
    console.log('removeActiveChannel: ' + channelID);
    dispatch({
        type: Constants.ACTIONS.REMOVE_ACTIVE_CHANNEL,
        channelID,
    });
};
