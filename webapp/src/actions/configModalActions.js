import Constants from '../constants';

export const openConfigModal = (channelId) => (dispatch) => {
    dispatch({
        type: Constants.ACTIONS.OPEN_CONFIG_MODAL,
        channelId,
    });
};

export const closeConfigModal = (channelId) => (dispatch) => {
    dispatch({
        type: Constants.ACTIONS.CLOSE_CONFIG_MODAL,
        channelId,
    });
};
