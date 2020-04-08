import Constants from '../constants';

let addedActiveChannelPrevState;
let removedActiveChannelPrevState;

export const addedActiveChannel = (channelID, action) => {
    switch (action.type) {
        case Constants.ACTIONS.ADD_ACTIVE_CHANNEL:
            addedActiveChannelPrevState = action.channelID;
            return action.channelID;
        default:
            return addedActiveChannelPrevState || '';
    }
};

export const removedActiveChannel = (channelID, action) => {
    switch (action.type) {
        case Constants.ACTIONS.REMOVE_ACTIVE_CHANNEL:
            removedActiveChannelPrevState = action.channelID;
            return action.channelID;
        default:
            return removedActiveChannelPrevState || '';
    }
};
