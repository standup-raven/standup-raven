import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import Actions from '../../actions';
import Selectors from '../../selectors';

import ConfigModal from './configModal';
import util from '../../utils';

const mapStateToProps = (state) => ({
    currentUserId: state.entities.users.currentUserId,
    userRoles: Object.keys(state.entities.roles.roles),
    channelID: state.entities.channels.currentChannelId,
    visible: Selectors.isConfigModalVisible(state),
    siteURL: util.trimTrailingSlash(util.getValueSafely(state, 'entities.general.config.SiteURL', '')),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: Actions.closeConfigModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(ConfigModal);
