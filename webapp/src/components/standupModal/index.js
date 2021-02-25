import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import Actions from '../../actions';
import Selectors from '../../selectors';

import Modal from './standupModal';
import util from '../../utils';

const mapStateToProps = (state) => ({
    currentUserId: state.entities.users.currentUserId,
    channelID: state.entities.channels.currentChannelId,
    visible: Selectors.isStandupModalVisible(state),
    siteURL: util.trimTrailingSlash(util.getValueSafely(state, 'entities.general.config.SiteURL', '')),
    isGuest: util.isGuestUser(state, state.entities.users.currentUserId),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: Actions.closeStandupModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Modal);
