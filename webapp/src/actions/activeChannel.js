import Constants from '../constants';

export const addActiveChannel = (channelID) => (dispatch) => {
    dispatch({
        type: Constants.ACTIONS.ADD_ACTIVE_CHANNEL,
        channelID,
    });
};

export const removeActiveChannel = (channelID) => (dispatch) => {
    dispatch({
        type: Constants.ACTIONS.REMOVE_ACTIVE_CHANNEL,
        channelID,
    });
};
