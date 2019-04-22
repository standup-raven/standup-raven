import {combineReducers} from 'redux';
import {standupModalVisible} from './standupModalReducer';
import {configModalVisible} from './configModalReducer';

export default combineReducers({
    standupModalVisible,
    configModalVisible,
});
