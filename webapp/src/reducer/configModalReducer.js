import Constants from '../constants';

let prevState;

export const configModalVisible = (state = false, action) => {
    switch (action.type) {
        case Constants.ACTIONS.OPEN_CONFIG_MODAL:
            prevState = true;
            return true;
        case Constants.ACTIONS.CLOSE_CONFIG_MODAL:
            prevState = false;
            return false;
        default:
            return prevState || false;
    }
};
