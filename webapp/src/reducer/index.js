import {combineReducers} from 'redux';
import {standupModalVisible} from './standupModalReducer';
import {configModalVisible} from './configModalReducer';
import {addedActiveChannel, removedActiveChannel} from './activeChannel';

export default combineReducers({
    standupModalVisible,
    configModalVisible,
    addedActiveChannel,
    removedActiveChannel,
});
