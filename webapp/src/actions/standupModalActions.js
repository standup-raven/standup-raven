import Constants from '../constants';

export const openStandupModal = (channelId) => (dispatch) => {
    dispatch({
        type: Constants.ACTIONS.OPEN_STANDUP_MODAL,
        channelId,
    });
};

export const closeStandupModal = (channelId) => (dispatch) => {
    dispatch({
        type: Constants.ACTIONS.CLOSE_STANDUP_MODAL,
        channelId,
    });
};
