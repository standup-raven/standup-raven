import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import Actions from '../../actions';
import Selectors from '../../selectors';

import ConfigModal from './configModal';

const mapStateToProps = (state) => ({
    currentUserId: state.entities.users.currentUserId,
    channelID: state.entities.channels.currentChannelId,
    visible: Selectors.isConfigModalVisible(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: Actions.closeConfigModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(ConfigModal);
