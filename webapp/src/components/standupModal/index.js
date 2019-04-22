import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import Actions from '../../actions';
import Selectors from '../../selectors';

import Modal from './standupModal';

const mapStateToProps = (state) => ({
    currentUserId: state.entities.users.currentUserId,
    channelID: state.entities.channels.currentChannelId,
    visible: Selectors.isStandupModalVisible(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: Actions.closeStandupModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Modal);
