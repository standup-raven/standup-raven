import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import Actions from '../../actions';
import Selectors from '../../selectors';

import ConfigModal from './configModal';
import util from '../../utils';

const mapStateToProps = (state) => {
    return {
        currentUserId: state.entities.users.currentUserId,
        userRoles: util.getCurrentUserRoles(state),
        channelID: state.entities.channels.currentChannelId,
        visible: Selectors.isConfigModalVisible(state),
        siteURL: util.trimTrailingSlash(util.getValueSafely(state, 'entities.general.config.SiteURL', '')),
        isGuest: util.isGuestUser(state, state.entities.users.currentUserId),
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: Actions.closeConfigModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(ConfigModal);
