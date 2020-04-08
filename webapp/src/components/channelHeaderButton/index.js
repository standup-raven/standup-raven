import {connect} from 'react-redux';

import ChannelHeaderButtonIcon from './button';
import util from '../../utils';
import Selectors from '../../selectors';

const mapStateToProps = (state) => ({
    channelID: state.entities.channels.currentChannelId,
    siteURL: util.trimTrailingSlash(util.getValueSafely(state, 'entities.general.config.SiteURL', '')),
    added: Selectors.addedActiveChannel(state),
    removed: Selectors.removedActiveChannel(state),
});

export default connect(mapStateToProps)(ChannelHeaderButtonIcon);
