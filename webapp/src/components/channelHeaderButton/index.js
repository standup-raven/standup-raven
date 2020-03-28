import {connect} from 'react-redux';

import ChannelHeaderButtonIcon from './button';

const mapStateToProps = (state) => ({
    channelID: state.entities.channels.currentChannelId,
});

export default connect(mapStateToProps)(ChannelHeaderButtonIcon);
