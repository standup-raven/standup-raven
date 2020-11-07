import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import Actions from '../../actions';
import Selectors from '../../selectors';

import ConfigModal from './configModal';
import util from '../../utils';

const mapStateToProps = (state) => {
    return {
        currentUserId: state.entities.users.currentUserId,
        userRoles: getCurrentUserRoles(state),
        channelID: state.entities.channels.currentChannelId,
        visible: Selectors.isConfigModalVisible(state),
        siteURL: util.trimTrailingSlash(util.getValueSafely(state, 'entities.general.config.SiteURL', '')),
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: Actions.closeConfigModal,
}, dispatch);

/**
 * Get current user's system and channel level roles.
 *
 * @param state the global application state.
 */
function getCurrentUserRoles(state) {
    const userRoles = [];

    // check for system admin
    if (state.entities.roles.roles[util.systemAdminRole] !== undefined) {
        userRoles.push(util.systemAdminRole);
    }

    // check channel-level roles
    const currentChannelID = state.entities.channels.currentChannelId;
    const currentChannel = state.entities.channels.myMembers[currentChannelID];
    if (currentChannel !== undefined) {
        // this is undefined at Mattermost initialization
        userRoles.push(...(currentChannel.roles.split(' ')));
    }

    console.log(userRoles);
    return userRoles;
}

export default connect(mapStateToProps, mapDispatchToProps)(ConfigModal);
