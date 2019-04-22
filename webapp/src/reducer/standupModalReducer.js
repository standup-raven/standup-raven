import Constants from '../constants';

let prevState;

export const standupModalVisible = (state = false, action) => {
    switch (action.type) {
        case Constants.ACTIONS.OPEN_STANDUP_MODAL:
            prevState = true;
            return true;
        case Constants.ACTIONS.CLOSE_STANDUP_MODAL:
            prevState = false;
            return false;
        default:
            return prevState || false;
    }
};
